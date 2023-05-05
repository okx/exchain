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
