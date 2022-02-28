package exported

import (
	storetypes "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var kvGasConfig storetypes.GasConfig

func init() {
	kvGasConfig = storetypes.KVGasConfig()
}

type SizerAccountKeeper interface {
	GetEncodedAccountSize(acc Account) (int, error)
}

func TryAddGetAccountGas(gasMeter sdk.GasMeter, ak SizerAccountKeeper, acc Account) bool {
	if ak == nil || gasMeter == nil {
		return false
	}
	if acc == nil {
		return false
	}
	size, err := ak.GetEncodedAccountSize(acc)
	if err != nil || size == 0 {
		return false
	}
	gas := kvGasConfig.ReadCostFlat + storetypes.Gas(size)*kvGasConfig.ReadCostPerByte
	gasMeter.ConsumeGas(gas, "x/bank/internal/keeper/keeper.BaseSendKeeper")
	return true
}
