package types

import (
	"errors"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetPreCompileABI(t *testing.T) {
	testCase := []struct {
		name      string
		data      []byte
		isErr     bool
		expectErr string
		isequal   bool
	}{
		{
			name:    "normal abi json",
			data:    preCompileJson,
			isErr:   false,
			isequal: true,
		},
		{
			name:    "normal abi json have more func",
			data:    []byte("[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"wasmAddr\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"callToWasm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"response\",\"type\":\"string\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"queryToWasm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"response\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"queryToWasm1\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"response\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]\n"),
			isErr:   false,
			isequal: false,
		},
		{
			name:    "normal abi json have more event",
			data:    []byte("[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"wasmAddr\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"callToWasm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"response\",\"type\":\"string\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"wasmAddr\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"__OKCSendToWasm\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"queryToWasm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"response\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]\n"),
			isErr:   false,
			isequal: false,
		},
		{
			name:    "normal abi json have less call to evm func",
			data:    []byte("[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"data\",\"type\":\"string\"}],\"name\":\"queryToWasm\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"response\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]\n"),
			isErr:   false,
			isequal: false,
		},
		{
			name:      "error abi json",
			data:      []byte("[\n  {\n    \"anonymous\": false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"wasmAddr\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"value\",\n        \"type\": \"uint256\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"calldata\",\n        \"type\": \"string\"\n      }\n    ],\n    \"name\": \"__OKCCallToWasm\",\n    \"type\": \"event\"\n  },\n  {\n    \"inputs\": [\n      {\n        \"internalType\": \"string\",\n        \"name\": \"callerWasmAddr\",\n        \"type\": \"string\"\n      },\n      {\n        \"internalType\": \"string\",\n        \"name\": \"data\",\n        \"type\": \"string\"\n      }\n    ],\n    \"name\": \"callByWasm\",\n    \"outputs\": [\n      {\n        \"internalType\": \"string\",\n        \"name\": \"response\",\n        \"type\": \"string\"\n      }\n    ],\n    \"stateMutability\": \"payable\",\n    \"type\": \"function\"\n  },{\n    : false,\n    \"inputs\": [\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"wasmAddr\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"string\",\n        \"name\": \"recipient\",\n        \"type\": \"string\"\n      },\n      {\n        \"indexed\": false,\n        \"internalType\": \"uint256\",\n        \"name\": \"amount\",\n        \"type\": \"uint256\"\n      }\n    ],\n    \"name\": \"__OKCSendToWasm\",\n    \"type\": \"event\"\n  }\n]"),
			isErr:     true,
			expectErr: "invalid character ':' looking for beginning of object key string",
			isequal:   false,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(tt *testing.T) {
			if tc.isErr {
				defer func() {
					r := recover()
					require.NotNil(tt, r)
					err := r.(error)
					require.ErrorContains(tt, err, tc.expectErr)
				}()
			}
			abi := GetPreCompileABI(tc.data)
			if tc.isequal {
				require.Equal(tt, *PreCompileABI.ABI, *abi.ABI)
			} else {
				require.NotEqual(tt, *PreCompileABI.ABI, *abi.ABI)
			}
		})
	}
}

func TestDecodePrecompileCallToWasmInput(t *testing.T) {
	testaddr := "0x123"
	testCalldata := "test call data"
	testCase := []struct {
		name      string
		data      []byte
		isErr     bool
		expectErr string
		postcheck func(wasmAddr, calldata string)
	}{
		{
			name: "normal",
			data: func() []byte {
				buff, err := encodeCallToWasmInput(testaddr, testCalldata)
				require.NoError(t, err)
				return buff
			}(),
			isErr: false,
			postcheck: func(wasmAddr, calldata string) {
				require.Equal(t, testaddr, wasmAddr)
				require.Equal(t, testCalldata, calldata)
			},
		},
		{
			name: "input is nil",
			data: func() []byte {
				buff, err := PreCompileABI.Methods[PrecompileCallToWasm].Inputs.Pack(testaddr)
				require.Error(t, err)
				require.ErrorContains(t, err, "argument count mismatch: got 1 for 2")

				return append(PreCompileABI.Methods[PrecompileCallToWasm].ID, buff...)
			}(),
			isErr: true,
			postcheck: func(wasmAddr, calldata string) {
				require.Equal(t, testaddr, wasmAddr)
				require.Equal(t, testCalldata, calldata)
			},
			expectErr: "decode precomplie call to wasm input unpack err :  method callToWasm data is nil",
		},
		{
			name: "input add one byte",
			data: func() []byte {
				buff, err := encodeCallToWasmInput(testaddr, testCalldata)
				require.NoError(t, err)
				buff = append(buff, 0x1)
				return append(PreCompileABI.Methods[PrecompileCallToWasm].ID, buff...)
			}(),
			isErr: true,
			postcheck: func(wasmAddr, calldata string) {
				require.Equal(t, testaddr, wasmAddr)
				require.Equal(t, testCalldata, calldata)
			},
			expectErr: "decode precomplie call to wasm input unpack err :  abi: cannot marshal in to go slice: offset 41246180337441339990308983541005712728593953418436849081514007625224398307360 would go over slice boundary (len=197)",
		},
		{
			name: "input less one byte",
			data: func() []byte {
				buff, err := encodeCallToWasmInput(testaddr, testCalldata)
				require.NoError(t, err)
				length := len(buff)
				buff = buff[:length-1]
				return append(PreCompileABI.Methods[PrecompileCallToWasm].ID, buff...)
			}(),
			isErr: true,
			postcheck: func(wasmAddr, calldata string) {
				require.Equal(t, testaddr, wasmAddr)
				require.Equal(t, testCalldata, calldata)
			},
			expectErr: "decode precomplie call to wasm input unpack err :  abi: cannot marshal in to go slice: offset 41246180337441339990308983541005712728593953418436849081514007625224398307360 would go over slice boundary (len=195)",
		},
		{
			name: "input update one byte",
			data: func() []byte {
				buff, err := encodeCallToWasmInput(testaddr, testCalldata)
				require.NoError(t, err)
				length := len(buff)
				buff[length-1] += 0x1
				return append(PreCompileABI.Methods[PrecompileCallToWasm].ID, buff...)
			}(),
			isErr: true,
			postcheck: func(wasmAddr, calldata string) {
				require.Equal(t, testaddr, wasmAddr)
				require.Equal(t, testCalldata, calldata)
			},
			expectErr: "decode precomplie call to wasm input unpack err :  abi: cannot marshal in to go slice: offset 41246180337441339990308983541005712728593953418436849081514007625224398307360 would go over slice boundary (len=196)",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(tt *testing.T) {
			wasmAddr, calldata, err := DecodePrecompileCallToWasmInput(tc.data)
			if tc.isErr {
				require.Error(tt, err)
				require.ErrorContains(tt, err, tc.expectErr)
			} else {
				require.NoError(tt, err)
				tc.postcheck(wasmAddr, calldata)
			}
		})
	}
}
func TestEncodePrecompileCallToWasmOutput(t *testing.T) {
	testResult := "test result"

	input, err := EncodePrecompileCallToWasmOutput(testResult)
	require.NoError(t, err)
	result, err := decodeCallToWasmOutput(input)
	require.NoError(t, err)
	require.Equal(t, testResult, result)
}

func TestDecodePrecompileQueryToWasmInput(t *testing.T) {
	testCalldata := "test call data"
	testCase := []struct {
		name      string
		data      []byte
		isErr     bool
		expectErr string
		postcheck func(calldata string)
	}{
		{
			name: "normal",
			data: func() []byte {
				buff, err := encodeQueryToWasmInput(testCalldata)
				require.NoError(t, err)
				return buff
			}(),
			isErr: false,
			postcheck: func(calldata string) {
				require.Equal(t, testCalldata, calldata)
			},
		},
		{
			name: "input is nil",
			data: func() []byte {
				buff, err := PreCompileABI.Methods[PrecompileQueryToWasm].Inputs.Pack()
				require.Error(t, err)
				require.ErrorContains(t, err, "argument count mismatch: got 0 for 1")

				return append(PreCompileABI.Methods[PrecompileQueryToWasm].ID, buff...)
			}(),
			isErr: true,
			postcheck: func(calldata string) {
				require.Equal(t, testCalldata, calldata)
			},
			expectErr: "decode precomplie query to wasm input unpack err :  method queryToWasm data is nil",
		},
		{
			name: "input add one byte",
			data: func() []byte {
				buff, err := encodeQueryToWasmInput(testCalldata)
				require.NoError(t, err)
				buff = append(buff, 0x1)
				return append(PreCompileABI.Methods[PrecompileQueryToWasm].ID, buff...)
			}(),
			isErr: true,
			postcheck: func(calldata string) {
				require.Equal(t, testCalldata, calldata)
			},
			expectErr: "decode precomplie query to wasm input unpack err :  abi: cannot marshal in to go slice: offset 86015489902299205649965666741627051758092273748291497414970767656871234895904 would go over slice boundary (len=101)",
		},
		{
			name: "input less one byte",
			data: func() []byte {
				buff, err := encodeQueryToWasmInput(testCalldata)
				require.NoError(t, err)
				length := len(buff)
				buff = buff[:length-1]
				return append(PreCompileABI.Methods[PrecompileQueryToWasm].ID, buff...)
			}(),
			isErr: true,
			postcheck: func(calldata string) {
				require.Equal(t, testCalldata, calldata)
			},
			expectErr: "decode precomplie query to wasm input unpack err :  abi: cannot marshal in to go slice: offset 86015489902299205649965666741627051758092273748291497414970767656871234895904 would go over slice boundary (len=99)",
		},
		{
			name: "input update one byte",
			data: func() []byte {
				buff, err := encodeQueryToWasmInput(testCalldata)
				require.NoError(t, err)
				length := len(buff)
				buff[length-1] += 0x1
				return append(PreCompileABI.Methods[PrecompileQueryToWasm].ID, buff...)
			}(),
			isErr: true,
			postcheck: func(calldata string) {
				require.Equal(t, testCalldata, calldata)
			},
			expectErr: "decode precomplie query to wasm input unpack err :  abi: cannot marshal in to go slice: offset 86015489902299205649965666741627051758092273748291497414970767656871234895904 would go over slice boundary (len=100)",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(tt *testing.T) {
			calldata, err := DecodePrecompileQueryToWasmInput(tc.data)
			if tc.isErr {
				require.Error(tt, err)
				require.ErrorContains(tt, err, tc.expectErr)
			} else {
				require.NoError(tt, err)
				tc.postcheck(calldata)
			}
		})
	}
}

func TestEncodePrecompileQueryToWasmOutput(t *testing.T) {
	testResult := "test result"

	input, err := EncodePrecompileQueryToWasmOutput(testResult)
	require.NoError(t, err)
	result, err := decodeQueryToWasmOutput(input)
	require.NoError(t, err)
	require.Equal(t, testResult, result)
}

func TestGetMethodByIdFromCallData(t *testing.T) {
	testaddr := "0x123"
	testCalldata := "test call data"
	testCase := []struct {
		name      string
		data      []byte
		isErr     bool
		expectErr string
		postcheck func(method *abi.Method)
	}{
		{
			name: "normal call to wasm",
			data: func() []byte {
				buff, err := encodeCallToWasmInput(testaddr, testCalldata)
				require.NoError(t, err)
				return buff
			}(),
			isErr: false,
			postcheck: func(method *abi.Method) {
				require.Equal(t, PreCompileABI.Methods[PrecompileCallToWasm], *method)
			},
		},
		{
			name: "normal query to wasm",
			data: func() []byte {
				buff, err := encodeQueryToWasmInput(testCalldata)
				require.NoError(t, err)
				return buff
			}(),
			isErr: false,
			postcheck: func(method *abi.Method) {
				require.Equal(t, PreCompileABI.Methods[PrecompileQueryToWasm], *method)
			},
		},
		{
			name: "error method not found",
			data: func() []byte {
				buff, err := encodeQueryToWasmInput(testCalldata)
				require.NoError(t, err)
				buff[0] += 0x1
				return buff
			}(),
			isErr:     true,
			expectErr: "no method with id: 0xbf2b0ac2",
		},
		{
			name: "error input is error",
			data: func() []byte {
				buff, err := encodeQueryToWasmInput(testCalldata)
				require.NoError(t, err)
				buff = append(buff, 0x1)
				return buff
			}(),
			isErr:     true,
			expectErr: "invalid call data; length should be a multiple of 32 bytes (was 97)",
		},
		{
			name: "error input is less than 4",
			data: func() []byte {
				buff, err := encodeQueryToWasmInput(testCalldata)
				require.NoError(t, err)
				buff = buff[:3]
				return buff
			}(),
			isErr:     true,
			expectErr: "the calldata length must more than 4",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(tt *testing.T) {
			method, err := GetMethodByIdFromCallData(tc.data)
			if tc.isErr {
				require.Error(tt, err)
				require.ErrorContains(tt, err, tc.expectErr)
			} else {
				require.NoError(tt, err)
				tc.postcheck(method)
			}
		})
	}

}

func encodeCallToWasmInput(wasmAddr, calldata string) ([]byte, error) {
	return PreCompileABI.Pack(PrecompileCallToWasm, wasmAddr, calldata)
}

func decodeCallToWasmOutput(input []byte) (string, error) {
	pack, err := PreCompileABI.Methods[PrecompileCallToWasm].Outputs.Unpack(input)
	if err != nil {
		return "", err
	}
	if len(pack) != 1 {
		return "", errors.New("decodeCallToWasmOutput failed: got multi result")
	}

	return pack[0].(string), nil
}

func encodeQueryToWasmInput(calldata string) ([]byte, error) {
	return PreCompileABI.Pack(PrecompileQueryToWasm, calldata)
}

func decodeQueryToWasmOutput(input []byte) (string, error) {
	pack, err := PreCompileABI.Methods[PrecompileQueryToWasm].Outputs.Unpack(input)
	if err != nil {
		return "", err
	}
	if len(pack) != 1 {
		return "", errors.New("decodeCallToWasmOutput failed: got multi result")
	}

	return pack[0].(string), nil
}
