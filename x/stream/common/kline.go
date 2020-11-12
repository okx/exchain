package common

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/backend"
	"github.com/okex/okexchain/x/dex"
	"github.com/okex/okexchain/x/stream/eureka"
	"github.com/okex/okexchain/x/stream/types"
	"github.com/tendermint/tendermint/libs/log"
	"sync"
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

func GetMarketIDMap() map[string]int64 {
	return marketIDMap
}

type MarketConfig struct {
	MarketServiceEnable           bool
	MarketEurekaURL               string
	MarketEurekaRegisteredAppName string
}

func NewMarketConfig(marketServiceEnable bool, marketEurekaURL, marketEurekaRegisteredAppName string) MarketConfig {
	return MarketConfig{
		MarketServiceEnable:           marketServiceEnable,
		MarketEurekaURL:               marketEurekaURL,
		MarketEurekaRegisteredAppName: marketEurekaRegisteredAppName,
	}
}

type KlineData struct {
	Height        int64
	matchResults  []*backend.MatchResult
	newTokenPairs []*dex.TokenPair
}

func NewKlineData() *KlineData {
	return &KlineData{
		matchResults: make([]*backend.MatchResult, 0),
	}
}

func (kd KlineData) BlockHeight() int64 {
	return kd.Height
}

func (kd KlineData) DataType() types.StreamDataKind {
	return types.StreamDataKlineKind
}

func (kd *KlineData) SetData(ctx sdk.Context, orderKeeper types.OrderKeeper, cache *Cache) {
	kd.Height = ctx.BlockHeight()
	kd.matchResults = GetMatchResults(ctx, orderKeeper)
	kd.newTokenPairs = cache.GetNewTokenPairs()
}

func (kd *KlineData) GetNewTokenPairs() []*dex.TokenPair {
	return kd.newTokenPairs
}

func (kd *KlineData) GetMatchResults() []*backend.MatchResult {
	return kd.matchResults
}

func GetMarketServiceURL(eurekaURL, registerAppName string) (string, error) {
	k, err := eureka.GetOneInstance(eurekaURL, registerAppName)
	if err != nil {
		return "", err
	}
	if len(k.Instances) == 0 {
		return "", fmt.Errorf("failed to find instance %s in eureka-server %s", registerAppName, eurekaURL)
	}
	return k.Instances[0].HomePageURL, nil
}

func RegisterNewTokenPair(tokenPairID int64, tokenPairName string, marketServiceURL string, logger log.Logger) (err error) {
	return nil
}
