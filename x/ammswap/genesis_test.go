package ammswap

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestValidateGenesis(t *testing.T) {
	defaultGenesisState := DefaultGenesisState()
	testSwapTokenPair := types.GetTestSwapTokenPair()
	defaultGenesisState.SwapTokenPairRecords = []SwapTokenPair{
		testSwapTokenPair,
	}
	err := ValidateGenesis(defaultGenesisState)
	require.Nil(t, err)

	invalidBaseAmount := sdk.NewDecCoinFromDec("bsa", sdk.NewDec(10000))
	invalidBaseAmount.Denom = "1add"
	invalidQuoteAmount := sdk.NewDecCoinFromDec("bsa", sdk.NewDec(10000))
	invalidQuoteAmount.Denom = "1dfdf"
	invalidPoolTokenName := "abc"

	defaultGenesisState.SwapTokenPairRecords = []SwapTokenPair{
		*types.NewSwapTokenPair(invalidBaseAmount, testSwapTokenPair.QuotePooledCoin, testSwapTokenPair.PoolTokenName),
	}
	err = ValidateGenesis(defaultGenesisState)
	require.NotNil(t, err)

	defaultGenesisState.SwapTokenPairRecords = []SwapTokenPair{
		*types.NewSwapTokenPair(testSwapTokenPair.BasePooledCoin, invalidQuoteAmount, testSwapTokenPair.PoolTokenName),
	}
	err = ValidateGenesis(defaultGenesisState)
	require.NotNil(t, err)

	defaultGenesisState.SwapTokenPairRecords = []SwapTokenPair{
		*types.NewSwapTokenPair(testSwapTokenPair.BasePooledCoin, testSwapTokenPair.QuotePooledCoin, invalidPoolTokenName),
	}
	err = ValidateGenesis(defaultGenesisState)
	require.NotNil(t, err)
}

func TestInitAndExportGenesis(t *testing.T) {
	mapp, _ := getMockApp(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	defaultGenesisState := DefaultGenesisState()
	testSwapTokenPair := types.GetTestSwapTokenPair()
	defaultGenesisState.SwapTokenPairRecords = []SwapTokenPair{
		testSwapTokenPair,
	}
	InitGenesis(ctx, keeper, defaultGenesisState)
	exportedGenesis := ExportGenesis(ctx, keeper)
	require.Equal(t, defaultGenesisState, exportedGenesis)

}

func TestInitAndExportGenesisWithZeroLiquidity(t *testing.T) {
	mapp, addrKeysSlice := getMockApp(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	err := types.SetTokens(ctx, mapp.tokenKeeper, mapp.supplyKeeper, addrKeysSlice[0].Address)
	require.NoError(t, err)

	// test ammswap InitGenesis: init 3 new ammswap tokens
	defaultGenesisState := DefaultGenesisState()
	defaultGenesisState.SwapTokenPairRecords = []SwapTokenPair{
		types.GetTestSwapTokenPair(), types.GetTestSwapTokenPairWithLargeLiquidity(), types.GetTestSwapTokenPairWithZeroLiquidity(),
	}
	InitGenesis(ctx, keeper, defaultGenesisState)
	swapTokenPairs := keeper.GetSwapTokenPairs(ctx)
	require.Equal(t, 2, len(swapTokenPairs))
	require.EqualValues(t, defaultGenesisState.SwapTokenPairRecords[:2], swapTokenPairs)

	// test ammswap ExportGenesis: create 2 new ammswap tokens
	handler := NewHandler(keeper)
	_, err = handler(ctx, types.GetCreateExchangeMsg4(addrKeysSlice[0].Address))
	require.NoError(t, err)
	_, err = handler(ctx, types.GetCreateExchangeMsg5(addrKeysSlice[0].Address))
	require.NoError(t, err)
	_, err = handler(ctx, types.NewMsgAddLiquidity(sdk.ZeroDec(),
		sdk.NewDecCoin(types.TestBasePooledToken4, sdk.OneInt()), sdk.NewDecCoin(types.TestQuotePooledToken, sdk.OneInt()),
		time.Now().Add(time.Hour).Unix(), addrKeysSlice[0].Address))
	require.NoError(t, err)

	exportedGenesis := ExportGenesis(ctx, keeper)
	require.EqualValues(t, defaultGenesisState.Params, exportedGenesis.Params)
	require.Equal(t, 3, len(exportedGenesis.SwapTokenPairRecords))
	expectedSwapTokenPairRecords := []SwapTokenPair{
		types.GetTestSwapTokenPair(),
		types.GetTestSwapTokenPairWithLargeLiquidity(),
		*types.NewSwapTokenPair(
			sdk.NewDecCoin(types.TestQuotePooledToken, sdk.OneInt()),
			sdk.NewDecCoin(types.TestBasePooledToken4, sdk.OneInt()),
			types.GetPoolTokenName(types.TestQuotePooledToken, types.TestBasePooledToken4)),
	}
	require.EqualValues(t, expectedSwapTokenPairRecords, exportedGenesis.SwapTokenPairRecords)

	// test supply Invariant & ExportGenesis
	supplyinvariant := supply.AllInvariants(mapp.supplyKeeper)
	_, broken := supplyinvariant(ctx)
	require.False(t, broken)
	var expectedCoins sdk.DecCoins
	mapp.AccountKeeper.IterateAccounts(ctx, func(acc exported.Account) bool {
		expectedCoins = expectedCoins.Add(acc.GetCoins()...)
		return false
	})
	supplyExportGenesis := supply.ExportGenesis(ctx, mapp.supplyKeeper)
	require.EqualValues(t, expectedCoins, supplyExportGenesis.Supply )
}