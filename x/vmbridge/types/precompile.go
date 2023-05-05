package types

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	evm_types "github.com/okex/exchain/x/evm/types"
)

const (
	PrecompileCallToWasm = "callToWasm"
)

var (
	PreCompileABI evm_types.ABI

	//go:embed precompile.json
	preCompileJson []byte
)

func init() {
	PreCompileABI = GetPreCompileABI(preCompileJson)
}

func GetPreCompileABI(data []byte) evm_types.ABI {
	ret, err := abi.JSON(bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return evm_types.ABI{ABI: &ret}
}

func DecodePrecompileCallToWasmInput(input []byte) (wasmAddr, calldata string, err error) {
	if !PreCompileABI.IsMatchFunction(PrecompileCallToWasm, input) {
		return "", "", fmt.Errorf("decode precomplie call to wasm input :  input sginature is not %s", PrecompileCallToWasm)
	}
	unpacked, err := PreCompileABI.DecodeInputParam(PrecompileCallToWasm, input)
	if err != nil {
		return "", "", fmt.Errorf("decode precomplie call to wasm input unpack err :  %s", err)
	}
	if len(unpacked) != 2 {
		return "", "", fmt.Errorf("decode precomplie call to wasm input unpack err :  unpack data len expect 2 but got %v", len(unpacked))
	}
	wasmAddr, ok := unpacked[0].(string)
	if !ok {
		return "", "", fmt.Errorf("decode precomplie call to wasm input unpack err : wasmAddr is not type of string")
	}
	calldata, ok = unpacked[1].(string)
	if !ok {
		return "", "", fmt.Errorf("decode precomplie call to wasm input unpack err : calldata is not type of string")
	}
	return wasmAddr, calldata, nil
}

func EncodePrecompileCallToWasmOutput(response string) ([]byte, error) {
	return PreCompileABI.EncodeOutput(PrecompileCallToWasm, []byte(response))
}

func GetMethodByIdFromCallData(calldata []byte) (*abi.Method, error) {
	return PreCompileABI.GetMethodById(calldata)
}
