package types

import "math/big"

type Tx interface {
	GetRaw() []byte
	TxHash() []byte
	GetFrom() string
	GetNonce() uint64
	GetGasPrice() *big.Int
}
