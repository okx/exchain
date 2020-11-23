package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/params"
)

// nolint
const (
	// System param
	DefaultOrderExpireBlocks = 259200 // order will be expired after 86400 blocks.
	DefaultMaxDealsPerBlock  = 1000   // deals limit per block

	// Fee param
	DefaultFeeAmountPerBlock     = "0" // okt
	DefaultFeeDenomPerBlock      = common.NativeToken
	DefaultFeeRateTrade          = "0.001" // percentage
	DefaultNewOrderMsgGasUnit    = 40000
	DefaultCancelOrderMsgGasUnit = 30000
)

// nolint : Parameter keys
var (
	KeyOrderExpireBlocks     = []byte("OrderExpireBlocks")
	KeyMaxDealsPerBlock      = []byte("MaxDealsPerBlock")
	KeyFeePerBlock           = []byte("FeePerBlock")
	KeyTradeFeeRate          = []byte("TradeFeeRate")
	KeyNewOrderMsgGasUnit    = []byte("NewOrderMsgGasUnit")
	KeyCancelOrderMsgGasUnit = []byte("CancelOrderMsgGasUnit")
	DefaultFeePerBlock       = sdk.NewDecCoinFromDec(DefaultFeeDenomPerBlock, sdk.MustNewDecFromStr(DefaultFeeAmountPerBlock))
)

// nolint
var _ params.ParamSet = &Params{}

// nolint : order parameters
type Params struct {
	OrderExpireBlocks     int64       `json:"order_expire_blocks"`
	MaxDealsPerBlock      int64       `json:"max_deals_per_block"`
	FeePerBlock           sdk.SysCoin `json:"fee_per_block"`
	TradeFeeRate          sdk.Dec     `json:"trade_fee_rate"`
	NewOrderMsgGasUnit    uint64      `json:"new_order_msg_gas_unit"`
	CancelOrderMsgGasUnit uint64      `json:"cancel_order_msg_gas_unit"`
}

// ParamKeyTable for auth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// TODO: to supplement the validate function for every pair of param
func validateParams(value interface{}) error {
	return nil
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyOrderExpireBlocks, &p.OrderExpireBlocks, common.ValidateInt64Positive("order expire blocks")},
		{KeyMaxDealsPerBlock, &p.MaxDealsPerBlock, common.ValidateInt64Positive("max deals per block")},
		{KeyFeePerBlock, &p.FeePerBlock, common.ValidateSysCoin("fee per block")},
		{KeyTradeFeeRate, &p.TradeFeeRate, common.ValidateRateNotNeg("trade fee rate")},
		{KeyNewOrderMsgGasUnit, &p.NewOrderMsgGasUnit, common.ValidateUint64Positive("new order msg gas unit")},
		{KeyCancelOrderMsgGasUnit, &p.CancelOrderMsgGasUnit, common.ValidateUint64Positive("cancel order msg gas unit")},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		OrderExpireBlocks:     DefaultOrderExpireBlocks,
		MaxDealsPerBlock:      DefaultMaxDealsPerBlock,
		FeePerBlock:           DefaultFeePerBlock,
		TradeFeeRate:          sdk.MustNewDecFromStr(DefaultFeeRateTrade),
		NewOrderMsgGasUnit:    DefaultNewOrderMsgGasUnit,
		CancelOrderMsgGasUnit: DefaultCancelOrderMsgGasUnit,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	return fmt.Sprintf(`Order Params:
  OrderExpireBlocks: %d
  MaxDealsPerBlock: %d
  FeePerBlock: %s
  TradeFeeRate: %s
  NewOrderMsgGasUnit: %d
  CancelOrderMsgGasUnit: %d`, p.OrderExpireBlocks,
		p.MaxDealsPerBlock, p.FeePerBlock,
		p.TradeFeeRate, p.NewOrderMsgGasUnit, p.CancelOrderMsgGasUnit)
}
