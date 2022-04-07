package base

import (
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/libs/tendermint/types"
)

var (
	_ upgrade.UpgradeModule = (*BaseIBCUpgradeModule)(nil)
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

func (b *BaseIBCUpgradeModule) BlockStoreModules() []string {
	return []string{"ibc", "mem_capability", "capability", "transfer", "erc20"}
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
