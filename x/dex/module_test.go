package dex

import (
	cliLcd "github.com/okex/exchain/libs/cosmos-sdk/client/lcd"
	interfacetypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/common/version"
	"github.com/okex/exchain/x/dex/types"
	"github.com/stretchr/testify/require"
)

func TestAppModule_Smoke(t *testing.T) {
	_, _, spKeeper, dexKeeper, ctx := getMockTestCaseEvn(t)

	//func NewAppModule(version ProtocolVersionType, keeper Keeper, supplyKeeper SupplyKeeper) AppModule {
	appModule := NewAppModule(version.CurrentProtocolVersion, dexKeeper, spKeeper)

	// Const Info
	require.True(t, appModule.Name() == ModuleName)
	require.True(t, appModule.Route() == RouterKey)
	require.True(t, appModule.QuerierRoute() == QuerierRoute)
	require.True(t, appModule.GetQueryCmd(types.ModuleCdc) != nil)
	require.True(t, appModule.GetTxCmd(types.ModuleCdc) != nil)

	// RegisterCodec
	appModule.RegisterCodec(codec.New())

	appModule.RegisterInvariants(MockInvariantRegistry{})
	reg := interfacetypes.NewInterfaceRegistry()
	proc := codec.NewProtoCodec(reg)
	cdcProxy := codec.NewCodecProxy(proc, types.ModuleCdc)

	rs := cliLcd.NewRestServer(cdcProxy, nil)
	appModule.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
	handler := appModule.NewHandler()
	require.True(t, handler != nil)
	querior := appModule.NewQuerierHandler()
	require.True(t, querior != nil)

	// Initialization for genesis
	defaultGen := appModule.DefaultGenesis()
	err := appModule.ValidateGenesis(defaultGen)
	require.True(t, err == nil)

	illegalData := []byte{}
	err = appModule.ValidateGenesis(illegalData)
	require.Error(t, err)

	validatorUpdates := appModule.InitGenesis(ctx, defaultGen)
	require.True(t, len(validatorUpdates) == 0)
	exportedGenesis := appModule.ExportGenesis(ctx)
	require.True(t, exportedGenesis != nil)

	// Begin Block
	appModule.BeginBlock(ctx, abci.RequestBeginBlock{})

	// EndBlock : add data for execute in EndBlock
	tokenPair := GetBuiltInTokenPair()
	withdrawInfo := types.WithdrawInfo{
		Owner:    tokenPair.Owner,
		Deposits: tokenPair.Deposits,
	}
	dexKeeper.SetWithdrawInfo(ctx, withdrawInfo)
	dexKeeper.SetWithdrawCompleteTimeAddress(ctx, ctx.BlockHeader().Time, tokenPair.Owner)

	// fail case : failed to SendCoinsFromModuleToAccount return error
	spKeeper.behaveEvil = true
	appModule.EndBlock(ctx, abci.RequestEndBlock{})

	// successful case : success to SendCoinsFromModuleToAccount which return nil
	spKeeper.behaveEvil = false
	appModule.EndBlock(ctx, abci.RequestEndBlock{})
}

type MockInvariantRegistry struct{}

func (ir MockInvariantRegistry) RegisterRoute(moduleName, route string, invar sdk.Invariant) {}
