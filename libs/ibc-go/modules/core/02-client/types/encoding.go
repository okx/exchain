package types

import (
	"fmt"

	"github.com/okex/exchain/common"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/ibc-go/modules/core/exported"
)

// MustUnmarshalClientState attempts to decode and return an ClientState object from
// raw encoded bytes. It panics on error.
func MustUnmarshalClientState(cdc *codec.MarshalProxy, bz []byte) exported.ClientState {
	clientState, err := UnmarshalClientState(cdc, bz)
	if err != nil {
		panic(fmt.Errorf("failed to decode client state: %w", err))
	}

	return clientState
}

// UnmarshalClientState returns an ClientState interface from raw encoded clientState
// bytes of a Proto-based ClientState type. An error is returned upon decoding
// failure.
func UnmarshalClientState(cdc *codec.MarshalProxy, bz []byte) (exported.ClientState, error) {
	var clientState exported.ClientState
	err := cdc.GetProtocMarshal().UnmarshalInterface(bz, &clientState)
	return clientState, err
}

// MustMarshalClientState attempts to encode an ClientState object and returns the
// raw encoded bytes. It panics on error.
func MustMarshalClientState(cdc *codec.MarshalProxy, clientState exported.ClientState) []byte {
	bz, err := MarshalClientState(cdc, clientState)
	if err != nil {
		panic(fmt.Errorf("failed to encode client state: %w", err))
	}

	return bz
}

// MarshalClientState protobuf serializes an ClientState interface
func MarshalClientState(cdc *codec.MarshalProxy, clientStateI exported.ClientState) ([]byte, error) {
	return cdc.GetProtocMarshal().MarshalInterface(clientStateI)
}

// MustUnmarshalConsensusState attempts to decode and return an ConsensusState object from
// raw encoded bytes. It panics on error.
func MustUnmarshalConsensusState(cdc *codec.MarshalProxy, bz []byte) exported.ConsensusState {
	consensusState, err := UnmarshalConsensusState(cdc, bz)
	if err != nil {
		panic(fmt.Errorf("failed to decode consensus state: %w", err))
	}

	return consensusState
}

// UnmarshalConsensusState returns a ConsensusState interface from raw encoded consensus state
// bytes of a Proto-based ConsensusState type. An error is returned upon decoding
// failure.
func UnmarshalConsensusState(cdc *codec.MarshalProxy, bz []byte) (exported.ConsensusState, error) {
	var consensusState exported.ConsensusState
	if err := cdc.UnMarshal(bz, &consensusState); err != nil {
		return nil, err
	}

	return consensusState, nil
}

// MustMarshalConsensusState attempts to encode a ConsensusState object and returns the
// raw encoded bytes. It panics on error.
func MustMarshalConsensusState(cdc *codec.MarshalProxy, consensusState exported.ConsensusState) []byte {
	bz, err := MarshalConsensusState(cdc, consensusState)
	if err != nil {
		panic(fmt.Errorf("failed to encode consensus state: %w", err))
	}

	return bz
}

// MarshalConsensusState protobuf serializes a ConsensusState interface
func MarshalConsensusState(cdc *codec.MarshalProxy, cs exported.ConsensusState) ([]byte, error) {
	return cdc.GetProtocMarshal().MarshalInterface(cs)
}

// MarshalHeader protobuf serializes a Header interface
func MarshalHeader(cdc *codec.MarshalProxy, h exported.Header) ([]byte, error) {
	return common.DefaultMarshal(cdc, h)
}

// MustMarshalHeader attempts to encode a Header object and returns the
// raw encoded bytes. It panics on error.
func MustMarshalHeader(cdc *codec.MarshalProxy, header exported.Header) []byte {
	bz, err := MarshalHeader(cdc, header)
	if err != nil {
		panic(fmt.Errorf("failed to encode header: %w", err))
	}

	return bz
}

// UnmarshalHeader returns a Header interface from raw proto encoded header bytes.
// An error is returned upon decoding failure.
func UnmarshalHeader(cdc codec.BinaryMarshaler, bz []byte) (exported.Header, error) {
	var header exported.Header
	if err := cdc.UnmarshalInterface(bz, &header); err != nil {
		return nil, err
	}

	return header, nil
}
