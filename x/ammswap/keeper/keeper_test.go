package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/okex/okexchain/x/common"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestKeeper_GetPoolTokenInfo(t *testing.T) {
	addrTest := "okexchain1a20d4xmqj4m9shtm0skt0aaahsgeu4h6746fs2"
	mapp, _ := GetTestInput(t, 1)
	keeper := mapp.swapKeeper
	mapp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mapp.BaseApp.NewContext(false, abci.Header{}).WithBlockHeight(10)
	mapp.supplyKeeper.SetSupply(ctx, supply.NewSupply(mapp.TotalCoinsSupply))

	// init a pool token
	symbol := types.PoolTokenPrefix + common.TestToken
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
