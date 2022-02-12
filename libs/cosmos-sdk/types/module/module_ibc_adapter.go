package module

import codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"

// RegisterInterfaces registers all module interface types
func (bm BasicManager) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	for _, m := range bm {
		if ada, ok := m.(AppModuleBasicAdapter); ok {
			ada.RegisterInterfaces(registry)
		}
	}
}
