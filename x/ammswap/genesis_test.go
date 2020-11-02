package ammswap

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	mapp.supplyKeeper.SetTokensSupply(ctx, mapp.TotalCoinsSupply)

	defaultGenesisState := DefaultGenesisState()
	testSwapTokenPair := types.GetTestSwapTokenPair()
	defaultGenesisState.SwapTokenPairRecords = []SwapTokenPair{
		testSwapTokenPair,
	}
	InitGenesis(ctx, keeper, defaultGenesisState)
	exportedGenesis := ExportGenesis(ctx, keeper)
	require.Equal(t, defaultGenesisState, exportedGenesis)

}
