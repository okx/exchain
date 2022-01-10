package types

import (
	"errors"
	"math/big"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/tlv"
	"github.com/okex/exchain/libs/tendermint/mempool"
)

var (
	_ sdk.Tx = (*WrappedTx)(nil)
)

// ExtraType identify the extra type
type ExtraType uint32

// Extra type consts
const (
	CheckedSignatureExtra ExtraType = iota

	CheckedStdTxSignature      uint32 = 0
	CheckedEthereumTxSignature uint32 = 1
)

// WrappedExtra interface for inner data
type WrappedExtra interface {
	Marshal() []byte
	Unmarshal(data []byte) error
}

// RawCheckedSignature used as the checked tx to transfer in confident nodes
type RawCheckedSignature struct {
	Type      uint32 // 0 StdTx 1 EthereumTx
	Signature []byte // node.key signature
	NodeKey   []byte // node public key
}

func NewRawCheckedSignature(ty uint32) *RawCheckedSignature {
	return &RawCheckedSignature{
		Type:      ty,
		Signature: []byte{},
		NodeKey:   []byte{},
	}
}

func (m *RawCheckedSignature) Marshal() []byte {
	buffer := tlv.NewBuffer()
	buffer.WriteUint32(m.Type)
	buffer.Write(m.Signature)
	buffer.Write(m.NodeKey)
	return nil
}
func (m *RawCheckedSignature) Unmarshal(data []byte) error {
	buffer := tlv.With(data)
	ty, t := buffer.Read()
	if t != tlv.Uint32 {
		return errors.New("failed to unmarshal the type")
	}
	m.Type = ty.(uint32)
	sig, t := buffer.Read()
	if t != tlv.Bytes {
		return errors.New("failed to unmarshal the signature of the buffer type")
	}
	m.Signature = sig.([]byte)
	key, t := buffer.Read()
	if t != tlv.Bytes {
		return errors.New("failed to unmarshal the signature of the buffer type")
	}
	m.NodeKey = key.([]byte)
	return nil
}

// RawWrappedExtra as a wrapper to wrap the extra to []byte
type RawWrappedExtra struct {
	Type  ExtraType
	Inner WrappedExtra
}

func (m *RawWrappedExtra) Marshal() []byte {
	buffer := tlv.NewBuffer()
	buffer.WriteUint32(uint32(m.Type))
	buffer.Write(m.Inner.Marshal())
	return buffer.Bytes()
}

func (m *RawWrappedExtra) Unmarshal(data []byte) error {
	buffer := tlv.With(data)
	ty, t := buffer.Read()
	if t != tlv.Uint32 {
		return errors.New("failed to unmarshal the inner raw part of the type")
	}
	m.Type = ExtraType(ty.(uint32)) // FIXME: need a type check
	bd, t := buffer.Read()
	if t != tlv.Bytes {
		return errors.New("failed to unmarshal the inner raw part of the buffer type")
	}
	switch m.Type {
	case CheckedSignatureExtra:
		{
			inner := NewRawCheckedSignature(CheckedEthereumTxSignature)
			err := inner.Unmarshal(bd.([]byte))
			if err != nil {
				return err
			}
			m.Inner = inner
		}
	default:
		return errors.New("unsupported inner type of the wrapped tx ")
	}
	return nil
}

// WrappedTx wrap the inner tx for more extensive
type WrappedTx struct {
	Inner sdk.Tx `json:"inner"`
	Extra []byte `json:"extra"`
}

// WithRawWrappedExtra to generate the carried data
func (tx WrappedTx) WithRawWrappedExtra(ty ExtraType, extra WrappedExtra) WrappedTx {
	raw := RawWrappedExtra{
		Type:  ty,
		Inner: extra,
	}
	tx.Extra = raw.Marshal()
	return tx
}

// GetOriginTx return the origin tx
func (tx WrappedTx) GetOriginTx() sdk.Tx {
	return tx.Inner
}

// Gets the all the transaction's messages.
func (tx *WrappedTx) GetMsgs() []sdk.Msg {
	return tx.Inner.GetMsgs()
}

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx *WrappedTx) ValidateBasic() error {
	return tx.Inner.ValidateBasic()
}

// Return tx sender and gas price
func (tx *WrappedTx) GetTxInfo(ctx sdk.Context) mempool.ExTxInfo {
	return tx.Inner.GetTxInfo(ctx)
}

// Return tx gas price
func (tx *WrappedTx) GetGasPrice() *big.Int {
	return tx.Inner.GetGasPrice()
}

// Return tx call function signature
func (tx *WrappedTx) GetTxFnSignatureInfo() ([]byte, int) {
	return tx.Inner.GetTxFnSignatureInfo()
}

// Return the data carried by multi type Tx
// StdTx, EthereumTx
// the return value is format with tlv, see more details in the libs/tendermint/libs/tlv
func (tx *WrappedTx) GetTxCarriedData() []byte {
	return tx.Extra
}
