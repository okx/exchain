package types

import (
	"math/big"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/mempool"
)

var (
	_ sdk.Tx = (*WrappedTx)(nil)
)

const (
	EthereumTransaction uint32 = iota
	StdTransaction
)

// WrappedTx wrap the inner tx for more extensive
type WrappedTx struct {
	Inner sdk.Tx `json:"inner"`
	// extra data carried by the tx and have to use another struct to represent this tx for operation
	Extra     []byte `json:"extra"`
	Signature []byte `json:"siganture"`
	NodeKey   []byte `json:"node_key"`
	Type      uint32 `json:"type"`
}

// NewWrappedTx create a new wrapped transaction with tx and type
func NewWrappedTx(tx sdk.Tx, ty uint32) WrappedTx {
	return WrappedTx{
		Inner: tx,
		Type:  ty,
	}
}

// WithSignature  generate the signature in the struct
func (tx WrappedTx) WithSignature(signature, key []byte) WrappedTx {
	tx.Signature = signature
	tx.NodeKey = key
	return tx
}

// GetOriginTx return the origin tx
func (tx WrappedTx) GetOriginTx() sdk.Tx {
	return tx.Inner
}

func (tx *WrappedTx) IsSigned() bool {
	return len(tx.Signature) > 0 && len(tx.NodeKey) > 0
}

// Gets the all the transaction's messages.
func (tx WrappedTx) GetMsgs() []sdk.Msg {
	return tx.Inner.GetMsgs()
}

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx WrappedTx) ValidateBasic() error {
	return tx.Inner.ValidateBasic()
}

// Return tx sender and gas price
func (tx WrappedTx) GetTxInfo(ctx sdk.Context) mempool.ExTxInfo {
	return tx.Inner.GetTxInfo(ctx)
}

// Return tx gas price
func (tx WrappedTx) GetGasPrice() *big.Int {
	return tx.Inner.GetGasPrice()
}

// Return tx call function signature
func (tx WrappedTx) GetTxFnSignatureInfo() ([]byte, int) {
	return tx.Inner.GetTxFnSignatureInfo()
}
