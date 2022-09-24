package baseapp

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var (
	cmHandles          = make(map[string]map[string]func(data []byte, signers []sdk.AccAddress) (sdk.Msg, error))
	evmMsgParser       func(msg sdk.Msg) (*CMTxParam, error)
	evmResultConverter func(data []byte) ([]byte, error)
	evmConvertJudge    func(filter map[string]struct{}, msg sdk.Msg) bool
	ContractAddr       = make(map[string]struct{}) // evm to cm contract address
)

type CMTxParam struct {
	Module   string `json:"module"`
	Function string `json:"function"`
	Data     string `json:"data"`
}

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

func ConvertMsg(msg sdk.Msg) (sdk.Msg, error) {
	cmtx, err := evmMsgParser(msg)
	if err != nil {
		return nil, err
	}
	if module, ok := cmHandles[cmtx.Module]; ok {
		if fn, ok := module[cmtx.Function]; ok {
			return fn([]byte(cmtx.Data), msg.GetSigners())
		}
	}
	return nil, fmt.Errorf("not find handle")
}

func RegisterEvmResultConverter(create func(data []byte) ([]byte, error)) {
	if create == nil {
		panic("Execute: Register EvmResultConverter is nil")
	}
	evmResultConverter = create
}

func EvmResultConvert(data []byte) ([]byte, error) {
	return evmResultConverter(data)
}

func RegisterEvmMsgParser(create func(msg sdk.Msg) (*CMTxParam, error)) {
	if create == nil {
		panic("Execute: Register EvmMsgParser is nil")
	}
	evmMsgParser = create
}

func RegisterEvmConvertJudge(create func(filter map[string]struct{}, msg sdk.Msg) bool) {
	if create == nil {
		panic("Execute: Register EvmConvertJudge is nil")
	}
	evmConvertJudge = create
}

// Todo add fork height
func IsNeedEvmConvert(msg sdk.Msg) bool {
	return evmConvertJudge(ContractAddr, msg)
}
