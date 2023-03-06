package rest

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/types/rest"
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

	// WithdrawRewardEnabledProposalReq defines a set withdraw reward enabled proposal request body.
	WithdrawRewardEnabledProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Enabled     bool           `json:"enabled" yaml:"enabled"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.SysCoins   `json:"deposit" yaml:"deposit"`
	}

	// RewardTruncatePrecisionProposalReq defines a set reward truncate precision proposal request body.
	RewardTruncatePrecisionProposalReq struct {
		BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

		Title       string         `json:"title" yaml:"title"`
		Description string         `json:"description" yaml:"description"`
		Proposer    sdk.AccAddress `json:"proposer" yaml:"proposer"`
		Deposit     sdk.SysCoins   `json:"deposit" yaml:"deposit"`
		Precision   int64          `json:"precision" yaml:"precision"`
	}
)
