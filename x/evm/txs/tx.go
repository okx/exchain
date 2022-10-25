package txs

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	bam "github.com/okex/exchain/libs/system/trace"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
)

type Tx interface {
	// Prepare convert msg to tx
	Prepare(msg *types.MsgEthereumTx) (err error)

	// SaveTx since the txCount is used by the stateDB, and a simulated tx is run only on the node it's submitted to,
	// then this will cause the txCount/stateDB of the node that ran the simulated tx to be different with the
	// other nodes, causing a consensus error
	SaveTx(msg *types.MsgEthereumTx)

	// GetChainConfig get chain config(the chain config may cached)
	GetChainConfig() (types.ChainConfig, bool)

	// GetSenderAccount get sender account
	GetSenderAccount() authexported.Account

	// Transition execute evm tx
	Transition(config types.ChainConfig) (result base.Result, err error)

	// DecorateResult some case(trace tx log) will modify the inResult to log and swallow inErr
	DecorateResult(inResult *base.Result, inErr error) (result *sdk.Result, err error)

	// Commit save the inner tx and contracts
	Commit(msg *types.MsgEthereumTx, result *base.Result)

	// EmitEvent emit event
	EmitEvent(msg *types.MsgEthereumTx, result *base.Result)

	// FinalizeWatcher after execute evm tx run here
	FinalizeWatcher(msg *types.MsgEthereumTx, err error, panic bool)

	// AnalyzeStart start record tag
	AnalyzeStart(tag string)

	// AnalyzeStop stop record tag
	AnalyzeStop(tag string)

	// Dispose release the resources of the tx, should be called after the tx is unused
	Dispose()
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
		e := recover()
		isPanic := e != nil
		tx.FinalizeWatcher(msg, err, isPanic)
		if isPanic {
			panic(e)
		}
	}()

	// execute evm tx
	var baseResult base.Result
	baseResult, err = tx.Transition(config)
	if err == nil {
		// Commit save the inner tx and contracts
		tx.Commit(msg, &baseResult)
		tx.EmitEvent(msg, &baseResult)
	}

	return tx.DecorateResult(&baseResult, err)
}
