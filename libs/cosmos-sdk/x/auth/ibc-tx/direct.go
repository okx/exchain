package ibc_tx

import (
	"fmt"
	ibctx "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	signing2 "github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
	signing "github.com/okex/exchain/libs/cosmos-sdk/x/auth/ibcsigning"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

// signModeDirectHandler defines the SIGN_MODE_DIRECT SignModeHandler
type signModeDirectHandler struct{}

var _ signing.SignModeHandler = signModeDirectHandler{}

// DefaultMode implements SignModeHandler.DefaultMode
func (signModeDirectHandler) DefaultMode() signing2.SignMode {
	return signing2.SignMode_SIGN_MODE_DIRECT
}

// Modes implements SignModeHandler.Modes
func (signModeDirectHandler) Modes() []signing2.SignMode {
	return []signing2.SignMode{signing2.SignMode_SIGN_MODE_DIRECT}
}

// GetSignBytes implements SignModeHandler.GetSignBytes
func (signModeDirectHandler) GetSignBytes(mode signing2.SignMode, data signing.SignerData, tx ibctx.Tx) ([]byte, error) {
	if mode != signing2.SignMode_SIGN_MODE_DIRECT {
		return nil, fmt.Errorf("expected %s, got %s", signing2.SignMode_SIGN_MODE_DIRECT, mode)
	}

	protoTx, ok := tx.(*wrapper)
	if !ok {
		return nil, fmt.Errorf("can only handle a protobuf Tx, got %T", tx)
	}

	bodyBz := protoTx.getBodyBytes()
	authInfoBz := protoTx.getAuthInfoBytes()

	return DirectSignBytes(bodyBz, authInfoBz, data.ChainID, data.AccountNumber)
}

// DirectSignBytes returns the SIGN_MODE_DIRECT sign bytes for the provided TxBody bytes, AuthInfo bytes, chain ID,
// account number and sequence.
func DirectSignBytes(bodyBytes, authInfoBytes []byte, chainID string, accnum uint64) ([]byte, error) {
	signDoc := types.SignDoc{
		BodyBytes:     bodyBytes,
		AuthInfoBytes: authInfoBytes,
		ChainId:       chainID,
		AccountNumber: accnum,
	}
	return signDoc.Marshal()
}
