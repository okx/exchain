package base

import (
	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
	"math/big"
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


func NewTx(config Config) *Tx {
	return &Tx{
		Ctx:    config.Ctx,
		Keeper: config.Keeper,
	}
}


// Prepare convert msg to state transition
func (tx *Tx) Prepare(msg *types.MsgEthereumTx) (err error) {
	tx.AnalyzeStart(bam.Txhash)
	defer tx.AnalyzeStop(bam.Txhash)

	tx.StateTransition, err = msg2st(&tx.Ctx, tx.Keeper, msg)
	return
}

func (tx *Tx) GetChainConfig() (types.ChainConfig, bool) {
	return tx.Keeper.GetChainConfig(tx.Ctx)
}


// Transition execute evm tx
func (tx *Tx) Transition(config types.ChainConfig) (result Result, err error) {
	result.ExecResult, result.ResultData, err, result.InnerTxs, result.Erc20Contracts = tx.StateTransition.TransitionDb(tx.Ctx, config)
	// async mod goes immediately
	if tx.Ctx.IsAsync() {
		tx.Keeper.LogsManages.Set(string(tx.Ctx.TxBytes()), keeper.TxResult{
			ResultData: result.ResultData,
			Err:        err,
		})
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

// SaveTx check Tx do not transition state db
func (t *Tx) SaveTx(msg *types.MsgEthereumTx) {}

func (t *Tx) GetSenderAccount() authexported.Account { return nil }

func (t *Tx) ResetWatcher(account authexported.Account) {}

// RefundFeesWatcher refund the watcher, check Tx do not save state so. skip
func (t *Tx) RefundFeesWatcher(account authexported.Account, coins sdk.Coins, price *big.Int) {}

// RestoreWatcherTransactionReceipt check Tx do not need restore
func (t *Tx) RestoreWatcherTransactionReceipt(msg *types.MsgEthereumTx) {}

// Commit check Tx do not need
func (t *Tx) Commit(msg *types.MsgEthereumTx, result *Result) {}

// FinalizeWatcher check Tx do not need this
func (t *Tx) FinalizeWatcher(account authexported.Account, err error) {}

// AnalyzeStart check Tx do not analyze start
func (t *Tx) AnalyzeStart(tag string) {}

// AnalyzeStop check Tx do not analyze stop
func (t *Tx) AnalyzeStop(tag string) {}

