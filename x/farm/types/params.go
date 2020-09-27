package types

import (
	"fmt"

	"github.com/okex/okexchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace        = ModuleName
	defaultQuoteSymbol       = "usdk"
	defaultCreatePoolFee     = "0"
	defaultCreatePoolDeposit = "10"
)

// Parameter store keys
var (
	KeyQuoteSymbol       = []byte("QuoteSymbol")
	KeyCreatePoolFee     = []byte("CreatePoolFee")
	KeyCreatePoolDeposit = []byte("CreatePoolDeposit")
)

// ParamKeyTable for farm module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for farm at genesis
type Params struct {
	QuoteSymbol       string      `json:"quote_symbol"`
	CreatePoolFee     sdk.DecCoin `json:"create_pool_fee"`
	CreatePoolDeposit sdk.DecCoin `json:"create_pool_deposit"`
}

// NewParams creates a new Params object
func NewParams(quoteToken string, createPoolFee sdk.DecCoin, createPoolDeposit sdk.DecCoin) Params {
	return Params{
		QuoteSymbol:       quoteToken,
		CreatePoolFee:     createPoolFee,
		CreatePoolDeposit: createPoolDeposit,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf("Params:\nQuoteSymbol:%s\nCreatePoolFee:%s\nCreatePoolDeposit:%s\n", p.QuoteSymbol, p.CreatePoolFee, p.CreatePoolDeposit)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyQuoteSymbol, Value: &p.QuoteSymbol},
		{Key: KeyCreatePoolFee, Value: &p.CreatePoolFee},
		{Key: KeyCreatePoolDeposit, Value: &p.CreatePoolDeposit},
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	createPoolFee := sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultCreatePoolFee))
	createPoolDeposit := sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultCreatePoolDeposit))
	return NewParams(defaultQuoteSymbol, createPoolFee, createPoolDeposit)
}
