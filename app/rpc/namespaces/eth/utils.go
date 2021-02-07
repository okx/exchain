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
	Code      uint32 `json:"code"`
	Log       string `json:"log"`
	Codespace string `json:"codespace"`
}

func (c cosmosError) Error() string {
	return c.Log
}

type wrapedEthError struct {
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

func newDataError(revert string, data string) *wrapedEthError {
	return &wrapedEthError{
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
				data: "null",
			}
		}
		lastSeg := strings.LastIndexAny(realErr.Log, "]")
		if lastSeg < 0 {
			return DataError{
				code: DefaultEVMErrorCode,
				Msg:  err.Error(),
				data: "null",
			}
		}
		marshaler := realErr.Log[0 : lastSeg+1]
		e = json.Unmarshal([]byte(marshaler), &logs)
		if e != nil {
			return DataError{
				code: DefaultEVMErrorCode,
				Msg:  err.Error(),
				data: "null",
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
			revert = "unknow"
		}
		data, f := m[types.ErrorHexData]
		if !f {
			data = "null"
		}
		switch method {
		case "eth_estimateGas":
			return DataError{
				code: VMExecuteExceptionInEstimate,
				Msg:  revert,
				data: data,
			}
		case "eth_call":
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
		data: "null",
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
