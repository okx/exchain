package pulsarclient

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Comcast/pulsar-client-go"
	"github.com/google/uuid"
	"github.com/nacos-group/nacos-sdk-go/vo"
	appCfg "github.com/okex/exchain/libs/cosmos-sdk/server/config"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/stream/common/kline"
)

type PulsarProducer struct {
	kline.MarketConfig
	producers []*pulsar.ManagedProducer
	partion   int64
}

func NewPulsarProducer(url string, cfg *appCfg.StreamConfig, logger log.Logger, asyncErrs *chan error) *PulsarProducer {
	var mp = &PulsarProducer{
		MarketConfig: kline.NewMarketConfig(cfg.MarketServiceEnable, cfg.MarketNacosUrls, cfg.MarketNacosNamespaceId,
			cfg.MarketNacosClusters, cfg.MarketNacosServiceName, cfg.MarketNacosGroupName),
		producers: make([]*pulsar.ManagedProducer, 0, cfg.MarketPartition),
		partion:   int64(cfg.MarketPartition),
	}

	for i := 0; i < cfg.MarketPartition; i++ {
		mcp := pulsar.NewManagedClientPool()
		mpCfg := pulsar.ManagedProducerConfig{
			Name:                  uuid.New().String() + "-subs_standard_dex_spot-" + strconv.Itoa(i),
			Topic:                 cfg.MarketTopic + "-partition-" + strconv.Itoa(i),
			NewProducerTimeout:    time.Second * 3,
			InitialReconnectDelay: time.Second,
			MaxReconnectDelay:     time.Minute,
			ManagedClientConfig: pulsar.ManagedClientConfig{
				ClientConfig: pulsar.ClientConfig{
					Addr: url,
					Errs: *asyncErrs,
				},
			},
		}
		mp.producers = append(mp.producers, pulsar.NewManagedProducer(mcp, mpCfg))
		logger.Info(fmt.Sprintf("%s try to create producer on topic %s on url:%s", mpCfg.Name, mpCfg.Topic, url))
	}
	return mp
}

func (pp *PulsarProducer) RefreshMarketIDMap(data *kline.KlineData, logger log.Logger) error {
	logger.Debug(fmt.Sprintf("marketServiceEnable:%v, nacosUrls:%s, marketNacosServiceName:%s",
		pp.MarketServiceEnable, pp.MarketNacosUrls, pp.MarketNacosServiceName))
	for _, tokenPair := range data.GetNewTokenPairs() {
		tokenPairName := tokenPair.Name()
		marketIDMap := kline.GetMarketIDMap()
		marketIDMap[tokenPairName] = int64(tokenPair.ID)
		logger.Debug(fmt.Sprintf("set new tokenpair %+v in map, MarketIdMap: %+v", tokenPair, marketIDMap))

		if pp.MarketServiceEnable {
			param := vo.SelectOneHealthInstanceParam{Clusters: pp.MarketNacosClusters, ServiceName: pp.MarketNacosServiceName, GroupName: pp.MarketNacosGroupName}
			marketServiceURL, err := kline.GetMarketServiceURL(pp.MarketNacosUrls, pp.MarketNacosNamespaceId, param)
			if err == nil {
				logger.Debug(fmt.Sprintf("successfully get the market service url [%s]", marketServiceURL))
			} else {
				logger.Error(fmt.Sprintf("failed to get the market service url [%s]. error: %s", marketServiceURL, err))
			}

			err = kline.RegisterNewTokenPair(int64(tokenPair.ID), tokenPairName, marketServiceURL, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("failed register tokenpair %+v in market service. error: %s", tokenPair, err))
				return err
			}
		}
	}
	return nil
}

func (pp *PulsarProducer) SendAllMsg(data *kline.KlineData, logger log.Logger) (map[string]int, error) {
	// log := logger.With("module", "pulsar")
	result := make(map[string]int)

	return result, nil
}
