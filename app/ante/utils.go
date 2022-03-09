package ante

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

func getAccount(ak *auth.AccountKeeper, ctx *sdk.Context, addr sdk.AccAddress, accCache auth.Account) (auth.Account, sdk.Gas) {
	var acc auth.Account
	gasMeter := ctx.GasMeter()
	var gasUsed sdk.Gas
	if accCache != nil {
		if ok, gass := exported.TryAddGetAccountGas(gasMeter, ak, accCache); ok {
			acc = accCache
			gasUsed = gass
		}
	}
	if acc == nil {
		acc = ak.GetAccount(*ctx, addr)
		gass, ok := exported.GetAccountGas(ak, acc)
		if ok {
			gasUsed = gass
		} else {
			gasUsed = 0
		}
	}
	return acc, gasUsed
}
