package keeper

import (
	"encoding/hex"
	"errors"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/vmbridge/types"
	"math/big"
)

var (
	big0 = big.NewInt(0)
)

func PrecompileHooks(k *Keeper) vm.CallToWasmByPrecompile {
	return func(ctx vm.OKContext, caller, to common.Address, value *big.Int, input []byte, remainGas uint64) ([]byte, uint64, error) {
		sdkCtx, ok := ctx.(*sdk.Context)
		if !ok {
			return nil, 0, errors.New("VMBridge use context is not type of sdk.Context ")
		}
		if value.Cmp(big0) < 0 {
			return nil, 0, errors.New("VMBridge call value is negative")
		}
		wasmContractAddr, calldata, err := types.DecodePrecompileCallToWasmInput(input)
		if err != nil {
			return nil, 0, err
		}

		if ctx.GetEVMStateDB() == nil {
			return nil, 0, errors.New("VMBridge use context have not evm statedb")
		}
		csdb, ok := ctx.GetEVMStateDB().(*evmtypes.CommitStateDB)
		if !ok {
			return nil, 0, errors.New("VMBridge context's statedb is not *evmtypes.CommitStateDB ")
		}
		csdb.ProtectStateDBEnvironment(*sdkCtx)

		buff, err := hex.DecodeString(calldata)
		if err != nil {
			return nil, 0, err
		}

		subCtx, commit := sdkCtx.CacheContextWithMultiSnapShotRWSet()
		currentGasMeter := subCtx.GasMeter()
		gasMeter := sdk.NewGasMeter(remainGas)
		subCtx.SetGasMeter(gasMeter)
		ret, err := k.CallToWasm(subCtx, sdk.AccAddress(caller.Bytes()), wasmContractAddr, sdk.NewIntFromBigInt(value), string(buff))
		left := gasMeter.Limit() - gasMeter.GasConsumed()
		subCtx.SetGasMeter(currentGasMeter)
		if err != nil {
			return nil, left, err
		}

		csdb.CMChangeCommit(commit)
		result, err := types.EncodePrecompileCallToWasmOutput(string(ret))
		return result, left, err
	}
}
