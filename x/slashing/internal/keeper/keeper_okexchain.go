package keeper

import (
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/x/slashing/internal/types"
)

func (k Keeper) modifyValidatorStatus(ctx sdk.Context, consAddress sdk.ConsAddress, status types.ValStatus) {
	signingInfo, found := k.GetValidatorSigningInfo(ctx, consAddress)
	if found {
		//update validator status to Created
		signingInfo.ValidatorStatus = status
		k.SetValidatorSigningInfo(ctx, consAddress, signingInfo)
	}
}
