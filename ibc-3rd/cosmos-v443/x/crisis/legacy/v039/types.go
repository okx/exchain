package v039

import sdk "github.com/okex/exchain/ibc-3rd/cosmos-v443/types"

const (
	ModuleName = "crisis"
)

type (
	GenesisState struct {
		ConstantFee sdk.Coin `json:"constant_fee" yaml:"constant_fee"`
	}
)
