package module

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	clictx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

// AppModuleBasic is the standard form for basic non-dependant elements of an application module.
type AppModuleBasicAdapter interface {
	AppModuleBasic
	RegisterInterfaces(codectypes.InterfaceRegistry)
	// client functionality
	RegisterGRPCGatewayRoutes(clictx.CLIContext, *runtime.ServeMux)
}

// AppModuleGenesis is the standard form for an application module genesis functions
type AppModuleGenesisAdapter interface {
	AppModuleGenesis
	AppModuleBasicAdapter
}

// AppModule is the standard form for an application module
type AppModuleAdapter interface {
	AppModule
	AppModuleGenesisAdapter
	// registers
	RegisterInvariants(sdk.InvariantRegistry)
	// RegisterServices allows a module to register services
	RegisterServices(Configurator)
}
