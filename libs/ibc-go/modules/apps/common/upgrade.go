package common

import (
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/ibc-go/modules/core/base"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

var (
	_ upgrade.UpgradeModule = (*Veneus3BaseUpgradeModule)(nil)

	ibcV4Map = map[string]struct{}{
		"feeibc":             {},
		"interchainaccounts": {},
		"icacontroller":      {},
		"icahost":            {},
		"icamauth":           {},
	}

	defaultDenyFilter cosmost.StoreFilter = func(module string, h int64, store cosmost.CommitKVStore) bool {
		_, exist := ibcV4Map[module]
		if !exist {
			return false
		}
		return true
	}
	defaultIBCCommitFilter cosmost.StoreFilter = func(module string, h int64, store cosmost.CommitKVStore) bool {
		_, exist := ibcV4Map[module]
		if !exist {
			return false
		}

		// ==veneus1
		if h == tmtypes.GetVenus4Height() {
			if store != nil {
				store.SetUpgradeVersion(h)
			}
			return false
		}

		// ibc modules
		if tmtypes.HigherThanVenus4(h) {
			return false
		}

		// < veneus1
		return true
	}
	defaultIBCPruneFilter cosmost.StoreFilter = func(module string, h int64, store cosmost.CommitKVStore) bool {
		_, exist := ibcV4Map[module]
		if !exist {
			return false
		}

		// ibc modulee && >=veneus1
		if tmtypes.HigherThanVenus4(h) {
			return false
		}

		// < veneus1
		return true
	}
	defaultIBCVersionFilter cosmost.VersionFilter = func(h int64) func(callback cosmost.VersionCallback) {
		if h < 0 {
			return func(callback cosmost.VersionCallback) {}
		}
		return func(callback cosmost.VersionCallback) {
			for name, _ := range ibcV4Map {
				hh := tmtypes.GetVenus4Height()
				callback(name, hh)
			}
		}
	}
)

type Veneus3BaseUpgradeModule struct {
	*base.BaseIBCUpgradeModule
}

func NewVeneus3BaseUpgradeModule(m module.AppModuleBasic) *Veneus3BaseUpgradeModule {
	ret := &Veneus3BaseUpgradeModule{}
	ret.BaseIBCUpgradeModule = base.NewBaseIBCUpgradeModule(m)

	return ret
}

func (v *Veneus3BaseUpgradeModule) CommitFilter() *cosmost.StoreFilter {
	if v.UpgradeHeight() == 0 {
		return &defaultDenyFilter
	}
	return &defaultIBCCommitFilter
}

func (v *Veneus3BaseUpgradeModule) PruneFilter() *cosmost.StoreFilter {
	if v.UpgradeHeight() == 0 {
		return &defaultDenyFilter
	}
	return &defaultIBCPruneFilter
}

func (v *Veneus3BaseUpgradeModule) VersionFilter() *cosmost.VersionFilter {
	return &defaultIBCVersionFilter
}

func (v *Veneus3BaseUpgradeModule) UpgradeHeight() int64 {
	return tmtypes.GetVenus4Height()
}
