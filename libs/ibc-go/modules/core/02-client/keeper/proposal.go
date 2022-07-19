package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/upgrade"
	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

// ClientUpdateProposal will try to update the client with the new header if and only if
// the proposal passes. The localhost client is not allowed to be modified with a proposal.
func (k Keeper) ClientUpdateProposal(ctx sdk.Context, p *types.ClientUpdateProposal) error {
	if p.SubjectClientId == exported.Localhost || p.SubstituteClientId == exported.Localhost {
		return sdkerrors.Wrap(types.ErrInvalidUpdateClientProposal, "cannot update localhost client with proposal")
	}

	subjectClientState, found := k.GetClientState(ctx, p.SubjectClientId)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "subject client with ID %s", p.SubjectClientId)
	}

	subjectClientStore := k.ClientStore(ctx, p.SubjectClientId)

	if status := subjectClientState.Status(ctx, subjectClientStore, k.cdc); status == exported.Active {
		return sdkerrors.Wrap(types.ErrInvalidUpdateClientProposal, "cannot update Active subject client")
	}

	substituteClientState, found := k.GetClientState(ctx, p.SubstituteClientId)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "substitute client with ID %s", p.SubstituteClientId)
	}

	if subjectClientState.GetLatestHeight().GTE(substituteClientState.GetLatestHeight()) {
		return sdkerrors.Wrapf(types.ErrInvalidHeight, "subject client state latest height is greater or equal to substitute client state latest height (%s >= %s)", subjectClientState.GetLatestHeight(), substituteClientState.GetLatestHeight())
	}

	substituteClientStore := k.ClientStore(ctx, p.SubstituteClientId)

	if status := substituteClientState.Status(ctx, substituteClientStore, k.cdc); status != exported.Active {
		return sdkerrors.Wrapf(types.ErrClientNotActive, "substitute client is not Active, status is %s", status)
	}

	clientState, err := subjectClientState.CheckSubstituteAndUpdateState(ctx, k.cdc, subjectClientStore, substituteClientStore, substituteClientState)
	if err != nil {
		return err
	}
	k.SetClientState(ctx, p.SubjectClientId, clientState)

	k.Logger(ctx).Info("client updated after governance proposal passed", "client-id", p.SubjectClientId, "height", clientState.GetLatestHeight().String())

	// emitting events in the keeper for proposal updates to clients
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateClientProposal,
			sdk.NewAttribute(types.AttributeKeySubjectClientID, p.SubjectClientId),
			//Ywmet todo add
			sdk.NewAttribute(types.AttributeKeyClientType, clientState.ClientType()),
			sdk.NewAttribute(types.AttributeKeyConsensusHeight, clientState.GetLatestHeight().String()),
		),
	)

	return nil
}

func (k Keeper) HandleUpgradeProposal(ctx sdk.Context, p *types.UpgradeProposal) error {
	clientState, err := types.UnpackClientState(p.UpgradedClientState)
	if err != nil {
		return sdkerrors.Wrap(err, "could not unpack UpgradedClientState")
	}

	// zero out any custom fields before setting
	cs := clientState.ZeroCustomFields()
	bz, err := types.MarshalClientState(k.cdc, cs)
	if err != nil {
		return sdkerrors.Wrap(err, "could not marshal UpgradedClientState")
	}

	if err := k.upgradeKeeper.ScheduleUpgrade(ctx, upgrade.Plan{
		Name:   p.Plan.Name,
		Time:   p.Plan.Time,
		Height: p.Plan.Height,
		Info:   p.Plan.Info,
	}); err != nil {
		return err
	}

	// sets the new upgraded client in last height committed on this chain is at plan.Height,
	// since the chain will panic at plan.Height and new chain will resume at plan.Height
	return k.upgradeKeeper.SetUpgradedClient(ctx, p.Plan.Height, bz)
}
