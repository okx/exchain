package types

import (
	amino "github.com/tendermint/go-amino"

	cryptoamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
)

var cdc = amino.NewCodec()

func init() {
	RegisterBlockAmino(cdc)
}

func RegisterBlockAmino(cdc *amino.Codec) {
	cryptoamino.RegisterAmino(cdc)
	RegisterEvidences(cdc)

	cdc.EnableBufferMarshaler(&Part{})
	cdc.EnableBufferMarshaler(PartSetHeader{})
	cdc.EnableBufferMarshaler(BlockID{})
	cdc.EnableBufferMarshaler(&Vote{})
	cdc.EnableBufferMarshaler(&Proposal{})

	cdc.EnableBufferMarshaler(&DeltaPayload{})
	cdc.EnableBufferMarshaler(&Deltas{})
}

// GetCodec returns a codec used by the package. For testing purposes only.
func GetCodec() *amino.Codec {
	return cdc
}

// For testing purposes only
func RegisterMockEvidencesGlobal() {
	RegisterMockEvidences(cdc)
}
