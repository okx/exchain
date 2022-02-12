package types

import (
	codectypes "github.com/okex/exchain/libs/cosmos-sdk/codec/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/exported"
	"sort"

)

// ClientsConsensusStates defines a slice of ClientConsensusStates that supports the sort interface
type ClientsConsensusStates []ClientConsensusStates

// Len implements sort.Interface
func (ccs ClientsConsensusStates) Len() int { return len(ccs) }

// Less implements sort.Interface
func (ccs ClientsConsensusStates) Less(i, j int) bool { return ccs[i].ClientId < ccs[j].ClientId }

// Swap implements sort.Interface
func (ccs ClientsConsensusStates) Swap(i, j int) { ccs[i], ccs[j] = ccs[j], ccs[i] }

// Sort is a helper function to sort the set of ClientsConsensusStates in place
func (ccs ClientsConsensusStates) Sort() ClientsConsensusStates {
	sort.Sort(ccs)
	return ccs
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (ccs ClientsConsensusStates) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, clientConsensus := range ccs {
		if err := clientConsensus.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}
	return nil
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (ccs ClientConsensusStates) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	for _, consStateWithHeight := range ccs.ConsensusStates {
		if err := consStateWithHeight.UnpackInterfaces(unpacker); err != nil {
			return err
		}
	}
	return nil
}


// UnpackInterfaces implements UnpackInterfacesMesssage.UnpackInterfaces
func (cswh ConsensusStateWithHeight) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return unpacker.UnpackAny(cswh.ConsensusState, new(exported.ConsensusState))
}

// NewClientConsensusStates creates a new ClientConsensusStates instance.
func NewClientConsensusStates(clientID string, consensusStates []ConsensusStateWithHeight) ClientConsensusStates {
	return ClientConsensusStates{
		ClientId:        clientID,
		ConsensusStates: consensusStates,
	}
}


// GetKey returns the key of metadata. Implements exported.GenesisMetadata interface.
func (gm GenesisMetadata) GetKey() []byte {
	return gm.Key
}

// GetValue returns the value of metadata. Implements exported.GenesisMetadata interface.
func (gm GenesisMetadata) GetValue() []byte {
	return gm.Value
}


// NewIdentifiedGenesisMetadata takes in a client ID and list of genesis metadata for that client
// and constructs a new IdentifiedGenesisMetadata.
func NewIdentifiedGenesisMetadata(clientID string, gms []GenesisMetadata) IdentifiedGenesisMetadata {
	return IdentifiedGenesisMetadata{
		ClientId:       clientID,
		ClientMetadata: gms,
	}
}


// NewGenesisMetadata is a constructor for GenesisMetadata
func NewGenesisMetadata(key, val []byte) GenesisMetadata {
	return GenesisMetadata{
		Key:   key,
		Value: val,
	}
}

