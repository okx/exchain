package stream

import (
	"fmt"
	"strings"
	"time"

	"github.com/okex/okchain/x/stream/quoteslite"

	"github.com/pkg/errors"

	"github.com/go-sql-driver/mysql"

	"github.com/tendermint/tendermint/libs/log"

	appCfg "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/okex/okchain/x/backend"
	"github.com/okex/okchain/x/stream/analyservice"
	"github.com/okex/okchain/x/stream/common"
	"github.com/okex/okchain/x/stream/pulsarclient"
	"github.com/okex/okchain/x/stream/pushservice"
	pushservicetypes "github.com/okex/okchain/x/stream/pushservice/types"
	"github.com/okex/okchain/x/stream/types"
)

type StreamKind byte
type EngineKind byte

const (
	StreamNilKind       StreamKind = 0x00
	StreamMysqlKind     StreamKind = 0x01
	StreamRedisKind     StreamKind = 0x02
	StreamPulsarKind    StreamKind = 0x03
	StreamWebSocketKind StreamKind = 0x04

	EngineNilKind       EngineKind = 0x00
	EngineAnalysisKind  EngineKind = 0x01
	EngineNotifyKind    EngineKind = 0x02
	EngineKlineKind     EngineKind = 0x03
	EngineWebSocketKind EngineKind = 0x04
)

var StreamKind2EngineKindMap = map[StreamKind]EngineKind{
	StreamMysqlKind:     EngineAnalysisKind,
	StreamRedisKind:     EngineNotifyKind,
	StreamPulsarKind:    EngineKlineKind,
	StreamWebSocketKind: EngineWebSocketKind,
}

var EngineKind2StreamKindMap = map[EngineKind]StreamKind{
	EngineAnalysisKind:  StreamMysqlKind,
	EngineNotifyKind:    StreamRedisKind,
	EngineKlineKind:     StreamPulsarKind,
	EngineWebSocketKind: StreamWebSocketKind,
}

// ***********************************
type MySqlEngine struct {
	url    string
	logger log.Logger
	orm    *backend.ORM
}

func NewMySqlEngine(url string, log log.Logger, cfg *appCfg.StreamConfig) (types.IStreamEngine, error) {
	ormTmp := analyservice.NewMysqlORM(url)
	log.Info("NewAnalysisService succeed")
	// connect mysql through streamUrl
	return &MySqlEngine{
		url:    url,
		logger: log,
		orm:    ormTmp,
	}, nil
}

func (e *MySqlEngine) Url() string {
	return e.url
}

func (e *MySqlEngine) Write(data types.IStreamData, success *bool) {
	e.logger.Debug("Entering MySqlEngine write")
	enData, ok := data.(*analyservice.DataAnalysis)
	if !ok {
		panic(fmt.Sprintf("MySqlEngine Convert data %+v to DataAnalysis failed", data))
	}

	results, err := e.orm.BatchInsertOrUpdate(enData.NewOrders, enData.UpdatedOrders, enData.Deals, enData.MatchResults, enData.FeeDetails, enData.Trans)
	if err != nil {
		e.logger.Error(fmt.Sprintf("MySqlEngine write failed: %s, results: %v", err.Error(), results))
		*success = false

		if mysqlerr, ok := err.(*mysql.MySQLError); ok {
			e.logger.Error(fmt.Sprintf("MySQLError: %+v", err.Error()))
			// Duplicate entry 'XXX' for key 'PRIMARY'
			if mysqlerr.Number == 1062 {
				e.logger.Error(fmt.Sprintf("MySqlEngine write failed becoz 1062: %s, considered success, result: %+v", err.Error(), results))
				*success = true
			}
		}
	} else {
		e.logger.Debug(fmt.Sprintf("MySqlEngine write result: %+v", results))
		*success = true
	}

}

// ***********************************
type PulsarEngine struct {
	url            string
	logger         log.Logger
	pulsarProducer *pulsarclient.PulsarProducer
}

func NewPulsarEngine(url string, logger log.Logger, cfg *appCfg.StreamConfig) (types.IStreamEngine, error) {
	asyncErrs := make(chan error, 16)
	producer := pulsarclient.NewPulsarProducer(url, cfg, logger, &asyncErrs)

	time.Sleep(time.Second * 3)
	if len(asyncErrs) != 0 {
		err := <-asyncErrs
		logger.Error(fmt.Sprintf("create pulsar producer failed: %s", err.Error()))
		return nil, err
	}

	logger.Info("create pulsar producer succeed")

	// connect mysql through streamUrl()
	return &PulsarEngine{
		url:            url,
		logger:         logger,
		pulsarProducer: producer,
	}, nil
}

func (e *PulsarEngine) Url() string {
	return e.url
}

func (e *PulsarEngine) Write(data types.IStreamData, success *bool) {
	e.logger.Debug("Entering PulsarEngine Write")
	enData, ok := data.(*pulsarclient.PulsarData)
	if !ok {
		panic(fmt.Sprintf("Convert data %+v to PulsarData failed", data))
	}

	err := e.pulsarProducer.RefreshMarketIdMap(enData, e.logger)
	if err != nil {
		e.logger.Error(fmt.Sprintf("pulsar engine RefreshMarketIdMap failed: %s", err.Error()))
		*success = false
		return
	}

	results, err := e.pulsarProducer.SendAllMsg(enData, e.logger)
	if err != nil {
		e.logger.Error(fmt.Sprintf("pulsar engine write failed: %s, results: %v", err.Error(), results))
		*success = false
	} else {
		e.logger.Debug(fmt.Sprintf("PulsarEngine write result: %+v", results))
		*success = true
	}
}

// ***********************************
type RedisEngine struct {
	url    string
	logger log.Logger
	srv    *pushservice.PushService
}

func NewRedisEngine(url string, logger log.Logger, cfg *appCfg.StreamConfig) (types.IStreamEngine, error) {
	redisUrl, redisPassword, err := common.ParseRedisUrl(url, cfg.RedisRequirePass)
	if err != nil {
		return nil, err
	}
	srv, err := pushservice.NewPushService(redisUrl, redisPassword, 0, logger)
	if err != nil {
		logger.Error("NewPushService failed ", err.Error())
		return nil, err
	}
	logger.Info("NewPushService succeed")
	return &RedisEngine{
		url:    url,
		logger: logger,
		srv:    srv,
	}, nil
}

func (e *RedisEngine) Url() string {
	return e.url
}

func (e *RedisEngine) Write(data types.IStreamData, success *bool) {
	e.logger.Debug("Entering RedisEngine Write")
	enData, ok := data.(*pushservicetypes.RedisBlock)
	if !ok {
		panic(fmt.Sprintf("Convert data %+v to PulsarData failed", data))
	}

	results, err := e.srv.WriteSync(enData)

	if err != nil {
		e.logger.Error(fmt.Sprintf("redis engine write failed: %s, results: %+v", err.Error(), results))
		*success = false
	} else {
		e.logger.Debug(fmt.Sprintf("RedisEngine write result: %+v", results))
		*success = true
	}
}

type EngineCreator func(url string, logger log.Logger, cfg *appCfg.StreamConfig) (types.IStreamEngine, error)

func GetEngineCreator(eKind EngineKind, sKind StreamKind) (EngineCreator, error) {
	m := map[string]EngineCreator{
		fmt.Sprintf("%d_%d", EngineAnalysisKind, StreamMysqlKind):      NewMySqlEngine,
		fmt.Sprintf("%d_%d", EngineNotifyKind, StreamRedisKind):        NewRedisEngine,
		fmt.Sprintf("%d_%d", EngineKlineKind, StreamPulsarKind):        NewPulsarEngine,
		fmt.Sprintf("%d_%d", EngineWebSocketKind, StreamWebSocketKind): quoteslite.NewWebSocketEngine,
	}

	key := fmt.Sprintf("%d_%d", eKind, sKind)
	c, ok := m[key]
	if ok {
		return c, nil
	} else {
		return nil, fmt.Errorf("No EngineCreator found for EngineKine %d & StreamKine %d ", eKind, sKind)
	}
}

// ***********************************
func ParseStreamEngineConfig(logger log.Logger, cfg *appCfg.StreamConfig) (map[EngineKind]types.IStreamEngine, error) {
	if cfg.Engine == "" {
		return nil, errors.New("stream engine config is empty")
	}
	engines := make(map[EngineKind]types.IStreamEngine)
	list := strings.Split(cfg.Engine, ",")
	for _, item := range list {
		enginesConf := strings.Split(item, "|")

		// Desktop Stream Engine Mode: mysql | websocket
		// HA Stream Engine Mode: mysql | redis | pulsar(kafka)

		if len(enginesConf) != 3 {
			return nil, fmt.Errorf("expected list in a form of \"engine_type:stream_type:stream_url\" pairs, given pair %s, list %s", item, list)
		}

		engineType := StringToEngineKind(enginesConf[0])
		streamType := StringToStreamKind(enginesConf[1])
		streamUrl := enginesConf[2]

		creatorFunc, err := GetEngineCreator(engineType, streamType)
		if err != nil {
			return nil, err
		} else {
			engine, err := creatorFunc(streamUrl, logger, cfg)
			if err != nil {
				return nil, err
			}
			engines[engineType] = engine

		}
	}

	return engines, nil
}

func StringToEngineKind(kind string) EngineKind {
	kind = strings.ToLower(kind)
	switch kind {
	case "analysis":
		return EngineAnalysisKind
	case "notify":
		return EngineNotifyKind
	case "kline":
		return EngineKlineKind
	case "websocket":
		return EngineWebSocketKind
	default:
		return EngineNilKind
	}
}

func StringToStreamKind(kind string) StreamKind {
	kind = strings.ToLower(kind)
	switch kind {
	case "mysql":
		return StreamMysqlKind
	case "redis":
		return StreamRedisKind
	case "pulsar":
		return StreamPulsarKind
	case "websocket":
		return StreamWebSocketKind
	default:
		return StreamNilKind
	}
}
