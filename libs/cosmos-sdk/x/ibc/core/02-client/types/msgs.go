package types

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	host "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/24-host"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/exported"
)
// message types for the IBC client
const (
	TypeMsgCreateClient       string = "create_client"
	TypeMsgUpdateClient       string = "update_client"
	TypeMsgUpgradeClient      string = "upgrade_client"
	TypeMsgSubmitMisbehaviour string = "submit_misbehaviour"
)

// NewMsgCreateClient creates a new MsgCreateClient instance
//nolint:interfacer
func NewMsgCreateClient(
	clientState exported.ClientState, consensusState exported.ConsensusState, signer sdk.AccAddress,
) (*MsgCreateClient, error) {

	anyClientState, err := PackClientState(clientState)
	if err != nil {
		return nil, err
	}

	anyConsensusState, err := PackConsensusState(consensusState)
	if err != nil {
		return nil, err
	}

	return &MsgCreateClient{
		ClientState:    anyClientState,
		ConsensusState: anyConsensusState,
		Signer:         signer.String(),
	}, nil
}

// Route implements sdk.Msg
func (msg MsgCreateClient) Route() string {
	return host.RouterKey
}

// Type implements sdk.Msg
func (msg MsgCreateClient) Type() string {
	return TypeMsgCreateClient
}

