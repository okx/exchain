package eth

import (
	"fmt"
	eabi "github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

type ABI struct {
	abi *eabi.ABI
}

func NewABI(data string) (*ABI, error) {
	parsed, err := eabi.JSON(strings.NewReader(data))
	if err != nil {
		fmt.Println("Can't generate ABI struct ", err)
		return nil, err
	}
	return &ABI{abi: &parsed}, nil
}

func (a *ABI) DecodeInputParam(methodName string, data []byte) (map[string]interface{}, error) {
	if len(data) <= 4 {
		return nil, fmt.Errorf("method %s data is nil", methodName)
	}
	method, ok := a.abi.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s is not exist in abi", methodName)
	}
	resmp := make(map[string]interface{})
	method.Inputs.UnpackIntoMap(resmp, data[4:])
	return resmp, nil
}

func (a *ABI) EncodeOutput(methodName string, data []byte) ([]byte, error) {
	method, ok := a.abi.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s is not exist in abi", methodName)
	}
	return method.Outputs.PackValues([]interface{}{string(data)})
}
