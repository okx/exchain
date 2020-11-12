package stream

import (
	"fmt"
	"github.com/okex/okexchain/x/stream/kafkaclient"
	"strings"
	"time"

	"github.com/okex/okexchain/x/stream/websocket"

	"github.com/pkg/errors"

	"github.com/go-sql-driver/mysql"

	"github.com/tendermint/tendermint/libs/log"

	appCfg "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/okex/okexchain/x/backend"
	"github.com/okex/okexchain/x/stream/analyservice"
	"github.com/okex/okexchain/x/stream/common"
	"github.com/okex/okexchain/x/stream/pulsarclient"
	"github.com/okex/okexchain/x/stream/pushservice"
	pushservicetypes "github.com/okex/okexchain/x/stream/pushservice/types"
	"github.com/okex/okexchain/x/stream/types"
)

type Kind byte
type EngineKind byte

const (
	StreamNilKind       Kind = 0x00
	StreamMysqlKind     Kind = 0x01
	StreamRedisKind     Kind = 0x02
	StreamPulsarKind    Kind = 0x03
	StreamWebSocketKind Kind = 0x04
	StreamKafkaKind     Kind = 0x05

	EngineNilKind       EngineKind = 0x00
	EngineAnalysisKind  EngineKind = 0x01
	EngineNotifyKind    EngineKind = 0x02
	EngineKlineKind     EngineKind = 0x03
	EngineWebSocketKind EngineKind = 0x04
)

var StreamKind2EngineKindMap = map[Kind]EngineKind{
	StreamMysqlKind:     EngineAnalysisKind,
	StreamRedisKind:     EngineNotifyKind,
	StreamPulsarKind:    EngineKlineKind,
	StreamKafkaKind:     EngineKlineKind,
	StreamWebSocketKind: EngineWebSocketKind,
}

var EngineKind2StreamKindMap = map[EngineKind]Kind{
	EngineAnalysisKind:  StreamMysqlKind,
	EngineNotifyKind:    StreamRedisKind,
	EngineWebSocketKind: StreamWebSocketKind,
}

type MySQLEngine struct {
	url    string
	logger log.Logger
	orm    *backend.ORM
}

func NewMySQLEngine(url string, log log.Logger, cfg *appCfg.StreamConfig) (types.IStreamEngine, error) {
	ormTmp := analyservice.NewMysqlORM(url)
	log.Info("NewAnalysisService succeed")
	// connect mysql through streamUrl
	return &MySQLEngine{
		url:    url,
		logger: log,
		orm:    ormTmp,
	}, nil
}

func (e *MySQLEngine) URL() string {
	return e.url
}

func (e *MySQLEngine) Write(data types.IStreamData, success *bool) {
	e.logger.Debug("Entering MySqlEngine write")
	enData, ok := data.(*analyservice.DataAnalysis)
	if !ok {
		panic(fmt.Sprintf("MySqlEngine Convert data %+v to DataAnalysis failed", data))
	}

	results, err := e.orm.BatchInsertOrUpdate(enData.NewOrders, enData.UpdatedOrders, enData.Deals, enData.MatchResults,
		enData.FeeDetails, enData.Trans, enData.SwapInfos)
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

func (e *PulsarEngine) URL() string {
	return e.url
}

func (e *PulsarEngine) Write(data types.IStreamData, success *bool) {
	e.logger.Debug("Entering PulsarEngine Write")
	enData, ok := data.(*common.KlineData)
	if !ok {
		panic(fmt.Sprintf("Convert data %+v to KlineData failed", data))
	}

	err := e.pulsarProducer.RefreshMarketIDMap(enData, e.logger)
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

type KafkaEngine struct {
	url           string
	logger        log.Logger
	kafkaProducer *kafkaclient.KafkaProducer
}

func NewKafkaEngine(url string, logger log.Logger, cfg *appCfg.StreamConfig) (types.IStreamEngine, error) {
	return &KafkaEngine{
		url:           url,
		logger:        logger,
		kafkaProducer: kafkaclient.NewKafkaProducer(url, cfg),
	}, nil
}

func (ke *KafkaEngine) URL() string {
	return ke.url
}

func (ke *KafkaEngine) Write(data types.IStreamData, success *bool) {
	ke.logger.Debug("Entering KafkaEngine Write")
	enData, ok := data.(*common.KlineData)
	if !ok {
		panic(fmt.Sprintf("Convert data %+v to KlineData failed", data))
	}

	if err := ke.kafkaProducer.RefreshMarketIDMap(enData, ke.logger); err != nil {
		ke.logger.Error(fmt.Sprintf("kafka engine RefreshMarketIdMap failed: %s", err.Error()))
		*success = false
		return
	}

	results, err := ke.kafkaProducer.SendAllMsg(enData, ke.logger)
	if err != nil {
		ke.logger.Error(fmt.Sprintf("kafka engine write failed: %s, results: %v", err.Error(), results))
		*success = false
	} else {
		ke.logger.Debug(fmt.Sprintf("kafka engine write result: %+v", results))
		*success = true
	}
}

type RedisEngine struct {
	url    string
	logger log.Logger
	srv    *pushservice.PushService
}

func NewRedisEngine(url string, logger log.Logger, cfg *appCfg.StreamConfig) (types.IStreamEngine, error) {
	redisURL, redisPassword, err := common.ParseRedisURL(url, cfg.RedisRequirePass)
	if err != nil {
		return nil, err
	}
	srv, err := pushservice.NewPushService(redisURL, redisPassword, 0, logger)
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

func (e *RedisEngine) URL() string {
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

func GetEngineCreator(eKind EngineKind, sKind Kind) (EngineCreator, error) {
	m := map[string]EngineCreator{
		fmt.Sprintf("%d_%d", EngineAnalysisKind, StreamMysqlKind):      NewMySQLEngine,
		fmt.Sprintf("%d_%d", EngineNotifyKind, StreamRedisKind):        NewRedisEngine,
		fmt.Sprintf("%d_%d", EngineKlineKind, StreamPulsarKind):        NewPulsarEngine,
		fmt.Sprintf("%d_%d", EngineWebSocketKind, StreamWebSocketKind): websocket.NewEngine,
		fmt.Sprintf("%d_%d", EngineKlineKind, StreamKafkaKind):         NewKafkaEngine,
	}

	key := fmt.Sprintf("%d_%d", eKind, sKind)
	c, ok := m[key]
	if ok {
		return c, nil
	}
	return nil, fmt.Errorf("no EngineCreator found for EngineKine %d & StreamKine %d ", eKind, sKind)
}

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
		streamURL := enginesConf[2]

		creatorFunc, err := GetEngineCreator(engineType, streamType)
		if err != nil {
			return nil, err
		}

		engine, err := creatorFunc(streamURL, logger, cfg)
		if err != nil {
			return nil, err
		}
		engines[engineType] = engine
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

func StringToStreamKind(kind string) Kind {
	kind = strings.ToLower(kind)
	switch kind {
	case "mysql":
		return StreamMysqlKind
	case "redis":
		return StreamRedisKind
	case "pulsar":
		EngineKind2StreamKindMap[EngineKlineKind] = StreamPulsarKind
		return StreamPulsarKind
	case "websocket":
		return StreamWebSocketKind
	case "kafka":
		EngineKind2StreamKindMap[EngineKlineKind] = StreamKafkaKind
		return StreamKafkaKind
	default:
		return StreamNilKind
	}
}
