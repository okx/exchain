package deliver

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/app/refund"
	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"math/big"
)

type Tx struct {
	*base.Tx
}

func NewTx(config base.Config) *Tx {
	return &Tx{
		Tx: base.NewTx(config),
	}
}

// SaveTx since the txCount is used by the stateDB, and a simulated tx is run only on the node it's submitted to,
// then this will cause the txCount/stateDB of the node that ran the simulated tx to be different with the
// other nodes, causing a consensus error
func (tx *Tx) SaveTx(msg *types.MsgEthereumTx) {
	tx.AnalyzeStart(bam.SaveTx)
	defer tx.AnalyzeStop(bam.SaveTx)

	tx.Keeper.Watcher.SaveEthereumTx(msg, *tx.StateTransition.TxHash, uint64(tx.Keeper.TxCount))
	// Prepare db for logs
	tx.StateTransition.Csdb.Prepare(*tx.StateTransition.TxHash, tx.Keeper.Bhash, tx.Keeper.TxCount)
	tx.StateTransition.Csdb.SetLogSize(tx.Keeper.LogSize)
	tx.Keeper.TxCount++
}

func (tx *Tx) GetSenderAccount() authexported.Account {
	pm := tx.Keeper.GenerateCSDBParams()
	infCtx := tx.Ctx.WithGasMeter(sdk.NewInfiniteGasMeter())

	return pm.AccountKeeper.GetAccount(infCtx, tx.StateTransition.Sender.Bytes())
}

func (tx *Tx) ResetWatcher(account authexported.Account) {
	tx.Keeper.Watcher.Reset()
	// delete account which is already in Watcher.batch
	if account != nil {
		tx.Keeper.Watcher.AddDelAccMsg(account, true)
	}
}

func (tx *Tx) RefundFeesWatcher(account authexported.Account, coin sdk.Coins, price *big.Int) {
	// fix account balance in watcher with refund fees
	if account == nil || !tx.Keeper.Watcher.Enabled() {
		return
	}
	gasConsumed := tx.Ctx.GasMeter().GasConsumed()
	fixedFees := refund.CaculateRefundFees(gasConsumed, coin, price)
	coins := account.GetCoins().Add2(fixedFees)
	account.SetCoins(coins) //ignore err, no err will be returned in SetCoins

	pm := tx.Keeper.GenerateCSDBParams()
	pm.Watcher.SaveAccount(account, false)
}

func (tx *Tx) RestoreWatcherTransactionReceipt(msg *types.MsgEthereumTx) {
	tx.Keeper.Watcher.SaveTransactionReceipt(
		watcher.TransactionFailed,
		msg,
		*tx.StateTransition.TxHash,
		uint64(tx.Keeper.TxCount-1),
		&types.ResultData{}, tx.Ctx.GasMeter().GasConsumed())
}

func (tx *Tx) Commit(msg *types.MsgEthereumTx, result *base.Result) {
	if result.InnerTxs != nil {
		tx.Keeper.AddInnerTx(tx.StateTransition.TxHash.Hex(), result.InnerTxs)
	}
	if result.Erc20Contracts != nil {
		tx.Keeper.AddContract(result.Erc20Contracts)
	}

	// update block bloom filter
	if !tx.Ctx.IsAsync() {
		tx.Keeper.Bloom.Or(tx.Keeper.Bloom, result.ExecResult.Bloom)
	}
	tx.Keeper.LogSize = tx.StateTransition.Csdb.GetLogSize()
	tx.Keeper.Watcher.SaveTransactionReceipt(watcher.TransactionSuccess,
		msg, *tx.StateTransition.TxHash,
		uint64(tx.Keeper.TxCount-1), result.ResultData, tx.Ctx.GasMeter().GasConsumed())
	if msg.Data.Recipient == nil {
		tx.StateTransition.Csdb.IteratorCode(func(addr common.Address, c types.CacheCode) bool {
			tx.Keeper.Watcher.SaveContractCode(addr, c.Code)
			tx.Keeper.Watcher.SaveContractCodeByHash(c.CodeHash, c.Code)
			return true
		})
	}
}

func (tx *Tx) FinalizeWatcher(account authexported.Account, err error) {
	if err != nil {
		tx.ResetWatcher(account)
		return
	}
	tx.Keeper.Watcher.Finalize()
}
