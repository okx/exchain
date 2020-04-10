package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/okex/okchain/x/distribution/types"
	"github.com/okex/okchain/x/staking/exported"
)

// GetCdc returns cdc
func (k Keeper) GetCdc() *codec.Codec {
	return k.cdc
}

// GetCodespace returns the name of codespace
func (k Keeper) GetCodespace() sdk.CodespaceType {
	return k.codespace
}

// GetParamSpace returns the name of codespace
func (k Keeper) GetParamSpace() params.Subspace {
	return k.paramSpace
}

// GetStoreKey returns store key
func (k Keeper) GetStoreKey() sdk.StoreKey {
	return k.storeKey
}

// GetStakingKeeper returns staking keeper
func (k Keeper) GetStakingKeeper() types.StakingKeeper {
	return k.stakingKeeper
}

// GetSupplyKeeper returns supply keeper
func (k Keeper) GetSupplyKeeper() types.SupplyKeeper {
	return k.supplyKeeper
}

// GetFeeCollectorName returns the name of fee_collector
func (k Keeper) GetFeeCollectorName() string {
	return k.feeCollectorName
}

// GetBlackListedAddrs returns the back list of address
func (k Keeper) GetBlackListedAddrs() map[string]bool {
	return k.blacklistedAddrs
}

// InitializeValidator initialize validator distribution record
func (k Keeper) InitializeValidator(ctx sdk.Context, val exported.ValidatorI) {
	k.initializeValidator(ctx, val)
}
