package types

import (
	"github.com/gogo/protobuf/proto"

	cryptotypes "github.com/okex/exchain/libs/cosmos-sdk/crypto/types"
	costypes "github.com/okex/exchain/libs/cosmos-sdk/types"
	stdtx "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

type (
	// Msg defines the interface a transaction message must fulfill.
	Msg interface {
		proto.Message

		// ValidateBasic does a simple validation check that
		// doesn't require access to any other information.
		ValidateBasic() error

		// Signers returns the addrs of signers that must sign.
		// CONTRACT: All signatures must be present to be valid.
		// CONTRACT: Returns addrs in some deterministic order.
		GetSigners() []costypes.AccAddress
	}

	// Fee defines an interface for an application application-defined concrete
	// transaction type to be able to set and return the transaction fee.
	Fee interface {
		GetGas() uint64
		GetAmount() costypes.CoinAdapters
	}

	// Signature defines an interface for an application application-defined
	// concrete transaction type to be able to set and return transaction signatures.
	Signature interface {
		GetPubKey() cryptotypes.PubKey
		GetSignature() []byte
	}

	// Tx defines the interface a transaction must fulfill.
	Tx interface {
		// Gets the all the transaction's messages.
		GetMsgs() []Msg

		// ValidateBasic does a simple and lightweight validation check that doesn't
		// require access to any other information.
		ValidateBasic() error
	}

	// FeeTx defines the interface to be implemented by Tx to use the FeeDecorators
	FeeTx interface {
		Tx
		GetGas() uint64
		GetFee() costypes.CoinAdapters
		FeePayer() costypes.AccAddress
		FeeGranter() costypes.AccAddress
	}

	// Tx must have GetMemo() method to use ValidateMemoDecorator
	TxWithMemo interface {
		Tx
		GetMemo() string
	}

	// TxWithTimeoutHeight extends the Tx interface by allowing a transaction to
	// set a height timeout.
	TxWithTimeoutHeight interface {
		Tx

		GetTimeoutHeight() uint64
	}
)

// TxDecoder unmarshals transaction bytes
//type TxDecoder func(txBytes []byte) (costypes.Tx, error)
type TxDecoder func(txBytes []byte) (stdtx.StdTx, error)
type IbcTxDecoder func(txBytes []byte) (*stdtx.IbcTx, error)

// TxEncoder marshals transaction to bytes
type TxEncoder func(tx Tx) ([]byte, error)

// MsgTypeURL returns the TypeURL of a `sdk.Msg`.
func MsgTypeURL(msg Msg) string {
	return "/" + proto.MessageName(msg)
}
