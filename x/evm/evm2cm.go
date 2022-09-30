package evm

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
)

var (
	ErrInputDataSize = errors.New("the input data size is error")
	ErrNoMatchParam  = errors.New("no match the abi param")
	sysABIParser     *types.ABI
)

const (
	sysContractABI            = `[{"inputs":[{"internalType":"string","name":"data","type":"string"}],"name":"invoke","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
	sysContractInvokeFunction = "invoke"
)

func init() {
	RegisterHandle()
	var err error
	sysABIParser, err = types.NewABI(sysContractABI)
	if err != nil {
		panic(fmt.Sprintln("init system abi fail", err.Error()))
	}
}

func RegisterHandle() {
	baseapp.RegisterEvmResultConverter(EncodeResultData)
	baseapp.RegisterEvmConvertJudge(EvmConvertJudge)
	baseapp.RegisterEvmParamParse(EvmParamParse)
}

func EvmParamParse(msg sdk.Msg) (*baseapp.CMTxParam, error) {
	evmTx, ok := msg.(*types.MsgEthereumTx)
	if !ok {
		return nil, fmt.Errorf("msg type is not a MsgEthereumTx")
	}
	cmtp, err := ParseContractParam(evmTx.Data.Payload)
	if err != nil {
		return nil, err
	}
	return cmtp, nil
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

func ParseContractParam(input []byte) (*baseapp.CMTxParam, error) {
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

func DecodeParam(data []byte) (*baseapp.CMTxParam, error) {
	value, err := hex.DecodeString(string(data)) // this is json fmt
	if err != nil {
		return nil, err
	}
	cmtx := &baseapp.CMTxParam{}
	err = json.Unmarshal(value, cmtx)
	if err != nil {
		return nil, err
	}
	return cmtx, nil
}

func EncodeResultData(data []byte) ([]byte, error) {
	ethHash := common.BytesToHash(data)
	return types.EncodeResultData(&types.ResultData{TxHash: ethHash})
}
