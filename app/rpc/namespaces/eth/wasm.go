package eth

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/x/wasm/ioutils"

	clientcontext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"

	wasmtypes "github.com/okex/exchain/x/wasm/types"

	"github.com/ethereum/go-ethereum/common/hexutil"
	rpctypes "github.com/okex/exchain/app/rpc/types"
)

const (
	genMsgStoreCode          = "genMsgStoreCode"
	wasmHelperABIStr         = `[{"inputs":[{"internalType":"address","name":"_contract","type":"address"},{"internalType":"string","name":"_msg","type":"string"}],"name":"genMsgExecuteContract","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_contract","type":"address"},{"internalType":"string","name":"_msg","type":"string"},{"internalType":"string","name":"amount","type":"string"}],"name":"genMsgExecuteContractWithOKT","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_admin","type":"address"},{"internalType":"uint256","name":"_codeID","type":"uint256"},{"internalType":"string","name":"_label","type":"string"},{"internalType":"string","name":"_msg","type":"string"}],"name":"genMsgInstantiateContract","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_admin","type":"address"},{"internalType":"uint256","name":"_codeID","type":"uint256"},{"internalType":"string","name":"_label","type":"string"},{"internalType":"string","name":"_msg","type":"string"},{"internalType":"string","name":"amount","type":"string"}],"name":"genMsgInstantiateContractWithOKT","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_contract","type":"address"},{"internalType":"uint256","name":"_codeID","type":"uint256"}],"name":"genMsgMigrateContract","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_contract","type":"address"},{"internalType":"uint256","name":"_codeID","type":"uint256"},{"internalType":"string","name":"_msg","type":"string"}],"name":"genMsgMigrateContractWithMSG","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"_wasmBytecode","type":"bytes"},{"internalType":"string","name":"_permission","type":"string"},{"internalType":"address","name":"_addr","type":"address"}],"name":"genMsgStoreCode","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"_newAdmin","type":"address"},{"internalType":"address","name":"_contract","type":"address"}],"name":"genMsgUpdateAdmin","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"input","type":"string"}],"name":"invoke","outputs":[],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"string","name":"_str","type":"string"}],"name":"stringToHexString","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"pure","type":"function"}]`
	FlagE2cWasmMsgHelperAddr = "e2c-wasm-msg-helper-addr"
)

var (
	wasmQueryParam = "input"
	wasmInvalidErr = fmt.Errorf("invalid input data")
	wasmHelperABI  *evmtypes.ABI
)

func init() {
	abi, err := evmtypes.NewABI(wasmHelperABIStr)
	if err != nil {
		panic(fmt.Errorf("wasm abi json decode failed: %s", err.Error()))
	}
	wasmHelperABI = abi
}

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

func (api *PublicEthereumAPI) isEvm2CmTx(to *common.Address) bool {
	if to == nil {
		return false
	}
	return bytes.Equal(api.systemContract, to.Bytes())
}

func (api *PublicEthereumAPI) isLargeWasmMsgStoreCode(args rpctypes.CallArgs) (code, newparam []byte, is bool) {
	if args.To == nil || args.Data == nil || len(*args.Data) <= int(api.e2cWasmCodeLimit) {
		return nil, nil, false
	}
	// set the e2cWasmMsgHelperAddr should only this contract address judge the large msg store code
	if api.e2cWasmMsgHelperAddr != "" && !bytes.Equal(common.HexToAddress(api.e2cWasmMsgHelperAddr).Bytes(), args.To.Bytes()) {
		return nil, nil, false
	}
	if !wasmHelperABI.IsMatchFunction(genMsgStoreCode, *args.Data) {
		return nil, nil, false
	}
	data, res, err := ParseMsgStoreCodeParam(*args.Data)
	if err != nil {
		return nil, nil, false
	}
	newparam, err = genNullCodeMsgStoreCodeParam(res)
	if err != nil {
		return nil, nil, false
	}
	return data, newparam, true
}

func ParseMsgStoreCodeParam(input []byte) ([]byte, []interface{}, error) {
	res, err := wasmHelperABI.DecodeInputParam(genMsgStoreCode, input)
	if err != nil {
		return nil, nil, err
	}
	if len(res) < 1 {
		return nil, nil, wasmInvalidErr
	}
	v, ok := res[0].([]byte)
	if !ok {
		return nil, nil, wasmInvalidErr
	}
	return v, res, nil
}

func genNullCodeMsgStoreCodeParam(input []interface{}) ([]byte, error) {
	if len(input) == 0 {
		return nil, wasmInvalidErr
	}
	_, ok := input[0].([]byte)
	if !ok {
		return nil, wasmInvalidErr
	}
	input[0] = []byte{}
	return wasmHelperABI.Pack(genMsgStoreCode, input...)
}

type MsgWrapper struct {
	Name string          `json:"type"`
	Data json.RawMessage `json:"value"`
}

func replaceToRealWasmCode(ret, code []byte) ([]byte, error) {
	re, err := wasmHelperABI.Unpack(genMsgStoreCode, ret)
	if err != nil || len(re) != 1 {
		return nil, wasmInvalidErr
	}
	hexdata, ok := re[0].(string)
	if !ok {
		return nil, wasmInvalidErr
	}

	// decode
	msgWrap, msc, err := hexDecodeToMsgStoreCode(hexdata)
	if err != nil {
		return nil, err
	}

	// replace and encode
	rstr, err := msgStoreCodeToHexDecode(msgWrap, msc, code)
	if err != nil {
		return nil, err
	}

	rret, err := wasmHelperABI.EncodeOutput(genMsgStoreCode, []byte(rstr))
	if err != nil {
		return nil, err
	}
	return rret, nil
}

func hexDecodeToMsgStoreCode(input string) (*MsgWrapper, *wasmtypes.MsgStoreCode, error) {
	value, err := hex.DecodeString(input)
	if err != nil {
		return nil, nil, err
	}
	var msgWrap MsgWrapper
	if err := json.Unmarshal(value, &msgWrap); err != nil {
		return nil, nil, err
	}
	var msc wasmtypes.MsgStoreCode
	if err := json.Unmarshal(msgWrap.Data, &msc); err != nil {
		return nil, nil, err
	}
	return &msgWrap, &msc, nil
}

func msgStoreCodeToHexDecode(msgWrap *MsgWrapper, msc *wasmtypes.MsgStoreCode, code []byte) (string, error) {
	msc.WASMByteCode = code
	v, err := json.Marshal(msc)
	if err != nil {
		return "", err
	}
	msgWrap.Data = v
	rData, err := json.Marshal(msgWrap)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(rData), nil
}

// get from the cli
func judgeWasmCode(input []byte) ([]byte, error) {
	// gzip the wasm file
	if ioutils.IsWasm(input) {
		wasm, err := ioutils.GzipIt(input)
		if err != nil {
			return nil, err
		}
		return wasm, nil
	} else if !ioutils.IsGzip(input) {
		return nil, fmt.Errorf("invalid input file. Use wasm binary or gzip")
	}
	return input, nil
}
