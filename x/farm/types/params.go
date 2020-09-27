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
	defaultQuoteToken        = "usdk"
	defaultCreatePoolFee     = "0"
	defaultCreatePoolDeposit = "0"
)

// Parameter store keys
var (
	KeyQuoteToken        = []byte("QuoteToken")
	KeyCreatePoolFee     = []byte("CreatePoolFee")
	KeyCreatePoolDeposit = []byte("CreatePoolDeposit")
)

// ParamKeyTable for farm module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for farm at genesis
type Params struct {
	QuoteToken        string      `json:"quote_token"`
	CreatePoolFee     sdk.DecCoin `json:"create_pool_fee"`
	CreatePoolDeposit sdk.DecCoin `json:"create_pool_deposit"`
}

// NewParams creates a new Params object
func NewParams(quoteToken string, createPoolFee sdk.DecCoin, createPoolDeposit sdk.DecCoin) Params {
	return Params{
		QuoteToken:        quoteToken,
		CreatePoolFee:     createPoolFee,
		CreatePoolDeposit: createPoolFee,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf("Params:\nQuoteToken:%s\nCreatePoolFee:%s\nCreatePoolDeposit:%s\n", p.QuoteToken, p.CreatePoolFee, p.CreatePoolDeposit)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyQuoteToken, Value: &p.QuoteToken},
		{Key: KeyCreatePoolFee, Value: &p.CreatePoolFee},
		{Key: KeyCreatePoolDeposit, Value: &p.CreatePoolDeposit},
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	createPoolFee := sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultCreatePoolFee))
	createPoolDeposit := sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultCreatePoolDeposit))
	return NewParams(defaultQuoteToken, createPoolFee, createPoolDeposit)
}
