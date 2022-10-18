package types

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type ABI struct {
	*abi.ABI
}

func NewABI(data string) (*ABI, error) {
	parsed, err := abi.JSON(strings.NewReader(data))
	if err != nil {
		return nil, err
	}
	return &ABI{ABI: &parsed}, nil
}

func (a *ABI) DecodeInputParam(methodName string, data []byte) ([]interface{}, error) {
	if len(data) <= 4 {
		return nil, fmt.Errorf("method %s data is nil", methodName)
	}
	method, ok := a.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s is not exist in abi", methodName)
	}
	return method.Inputs.Unpack(data[4:])
}

func (a *ABI) IsMatchFunction(methodName string, data []byte) bool {
	if len(data) < 4 {
		return false
	}
	method, ok := a.Methods[methodName]
	if !ok {
		return false
	}
	if bytes.Equal(method.ID, data[:4]) {
		return true
	}
	return false
}

func (a *ABI) EncodeOutput(methodName string, data []byte) ([]byte, error) {
	method, ok := a.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s is not exist in abi", methodName)
	}
	return method.Outputs.PackValues([]interface{}{string(data)})
}
