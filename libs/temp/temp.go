package temp

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var (
	cmHandles  = make(map[string]map[string]func(data []byte, signers []sdk.AccAddress) (sdk.Msg, error))
	evmHandles = make(map[string]func(data []byte) ([]byte, error))
	evmParser  func(msg sdk.Msg) (string, string, []byte, error)
)

func Register(moduleName, funcName string, create func(data []byte, signers []sdk.AccAddress) (sdk.Msg, error)) {
	if create == nil {
		panic("Execute: Register driver is nil")
	}
	v, ok := cmHandles[moduleName]
	if !ok {
		v = make(map[string]func(data []byte, signers []sdk.AccAddress) (sdk.Msg, error))
	}
	v[funcName] = create
	cmHandles[moduleName] = v
}

func LoadNewMsg(msg sdk.Msg) (sdk.Msg, error) {
	moduleName, funcName, data, err := evmParser(msg)
	if err != nil {
		return nil, err
	}
	if module, ok := cmHandles[moduleName]; ok {
		if fn, ok := module[funcName]; ok {
			return fn(data, msg.GetSigners())
		}
	}
	return nil, fmt.Errorf("not find handle")
}

func RegisterEvm(name string, create func(data []byte) ([]byte, error)) {
	if create == nil {
		panic("Execute: Register driver is nil")
	}
	evmHandles[name] = create
}

func RegisterEvmParser(create func(msg sdk.Msg) (string, string, []byte, error)) {
	if create == nil {
		panic("Execute: Register driver is nil")
	}
	evmParser = create
}

func LoadEvmFunc(name string, data []byte) ([]byte, error) {
	if v, ok := evmHandles[name]; ok {
		return v(data)
	}
	return nil, fmt.Errorf("not find evm handle")
}
