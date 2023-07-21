package keeper

import (
	"fmt"
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
		code := request.WasmCode
		request.WasmCode = nil
		fmt.Printf("create contract start：height=%d, gasbefore=%d, IsCheckTx=%v, IsMempoolSimulate=%v\n", ctx.BlockHeight(), gasBefore, ctx.IsCheckTx(), ctx.IsMempoolSimulate())

		request.WasmCode = code
		creator, err := sdk.WasmAddressFromBech32(request.Creator)
		if err != nil {
			return "", ctx.GasMeter().GasConsumed() - gasBefore, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.Creator)
		}
		admin, err := sdk.WasmAddressFromBech32(request.AdminAddr)
		if err != nil {
			return "", ctx.GasMeter().GasConsumed() - gasBefore, sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, request.AdminAddr)
		}
		addr, _, err := keeper.CreateByContract(ctx, creator, request.WasmCode, request.CodeID, request.InitMsg, admin, request.Label, request.IsCreate2, request.Salt, nil)
		if err != nil {
			fmt.Println(9, err)
			return "", ctx.GasMeter().GasConsumed() - gasBefore, err
		}

		fmt.Printf("创建合约end：height=%d, IsCheckTx=%v, IsMempoolSimulate=%v,contractAddr=%s, gas=%d\n", ctx.BlockHeight(), ctx.IsCheckTx(), ctx.IsMempoolSimulate(), addr.String(), ctx.GasMeter().GasConsumed()-gasBefore)
		return addr.String(), ctx.GasMeter().GasConsumed() - gasBefore, nil
	}
}
