package types

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/dependence/cosmos-sdk/x/auth"
	"github.com/okex/exchain/x/common"
	orderTypes "github.com/okex/exchain/x/order/types"
	tokenTypes "github.com/okex/exchain/x/token/types"
	"github.com/willf/bitset"
)

// GenerateTx return transaction, called at DeliverTx
func GenerateTx(tx *auth.StdTx, txHash string, ctx sdk.Context, orderKeeper OrderKeeper, timestamp int64) []*Transaction {
	orderHandlerTxResult := orderKeeper.GetTxHandlerMsgResult()
	idx := int(0)
	var txs []*Transaction

	for _, msg := range tx.GetMsgs() {
		switch msg.Type() {
		case "send": // token/send
			if sendMsg, ok := msg.(tokenTypes.MsgSend); ok {
				txFrom, txTo := buildTransactionsTransfer(tx, sendMsg, txHash, timestamp)
				txs = append(txs, txFrom, txTo)
			}
		case "new": // order/new
			if orderMsg, ok := msg.(orderTypes.MsgNewOrders); ok {
				transaction := buildTransactionNew(orderHandlerTxResult[idx], orderMsg,
					txHash, ctx, timestamp)
				txs = append(txs, transaction...)
				idx++
			}
		case "cancel": // order/cancel
			if cancelMsg, ok := msg.(orderTypes.MsgCancelOrders); ok {
				transaction := buildTransactionCancel(orderHandlerTxResult[idx], cancelMsg,
					txHash, ctx, orderKeeper, timestamp)
				txs = append(txs, transaction...)
				idx++
			}
		default: // In other cases, do nothing
			continue
		}
	}
	return txs
}

func buildTransactionsTransfer(tx *auth.StdTx, msg tokenTypes.MsgSend, txHash string, timestamp int64) (*Transaction, *Transaction) {
	decCoins := msg.Amount

	txFrom := &Transaction{
		TxHash:    txHash,
		Address:   msg.FromAddress.String(),
		Type:      TxTypeTransfer,
		Side:      TxSideFrom,
		Symbol:    decCoins[0].Denom,
		Quantity:  decCoins[0].Amount.String(),
		Fee:       tx.Fee.Amount.String(),
		Timestamp: timestamp,
	}
	txTo := &Transaction{
		TxHash:    txHash,
		Address:   msg.ToAddress.String(),
		Type:      TxTypeTransfer,
		Side:      TxSideTo,
		Symbol:    decCoins[0].Denom,
		Quantity:  decCoins[0].Amount.String(),
		Fee:       sdk.NewDecCoin(common.NativeToken, sdk.ZeroInt()).String(),
		Timestamp: timestamp,
	}
	return txFrom, txTo
}

func buildTransactionNew(handlerMsgResult bitset.BitSet, msg orderTypes.MsgNewOrders, txHash string, ctx sdk.Context, timestamp int64) []*Transaction {
	var result []*Transaction

	for idx, item := range msg.OrderItems {
		if !handlerMsgResult.Test(uint(idx)) {
			continue
		}

		side := TxSideBuy
		if item.Side == orderTypes.SellOrder {
			side = TxSideSell
		}

		tx := Transaction{
			TxHash:    txHash,
			Address:   msg.Sender.String(),
			Type:      TxTypeOrderNew,
			Side:      int64(side),
			Symbol:    item.Product,
			Quantity:  item.Quantity.String(),
			Fee:       sdk.NewDecCoin(common.NativeToken, sdk.ZeroInt()).String(),
			Timestamp: timestamp,
		}

		result = append(result, &tx)
	}

	return result
}

func buildTransactionCancel(handlerMsgResult bitset.BitSet, msg orderTypes.MsgCancelOrders, txHash string, ctx sdk.Context, orderKeeper OrderKeeper, timestamp int64) []*Transaction {
	var result []*Transaction

	for idx, orderID := range msg.OrderIDs {
		if !handlerMsgResult.Test(uint(idx)) {
			continue
		}

		order := orderKeeper.GetOrder(ctx, orderID)
		if order == nil {
			continue
		}
		side := TxSideBuy
		if order.Side == orderTypes.SellOrder {
			side = TxSideSell
		}
		cancelFeeStr := order.GetExtraInfoWithKey(orderTypes.OrderExtraInfoKeyCancelFee)
		if cancelFeeStr == "" {
			cancelFeeStr = sdk.NewDecCoin(common.NativeToken, sdk.ZeroInt()).String()
		}
		tx := Transaction{
			TxHash:    txHash,
			Address:   msg.Sender.String(),
			Type:      TxTypeOrderCancel,
			Side:      int64(side),
			Symbol:    order.Product,
			Quantity:  order.Quantity.String(),
			Fee:       cancelFeeStr,
			Timestamp: timestamp,
		}

		result = append(result, &tx)
	}

	return result
}
