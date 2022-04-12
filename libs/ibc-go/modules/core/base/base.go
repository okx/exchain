package base

import (
	"github.com/okex/exchain/libs/cosmos-sdk/store"
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/libs/tendermint/types"
)

var (
	_ upgrade.UpgradeModule = (*BaseIBCUpgradeModule)(nil)

	defaultHandleStore upgrade.HandleStore = func(st cosmost.CommitKVStore, h int64) {
		st.SetUpgradeVersion(h)
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

func (b *BaseIBCUpgradeModule) BlockStoreModules() map[string]upgrade.HandleStore {
	return map[string]upgrade.HandleStore{
		"ibc":            defaultHandleStore,
		"mem_capability": defaultHandleStore,
		"capability":     defaultHandleStore,
		"transfer":       defaultHandleStore,
		"erc20":          defaultHandleStore,
	}
}

func (b *BaseIBCUpgradeModule) RegisterParam() params.ParamSet {
	return nil
}

func (b *BaseIBCUpgradeModule) HandleStoreWhenMeetUpgradeHeight() upgrade.HandleStore {
	return func(st store.CommitKVStore, h int64) {
		st.SetUpgradeVersion(h)
	}
}

func (b *BaseIBCUpgradeModule) Seal() {
	b.Inited = true
}
func (b *BaseIBCUpgradeModule) Sealed() bool {
	return b.Inited
}
