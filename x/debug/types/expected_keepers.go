package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type OrderKeeper interface {
	DumpStore(ctx sdk.Context)
}

type StakingKeeper interface {
	SanityCheck(ctx sdk.Context) error
}
