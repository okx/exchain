package staking

import (
	"bufio"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"

	"github.com/okex/exchain/x/staking/keeper"
	"github.com/okex/exchain/x/staking/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	cliLcd "github.com/cosmos/cosmos-sdk/client/lcd"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
)

// getMockApp returns an initialized mock application for this module.
func getMockApp(t *testing.T) (*mock.App, keeper.MockStakingKeeper) {
	mApp := mock.NewApp()

	//RegisterCodec(mApp.Cdc)

	_, accKeeper, mKeeper := CreateTestInput(t, false, SufficientInitPower)
	keeper := mKeeper.Keeper

	mApp.Router().AddRoute(RouterKey, NewHandler(keeper))
	mApp.SetEndBlocker(getEndBlocker(keeper))
	mApp.SetInitChainer(getInitChainer(mApp, keeper, accKeeper, mKeeper.SupplyKeeper))

	require.NoError(t, mApp.CompleteSetup(mKeeper.StoreKey, mKeeper.TkeyStoreKey))
	return mApp, mKeeper
}

// getEndBlocker returns a staking endblocker.
func getEndBlocker(keeper Keeper) sdk.EndBlocker {
	return func(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
		validatorUpdates := EndBlocker(ctx, keeper)

		return abci.ResponseEndBlock{
			ValidatorUpdates: validatorUpdates,
		}
	}
}

// getInitChainer initializes the chainer of the mock app and sets the genesis
// state. It returns an empty ResponseInitChain.
func getInitChainer(mapp *mock.App, keeper Keeper, accKeeper types.AccountKeeper,
	supKeeper types.SupplyKeeper) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
		mapp.InitChainer(ctx, req)

		stakingGenesis := DefaultGenesisState()
		validators := InitGenesis(ctx, keeper, accKeeper, supKeeper, stakingGenesis)

		return abci.ResponseInitChain{
			Validators: validators,
		}
	}
}

type MockInvariantRegistry struct{}

func (ir MockInvariantRegistry) RegisterRoute(moduleName, route string, invar sdk.Invariant) {}

//__________________________________________________________________________________________

func TestAppSmoke(t *testing.T) {
	mApp, mKeeper := getMockApp(t)
	appModule := NewAppModule(mKeeper.Keeper, mKeeper.AccKeeper, mKeeper.SupplyKeeper)

	// Const Info
	require.True(t, appModule.Name() == ModuleName)
	require.True(t, appModule.Route() == RouterKey)
	require.True(t, appModule.QuerierRoute() == QuerierRoute)
	require.True(t, appModule.GetQueryCmd(mApp.Cdc) != nil)
	require.True(t, appModule.GetTxCmd(mApp.Cdc) != nil)

	appModule.RegisterCodec(mApp.Cdc)
	appModule.RegisterInvariants(MockInvariantRegistry{})
	rs := cliLcd.NewRestServer(mApp.Cdc, nil)
	appModule.RegisterRESTRoutes(rs.CliCtx, rs.Mux)
	handler := appModule.NewHandler()
	require.True(t, handler != nil)
	querior := appModule.NewQuerierHandler()
	require.True(t, querior != nil)

	// Extra Helper
	appModule.CreateValidatorMsgHelpers("0.0.0.0")
	cliCtx := context.NewCLIContext().WithCodec(mApp.Cdc)
	inBuf := bufio.NewReader(os.Stdin)
	txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(mApp.Cdc))
	appModule.BuildCreateValidatorMsg(cliCtx, txBldr)

	// Initialization for genesis
	defaultGen := appModule.DefaultGenesis()
	err := appModule.ValidateGenesis(defaultGen)
	require.True(t, err == nil)

	illegalData := []byte{}
	err = appModule.ValidateGenesis(illegalData)
	require.Error(t, err)

	// Basic abci test
	header := abci.Header{ChainID: keeper.TestChainID, Height: 0}
	ctx := sdk.NewContext(mKeeper.MountedStore, header, false, log.NewNopLogger())
	validatorUpdates := appModule.InitGenesis(ctx, defaultGen)
	require.True(t, len(validatorUpdates) == 0)
	exportedGenesis := appModule.ExportGenesis(ctx)
	require.True(t, exportedGenesis != nil)

	// Begin & End Block
	appModule.BeginBlock(ctx, abci.RequestBeginBlock{})
	appModule.EndBlock(ctx, abci.RequestEndBlock{})

}
