package module

import (
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/tendermint/abci/types"
)

// RegisterInterfaces registers all module interface types
func (bm BasicManager) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	for _, m := range bm {
		if ada, ok := m.(AppModuleBasicAdapter); ok {
			ada.RegisterInterfaces(registry)
		}
	}
}

// SetOrderBeginBlockers sets the order of set begin-blocker calls
func (m *Manager) SetOrderUpgrade(moduleNames ...string) {
	m.OrderUpgrades = moduleNames
}

func (bm Manager) Upgrade(req *types.UpgradeReq) (*types.UpgradeResp, error) {
	ret := new(types.UpgradeResp)
	for _, moduleName := range bm.OrderUpgrades {
		m, exist := bm.Modules[moduleName]
		if !exist {
			continue
		}
		ada, ok := m.(AppModuleAdapter)
		if !ok {
			continue
		}
		resp, err := ada.Upgrade(req)
		if nil != err {
			return nil, err
		}
		ret.ModuleResults = append(ret.ModuleResults, resp)
	}
	return ret, nil
}
