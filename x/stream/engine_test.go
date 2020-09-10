package stream

import (
	"os"
	"testing"

	appCfg "github.com/cosmos/cosmos-sdk/server/config"

	"github.com/okex/okchain/x/stream/common"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	REDISURL  = "redis://127.0.0.1:16379"
	PULSARURL = "127.0.0.1:6650"
	MYSQLURL  = "okdexer:okdex123!@tcp(127.0.0.1:13306)/okdex"
)

func TestParseStreamEngineConfig(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	// clear redis
	pool, err := common.NewPool(REDISURL, "", logger)
	require.Nil(t, err)
	_, err = pool.Get().Do("FLUSHALL")
	require.Nil(t, err)
	// clear viper
	viper.Reset()

	cfg := appCfg.DefaultStreamConfig()

	engines, err := ParseStreamEngineConfig(logger, cfg)
	require.NotNil(t, err)
	require.Nil(t, engines)

	cfg.Engine = "analysismysql|" + MYSQLURL
	engines, err = ParseStreamEngineConfig(logger, cfg)
	require.NotNil(t, err)
	require.Nil(t, engines)

	cfg.Engine = "analysis|xxxx|" + MYSQLURL
	engines, err = ParseStreamEngineConfig(logger, cfg)
	require.NotNil(t, err)
	require.Nil(t, engines)

	cfg.Engine = "notify|redis|" + "127.0.0.1:99999"
	engines, err = ParseStreamEngineConfig(logger, cfg)
	require.NotNil(t, err)
	require.Nil(t, engines)

	cfg.Engine = "notify|xxxx|" + REDISURL
	engines, err = ParseStreamEngineConfig(logger, cfg)
	require.NotNil(t, err)
	require.Nil(t, engines)

	cfg.Engine = "kline|pulsar|" + "127.0.0.1:99999"
	engines, err = ParseStreamEngineConfig(logger, cfg)
	require.NotNil(t, err)
	require.Nil(t, engines)

	cfg.Engine = "kline|xxxxx|" + PULSARURL
	engines, err = ParseStreamEngineConfig(logger, cfg)
	require.NotNil(t, err)
	require.Nil(t, engines)

	cfg.Engine = "analysis|mysql|" + MYSQLURL + ",notify|redis|" + REDISURL + ",kline|pulsar|" + PULSARURL
	engines, err = ParseStreamEngineConfig(logger, cfg)
	require.Nil(t, err)
	require.Equal(t, 3, len(engines))
}

func TestNewMySqlEngine(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	// clear redis
	pool, err := common.NewPool(REDISURL, "", logger)
	require.Nil(t, err)
	_, err = pool.Get().Do("FLUSHALL")
	require.Nil(t, err)
	// clear viper
	viper.Reset()

	engine, err := NewMySQLEngine(MYSQLURL, logger, nil)
	require.Nil(t, err)
	require.NotNil(t, engine)
}

func TestNewRedisEngine(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	// clear redis
	pool, err := common.NewPool(REDISURL, "", logger)
	require.Nil(t, err)
	_, err = pool.Get().Do("FLUSHALL")
	require.Nil(t, err)
	// clear viper
	viper.Reset()

	engine, err := NewRedisEngine("", logger, nil)
	require.NotNil(t, err)
	require.Nil(t, engine)

	engine, err = NewRedisEngine(REDISURL, logger, nil)
	require.Nil(t, err)
	require.NotNil(t, engine)
}

func TestNewPulsarEngine(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	// clear redis
	pool, err := common.NewPool(REDISURL, "", logger)
	require.Nil(t, err)
	_, err = pool.Get().Do("FLUSHALL")
	require.Nil(t, err)
	// clear viper
	viper.Reset()
	cfg := appCfg.DefaultStreamConfig()
	engine, err := NewPulsarEngine(PULSARURL, logger, cfg)
	require.Nil(t, err)
	require.NotNil(t, engine)
}

func TestStringToEngineKind(t *testing.T) {
	kind := "Analysis"
	require.Equal(t, EngineAnalysisKind, StringToEngineKind(kind))
	kind = "notify"
	require.Equal(t, EngineNotifyKind, StringToEngineKind(kind))
	kind = "kline"
	require.Equal(t, EngineKlineKind, StringToEngineKind(kind))
	kind = ""
	require.Equal(t, EngineNilKind, StringToEngineKind(kind))
}

func TestStringToStreamKind(t *testing.T) {
	kind := "Mysql"
	require.Equal(t, StreamMysqlKind, StringToStreamKind(kind))
	kind = "redis"
	require.Equal(t, StreamRedisKind, StringToStreamKind(kind))
	kind = "pulsar"
	require.Equal(t, StreamPulsarKind, StringToStreamKind(kind))
	kind = ""
	require.Equal(t, StreamNilKind, StringToStreamKind(kind))
}
