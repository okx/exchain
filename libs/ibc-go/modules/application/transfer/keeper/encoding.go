package keeper

import (
	"github.com/okex/exchain/libs/ibc-go/modules/application/transfer/types"
)

// UnmarshalDenomTrace attempts to decode and return an DenomTrace object from
// raw encoded bytes.
func (k Keeper) UnmarshalDenomTrace(bz []byte) (types.DenomTrace, error) {
	var denomTrace types.DenomTrace
	if err := k.cdc.GetProtocMarshal().UnmarshalBinaryBare(bz, &denomTrace); err != nil {
		return types.DenomTrace{}, err
	}
	return denomTrace, nil
}

// MustUnmarshalDenomTrace attempts to decode and return an DenomTrace object from
// raw encoded bytes. It panics on error.
func (k Keeper) MustUnmarshalDenomTrace(bz []byte) types.DenomTrace {
	var denomTrace types.DenomTrace
	k.cdc.GetProtocMarshal().MustUnmarshalBinaryBare(bz, &denomTrace)
	return denomTrace
}

// MarshalDenomTrace attempts to encode an DenomTrace object and returns the
// raw encoded bytes.
func (k Keeper) MarshalDenomTrace(denomTrace types.DenomTrace) ([]byte, error) {
	return k.cdc.GetProtocMarshal().MarshalBinaryBare(&denomTrace)
}

// MustMarshalDenomTrace attempts to encode an DenomTrace object and returns the
// raw encoded bytes. It panics on error.
func (k Keeper) MustMarshalDenomTrace(denomTrace types.DenomTrace) []byte {
	return k.cdc.GetProtocMarshal().MustMarshalBinaryBare(&denomTrace)
}
