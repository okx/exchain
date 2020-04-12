package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OrderKeeper interface {
	DumpStore(ctx sdk.Context)
}
