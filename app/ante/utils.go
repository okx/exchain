package ante

import (
	storetypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

var kvGasConfig storetypes.GasConfig

func init() {
	kvGasConfig = storetypes.KVGasConfig()
}

func getAccount(ak *auth.AccountKeeper, ctx *sdk.Context, addr sdk.AccAddress, accCache auth.Account) (auth.Account, sdk.Gas) {
	var acc auth.Account
	gasMeter := ctx.GasMeter()
	gasBefore := gasMeter.GasConsumed()
	var gasUsed sdk.Gas
	if accCache != nil {
		if exported.TryAddGetAccountGas(gasMeter, ak, accCache) {
			acc = accCache
			gasUsed = gasMeter.GasConsumed() - gasBefore
		}
	}
	if acc == nil {
		acc = ak.GetAccount(*ctx, addr)
		gasUsed = gasMeter.GasConsumed() - gasBefore
	}
	return acc, gasUsed
}
