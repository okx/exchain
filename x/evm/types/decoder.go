package types

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	typestx "github.com/okex/exchain/libs/cosmos-sdk/types/tx"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	ibctxdecoder "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibc-tx"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/types"
)

const IGNORE_HEIGHT_CHECKING = -1

// TxDecoder returns an sdk.TxDecoder that can decode both auth.StdTx and
// MsgEthereumTx transactions.
func TxDecoder(cdc *codec.Codec, proxy ...*codec.CodecProxy) sdk.TxDecoder {

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
		proxyCodec := &codec.CodecProxy{}
		if len(proxy) > 0 {
			proxyCodec = proxy[0]
		}

		for _, f := range []decodeFunc{
			evmDecoder,
			ubruDecoder,
			ubDecoder,
			byteTx,
			relayTx,
		} {
			if tx, err = f(cdc, proxyCodec, txBytes, height); err == nil {
				switch realTx := tx.(type) {
				case authtypes.StdTx:
					realTx.Raw = txBytes
					realTx.Hash = types.Tx(txBytes).Hash(height)
					return realTx, nil
				case *MsgEthereumTx:
					realTx.Raw = txBytes
					realTx.Hash = types.Tx(txBytes).Hash(height)
					return realTx, nil
				case *authtypes.IbcTx:
					realTx.Raw = txBytes
					realTx.Hash = types.Tx(txBytes).Hash(height)
					return realTx, nil
				}
			}
		}

		return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}
}

// Unmarshaler is a generic type for Unmarshal functions
type Unmarshaler func(bytes []byte, ptr interface{}) error

var byteTx decodeFunc = func(c *codec.Codec, proxy *codec.CodecProxy, bytes []byte, i int64) (sdk.Tx, error) {
	bw := new(sdk.BytesWrapper)
	txBytes, err := bw.UnmarshalToTx(bytes)
	if nil != err {
		return nil, err
	}
	tt := new(auth.StdTx)
	err = c.UnmarshalJSON(txBytes, &tt)
	if len(tt.GetMsgs()) == 0 {
		return nil, errors.New("asd")
	}
	return *tt, err
}

var relayTx decodeFunc = func(c *codec.Codec, proxy *codec.CodecProxy, bytes []byte, i int64) (sdk.Tx, error) {
	//bytes maybe SimulateRequest or BroadcastTxRequest
	tx := &typestx.Tx{}
	simReq := &typestx.SimulateRequest{}
	txBytes := bytes

	err := simReq.Unmarshal(bytes)
	if err != nil {
		broadcastReq := &typestx.BroadcastTxRequest{}
		err = broadcastReq.Unmarshal(bytes)
		if err != nil {
			return authtypes.StdTx{}, err
		}
	} else {
		tx = simReq.Tx
		txBytes = simReq.TxBytes
	}

	if txBytes == nil && simReq.Tx != nil {
		txBytes, err = proto.Marshal(tx)
		if err != nil {
			return nil, fmt.Errorf("relayTx invalid tx Marshal err %v", err)
		}
	}

	if txBytes == nil {
		return nil, errors.New("relayTx empty txBytes is not allowed")
	}

	if proxy == nil {
		return nil, errors.New("relayTx proxy decoder not provided")
	}
	marshaler := proxy.GetProtocMarshal()
	decode := ibctxdecoder.IbcTxDecoder(marshaler)
	txdata, err := decode(txBytes)
	if err != nil {
		return nil, fmt.Errorf("IbcTxDecoder decode tx err %v", err)
	}

	return txdata, nil
}

// func validateBasicTxMsgs(msgs []ibcsdk.Msg) error {
// 	if len(msgs) == 0 {
// 		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "must contain at least one message")
// 	}

// 	for _, msg := range msgs {
// 		err := msg.ValidateBasic()
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

type decodeFunc func(*codec.Codec, *codec.CodecProxy, []byte, int64) (sdk.Tx, error)

// 1. Try to decode as MsgEthereumTx by RLP
func evmDecoder(_ *codec.Codec, proxy *codec.CodecProxy, txBytes []byte, height int64) (tx sdk.Tx, err error) {

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
func ubruDecoder(cdc *codec.Codec, proxy *codec.CodecProxy, txBytes []byte, height int64) (tx sdk.Tx, err error) {
	var v interface{}
	if v, err = cdc.UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller(txBytes, &tx); err != nil {
		return nil, err
	}
	return sanityCheck(v.(sdk.Tx), height)
}

// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
// 3. the original amino way, decode by reflection.
func ubDecoder(cdc *codec.Codec, proxy *codec.CodecProxy, txBytes []byte, height int64) (tx sdk.Tx, err error) {
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
