package types

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	SendToWasmEventName  = "__OKCSendToWasm"
	WasmCalledMethodName = "mintCW20"

	SendToEvmSubMsgName = "send-to-evm"
	EvmCalledMethodName = "mintERC20"
)

var (
	// SendToWasmEventName represent the signature of
	// `event __SendToWasmEventName(string wasmAddr,string recipient, string amount)`
	SendToWasmEvent abi.Event

	EvmABI abi.ABI
	//go:embed abi.json
	abiJson []byte
)

func init() {
	EvmABI, SendToWasmEvent = GetEVMABIConfig(abiJson)
}

type MintCW20Method struct {
	Amount    string `json:"amount"`
	Recipient string `json:"recipient"`
}

func GetMintCW20Input(amount, recipient string) ([]byte, error) {
	method := MintCW20Method{
		Amount:    amount,
		Recipient: recipient,
	}
	input := struct {
		Method MintCW20Method `json:"mint_c_w20"`
	}{
		Method: method,
	}
	return json.Marshal(input)
}

type MintERC20Method struct {
	ABI abi.ABI
}

func GetMintERC20Input(callerAddr string, recipient common.Address, amount *big.Int) ([]byte, error) {
	data, err := EvmABI.Pack(EvmCalledMethodName, callerAddr, recipient, amount)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetMintERC20Output(data []byte) (bool, error) {
	result, err := EvmABI.Unpack(EvmCalledMethodName, data)
	if err != nil {
		return false, err
	}
	if len(result) != 1 {
		return false, fmt.Errorf("%s method outputs must be one output", EvmCalledMethodName)
	}
	return result[0].(bool), nil
}

func GetEVMABIConfig(data []byte) (abi.ABI, abi.Event) {
	ret, err := abi.JSON(bytes.NewReader(data))
	if err != nil {
		panic(fmt.Errorf("json decode failed: %s", err.Error()))
	}
	event, ok := ret.Events[SendToWasmEventName]
	if !ok {
		panic(fmt.Errorf("abi must have event %s,%s,%s", SendToWasmEvent, ret, string(data)))
	}
	return ret, event
}
