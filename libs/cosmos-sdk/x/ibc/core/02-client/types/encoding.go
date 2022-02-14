package types

import (
	"fmt"
	"github.com/okex/exchain/common"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/exported"
)

// MustUnmarshalClientState attempts to decode and return an ClientState object from
// raw encoded bytes. It panics on error.
func MustUnmarshalClientState(cdc codec.Codec, bz []byte) exported.ClientState {
	clientState, err := UnmarshalClientState(cdc, bz)
	if err != nil {
		panic(fmt.Errorf("failed to decode client state: %w", err))
	}

	return clientState
}



// UnmarshalClientState returns an ClientState interface from raw encoded clientState
// bytes of a Proto-based ClientState type. An error is returned upon decoding
// failure.
func UnmarshalClientState(cdc codec.Codec, bz []byte) (exported.ClientState, error) {
	var clientState exported.ClientState
	if err := cdc.UnmarshalBinaryBare(bz,&clientState); err != nil {
		return nil, err
	}

	return clientState, nil
}




// MustMarshalClientState attempts to encode an ClientState object and returns the
// raw encoded bytes. It panics on error.
func MustMarshalClientState(cdc codec.Codec, clientState exported.ClientState) []byte {
	bz, err := MarshalClientState(cdc, clientState)
	if err != nil {
		panic(fmt.Errorf("failed to encode client state: %w", err))
	}

	return bz
}

// MarshalClientState protobuf serializes an ClientState interface
func MarshalClientState(cdc codec.Codec, clientStateI exported.ClientState) ([]byte, error) {
	return cdc.MarshalBinaryBare(clientStateI)
}




// MustUnmarshalConsensusState attempts to decode and return an ConsensusState object from
// raw encoded bytes. It panics on error.
func MustUnmarshalConsensusState(cdc codec.Codec, bz []byte) exported.ConsensusState {
	consensusState, err := UnmarshalConsensusState(cdc, bz)
	if err != nil {
		panic(fmt.Errorf("failed to decode consensus state: %w", err))
	}

	return consensusState
}


// UnmarshalConsensusState returns a ConsensusState interface from raw encoded consensus state
// bytes of a Proto-based ConsensusState type. An error is returned upon decoding
// failure.
func UnmarshalConsensusState(cdc codec.Codec, bz []byte) (exported.ConsensusState, error) {
	var consensusState exported.ConsensusState
	if err := cdc.UnmarshalBinaryBare(bz, &consensusState); err != nil {
		return nil, err
	}

	return consensusState, nil
}

// MustMarshalConsensusState attempts to encode a ConsensusState object and returns the
// raw encoded bytes. It panics on error.
func MustMarshalConsensusState(cdc codec.Codec, consensusState exported.ConsensusState) []byte {
	bz, err := MarshalConsensusState(cdc, consensusState)
	if err != nil {
		panic(fmt.Errorf("failed to encode consensus state: %w", err))
	}

	return bz
}

// MarshalConsensusState protobuf serializes a ConsensusState interface
func MarshalConsensusState(cdc codec.Codec, cs exported.ConsensusState) ([]byte, error) {
	return cdc.MarshalBinaryBare(cs)
}



// MarshalHeader protobuf serializes a Header interface
func MarshalHeader(cdc codec.Codec, h exported.Header) ([]byte, error) {
	return common.DefaultMarshal(cdc,h)
}

// MustMarshalHeader attempts to encode a Header object and returns the
// raw encoded bytes. It panics on error.
func MustMarshalHeader(cdc codec.Codec, header exported.Header) []byte {
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
