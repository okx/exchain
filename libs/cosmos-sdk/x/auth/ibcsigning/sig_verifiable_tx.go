package signing

import (
	types3 "github.com/okex/exchain/libs/cosmos-sdk/crypto/types"
	types2 "github.com/okex/exchain/libs/cosmos-sdk/types"
	txmsg "github.com/okex/exchain/libs/cosmos-sdk/types/ibc-adapter"
	signing2 "github.com/okex/exchain/libs/cosmos-sdk/types/tx/signing"
)

// SigVerifiableTx defines a transaction interface for all signature verification
// handlers.
type SigVerifiableTx interface {
	txmsg.Tx
	GetSigners() []types2.AccAddress
	GetPubKeys() ([]types3.PubKey, error) // If signer already has pubkey in context, this list will have nil in its place
	GetSignaturesV2() ([]signing2.SignatureV2, error)
}

// Tx defines a transaction interface that supports all standard message, signature
// fee, memo, and auxiliary interfaces.
type Tx interface {
	SigVerifiableTx

	txmsg.TxWithMemo
	txmsg.FeeTx
	txmsg.TxWithTimeoutHeight
}
