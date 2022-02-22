package evm

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/app/refund"
	ethermint "github.com/okex/exchain/app/types"
	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	common2 "github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/common/analyzer"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"math/big"
)

// NewHandler returns a handler for Ethermint type messages.
func NewHandler(k *Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (result *sdk.Result, err error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		defer func() {
			if cfg.DynamicConfig.GetMaxGasUsedPerBlock() < 0 {
				return
			}

			if err != nil {
				return
			}

			db := bam.InstanceOfHistoryGasUsedRecordDB()
			msgFnSignature, toDeployContractSize := getMsgCallFnSignature(msg)

			if msgFnSignature == nil {
				return
			}

			hisGu, err := db.Get(msgFnSignature)
			if err != nil {
				return
			}

			gc := int64(ctx.GasMeter().GasConsumed())
			if toDeployContractSize > 0 {
				// calculate average gas consume for deploy contract case
				gc = gc / int64(toDeployContractSize)
			}

			var avgGas int64
			if hisGu != nil {
				hgu := common2.BytesToInt64(hisGu)
				avgGas = int64(bam.GasUsedFactor*float64(gc) + (1.0-bam.GasUsedFactor)*float64(hgu))
			} else {
				avgGas = gc
			}

			err = db.Set(msgFnSignature, common2.Int64ToBytes(avgGas))
			if err != nil {
				return
			}
		}()

		evmtx, ok := msg.(types.MsgEthereumTx)
		if ok {
			result, err = handleMsgEthereumTx(ctx, k, &evmtx)
			if err != nil {
				err = sdkerrors.New(types.ModuleName, types.CodeSpaceEvmCallFailed, err.Error())
			}
		} else {
			err = sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", ModuleName, msg)
		}

		return result, err
	}
}

func getMsgCallFnSignature(msg sdk.Msg) ([]byte, int) {
	switch msg := msg.(type) {
	case types.MsgEthereumTx:
		return msg.GetTxFnSignatureInfo()
	default:
		return nil, 0
	}
}

// handleMsgEthereumTx handles an Ethereum specific tx
func handleMsgEthereumTx(ctx sdk.Context, k *Keeper, msg *types.MsgEthereumTx) (*sdk.Result, error) {
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

	StartTxLog(bam.Txhash)
	st, err := msg2st(&ctx, k, msg)
	StopTxLog(bam.Txhash)

	if err != nil {
		return nil, err
	}


	StartTxLog(bam.SaveTx)
	// since the txCount is used by the stateDB, and a simulated tx is run only on the node it's submitted to,
	// then this will cause the txCount/stateDB of the node that ran the simulated tx to be different than the
	// other nodes, causing a consensus error
	if !st.Simulate {
		k.Watcher.SaveEthereumTx(msg, *st.TxHash, uint64(k.TxCount))
		// Prepare db for logs
		st.Csdb.Prepare(*st.TxHash, k.Bhash, k.TxCount)
		st.Csdb.SetLogSize(k.LogSize)
		k.TxCount++
	}
	config, found := k.GetChainConfig(ctx)
	StopTxLog(bam.SaveTx)

	if !found {
		return nil, types.ErrChainConfigNotFound
	}

	defer func() {
		pm := k.GenerateCSDBParams()
		infCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
		sendAcc := pm.AccountKeeper.GetAccount(infCtx, st.Sender.Bytes())
		if !st.Simulate && k.Watcher.Enabled() && sendAcc != nil {
			//fix sender's balance in watcher with refund fees
			gasConsumed := ctx.GasMeter().GasConsumed()
			fixedFees := refund.CaculateRefundFees(gasConsumed, msg.GetFee(), msg.Data.Price)
			coins := sendAcc.GetCoins().Add2(fixedFees)
			sendAcc.SetCoins(coins) //ignore err, no err will be returned in SetCoins
			pm.Watcher.SaveAccount(sendAcc, false)
		}
		if e := recover(); e != nil {
			k.Watcher.Reset()
			// delete account which is already in Watcher.batch
			if sendAcc != nil {
				k.Watcher.AddDelAccMsg(sendAcc, true)
			}
			panic(e)
		}
		if !st.Simulate {
			if err != nil {
				k.Watcher.Reset()
				if sendAcc != nil {
					// delete account which is already in Watcher.batch
					k.Watcher.AddDelAccMsg(sendAcc, true)
				}
			} else {
				//save state and account data into batch
				k.Watcher.Finalize()
			}
		}
	}()

	StartTxLog(bam.TransitionDb)
	defer StopTxLog(bam.TransitionDb)
	executionResult, resultData, err, innerTxs, erc20s := st.TransitionDb(ctx, config)
	if ctx.IsAsync() {
		k.LogsManages.Set(string(ctx.TxBytes()), keeper.TxResult{
			ResultData: resultData,
			Err:        err,
		})
	}

	if err != nil {
		if !st.Simulate {
			k.Watcher.SaveTransactionReceipt(watcher.TransactionFailed, msg, *st.TxHash, uint64(k.TxCount-1), &types.ResultData{}, ctx.GasMeter().GasConsumed())
		}
		if ctx.IsTraceTxLog() {
			// the result was replaced to trace logs when trace tx even if err != nil
			executionResult.Result.Data = executionResult.TraceLogs
			return executionResult.Result, nil
		}
		return nil, err
	}

	if !st.Simulate {
		if innerTxs != nil {
			k.AddInnerTx(st.TxHash.Hex(), innerTxs)
		}
		if erc20s != nil {
			k.AddContract(erc20s)
		}

		// update block bloom filter
		if !ctx.IsAsync() {
			k.Bloom.Or(k.Bloom, executionResult.Bloom) // not support paralleled-tx´
		}
		k.LogSize = st.Csdb.GetLogSize()
		k.Watcher.SaveTransactionReceipt(watcher.TransactionSuccess, msg, *st.TxHash, uint64(k.TxCount-1), resultData, ctx.GasMeter().GasConsumed())
		if msg.Data.Recipient == nil {
			st.Csdb.IteratorCode(func(addr common.Address, c types.CacheCode) bool {
				k.Watcher.SaveContractCode(addr, c.Code)
				k.Watcher.SaveContractCodeByHash(c.CodeHash, c.Code)
				return true
			})
		}
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeEthereumTx,
			sdk.NewAttribute(sdk.AttributeKeyAmount, msg.Data.Amount.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, st.Sender.String()),
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
	if ctx.IsTraceTxLog() {
		// the result was replaced to trace logs when trace tx
		executionResult.Result.Data = executionResult.TraceLogs
	}
	return executionResult.Result, nil
}



func getSender(ctx *sdk.Context, chainIDEpoch *big.Int, msg *types.MsgEthereumTx) (sender common.Address, err error) {
	if ctx.IsCheckTx() {
		if from := ctx.From(); len(from) > 0 {
			sender = common.HexToAddress(from)
		}
	}
	if len(sender) == 0 {
		senderSigCache, err := msg.VerifySig(chainIDEpoch, ctx.BlockHeight(), ctx.SigCache())
		if err == nil {
			sender = senderSigCache.GetFrom()
		}
	}

	return
}

func msg2st(ctx *sdk.Context, k *Keeper, msg *types.MsgEthereumTx) (st types.StateTransition, err error) {

	var chainIDEpoch *big.Int
	chainIDEpoch, err = ethermint.ParseChainID(ctx.ChainID())
	if err != nil {
		return
	}

	var sender common.Address
	// Verify signature and retrieve sender address
	sender, err = getSender(ctx, chainIDEpoch, msg)
	if err != nil {
		return
	}

	txHash := tmtypes.Tx(ctx.TxBytes()).Hash(ctx.BlockHeight())
	ethHash := common.BytesToHash(txHash)

	st = types.StateTransition{
		AccountNonce: msg.Data.AccountNonce,
		Price:        msg.Data.Price,
		GasLimit:     msg.Data.GasLimit,
		Recipient:    msg.Data.Recipient,
		Amount:       msg.Data.Amount,
		Payload:      msg.Data.Payload,
		Csdb:         types.CreateEmptyCommitStateDB(k.GenerateCSDBParams(), *ctx),
		ChainID:      chainIDEpoch,
		TxHash:       &ethHash,
		Sender:       sender,
		Simulate:     ctx.IsCheckTx(),
		TraceTx:      ctx.IsTraceTx(),
		TraceTxLog:   ctx.IsTraceTxLog(),
	}

	return
}