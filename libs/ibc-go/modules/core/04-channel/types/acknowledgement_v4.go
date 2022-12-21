package types

import (
	"fmt"

	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

const (
	// ackErrorString defines a string constant included in error acknowledgements
	// NOTE: Changing this const is state machine breaking as acknowledgements are written into state.
	ackErrorString = "error handling packet: see events for details"
)

// NewErrorAcknowledgement returns a new instance of Acknowledgement using an Acknowledgement_Error
// type in the Response field.
// NOTE: Acknowledgements are written into state and thus, changes made to error strings included in packet acknowledgements
// risk an app hash divergence when nodes in a network are running different patch versions of software.
func NewErrorAcknowledgementV4(err error) Acknowledgement {
	// the ABCI code is included in the abcitypes.ResponseDeliverTx hash
	// constructed in Tendermint and is therefore deterministic
	_, code, _ := sdkerrors.ABCIInfo(err, false) // discard non-determinstic codespace and log values

	return Acknowledgement{
		Response: &Acknowledgement_Error{
			Error: fmt.Sprintf("ABCI code: %d: %s", code, ackErrorString),
		},
	}
}
