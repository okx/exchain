package types

import (
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/exported"
)

// UnpackClientState unpacks an Any into a ClientState. It returns an error if the
// client state can't be unpacked into a ClientState.
func UnpackClientState(any *codectypes.Any) (exported.ClientState, error) {
	if any == nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	clientState, ok := any.GetCachedValue().(exported.ClientState)
	if !ok {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into ClientState %T", any)
	}

	return clientState, nil
}
