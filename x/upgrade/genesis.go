package upgrade

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/okex/okexchain/x/common/proto"
	"github.com/okex/okexchain/x/upgrade/types"
)

// InitGenesis builds the genesis version for first version
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	genesisVersion := data.GenesisVersion
	k.SetParams(ctx, data.Params)
	k.AddNewVersionInfo(ctx, genesisVersion)
	k.GetProtocolKeeper().ClearUpgradeConfig(ctx)
	k.GetProtocolKeeper().SetCurrentVersion(ctx, genesisVersion.UpgradeInfo.ProtocolDef.Version)
}

// ExportGenesis outputs genesis state
func ExportGenesis(_ sdk.Context) types.GenesisState {
	return types.GenesisState{
		GenesisVersion: types.NewVersionInfo(
			proto.DefaultUpgradeConfig("https://github.com/okex/okexchain/releases/tag/v"+version.Version), true),
		Params: types.DefaultParams(),
	}
}
