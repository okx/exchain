package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/okex/okchain/x/common/proto"
	"github.com/okex/okchain/x/upgrade/types"
)

// GenesisState - all upgrade state that must be provided at genesis
type GenesisState struct {
	GenesisVersion VersionInfo         `json:"genesis_version"`
	Params         types.UpgradeParams `json:"params"`
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{
		NewVersionInfo(proto.DefaultUpgradeConfig("https://github.com/okex/okchain/releases/tag/v"+""), true),
		types.DefaultParams(),
	}
}

// the method not invoked in v0.6 branch
// InitGenesis - build the genesis version For first Version
func InitGenesis(ctx sdk.Context, k Keeper, data GenesisState) {
	genesisVersion := data.GenesisVersion
	k.SetParams(ctx, data.Params)
	k.AddNewVersionInfo(ctx, genesisVersion)
	k.GetProtocolKeeper().ClearUpgradeConfig(ctx)
	k.GetProtocolKeeper().SetCurrentVersion(ctx, genesisVersion.UpgradeInfo.ProtocolDef.Version)
}

// the method not invoked in v0.6 branch
// WriteGenesis - output genesis parameters
func ExportGenesis(ctx sdk.Context) GenesisState {
	return GenesisState{
		NewVersionInfo(proto.DefaultUpgradeConfig("https://github.com/okex/okchain/releases/tag/v"+version.Version), true),
		types.DefaultParams(),
	}
}
