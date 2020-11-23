package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/common"
	"github.com/okex/okexchain/x/params"
)

const (
	DefaultFeeIssue  = "2500"
	DefaultFeeMint   = "10"
	DefaultFeeBurn   = "10"
	DefaultFeeModify = "0"
	DefaultFeeChown  = "10"
)

var (
	KeyFeeIssue               = []byte("FeeIssue")
	KeyFeeMint                = []byte("FeeMint")
	KeyFeeBurn                = []byte("FeeBurn")
	KeyFeeModify              = []byte("FeeModify")
	KeyFeeChown               = []byte("FeeChown")
	KeyOwnershipConfirmWindow = []byte("OwnershipConfirmWindow")
)

var _ params.ParamSet = &Params{}

// mint parameters
type Params struct {
	FeeIssue               sdk.SysCoin   `json:"issue_fee"`
	FeeMint                sdk.SysCoin   `json:"mint_fee"`
	FeeBurn                sdk.SysCoin   `json:"burn_fee"`
	FeeModify              sdk.SysCoin   `json:"modify_fee"`
	FeeChown               sdk.SysCoin   `json:"transfer_ownership_fee"`
	OwnershipConfirmWindow time.Duration `json:"ownership_confirm_window"`
}

// ParamKeyTable for auth module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

func validateParams(value interface{}) error {
	return nil
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		{KeyFeeIssue, &p.FeeIssue, common.ValidateSysCoin("issue fee")},
		{KeyFeeMint, &p.FeeMint, common.ValidateSysCoin("mint fee")},
		{KeyFeeBurn, &p.FeeBurn, common.ValidateSysCoin("burn fee")},
		{KeyFeeModify, &p.FeeModify, common.ValidateSysCoin("modify fee")},
		{KeyFeeChown, &p.FeeChown, common.ValidateSysCoin("change ownership fee")},
		{KeyOwnershipConfirmWindow, &p.OwnershipConfirmWindow, common.ValidateDurationPositive("confirm ownership window")},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		FeeIssue:               sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeIssue)),
		FeeMint:                sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeMint)),
		FeeBurn:                sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeBurn)),
		FeeModify:              sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeModify)),
		FeeChown:               sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeChown)),
		OwnershipConfirmWindow: DefaultOwnershipConfirmWindow,
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("FeeIssue: %s\n", p.FeeIssue))
	sb.WriteString(fmt.Sprintf("FeeMint: %s\n", p.FeeMint))
	sb.WriteString(fmt.Sprintf("FeeBurn: %s\n", p.FeeBurn))
	sb.WriteString(fmt.Sprintf("FeeModify: %s\n", p.FeeModify))
	sb.WriteString(fmt.Sprintf("FeeChown: %s\n", p.FeeChown))
	sb.WriteString(fmt.Sprintf("OwnershipConfirmWindow: %s\n", p.OwnershipConfirmWindow))
	return sb.String()
}
