package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

const addrTest = "okexchain1a20d4xmqj4m9shtm0skt0aaahsgeu4h6746fs2"

func TestKeeper_GetPoolTokenInfo(t *testing.T) {
	mapp, _ := GetTestInput(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetTokensSupply(ctx, mapp.TotalCoinsSupply)

	// init a pool token
	symbol := types.GetPoolTokenName(types.TestBasePooledToken, types.TestQuotePooledToken)
	keeper.NewPoolToken(ctx, symbol)
	poolToken, err := keeper.GetPoolTokenInfo(ctx, symbol)
	require.Nil(t, err)
	require.EqualValues(t, symbol, poolToken.WholeName)

	// pool token is Interest token
	require.EqualValues(t, types.GenerateTokenType, poolToken.Type)

	// check pool token total supply
	amount := keeper.GetPoolTokenAmount(ctx, symbol)
	require.EqualValues(t, sdk.MustNewDecFromStr("0"), amount)

	mintToken := sdk.NewDecCoinFromDec(symbol, sdk.NewDec(1000000))
	err = keeper.MintPoolCoinsToUser(ctx, sdk.DecCoins{mintToken}, sdk.AccAddress(addrTest))
	require.Nil(t, err)

	balance := mapp.bankKeeper.GetCoins(ctx, sdk.AccAddress(addrTest))
	require.NotNil(t, balance)
}

func TestKeeper_GetSwapTokenPairs(t *testing.T) {
	mapp, _ := GetTestInput(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetTokensSupply(ctx, mapp.TotalCoinsSupply)

	swapTokenPair := types.GetTestSwapTokenPair()
	keeper.SetSwapTokenPair(ctx, types.TestSwapTokenPairName, swapTokenPair)

	expectedSwapTokenPairList := []types.SwapTokenPair{swapTokenPair}
	swapTokenPairList := keeper.GetSwapTokenPairs(ctx)
	require.Equal(t, expectedSwapTokenPairList, swapTokenPairList)
}

func TestKeeper_GetRedeemableAssets(t *testing.T) {
	mapp, _ := GetTestInput(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetTokensSupply(ctx, mapp.TotalCoinsSupply)

	swapTokenPair := types.GetTestSwapTokenPair()
	tokenNumber := sdk.NewDec(100)
	swapTokenPair.QuotePooledCoin.Amount = tokenNumber
	swapTokenPair.BasePooledCoin.Amount = tokenNumber
	keeper.SetSwapTokenPair(ctx, types.TestSwapTokenPairName, swapTokenPair)
	poolToken := types.InitPoolToken(swapTokenPair.PoolTokenName)
	initPoolTokenAmount := sdk.NewDecCoinFromDec(swapTokenPair.PoolTokenName, sdk.NewDec(1))
	err := keeper.MintPoolCoinsToUser(ctx, sdk.DecCoins{initPoolTokenAmount}, sdk.AccAddress(addrTest))
	require.Nil(t, err)
	mapp.tokenKeeper.NewToken(ctx, poolToken)

	expectedBaseAmount, expectedQuoteAmount := swapTokenPair.BasePooledCoin, swapTokenPair.QuotePooledCoin
	baseAmount, quoteAmount, err := keeper.GetRedeemableAssets(ctx, swapTokenPair.BasePooledCoin.Denom, swapTokenPair.QuotePooledCoin.Denom, initPoolTokenAmount.Amount)
	require.Equal(t, expectedBaseAmount, baseAmount)
	require.Equal(t, expectedQuoteAmount, quoteAmount)
}

func TestGetInputPrice(t *testing.T) {
	inputAmount := sdk.NewDecWithPrec(1, 8)
	inputReserve := sdk.NewDec(1)
	outputReserve := sdk.NewDec(1)
	feeRate := sdk.NewDecWithPrec(3, 3)
	outputAmount := GetInputPrice(inputAmount, inputReserve, outputReserve, feeRate)
	expectedAmount := sdk.NewDec(0)
	require.Equal(t, expectedAmount.String(), outputAmount.String())
}
