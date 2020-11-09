package types

import (
	"fmt"
	"time"

	"github.com/okex/okexchain/x/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/params"
)

// Default parameter namespace
const (
	DefaultParamspace                = ModuleName
	defaultQuoteSymbol               = "usdk"
	defaultCreatePoolFee             = "0"
	defaultCreatePoolDeposit         = "10"
	defaultManageWhiteListMinDeposit = "100"
)

// Parameter store keys
var (
	KeyQuoteSymbol                     = []byte("QuoteSymbol")
	KeyCreatePoolFee                   = []byte("CreatePoolFee")
	KeyCreatePoolDeposit               = []byte("CreatePoolDeposit")
	keyManageWhiteListMaxDepositPeriod = []byte("ManageWhiteListMaxDepositPeriod")
	keyManageWhiteListMinDeposit       = []byte("ManageWhiteListMinDeposit")
	keyManageWhiteListVotingPeriod     = []byte("ManageWhiteListVotingPeriod")
	keyYieldNativeToken                = []byte("YieldNativeToken")
)

// ParamKeyTable for farm module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for farm at genesis
type Params struct {
	QuoteSymbol       string      `json:"quote_symbol"`
	CreatePoolFee     sdk.SysCoin `json:"create_pool_fee"`
	CreatePoolDeposit sdk.SysCoin `json:"create_pool_deposit"`
	// proposal params
	ManageWhiteListMaxDepositPeriod time.Duration `json:"manage_white_list_max_deposit_period"`
	ManageWhiteListMinDeposit       sdk.SysCoins  `json:"manage_white_list_min_deposit"`
	ManageWhiteListVotingPeriod     time.Duration `json:"manage_white_list_voting_period"`
	YieldNativeToken                bool          `json:"yield_native_token"`
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf(`Params:
  Quote Symbol:								%s
  Create Pool Fee:							%s
  Create Pool Deposit:						%s
  Manage White List Max Deposit Period:		%s
  Manage White List Min Deposit:			%s
  Manage White List Voting Period:			%s
  Yield Native Token Enabled:               %v`,
		p.QuoteSymbol, p.CreatePoolFee, p.CreatePoolDeposit,
		p.ManageWhiteListMaxDepositPeriod, p.ManageWhiteListMinDeposit, p.ManageWhiteListVotingPeriod,
		p.YieldNativeToken)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{Key: KeyQuoteSymbol, Value: &p.QuoteSymbol},
		{Key: KeyCreatePoolFee, Value: &p.CreatePoolFee},
		{Key: KeyCreatePoolDeposit, Value: &p.CreatePoolDeposit},
		{Key: keyManageWhiteListMaxDepositPeriod, Value: &p.ManageWhiteListMaxDepositPeriod},
		{Key: keyManageWhiteListMinDeposit, Value: &p.ManageWhiteListMinDeposit},
		{Key: keyManageWhiteListVotingPeriod, Value: &p.ManageWhiteListVotingPeriod},
		{Key: keyYieldNativeToken, Value: &p.YieldNativeToken},
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return Params{
		QuoteSymbol:                     defaultQuoteSymbol,
		CreatePoolFee:                   sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultCreatePoolFee)),
		CreatePoolDeposit:               sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultCreatePoolDeposit)),
		ManageWhiteListMaxDepositPeriod: time.Hour * 24,
		ManageWhiteListMinDeposit:       sdk.NewDecCoinsFromDec(common.NativeToken, sdk.MustNewDecFromStr(defaultManageWhiteListMinDeposit)),
		ManageWhiteListVotingPeriod:     time.Hour * 72,
		YieldNativeToken:                false,
	}
}
