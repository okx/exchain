package auth

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/types/upgrade"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	authinternaltypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/internal"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

var (
	_ module.AppModuleAdapter = AppModule{}
)

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (am AppModule) RegisterGRPCGatewayRoutes(cliContext context.CLIContext, serveMux *runtime.ServeMux) {
}

func (am AppModule) RegisterTask() upgrade.HeightTask {
	return nil
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	a := &am.accountKeeper
	types.RegisterQueryServer(cfg.QueryServer(), a)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*exported.Account)(nil),
		&authinternaltypes.BaseAccount{},
	)
}
