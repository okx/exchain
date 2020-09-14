package pulsarclient

import (
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/stream/eureka"
	"github.com/okex/okexchain/x/stream/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	marketIDMap = make(map[string]int64, 200)
	initMapOnce sync.Once
)

func InitTokenPairMap(ctx sdk.Context, dexKeeper types.DexKeeper) {
	initMapOnce.Do(func() {
		tokenPairs := dexKeeper.GetTokenPairs(ctx)
		for i := 0; i < len(tokenPairs); i++ {
			marketIDMap[tokenPairs[i].Name()] = int64(tokenPairs[i].ID)
		}
	})
}

func (p *PulsarProducer) RefreshMarketIDMap(data *PulsarData, logger log.Logger) error {
	logger.Debug(fmt.Sprintf("marketServiceEnable:%v, eurekaUrl:%s, registerAppName:%s", p.marketServiceEnable, p.marketEurekaURL, p.marketEurekaRegisteredAppName))
	for _, tokenPair := range data.newTokenPairs {
		tokenPairName := tokenPair.Name()
		marketIDMap[tokenPairName] = int64(tokenPair.ID)
		logger.Debug(fmt.Sprintf("set new tokenpair %+v in map, MarketIdMap: %+v", tokenPair, marketIDMap))

		if p.marketServiceEnable {
			marketServiceURL, err := getMarketServiceURL(p.marketEurekaURL, p.marketEurekaRegisteredAppName)
			if err == nil {
				logger.Debug(fmt.Sprintf("successfully get the market service url [%s]", marketServiceURL))
			} else {
				logger.Error(fmt.Sprintf("failed to get the market service url [%s]. error: %s", marketServiceURL, err))
			}

			err = RegisterNewTokenPair(int64(tokenPair.ID), tokenPairName, marketServiceURL, logger)
			if err != nil {
				logger.Error(fmt.Sprintf("failed register tokenpair %+v in market service. error: %s", tokenPair, err))
				return err
			}
		}
	}
	return nil
}

func getMarketServiceURL(eurekaURL, registerAppName string) (string, error) {
	k, err := eureka.GetOneInstance(eurekaURL, registerAppName)
	if err != nil {
		return "", err
	}
	if len(k.Instances) == 0 {
		return "", fmt.Errorf("failed to find instance %s in eureka-server %s", registerAppName, eurekaURL)
	}
	return k.Instances[0].HomePageURL + "manager/add", nil
}
