package infura

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/google/uuid"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	"github.com/okx/okbchain/x/infura/distrlock"
	"github.com/okx/okbchain/x/infura/types"
)

const (
	latestTaskKey         = "infura_latest_task"
	distributeLock        = "infura_lock"
	distributeLockTimeout = 300000
	taskTimeout           = distributeLockTimeout * 0.98

	FlagEnable         = "infura.enable"
	FlagRedisUrl       = "infura.redis-url"
	FlagRedisAuth      = "infura.redis-auth"
	FlagRedisDB        = "infura.redis-db"
	FlagMysqlUrl       = "infura.mysql-url"
	FlagMysqlUser      = "infura.mysql-user"
	FlagMysqlPass      = "infura.mysql-pass"
	FlagMysqlDB        = "infura.mysql-db"
	FlagCacheQueueSize = "infura.cache-queue-size"
)

// Stream maintains the infura engine
type Stream struct {
	enable     bool
	logger     log.Logger
	cfg        *types.Config
	cache      *Cache
	engine     types.IStreamEngine
	scheduler  types.IDistributeStateService
	cacheQueue *CacheQueue
}

func NewStream(logger log.Logger) *Stream {
	logger.Info("entering NewStream")
	se := &Stream{
		enable: viper.GetBool(FlagEnable),
		logger: logger,
	}
	if !se.enable {
		return se
	}
	// initialize
	se.cache = NewCache()
	se.cfg = initConfig()
	engine, err := newStreamEngine(se.cfg, logger)
	if err != nil {
		panic(fmt.Sprintf("ParseStreamEngineConfig failed: %+v", err))
	}
	se.engine = engine

	scheduler, err := newRedisLockService(se.cfg.RedisUrl, se.cfg.RedisAuth, se.cfg.RedisDB, logger)
	if err != nil {
		errStr := fmt.Sprintf("parse redis scheduler failed error: %s, redis url: %s", err.Error(), se.cfg.RedisUrl)
		logger.Error(errStr)
		panic(errStr)
	}
	se.scheduler = scheduler

	// start cache queue
	if se.cfg.CacheQueueSize > 0 {
		se.cacheQueue = newCacheQueue(se.cfg.CacheQueueSize)
		go se.cacheQueue.Start()
	}

	se.logger.Info("NewStream success.")
	return se
}

func initConfig() *types.Config {
	return &types.Config{
		RedisUrl:       viper.GetString(FlagRedisUrl),
		RedisAuth:      viper.GetString(FlagRedisAuth),
		RedisDB:        viper.GetInt(FlagRedisDB),
		MysqlUrl:       viper.GetString(FlagMysqlUrl),
		MysqlUser:      viper.GetString(FlagMysqlUser),
		MysqlPass:      viper.GetString(FlagMysqlPass),
		MysqlDB:        viper.GetString(FlagMysqlDB),
		CacheQueueSize: viper.GetInt(FlagCacheQueueSize),
	}
}

func newRedisLockService(redisURL string, redisPass string, db int, logger log.Logger) (types.IDistributeStateService, error) {
	if redisURL == "" {
		return nil, fmt.Errorf("no valid redisUrl found, no IDistributeStateService is created, redisUrl: %s", redisURL)
	}
	workerID := uuid.New().String()

	scheduler, err := distrlock.NewRedisDistributeStateService(redisURL, redisPass, db, logger, workerID)
	return scheduler, err
}
