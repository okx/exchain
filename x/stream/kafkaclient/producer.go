package kafkaclient

import (
	"context"
	"encoding/json"
	"fmt"
	appcfg "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/okex/okexchain/x/backend"
	"github.com/okex/okexchain/x/stream/common/kline"
	"github.com/segmentio/kafka-go"
	"github.com/tendermint/tendermint/libs/log"
	"sync"
	"time"
)

type KafkaProducer struct {
	kline.MarketConfig
	*kafka.Writer
}

func NewKafkaProducer(url string, cfg *appcfg.StreamConfig) *KafkaProducer {
	return &KafkaProducer{
		MarketConfig: kline.NewMarketConfig(cfg.MarketServiceEnable, cfg.EurekaServerUrl, cfg.MarketQuotationsEurekaName),
		Writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{url},
			Topic:    cfg.MarketPulsarTopic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (kp *KafkaProducer) RefreshMarketIDMap(data *kline.KlineData, logger log.Logger) error {
	logger.Debug(fmt.Sprintf(
		"marketServiceEnable:%v, eurekaUrl:%s, registerAppName:%s",
		kp.MarketServiceEnable, kp.MarketEurekaURL, kp.MarketEurekaRegisteredAppName,
	))

	for _, tokenPair := range data.GetNewTokenPairs() {
		tokenPairName := tokenPair.Name()
		marketIDMap := kline.GetMarketIDMap()
		marketIDMap[tokenPairName] = int64(tokenPair.ID)
		logger.Debug(fmt.Sprintf("set new tokenpair %+v in map, MarketIdMap: %+v", tokenPair, marketIDMap))

		if kp.MarketServiceEnable {
			marketServiceURL, err := kline.GetMarketServiceURL(kp.MarketEurekaURL, kp.MarketEurekaRegisteredAppName)
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

func (kp *KafkaProducer) SendAllMsg(data *kline.KlineData, logger log.Logger) (map[string]int, error) {
	// log := logger.With("module", "kafka")
	result := make(map[string]int)
	matchResults := data.GetMatchResults()
	result["matchResults"] = len(matchResults)
	if len(matchResults) == 0 {
		return result, nil
	}

	var errChan = make(chan error, len(matchResults))
	var wg sync.WaitGroup
	wg.Add(len(matchResults))
	for _, matchResult := range matchResults {
		go func(matchResult backend.MatchResult) {
			defer wg.Done()
			marketID, ok := kline.GetMarketIDMap()[matchResult.Product]
			if !ok {
				err := fmt.Errorf("failed to find %s marketId", matchResult.Product)
				errChan <- err
				return
			}

			msg, err := json.Marshal(&matchResult)
			if err != nil {
				errChan <- err
				return
			}

			if err = kp.WriteMessages(context.Background(),
				kafka.Message{
					Key:   getKafkaMsgKey(marketID),
					Value: msg,
				},
			); err != nil {
				errChan <- err
				return
			}

			logger.Debug(
				fmt.Sprintf("successfully send matchResult [marketId:%d, CreatedTime:%s, BlockHeight:%d, Quantity:%f, Price:%f, InstrumentName:%s]",
					marketID, time.Unix(matchResult.Timestamp, 0).Format("2006-01-02 15:04:05"), matchResult.BlockHeight,
					matchResult.Quantity, matchResult.Price, matchResult.Product,
				),
			)
		}(*matchResult)
	}
	wg.Wait()

	if len(errChan) != 0 {
		err := <-errChan
		return result, err
	}
	return result, nil
}
