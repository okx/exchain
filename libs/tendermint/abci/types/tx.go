package types

import "math/big"

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

type TxEssentials interface {
	GetRaw() []byte
	TxHash() []byte
	GetFrom() string
	GetNonce() uint64
	GetGasPrice() *big.Int
	GetType() TransactionType
}

type MockTx struct {
	Raw      []byte
	Hash     []byte
	From     string
	Nonce    uint64
	GasPrice *big.Int
}

func (tx MockTx) GetRaw() []byte {
	return tx.Raw
}

func (tx MockTx) TxHash() []byte {
	return tx.Hash
}

func (tx MockTx) GetFrom() string {
	return tx.From
}

func (tx MockTx) GetNonce() uint64 {
	return tx.Nonce
}

func (tx MockTx) GetGasPrice() *big.Int {
	return tx.GasPrice
}

func (tx MockTx) GetType() TransactionType {
	return UnknownType
}
