package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/staking/types"
)

func calculateWeight(tokens sdk.Dec) types.Shares {
	return tokens
}

func SimulateWeight(tokens sdk.Dec) types.Shares {
	return calculateWeight(tokens)
}
