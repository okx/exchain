package evm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/x/common/perf"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	tmtypes "github.com/tendermint/tendermint/types"
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

		seq := perf.GetPerf().OnDeliverTxEnter(ctx, types.ModuleName, name)
		defer perf.GetPerf().OnDeliverTxExit(ctx, types.ModuleName, name, seq)

		result, err = handlerFun()
		if err != nil {
			err = sdkerrors.New(types.ModuleName, types.CodeSpaceEvmCallFailed, err.Error())
		}

		return result, err
	}
}

// handleMsgEthereumTx handles an Ethereum specific tx
func handleMsgEthereumTx(ctx sdk.Context, k *Keeper, msg types.MsgEthereumTx) (*sdk.Result, error) {
	// parse the chainID from a string to a base-10 integer
	chainIDEpoch, err := ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return nil, err
	}

	// Verify signature and retrieve sender address
	sender, err := msg.VerifySig(chainIDEpoch, ctx.BlockHeight())
	if err != nil {
		return nil, err
	}

	txHash := tmtypes.Tx(ctx.TxBytes()).Hash()
	ethHash := common.BytesToHash(txHash)

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
		k.Watcher.SaveEthereumTx(msg, common.BytesToHash(txHash), uint64(ctx.EvmTransactionIndex()))
		// Prepare db for logs
		st.Csdb.Prepare(ethHash, k.Bhash, int(ctx.EvmTransactionIndex()))
		st.Csdb.SetLogSize(k.LogSize)
	}

	config, found := k.GetChainConfig(ctx)
	if !found {
		return nil, types.ErrChainConfigNotFound
	}

	defer func() {
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

	executionResult, resultData, err := st.TransitionDb(ctx, config)
	if ctx.IsAsync() {
		tmp := keeper.TxMapping{
			ResultData: resultData,
		}
		k.Mmpp[ctx.EvmTransactionIndex()] = tmp
	}

	if err != nil {
		if !st.Simulate {
			k.Watcher.SaveTransactionReceipt(watcher.TransactionFailed, msg, common.BytesToHash(txHash), uint64(ctx.EvmTransactionIndex()-1), &types.ResultData{}, ctx.GasMeter().GasConsumed())
		}
		return nil, err
	}

	if !st.Simulate {
		// update block bloom filter
		k.Bloom.Or(k.Bloom, executionResult.Bloom)
		k.LogSize = st.Csdb.GetLogSize()
		k.Watcher.SaveTransactionReceipt(watcher.TransactionSuccess, msg, common.BytesToHash(txHash), uint64(ctx.EvmTransactionIndex()-1), resultData, ctx.GasMeter().GasConsumed())
		if msg.Data.Recipient == nil {
			st.Csdb.IteratorCode(func(addr common.Address, c types.CacheCode) bool {
				k.Watcher.SaveContractCode(addr, c.Code)
				k.Watcher.SaveContractCodeByHash(c.CodeHash, c.Code)
				return true
			})
		}
	}

	// log successful execution
	k.Logger(ctx).Info(executionResult.Result.Log)

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
		st.Csdb.Prepare(ethHash, k.Bhash, int(ctx.EvmTransactionIndex()))
		st.Csdb.SetLogSize(k.LogSize)
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

	// log successful execution
	k.Logger(ctx).Info(executionResult.Result.Log)

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
