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
		logger.Logger = l.With("module", "txdecoder")
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
		// 0. try sdk.WrappedTx
		if tx, err = authtypes.DecodeWrappedTx(txBytes, payloadDecoder, heights...); err == nil {
			return
		} else {
			dumpErr(txBytes, "DecodeWrappedTx", err)
		}

		tx, err = payloadDecoder(txBytes, heights...)
		return
	}
}

type decodeFunc func(*codec.Codec, []byte, int64)(sdk.Tx, error)

func payloadTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte, heights ...int64) (tx sdk.Tx, err error) {
		if len(txBytes) == 0 {
			err = sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
			return
		}
		height := global.GetGlobalHeight()
		if len(heights) >= 1 {
			height = heights[0]
		}

		decoders := []decodeFunc {
			evmDecoder,
			ubruDecoder,
			ubDecoder,
		}

		for _, decoder := range decoders {
			tx, err = decoder(cdc, txBytes, height)
			if err == nil {
				return
			}
		}
		err = sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		return
	}
}

// 1. Try to decode as MsgEthereumTx by RLP
func evmDecoder(_ *codec.Codec, txBytes []byte, height int64) (tx sdk.Tx, err error) {

	// bypass height checking in case of a negative number
	if height >= 0 {
		if !types.HigherThanVenus(height) {
			err = sdkerrors.Wrap(sdkerrors.ErrTxDecode, "lower than Venus")
			return
		}
	}

	var ethTx MsgEthereumTx
	if err = authtypes.EthereumTxDecode(txBytes, &ethTx); err == nil {
		tx = ethTx
		return
	} else {
		dumpErr(txBytes, "decoderEvmtx", err)
	}
	return
}

// 2. try customized unmarshalling implemented by UnmarshalFromAmino. higher performance!
func ubruDecoder(cdc *codec.Codec, txBytes []byte, height int64) (tx sdk.Tx, err error) {
	var v interface{}
	if v, err = cdc.UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller(txBytes, &tx); err == nil {
		return sanityCheck(v.(sdk.Tx), height)
	} else {
		dumpErr(txBytes, "UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller", err)
	}
	return
}

// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
// 3. the original amino way, decode by reflection.
func ubDecoder(cdc *codec.Codec, txBytes []byte, height int64) (tx sdk.Tx, err error) {
	if err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx); err == nil {
		return sanityCheck(tx, height)
	} else {
		dumpErr(txBytes, "UnmarshalBinaryLengthPrefixed", err)
	}
	return
}


func sanityCheck(tx sdk.Tx, height int64) (output sdk.Tx, err error) {
	output = tx
	// bypass height checking in case of a negative number
	if height >= 0 {
		if tx.GetType() == sdk.EvmTxType && types.HigherThanVenus(height) {
			output = nil
			err = sdkerrors.Wrap(sdkerrors.ErrTxDecode, "amino decode is not allowed for MsgEthereumTx")
		}
	}
	return
}


