package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/params"
)

const (
	DefaultFeeIssue  = "2500"
	DefaultFeeMint   = "10"
	DefaultFeeBurn   = "10"
	DefaultFeeModify = "0"
	DefaultFeeChown  = "10"
)

var (
	KeyFeeIssue                       = []byte("FeeIssue")
	KeyFeeMint                        = []byte("FeeMint")
	KeyFeeBurn                        = []byte("FeeBurn")
	KeyFeeModify                      = []byte("FeeModify")
	KeyFeeChown                       = []byte("FeeChown")
	KeyCertifiedTokenMinDeposit       = []byte("CertifiedTokenMinDeposit")
	KeyCertifiedTokenMaxDepositPeriod = []byte("CertifiedTokenMaxDepositPeriod")
	KeyCertifiedTokenVotingPeriod     = []byte("CertifiedTokenVotingPeriod")
)

var _ params.ParamSet = &Params{}

// mint parameters
type Params struct {
	FeeIssue  sdk.DecCoin `json:"issue_fee"`
	FeeMint   sdk.DecCoin `json:"mint_fee"`
	FeeBurn   sdk.DecCoin `json:"burn_fee"`
	FeeModify sdk.DecCoin `json:"modify_fee"`
	FeeChown  sdk.DecCoin `json:"transfer_ownership_fee"`

	CertifiedTokenMaxDepositPeriod time.Duration `json:"no_suffix_token_max_deposit_period"`
	CertifiedTokenMinDeposit       sdk.DecCoins  `json:"no_suffix_token_min_deposit"`
	CertifiedTokenVotingPeriod     time.Duration `json:"no_suffix_token_voting_period"`
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
		{KeyFeeIssue, &p.FeeIssue},
		{KeyFeeMint, &p.FeeMint},
		{KeyFeeBurn, &p.FeeBurn},
		{KeyFeeModify, &p.FeeModify},
		{KeyFeeChown, &p.FeeChown},
		{KeyCertifiedTokenMinDeposit, &p.CertifiedTokenMinDeposit},
		{KeyCertifiedTokenMaxDepositPeriod, &p.CertifiedTokenMaxDepositPeriod},
		{KeyCertifiedTokenVotingPeriod, &p.CertifiedTokenVotingPeriod},
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	var minDeposit = sdk.DecCoins{sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100))}

	return Params{
		FeeIssue:  sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeIssue)),
		FeeMint:   sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeMint)),
		FeeBurn:   sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeBurn)),
		FeeModify: sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeModify)),
		FeeChown:  sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeChown)),

		CertifiedTokenMinDeposit:       minDeposit,
		CertifiedTokenMaxDepositPeriod: time.Hour * 24,
		CertifiedTokenVotingPeriod:     time.Hour * 72,
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
	sb.WriteString(fmt.Sprintf("CertifiedTokenMinDeposit: %s\n", p.CertifiedTokenMinDeposit))
	sb.WriteString(fmt.Sprintf("CertifiedTokenMaxDepositPeriod: %s\n", p.CertifiedTokenMaxDepositPeriod))
	sb.WriteString(fmt.Sprintf("CertifiedTokenVotingPeriod: %s\n", p.CertifiedTokenVotingPeriod))

	return sb.String()
}
