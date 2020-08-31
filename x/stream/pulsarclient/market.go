package pulsarclient

import (
	"errors"
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/stream/eureka"
	"github.com/okex/okchain/x/stream/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	marketIdMap = make(map[string]int64, 500)
	initMapOnce sync.Once
)

func InitTokenPairMap(ctx sdk.Context, dexKeeper types.DexKeeper) {
	initMapOnce.Do(func() {
		tokenPairs := dexKeeper.GetTokenPairs(ctx)
		for i := 0; i < len(tokenPairs); i++ {
			marketIdMap[tokenPairs[i].Name()] = int64(tokenPairs[i].ID)
		}
	})
}

func (p *PulsarProducer) RefreshMarketIdMap(data *PulsarData, logger log.Logger) error {
	logger.Debug(fmt.Sprintf("marketServiceEnable:%v, eurekaUrl:%s, registerAppName:%s", p.marketServiceEnable, p.marketEurekaUrl, p.marketEurekaRegisteredAppName))
	for _, tokenPair := range data.newTokenPairs {
		tokenPairName := tokenPair.Name()
		marketIdMap[tokenPairName] = int64(tokenPair.ID)
		logger.Debug(fmt.Sprintf("set new tokenpair %+v in map, MarketIdMap: %+v", tokenPair, marketIdMap))

		if p.marketServiceEnable {
			marketServiceUrl, err := getMarketServiceUrl(p.marketEurekaUrl, p.marketEurekaRegisteredAppName)
			if err == nil {
				logger.Debug(fmt.Sprintf("successfully get the market service url [%s]", marketServiceUrl))
			} else {
				logger.Error(fmt.Sprintf("failed to get the market service url [%s]. error: %s", marketServiceUrl, err))
			}

			err = RegisterNewTokenPair(int64(tokenPair.ID), tokenPairName, marketServiceUrl, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("failed register tokenpair %+v in market service. error: %s", tokenPair, err))
				return err
			}
		}
	}
	return nil
}

func getTokenPairName(BaseAssetSymbol, QuoteAssetSymbol string) string {
	pairName := BaseAssetSymbol + "_" + QuoteAssetSymbol
	return pairName
}

func getMarketServiceUrl(eurekaUrl, registerAppName string) (string, error) {
	k, err := eureka.GetOneInstance(eurekaUrl, registerAppName)
	if err != nil {
		return "", err
	}
	if len(k.Instances) == 0 {
		return "", errors.New(fmt.Sprintf("failed to find instance %s in eureka-server %s", registerAppName, eurekaUrl))
	} else {
		return k.Instances[0].HomePageURL + "manager/add", nil
	}
}
