package base

import (
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/libs/tendermint/types"
)

var (
	_      upgrade.UpgradeModule = (*BaseIBCUpgradeModule)(nil)
	ibcMap                       = map[string]struct{}{
		"ibc":            struct{}{},
		"mem_capability": struct{}{},
		"capability":     struct{}{},
		"transfer":       struct{}{},
		"erc20":          struct{}{},
	}
	defaultDenyFilter cosmost.StoreFilter = func(module string, h int64, store cosmost.CommitKVStore) bool {
		_, exist := ibcMap[module]
		if !exist {
			return false
		}
		return true
	}
	defaultIBCCommitFilter cosmost.StoreFilter = func(module string, h int64, store cosmost.CommitKVStore) bool {
		_, exist := ibcMap[module]
		if !exist {
			return false
		}

		// ==veneus1
		if h == types.GetVenus1Height() {
			if store != nil {
				store.SetUpgradeVersion(h)
			}
			return false
		}

		// ibc modules
		if types.HigherThanVenus1(h) {
			return false
		}

		// < veneus1
		return true
	}
	defaultIBCPruneFilter cosmost.StoreFilter = func(module string, h int64, store cosmost.CommitKVStore) bool {
		_, exist := ibcMap[module]
		if !exist {
			return false
		}

		// ibc modulee && >=veneus1
		if types.HigherThanVenus1(h) {
			return false
		}

		// < veneus1
		return true
	}
	defaultIBCVersionFilter cosmost.VersionFilter = func(h int64) func(cb func(name string, version int64)) {
		if h < 0 {
			return func(cb func(name string, version int64)) {}
		}

		return func(cb func(name string, version int64)) {
			for name, _ := range ibcMap {
				hh := types.GetVenus1Height()
				cb(name, hh)
			}
		}
	}
)

type BaseIBCUpgradeModule struct {
	appModule module.AppModuleBasic
	Inited    bool
}

func NewBaseIBCUpgradeModule(appModule module.AppModuleBasic) *BaseIBCUpgradeModule {
	return &BaseIBCUpgradeModule{appModule: appModule}
}

func (b *BaseIBCUpgradeModule) ModuleName() string {
	return b.appModule.Name()
}

func (b *BaseIBCUpgradeModule) RegisterTask() upgrade.HeightTask {
	panic("override")
}

func (b *BaseIBCUpgradeModule) UpgradeHeight() int64 {
	return types.GetVenus1Height()
}

func (b *BaseIBCUpgradeModule) CommitFilter() *cosmost.StoreFilter {
	if b.UpgradeHeight() == 0 {
		return &defaultDenyFilter
	}
	return &defaultIBCCommitFilter
}
func (b *BaseIBCUpgradeModule) PruneFilter() *cosmost.StoreFilter {
	if b.UpgradeHeight() == 0 {
		return &defaultDenyFilter
	}
	return &defaultIBCPruneFilter
}

func (b *BaseIBCUpgradeModule) VersionFilter() *cosmost.VersionFilter {
	return &defaultIBCVersionFilter
}

func (b *BaseIBCUpgradeModule) RegisterParam() params.ParamSet {
	return nil
}

func (b *BaseIBCUpgradeModule) Seal() {
	b.Inited = true
}
func (b *BaseIBCUpgradeModule) Sealed() bool {
	return b.Inited
}
