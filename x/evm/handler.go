package evm

import (
	bam "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	common2 "github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/evm/txs"
	"github.com/okex/exchain/x/evm/txs/base"
	"github.com/okex/exchain/x/evm/types"
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

		evmtx, ok := msg.(*types.MsgEthereumTx)
		if ok {
			result, err = handleMsgEthereumTx(ctx, k, evmtx)
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
	case *types.MsgEthereumTx:
		return msg.GetTxFnSignatureInfo()
	default:
		return nil, 0
	}
}

// handleMsgEthereumTx handles an Ethereum specific tx
// 1. txs can be divided into TraceTxLog/CheckTx/DeliverTx
func handleMsgEthereumTx(ctx sdk.Context, k *Keeper, msg *types.MsgEthereumTx) (*sdk.Result, error) {
	txFactory := txs.NewFactory(base.Config{
		Ctx:    ctx,
		Keeper: k,
	})
	tx, err := txFactory.CreateTx()
	if err != nil {
		return nil, err
	}

	// core logical to handle ethereum tx
	return txs.TransitionEvmTx(tx, msg)
}


