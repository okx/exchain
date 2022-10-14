package feesplit

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

	"github.com/okex/exchain/x/feesplit/keeper"
	"github.com/okex/exchain/x/feesplit/types"
)

// InitGenesis import module genesis
func InitGenesis(
	ctx sdk.Context,
	k keeper.Keeper,
	data types.GenesisState,
) {
	k.SetParams(ctx, data.Params)

	for _, feeSplit := range data.FeeSplits {
		contract := feeSplit.ContractAddress
		deployer := feeSplit.DeployerAddress
		withdrawer := feeSplit.WithdrawerAddress

		// Set initial contracts receiving transaction fees
		k.SetFeeSplit(ctx, feeSplit)
		k.SetDeployerMap(ctx, deployer, contract)

		if len(withdrawer) != 0 {
			k.SetWithdrawerMap(ctx, withdrawer, contract)
		}
	}
}

// ExportGenesis export module state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return &types.GenesisState{
		Params:    k.GetParams(ctx),
		FeeSplits: k.GetFeeSplits(ctx),
	}
}
