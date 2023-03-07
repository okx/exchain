package types

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okx/okbchain/libs/cosmos-sdk/types/errors"
	typestx "github.com/okx/okbchain/libs/cosmos-sdk/types/tx"
	ibctxdecoder "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/ibc-tx"
	authtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okx/okbchain/libs/tendermint/global"
	"github.com/okx/okbchain/libs/tendermint/types"
)

const IGNORE_HEIGHT_CHECKING = -1

// evmDecoder:  MsgEthereumTx decoder by Ethereum RLP
// aminoDecoder: decode bytes to stdTx with amino
// ibcDecoder:  Protobuf decoder

// When and which decoder decoding what kind of tx:
// | ------------| --------------------|
// |             | Tx Type             |
// | ------------|---------------------|
// | evmDecoder  |   evmTx             |
// | aminoDecoder|   stdTx             |
// | ibcDecoder  |   ibcTx             |
// | ------------| --------------------|

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

		for index, f := range []decodeFunc{
			evmDecoder,
			aminoDecoder,
			ibcDecoder,
		} {
			if tx, err = f(cdc, txBytes); err == nil {
				tx.SetRaw(txBytes)
				tx.SetTxHash(types.Tx(txBytes).Hash())
				// index=0 means it is a evmtx(evmDecoder) ,we wont verify again
				// height > IGNORE_HEIGHT_CHECKING means it is a query request
				if index > 0 && height > IGNORE_HEIGHT_CHECKING {
					if sensitive, ok := tx.(sdk.HeightSensitive); ok {
						if err := sensitive.ValidWithHeight(height); err != nil {
							return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
						}
					}
				}

				return tx, nil
			}
		}

		return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
	}
}

// Unmarshaler is a generic type for Unmarshal functions
type Unmarshaler func(bytes []byte, ptr interface{}) error

// 3. Try to decode with protobuf
func ibcDecoder(cdcWrapper codec.CdcAbstraction, bytes []byte) (tx sdk.Tx, err error) {
	simReq := &typestx.SimulateRequest{}
	txBytes := bytes

	err = simReq.Unmarshal(bytes)
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
	tx, err = decode(txBytes)
	if err != nil {
		return nil, fmt.Errorf("IbcTxDecoder decode tx err %v", err)
	}

	return
}

type decodeFunc func(codec.CdcAbstraction, []byte) (sdk.Tx, error)

// 1. Try to decode as MsgEthereumTx by RLP
func evmDecoder(_ codec.CdcAbstraction, txBytes []byte) (tx sdk.Tx, err error) {
	var ethTx MsgEthereumTx
	if err = authtypes.EthereumTxDecode(txBytes, &ethTx); err == nil {
		tx = &ethTx
	}
	return
}

// 2. Try to decode Tx by amino
func aminoDecoder(cdc codec.CdcAbstraction, txBytes []byte) (tx sdk.Tx, err error) {
	var v interface{}
	if v, err = cdc.UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller(txBytes, &tx); err == nil {
		return aminoSanityCheck(v.(sdk.Tx))
	}
	err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	if err != nil {
		return nil, err
	}
	return aminoSanityCheck(tx)
}

func aminoSanityCheck(tx sdk.Tx) (sdk.Tx, error) {
	if tx.GetType() != sdk.StdTxType {
		return nil, fmt.Errorf("amino decode is not allowed for %s", tx.GetType())
	}
	return tx, nil
}
