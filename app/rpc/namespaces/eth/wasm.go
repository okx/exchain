package eth

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"

	clientcontext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"

	wasmtypes "github.com/okex/exchain/x/wasm/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	rpctypes "github.com/okex/exchain/app/rpc/types"
)

var (
	wasmQueryParam = "input"
	wasmInvalidErr = fmt.Errorf("invalid input data")
)

func getSystemContractAddr(clientCtx clientcontext.CLIContext) []byte {
	route := fmt.Sprintf("custom/%s/%s", evmtypes.ModuleName, evmtypes.QuerySysContractAddress)
	addr, _, err := clientCtx.QueryWithData(route, nil)
	if err != nil {
		return nil
	}
	return addr
}

type SmartContractStateRequest struct {
	// address is the address of the contract
	Address string `json:"address"`
	// QueryData contains the query data passed to the contract
	QueryData string `json:"query_data"`
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
	method, err := evm.SysABI().MethodById(methodSigData)
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

	var stateReq SmartContractStateRequest
	if err := json.Unmarshal(inputData, &stateReq); err != nil {
		return nil, err
	}

	queryData, err := hex.DecodeString(stateReq.QueryData)
	if err != nil {
		return nil, wasmInvalidErr
	}

	queryClient := wasmtypes.NewQueryClient(clientCtx)
	res, err := queryClient.SmartContractState(context.Background(), &wasmtypes.QuerySmartContractStateRequest{
		Address:   stateReq.Address,
		QueryData: queryData,
	})
	if err != nil {
		return nil, err
	}

	out, err := clientCtx.CodecProy.GetProtocMarshal().MarshalJSON(res)
	if err != nil {
		return nil, err
	}
	result, err := evm.EncodeQueryOutput(out)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (api *PublicEthereumAPI) isWasmCall(args rpctypes.CallArgs) bool {
	if args.To == nil || !bytes.Equal(args.To.Bytes(), api.systemContract) {
		return false
	}
	return args.Data != nil && evm.IsMatchSystemContractQuery(*args.Data)
}
