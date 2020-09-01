package stream

import (
	"errors"
	"fmt"

	"github.com/okex/okchain/x/stream/eureka"
	"github.com/okex/okchain/x/stream/nacos"
	"github.com/okex/okchain/x/stream/websocket"

	appCfg "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/google/uuid"
	"github.com/okex/okchain/x/stream/distrlock"

	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/okex/okchain/x/backend"
	"github.com/okex/okchain/x/stream/common"
	"github.com/okex/okchain/x/stream/pushservice"
	"github.com/okex/okchain/x/stream/types"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	latestTaskKey         = "latest_stream_task"
	distributeLock        = "stream_lock"
	distributeLockTimeout = 30000
)

// Stream maintains the engines
type Stream struct {
	orderKeeper    types.OrderKeeper    // The reference to the OrderKeeper to get deals
	tokenKeeper    types.TokenKeeper    // The reference to the TokenKeeper to get fee details
	accountKeeper  types.AccountKeeper  //for account
	marketKeeper   backend.MarketKeeper // The reference to MarketKeeper to get ticker/klines
	dexKeeper      types.DexKeeper
	cdc            *codec.Codec // The wire codec for binary encoding/decoding.
	logger         log.Logger
	engines        map[EngineKind]types.IStreamEngine
	Cache          *common.Cache
	AnalysisEnable bool

	// Fore. 20190809
	scheduler       types.IDistributeStateService
	distrLatestTask *Task
	taskChan        chan TaskWithData
	resultChan      chan Task
	coordinator     *Coordinator
	cacheQueue      *CacheQueue
	cfg             *appCfg.StreamConfig
}

func NewStream(orderKeeper types.OrderKeeper, tokenKeeper types.TokenKeeper, dexKeeper types.DexKeeper, cdc *codec.Codec, logger log.Logger, cfg *appCfg.Config) *Stream {

	logger.Info("entering NewStreamEngine")

	se := &Stream{
		orderKeeper: orderKeeper,
		tokenKeeper: tokenKeeper,
		dexKeeper:   dexKeeper,
		cdc:         cdc,
		logger:      logger,
		Cache:       common.NewCache(),
	}
	// read config
	se.cfg = cfg.StreamConfig
	logger.Debug("NewStream", "config", *se.cfg)

	// start eureka client
	if cfg.BackendConfig.EnableBackend == true && se.cfg.EurekaServerUrl != "" {
		eureka.StartEurekaClient(logger, se.cfg.EurekaServerUrl, se.cfg.RestApplicationName)
	}
	// start nacos client
	if cfg.BackendConfig.EnableBackend == true && se.cfg.NacosServerUrl != "" {
		nacos.StartNacosClient(logger, se.cfg.NacosServerUrl, se.cfg.NacosNamespaceId, se.cfg.RestApplicationName)
	}
	// Enable marketKeeper if KlineQueryConnect is set.
	if se.cfg.KlineQueryConnect != "" {
		address, password, err := common.ParseRedisUrl(se.cfg.KlineQueryConnect, se.cfg.RedisRequirePass)
		if err != nil {
			logger.Error("Fail to parse redis url ", se.cfg.KlineQueryConnect, " error: ", err.Error())
		} else {
			srv, err := pushservice.NewPushService(address, password, 0, logger)
			if err != nil {
				logger.Error("NewPushService failed ", err.Error())
			} else {
				se.marketKeeper = NewRedisMarketKeeper(srv.GetConnCli(), logger)
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

	se.taskChan = make(chan TaskWithData, 1)
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
		go websocket.StartWSServer(logger, se.engines[EngineWebSocketKind].Url())
	}

	return se
}

func newRedisLockServiceWithConf(redisUrl string, redisPass string, workerId string, logger log.Logger) (types.IDistributeStateService, error) {
	if redisUrl == "" {
		return nil, errors.New(fmt.Sprintf("no valid redisUrl found, no IDistributeStateService is created, redisUrl: %s", redisUrl))
	}
	if workerId == "" {
		workerId = uuid.New().String()
	} else {
		workerId = workerId + "-" + uuid.New().String()
	}

	scheduler, err := distrlock.NewRedisDistributeStateService(redisUrl, redisPass, logger, workerId)
	return scheduler, err
}
