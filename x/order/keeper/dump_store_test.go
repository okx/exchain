package keeper

import (
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/dex"
	"github.com/okex/exchain/x/order/types"
	"github.com/stretchr/testify/require"

	"testing"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

func TestDumpStore(t *testing.T) {
	common.InitConfig()
	testInput := CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx
	querier := NewQuerier(keeper)

	tokenPair := dex.GetBuiltInTokenPair()
	err := testInput.DexKeeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	product := types.TestTokenPair
	orders := []*types.Order{
		mockOrder("", product, types.SellOrder, "0.6", "1.1"),
		mockOrder("", product, types.SellOrder, "0.5", "1.2"),
		mockOrder("", product, types.BuyOrder, "0.4", "1.3"),
		mockOrder("", product, types.BuyOrder, "0.3", "1.4"),
	}
	for i := 0; i < 4; i++ {
		orders[i].Sender = testInput.TestAddrs[0]
		err := keeper.PlaceOrder(ctx, orders[i])
		require.NoError(t, err)
	}
	keeper.Cache2Disk(ctx)

	// Default query
	path := []string{types.QueryStore}
	bz, err := querier(ctx, path, abci.RequestQuery{})
	require.Nil(t, err)
	ss := &StoreStatistic{}
	keeper.cdc.MustUnmarshalJSON(bz, ss)
	require.EqualValues(t, 4, ss.StoreOrderNum)
	require.EqualValues(t, 4, ss.DepthBookNum[product])

	// Test DumpStore
	keeper.SetLastPrice(ctx, product, sdk.MustNewDecFromStr("1.5"))
	keeper.DumpStore(ctx)
}
