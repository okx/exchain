package v0_9

import (
	"github.com/okex/okchain/x/upgrade/types"
)

// const
const (
	ModuleName = types.ModuleName
)

// GenesisState is the strcut of genesis state for migrating
type GenesisState struct {
	GenesisVersion types.VersionInfo   `json:"genesis_version"`
	Params         types.UpgradeParams `json:"params"`
}
