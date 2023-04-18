package eth

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	wasmtypes "github.com/okex/exchain/x/wasm/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	rpctypes "github.com/okex/exchain/app/rpc/types"
)

var (
	wasmContractABI = `
[
	{
		"constant": true,
		"inputs": [
			{
				"name": "input",
				"type": "string"
			}
		],
		"name": "smart_query",
		"outputs": [
			{
				"name": "",
				"type": "string"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "x",
				"type": "uint256"
			}
		],
		"name": "set",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "get",
		"outputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	}
]
`
	wasmQueryParam = "input"
	wasmInvalidErr = fmt.Errorf("invalid input data")
)

func newWasmAbi() abi.ABI {
	wasmABI, err := abi.JSON(strings.NewReader(wasmContractABI))
	if err != nil {
		panic(fmt.Errorf("wasm abi json decode failed: %s", err.Error()))
	}
	return wasmABI
}

func (api *PublicEthereumAPI) wasmCall(args rpctypes.CallArgs, blockNum rpctypes.BlockNumber) (hexutil.Bytes, error) {
	clientCtx := api.clientCtx
	// pass the given block height to the context if the height is not pending or latest
	if !(blockNum == rpctypes.PendingBlockNumber || blockNum == rpctypes.LatestBlockNumber) {
		clientCtx = api.clientCtx.WithHeight(blockNum.Int64())
	}

	if args.Data == nil {
		return nil, wasmInvalidErr
	}
	data := *args.Data

	methodSigData := data[:4]
	inputsSigData := data[4:]
	method, err := api.wasmABI.MethodById(methodSigData)
	if err != nil {
		return nil, err
	}
	inputsMap := make(map[string]interface{})
	if err := method.Inputs.UnpackIntoMap(inputsMap, inputsSigData); err != nil {
		return nil, err
	}

	inputData, err := hex.DecodeString(inputsMap[wasmQueryParam].(string))
	if err != nil {
		return nil, err
	}

	var stateReq wasmtypes.QuerySmartContractStateRequest
	if err := json.Unmarshal(inputData, &stateReq); err != nil {
		return nil, err
	}

	queryClient := wasmtypes.NewQueryClient(clientCtx)
	res, err := queryClient.SmartContractState(context.Background(), &stateReq)
	if err != nil {
		return nil, err
	}

	out, err := clientCtx.CodecProy.GetProtocMarshal().MarshalJSON(res)
	if err != nil {
		return nil, err
	}
	return out, nil
}
