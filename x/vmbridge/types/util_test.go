package types

import (
	"encoding/hex"
	"encoding/json"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueryMessage(t *testing.T) {
	bank := wasmvmtypes.BankQuery{Balance: &wasmvmtypes.BalanceQuery{Address: "0xbbE4733d85bc2b90682147779DA49caB38C0aA1F", Denom: "okt"}}
	buff, err := json.Marshal(bank)
	require.NoError(t, err)
	t.Log("bank balance", string(buff))
	t.Log("bank balance", hex.EncodeToString(buff))

	wasmsmartRequest := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: "0x5A8D648DEE57b2fc90D98DC17fa887159b69638b", Msg: []byte("{\"balance\":{\"address\":\"0xbbE4733d85bc2b90682147779DA49caB38C0aA1F\"}}")}}
	buff, err = json.Marshal(wasmsmartRequest)
	require.NoError(t, err)
	t.Log("wasm smart", string(buff))
	t.Log("wasm smart", hex.EncodeToString(buff))

	key, err := hex.DecodeString("0006636F6E666967636F6E7374616E7473")
	require.NoError(t, err)
	wasmsmartRequest = wasmvmtypes.WasmQuery{Raw: &wasmvmtypes.RawQuery{ContractAddr: "0x5A8D648DEE57b2fc90D98DC17fa887159b69638b", Key: key}}
	buff, err = json.Marshal(wasmsmartRequest)
	require.NoError(t, err)
	t.Log("wasm raw", string(buff))
	t.Log("wasm raw", hex.EncodeToString(buff))

	wasmsmartRequest = wasmvmtypes.WasmQuery{ContractInfo: &wasmvmtypes.ContractInfoQuery{ContractAddr: "0x5A8D648DEE57b2fc90D98DC17fa887159b69638b"}}
	buff, err = json.Marshal(wasmsmartRequest)
	require.NoError(t, err)
	t.Log("wasm info", string(buff))
	t.Log("wasm info", hex.EncodeToString(buff))
}

func TestGetWasmVMQueryRequest(t *testing.T) {
	testaddr := "0x123"
	testQueryMsg := "{\"balance\":{\"address\":\"0xbbE4733d85bc2b90682147779DA49caB38C0aA1F\"}}"

	testCase := []struct {
		name      string
		data      []byte
		isErr     bool
		expectErr string
		postcheck func(request *wasmvmtypes.QueryRequest)
	}{
		{
			name: "normal smart",
			data: func() []byte {
				wasmsmartRequest := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: testaddr, Msg: []byte(testQueryMsg)}}
				buff, err := json.Marshal(wasmsmartRequest)
				require.NoError(t, err)
				return buff
			}(),
			isErr: false,
			postcheck: func(request *wasmvmtypes.QueryRequest) {
				expectWasm := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: testaddr, Msg: []byte(testQueryMsg)}}
				expect := wasmvmtypes.QueryRequest{Wasm: &expectWasm}
				expectBuff, err := json.Marshal(expect)
				require.NoError(t, err)

				actualBuff, err := json.Marshal(request)
				require.NoError(t, err)
				require.Equal(t, string(expectBuff), string(actualBuff))
			},
		},
		{
			name: "normal raw",
			data: func() []byte {
				wasmsmartRequest := wasmvmtypes.WasmQuery{Raw: &wasmvmtypes.RawQuery{ContractAddr: testaddr, Key: []byte(testQueryMsg)}}
				buff, err := json.Marshal(wasmsmartRequest)
				require.NoError(t, err)
				return buff
			}(),
			isErr: false,
			postcheck: func(request *wasmvmtypes.QueryRequest) {
				expectWasm := wasmvmtypes.WasmQuery{Raw: &wasmvmtypes.RawQuery{ContractAddr: testaddr, Key: []byte(testQueryMsg)}}
				expect := wasmvmtypes.QueryRequest{Wasm: &expectWasm}
				expectBuff, err := json.Marshal(expect)
				require.NoError(t, err)

				actualBuff, err := json.Marshal(request)
				require.NoError(t, err)
				require.Equal(t, string(expectBuff), string(actualBuff))
			},
		},
		{
			name: "normal contract info",
			data: func() []byte {
				wasmsmartRequest := wasmvmtypes.WasmQuery{ContractInfo: &wasmvmtypes.ContractInfoQuery{ContractAddr: testaddr}}
				buff, err := json.Marshal(wasmsmartRequest)
				require.NoError(t, err)
				return buff
			}(),
			isErr: false,
			postcheck: func(request *wasmvmtypes.QueryRequest) {
				expectWasm := wasmvmtypes.WasmQuery{ContractInfo: &wasmvmtypes.ContractInfoQuery{ContractAddr: testaddr}}
				expect := wasmvmtypes.QueryRequest{Wasm: &expectWasm}
				expectBuff, err := json.Marshal(expect)
				require.NoError(t, err)

				actualBuff, err := json.Marshal(request)
				require.NoError(t, err)
				require.Equal(t, string(expectBuff), string(actualBuff))
			},
		},
		{
			name: "error json",
			data: func() []byte {
				wasmsmartRequest := wasmvmtypes.WasmQuery{Smart: &wasmvmtypes.SmartQuery{ContractAddr: testaddr, Msg: []byte(testQueryMsg)}}
				buff, err := json.Marshal(wasmsmartRequest)
				require.NoError(t, err)
				buff[0] += 0x1
				return buff
			}(),
			isErr:     true,
			expectErr: "invalid character '|' looking for beginning of value",
		},
		{
			name: "empty json",
			data: func() []byte {
				wasmsmartRequest := wasmvmtypes.WasmQuery{}
				buff, err := json.Marshal(wasmsmartRequest)
				require.NoError(t, err)
				return buff
			}(),
			isErr:     true,
			expectErr: "query request only need one but got 0",
		},
		{
			name: "other json",
			data: func() []byte {
				wasmsmartRequest := wasmvmtypes.BankQuery{Balance: &wasmvmtypes.BalanceQuery{Address: "0xbbE4733d85bc2b90682147779DA49caB38C0aA1F", Denom: "okt"}}
				buff, err := json.Marshal(wasmsmartRequest)
				require.NoError(t, err)
				return buff
			}(),
			isErr:     true,
			expectErr: "query request only need one but got 0",
		},
		{
			name: "mutli smart and raw query",
			data: func() []byte {
				wasmsmartRequest := wasmvmtypes.WasmQuery{
					Smart: &wasmvmtypes.SmartQuery{ContractAddr: testaddr, Msg: []byte(testQueryMsg)},
					Raw:   &wasmvmtypes.RawQuery{ContractAddr: testaddr, Key: []byte(testQueryMsg)},
				}
				buff, err := json.Marshal(wasmsmartRequest)
				require.NoError(t, err)
				return buff
			}(),
			isErr:     true,
			expectErr: "query request only need one but got 2",
		},
		{
			name: "mutli smart and raw query",
			data: func() []byte {
				wasmsmartRequest := wasmvmtypes.WasmQuery{
					Smart:        &wasmvmtypes.SmartQuery{ContractAddr: testaddr, Msg: []byte(testQueryMsg)},
					ContractInfo: &wasmvmtypes.ContractInfoQuery{ContractAddr: testaddr},
				}
				buff, err := json.Marshal(wasmsmartRequest)
				require.NoError(t, err)
				return buff
			}(),
			isErr:     true,
			expectErr: "query request only need one but got 2",
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(tt *testing.T) {
			request, err := GetWasmVMQueryRequest(tc.data)
			if tc.isErr {
				require.Error(tt, err)
				require.ErrorContains(tt, err, tc.expectErr)
			} else {
				require.NoError(tt, err)
				tc.postcheck(request)
			}
		})
	}
}
