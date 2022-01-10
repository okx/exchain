package types

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

// ----------------------------------------------------------------------------
// Auxiliary

func TxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var err error
		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}
		payloadDecoder := payloadTxDecoder(cdc)

		//----------------------------------------------
		//----------------------------------------------
		// 1. try sdk.CheckedTx
		var chkTx sdk.Tx
		if chkTx, err = authtypes.DecodeWrappedTx(txBytes, payloadDecoder); err == nil {
			dumpTxType(chkTx)

			return chkTx, nil
		} else {
			fmt.Printf("DecodeWrappedTx failed:%p %s\n", txBytes, err)
		}

		return payloadDecoder(txBytes)
	}
}

func dumpTxType(tx sdk.Tx)  {
	fmt.Printf("DecodeTx: %T\n", tx)

	//switch t tx.(type) {
	//
	//}
}
func payloadTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (tx sdk.Tx, err error) {
		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}

		//----------------------------------------------
		//----------------------------------------------
		// 2. Try to decode as MsgEthereumTx by RLP
		if tx, err = decoderEvmtx(txBytes); err == nil {
			dumpTxType(tx)

			return tx, nil
		} else {
			fmt.Printf("EthereumTxDecode failed: %p %s\n", txBytes, err)
		}

		//----------------------------------------------
		//----------------------------------------------
		// 3. try other concrete message types registered by MakeTxCodec
		// TODO: switch to UnmarshalBinaryBare on SDK v0.40.0
		if v, err := cdc.UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller(txBytes, &tx); err == nil {

			//debug.PrintStack()
			dumpTxType(v.(sdk.Tx))
			return v.(sdk.Tx), nil
		} else {
			fmt.Printf("UnmarshalBinaryLengthPrefixedWithRegisteredUbmarshaller failed:%p %s\n", txBytes, err)
		}

		//----------------------------------------------
		//----------------------------------------------
		// 4. try others
		if err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx); err == nil {
			dumpTxType(tx)

			//debug.PrintStack()
			return tx, nil
		} else {
			fmt.Printf("UnmarshalBinaryLengthPrefixed failed: %p %s\n", txBytes, err)
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