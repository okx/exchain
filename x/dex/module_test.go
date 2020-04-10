package dex

import (
	cliLcd "github.com/cosmos/cosmos-sdk/client/lcd"

	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okchain/x/common/version"
	"github.com/okex/okchain/x/dex/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAppModule_Smoke(t *testing.T) {
	_, _, spKeeper, dexKeeper, ctx := getMockTestCaseEvn(t)

	//func NewAppModule(version ProtocolVersionType, keeper Keeper, supplyKeeper SupplyKeeper) AppModule {
	appModule := NewAppModule(version.CurrentProtocolVersion, dexKeeper, spKeeper)

	// Const Info
	require.True(t, appModule.Name() == ModuleName)
	require.True(t, appModule.Route() == RouterKey)
	require.True(t, appModule.QuerierRoute() == QuerierRoute)
	require.True(t, appModule.GetQueryCmd(dexKeeper.GetCDC()) != nil)
	require.True(t, appModule.GetTxCmd(dexKeeper.GetCDC()) != nil)

	// RegisterCodec
	appModule.RegisterCodec(codec.New())

	appModule.RegisterInvariants(nil)
	rs := cliLcd.NewRestServer(dexKeeper.GetCDC(), nil)
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
