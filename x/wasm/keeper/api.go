package keeper

import (
	wasmvm "github.com/CosmWasm/wasmvm"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/x/wasm/types"
	"strconv"
)

const (
	// DefaultGasCostHumanAddress is how moch SDK gas we charge to convert to a human address format
	DefaultGasCostHumanAddress = 5
	// DefaultGasCostCanonicalAddress is how moch SDK gas we charge to convert to a canonical address format
	DefaultGasCostCanonicalAddress = 4

	// DefaultDeserializationCostPerByte The formular should be `len(data) * deserializationCostPerByte`
	DefaultDeserializationCostPerByte = 1

	CallCreateDepth = 20
)

var (
	costHumanize            = DefaultGasCostHumanAddress * DefaultGasMultiplier
	costCanonical           = DefaultGasCostCanonicalAddress * DefaultGasMultiplier
	costJSONDeserialization = wasmvmtypes.UFraction{
		Numerator:   DefaultDeserializationCostPerByte * DefaultGasMultiplier,
		Denominator: 1,
	}
)

func humanAddress(canon []byte) (string, uint64, error) {
	if err := sdk.WasmVerifyAddress(canon); err != nil {
		return "", costHumanize, err
	}
	return sdk.WasmAddress(canon).String(), costHumanize, nil
}

func canonicalAddress(human string) ([]byte, uint64, error) {
	bz, err := sdk.WasmAddressFromBech32(human)
	return bz, costCanonical, err
}

var cosmwasmAPI = wasmvm.GoAPI{
	HumanAddress:     humanAddress,
	CanonicalAddress: canonicalAddress,
}

func contractExternal(ctx sdk.Context, k Keeper) func(request wasmvmtypes.ContractCreateRequest, gasLimit uint64) (string, uint64, error) {
	return func(request wasmvmtypes.ContractCreateRequest, gasLimit uint64) (string, uint64, error) {
		ctx.IncrementCallDepth()
		if ctx.CallDepth() >= CallCreateDepth {
			return "", 0, sdkerrors.Wrap(types.ErrExceedCallDepth, strconv.Itoa(int(ctx.CallDepth())))
		}

		//gasMeter := ctx.GasMeter()
		//ctx.SetGasMeter(sdk.NewGasMeter(k.gasRegister.FromWasmVMGas(gasLimit)))
		gasBefore := ctx.GasMeter().GasConsumed()

		//defer func() {
		//	ctx.DecrementCallDepth()
		//
		//	// reset gas meter
		//	gasCost := ctx.GasMeter().GasConsumed() - gasBefore
		//	ctx.SetGasMeter(gasMeter)
		//	ctx.GasMeter().ConsumeGas(gasCost, "contract sub-create")
		//}()

		creator, err := sdk.WasmAddressFromBech32(request.Creator)
		if err != nil {
			return "", 0, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Creator)
		}
		admin, err := sdk.WasmAddressFromBech32(request.AdminAddr)
		if err != nil {
			return "", 0, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.AdminAddr)
		}
		addr, _, err := k.CreateByContract(ctx, creator, request.WasmCode, request.CodeID, request.InitMsg, admin, request.Label, request.IsCreate2, request.Salt, nil)
		if err != nil {
			return "", k.gasRegister.ToWasmVMGas(ctx.GasMeter().GasConsumed()) - gasBefore, err
		}

		return addr.String(), k.gasRegister.ToWasmVMGas(ctx.GasMeter().GasConsumed() - gasBefore), nil
	}
}
