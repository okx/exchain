// nolint
package v0_9

import (
	"github.com/okex/okchain/x/dex/types"
)

const (
	ModuleName = types.ModuleName
)

type (
	// GenesisState - all slashing state that must be provided at genesis
	GenesisState struct {
		Params        types.Params        `json:"params"`
		TokenPairs    []*types.TokenPair  `json:"token_pairs"`
		WithdrawInfos types.WithdrawInfos `json:"withdraw_infos"`
	}
)
