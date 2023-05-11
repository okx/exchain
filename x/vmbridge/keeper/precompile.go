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

		if ctx.GetEVMStateDB() == nil {
			return nil, 0, errors.New("VMBridge use context have not evm statedb")
		}
		csdb, ok := ctx.GetEVMStateDB().(*evmtypes.CommitStateDB)
		if !ok {
			return nil, 0, errors.New("VMBridge context's statedb is not *evmtypes.CommitStateDB ")
		}
		csdb.ProtectStateDBEnvironment(*sdkCtx)
		return methodDispatch(k, csdb, *sdkCtx, caller, to, value, input, remainGas)
	}
}

func methodDispatch(k *Keeper, csdb *evmtypes.CommitStateDB, sdkCtx sdk.Context, caller, to common.Address, value *big.Int, input []byte, remainGas uint64) (result []byte, leftGas uint64, err error) {
	method, err := types.GetMethodByIdFromCallData(input)
	if err != nil {
		return nil, 0, err
	}

	params := k.wasmKeeper.GetParams(sdkCtx)
	if !params.VmbridgeEnable {
		return nil, 0, types.ErrVMBridgeEnable
	}
	// prepare subctx for execute cm msg
	subCtx, commit := sdkCtx.CacheContextWithMultiSnapshotRWSet()
	currentGasMeter := subCtx.GasMeter()
	gasMeter := sdk.NewGasMeter(remainGas)
	subCtx.SetGasMeter(gasMeter)

	switch method.Name {
	case types.PrecompileCallToWasm:
		result, leftGas, err = callToWasm(k, subCtx, caller, to, value, input)
	case types.PrecompileQueryToWasm:
		result, leftGas, err = queryToWasm(k, subCtx, caller, to, value, input)
	default:
		result, leftGas, err = nil, 0, errors.New("methodDispatch failed: unknown method")
	}
	subCtx.SetGasMeter(currentGasMeter)
	if err != nil {
		return result, leftGas, err
	}

	//if the result of executing cm msg if success, then update rwset to parent ctx and add cmchange to journal for reverting snapshot in the future
	csdb.CMChangeCommit(commit)
	return result, leftGas, nil
}

func callToWasm(k *Keeper, sdkCtx sdk.Context, caller, to common.Address, value *big.Int, input []byte) ([]byte, uint64, error) {
	wasmContractAddr, calldata, err := types.DecodePrecompileCallToWasmInput(input)
	if err != nil {
		return nil, 0, err
	}
	buff, err := hex.DecodeString(calldata)
	if err != nil {
		return nil, 0, err
	}

	ret, err := k.CallToWasm(sdkCtx, sdk.AccAddress(caller.Bytes()), wasmContractAddr, sdk.NewIntFromBigInt(value), string(buff))
	gasMeter := sdkCtx.GasMeter()
	left := gasMeter.Limit() - gasMeter.GasConsumed()
	if err != nil {
		return nil, left, err
	}

	result, err := types.EncodePrecompileCallToWasmOutput(string(ret))
	return result, left, err
}

func queryToWasm(k *Keeper, sdkCtx sdk.Context, caller, to common.Address, value *big.Int, input []byte) ([]byte, uint64, error) {
	if value.Sign() != 0 {
		return nil, 0, errors.New("queryToWasm can not be send token")
	}
	calldata, err := types.DecodePrecompileQueryToWasmInput(input)
	if err != nil {
		return nil, 0, err
	}
	buff, err := hex.DecodeString(calldata)
	if err != nil {
		return nil, 0, err
	}

	ret, err := k.QueryToWasm(sdkCtx, caller.String(), buff)
	gasMeter := sdkCtx.GasMeter()
	left := gasMeter.Limit() - gasMeter.GasConsumed()
	if err != nil {
		return nil, left, err
	}

	result, err := types.EncodePrecompileQueryToWasmOutput(string(ret))
	return result, left, err
}
