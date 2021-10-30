package types

import (
	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/x/staking"
)

type OrderKeeper interface {
	DumpStore(ctx sdk.Context)
}

type StakingKeeper interface {
	GetAllValidators(ctx sdk.Context) (validators staking.Validators)
	GetValidatorAllShares(ctx sdk.Context, valAddr sdk.ValAddress) staking.SharesResponses
}

type CrisisKeeper interface {
	AssertInvariants(ctx sdk.Context)
	Invariants() []sdk.Invariant
}
