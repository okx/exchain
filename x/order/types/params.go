package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/params"
)

// nolint
const (
	// System param
	DefaultOrderExpireBlocks = 259200 // order will be expired after 86400 blocks.
	DefaultMaxDealsPerBlock  = 1000   // deals limit per block

	// Fee param
	DefaultFeeAmountPerBlock = "0.000001" // okt
	DefaultFeeDenomPerBlock  = common.NativeToken
	DefaultFeeRateTrade      = "0.001" // percentage
)

// nolint : Parameter keys
var (
	KeyOrderExpireBlocks = []byte("OrderExpireBlocks")
	KeyMaxDealsPerBlock  = []byte("MaxDealsPerBlock")
	KeyFeePerBlock       = []byte("FeePerBlock")
	KeyTradeFeeRate      = []byte("TradeFeeRate")
	DefaultFeePerBlock   = sdk.NewDecCoinFromDec(DefaultFeeDenomPerBlock, sdk.MustNewDecFromStr(DefaultFeeAmountPerBlock))
)

// nolint
var _ params.ParamSet = &Params{}

// nolint : order parameters
type Params struct {
	OrderExpireBlocks int64       `json:"order_expire_blocks"`
	MaxDealsPerBlock  int64       `json:"max_deals_per_block"`
	FeePerBlock       sdk.DecCoin `json:"fee_per_block"`
	TradeFeeRate      sdk.Dec     `json:"trade_fee_rate"`
}

// ParamKeyTable for auth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyOrderExpireBlocks, &p.OrderExpireBlocks},
		{KeyMaxDealsPerBlock, &p.MaxDealsPerBlock},
		{KeyFeePerBlock, &p.FeePerBlock},
		{KeyTradeFeeRate, &p.TradeFeeRate},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		OrderExpireBlocks: DefaultOrderExpireBlocks,
		MaxDealsPerBlock:  DefaultMaxDealsPerBlock,
		FeePerBlock:       DefaultFeePerBlock,
		TradeFeeRate:      sdk.MustNewDecFromStr(DefaultFeeRateTrade),
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	return fmt.Sprintf(`Order Params:
  OrderExpireBlocks: %d
  MaxDealsPerBlock: %d
  FeePerBlock: %s
  TradeFeeRate: %s`, p.OrderExpireBlocks,
		p.MaxDealsPerBlock, p.FeePerBlock,
		p.TradeFeeRate)
}
