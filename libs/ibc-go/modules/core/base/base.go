package base

import (
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/tendermint/types"
)

var (
	_ module.UpgradeModule = (*BaseIBCUpgradeModule)(nil)
)

type BaseIBCUpgradeModule struct {
	appModule module.AppModuleBasicAdapter
}

func NewBaseIBCUpgradeModule(appModule module.AppModuleBasicAdapter) *BaseIBCUpgradeModule {
	return &BaseIBCUpgradeModule{appModule: appModule}
}

func (b *BaseIBCUpgradeModule) ModuleName() string {
	return b.appModule.Name()
}

func (b *BaseIBCUpgradeModule) RegisterTask() module.HeightTask {
	panic("override")
}

func (b *BaseIBCUpgradeModule) UpgradeHeight() int64 {
	if types.GetIBCHeight() == 1 {
		return 1
	}
	return types.GetIBCHeight() + 1
}

func (b *BaseIBCUpgradeModule) BlockStoreModules() []string {
	return []string{"ibc", "mem_capability", "capability", "transfer"}
}
