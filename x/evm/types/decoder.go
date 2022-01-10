package types

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
	"sync"
)

var logger decoderLogger
var loggerOnce sync.Once
func SetLogger(l log.Logger) {
	loggerOnce.Do(func() {
		logger.Logger = l.With("module", "main")
	})
}

type decoderLogger struct {
	log.Logger
}

func (l decoderLogger) Info(msg string, keyvals ...interface{}) {
	if l.Logger == nil {
		return
	}
	l.Logger.Info(msg, keyvals...)
}

func dumpTxType(tx sdk.Tx, txBytes []byte)  {
	var msgType string
	if tx != nil {
		if len(tx.GetMsgs()) > 0 {
			msgType = fmt.Sprintf("%T", tx.GetMsgs()[0])
		} else {
			msgType = "empty"
		}
	}

	logger.Info("------> succeeded",
		"tx-type", fmt.Sprintf("%T", tx),
		"msg-type", msgType,
		"address", fmt.Sprintf("%p", txBytes),
		"tx", tx,
	)
}

func dumpErr(txBytes []byte, caller string, err error)  {
	logger.Info("------> failed to attempt",
		caller, err,
		"address", fmt.Sprintf("%p", txBytes),
	)
}

// ----------------------------------------------------------------------------
// Auxiliary

func TxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte, heights ...int64) (tx sdk.Tx, err error) {
		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}

		logger.Info("TxDecoder started", "address", fmt.Sprintf("%p", txBytes))
		//debug.PrintStack()
		defer logger.Info("TxDecoder finished", "address", fmt.Sprintf("%p", txBytes))

		defer func() {
			dumpTxType(tx, txBytes)
		}()

		payloadDecoder := payloadTxDecoder(cdc)

		//----------------------------------------------
		//----------------------------------------------
		// 1. try sdk.WrappedTx
		if tx, err = authtypes.DecodeWrappedTx(txBytes, payloadDecoder, heights...); err == nil {
			return
		} else {
			dumpErr(txBytes, "DecodeWrappedTx", err)
		}

		tx, err = payloadDecoder(txBytes, heights...)
		return
	}
}

func payloadTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte, heights ...int64) (tx sdk.Tx, err error) {
		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}

		var height int64
		if len(heights) >= 1 {
			height = heights[0]
		} else {
			height = global.GetGlobalHeight()
		}

		if types.HigherThanVenus(height) {
			//----------------------------------------------
			//----------------------------------------------
			// 2. Try to decode as MsgEthereumTx by RLP
			if tx, err = decoderEvmtx(txBytes); err == nil {
				return
			} else {
				dumpErr(txBytes, "decoderEvmtx", err)
			}
		}
		//----------------------------------------------
		//----------------------------------------------
		// 3. try other concrete message types registered by MakeTxCodec
		// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
		var v interface{}
		if v, err = cdc.UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller(txBytes, &tx); err == nil {
			tx = v.(sdk.Tx)
			return
		} else {
			dumpErr(txBytes, "UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller", err)
		}

		//----------------------------------------------
		//----------------------------------------------
		// 4. try others
		if err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx); err == nil {
			return
		} else {
			dumpErr(txBytes, "UnmarshalBinaryLengthPrefixed", err)
		}
		return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}
}


func decoderEvmtx(txBytes []byte) (sdk.Tx, error) {
	var ethTx MsgEthereumTx
	if err := authtypes.EthereumTxDecode(txBytes, &ethTx); err != nil {
		return nil, err
	}
	return ethTx, nil
}