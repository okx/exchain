package keeper

import (
	wasmvm "github.com/CosmWasm/wasmvm"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
)

const (
	// DefaultGasCostHumanAddress is how moch SDK gas we charge to convert to a human address format
	DefaultGasCostHumanAddress = 5
	// DefaultGasCostCanonicalAddress is how moch SDK gas we charge to convert to a canonical address format
	DefaultGasCostCanonicalAddress = 4

	// DefaultDeserializationCostPerByte The formular should be `len(data) * deserializationCostPerByte`
	DefaultDeserializationCostPerByte = 1
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

func contractExternal(ctx sdk.Context, keeper Keeper) func(request wasmvmtypes.ContractCreateRequest, gasLimit uint64) (string, uint64, error) {
	return func(request wasmvmtypes.ContractCreateRequest, gasLimit uint64) (string, uint64, error) {
		gasBefore := ctx.GasMeter().GasConsumed()
		creator, err := sdk.WasmAddressFromBech32(request.Creator)
		if err != nil {
			return "", ctx.GasMeter().GasConsumed() - gasBefore, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Creator)
		}
		admin, err := sdk.WasmAddressFromBech32(request.AdminAddr)
		if err != nil {
			return "", ctx.GasMeter().GasConsumed() - gasBefore, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.AdminAddr)
		}
		addr, _, err := keeper.CreateByContract(ctx, creator, request.WasmCode, request.InitMsg, admin, request.Label, nil)
		if err != nil {
			return "", ctx.GasMeter().GasConsumed() - gasBefore, err
		}

		return addr.String(), ctx.GasMeter().GasConsumed() - gasBefore, nil
	}
}
