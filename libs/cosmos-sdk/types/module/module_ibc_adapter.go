package module

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
)

// RegisterInterfaces registers all module interface types
func (bm BasicManager) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	for _, m := range bm {
		if ada, ok := m.(AppModuleBasicAdapter); ok {
			ada.RegisterInterfaces(registry)
		}
	}
}

// RegisterGRPCGatewayRoutes registers all module rest routes
func (bm BasicManager) RegisterGRPCGatewayRoutes(clientCtx clientCtx.CLIContext, rtr *runtime.ServeMux) {
	for _, m := range bm {
		if ada, ok := m.(AppModuleBasicAdapter); ok {
			ada.RegisterGRPCGatewayRoutes(clientCtx, rtr)
		}
	}
}
