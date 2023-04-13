package feesplit

import (
	store "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/feesplit/types"
)

var (
	defaultVersionFilter store.VersionFilter = func(h int64) func(cb func(name string, version int64)) {
		if h < 0 {
			return func(cb func(name string, version int64)) {}
		}

		return func(cb func(name string, version int64)) {
			cb(ModuleName, tmtypes.GetVenus3Height())
		}
	}
)

func (am AppModule) RegisterTask() upgrade.HeightTask {
	return upgrade.NewHeightTask(
		0, func(ctx sdk.Context) error {
			if am.Sealed() {
				return nil
			}
			InitGenesis(ctx, am.keeper, types.DefaultGenesisState())
			return nil
		})
}

func (am AppModule) CommitFilter() *store.StoreFilter {
	var filter store.StoreFilter
	filter = func(module string, h int64, s store.CommitKVStore) bool {
		if module != ModuleName {
			return false
		}
		if am.UpgradeHeight() == 0 {
			return true
		}
		if h == tmtypes.GetVenus3Height() {
			if s != nil {
				s.SetUpgradeVersion(h)
			}
			return false
		}

		if tmtypes.HigherThanVenus3(h) {
			return false
		}

		return true
	}
	return &filter
}

func (am AppModule) PruneFilter() *store.StoreFilter {
	var filter store.StoreFilter
	filter = func(module string, h int64, s store.CommitKVStore) bool {
		if module != ModuleName {
			return false
		}

		if am.UpgradeHeight() == 0 {
			return true
		}
		if tmtypes.HigherThanVenus3(h) {
			return false
		}

		return true
	}
	return &filter
}

func (am AppModule) VersionFilter() *store.VersionFilter {
	return &defaultVersionFilter
}

func (am AppModule) UpgradeHeight() int64 {
	return tmtypes.GetVenus3Height()
}
