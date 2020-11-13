package common

import (
	"fmt"
	"strconv"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/okex/okexchain/x/backend"
	"github.com/okex/okexchain/x/dex"
	"github.com/okex/okexchain/x/stream/nacos"
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

func (kd *KlineData) SetMatchResults(matchResults []*backend.MatchResult) {
	kd.matchResults = matchResults
}

func GetMarketServiceURL(urls string, nameSpace string, param vo.SelectOneHealthInstanceParam) (string, error) {
	k, err := nacos.GetOneInstance(urls, nameSpace, param)
	if err != nil {
		return "", err
	}
	if k == nil {
		return "", fmt.Errorf("there is no %s service in nacos-server %s", param.ServiceName, urls)
	}
	port := strconv.FormatUint(k.Port, 10)
	return k.Ip + ":" + port + "/manager/add", nil
}

func RegisterNewTokenPair(tokenPairID int64, tokenPairName string, marketServiceURL string, logger log.Logger) (err error) {
	return nil
}
