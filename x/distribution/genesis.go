package distribution

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/distribution/types"
)

// InitGenesis sets distribution information for genesis
func InitGenesis(ctx sdk.Context, keeper Keeper, supplyKeeper types.SupplyKeeper, data types.GenesisState) {
	moduleHoldings := sdk.DecCoins{}

	keeper.SetPreviousProposerConsAddr(ctx, data.PreviousProposer)
	keeper.SetWithdrawAddrEnabled(ctx, data.WithdrawAddrEnabled)
	for _, dwi := range data.DelegatorWithdrawInfos {
		keeper.SetDelegatorWithdrawAddr(ctx, dwi.DelegatorAddress, dwi.WithdrawAddress)
	}
	for _, acc := range data.ValidatorAccumulatedCommissions {
		keeper.SetValidatorAccumulatedCommission(ctx, acc.ValidatorAddress, acc.Accumulated)
		moduleHoldings = moduleHoldings.Add(acc.Accumulated)
	}

	// check if the module account exists
	moduleAcc := keeper.GetDistributionAccount(ctx)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}
	moduleHoldingsInt, _ := moduleHoldings.TruncateDecimal()
	if moduleAcc.GetCoins().IsZero() {
		if err := moduleAcc.SetCoins(moduleHoldingsInt); err != nil {
			panic(err)
		}
		supplyKeeper.SetModuleAccount(ctx, moduleAcc)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	withdrawAddrEnabled := keeper.GetWithdrawAddrEnabled(ctx)
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

	return types.NewGenesisState(withdrawAddrEnabled, dwi, pp, acc)
}
