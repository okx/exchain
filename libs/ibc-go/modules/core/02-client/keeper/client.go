package keeper

import (
	"encoding/hex"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

// CreateClient creates a new client state and populates it with a given consensus
// state as defined in https://github.com/cosmos/ics/tree/master/spec/ics-002-client-semantics#create
func (k Keeper) CreateClient(
	ctx sdk.Context, clientState exported.ClientState, consensusState exported.ConsensusState,
) (string, error) {
	params := k.GetParams(ctx)
	if !params.IsAllowedClient(clientState.ClientType()) {
		return "", sdkerrors.Wrapf(
			types.ErrInvalidClientType,
			"client state type %s is not registered in the allowlist", clientState.ClientType(),
		)
	}

	clientID := k.GenerateClientIdentifier(ctx, clientState.ClientType())

	k.SetClientState(ctx, clientID, clientState)
	k.Logger(ctx).Info("client created at height", "client-id", clientID, "height", clientState.GetLatestHeight().String())

	// verifies initial consensus state against client state and initializes client store with any client-specific metadata
	// e.g. set ProcessedTime in Tendermint clients
	if err := clientState.Initialize(ctx, k.cdc, k.ClientStore(ctx, clientID), consensusState); err != nil {
		return "", err
	}

	// check if consensus state is nil in case the created client is Localhost
	if consensusState != nil {
		k.SetClientConsensusState(ctx, clientID, clientState.GetLatestHeight(), consensusState)
	}

	k.Logger(ctx).Info("client created at height", "client-id", clientID, "height", clientState.GetLatestHeight().String())

	return clientID, nil
}

// UpdateClient updates the consensus state and the state root from a provided header.
func (k Keeper) UpdateClient(ctx sdk.Context, clientID string, header exported.Header) error {
	clientState, found := k.GetClientState(ctx, clientID)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot update client with ID %s", clientID)
	}
	clientStore := k.ClientStore(ctx, clientID)
	// prevent update if the client is frozen before or at header height
	if status := clientState.Status(ctx, clientStore, k.cdc); status != exported.Active {
		return sdkerrors.Wrapf(types.ErrClientNotActive, "cannot update client (%s) with status %s", clientID, status)
	}

	eventType := types.EventTypeUpdateClient

	// Any writes made in CheckHeaderAndUpdateState are persisted on both valid updates and misbehaviour updates.
	// Light client implementations are responsible for writing the correct metadata (if any) in either case.
	newClientState, newConsensusState, err := clientState.CheckHeaderAndUpdateState(ctx, k.cdc, clientStore, header)
	if err != nil {
		return sdkerrors.Wrapf(err, "cannot update client with ID %s", clientID)
	}

	// emit the full header in events
	var (
		headerStr       string
		consensusHeight exported.Height
	)
	if header != nil {
		// Marshal the Header as an Any and encode the resulting bytes to hex.
		// This prevents the event value from containing invalid UTF-8 characters
		// which may cause data to be lost when JSON encoding/decoding.
		headerStr = hex.EncodeToString(types.MustMarshalHeader(k.cdc, header))
		// set default consensus height with header height
		consensusHeight = header.GetHeight()

	}

	// set new client state regardless of if update is valid update or misbehaviour
	k.SetClientState(ctx, clientID, newClientState)
	// If client state is not frozen after clientState CheckHeaderAndUpdateState,
	// then update was valid. Write the update state changes, and set new consensus state.
	// Else the update was proof of misbehaviour and we must emit appropriate misbehaviour events.
	if status := newClientState.Status(ctx, clientStore, k.cdc); status != exported.Frozen {
		// if update is not misbehaviour then update the consensus state
		// we don't set consensus state for localhost client
		if header != nil && clientID != exported.Localhost {
			k.SetClientConsensusState(ctx, clientID, header.GetHeight(), newConsensusState)
		} else {
			consensusHeight = types.GetSelfHeight(ctx)
		}

		k.Logger(ctx).Info("client state updated", "client-id", clientID, "height", consensusHeight.String())

	} else {
		// set eventType to SubmitMisbehaviour
		eventType = types.EventTypeSubmitMisbehaviour

		k.Logger(ctx).Info("client frozen due to misbehaviour", "client-id", clientID)
	}
	// emit the full header in events

	// emitting events in the keeper emits for both begin block and handler client updates
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			eventType,
			sdk.NewAttribute(types.AttributeKeyClientID, clientID),
			sdk.NewAttribute(types.AttributeKeyClientType, clientState.ClientType()),
			sdk.NewAttribute(types.AttributeKeyConsensusHeight, consensusHeight.String()),
			sdk.NewAttribute(types.AttributeKeyHeader, headerStr),
		),
	)

	return nil
}

// UpgradeClient upgrades the client to a new client state if this new client was committed to
// by the old client at the specified upgrade height
func (k Keeper) UpgradeClient(ctx sdk.Context, clientID string, upgradedClient exported.ClientState, upgradedConsState exported.ConsensusState,
	proofUpgradeClient, proofUpgradeConsState []byte) error {
	clientState, found := k.GetClientState(ctx, clientID)
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot update client with ID %s", clientID)
	}

	clientStore := k.ClientStore(ctx, clientID)

	if status := clientState.Status(ctx, clientStore, k.cdc); status != exported.Active {
		return sdkerrors.Wrapf(types.ErrClientNotActive, "cannot upgrade client (%s) with status %s", clientID, status)
	}

	updatedClientState, updatedConsState, err := clientState.VerifyUpgradeAndUpdateState(ctx, k.cdc, clientStore,
		upgradedClient, upgradedConsState, proofUpgradeClient, proofUpgradeConsState)
	if err != nil {
		return sdkerrors.Wrapf(err, "cannot upgrade client with ID %s", clientID)
	}

	k.SetClientState(ctx, clientID, updatedClientState)
	k.SetClientConsensusState(ctx, clientID, updatedClientState.GetLatestHeight(), updatedConsState)

	k.Logger(ctx).Info("client state upgraded", "client-id", clientID, "height", updatedClientState.GetLatestHeight().String())

	// emitting events in the keeper emits for client upgrades
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpgradeClient,
			sdk.NewAttribute(types.AttributeKeyClientID, clientID),
			sdk.NewAttribute(types.AttributeKeyClientType, updatedClientState.ClientType()),
			sdk.NewAttribute(types.AttributeKeyConsensusHeight, updatedClientState.GetLatestHeight().String()),
		),
	)

	return nil
}

// CheckMisbehaviourAndUpdateState checks for client misbehaviour and freezes the
// client if so.
func (k Keeper) CheckMisbehaviourAndUpdateState(ctx sdk.Context, misbehaviour exported.Misbehaviour) error {
	clientState, found := k.GetClientState(ctx, misbehaviour.GetClientID())
	if !found {
		return sdkerrors.Wrapf(types.ErrClientNotFound, "cannot check misbehaviour for client with ID %s", misbehaviour.GetClientID())
	}

	clientStore := k.ClientStore(ctx, misbehaviour.GetClientID())

	if status := clientState.Status(ctx, clientStore, k.cdc); status != exported.Active {
		return sdkerrors.Wrapf(types.ErrClientNotActive, "cannot process misbehaviour for client (%s) with status %s", misbehaviour.GetClientID(), status)
	}

	if err := misbehaviour.ValidateBasic(); err != nil {
		return err
	}

	clientState, err := clientState.CheckMisbehaviourAndUpdateState(ctx, k.cdc, clientStore, misbehaviour)
	if err != nil {
		return err
	}

	k.SetClientState(ctx, misbehaviour.GetClientID(), clientState)
	k.Logger(ctx).Info("client frozen due to misbehaviour", "client-id", misbehaviour.GetClientID())

	return nil
}
