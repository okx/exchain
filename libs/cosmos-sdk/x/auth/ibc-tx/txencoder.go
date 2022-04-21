package ibc_tx

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	ibctx "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

func IbcTxEncoder() ibctx.IBCTxEncoder {
	return func(tx ibctx.Tx) ([]byte, error) {
		txWrapper, ok := tx.(*wrapper)
		if !ok {
			return nil, fmt.Errorf("expected %T, got %T", &wrapper{}, tx)
		}

		raw := &types.TxRaw{
			BodyBytes:     txWrapper.getBodyBytes(),
			AuthInfoBytes: txWrapper.getAuthInfoBytes(),
			Signatures:    txWrapper.tx.Signatures,
		}

		return proto.Marshal(raw)
	}
}
