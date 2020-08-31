package pushservice

import (
	"fmt"
	"testing"

	"github.com/okex/okchain/x/stream/pushservice/conn"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/stream/common"
	"github.com/okex/okchain/x/stream/pushservice/channels"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/okex/okchain/x/order"
	"github.com/okex/okchain/x/stream/pushservice/types"
	"github.com/okex/okchain/x/token"
	"github.com/stretchr/testify/require"
)

const (
	IP_PORT      = "localhost:16379"
	PASSWD       = ""
	DB           = 0
	TestRedisUrl = "redis://127.0.0.1:16379"
)

func getPushService(ctx sdk.Context, t *testing.T) *PushService {
	srv, err := NewPushService(IP_PORT, PASSWD, DB, ctx.Logger())
	require.Nil(t, err)
	return srv
}

//func TestService(t *testing.T) {
//	srv, err := NewPushService("1.1.1.1:6371", PASSWD, DB, nil)
//	require.Nil(t, srv)
//	require.Error(t, err)
//
//	app, _ := common.GetMockApp(t, 2)
//	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
//	ctx := app.NewContext(false, abci.Header{})
//	srv, err = NewPushService("1.1.1.1:6371", PASSWD, DB, ctx.Logger())
//	require.Nil(t, srv)
//	require.Error(t, err)
//}

func TestSetData(t *testing.T) {
	app, addrKeysSlice := common.GetMockApp(t, 2)

	pool, err := common.NewPool(TestRedisUrl, "", app.Logger())
	require.Nil(t, err)
	pool.Get().Do("FLUSHALL")

	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := app.NewContext(false, abci.Header{})

	//init push_service
	srv := getPushService(ctx, t)
	pBlock := types.NewRedisBlock()
	require.True(t, pBlock.Empty())

	//push
	pBlock.SetData(ctx, app.OrderKeeper, app.TokenKeeper, app.AccountKeeper)
	require.Equal(t, 0, len(pBlock.MatchesMap))
	require.Equal(t, 4, len(pBlock.AccountsMap)) //addr:coin
	require.Equal(t, 0, len(pBlock.Instruments))
	srv.WriteSync(pBlock)

	tokenPair := token.GetBuiltInTokenPair()

	app.TokenKeeper.SaveTokenPair(ctx, tokenPair)

	quantity := "1"
	price := "0.1"

	orderMsg0 := order.NewMsgNewOrder(nil, types.TestTokenPair, types.BuyOrder, price, quantity)
	ctx = app.NewContext(true, abci.Header{Height: 2})
	common.MockApplyBlock(app, int64(2), common.ProduceOrderTxs(app, ctx, 10, addrKeysSlice[0], &orderMsg0))
	//push
	pBlock.SetData(ctx, app.OrderKeeper, app.TokenKeeper, app.AccountKeeper)
	require.Equal(t, 0, len(pBlock.MatchesMap))
	require.Equal(t, 2, len(pBlock.AccountsMap)) //addr:coin
	require.Equal(t, 0, len(pBlock.Instruments)) //len(Instruments) is 0 because of unittest don't set the stream_engine config.
	srv.WriteSync(pBlock)

	ctx = app.NewContext(true, abci.Header{Height: 3})
	orderMsg1 := order.NewMsgNewOrder(nil, types.TestTokenPair, types.SellOrder, price, quantity)
	common.MockApplyBlock(app, int64(3), common.ProduceOrderTxs(app, ctx, 10, addrKeysSlice[1], &orderMsg1))
	//push
	pBlock.SetData(ctx, app.OrderKeeper, app.TokenKeeper, app.AccountKeeper)
	require.Equal(t, 1, len(pBlock.MatchesMap))
	require.Equal(t, 5, len(pBlock.AccountsMap))
	require.Equal(t, 0, len(pBlock.Instruments))
	srv.WriteSync(pBlock)

	ctx = app.NewContext(true, abci.Header{Height: 4})
	pBlock.SetData(ctx, app.OrderKeeper, app.TokenKeeper, app.AccountKeeper)
	require.Equal(t, 2, len(pBlock.OrdersMap))
	require.Equal(t, 0, len(pBlock.AccountsMap))
	require.Equal(t, 0, len(pBlock.Instruments))
	require.Equal(t, int64(4), pBlock.Height)

	require.Equal(t, 0, len(pBlock.MatchesMap))
	srv.WriteSync(pBlock)
	require.True(t, pBlock.Empty())
	require.NoError(t, srv.Close())
}

func TestSetDataErr(t *testing.T) {
	app, addrKeysSlice := common.GetMockApp(t, 2)

	pool, _ := common.NewPool(TestRedisUrl, "", app.Logger())
	pool.Get().Do("FLUSHALL")

	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := app.NewContext(false, abci.Header{})

	//init push service
	srv := getPushService(ctx, t)
	pBlock := types.NewRedisBlock()
	srv.Close() //close redis connect, to test set err

	tokenPair := token.GetBuiltInTokenPair()
	app.TokenKeeper.SaveTokenPair(ctx, tokenPair)

	quantity := "1"
	price := "0.1"

	//orders, account
	orderMsg0 := order.NewMsgNewOrder(nil, types.TestTokenPair, types.BuyOrder, price, quantity)
	ctx = app.NewContext(true, abci.Header{Height: 2})
	common.MockApplyBlock(app, int64(2), common.ProduceOrderTxs(app, ctx, 10, addrKeysSlice[0], &orderMsg0))
	//push
	pBlock.SetData(ctx, app.OrderKeeper, app.TokenKeeper, app.AccountKeeper)
	_, err := srv.WriteSync(pBlock)
	require.Error(t, err)

	//orders, depth, matches
	orderMsg0 = order.NewMsgNewOrder(nil, types.TestTokenPair, types.BuyOrder, price, quantity)
	ctx = app.NewContext(true, abci.Header{Height: 3})
	common.MockApplyBlock(app, int64(3), common.ProduceOrderTxs(app, ctx, 10, addrKeysSlice[0], &orderMsg0))

	ctx = app.NewContext(true, abci.Header{Height: 4})
	orderMsg1 := order.NewMsgNewOrder(nil, types.TestTokenPair, types.SellOrder, price, quantity)
	common.MockApplyBlock(app, int64(4), common.ProduceOrderTxs(app, ctx, 10, addrKeysSlice[1], &orderMsg1))
	//push
	pBlock.SetData(ctx, app.OrderKeeper, app.TokenKeeper, app.AccountKeeper)
	_, err = srv.WriteSync(pBlock)
	require.Error(t, err)

	ctx = app.NewContext(true, abci.Header{Height: 5})
	pBlock.SetData(ctx, app.OrderKeeper, app.TokenKeeper, app.AccountKeeper)
	_, err = srv.WriteSync(pBlock)
	require.Error(t, err)
}

func TestMGet(t *testing.T) {
	app, _ := common.GetMockApp(t, 2)
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := app.NewContext(false, abci.Header{})

	c, err := conn.NewClient(IP_PORT, PASSWD, DB, ctx.Logger())
	require.NoError(t, err)

	vals, err := c.MGet([]string{channels.GetCSpotMetaKey(), channels.GetCSpotDepthKey(types.TestTokenPair)})
	require.NoError(t, err)
	require.Equal(t, 2, len(vals))

	fmt.Println(vals)
}
