package types

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	typestx "github.com/okex/exchain/libs/cosmos-sdk/types/tx"
	ibctxdecoder "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/types"
)

const IGNORE_HEIGHT_CHECKING = -1

// TxDecoder returns an sdk.TxDecoder that can decode both auth.StdTx and
// MsgEthereumTx transactions.
func TxDecoder(cdc codec.CdcAbstraction) sdk.TxDecoder {

	return func(txBytes []byte, heights ...int64) (sdk.Tx, error) {
		if len(heights) > 1 {
			return nil, fmt.Errorf("to many height parameters")
		}
		var tx sdk.Tx
		var err error
		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}

		var height int64
		if len(heights) == 1 {
			height = heights[0]
		} else {
			height = global.GetGlobalHeight()
		}

		for _, f := range []decodeFunc{
			evmDecoder,
			ubruDecoder,
			ubDecoder,
			relayTx,
		} {
			if tx, err = f(cdc, txBytes, height); err == nil {
				tx.SetRaw(txBytes)
				tx.SetTxHash(types.Tx(txBytes).Hash(height))
				return tx, nil
			}
		}

		return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}
}

// Unmarshaler is a generic type for Unmarshal functions
type Unmarshaler func(bytes []byte, ptr interface{}) error

func relayTx(cdcWrapper codec.CdcAbstraction, bytes []byte, i int64) (sdk.Tx, error) {
	simReq := &typestx.SimulateRequest{}
	txBytes := bytes

	err := simReq.Unmarshal(bytes)
	if err == nil && simReq.Tx != nil {
		txBytes, err = proto.Marshal(simReq.Tx)
		if err != nil {
			return nil, fmt.Errorf("relayTx invalid tx Marshal err %v", err)
		}
	}

	if txBytes == nil {
		return nil, errors.New("relayTx empty txBytes is not allowed")
	}

	cdc, ok := cdcWrapper.(*codec.CodecProxy)
	if !ok {
		return nil, errors.New("Invalid cdc abstraction!")
	}
	marshaler := cdc.GetProtocMarshal()
	decode := ibctxdecoder.IbcTxDecoder(marshaler)
	txdata, err := decode(txBytes)
	if err != nil {
		return nil, fmt.Errorf("IbcTxDecoder decode tx err %v", err)
	}

	return txdata, nil
}

type decodeFunc func(codec.CdcAbstraction, []byte, int64) (sdk.Tx, error)

// 1. Try to decode as MsgEthereumTx by RLP
func evmDecoder(_ codec.CdcAbstraction, txBytes []byte, height int64) (tx sdk.Tx, err error) {

	// bypass height checking in case of a negative number
	if height >= 0 && !types.HigherThanVenus(height) {
		err = fmt.Errorf("lower than Venus")
		return
	}

	var ethTx MsgEthereumTx
	if err = authtypes.EthereumTxDecode(txBytes, &ethTx); err == nil {
		tx = &ethTx
	}
	return
}

// 2. try customized unmarshalling implemented by UnmarshalFromAmino. higher performance!
func ubruDecoder(cdc codec.CdcAbstraction, txBytes []byte, height int64) (tx sdk.Tx, err error) {
	var v interface{}
	if v, err = cdc.UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller(txBytes, &tx); err != nil {
		return nil, err
	}
	return sanityCheck(v.(sdk.Tx), height)
}

// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
// 3. the original amino way, decode by reflection.
func ubDecoder(cdc codec.CdcAbstraction, txBytes []byte, height int64) (tx sdk.Tx, err error) {
	err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return nil, err
	}
	return sanityCheck(tx, height)
}

func sanityCheck(tx sdk.Tx, height int64) (sdk.Tx, error) {
	if tx.GetType() == sdk.EvmTxType && types.HigherThanVenus(height) {
		return nil, fmt.Errorf("amino decode is not allowed for MsgEthereumTx")
	}
	return tx, nil
}
