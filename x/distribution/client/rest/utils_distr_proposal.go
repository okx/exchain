package rest

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
)

type (
	// ChangeDistributionTypeProposalReq defines a change distribution type proposal request body.
	ChangeDistributionTypeProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Type        uint32         `json:"type" yaml:"type"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.SysCoins   `json:"deposit" yaml:"deposit"`
	}
)
