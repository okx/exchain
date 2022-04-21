package staking

import (
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	params2 "github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/x/params"
	"github.com/okex/exchain/x/staking/types"
)

var (
	_ upgrade.UpgradeModule = AppModule{}
)

func (am AppModule) RegisterTask() upgrade.HeightTask {
	return nil
}

func (am AppModule) UpgradeHeight() int64 {
	return -1
}

func (am AppModule) BlockStoreModules() map[string]upgrade.HandleStore {
	return nil
}

func (am AppModule) RegisterParam() params.ParamSet {
	v := types.KeyHistoricalEntriesParams(7)
	return params2.ParamSet(v)
}

func (am AppModule) ModuleName() string {
	return ModuleName
}

func (b AppModule) CommitFilter() *cosmost.StoreFilter {
	return nil
}

func (b AppModule) PruneFilter() *cosmost.StoreFilter {
	return nil
}
