package types

import (
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	clienttypes"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/02-client/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/exported"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

// ClientType returns Tendermint
func (ConsensusState) ClientType() string {
	return exported.Tendermint
}

// GetRoot returns the commitment Root for the specific
func (cs ConsensusState) GetRoot() exported.Root {
	return cs.Root
}

// GetTimestamp returns block time in nanoseconds of the header that created consensus state
func (cs ConsensusState) GetTimestamp() uint64 {
	return uint64(cs.Timestamp.UnixNano())
}

// ValidateBasic defines a basic validation for the tendermint consensus state.
// NOTE: ProcessedTimestamp may be zero if this is an initial consensus state passed in by relayer
// as opposed to a consensus state constructed by the chain.
func (cs ConsensusState) ValidateBasic() error {
	if cs.Root.Empty() {
		return sdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "root cannot be empty")
	}
	if err := tmtypes.ValidateHash(cs.NextValidatorsHash); err != nil {
		return sdkerrors.Wrap(err, "next validators hash is invalid")
	}
	if cs.Timestamp.Unix() <= 0 {
		return sdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "timestamp must be a positive Unix time")
	}
	return nil
}
