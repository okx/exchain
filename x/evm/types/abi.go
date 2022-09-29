package types

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"strings"
)

type ABI struct {
	*abi.ABI
}

func NewABI(data string) (*ABI, error) {
	parsed, err := abi.JSON(strings.NewReader(data))
	if err != nil {
		fmt.Println("Can't generate ABI struct ", err)
		return nil, err
	}
	return &ABI{ABI: &parsed}, nil
}

func (a *ABI) DecodeInputParam(methodName string, data []byte) (map[string]interface{}, error) {
	if len(data) <= 4 {
		return nil, fmt.Errorf("method %s data is nil", methodName)
	}
	method, ok := a.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s is not exist in abi", methodName)
	}
	resmp := make(map[string]interface{})
	err := method.Inputs.UnpackIntoMap(resmp, data[4:])
	return resmp, err
}

//func DecodeOneInputParam() ([]byte, error) {
//	abi.ArgumentMarshaling{Name: "a", Type: "uint256"}
//	abi.
//	abi.Argument{
//		Name: "",
//		Type: "",
//	}
//	abi.Arguments{Argument{Argument}}
//}

func (a *ABI) EncodeOutput(methodName string, data []byte) ([]byte, error) {
	method, ok := a.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s is not exist in abi", methodName)
	}
	return method.Outputs.PackValues([]interface{}{string(data)})
}
