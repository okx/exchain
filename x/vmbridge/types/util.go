package types

import (
	"encoding/json"
	"fmt"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
)

func GetWasmVMQueryRequest(requestData []byte) (*wasmvmtypes.QueryRequest, error) {
	var request wasmvmtypes.WasmQuery
	if err := json.Unmarshal(requestData, &request); err != nil {
		return nil, err
	}

	requestNum := 0
	if request.Smart != nil {
		requestNum++
	}
	if request.Raw != nil {
		requestNum++
	}
	if request.ContractInfo != nil {
		requestNum++
	}
	if requestNum != 1 {
		return nil, fmt.Errorf("query request only need one but got %d", requestNum)
	}
	result := wasmvmtypes.QueryRequest{Wasm: &request}
	return &result, nil
}
