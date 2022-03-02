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

func getAccountGas(ak *auth.AccountKeeper, acc auth.Account) sdk.Gas {
	if acc == nil || ak == nil {
		panic("nil pinter")
	}
	size := ak.GetEncodedAccountSize(acc)
	gas := kvGasConfig.ReadCostFlat + storetypes.Gas(size)*kvGasConfig.ReadCostPerByte
	return gas
}
