package v0_8

import (
	"github.com/okex/okchain/x/upgrade/types"
)

const (
	ModuleName = types.ModuleName
)

type (
	GenesisState struct {
		GenesisVersion types.VersionInfo `json:"genesis_version"`
	}
)
