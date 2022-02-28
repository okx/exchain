package types

import (
	sdkerrors "github.com/okex/exchain/ibc-3rd/cosmos-v443/types/errors"
)

const StoreCodespace = "store"

var (
	ErrInvalidProof = sdkerrors.Register(StoreCodespace, 2, "invalid proof")
)
