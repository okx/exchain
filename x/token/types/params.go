package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/params"
)

const (
	DefaultFeeBase   = "0.0125"
	DefaultFeeIssue  = "20000"
	DefaultFeeMint   = "2000"
	DefaultFeeBurn   = "10"
	DefaultFeeModify = "0"
	DefaultFeeSend   = "0"
	// 0.0125 * 0.8
	DefaultFeeMultiSend = "0.01"
	DefaultFeeChown     = "10"
)

var (
	KeyFeeBase      = []byte("FeeBase")
	KeyFeeIssue     = []byte("FeeIssue")
	KeyFeeMint      = []byte("FeeMint")
	KeyFeeBurn      = []byte("FeeBurn")
	KeyFeeModify    = []byte("FeeModify")
	KeyFeeSend      = []byte("FeeSend")
	KeyFeeMultiSend = []byte("FeeMultiSend")
	KeyFeeChown     = []byte("FeeChown")
)

var _ params.ParamSet = &Params{}

// mint parameters
type Params struct {
	FeeBase      sdk.DecCoin `json:"base_fee"`
	FeeIssue     sdk.DecCoin `json:"issue_fee"`
	FeeMint      sdk.DecCoin `json:"mint_fee"`
	FeeBurn      sdk.DecCoin `json:"burn_fee"`
	FeeModify    sdk.DecCoin `json:"modify_fee"`
	FeeSend      sdk.DecCoin `json:"send_fee"`
	FeeMultiSend sdk.DecCoin `json:"multi_send_fee"`
	FeeChown     sdk.DecCoin `json:"transfer_ownership_fee"`
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
		{KeyFeeBase, &p.FeeBase},
		{KeyFeeIssue, &p.FeeIssue},
		{KeyFeeMint, &p.FeeMint},
		{KeyFeeBurn, &p.FeeBurn},
		{KeyFeeModify, &p.FeeModify},
		{KeyFeeSend, &p.FeeSend},
		{KeyFeeMultiSend, &p.FeeMultiSend},
		{KeyFeeChown, &p.FeeChown},
	}
}

// // Equal returns a boolean determining if two Params types are identical.
// func (p Params) Equal(p2 Params) bool {
// 	bz1 := msgCdc.MustMarshalBinaryLengthPrefixed(&p)
// 	bz2 := msgCdc.MustMarshalBinaryLengthPrefixed(&p2)
// 	return bytes.Equal(bz1, bz2)
// }

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		FeeBase:      sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeBase)),
		FeeIssue:     sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeIssue)),
		FeeMint:      sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeMint)),
		FeeBurn:      sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeBurn)),
		FeeModify:    sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeModify)),
		FeeSend:      sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeSend)),
		FeeMultiSend: sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeMultiSend)),
		FeeChown:     sdk.NewDecCoinFromDec(common.NativeToken, sdk.MustNewDecFromStr(DefaultFeeChown)),
	}
}

// String implements the stringer interface.
func (p Params) String() string {
	var sb strings.Builder
	sb.WriteString("Params: \n")
	sb.WriteString(fmt.Sprintf("FeeBase: %s\n", p.FeeBase))
	sb.WriteString(fmt.Sprintf("FeeIssue: %s\n", p.FeeIssue))
	sb.WriteString(fmt.Sprintf("FeeMint: %s\n", p.FeeMint))
	sb.WriteString(fmt.Sprintf("FeeBurn: %s\n", p.FeeBurn))
	sb.WriteString(fmt.Sprintf("FeeModify: %s\n", p.FeeModify))
	sb.WriteString(fmt.Sprintf("FeeSend: %s\n", p.FeeSend))
	sb.WriteString(fmt.Sprintf("FeeMultiSend: %s\n", p.FeeMultiSend))
	sb.WriteString(fmt.Sprintf("FeeChown: %s\n", p.FeeChown))

	return sb.String()
}
