package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	icatypes "github.com/okex/exchain/libs/ibc-go/modules/apps/27-interchain-accounts/types"
)

// RegisterInterchainAccount attempts to create a new account using the provided address and
// stores it in state keyed by the provided connection and port identifiers
// If an account for the provided address already exists this function returns early (no-op)
func (k Keeper) RegisterInterchainAccount(ctx sdk.Context, connectionID, controllerPortID string, accAddress sdk.AccAddress) {
	if acc := k.accountKeeper.GetAccount(ctx, accAddress); acc != nil {
		return
	}
	baseAcc := auth.NewBaseAccountWithAddress(accAddress)
	interchainAccount := icatypes.NewAminoInterchainAccount(
		&baseAcc,
		controllerPortID,
	)

	k.accountKeeper.NewAccount(ctx, interchainAccount)
	k.accountKeeper.SetAccount(ctx, interchainAccount)

	k.SetInterchainAccountAddress(ctx, connectionID, controllerPortID, interchainAccount.Address.String())
}
