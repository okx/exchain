package types

import "github.com/okex/okexchain/x/common/proto"

// GenesisState contains all upgrade state that must be provided at genesis
type GenesisState struct {
	GenesisVersion VersionInfo   `json:"genesis_version"`
	Params         UpgradeParams `json:"params"`
}

// DefaultGenesisState returns default raw genesis raw message
func DefaultGenesisState() GenesisState {
	return GenesisState{
		NewVersionInfo(proto.DefaultUpgradeConfig("https://github.com/okex/okexchain/releases/tag/v"), true),
		DefaultParams(),
	}
}
