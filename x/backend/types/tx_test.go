package types

import (
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/okex/okchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okchain/x/order"
	orderKeeper "github.com/okex/okchain/x/order/keeper"
	orderTypes "github.com/okex/okchain/x/order/types"
	tokenKeeper "github.com/okex/okchain/x/token"
	token "github.com/okex/okchain/x/token/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestGenerateTx(t *testing.T) {
	txbldr := auth.NewTxBuilder(auth.DefaultTxEncoder(auth.ModuleCdc), 1, 2, 3, 4, false, "okchain", "memo", nil, nil)

	priKeyFrom := secp256k1.GenPrivKey()
	pubKeyFrom := priKeyFrom.PubKey()
	accFrom := sdk.AccAddress(pubKeyFrom.Address())

	priKeyTo := secp256k1.GenPrivKey()
	pubKeyTo := priKeyTo.PubKey()
	accTo := sdk.AccAddress(pubKeyTo.Address())

	// send
	decCoins, err := sdk.ParseDecCoins("100" + common.NativeToken)
	require.Nil(t, err)
	sendMsg := token.NewMsgTokenSend(accFrom, accTo, decCoins)

	sendMsgSig, _ := priKeyFrom.Sign(sendMsg.GetSignBytes())
	sigs := []auth.StdSignature{
		{
			PubKey:    pubKeyFrom,
			Signature: sendMsgSig,
		},
	}
	txSigMsg, _ := txbldr.BuildSignMsg([]sdk.Msg{sendMsg})
	tx := auth.NewStdTx(txSigMsg.Msgs, txSigMsg.Fee, sigs, "")
	ctx0, keeper0, _, _ := tokenKeeper.CreateParam(t, false)
	GenerateTx(&tx, "", ctx0, nil, keeper0, time.Now().Unix())

	// order/new
	orderNewMsg := order.NewMsgNewOrder(accFrom, "btc_"+common.NativeToken, SellOrder, "23.76", "289")
	orderNewMsgSig, _ := priKeyFrom.Sign(orderNewMsg.GetSignBytes())
	sigs = []auth.StdSignature{
		{
			PubKey:    pubKeyFrom,
			Signature: orderNewMsgSig,
		},
	}
	txSigMsg, _ = txbldr.BuildSignMsg([]sdk.Msg{orderNewMsg})
	tx = auth.NewStdTx(txSigMsg.Msgs, txSigMsg.Fee, sigs, "")
	GenerateTx(&tx, "", sdk.Context{}, nil, nil, time.Now().Unix())

	// order/cancel
	orderCancelMsg := order.NewMsgCancelOrder(accFrom, "ORDER-123")
	orderCancelMsgSig, _ := priKeyFrom.Sign(orderCancelMsg.GetSignBytes())
	sigs = []auth.StdSignature{
		{
			PubKey:    pubKeyFrom,
			Signature: orderCancelMsgSig,
		},
	}
	txSigMsg, _ = txbldr.BuildSignMsg([]sdk.Msg{orderCancelMsg})
	tx = auth.NewStdTx(txSigMsg.Msgs, txSigMsg.Fee, sigs, "")

	testInput := orderKeeper.CreateTestInput(t)
	keeper := testInput.OrderKeeper
	ctx := testInput.Ctx.WithBlockHeight(10)
	or := &order.Order{
		OrderID: orderCancelMsg.OrderIDs[0],
		Side:    SellOrder,
	}
	keeper.SetOrder(ctx, or.OrderID, or)
	or.SetExtraInfoWithKeyValue(orderTypes.OrderExtraInfoKeyCancelFee, "1"+common.NativeToken)
	GenerateTx(&tx, "", ctx, keeper, nil, time.Now().Unix())
}

func TestTicker(t *testing.T) {
	tiker1 := Ticker{
		Symbol:           "btc",
		Product:          "btc_" + common.NativeToken,
		Timestamp:        0,
		Open:             10.5,
		Close:            53.5,
		High:             100,
		Low:              6.66,
		Price:            2.46,
		Volume:           3000,
		Change:           43,
		ChangePercentage: "409.52%",
	}
	tiker2 := Ticker{
		Symbol:           "eth",
		Product:          "eth_" + common.NativeToken,
		Timestamp:        0,
		Open:             3.8,
		Close:            15.9,
		High:             200,
		Low:              2,
		Price:            9.6,
		Volume:           110,
		Change:           12.1,
		ChangePercentage: "318.42%",
	}

	tikerStr := tiker1.PrettyString()
	str := fmt.Sprintf("[Ticker] Symbol: %s, Price: %f, TStr: %s, Timestamp: %d, OCHLV(%f, %f, %f, %f, %f) [%f, %s])",
		tiker1.Symbol, tiker1.Price, TimeString(tiker1.Timestamp), tiker1.Timestamp, tiker1.Open, tiker1.Close, tiker1.High, tiker1.Low, tiker1.Volume, tiker1.Change, tiker1.ChangePercentage)

	require.Equal(t, str, tikerStr)

	tikers := Tickers{tiker1, tiker2}
	sort.Sort(tikers)
	require.Equal(t, tiker2.Symbol, tikers[0].Symbol)
	require.Equal(t, tiker1.Symbol, tikers[1].Symbol)
}
