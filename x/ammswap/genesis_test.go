package ammswap

import (
	"testing"

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

func TestExportSupplyGenesisWithZeroLiquidity(t *testing.T) {
	// init
	mapp, addrKeysSlice := getMockApp(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))
	// set test tokens
	err := types.SetTestTokens(ctx, mapp.tokenKeeper, mapp.supplyKeeper, addrKeysSlice[0].Address, mapp.TotalCoinsSupply)
	require.NoError(t, err)

	handler := NewHandler(keeper)
	// 1. create 2 new ammswap tokens
	// then add liquidity in 2 swap pairs
	// then remove liquidity totally in 1 swap pair
	for _, msg := range types.CreateTestMsgs(addrKeysSlice[0].Address) {
		_, err = handler(ctx, msg)
		require.NoError(t, err)
	}

	// 2.1 Test supply Invariant
	supplyinvariant := supply.AllInvariants(mapp.supplyKeeper)
	_, broken := supplyinvariant(ctx)
	require.False(t, broken)
	// 2.2 Test supply ExportGenesis: remove coin whose supply is zero
	var expectedCoins sdk.DecCoins
	mapp.AccountKeeper.IterateAccounts(ctx, func(acc exported.Account) bool {
		expectedCoins = expectedCoins.Add(acc.GetCoins()...)
		return false
	})
	supplyExportGenesis := supply.ExportGenesis(ctx, mapp.supplyKeeper)
	require.EqualValues(t, expectedCoins, supplyExportGenesis.Supply)
}