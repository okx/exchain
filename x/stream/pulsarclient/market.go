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
	marketIdMap    = make(map[string]int64, 500)
	marketIdMapMux sync.RWMutex
	initMapOnce    sync.Once
)

func InitTokenPairMap(ctx sdk.Context, dexKeeper types.DexKeeper) {
	initMapOnce.Do(func() {
		tokenpairs := dexKeeper.GetTokenPairs(ctx)
		for i := 0; i < len(tokenpairs); i++ {
			tokenpairname := getTokenPairName(tokenpairs[i].BaseAssetSymbol, tokenpairs[i].QuoteAssetSymbol)
			UpdateMarketIdMap(tokenpairname, int64(tokenpairs[i].ID))
		}
	})
}

// UpdateMarketIdMap updates market id map
func UpdateMarketIdMap(key string, value int64) {
	marketIdMapMux.Lock()
	defer marketIdMapMux.Unlock()
	marketIdMap[key] = value

}

func GetMarketIdFromMap(key string) (int64, bool) {
	marketIdMapMux.RLock()
	defer marketIdMapMux.RUnlock()
	id, ok := marketIdMap[key]
	fmt.Println("GetMarketIdFromMap", key, id, ok)
	return id, ok
}

func GetMarketIdMap() map[string]int64 {
	marketIdMapMux.RLock()
	defer marketIdMapMux.RUnlock()
	idMap := make(map[string]int64, len(marketIdMap))
	for k, v := range marketIdMap {
		idMap[k] = v
	}
	return idMap
}

func (p *PulsarProducer) RefreshMarketIdMap(data *PulsarData, logger log.Logger) error {
	logger.Debug(fmt.Sprintf("marketServiceEnable:%v, eurekaUrl:%s, registerAppName:%s", p.marketServiceEnable, p.marketEurekaUrl, p.marketEurekaRegisteredAppName))
	for _, tokenPair := range data.newTokenPairs {
		tokenpairName := getTokenPairName(tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)
		UpdateMarketIdMap(tokenpairName, int64(tokenPair.ID))
		logger.Debug(fmt.Sprintf("set new tokenpair %+v in map, MarketIdMap: %+v", tokenPair, marketIdMap))

		if p.marketServiceEnable {
			marketServiceUrl, err := getMarketServiceUrl(p.marketEurekaUrl, p.marketEurekaRegisteredAppName)
			if err == nil {
				logger.Debug(fmt.Sprintf("successfully get the market service url [%s]", marketServiceUrl))
			} else {
				logger.Error(fmt.Sprintf("failed to get the market service url [%s]. error: %s", marketServiceUrl, err))
			}

			err = RegisterNewTokenPair(int64(tokenPair.ID), tokenpairName, marketServiceUrl, logger)
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
