package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	feetypes "github.com/okex/exchain/libs/ibc-go/modules/apps/29-fee/types"
)

func (k Keeper) AllowAutoDispatch(ctx sdk.Context) bool {
	var res bool
	k.paramSpace.Get(ctx, feetypes.KeyAllowAutoDispatch, &res)
	return res
}

func (k Keeper) GetFeePercent(ctx sdk.Context) uint32 {
	var res uint32
	k.paramSpace.Get(ctx, feetypes.KeyFeePercent, &res)
	return res
}

func (k Keeper) GetRecvFeePercent(ctx sdk.Context) uint32 {
	var res uint32
	k.paramSpace.Get(ctx, feetypes.KeyRecvPercent, &res)
	return res
}

func (k Keeper) GetAckFeePercent(ctx sdk.Context) uint32 {
	var res uint32
	k.paramSpace.Get(ctx, feetypes.KeyAckPercent, &res)
	return res
}

func (k Keeper) GetTimeOutFeePercent(ctx sdk.Context) uint32 {
	var res uint32
	k.paramSpace.Get(ctx, feetypes.KeyTimeOutPercent, &res)
	return res
}
