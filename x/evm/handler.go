package evm

import (
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/exchain/app/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/common/analyzer"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
)

// NewHandler returns a handler for Ethermint type messages.
func NewHandler(k *Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, err error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		var handlerFun func() (*sdk.Result, error)
		var name string
		switch msg := msg.(type) {
		case types.MsgEthereumTx:
			name = "handleMsgEthereumTx"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgEthereumTx(ctx, k, msg)
			}
		case types.MsgEthermint:
			name = "handleMsgEthermint"
			handlerFun = func() (*sdk.Result, error) {
				return handleMsgEthermint(ctx, k, msg)
			}
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}

		_ = name

		result, err = handlerFun()
		if err != nil {
			err = sdkerrors.New(types.ModuleName, types.CodeSpaceEvmCallFailed, err.Error())
		}

		return result, err
	}
}

// handleMsgEthereumTx handles an Ethereum specific tx
func handleMsgEthereumTx(ctx sdk.Context, k *Keeper, msg types.MsgEthereumTx) (*sdk.Result, error) {
	StartTxLog := func(tag string) {
		if !ctx.IsCheckTx() {
			analyzer.StartTxLog(tag)
		}
	}
	StopTxLog := func(tag string) {
		if !ctx.IsCheckTx() {
			analyzer.StopTxLog(tag)
		}
	}

	// parse the chainID from a string to a base-10 integer
	StartTxLog(bam.EvmHandler)
	defer StopTxLog(bam.EvmHandler)

	StartTxLog(bam.ParseChainID)
	chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return nil, err
	}
	StopTxLog(bam.ParseChainID)

	// Verify signature and retrieve sender address

	StartTxLog(bam.VerifySig)
	senderSigCache, err := msg.VerifySig(chainIDEpoch, ctx.BlockHeight(), ctx.SigCache())
	if err != nil {
		return nil, err
	}
	StopTxLog(bam.VerifySig)

	StartTxLog(bam.Txhash)
	sender := senderSigCache.GetFrom()
	txHash := tmtypes.Tx(ctx.TxBytes()).Hash()
	ethHash := common.BytesToHash(txHash)
	StopTxLog(bam.Txhash)

	StartTxLog(bam.SaveTx)
	st := types.StateTransition{
		AccountNonce: msg.Data.AccountNonce,
		Price:        msg.Data.Price,
		GasLimit:     msg.Data.GasLimit,
		Recipient:    msg.Data.Recipient,
		Amount:       msg.Data.Amount,
		Payload:      msg.Data.Payload,
		Csdb:         types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx),
		ChainID:      chainIDEpoch,
		TxHash:       &ethHash,
		Sender:       sender,
		Simulate:     ctx.IsCheckTx(),
	}

	// since the txCount is used by the stateDB, and a simulated tx is run only on the node it's submitted to,
	// then this will cause the txCount/stateDB of the node that ran the simulated tx to be different than the
	// other nodes, causing a consensus error

	if !st.Simulate {
		k.Watcher.SaveEthereumTx(msg, common.BytesToHash(txHash), uint64(k.TxCount))
		// Prepare db for logs
		st.Csdb.Prepare(ethHash, k.Bhash, k.TxCount)
		st.Csdb.SetLogSize(k.LogSize)
		k.TxCount++
	}

	config, found := k.GetChainConfig(ctx)
	if !found {
		return nil, types.ErrChainConfigNotFound
	}

	StopTxLog(bam.SaveTx)

	defer func() {
		StartTxLog(bam.HandlerDefer)
		defer StopTxLog(bam.HandlerDefer)

		if !st.Simulate && k.Watcher.Enabled() {
			currentGasMeter := ctx.GasMeter()
			pm := k.GenerateCSDBParams()
			infCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
			sendAcc := pm.AccountKeeper.GetAccount(infCtx, sender.Bytes())
			if sendAcc != nil {
				pm.Watcher.SaveAccount(sendAcc, true)
			}
			ctx.WithGasMeter(currentGasMeter)
		}
		if e := recover(); e != nil {
			k.Watcher.Reset()
			panic(e)
		}
		if !st.Simulate {
			if err != nil {
				k.Watcher.Reset()
			} else {
				//save state and account data into batch
				k.Watcher.Finalize()
			}
		}
	}()

	StartTxLog(bam.TransitionDb)
	executionResult, resultData, err := st.TransitionDb(ctx, config)
	if ctx.IsAsync() {
		k.LogsManages.Set(string(ctx.TxBytes()), keeper.TxResult{
			ResultData: resultData,
			Err:        err,
		})
	}

	if err != nil {
		if !st.Simulate {
			k.Watcher.SaveTransactionReceipt(watcher.TransactionFailed, msg, common.BytesToHash(txHash), uint64(k.TxCount-1), &types.ResultData{}, ctx.GasMeter().GasConsumed())
		}
		return nil, err
	}
	StopTxLog(bam.TransitionDb)

	StartTxLog(bam.Bloomfilter)
	if !st.Simulate {
		// update block bloom filter
		if !ctx.IsAsync() {
			k.Bloom.Or(k.Bloom, executionResult.Bloom) // not support paralleled-tx´
		}
		k.LogSize = st.Csdb.GetLogSize()
		k.Watcher.SaveTransactionReceipt(watcher.TransactionSuccess, msg, common.BytesToHash(txHash), uint64(k.TxCount-1), resultData, ctx.GasMeter().GasConsumed())
		if msg.Data.Recipient == nil {
			st.Csdb.IteratorCode(func(addr common.Address, c types.CacheCode) bool {
				k.Watcher.SaveContractCode(addr, c.Code)
				k.Watcher.SaveContractCodeByHash(c.CodeHash, c.Code)
				return true
			})
		}
	}
	StopTxLog(bam.Bloomfilter)

	StartTxLog(bam.EmitEvents)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEthereumTx,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Data.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, sender.String()),
		),
	})

	if msg.Data.Recipient != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeEthereumTx,
				sdk.NewAttribute(types.AttributeKeyRecipient, msg.Data.Recipient.String()),
			),
		)
	}

	// set the events to the result
	executionResult.Result.Events = ctx.EventManager().Events()
	StopTxLog(bam.EmitEvents)
	return executionResult.Result, nil
}

// handleMsgEthermint handles an sdk.StdTx for an Ethereum state transition
func handleMsgEthermint(ctx sdk.Context, k *Keeper, msg types.MsgEthermint) (*sdk.Result, error) {

	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		return nil, sdkerrors.Wrap(ethermint.ErrInvalidMsgType, "Ethermint type message is not allowed.")
	}

	// parse the chainID from a string to a base-10 integer
	chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return nil, err
	}

	txHash := tmtypes.Tx(ctx.TxBytes()).Hash()
	ethHash := common.BytesToHash(txHash)

	st := types.StateTransition{
		AccountNonce: msg.AccountNonce,
		Price:        msg.Price.BigInt(),
		GasLimit:     msg.GasLimit,
		Amount:       msg.Amount.BigInt(),
		Payload:      msg.Payload,
		Csdb:         types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), ctx),
		ChainID:      chainIDEpoch,
		TxHash:       &ethHash,
		Sender:       common.BytesToAddress(msg.From.Bytes()),
		Simulate:     ctx.IsCheckTx(),
	}

	if msg.Recipient != nil {
		to := common.BytesToAddress(msg.Recipient.Bytes())
		st.Recipient = &to
	}

	if !st.Simulate {
		// Prepare db for logs
		st.Csdb.Prepare(ethHash, k.Bhash, k.TxCount)
		st.Csdb.SetLogSize(k.LogSize)
		k.TxCount++
	}

	config, found := k.GetChainConfig(ctx)
	if !found {
		return nil, types.ErrChainConfigNotFound
	}

	executionResult, _, err := st.TransitionDb(ctx, config)
	if err != nil {
		return nil, err
	}

	// update block bloom filter
	if !st.Simulate {
		k.Bloom.Or(k.Bloom, executionResult.Bloom)
		k.LogSize = st.Csdb.GetLogSize()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEthermint,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.From.String()),
		),
	})

	if msg.Recipient != nil {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeEthermint,
				sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient.String()),
			),
		)
	}

	// set the events to the result
	executionResult.Result.Events = ctx.EventManager().Events()
	return executionResult.Result, nil
}
