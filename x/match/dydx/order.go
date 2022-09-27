package dydx

import (
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type Order struct {
	IsBuy          bool
	IsDecreaseOnly bool
	Amount         *big.Float
	LimitPrice     Price
	TriggerPrice   Price
	LimitFee       Fee
	Maker          string
	Taker          string
	Expiration     uint64
	Salt           *big.Int
}

func (order *Order) GetSolFlags() string {
	var booleanFlag = 0
	if order.LimitFee.Value.Sign() == -1 {
		booleanFlag += int(IS_NEGATIVE_LIMIT_FEE)
	}
	if order.IsDecreaseOnly {
		booleanFlag += int(IS_DECREASE_ONLY)
	}
	if order.IsBuy {
		booleanFlag += int(IS_BUY)
	}
	saltBytes := BnToBytes32(order.Salt)
	return "0x" + saltBytes[len(saltBytes)-63:] + strconv.Itoa(booleanFlag)
}

func (order *Order) ToSolidity() *SolOrder {
	var solOrder SolOrder
	solOrder.Flags = order.GetSolFlags()
	solOrder.Amount = order.Amount.Text('f', 0)
	solOrder.LimitPrice = order.LimitPrice.ToSolidity()
	solOrder.TriggerPrice = order.TriggerPrice.ToSolidity()
	solOrder.LimitFee = order.LimitFee.ToSolidity()
	solOrder.Maker = order.Maker
	solOrder.Taker = order.Taker
	solOrder.Expiration = strconv.FormatUint(order.Expiration, 10)
	return &solOrder
}

func (order *Order) ToBytes() string {
	solOrder := order.ToSolidity()
	var args = abi.Arguments{
		{
			Type: SolTyBytes32,
		},
		{
			Type: SolTyUint256,
		},
		{
			Type: SolTyUint256,
		},
		{
			Type: SolTyUint256,
		},
		{
			Type: SolTyUint256,
		},
		{
			Type: SolTyAddress,
		},
		{
			Type: SolTyAddress,
		},
		{
			Type: SolTyUint256,
		},
	}
	bz, err := args.Pack(common.HexToHash(solOrder.Flags),
		string2Int(solOrder.Amount),
		string2Int(solOrder.LimitPrice),
		string2Int(solOrder.TriggerPrice),
		string2Int(solOrder.LimitFee),
		common.HexToAddress(solOrder.Maker),
		common.HexToAddress(solOrder.Taker),
		string2Int(solOrder.Expiration),
	)
	if err != nil {
		panic(err)
	}
	return common.Bytes2Hex(bz)
}

type SolOrder struct {
	Flags        string
	Amount       string
	LimitPrice   string
	TriggerPrice string
	LimitFee     string
	Maker        string
	Taker        string
	Expiration   string
}

type SignedOrder struct {
	Order
	TypedSignature string
}

func string2Int(s string) *big.Int {
	i := new(big.Int)
	_, ok := i.SetString(s, 10)
	if !ok {
		panic("string2Int: invalid string")
	}
	return i
}

func FillToTradeData(order *SignedOrder, amount *big.Int, price Price, fee Fee) string {
	orderData := order.ToBytes()
	signatureData := order.TypedSignature + strings.Repeat("0", 60)
	var args = abi.Arguments{
		{
			Type: SolTyUint256,
		},
		{
			Type: SolTyUint256,
		},
		{
			Type: SolTyUint256,
		},
		{
			Type: SolTyBool,
		},
	}
	fillData, err := args.Pack(
		amount,
		price,
		string2Int(fee.ToSolidity()),
		fee.Value.Sign() == -1,
	)
	if err != nil {
		panic(err)
	}
	return CombineHexString(orderData, common.Bytes2Hex(fillData), signatureData)
}
