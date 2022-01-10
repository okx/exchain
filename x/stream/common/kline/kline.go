package kline

import (
	"sync"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/dex"
	"github.com/okex/exchain/x/stream/common"
	"github.com/okex/exchain/x/stream/types"
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
	MarketServiceEnable    bool
	MarketNacosUrls        string
	MarketNacosNamespaceId string
	MarketNacosClusters    []string
	MarketNacosServiceName string
	MarketNacosGroupName   string
}

func NewMarketConfig(enable bool, urls, nameSpace string, clusters []string, serviceName, groupName string) MarketConfig {
	return MarketConfig{
		MarketServiceEnable:    enable,
		MarketNacosUrls:        urls,
		MarketNacosNamespaceId: nameSpace,
		MarketNacosClusters:    clusters,
		MarketNacosServiceName: serviceName,
		MarketNacosGroupName:   groupName,
	}
}

type KlineData struct {
	Height        int64
	newTokenPairs []*dex.TokenPair
}

func NewKlineData() *KlineData {
	return &KlineData{}
}

func (kd KlineData) BlockHeight() int64 {
	return kd.Height
}

func (kd KlineData) DataType() types.StreamDataKind {
	return types.StreamDataKlineKind
}

func (kd *KlineData) SetData(ctx sdk.Context, orderKeeper types.OrderKeeper, cache *common.Cache) {
	kd.Height = ctx.BlockHeight()
	kd.newTokenPairs = cache.GetNewTokenPairs()
}

func (kd *KlineData) GetNewTokenPairs() []*dex.TokenPair {
	return kd.newTokenPairs
}
