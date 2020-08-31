package pulsarclient

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Comcast/pulsar-client-go"
	appCfg "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/okex/okchain/x/backend"
	"github.com/tendermint/tendermint/libs/log"
)

type PulsarProducer struct {
	producers                     []*pulsar.ManagedProducer
	partion                       int64
	marketServiceEnable           bool
	marketEurekaUrl               string
	marketEurekaRegisteredAppName string
}

func NewPulsarProducer(url string, cfg *appCfg.StreamConfig, logger log.Logger, asyncErrs *chan error) *PulsarProducer {
	var mp = &PulsarProducer{producers: make([]*pulsar.ManagedProducer, 0, cfg.MarketPulsarPartition),
		partion: int64(cfg.MarketPulsarPartition), marketServiceEnable: cfg.MarketServiceEnable, marketEurekaUrl: cfg.EurekaServerUrl,
		marketEurekaRegisteredAppName: cfg.MarketQuotationsEurekaName}

	for i := 0; i < cfg.MarketPulsarPartition; i++ {
		mcp := pulsar.NewManagedClientPool()
		mpCfg := pulsar.ManagedProducerConfig{
			Name:                  uuid.New().String() + "-subs_standard_dex_spot-" + strconv.Itoa(i),
			Topic:                 cfg.MarketPulsarTopic + "-partition-" + strconv.Itoa(i),
			NewProducerTimeout:    time.Second * 3,
			InitialReconnectDelay: time.Second,
			MaxReconnectDelay:     time.Minute,
			ManagedClientConfig: pulsar.ManagedClientConfig{
				ClientConfig: pulsar.ClientConfig{
					Addr: url, //"192.168.80.186:6650",
					Errs: *asyncErrs,
				},
			},
		}
		mp.producers = append(mp.producers, pulsar.NewManagedProducer(mcp, mpCfg))
		logger.Info(fmt.Sprintf("%s try to create producer on topic %s on url:%s", mpCfg.Name, mpCfg.Topic, url))
	}
	return mp
}

var bizType = MARKET_CAL_SERVICE_DEX_SPOT_BIZ_TYPE
var marketType = MARKET_CAL_SERVICE_DEX_SPOT_MARKET_TYPE
var iscalc = true

func (mp *PulsarProducer) SendAllMsg(data *PulsarData, logger log.Logger) (map[string]int, error) {
	// log := logger.With("module", "pulsar")
	result := make(map[string]int, 0)
	matchResults := data.matchResults
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
			marketId, ok := marketIdMap[matchResult.Product] //attention,maybe marketId is 0
			if !ok {
				err := fmt.Errorf("failed to find %s marketId", matchResult.Product)
				errChan <- err
				return
			}

			timestamp := matchResult.Timestamp * 1000
			matchResultMsg := MatchResultMsg{
				BizType:        &bizType,
				MarketId:       &marketId,
				MarketType:     &marketType,
				Size:           &matchResult.Quantity,
				Price:          &matchResult.Price,
				CreatedTime:    &timestamp,
				InstrumentId:   &marketId,
				InstrumentName: &matchResult.Product,
				IsCalc:         &iscalc,
			}

			msg, err := proto.Marshal(&matchResultMsg)
			if err != nil {
				errChan <- err
				return
			}

			sctx, cancle := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancle()
			_, err = mp.producers[marketId%(mp.partion)].Send(sctx, msg)
			if err != nil {
				errChan <- err
				return
			}
			logger.Debug(fmt.Sprintf("successfully send matchResult [marketId:%d, CreatedTime:%s, BlockHeight:%d, Quantity:%f, Price:%f, InstrumentName:%s]",
				marketId, time.Unix(matchResult.Timestamp, 0).Format("2006-01-02 15:04:05"), matchResult.BlockHeight, matchResult.Quantity, matchResult.Price, matchResult.Product))

			// b, _ := json.Marshal(matchResultMsg) //removed in production,
			// log.Debug("Send", "matchResultMsg", string(b), "height", matchResult.BlockHeight)
		}(*matchResult)
	}
	wg.Wait()

	if len(errChan) != 0 {
		err := <-errChan
		return result, err
	}
	return result, nil
}
