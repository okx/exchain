package ammswap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	swapkeeper "github.com/okex/okexchain/x/ammswap/keeper"
	"github.com/okex/okexchain/x/ammswap/types"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)



func NewTestSwapTokenPairWithInitLiquidity(t *testing.T, ctx sdk.Context, k swapkeeper.Keeper,
	baseToken, quoteToken sdk.DecCoin, addrList []sdk.AccAddress) SwapTokenPair {
	handler := NewHandler(k)

	createExchangeMsg := types.NewMsgCreateExchange(baseToken.Denom, quoteToken.Denom, addrList[0])
	result := handler(ctx, createExchangeMsg)
	require.Equal(t, true, result.IsOK())
	deadLine := time.Now().Unix()
	addLiquidityMsg := types.NewMsgAddLiquidity(sdk.NewDec(0), baseToken, quoteToken, deadLine, addrList[0])
	result = handler(ctx, addLiquidityMsg)
	require.Equal(t, true, result.IsOK())
	for _, addr := range addrList {
		baseToken1 := sdk.NewDecCoinFromDec(baseToken.Denom, baseToken.Amount.Mul(sdk.NewDec(100)))
		quoteToken1:= sdk.NewDecCoinFromDec(quoteToken.Denom, quoteToken.Amount.Mul(sdk.NewDec(100)))
		addLiquidityMsg := types.NewMsgAddLiquidity(sdk.NewDec(0), baseToken1, quoteToken1, deadLine, addr)
		result = handler(ctx, addLiquidityMsg)
		require.Equal(t, true, result.IsOK())
	}


	swapTokenPair, err := k.GetSwapTokenPair(ctx, types.GetSwapTokenPairName(baseToken.Denom, quoteToken.Denom))
	require.Nil(t, err)

	return swapTokenPair
}

