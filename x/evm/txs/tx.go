package txs

import (
	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
	"math/big"
)

type Tx interface {
	// Prepare convert msg to tx
	Prepare(msg *types.MsgEthereumTx) (err error)

	// SaveTx since the txCount is used by the stateDB, and a simulated tx is run only on the node it's submitted to,
	// then this will cause the txCount/stateDB of the node that ran the simulated tx to be different with the
	// other nodes, causing a consensus error
	SaveTx(msg *types.MsgEthereumTx)

	// GetChainConfig get chain config
	GetChainConfig() (*types.ChainConfig, bool)

	// GetSenderAccount get sender account
	GetSenderAccount() authexported.Account

	// ResetWatcher when panic reset watcher
	ResetWatcher(account authexported.Account)

	// RefundFeesWatcher fix account balance in watcher with refund fees
	RefundFeesWatcher(account authexported.Account, coins sdk.Coins, price *big.Int)

	// Transition execute evm tx
	Transition(config *types.ChainConfig) (result base.Result, err error)

	// DecorateResult some case(trace tx log) will modify the inResult to log and swallow inErr
	DecorateResult(inResult *base.Result, inErr error) (result *sdk.Result, err error)

	// RestoreWatcherTransactionReceipt restore watcher TransactionReceipt
	RestoreWatcherTransactionReceipt(msg *types.MsgEthereumTx)

	// Commit save the inner tx and contracts
	Commit(msg *types.MsgEthereumTx, result *base.Result)

	// EmitEvent emit event
	EmitEvent(msg *types.MsgEthereumTx, result *base.Result)

	// FinalizeWatcher after execute evm tx run here
	FinalizeWatcher(account authexported.Account, err error)

	// AnalyzeStart start record tag
	AnalyzeStart(tag string)

	// AnalyzeStop stop record tag
	AnalyzeStop(tag string)
}

// TransitionEvmTx execute evm transition template
func TransitionEvmTx(tx Tx, msg *types.MsgEthereumTx) (result *sdk.Result, err error) {
	tx.AnalyzeStart(bam.EvmHandler)
	defer tx.AnalyzeStop(bam.EvmHandler)

	// Prepare convert msg to state transition
	err = tx.Prepare(msg)
	if err != nil {
		return nil, err
	}

	// save tx
	tx.SaveTx(msg)

	// execute transition, the result
	tx.AnalyzeStart(bam.TransitionDb)
	defer tx.AnalyzeStop(bam.TransitionDb)

	config, found := tx.GetChainConfig()
	if !found {
		return nil, types.ErrChainConfigNotFound
	}

	defer func() {
		senderAccount := tx.GetSenderAccount()
		tx.RefundFeesWatcher(senderAccount, msg.GetFee(), msg.Data.Price)
		if e := recover(); e != nil {
			tx.ResetWatcher(senderAccount)
			panic(e)
		}
		tx.FinalizeWatcher(senderAccount, err)
	}()

	// execute evm tx
	var baseResult base.Result
	baseResult, err = tx.Transition(config)
	if err == nil {
		// Commit save the inner tx and contracts
		tx.Commit(msg, &baseResult)
		tx.EmitEvent(msg, &baseResult)
	} else {
		tx.RestoreWatcherTransactionReceipt(msg)
	}

	return tx.DecorateResult(&baseResult, err)
}
