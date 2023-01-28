package keeper

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	clienttypes "github.com/okex/exchain/libs/ibc-go/modules/core/02-client/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/03-connection/types"
	commitmenttypes "github.com/okex/exchain/libs/ibc-go/modules/core/23-commitment/types"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

func (k Keeper) ConnOpenTryV4(
	ctx sdk.Context,
	counterparty types.Counterparty, // counterpartyConnectionIdentifier, counterpartyPrefix and counterpartyClientIdentifier
	delayPeriod uint64,
	clientID string, // clientID of chainA
	clientState exported.ClientState, // clientState that chainA has for chainB
	counterpartyVersions []exported.Version, // supported versions of chain A
	proofInit []byte, // proof that chainA stored connectionEnd in state (on ConnOpenInit)
	proofClient []byte, // proof that chainA stored a light client of chainB
	proofConsensus []byte, // proof that chainA stored chainB's consensus state at consensus height
	proofHeight exported.Height, // height at which relayer constructs proof of A storing connectionEnd in state
	consensusHeight exported.Height, // latest height of chain B which chain A has stored in its chain B client
) (string, error) {
	// generate a new connection
	connectionID := k.GenerateConnectionIdentifier(ctx)

	selfHeight := clienttypes.GetSelfHeight(ctx)
	if consensusHeight.GTE(selfHeight) {
		return "", sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"consensus height is greater than or equal to the current block height (%s >= %s)", consensusHeight, selfHeight,
		)
	}

	// validate client parameters of a chainB client stored on chainA
	if err := k.clientKeeper.ValidateSelfClient(ctx, clientState); err != nil {
		return "", err
	}

	expectedConsensusState, err := k.clientKeeper.GetSelfConsensusStateV4(ctx, consensusHeight)
	if err != nil {
		return "", sdkerrors.Wrapf(err, "self consensus state not found for height %s", consensusHeight.String())
	}

	// expectedConnection defines Chain A's ConnectionEnd
	// NOTE: chain A's counterparty is chain B (i.e where this code is executed)
	// NOTE: chainA and chainB must have the same delay period
	prefix := k.GetCommitmentPrefix()
	expectedCounterparty := types.NewCounterparty(clientID, "", commitmenttypes.NewMerklePrefix(prefix.Bytes()))
	expectedConnection := types.NewConnectionEnd(types.INIT, counterparty.ClientId, expectedCounterparty, types.ExportedVersionsToProto(counterpartyVersions), delayPeriod)

	// chain B picks a version from Chain A's available versions that is compatible
	// with Chain B's supported IBC versions. PickVersion will select the intersection
	// of the supported versions and the counterparty versions.
	version, err2 := types.PickVersion(types.GetCompatibleVersions(), counterpartyVersions)
	if err2 != nil {
		return "", err2
	}

	// connection defines chain B's ConnectionEnd
	connection := types.NewConnectionEnd(types.TRYOPEN, clientID, counterparty, []*types.Version{version}, delayPeriod)

	// Check that ChainA committed expectedConnectionEnd to its state
	if err := k.VerifyConnectionState(
		ctx, connection, proofHeight, proofInit, counterparty.ConnectionId,
		expectedConnection,
	); err != nil {
		return "", err
	}

	// Check that ChainA stored the clientState provided in the msg
	if err := k.VerifyClientState(ctx, connection, proofHeight, proofClient, clientState); err != nil {
		return "", err
	}

	// Check that ChainA stored the correct ConsensusState of chainB at the given consensusHeight
	if err := k.VerifyClientConsensusState(
		ctx, connection, proofHeight, consensusHeight, proofConsensus, expectedConsensusState,
	); err != nil {
		return "", err
	}

	// store connection in chainB state
	if err := k.addConnectionToClient(ctx, clientID, connectionID); err != nil {
		return "", sdkerrors.Wrapf(err, "failed to add connection with ID %s to client with ID %s", connectionID, clientID)
	}

	k.SetConnection(ctx, connectionID, connection)
	k.Logger(ctx).Info("connection state updated", "connection-id", connectionID, "previous-state", "NONE", "new-state", "TRYOPEN")

	EmitConnectionOpenTryEvent(ctx, connectionID, clientID, counterparty)

	return connectionID, nil
}

func (k Keeper) ConnOpenAckV4(
	ctx sdk.Context,
	connectionID string,
	clientState exported.ClientState, // client state for chainA on chainB
	version *types.Version, // version that ChainB chose in ConnOpenTry
	counterpartyConnectionID string,
	proofTry []byte, // proof that connectionEnd was added to ChainB state in ConnOpenTry
	proofClient []byte, // proof of client state on chainB for chainA
	proofConsensus []byte, // proof that chainB has stored ConsensusState of chainA on its client
	proofHeight exported.Height, // height that relayer constructed proofTry
	consensusHeight exported.Height, // latest height of chainA that chainB has stored on its chainA client
) error {
	// Check that chainB client hasn't stored invalid height
	selfHeight := clienttypes.GetSelfHeight(ctx)
	if consensusHeight.GTE(selfHeight) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidHeight,
			"consensus height is greater than or equal to the current block height (%s >= %s)", consensusHeight, selfHeight,
		)
	}

	// Retrieve connection
	connection, found := k.GetConnection(ctx, connectionID)
	if !found {
		return sdkerrors.Wrap(types.ErrConnectionNotFound, connectionID)
	}

	// verify the previously set connection state
	if connection.State != types.INIT {
		return sdkerrors.Wrapf(
			types.ErrInvalidConnectionState,
			"connection state is not INIT (got %s)", connection.State.String(),
		)
	}

	// ensure selected version is supported
	if !types.IsSupportedVersionV4(types.ProtoVersionsToExported(connection.Versions), version) {
		return sdkerrors.Wrapf(
			types.ErrInvalidConnectionState,
			"the counterparty selected version %s is not supported by versions selected on INIT", version,
		)
	}

	// validate client parameters of a chainA client stored on chainB
	if err := k.clientKeeper.ValidateSelfClient(ctx, clientState); err != nil {
		return err
	}

	// Retrieve chainA's consensus state at consensusheight
	expectedConsensusState, err := k.clientKeeper.GetSelfConsensusStateV4(ctx, consensusHeight)
	if err != nil {
		return sdkerrors.Wrapf(err, "self consensus state not found for height %s", consensusHeight.String())
	}

	prefix := k.GetCommitmentPrefix()
	expectedCounterparty := types.NewCounterparty(connection.ClientId, connectionID, commitmenttypes.NewMerklePrefix(prefix.Bytes()))
	expectedConnection := types.NewConnectionEnd(types.TRYOPEN, connection.Counterparty.ClientId, expectedCounterparty, []*types.Version{version}, connection.DelayPeriod)

	// Ensure that ChainB stored expected connectionEnd in its state during ConnOpenTry
	if err := k.VerifyConnectionState(
		ctx, connection, proofHeight, proofTry, counterpartyConnectionID,
		expectedConnection,
	); err != nil {
		return err
	}

	// Check that ChainB stored the clientState provided in the msg
	if err := k.VerifyClientState(ctx, connection, proofHeight, proofClient, clientState); err != nil {
		return err
	}

	// Ensure that ChainB has stored the correct ConsensusState for chainA at the consensusHeight
	if err := k.VerifyClientConsensusState(
		ctx, connection, proofHeight, consensusHeight, proofConsensus, expectedConsensusState,
	); err != nil {
		return err
	}

	k.Logger(ctx).Info("connection state updated", "connection-id", connectionID, "previous-state", "INIT", "new-state", "OPEN")

	// Update connection state to Open
	connection.State = types.OPEN
	connection.Versions = []*types.Version{version}
	connection.Counterparty.ConnectionId = counterpartyConnectionID
	k.SetConnection(ctx, connectionID, connection)

	EmitConnectionOpenAckEvent(ctx, connectionID, connection)

	return nil
}
