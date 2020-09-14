package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/staking"
)

type OrderKeeper interface {
	DumpStore(ctx sdk.Context)
}

type StakingKeeper interface {
	GetAllValidators(ctx sdk.Context) (validators staking.Validators)
	GetValidatorAllShares(ctx sdk.Context, valAddr sdk.ValAddress) staking.SharesResponses
}
