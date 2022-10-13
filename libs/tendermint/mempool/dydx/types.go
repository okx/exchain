package dydx

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/libs/dydx/contracts"
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
	EIP712_DOMAIN_SEPARATOR_SCHEMA = "EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"
	EIP712_DOMAIN_NAME             = "P1Orders"
	EIP712_DOMAIN_VERSION          = "1.0"
	EIP712_ORDER_STRUCT_SCHEMA     = "Order(bytes32 flags,uint256 amount,uint256 limitPrice,uint256 triggerPrice,uint256 limitFee,address maker,address taker,uint256 expiration)"

	//TODO, mock addr
	contractAddress = "f1730217Bd65f86D2F008f1821D8Ca9A26d64619"
	KeySize         = sha256.Size
)

var (
	zeroOrderHash = common.Hash{}
	ZeroKey       = [KeySize]byte{}
	callTypeABI   = abi.MustNewType("int32")
	orderTuple    = abi.MustNewType("tuple(bytes32 flags, uint256 amount, uint256 limitprice, uint256 triggerprice, uint256 limitfee, address maker, address taker, uint256 expiration)")
	signedTuple   = abi.MustNewType("tuple(bytes msg, bytes32 sig)")

	//TODO: get chainID
	chainID = big.NewInt(65)

	EIP191_HEADER                       = []byte{0x19, 0x01}
	EIP712_DOMAIN_SEPARATOR_SCHEMA_HASH = crypto.Keccak256Hash([]byte(EIP712_DOMAIN_SEPARATOR_SCHEMA))
	EIP712_DOMAIN_NAME_HASH             = crypto.Keccak256Hash([]byte(EIP712_DOMAIN_NAME))
	EIP712_DOMAIN_VERSION_HASH          = crypto.Keccak256Hash([]byte(EIP712_DOMAIN_VERSION))
	_EIP712_DOMAIN_HASH_                = common.Hash{}
	EIP712_ORDER_STRUCT_SCHEMA_HASH     = crypto.Keccak256Hash([]byte(EIP712_ORDER_STRUCT_SCHEMA))
)

func init() {
	addr, err := hex.DecodeString(contractAddress)
	if err != nil {
		panic(err)
	}
	chainIDBytes, err := callTypeABI.Encode(chainID)
	if err != nil {
		panic(err)
	}
	_EIP712_DOMAIN_HASH_ = crypto.Keccak256Hash(EIP712_DOMAIN_SEPARATOR_SCHEMA_HASH[:], EIP712_DOMAIN_NAME_HASH[:], EIP712_DOMAIN_VERSION_HASH[:], chainIDBytes, common.LeftPadBytes(addr, 32))

}

func (o OrderRaw) Key() [KeySize]byte {
	return sha256.Sum256(o)
}

type P1Order struct {
	CallType int32
	contracts.P1OrdersOrder
}

//TODO to verify
func (p *P1Order) VerifySignature(sig []byte) error {
	orderHash := p.Hash()
	pub, err := crypto.Ecrecover(orderHash[:], sig)
	if err != nil {
		return err
	}
	if !crypto.VerifySignature(pub, orderHash[:], sig) {
		return fmt.Errorf("invalid signature")
	}
	return nil
}

// Hash returns the EIP712 hash of an order.
//TODO to verify
func (p *P1Order) Hash() common.Hash {
	orderBytes, err := p.encodeOrder()
	if err != nil {
		return common.Hash{}
	}
	structHash := crypto.Keccak256Hash(EIP712_ORDER_STRUCT_SCHEMA_HASH[:], orderBytes[:])
	return crypto.Keccak256Hash(EIP191_HEADER, _EIP712_DOMAIN_HASH_[:], structHash[:])
}

func (p *P1Order) Type() OrderType {
	if p == nil {
		return UnknownOrderType
	}
	if p.Flags[31] == 1 {
		return SellOrderType
	}
	return BuyOrderType
}

func (p *P1Order) Price() *big.Int {
	return p.LimitPrice
}

func (p *P1Order) DecodeFrom(data []byte) error {
	err := callTypeABI.DecodeStruct(data[:32], &p.CallType)
	if err != nil {
		return err
	}
	return orderTuple.DecodeStruct(data[32:], &p.P1OrdersOrder)
}

func (p *P1Order) encodeCallType() ([]byte, error) {
	return callTypeABI.Encode(p.CallType)
}

func (p *P1Order) encodeOrder() ([]byte, error) {
	return orderTuple.Encode(p.P1OrdersOrder)
}

func (p *P1Order) Encode() ([]byte, error) {
	bs1, err := p.encodeCallType()
	if err != nil {
		return nil, err
	}
	bs2, err := p.encodeOrder()
	if err != nil {
		return nil, err
	}
	return append(bs1, bs2...), nil
}

func (p P1Order) clone() P1Order {
	return P1Order{
		CallType: p.CallType,
		P1OrdersOrder: contracts.P1OrdersOrder{
			Amount:       new(big.Int).Set(p.Amount),
			LimitPrice:   new(big.Int).Set(p.LimitPrice),
			TriggerPrice: new(big.Int).Set(p.TriggerPrice),
			LimitFee:     new(big.Int).Set(p.LimitFee),
			Maker:        p.Maker,
			Taker:        p.Taker,
			Expiration:   new(big.Int).Set(p.Expiration),
		},
	}
}

type WrapOrder struct {
	P1Order
	FrozenAmount *big.Int
	LeftAmount   *big.Int
	Raw          OrderRaw
	Sig          []byte
	orderHash    common.Hash
}

func (w *WrapOrder) Hash() common.Hash {
	if w.orderHash == zeroOrderHash {
		w.orderHash = w.P1Order.Hash()
	}
	return w.orderHash
}

func (w *WrapOrder) Frozen(amount *big.Int) {
	w.LeftAmount.Sub(w.LeftAmount, amount)
	w.FrozenAmount.Add(w.FrozenAmount, amount)
}

func (w *WrapOrder) Unfrozen(amount *big.Int) {
	w.LeftAmount.Add(w.LeftAmount, amount)
	w.FrozenAmount.Sub(w.FrozenAmount, amount)
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
