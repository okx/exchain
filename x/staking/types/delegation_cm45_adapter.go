package types

import (
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

type CM45Delegation struct {
	DelAddr sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
	ValAddr sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	Shares  sdk.Dec        `json:"shares" yaml:"shares"`
}

func NewCM45Delegation(delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdk.Dec) CM45Delegation {
	return CM45Delegation{
		DelAddr: delAddr,
		ValAddr: valAddr,
		Shares:  shares,
	}
}

type CM45DelegationResp struct {
	Delegation CM45Delegation `json:"delegation" yaml:"delegation"`
	Balance    sdk.DecCoin    `json:"balance" yaml:"balance"`
}

func NewCM45DelegationResp(cm45delegation CM45Delegation, tokens sdk.Dec) CM45DelegationResp {
	return CM45DelegationResp{
		Delegation: cm45delegation,
		Balance:    sdk.NewDecCoinFromDec("okt", tokens),
	}
}

type CM45DelegationResponses struct {
	DelResponses []CM45DelegationResp `json:"delegation_responses"`
}

func NewCM45DelegationResponses(ds []CM45DelegationResp) CM45DelegationResponses {
	return CM45DelegationResponses{
		DelResponses: ds,
	}
}

func FormatCM45DelegationResponses(delegator Delegator) CM45DelegationResponses {
	if delegator.ValidatorAddresses == nil {
		return NewCM45DelegationResponses(make([]CM45DelegationResp, 0))
	}
	delResps := make([]CM45DelegationResp, 0)
	delAddr := delegator.DelegatorAddress
	shares := delegator.Shares
	tokens := delegator.Tokens
	for _, valAddr := range delegator.ValidatorAddresses {
		cm45Delegation := NewCM45Delegation(delAddr, valAddr, shares)
		cm45DelegationResp := NewCM45DelegationResp(cm45Delegation, tokens)
		delResps = append(delResps, cm45DelegationResp)
	}
	return NewCM45DelegationResponses(delResps)
}

type CM45Entry struct {
	CompletionTime time.Time `json:"completion_time"`
	Balance        sdk.Dec   `json:"balance" yaml:"balance"`
}

func NewCM45Entry(ct time.Time, balance sdk.Dec) CM45Entry {
	return CM45Entry{
		CompletionTime: ct,
		Balance:        balance,
	}
}

type CM45UnbondingResp struct {
	DelAddr sdk.AccAddress
	Entries []CM45Entry
}

func NewCM45UnbondingResp(delAddr sdk.AccAddress, entry CM45Entry) CM45UnbondingResp {
	entries := make([]CM45Entry, 0)
	entries = append(entries, entry)
	return CM45UnbondingResp{
		DelAddr: delAddr,
		Entries: entries,
	}
}

type UnbondingResponses struct {
	UR []CM45UnbondingResp `json:"unbonding_responses"`
}

func NewUnbondingResponses(ur []CM45UnbondingResp) UnbondingResponses {
	return UnbondingResponses{
		UR: ur,
	}
}

func FormatCM45UnbondingResponses(ui UndelegationInfo) UnbondingResponses {
	cm45Entry := NewCM45Entry(ui.CompletionTime, ui.Quantity)
	cm45UnbondingResp := NewCM45UnbondingResp(ui.DelegatorAddress, cm45Entry)
	responses := make([]CM45UnbondingResp, 0)
	responses = append(responses, cm45UnbondingResp)
	return NewUnbondingResponses(responses)
}
