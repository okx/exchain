package keeper

import (
	"strings"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/x/erc20/types"
)

// GetParams returns the total set of erc20 parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// SetParams sets the erc20 parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetSourceChannelID returns the channel id for an ibc voucher
// The voucher has for format ibc/hash(path)
func (k Keeper) GetSourceChannelID(ctx sdk.Context, ibcVoucherDenom string) (channelID string, err error) {
	hash := strings.Split(ibcVoucherDenom, "/")
	if len(hash) != 2 {
		return "", sdkerrors.Wrapf(types.ErrIbcDenomInvalid, "%s is invalid", ibcVoucherDenom)
	}

	path, err := k.transferKeeper.DenomPathFromHash(ctx, ibcVoucherDenom)
	if err != nil {
		return "", err
	}

	// the path has for format port/channelId
	return strings.Split(path, "/")[1], nil
}
