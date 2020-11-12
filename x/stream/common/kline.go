package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/backend"
	"github.com/okex/okexchain/x/dex"
	"github.com/okex/okexchain/x/stream/types"
)

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
