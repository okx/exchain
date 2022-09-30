package dydx

import (
	"crypto/sha256"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/umbracle/ethgo/abi"
)

type (
	OrderRaw  []byte
	OrderType int8
)

const (
	UnknownOrderType OrderType = iota
	SellOrderType
	BuyOrderType
)

const (
	contractAddress = "0x5D64795f3f815924E607C7e9651e89Db4Dbddb62"
	KeySize         = sha256.Size
)

var (
	ZeroKey     = [KeySize]byte{}
	orderTuple  = abi.MustNewType("tuple(int32 calltype, bytes32 flags, uint256 amount, uint256 limitprice, uint256 triggerprice, uint256 limitfee, address maker, address taker, uint256 expiration)")
	signedTuple = abi.MustNewType("tuple(bytes msg, bytes32 sig)")
)

func (o OrderRaw) Key() [KeySize]byte {
	return sha256.Sum256(o)
}

type Order struct {
	CallType     int32
	Flags        [32]byte
	Maker        common.Address
	Taker        common.Address
	Amount       *big.Int
	LimitPrice   *big.Int
	TriggerPrice *big.Int
	LimitFee     *big.Int
	Expiration   *big.Int
}

func (o *Order) Type() OrderType {
	//TODO: order type
	if o == nil {
		return UnknownOrderType
	}
	if o.Flags[31] == 1 {
		return SellOrderType
	}
	return BuyOrderType
}

func (o *Order) Price() *big.Int {
	return o.LimitPrice
}

func (o *Order) DecodeFrom(data []byte) error {
	return orderTuple.DecodeStruct(data, o)
}

func (o *Order) clone() Order {
	return Order{
		Amount:       new(big.Int).Set(o.Amount),
		LimitPrice:   new(big.Int).Set(o.LimitPrice),
		TriggerPrice: new(big.Int).Set(o.TriggerPrice),
		LimitFee:     new(big.Int).Set(o.LimitFee),
		Maker:        o.Maker,
		Taker:        o.Taker,
		Expiration:   new(big.Int).Set(o.Expiration),
	}
}

type WrapOrder struct {
	Order
	LeftAmount *big.Int
	Raw        OrderRaw
	Sig        []byte
	OrderKey   [KeySize]byte
}

func (w *WrapOrder) Key() [KeySize]byte {
	if w.OrderKey == ZeroKey {
		w.OrderKey = w.Raw.Key()
	}
	return w.OrderKey
}

type SignedOrder struct {
	Msg []byte
	Sig [32]byte
}

func (s *SignedOrder) DecodeFrom(data []byte) error {
	return signedTuple.DecodeStruct(data, s)
}

type MempoolOrder struct {
	raw    OrderRaw
	height int64

	// ids of peers who've sent us this order (as a map for quick lookups).
	// senders: PeerID -> struct{}
	senders sync.Map
}

func NewMempoolOrder(order OrderRaw, height int64) *MempoolOrder {
	return &MempoolOrder{
		raw:    order,
		height: height,
	}
}

func (m *MempoolOrder) Key() [KeySize]byte {
	return m.raw.Key()
}

func (m *MempoolOrder) Raw() OrderRaw {
	return m.raw
}

func (m *MempoolOrder) Height() int64 {
	return m.height
}

func (m *MempoolOrder) StoreSender(senderID uint16) {
	m.senders.LoadOrStore(senderID, struct{}{})

}

func (m *MempoolOrder) HasSender(senderID uint16) bool {
	_, ok := m.senders.Load(senderID)
	return ok
}

type TxData struct {
	AccountNonce uint64          `json:"nonce"`
	Price        *big.Int        `json:"gasPrice"`
	GasLimit     uint64          `json:"gas"`
	Recipient    *common.Address `json:"to" rlp:"nil"` // nil means contract creation
	Amount       *big.Int        `json:"value"`
	Payload      []byte          `json:"input"`

	// signature values
	V *big.Int `json:"v"`
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

func extractOrder(tx types.Tx) []byte {
	var evmTx TxData
	if err := rlp.DecodeBytes(tx, &evmTx); err != nil {
		return nil
	}
	if evmTx.Recipient.Hex() == contractAddress {
		return evmTx.Payload
	}
	return nil
}
