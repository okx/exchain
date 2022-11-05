package dydx

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/libs/dydx/contracts"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/pkg/errors"
	"github.com/umbracle/ethgo/abi"
)

type (
	OrderRaw  []byte
	OrderType int8
)

type SignatureType uint8

const (
	NoPrepend   SignatureType = iota // No string was prepended.
	Decimal                          // PREPEND_DEC was prepended.
	Hexadecimal                      // PREPEND_HEX was prepended.
	Invalid                          // Not a valid type. Used for bound-checking.
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
	PREPEND_DEC                    = "\x19Ethereum Signed Message:\n32"
	PREPEND_HEX                    = "\x19Ethereum Signed Message:\n\x20"
	NUM_SIGNATURE_BYTES            = 66

	KeySize = sha256.Size
)

var (
	zeroOrderHash = common.Hash{}
	callTypeABI   = abi.MustNewType("int32")
	orderTuple    = abi.MustNewType("tuple(bytes32 flags, uint256 amount, uint256 limitprice, uint256 triggerprice, uint256 limitfee, address maker, address taker, uint256 expiration)")

	chainID         = big.NewInt(65)
	ContractAddress = Config.P1OrdersContractAddress

	EIP191_HEADER                       = []byte{0x19, 0x01}
	EIP712_DOMAIN_SEPARATOR_SCHEMA_HASH = crypto.Keccak256Hash([]byte(EIP712_DOMAIN_SEPARATOR_SCHEMA))
	EIP712_DOMAIN_NAME_HASH             = crypto.Keccak256Hash([]byte(EIP712_DOMAIN_NAME))
	EIP712_DOMAIN_VERSION_HASH          = crypto.Keccak256Hash([]byte(EIP712_DOMAIN_VERSION))
	EIP712_ORDER_STRUCT_SCHEMA_HASH     = crypto.Keccak256Hash([]byte(EIP712_ORDER_STRUCT_SCHEMA))
	_EIP712_DOMAIN_HASH_                = crypto.Keccak256Hash(EIP712_DOMAIN_SEPARATOR_SCHEMA_HASH[:], EIP712_DOMAIN_NAME_HASH[:], EIP712_DOMAIN_VERSION_HASH[:], common.LeftPadBytes(chainID.Bytes(), 32), common.LeftPadBytes(common.FromHex(ContractAddress), 32))
)

// InitWithChainID uses the chain-id of the node.
func InitWithChainID(id *big.Int) {
	return
	chainID = id
	_EIP712_DOMAIN_HASH_ = crypto.Keccak256Hash(EIP712_DOMAIN_SEPARATOR_SCHEMA_HASH[:], EIP712_DOMAIN_NAME_HASH[:], EIP712_DOMAIN_VERSION_HASH[:], common.LeftPadBytes(chainID.Bytes(), 32), common.LeftPadBytes(common.FromHex(ContractAddress), 32))
}

var (
	FlagMaskNull               = byte(0)
	FlagMaskIsBuy              = byte(1)
	FlagMaskIsDecreaseOnly     = byte(1 << 1)
	FlagMaskIsNegativeLimitFee = byte(1 << 2)
)

func (o OrderRaw) Key() [KeySize]byte {
	return sha256.Sum256(o)
}

type P1Order struct {
	CallType int32
	contracts.P1OrdersOrder
}

func (p *P1Order) VerifySignature(sig []byte) error {
	orderHash := p.Hash()
	addr, err := ecrecover(orderHash, sig)
	if err != nil {
		return err
	}
	if addr != p.Maker {
		return ErrInvalidSignature
	}
	return nil
}

func ecrecover(hash common.Hash, sig []byte) (common.Address, error) {
	if len(sig) != NUM_SIGNATURE_BYTES {
		return common.Address{}, ErrInvalidSignature
	}
	sigType := SignatureType(sig[NUM_SIGNATURE_BYTES-1])
	var signedHash common.Hash
	switch sigType {
	case NoPrepend:
		signedHash = hash
	case Decimal:
		signedHash = crypto.Keccak256Hash([]byte(PREPEND_DEC), hash[:])
	case Hexadecimal:
		signedHash = crypto.Keccak256Hash([]byte(PREPEND_HEX), hash[:])
	default:
		return common.Address{}, ErrInvalidSignature
	}

	// sig[NUM_SIGNATURE_BYTES-1] is sigType
	ethsig := make([]byte, NUM_SIGNATURE_BYTES-1)
	copy(ethsig, sig[:NUM_SIGNATURE_BYTES-1])
	// Convert to Ethereum signature format [R || S || V] where V is 0 or 1, from https://github.com/ethereum/go-ethereum/crypto/signature_nocgo.go Sign function
	ethsig[len(ethsig)-1] -= 27

	pub, err := crypto.SigToPub(signedHash[:], ethsig)
	if err != nil {
		return common.Address{}, ErrInvalidSignature
	}

	return crypto.PubkeyToAddress(*pub), nil
}

// Hash returns the EIP712 hash of an order.
func (p *P1Order) Hash() common.Hash {
	orderBytes, err := p.encodeOrder()
	if err != nil {
		return common.Hash{}
	}
	structHash := crypto.Keccak256Hash(EIP712_ORDER_STRUCT_SCHEMA_HASH[:], orderBytes[:])
	return crypto.Keccak256Hash(EIP191_HEADER, _EIP712_DOMAIN_HASH_[:], structHash[:])
}

func (p *P1Order) Hash2(chainId int64, orderContractAddr string) common.Hash {
	orderBytes, err := p.encodeOrder()
	if err != nil {
		return common.Hash{}
	}
	_EIP712_DOMAIN_HASH_ := crypto.Keccak256Hash(EIP712_DOMAIN_SEPARATOR_SCHEMA_HASH[:], EIP712_DOMAIN_NAME_HASH[:], EIP712_DOMAIN_VERSION_HASH[:], common.LeftPadBytes(big.NewInt(chainId).Bytes(), 32), common.LeftPadBytes(common.FromHex(orderContractAddr), 32))
	structHash := crypto.Keccak256Hash(EIP712_ORDER_STRUCT_SCHEMA_HASH[:], orderBytes[:])
	return crypto.Keccak256Hash(EIP191_HEADER, _EIP712_DOMAIN_HASH_[:], structHash[:])
}

func (p *P1Order) Type() OrderType {
	if p == nil {
		return UnknownOrderType
	}
	if p.Flags[31] == 1 {
		return BuyOrderType
	}
	return SellOrderType
}

func (p *P1Order) isBuy() bool {
	if p != nil && p.Flags[31]&FlagMaskIsBuy != FlagMaskNull {
		return true
	}
	return false
}

func (p *P1Order) isDecreaseOnly() bool {
	if p != nil && p.Flags[31]&FlagMaskIsDecreaseOnly != FlagMaskNull {
		return true
	}
	return false
}

func (p *P1Order) isNegativeLimitFee() bool {
	if p != nil && p.Flags[31]&FlagMaskIsNegativeLimitFee != FlagMaskNull {
		return true
	}
	return false
}

func (p *P1Order) Price() *big.Int {
	return p.LimitPrice
}

func (p *P1Order) DecodeFrom(data []byte) error {
	return orderTuple.DecodeStruct(data, &p.P1OrdersOrder)
}

func (p *P1Order) encodeCallType() ([]byte, error) {
	return callTypeABI.Encode(p.CallType)
}

func (p *P1Order) encodeOrder() ([]byte, error) {
	return orderTuple.Encode(p.P1OrdersOrder)
}

func (p *P1Order) Encode() ([]byte, error) {
	return p.encodeOrder()
	//bs1, err := p.encodeCallType()
	//if err != nil {
	//	return nil, err
	//}
	//bs2, err := p.encodeOrder()
	//if err != nil {
	//	return nil, err
	//}
	//return append(bs1, bs2...), nil
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
	sync.RWMutex

	P1Order
	FrozenAmount *big.Int
	LeftAmount   *big.Int
	Raw          OrderRaw
	Sig          []byte
	orderHash    common.Hash
}

func (w *WrapOrder) DecodeFrom(data []byte) error {
	if len(data) < NUM_SIGNATURE_BYTES {
		return ErrInvalidSignedOrder
	}
	err := w.P1Order.DecodeFrom(data[:len(data)-NUM_SIGNATURE_BYTES])
	if err != nil {
		return errors.Wrap(err, ErrInvalidOrder.Error())
	}
	w.LeftAmount = new(big.Int).Set(w.P1Order.Amount)
	w.FrozenAmount = big.NewInt(0)
	w.Sig = data[len(data)-NUM_SIGNATURE_BYTES:]
	w.Raw = data
	return nil
}

func (w *WrapOrder) Hash() common.Hash {
	w.Lock()
	defer w.Unlock()
	if w.orderHash == zeroOrderHash {
		w.orderHash = w.P1Order.Hash()
	}
	return w.orderHash
}

func (w *WrapOrder) GetLimitPrice() *big.Int {
	w.RLock()
	defer w.RUnlock()
	return w.LimitPrice
}

func (w *WrapOrder) GetLeftAmount() *big.Int {
	w.RLock()
	defer w.RUnlock()
	return w.LeftAmount
}

func (w *WrapOrder) LeftAndFrozen() *big.Int {
	w.RLock()
	defer w.RUnlock()
	return new(big.Int).Add(w.LeftAmount, w.FrozenAmount)
}

func (w *WrapOrder) Frozen(amount *big.Int) {
	w.Lock()
	defer w.Unlock()
	w.LeftAmount.Sub(w.LeftAmount, amount)
	w.FrozenAmount.Add(w.FrozenAmount, amount)
	if w.LeftAmount.Sign() < 0 {
		fmt.Println("WrapOrder Frozen error", w.orderHash, w.Amount, w.LeftAmount, w.FrozenAmount, amount)
	}
}

func (w *WrapOrder) Unfrozen(amount *big.Int) {
	w.Lock()
	defer w.Unlock()
	w.LeftAmount.Add(w.LeftAmount, amount)
	w.FrozenAmount.Sub(w.FrozenAmount, amount)
	if w.FrozenAmount.Sign() < 0 {
		fmt.Println("WrapOrder Unfrozen error", w.orderHash, w.Amount, w.LeftAmount, w.FrozenAmount, amount)
	}
}

func (w *WrapOrder) Done(amount *big.Int) {
	w.Lock()
	defer w.Unlock()
	if w.FrozenAmount.Cmp(amount) >= 0 {
		w.FrozenAmount.Sub(w.FrozenAmount, amount)
	} else {
		diff := new(big.Int).Sub(amount, w.FrozenAmount)
		w.FrozenAmount = big.NewInt(0)
		w.LeftAmount.Sub(w.LeftAmount, diff)
	}

	if w.FrozenAmount.Sign() < 0 {
		fmt.Println("WrapOrder Done error", w.orderHash, w.Amount, w.LeftAmount, w.FrozenAmount, amount)
	}
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

func ExtractOrder(tx types.Tx) []byte {
	var evmTx TxData
	if err := rlp.DecodeBytes(tx, &evmTx); err != nil {
		return nil
	}
	if evmTx.Recipient != nil && evmTx.Recipient.Hex() == AddressForOrder {
		return evmTx.Payload
	}
	return nil
}

type FilledP1Order struct {
	Filled *big.Int
	Time   time.Time
	contracts.P1OrdersOrder
}
