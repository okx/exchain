package common

import (
	cosmost "github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/module"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okx/okbchain/libs/ibc-go/modules/core/base"
	tmtypes "github.com/okx/okbchain/libs/tendermint/types"
)

var (
	_ upgrade.UpgradeModule = (*Venus3BaseUpgradeModule)(nil)

	ibcV4Map = map[string]struct{}{
		"feeibc":             {},
		"interchainaccounts": {},
		"icacontroller":      {},
		"icahost":            {},
		"icamauth":           {},
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

type Venus3BaseUpgradeModule struct {
	*base.BaseIBCUpgradeModule
}

func NewVenus3BaseUpgradeModule(m module.AppModuleBasic) *Venus3BaseUpgradeModule {
	ret := &Venus3BaseUpgradeModule{}
	ret.BaseIBCUpgradeModule = base.NewBaseIBCUpgradeModule(m)

	return ret
}

func (v *Venus3BaseUpgradeModule) CommitFilter() *cosmost.StoreFilter {
	var filter cosmost.StoreFilter
	filter = func(module string, h int64, s cosmost.CommitKVStore) bool {
		_, exist := ibcV4Map[module]
		if !exist {
			return false
		}

		if v.UpgradeHeight() == 0 {
			return true
		}

		if h == tmtypes.GetVenus4Height() {
			if s != nil {
				s.SetUpgradeVersion(h)
			}
			return false
		}

		if tmtypes.HigherThanVenus4(h) {
			return false
		}

		return true
	}
	return &filter
}

func (v *Venus3BaseUpgradeModule) PruneFilter() *cosmost.StoreFilter {
	var filter cosmost.StoreFilter
	filter = func(module string, h int64, s cosmost.CommitKVStore) bool {
		_, exist := ibcV4Map[module]
		if !exist {
			return false
		}

		if v.UpgradeHeight() == 0 {
			return true
		}
		// ibc module && >=venus4
		if tmtypes.HigherThanVenus4(h) {
			return false
		}

		return true
	}
	return &filter
}

func (v *Venus3BaseUpgradeModule) VersionFilter() *cosmost.VersionFilter {
	return &defaultIBCVersionFilter
}

func (v *Venus3BaseUpgradeModule) UpgradeHeight() int64 {
	return tmtypes.GetVenus4Height()
}
