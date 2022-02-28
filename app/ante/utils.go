package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

func getAccount(ak *auth.AccountKeeper, ctx *sdk.Context, addr sdk.AccAddress, accCache auth.Account) auth.Account {
	var acc auth.Account
	if accCache != nil {
		if exported.TryAddGetAccountGas(ctx.GasMeter(), ak, accCache) {
			acc = accCache
		}
	}
	if acc == nil {
		acc = ak.GetAccount(*ctx, addr)
	}
	return acc
}
