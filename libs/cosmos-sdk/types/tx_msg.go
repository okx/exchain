package types

import (
	"encoding/json"
	"math/big"
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

	// Return tx gas price
	GetGasPrice() *big.Int

	// Return tx call function signature
	GetTxFnSignatureInfo() ([]byte, int)

	GetPartnerInfo(ctx Context) (string, string)

	GetType() TransactionType

	GetSigners() []AccAddress

	GetGas() uint64

	GetRaw() []byte
	GetFrom() string
	GetNonce() uint64
	TxHash() []byte
	SetRaw([]byte)
	SetTxHash([]byte)
}

type BaseTx struct {
	Raw   []byte
	Hash  []byte
	From  string
	Nonce uint64
}

func (tx *BaseTx) GetMsgs() []Msg                      { return nil }
func (tx *BaseTx) ValidateBasic() error                { return nil }
func (tx *BaseTx) GetGasPrice() *big.Int               { return big.NewInt(0) }
func (tx *BaseTx) GetTxFnSignatureInfo() ([]byte, int) { return nil, 0 }
func (tx *BaseTx) GetType() TransactionType            { return UnknownType }
func (tx *BaseTx) GetSigners() []AccAddress            { return nil }
func (tx *BaseTx) GetGas() uint64                      { return 0 }
func (tx *BaseTx) GetNonce() uint64                    { return tx.Nonce }
func (tx *BaseTx) GetFrom() string                     { return tx.From }
func (tx *BaseTx) GetRaw() []byte                      { return tx.Raw }
func (tx *BaseTx) TxHash() []byte                    { return tx.Hash }
func (tx *BaseTx) GetPartnerInfo(ctx Context) (string, string) {return "", ""}
func (tx *BaseTx) SetRaw(raw []byte)                   { tx.Raw = raw }
func (tx *BaseTx) SetTxHash(hash []byte)               { tx.Hash = hash }

//__________________________________________________________

type TransactionType int

const (
	UnknownType TransactionType = iota
	StdTxType
	EvmTxType
)

func (t TransactionType) String() (res string) {
	switch t {
	case StdTxType:
		res = "StdTx"
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

//__________________________________________________________

var _ Msg = (*TestMsg)(nil)

// TestMsg is msg type for testing
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
