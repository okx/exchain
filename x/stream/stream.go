package stream

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/spf13/viper"

	"github.com/okex/exchain/x/stream/eureka"
	"github.com/okex/exchain/x/stream/nacos"
	"github.com/okex/exchain/x/stream/websocket"

	"github.com/google/uuid"
	appCfg "github.com/okex/exchain/libs/cosmos-sdk/server/config"
	"github.com/okex/exchain/x/stream/distrlock"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/stream/common"
	"github.com/okex/exchain/x/stream/pushservice"
	"github.com/okex/exchain/x/stream/types"
)

const (
	latestTaskKey         = "latest_stream_task"
	distributeLock        = "stream_lock"
	distributeLockTimeout = 30000

	NacosTmrpcUrls        = "stream.tmrpc_nacos_urls"
	NacosTmrpcNamespaceID = "stream.tmrpc_nacos_namespace_id"
	NacosTmrpcAppName     = "stream.tmrpc_application_name"
	RpcExternalAddr       = "rpc.external_laddr"
)

// Stream maintains the engines
type Stream struct {
	orderKeeper    types.OrderKeeper // The reference to the OrderKeeper to get deals
	tokenKeeper    types.TokenKeeper // The reference to the TokenKeeper to get fee details
	dexKeeper      types.DexKeeper
	swapKeeper     types.SwapKeeper
	farmKeeper     types.FarmKeeper
	cdc            *codec.Codec // The wire codec for binary encoding/decoding.
	logger         log.Logger
	engines        map[EngineKind]types.IStreamEngine
	Cache          *common.Cache
	AnalysisEnable bool

	// Fore. 20190809
	scheduler       types.IDistributeStateService
	distrLatestTask *Task
	taskChan        chan *TaskWithData
	resultChan      chan Task
	coordinator     *Coordinator
	cacheQueue      *CacheQueue
	cfg             *appCfg.StreamConfig
}

func NewStream(orderKeeper types.OrderKeeper, tokenKeeper types.TokenKeeper, dexKeeper types.DexKeeper, swapKeeper types.SwapKeeper, farmKeeper types.FarmKeeper, cdc *codec.Codec, logger log.Logger, cfg *appCfg.Config) *Stream {

	logger.Info("entering NewStreamEngine")

	se := &Stream{
		orderKeeper: orderKeeper,
		tokenKeeper: tokenKeeper,
		dexKeeper:   dexKeeper,
		swapKeeper:  swapKeeper,
		farmKeeper:  farmKeeper,
		cdc:         cdc,
		logger:      logger,
		Cache:       common.NewCache(),
	}
	// read config
	se.cfg = cfg.StreamConfig
	logger.Debug("NewStream", "config", *se.cfg)

	// start eureka client for registering restful service
	if cfg.BackendConfig.EnableBackend && se.cfg.EurekaServerUrl != "" {
		eureka.StartEurekaClient(logger, se.cfg.EurekaServerUrl, se.cfg.RestApplicationName, se.RestExternalAddr())
	}

	// start nacos client for registering restful service
	if cfg.BackendConfig.EnableBackend && se.cfg.RestNacosUrls != "" {
		nacos.StartNacosClient(logger, se.cfg.RestNacosUrls, se.cfg.RestNacosNamespaceId, se.cfg.RestApplicationName, se.RestExternalAddr())
	}

	// start nacos client for tmrpc service
	if se.NacosTmRpcUrls() != "" {
		nacos.StartNacosClient(logger, se.NacosTmRpcUrls(), se.NacosTmRpcNamespaceID(), se.NacosTmRpcAppName(), se.RpcExternalAddr())
	}

	// Enable marketKeeper if KlineQueryConnect is set.
	if se.cfg.KlineQueryConnect != "" {
		address, password, err := common.ParseRedisURL(se.cfg.KlineQueryConnect, se.cfg.RedisRequirePass)
		if err != nil {
			logger.Error("Fail to parse redis url ", se.cfg.KlineQueryConnect, " error: ", err.Error())
		} else {
			_, err := pushservice.NewPushService(address, password, 0, logger)
			if err != nil {
				logger.Error("NewPushService failed ", err.Error())
			} else {
				logger.Info("NewPushService succeed")
			}
		}
	}

	//
	if se.cfg.Engine == "" {
		return se
	}

	// LocalLockService is used for desktop environment
	var scheduler types.IDistributeStateService
	var err error
	if se.cfg.RedisScheduler != "" {
		scheduler, err = newRedisLockServiceWithConf(se.cfg.RedisScheduler, se.cfg.RedisRequirePass, se.cfg.WorkerId, logger)
	} else {
		scheduler, err = distrlock.NewLocalStateService(se.logger, se.cfg.WorkerId, se.cfg.LocalLockDir)
	}
	if err != nil {
		errStr := fmt.Sprintf("parse redis scheduler failed : %s", err.Error())
		logger.Error(errStr)
		panic(errStr)
	}
	se.scheduler = scheduler

	engines, err := ParseStreamEngineConfig(logger, se.cfg)
	if err != nil {
		errStr := fmt.Sprintf("ParseStreamEngineConfig failed: %+v", err)
		logger.Error(errStr)
		panic(errStr)
	}

	se.engines = engines
	se.logger.Info(fmt.Sprintf("%d engines created, verbose info: %+v", len(se.engines), se.engines))
	se.AnalysisEnable = se.engines[EngineAnalysisKind] != nil

	se.taskChan = make(chan *TaskWithData, 1)
	se.resultChan = make(chan Task, 1)
	se.distrLatestTask = nil
	se.coordinator = NewCoordinator(logger, se.taskChan, se.resultChan, distributeLockTimeout, se.engines)
	go se.coordinator.run()

	// start stream cache queue
	if se.cfg.CacheQueueCapacity > 0 {
		se.cacheQueue = newCacheQueue(se.cfg.CacheQueueCapacity)
		go se.cacheQueue.Start()
	}

	se.logger.Info("NewStreamEngine success.")

	// Enable websocket
	if se.engines[EngineWebSocketKind] != nil {
		go websocket.StartWSServer(logger, se.engines[EngineWebSocketKind].URL())
	}

	return se
}

func newRedisLockServiceWithConf(redisURL string, redisPass string, workerID string, logger log.Logger) (types.IDistributeStateService, error) {
	if redisURL == "" {
		return nil, fmt.Errorf("no valid redisUrl found, no IDistributeStateService is created, redisUrl: %s", redisURL)
	}
	if workerID == "" {
		workerID = uuid.New().String()
	} else {
		workerID = workerID + "-" + uuid.New().String()
	}

	scheduler, err := distrlock.NewRedisDistributeStateService(redisURL, redisPass, logger, workerID)
	return scheduler, err
}

func (s Stream) NacosTmRpcUrls() string {
	return viper.GetString(NacosTmrpcUrls)
}

func (s Stream) NacosTmRpcNamespaceID() string {
	return viper.GetString(NacosTmrpcNamespaceID)
}

func (s Stream) NacosTmRpcAppName() string {
	return viper.GetString(NacosTmrpcAppName)
}

func (s Stream) RpcExternalAddr() string {
	return viper.GetString(RpcExternalAddr)
}

func (s Stream) RestExternalAddr() string {
	return viper.GetString(server.FlagExternalListenAddr)
}
