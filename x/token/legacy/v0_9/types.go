package v0_9

import (
	"github.com/okex/okchain/x/token/types"
)

const (
	ModuleName = types.ModuleName
)

type (
	// GenesisState - all slashing state that must be provided at genesis
	GenesisState struct {
		Params      types.Params     `json:"params"`
		Tokens      []types.Token    `json:"tokens"`
		LockedAsset []types.AccCoins `json:"lock_asset"`
	}
)
