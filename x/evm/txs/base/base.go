package base

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/exchain/app/refund"
	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/x/common/analyzer"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
)

// Keeper alias of keeper.Keeper, to solve import circle. also evm.Keeper is alias keeper.Keeper
type Keeper = keeper.Keeper

// Config tx's needed ctx and keeper
type Config struct {
	Ctx    sdk.Context
	Keeper *Keeper
}

// Result evm execute result
type Result struct {
	ExecResult     *types.ExecutionResult
	ResultData     *types.ResultData
	InnerTxs       interface{}
	Erc20Contracts interface{}
}

// Tx evm tx
type Tx struct {
	Ctx    sdk.Context
	Keeper *Keeper

	StateTransition types.StateTransition
}

// Prepare convert msg to state transition
func (tx *Tx) Prepare(msg *types.MsgEthereumTx) (err error) {
	tx.AnalyzeStart(bam.Txhash)
	defer tx.AnalyzeStop(bam.Txhash)

	tx.StateTransition, err = msg2st(&tx.Ctx, tx.Keeper, msg)
	return
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

func (tx *Tx) GetChainConfig() (types.ChainConfig, bool) {
	return tx.Keeper.GetChainConfig(tx.Ctx)
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

// Transition execute evm tx
func (tx *Tx) Transition(config types.ChainConfig) (result Result, err error) {
	// snapshot to contain the tx processing and post processing in same scope
	var commit func()
	tmpCtx := tx.Ctx
	if tx.Keeper.GetHooks() != nil {
		// Create a cache context to revert state when tx hooks fails,
		// the cache context is only committed when both tx and hooks executed successfully.
		// Didn't use `Snapshot` because the context stack has exponential complexity on certain operations,
		// thus restricted to be used only inside `ApplyMessage`.
		tmpCtx, commit = tx.Ctx.CacheContext()
	}

	result.ExecResult, result.ResultData, err, result.InnerTxs, result.Erc20Contracts = tx.StateTransition.TransitionDb(tmpCtx, config)
	// async mod goes immediately
	if tx.Ctx.IsAsync() {
		tx.Keeper.LogsManages.Set(string(tx.Ctx.TxBytes()), keeper.TxResult{
			ResultData: result.ResultData,
			Err:        err,
		})
	}
	if err != nil {
		return
	}

	// call evm hooks
	receipt := &ethtypes.Receipt{
		//Type:              ethtypes.DynamicFeeTxType,// TODO: hardcode
		PostState:         nil, // TODO: intermediate state root
		Status:            ethtypes.ReceiptStatusSuccessful,
		CumulativeGasUsed: 0, // TODO: cumulativeGasUsed
		Bloom:             result.ResultData.Bloom,
		Logs:              result.ResultData.Logs,
		TxHash:            result.ResultData.TxHash,
		ContractAddress:   result.ResultData.ContractAddress,
		GasUsed:           result.ExecResult.GasInfo.GasConsumed,
		BlockHash:         tx.Keeper.GetHeightHash(tx.Ctx, uint64(tx.Ctx.BlockHeight())),
		BlockNumber:       big.NewInt(tx.Ctx.BlockHeight()),
		TransactionIndex:  uint(tx.Keeper.TxCount),
	}
	if err = tx.Keeper.CallEvmHooks(tmpCtx, tx.StateTransition.Sender, tx.StateTransition.Recipient, receipt); err != nil {
		tx.Keeper.Logger(tx.Ctx).Error("tx post processing failed", "error", err)
	} else if commit != nil {
		// PostTxProcessing is successful, commit the tmpCtx
		commit()
		tx.Ctx.EventManager().EmitEvents(tmpCtx.EventManager().Events())
	}

	return
}

// DecorateResult TraceTxLog situation Decorate the result
// it was replaced to trace logs when trace tx even if err != nil
func (tx *Tx) DecorateResult(inResult *Result, inErr error) (result *sdk.Result, err error) {
	if inErr != nil {
		return nil, inErr
	}
	return inResult.ExecResult.Result, inErr
}

func (tx *Tx) RestoreWatcherTransactionReceipt(msg *types.MsgEthereumTx) {
	tx.Keeper.Watcher.SaveTransactionReceipt(
		watcher.TransactionFailed,
		msg,
		*tx.StateTransition.TxHash,
		uint64(tx.Keeper.TxCount-1),
		&types.ResultData{}, tx.Ctx.GasMeter().GasConsumed())
}

func (tx *Tx) Commit(msg *types.MsgEthereumTx, result *Result) {
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

func (tx *Tx) EmitEvent(msg *types.MsgEthereumTx, result *Result) {
	tx.Ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEthereumTx,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Data.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, types.EthAddressStringer(tx.StateTransition.Sender).String()),
		),
	})

	if msg.Data.Recipient != nil {
		tx.Ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeEthereumTx,
				sdk.NewAttribute(types.AttributeKeyRecipient, types.EthAddressStringer(*msg.Data.Recipient).String()),
			),
		)
	}

	// set the events to the result
	result.ExecResult.Result.Events = tx.Ctx.EventManager().Events()
}

func (tx *Tx) FinalizeWatcher(account authexported.Account, err error) {
	if err != nil {
		tx.ResetWatcher(account)
		return
	}
	tx.Keeper.Watcher.Finalize()
}

func (tx *Tx) AnalyzeStart(tag string) {
	analyzer.StartTxLog(tag)
}

func (tx *Tx) AnalyzeStop(tag string) {
	analyzer.StopTxLog(tag)
}

func NewTx(config Config) *Tx {
	return &Tx{
		Ctx:    config.Ctx,
		Keeper: config.Keeper,
	}
}
