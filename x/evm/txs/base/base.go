package base

import (
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	bam "github.com/okex/exchain/libs/system/trace"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
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
	reuseCsdb       bool
}

// Prepare convert msg to state transition
func (tx *Tx) Prepare(msg *types.MsgEthereumTx) (err error) {
	tx.AnalyzeStart(bam.Txhash)
	defer tx.AnalyzeStop(bam.Txhash)

	tx.reuseCsdb, err = msg2st(&tx.Ctx, tx.Keeper, msg, &tx.StateTransition)
	return
}

// GetChainConfig get chain config, the chain config may cached
func (tx *Tx) GetChainConfig() (types.ChainConfig, bool) {
	return tx.Keeper.GetChainConfig(tx.Ctx)
}

// Transition execute evm tx
func (tx *Tx) Transition(config types.ChainConfig) (result Result, err error) {
	result.ExecResult, result.ResultData, err, result.InnerTxs, result.Erc20Contracts = tx.StateTransition.TransitionDb(tx.Ctx, config)

	if err != nil {
		return
	}

	// call evm hooks
	if tmtypes.HigherThanVenus1(tx.Ctx.BlockHeight()) && !tx.Ctx.IsCheckTx() {
		receipt := &ethtypes.Receipt{
			Status:           ethtypes.ReceiptStatusSuccessful,
			Bloom:            result.ResultData.Bloom,
			Logs:             result.ResultData.Logs,
			TxHash:           result.ResultData.TxHash,
			ContractAddress:  result.ResultData.ContractAddress,
			GasUsed:          result.ExecResult.GasInfo.GasConsumed,
			BlockNumber:      big.NewInt(tx.Ctx.BlockHeight()),
			TransactionIndex: uint(tx.Keeper.TxCount),
		}
		err = tx.Keeper.CallEvmHooks(tx.Ctx, tx.StateTransition.Sender, tx.StateTransition.Recipient, receipt)
		if err != nil {
			tx.Keeper.Logger().Error("tx call evm hooks failed", "error", err)
		}
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

func NewTx(config Config) *Tx {
	return &Tx{
		Ctx:    config.Ctx,
		Keeper: config.Keeper,
	}
}

func (tx *Tx) AnalyzeStart(tag string) {
	bam.StartTxLog(tag)
}

func (tx *Tx) AnalyzeStop(tag string) {
	bam.StopTxLog(tag)
}

// SaveTx check Tx do not transition state db
func (tx *Tx) SaveTx(msg *types.MsgEthereumTx) {}

// GetSenderAccount check Tx do not need this
func (tx *Tx) GetSenderAccount() authexported.Account { return nil }

// ResetWatcher check Tx do not need this
func (tx *Tx) ResetWatcher(account authexported.Account) {}

// RefundFeesWatcher refund the watcher, check Tx do not save state so. skip
func (tx *Tx) RefundFeesWatcher(account authexported.Account, ethereumTx *types.MsgEthereumTx) {}

// Commit check Tx do not need
func (tx *Tx) Commit(msg *types.MsgEthereumTx, result *Result) {}

// FinalizeWatcher check Tx do not need this
func (tx *Tx) FinalizeWatcher(msg *types.MsgEthereumTx, account authexported.Account, err error) {}

func (tx *Tx) Dispose() {
	if tx != nil && tx.reuseCsdb {
		tx.reuseCsdb = false
		if tx.StateTransition.Csdb != nil {
			putCommitStateDB(tx.StateTransition.Csdb)
			tx.StateTransition.Csdb = nil
		}
	}
}
