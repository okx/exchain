package v0_8

import (
	"github.com/okex/okchain/x/upgrade/types"
)

// ModuleName is a mark of okchain module
const ModuleName = types.ModuleName

// GenesisState is the strcut of genesis state for migrating
type GenesisState struct {
	GenesisVersion types.VersionInfo `json:"genesis_version"`
}
