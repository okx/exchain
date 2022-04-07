package types

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	host "github.com/okex/exchain/libs/ibc-go/modules/core/24-host"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

// UnpackInterfaces implements UnpackInterfacesMesssage.UnpackInterfaces
func (ics IdentifiedClientState) UnpackInterfaces(unpacker types.AnyUnpacker) error {
	return unpacker.UnpackAny(ics.ClientState, new(exported.ClientState))
}

// NewConsensusStateWithHeight creates a new ConsensusStateWithHeight instance
func NewConsensusStateWithHeight(height Height, consensusState exported.ConsensusState) ConsensusStateWithHeight {
	msg, ok := consensusState.(proto.Message)
	if !ok {
		panic(fmt.Errorf("cannot proto marshal %T", consensusState))
	}

	anyConsensusState, err := types.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}

	return ConsensusStateWithHeight{
		Height:         height,
		ConsensusState: anyConsensusState,
	}
}

// NewIdentifiedClientState creates a new IdentifiedClientState instance
func NewIdentifiedClientState(clientID string, clientState exported.ClientState) IdentifiedClientState {
	msg, ok := clientState.(proto.Message)
	if !ok {
		panic(fmt.Errorf("cannot proto marshal %T", clientState))
	}

	anyClientState, err := types.NewAnyWithValue(msg)
	if err != nil {
		panic(err)
	}

	return IdentifiedClientState{
		ClientId:    clientID,
		ClientState: anyClientState,
	}
}

var _ sort.Interface = IdentifiedClientStates{}

// IdentifiedClientStates defines a slice of ClientConsensusStates that supports the sort interface
type IdentifiedClientStates []IdentifiedClientState

// Len implements sort.Interface
func (ics IdentifiedClientStates) Len() int { return len(ics) }

// Less implements sort.Interface
func (ics IdentifiedClientStates) Less(i, j int) bool { return ics[i].ClientId < ics[j].ClientId }

// Swap implements sort.Interface
func (ics IdentifiedClientStates) Swap(i, j int) { ics[i], ics[j] = ics[j], ics[i] }

// Sort is a helper function to sort the set of IdentifiedClientStates in place
func (ics IdentifiedClientStates) Sort() IdentifiedClientStates {
	sort.Sort(ics)
	return ics
}

// IsValidClientID checks if the clientID is valid and can be parsed into the client
// identifier format.
func IsValidClientID(clientID string) bool {
	_, _, err := ParseClientIdentifier(clientID)
	return err == nil
}

// ValidateClientType validates the client type. It cannot be blank or empty. It must be a valid
// client identifier when used with '0' or the maximum uint64 as the sequence.
func ValidateClientType(clientType string) error {
	if strings.TrimSpace(clientType) == "" {
		return sdkerrors.Wrap(ErrInvalidClientType, "client type cannot be blank")
	}

	smallestPossibleClientID := FormatClientIdentifier(clientType, 0)
	largestPossibleClientID := FormatClientIdentifier(clientType, uint64(math.MaxUint64))

	// IsValidClientID will check client type format and if the sequence is a uint64
	if !IsValidClientID(smallestPossibleClientID) {
		return sdkerrors.Wrap(ErrInvalidClientType, "")
	}

	if err := host.ClientIdentifierValidator(smallestPossibleClientID); err != nil {
		return sdkerrors.Wrap(err, "client type results in smallest client identifier being invalid")
	}
	if err := host.ClientIdentifierValidator(largestPossibleClientID); err != nil {
		return sdkerrors.Wrap(err, "client type results in largest client identifier being invalid")
	}

	return nil
}
