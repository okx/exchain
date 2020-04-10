package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/okex/okchain/x/common"
	orderTypes "github.com/okex/okchain/x/order/types"
	tokenTypes "github.com/okex/okchain/x/token/types"
)

func GenerateTx(tx *auth.StdTx, txHash string, ctx sdk.Context, orderKeeper OrderKeeper, tokenKeeper TokenKeeper,
	timestamp int64) []*Transaction {
	txs := make([]*Transaction, 0, 2)
	for _, msg := range tx.GetMsgs() {
		switch msg.Type() {
		case "send": // token/send
			txFrom, txTo := buildTransactionsTransfer(msg.(tokenTypes.MsgSend), txHash, ctx, tokenKeeper,
				timestamp)
			txs = append(txs, txFrom, txTo)
		case "new": // order/new
			transaction := buildTransactionNew(msg.(orderTypes.MsgNewOrders), txHash, ctx, timestamp)
			txs = append(txs, transaction)
		case "cancel": // order/cancel
			transaction := buildTransactionCancel(msg.(orderTypes.MsgCancelOrders), txHash, ctx, orderKeeper,
				timestamp)
			if transaction != nil {
				txs = append(txs, transaction)
			}
		default: // In other cases, do nothing
			continue
		}
	}
	return txs
}

func buildTransactionsTransfer(msg tokenTypes.MsgSend, txHash string, ctx sdk.Context, tokenKeeper TokenKeeper,
	timestamp int64) (*Transaction, *Transaction) {
	decCoins := msg.Amount

	txFrom := &Transaction{
		TxHash:    txHash,
		Address:   msg.FromAddress.String(),
		Type:      TxTypeTransfer,
		Side:      TxSideFrom,
		Symbol:    decCoins[0].Denom,
		Quantity:  decCoins[0].Amount.String(),
		Fee:       tokenKeeper.GetParams(ctx).FeeBase.String(), // TODO: get fee from params
		Timestamp: timestamp,
	}
	txTo := &Transaction{
		TxHash:    txHash,
		Address:   msg.ToAddress.String(),
		Type:      TxTypeTransfer,
		Side:      TxSideTo,
		Symbol:    decCoins[0].Denom,
		Quantity:  decCoins[0].Amount.String(),
		Fee:       sdk.DecCoin{Denom: common.NativeToken, Amount: sdk.ZeroDec()}.String(),
		Timestamp: timestamp,
	}
	return txFrom, txTo
}

func buildTransactionNew(msg orderTypes.MsgNewOrders, txHash string, ctx sdk.Context, timestamp int64) *Transaction {
	side := TxSideBuy
	if msg.OrderItems[0].Side == orderTypes.SellOrder {
		side = TxSideSell
	}
	return &Transaction{
		TxHash:    txHash,
		Address:   msg.Sender.String(),
		Type:      TxTypeOrderNew,
		Side:      int64(side),
		Symbol:    msg.OrderItems[0].Product,
		Quantity:  msg.OrderItems[0].Quantity.String(),
		Fee:       sdk.DecCoin{Denom: common.NativeToken, Amount: sdk.ZeroDec()}.String(), // TODO: get fee from params
		Timestamp: timestamp,
	}
}

func buildTransactionCancel(msg orderTypes.MsgCancelOrders, txHash string, ctx sdk.Context, orderKeeper OrderKeeper, timestamp int64) *Transaction {
	order := orderKeeper.GetOrder(ctx, msg.OrderIDs[0])
	if order == nil {
		return nil
	}
	side := TxSideBuy
	if order.Side == orderTypes.SellOrder {
		side = TxSideSell
	}
	cancelFeeStr := order.GetExtraInfoWithKey(orderTypes.OrderExtraInfoKeyCancelFee)
	if cancelFeeStr == "" {
		cancelFeeStr = sdk.DecCoin{Denom: common.NativeToken, Amount: sdk.ZeroDec()}.String()
	}
	return &Transaction{
		TxHash:    txHash,
		Address:   msg.Sender.String(),
		Type:      TxTypeOrderCancel,
		Side:      int64(side),
		Symbol:    order.Product,
		Quantity:  order.Quantity.String(),
		Fee:       cancelFeeStr,
		Timestamp: timestamp,
	}
}
