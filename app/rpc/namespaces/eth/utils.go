package eth

import (
	"encoding/json"
	"math/big"
	"strings"

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

	RpcEthCall           = "eth_call"
	RpcEthEstimateGas    = "eth_estimateGas"
	RpcEthGetBlockByHash = "eth_getBlockByHash"

	RpcUnknowErr = "unknow"
	RpcNullData  = "null"
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

func NewWrappedCosmosError(code int, log, codeSpace string) cosmosError {
	e := newCosmosError(code, log, codeSpace)
	b, _ := json.Marshal(e)
	e.Log = string(b)
	return e
}

type wrappedEthError struct {
	Wrap ethDataError `json:"0x00000000000000000000000000000000"`
}

type ethDataError struct {
	Error           string `json:"error"`
	Program_counter int    `json:"program_counter"`
	Reason          string `json:"reason"`
	Ret             string `json:"return"`
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
			Error:           "revert",
			Program_counter: 0,
			Reason:          revert,
			Ret:             data,
		}}
}

func TransformDataError(err error, method string) DataError {
	msg := err.Error()
	var logs []string
	var realErr cosmosError
	if len(msg) > 0 {
		e := json.Unmarshal([]byte(msg), &realErr)
		if e != nil {
			return DataError{
				code: DefaultEVMErrorCode,
				Msg:  err.Error(),
				data: RpcNullData,
			}
		}
		if method == RpcEthGetBlockByHash {
			return DataError{
				code: DefaultEVMErrorCode,
				Msg:  realErr.Error(),
				data: RpcNullData,
			}
		}
		lastSeg := strings.LastIndexAny(realErr.Log, "]")
		if lastSeg < 0 {
			return DataError{
				code: DefaultEVMErrorCode,
				Msg:  err.Error(),
				data: RpcNullData,
			}
		}
		marshaler := realErr.Log[0 : lastSeg+1]
		e = json.Unmarshal([]byte(marshaler), &logs)
		if e != nil {
			return DataError{
				code: DefaultEVMErrorCode,
				Msg:  err.Error(),
				data: RpcNullData,
			}
		}
		m := genericStringMap(logs)
		if m == nil {
			return DataError{
				code: DefaultEVMErrorCode,
				Msg:  err.Error(),
				data: "null",
			}
		}
		//if there have multi error type of EVM, this need a reactor mode to process error
		revert, f := m[vm.ErrExecutionReverted.Error()]
		if !f {
			revert = RpcUnknowErr
		}
		data, f := m[types.ErrorHexData]
		if !f {
			data = RpcNullData
		}
		switch method {
		case RpcEthEstimateGas:
			return DataError{
				code: VMExecuteExceptionInEstimate,
				Msg:  revert,
				data: data,
			}
		case RpcEthCall:
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
		data: RpcNullData,
	}
}

func genericStringMap(s []string) map[string]string {
	var ret = make(map[string]string, 0)
	if len(s)%2 != 0 {
		return nil
	}
	for i := 0; i < len(s); i += 2 {
		ret[s[i]] = s[i+1]
	}
	return ret
}
