package eth

import (
	"encoding/json"
	"math/big"
	"strings"
	"time"

	"github.com/okex/okexchain/x/evm/types"

	"github.com/ethereum/go-ethereum/core/vm"

	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethermint "github.com/okex/okexchain/app/types"
	"github.com/spf13/viper"
)

const (
	DefaultEVMErrorCode          = -32000
	VMExecuteException           = -32015
	VMExecuteExceptionInEstimate = 3

	RPCEthCall           = "eth_call"
	RPCEthEstimateGas    = "eth_estimateGas"
	RPCEthGetBlockByHash = "eth_getBlockByHash"

	RPCUnknowErr = "unknow"
	RPCNullData  = "null"
)

//gasPrice: to get "minimum-gas-prices" config or to get ethermint.DefaultGasPrice
func ParseGasPrice() *hexutil.Big {
	gasPrices, err := sdk.ParseDecCoins(viper.GetString(server.FlagMinGasPrices))
	if err == nil && gasPrices != nil && len(gasPrices) > 0 {
		return (*hexutil.Big)(gasPrices[0].Amount.BigInt())
	}

	//return the default gas price : DefaultGasPrice
	return (*hexutil.Big)(sdk.NewDecFromBigIntWithPrec(big.NewInt(ethermint.DefaultGasPrice), sdk.Precision/2).BigInt())
}

type performanceSimulate struct {
	callTimeStamp int64
}

func (p *performanceSimulate) BeginSimulate() *performanceSimulate {
	p.callTimeStamp = time.Now().UnixNano()
	return p
}

func (p performanceSimulate) EndSimulate() uint64 {
	//return in ms
	return uint64((time.Now().UnixNano() - p.callTimeStamp) / int64(1e6))
}

func NewPerformanceSimulate() *performanceSimulate {
	return &performanceSimulate{
		callTimeStamp: 0,
	}
}

type cosmosError struct {
	Code      int    `json:"code"`
	Log       string `json:"log"`
	Codespace string `json:"codespace"`
}

func (c cosmosError) Error() string {
	return c.Log
}

func newCosmosError(code int, log, codeSpace string) cosmosError {
	return cosmosError{
		Code:      code,
		Log:       log,
		Codespace: codeSpace,
	}
}

func newWrappedCosmosError(code int, log, codeSpace string) cosmosError {
	e := newCosmosError(code, log, codeSpace)
	b, _ := json.Marshal(e)
	e.Log = string(b)
	return e
}

type wrappedEthError struct {
	Wrap ethDataError `json:"0x00000000000000000000000000000000"`
}

type ethDataError struct {
	Error          string `json:"error"`
	ProgramCounter int    `json:"program_counter"`
	Reason         string `json:"reason"`
	Ret            string `json:"return"`
}

type DataError struct {
	code int         `json:"code"`
	Msg  string      `json:"msg"`
	data interface{} `json:"data,omitempty"`
}

func (d DataError) Error() string {
	return d.Msg
}

func (d DataError) ErrorData() interface{} {
	return d.data
}

func (d DataError) ErrorCode() int {
	return d.code
}

func newDataError(revert string, data string) *wrappedEthError {
	return &wrappedEthError{
		Wrap: ethDataError{
			Error:          "revert",
			ProgramCounter: 0,
			Reason:         revert,
			Ret:            data,
		}}
}

func TransformDataError(err error, method string) error {
	msg := err.Error()
	var realErr cosmosError
	if len(msg) > 0 {
		e := json.Unmarshal([]byte(msg), &realErr)
		if e != nil {
			return DataError{
				code: DefaultEVMErrorCode,
				Msg:  err.Error(),
				data: RPCNullData,
			}
		}
		if method == RPCEthGetBlockByHash {
			return DataError{
				code: DefaultEVMErrorCode,
				Msg:  realErr.Error(),
				data: RPCNullData,
			}
		}
		m, retErr := preProcessError(realErr, err.Error())
		if retErr != nil {
			return realErr
		}
		//if there have multi error type of EVM, this need a reactor mode to process error
		revert, f := m[vm.ErrExecutionReverted.Error()]
		if !f {
			revert = RPCUnknowErr
		}
		data, f := m[types.ErrorHexData]
		if !f {
			data = RPCNullData
		}
		switch method {
		case RPCEthEstimateGas:
			return DataError{
				code: VMExecuteExceptionInEstimate,
				Msg:  revert,
				data: data,
			}
		case RPCEthCall:
			return DataError{
				code: VMExecuteException,
				Msg:  revert,
				data: newDataError(revert, data),
			}
		default:
			return DataError{
				code: DefaultEVMErrorCode,
				Msg:  revert,
				data: newDataError(revert, data),
			}
		}

	}
	return DataError{
		code: DefaultEVMErrorCode,
		Msg:  err.Error(),
		data: RPCNullData,
	}
}

//Preprocess error string, the string of realErr.Log is most like:
//`["execution reverted","message","HexData","0x00000000000"];some failed information`
//we need marshalled json slice from realErr.Log and using segment tag `[` and `]` to cut it
func preProcessError(realErr cosmosError, origErrorMsg string) (map[string]string, error) {
	var logs []string
	lastSeg := strings.LastIndexAny(realErr.Log, "]")
	if lastSeg < 0 {
		return nil, DataError{
			code: DefaultEVMErrorCode,
			Msg:  origErrorMsg,
			data: RPCNullData,
		}
	}
	marshaler := realErr.Log[0 : lastSeg+1]
	e := json.Unmarshal([]byte(marshaler), &logs)
	if e != nil {
		return nil, DataError{
			code: DefaultEVMErrorCode,
			Msg:  origErrorMsg,
			data: RPCNullData,
		}
	}
	m := genericStringMap(logs)
	if m == nil {
		return nil, DataError{
			code: DefaultEVMErrorCode,
			Msg:  origErrorMsg,
			data: RPCNullData,
		}
	}
	return m, nil
}

func genericStringMap(s []string) map[string]string {
	var ret = make(map[string]string)
	if len(s)%2 != 0 {
		return nil
	}
	for i := 0; i < len(s); i += 2 {
		ret[s[i]] = s[i+1]
	}
	return ret
}
