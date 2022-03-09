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
	GetEncodedAccountSize(acc Account) int
}

func TryAddGetAccountGas(gasMeter sdk.GasMeter, ak SizerAccountKeeper, acc Account) (bool, sdk.Gas) {
	if ak == nil || gasMeter == nil || acc == nil {
		return false, 0
	}
	size := ak.GetEncodedAccountSize(acc)
	if size == 0 {
		return false, 0
	}
	gas := kvGasConfig.ReadCostFlat + storetypes.Gas(size)*kvGasConfig.ReadCostPerByte
	gasMeter.ConsumeGas(gas, "x/bank/internal/keeper/keeper.BaseSendKeeper")
	return true, gas
}

func GetAccountGas(ak SizerAccountKeeper, acc Account) (sdk.Gas, bool) {
	if acc == nil || ak == nil {
		return 0, false
	}
	size := ak.GetEncodedAccountSize(acc)
	if size == 0 {
		return 0, false
	}
	gas := kvGasConfig.ReadCostFlat + storetypes.Gas(size)*kvGasConfig.ReadCostPerByte
	return gas, true
}
