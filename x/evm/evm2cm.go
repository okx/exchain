package evm

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

var (
	ErrNoMatchParam = errors.New("no match the abi param")
	sysABIParser    *types.ABI
)

const (
	sysContractABI            = `[{"inputs":[{"internalType":"string","name":"_data","type":"string"}],"name":"invoke","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"string","name":"input","type":"string"}],"name":"query","outputs":[{"internalType":"string","name":"","type":"string"}],"stateMutability":"view","type":"function"}]`
	sysContractInvokeFunction = "invoke"
	sysContractQueryFunction  = "query"
)

func init() {
	RegisterHandle()
	var err error
	sysABIParser, err = types.NewABI(sysContractABI)
	if err != nil {
		panic(fmt.Sprintln("init system abi fail", err.Error()))
	}
}

func SysABI() *types.ABI {
	return sysABIParser
}

func RegisterHandle() {
	baseapp.RegisterEvmResultConverter(EncodeResultData)
	baseapp.RegisterEvmConvertJudge(EvmConvertJudge)
	baseapp.RegisterEvmParamParse(EvmParamParse)
}

func EvmParamParse(msg sdk.Msg) ([]byte, error) {
	evmTx, ok := msg.(*types.MsgEthereumTx)
	if !ok {
		return nil, fmt.Errorf("msg type is not a MsgEthereumTx")
	}
	value, err := ParseContractParam(evmTx.Data.Payload)
	if err != nil {
		return nil, err
	}
	return value, nil
}

func EvmConvertJudge(msg sdk.Msg) ([]byte, bool) {
	if msg.Route() != types.ModuleName {
		return nil, false
	}
	evmTx, ok := msg.(*types.MsgEthereumTx)
	if !ok || evmTx.Data.Recipient == nil {
		return nil, false
	}
	if !sysABIParser.IsMatchFunction(sysContractInvokeFunction, evmTx.Data.Payload) {
		return nil, false
	}
	return evmTx.Data.Recipient[:], true
}

func ParseContractParam(input []byte) ([]byte, error) {
	res, err := sysABIParser.DecodeInputParam(sysContractInvokeFunction, input)
	if err != nil {
		return nil, err
	}
	if len(res) != 1 {
		return nil, ErrNoMatchParam
	}
	v, ok := res[0].(string)
	if !ok {
		return nil, ErrNoMatchParam
	}
	return DecodeParam([]byte(v))
}

func DecodeParam(data []byte) ([]byte, error) {
	value, err := hex.DecodeString(string(data)) // this is json fmt
	if err != nil {
		return nil, err
	}
	return value, nil
}

func EncodeResultData(txHash, data []byte) ([]byte, error) {
	ethHash := common.BytesToHash(txHash)
	return types.EncodeResultData(&types.ResultData{Ret: data, TxHash: ethHash})
}

func IsMatchSystemContractFunction(data []byte) bool {
	return sysABIParser.IsMatchFunction(sysContractInvokeFunction, data)
}

func IsMatchSystemContractQuery(data []byte) bool {
	return sysABIParser.IsMatchFunction(sysContractQueryFunction, data)
}

func EncodeQueryOutput(data []byte) ([]byte, error) {
	return sysABIParser.EncodeOutput(sysContractQueryFunction, data)
}
