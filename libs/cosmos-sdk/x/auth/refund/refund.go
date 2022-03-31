package refund

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
)

func RefundFees(supplyKeeper types.SupplyKeeper, ctx sdk.Context, acc sdk.AccAddress, refundFees sdk.Coins) error {
	coins := supplyKeeper.GetFee()

	if !refundFees.IsValid() {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "invalid refund fee amount: %s", refundFees)
	}

	// verify the account has enough funds to pay for fees
	fee, hasNeg := coins.SafeSub(refundFees)
	if hasNeg {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds,
			"insufficient funds to refund for fees; %s < %s", coins, refundFees)
	}

	err := supplyKeeper.AddCoins(ctx, acc, refundFees)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}
	err = supplyKeeper.AddCoinsToFeeCollector(ctx, fee)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFunds, err.Error())
	}
	supplyKeeper.ResetFee()

	return nil
}
