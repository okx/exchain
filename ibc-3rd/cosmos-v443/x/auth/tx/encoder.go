package tx

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/okex/exchain/ibc-3rd/cosmos-v443/codec"
	sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"
	txtypes "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/tx"
)

// DefaultTxEncoder returns a default protobuf TxEncoder using the provided Marshaler
func DefaultTxEncoder() sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		txWrapper, ok := tx.(*wrapper)
		if !ok {
			return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, tx)
		}

		raw := &txtypes.TxRaw{
			BodyBytes:     txWrapper.getBodyBytes(),
			AuthInfoBytes: txWrapper.getAuthInfoBytes(),
			Signatures:    txWrapper.tx.Signatures,
		}

		return proto.Marshal(raw)
	}
}

// DefaultJSONTxEncoder returns a default protobuf JSON TxEncoder using the provided Marshaler.
func DefaultJSONTxEncoder(cdc codec.ProtoCodecMarshaler) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		txWrapper, ok := tx.(*wrapper)
		if ok {
			return cdc.MarshalJSON(txWrapper.tx)
		}

		protoTx, ok := tx.(*txtypes.Tx)
		if ok {
			return cdc.MarshalJSON(protoTx)
		}

		return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, tx)

	}
}
