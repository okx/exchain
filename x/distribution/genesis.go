package distribution

import (
	"fmt"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"

	"github.com/okex/exchain/x/distribution/types"
)

// InitGenesis sets distribution information for genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, supplyKeeper types.SupplyKeeper, data types.GenesisState) {

	keeper.SetFeePool(ctx, data.FeePool)
	keeper.SetParams(ctx, data.Params)
	keeper.SetPreviousProposerConsAddr(ctx, data.PreviousProposer)

	for _, dwi := range data.DelegatorWithdrawInfos {
		keeper.SetDelegatorWithdrawAddr(ctx, dwi.DelegatorAddress, dwi.WithdrawAddress)
	}

	moduleHoldings := sdk.SysCoins{}
	for _, acc := range data.ValidatorAccumulatedCommissions {
		keeper.SetValidatorAccumulatedCommission(ctx, acc.ValidatorAddress, acc.Accumulated)
		moduleHoldings = moduleHoldings.Add(acc.Accumulated...)
	}
	moduleHoldings = moduleHoldings.Add(data.FeePool.CommunityPool...)

	// check if the module account exists
	moduleAcc := keeper.GetDistributionAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	if moduleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(moduleHoldings); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, moduleAcc)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	feePool := keeper.GetFeePool(ctx)
	params := keeper.GetParams(ctx)

	dwi := make([]types.DelegatorWithdrawInfo, 0)
	keeper.IterateDelegatorWithdrawAddrs(ctx, func(del sdk.AccAddress, addr sdk.AccAddress) (stop bool) {
		dwi = append(dwi, types.DelegatorWithdrawInfo{
			DelegatorAddress: del,
			WithdrawAddress:  addr,
		})
		return false
	})
	pp := keeper.GetPreviousProposerConsAddr(ctx)
	acc := make([]types.ValidatorAccumulatedCommissionRecord, 0)
	keeper.IterateValidatorAccumulatedCommissions(ctx,
		func(addr sdk.ValAddress, commission types.ValidatorAccumulatedCommission) (stop bool) {
			acc = append(acc, types.ValidatorAccumulatedCommissionRecord{
				ValidatorAddress: addr,
				Accumulated:      commission,
			})
			return false
		},
	)

	return types.NewGenesisState(params, feePool, dwi, pp, acc)
}
