package v0_9

import (
	"github.com/okex/okchain/x/upgrade/types"
)

const (
	ModuleName = types.ModuleName
)

type (
	GenesisState struct {
		GenesisVersion types.VersionInfo   `json:"genesis_version"`
		Params         types.UpgradeParams `json:"params"`
	}
)
