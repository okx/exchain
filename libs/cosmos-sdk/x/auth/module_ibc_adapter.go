package auth

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	types3 "github.com/okex/exchain/temp"
	//clienttypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/types"
	//connectiontypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/03-connection/types"
	//channeltypes "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/04-channel/types"
)

//// RegisterServices registers a GRPC query service to respond to the
//// module-specific GRPC queries.
//func (am AppModule) RegisterServices(cfg module.Configurator) {
//	types.RegisterQueryServer(cfg.QueryServer(), am.accountKeeper)
//	m := keeper.NewMigrator(am.accountKeeper, cfg.QueryServer())
//	err := cfg.RegisterMigration(types.ModuleName, 1, m.Migrate1to2)
//	if err != nil {
//		panic(err)
//	}
//}
var (
	_ module.AppModuleAdapter = AppModule{}
)

func (am AppModule) RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	RegisterInterfaces(registry)
}

func (am AppModule) RegisterGRPCGatewayRoutes(cliContext context.CLIContext, serveMux *runtime.ServeMux) {
}

func (am AppModule) Upgrade(req *abci.UpgradeReq) (*abci.ModuleUpgradeResp, error) {
	return nil, nil
}

// RegisterServices registers module services.
func (am AppModule) RegisterServices(cfg module.Configurator) {
	a := &am.accountKeeper
	types.RegisterQueryServer(cfg.QueryServer(), a)
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	//registry.RegisterInterface(
	//	"cosmos.vesting.v1beta1.VestingAccount",
	//	(*exported.VestingAccount)(nil),
	//	&ContinuousVestingAccount{},
	//	&DelayedVestingAccount{},
	//	&PeriodicVestingAccount{},
	//	&PermanentLockedAccount{},
	//)

	registry.RegisterImplementations(
		(*exported.Account)(nil),
		//&BaseVestingAccount{},
		//&DelayedVestingAccount{},
		//&ContinuousVestingAccount{},
		//&PeriodicVestingAccount{},
		//&PermanentLockedAccount{},
		&types3.BaseAccount{},
	)

	//registry.RegisterImplementations(
	//	(*authtypes.GenesisAccount)(nil),
	//	&BaseVestingAccount{},
	//	&DelayedVestingAccount{},
	//	&ContinuousVestingAccount{},
	//	&PeriodicVestingAccount{},
	//	&PermanentLockedAccount{},
	//)

	//registry.RegisterImplementations(
	//	(*sdk.Msg)(nil),
	//	&MsgCreateVestingAccount{},
	//)

	//msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
