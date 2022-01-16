package types

import (
	"encoding/json"
	"math/big"

	"github.com/okex/exchain/libs/tendermint/mempool"
)

// Transactions messages must fulfill the Msg
type Msg interface {

	// Return the message type.
	// Must be alphanumeric or empty.
	Route() string

	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Type() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() error

	// Get the canonical byte representation of the Msg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigners() []AccAddress
}

//__________________________________________________________

// Transactions objects must fulfill the Tx
type Tx interface {
	// Gets the all the transaction's messages.
	GetMsgs() []Msg

	// ValidateBasic does a simple and lightweight validation check that doesn't
	// require access to any other information.
	ValidateBasic() error

	// Return tx sender and gas price
	GetTxInfo(ctx Context) mempool.ExTxInfo

	// Return tx gas price
	GetGasPrice() *big.Int

	// Return tx call function signature
	GetTxFnSignatureInfo() ([]byte, int)

	// none nil for wrapped tx type, nil for other tx type
	GetPayloadTx() Tx

	// 0 for StdTxType, 1 for wrapped tx, 2 for evm tx
	GetType() TransactionType

	GetPayloadTxBytes() []byte
}

//__________________________________________________________

type TransactionType int
const (
	StdTxType       TransactionType = 0
	WrappedTxType   TransactionType = 1
	EvmTxType       TransactionType = 2
	UnknownType     TransactionType = 3
)

func (t TransactionType) String() (res string) {
	switch t {
	case StdTxType:
		res = "StdTx"
	case WrappedTxType:
		res = "WrappedTx"
	case EvmTxType:
		res = "EvmTx"
	default:
		res = "Unknown"
	}
	return res
}
//__________________________________________________________
// TxDecoder unmarshals transaction bytes
type TxDecoder func(txBytes []byte, height ...int64) (Tx, error)

// TxEncoder marshals transaction to bytes
type TxEncoder func(tx Tx) ([]byte, error)

type ExTxInfo struct {
	Metadata  []byte  `json:"metadata"`  // customized message from the node who signs the tx
	Signature []byte  `json:"signature"` // signature for payload+metadata
	NodeKey   []byte  `json:"nodeKey"`   // pub key of the node who signs the tx
}
type WrappedTxEncoder func(txBytes []byte, info *ExTxInfo, txtype TransactionType) ([]byte, error)

//__________________________________________________________

var _ Msg = (*TestMsg)(nil)

// msg type for testing
type TestMsg struct {
	signers []AccAddress
}

func NewTestMsg(addrs ...AccAddress) *TestMsg {
	return &TestMsg{
		signers: addrs,
	}
}

//nolint
func (msg *TestMsg) Route() string { return "TestMsg" }
func (msg *TestMsg) Type() string  { return "Test message" }
func (msg *TestMsg) GetSignBytes() []byte {
	bz, err := json.Marshal(msg.signers)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(bz)
}
func (msg *TestMsg) ValidateBasic() error { return nil }
func (msg *TestMsg) GetSigners() []AccAddress {
	return msg.signers
}

type TestMsg2 struct {
	Signers []AccAddress
}

func NewTestMsg2(addrs ...AccAddress) *TestMsg2 {
	return &TestMsg2{
		Signers: addrs,
	}
}

//nolint
func (msg TestMsg2) Route() string { return "TestMsg" }
func (msg TestMsg2) Type() string  { return "Test message" }
func (msg TestMsg2) GetSignBytes() []byte {
	bz, err := json.Marshal(msg.Signers)
	if err != nil {
		panic(err)
	}
	return MustSortJSON(bz)
}
func (msg TestMsg2) ValidateBasic() error { return nil }
func (msg TestMsg2) GetSigners() []AccAddress {
	return msg.Signers
}
