package ammswap

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/token"
	tokentypes "github.com/okex/okexchain/x/token/types"
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
	// init
	mapp, addrKeysSlice, keeper, tokenKeeper, supplyKeeper := getMockAppWithKeeper(t, 1)
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	// set test tokens
	err := types.SetTestTokens(ctx, tokenKeeper, supplyKeeper, addrKeysSlice[0].Address)
	require.NoError(t, err)

	// 1. test ammswap InitGenesis:
	// 1.1 add 3 new ammswap tokens in genesis
	defaultGenesisState := DefaultGenesisState()
	defaultGenesisState.SwapTokenPairRecords = []SwapTokenPair{
		types.GetTestSwapTokenPair(), types.GetTestSwapTokenPairWithLargeLiquidity(), types.GetTestSwapTokenPairWithZeroLiquidity(),
	}
	// 1.2 ammswap InitGenesis: should remove 1 swap pair whose value is zero
	InitGenesis(ctx, keeper, defaultGenesisState)
	swapTokenPairs := keeper.GetSwapTokenPairs(ctx)
	require.EqualValues(t, defaultGenesisState.SwapTokenPairRecords[:2], swapTokenPairs)

	// 2. test ammswap ExportGenesis:
	handler := NewHandler(keeper)
	// 2.1 create 2 new ammswap tokens
	// then add liquidity in 2 swap pairs
	// then remove liquidity in 1 swap pair
	for _, msg := range types.CreateTestMsgs(addrKeysSlice[0].Address) {
		_, err = handler(ctx, msg)
		require.NoError(t, err)
	}
	// 2.2 ammswap ExportGenesis: should remove 1 swap pair whose value is zero
	exportedGenesis := ExportGenesis(ctx, keeper)
	require.EqualValues(t, defaultGenesisState.Params, exportedGenesis.Params)
	expectedSwapTokenPairRecords := []SwapTokenPair{
		types.GetTestSwapTokenPair(),
		types.GetTestSwapTokenPairWithLargeLiquidity(),
		*types.NewSwapTokenPair(
			sdk.NewDecCoin(types.TestQuotePooledToken, sdk.OneInt()),
			sdk.NewDecCoin(types.TestBasePooledToken4, sdk.OneInt()),
			types.GetPoolTokenName(types.TestQuotePooledToken, types.TestBasePooledToken4)),
	}
	require.EqualValues(t, expectedSwapTokenPairRecords, exportedGenesis.SwapTokenPairRecords)

	// 3.1 test supply Invariant
	supplyinvariant := supply.AllInvariants(supplyKeeper)
	_, broken := supplyinvariant(ctx)
	require.False(t, broken)
	// 3.2 test supply ExportGenesis: remove coin whose supply is zero
	var expectedCoins sdk.DecCoins
	mapp.AccountKeeper.IterateAccounts(ctx, func(acc exported.Account) bool {
		expectedCoins = expectedCoins.Add(acc.GetCoins()...)
		return false
	})
	supplyExportGenesis := supply.ExportGenesis(ctx, supplyKeeper)
	require.EqualValues(t, expectedCoins, supplyExportGenesis.Supply)

	// 4.1 test token ExportGenesis: remove coin whose TotalSupply is zer
	tokenKeeper.SetParams(ctx, tokentypes.DefaultParams())
	tokenExportGenesis := token.ExportGenesis(ctx, tokenKeeper)
	fmt.Println(len(tokenExportGenesis.Tokens))
}