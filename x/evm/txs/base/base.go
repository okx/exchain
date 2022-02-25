package base

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
)

// Config tx's needed ctx and keeper
type Config struct {
	Ctx    sdk.Context
	Keeper *evm.Keeper
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
	Keeper *evm.Keeper

	StateTransition types.StateTransition
	// state transition result
	Result Result
}

func NewTx(config Config) *Tx {
	return &Tx{
		Ctx:    config.Ctx,
		Keeper: config.Keeper,
	}
}

// Prepare convert msg to state transition
func (tx *Tx) Prepare(msg *types.MsgEthereumTx) (err error) {
	tx.StateTransition, err = msg2st(&tx.Ctx, tx.Keeper, msg)
	return
}

// SaveTx since the txCount is used by the stateDB, and a simulated tx is run only on the node it's submitted to,
// then this will cause the txCount/stateDB of the node that ran the simulated tx to be different with the
// other nodes, causing a consensus error
func (tx *Tx) SaveTx(msg *types.MsgEthereumTx) {
	tx.Keeper.Watcher.SaveEthereumTx(msg, *tx.StateTransition.TxHash, uint64(tx.Keeper.TxCount))
	// Prepare db for logs
	tx.StateTransition.Csdb.Prepare(*tx.StateTransition.TxHash, tx.Keeper.Bhash, tx.Keeper.TxCount)
	tx.StateTransition.Csdb.SetLogSize(tx.Keeper.LogSize)
	tx.Keeper.TxCount++
}

// Transition execute evm tx
func (tx *Tx) Transition() (err error) {
	config, found := tx.Keeper.GetChainConfig(tx.Ctx)
	if !found {
		return types.ErrChainConfigNotFound
	}
	tx.Result.ExecResult, tx.Result.ResultData, err, tx.Result.InnerTxs, tx.Result.Erc20Contracts = tx.StateTransition.TransitionDb(tx.Ctx, config)
	// async mod goes immediately
	if tx.Ctx.IsAsync() {
		tx.Keeper.LogsManages.Set(string(tx.Ctx.TxBytes()), keeper.TxResult{
			ResultData: tx.Result.ResultData,
			Err:        err,
		})
	}

	return
}

// DecorateError TraceTxLog situation Decorate the result
// it was replaced to trace logs when trace tx even if err != nil
func (tx *Tx) DecorateError(err error) (*sdk.Result, error) {
	if tx.Ctx.IsTraceTxLog() {
		// the result was replaced to trace logs when trace tx even if err != nil
		tx.Result.ExecResult.Result.Data = tx.Result.ExecResult.TraceLogs
		return tx.Result.ExecResult.Result, nil
	}

	return nil, err
}

func (tx *Tx) Emit(msg *types.MsgEthereumTx) *sdk.Result {
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
	tx.Result.ExecResult.Result.Events = tx.Ctx.EventManager().Events()
	if tx.Ctx.IsTraceTxLog() {
		// the result was replaced to trace logs when trace tx
		tx.Result.ExecResult.Result.Data = tx.Result.ExecResult.TraceLogs
	}

	return tx.Result.ExecResult.Result
}

//
func (tx *Tx) Finalize() error {
	//TODO implement me
	panic("implement me")
}
