package kafkaclient

import (
	appcfg "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/okex/okexchain/x/stream/common"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	common.MarketConfig
	*kafka.Writer
}

func NewKafkaProducer(url string, cfg *appcfg.StreamConfig) *KafkaProducer {
	return &KafkaProducer{
		MarketConfig: common.NewMarketConfig(cfg.MarketServiceEnable, cfg.EurekaServerUrl, cfg.MarketQuotationsEurekaName),
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{url},
			Topic:    cfg.MarketPulsarTopic,
			Balancer: newOKDExBalancer(int64(cfg.MarketPulsarPartition)),
		}),
	}
}
