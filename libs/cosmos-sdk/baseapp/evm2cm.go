package baseapp

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var (
	cmHandles          = make(map[string]map[string]func(data []byte, signers []sdk.AccAddress) (sdk.Msg, error))
	evmResultConverter func(txHash, data []byte) ([]byte, error)
	evmConvertJudge    func(msg sdk.Msg) ([]byte, bool)
	evmParamParse      func(msg sdk.Msg) (*CMTxParam, error)
)

type CMTxParam struct {
	Module   string `json:"module"`
	Function string `json:"function"`
	Data     string `json:"data"`
}

func RegisterCmHandle(moduleName, funcName string, create func(data []byte, signers []sdk.AccAddress) (sdk.Msg, error)) {
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

func RegisterEvmResultConverter(create func(txHash, data []byte) ([]byte, error)) {
	if create == nil {
		panic("Execute: Register EvmResultConverter is nil")
	}
	evmResultConverter = create
}

func RegisterEvmParamParse(create func(msg sdk.Msg) (*CMTxParam, error)) {
	if create == nil {
		panic("Execute: Register EvmParamParse is nil")
	}
	evmParamParse = create
}

func RegisterEvmConvertJudge(create func(msg sdk.Msg) ([]byte, bool)) {
	if create == nil {
		panic("Execute: Register EvmConvertJudge is nil")
	}
	evmConvertJudge = create
}

func ConvertMsg(msg sdk.Msg) (sdk.Msg, error) {
	cmtx, err := evmParamParse(msg)
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

func EvmResultConvert(txHash, data []byte) ([]byte, error) {
	return evmResultConverter(txHash, data)
}

func (app *BaseApp) JudgeEvmConvert(ctx sdk.Context, msg sdk.Msg) bool {
	if app.EvmSysContractAddressHandler == nil ||
		evmConvertJudge == nil ||
		evmParamParse == nil ||
		evmResultConverter == nil {
		return false
	}
	addr, ok := evmConvertJudge(msg)
	if !ok || len(addr) == 0 {
		return false
	}
	return app.EvmSysContractAddressHandler(ctx, addr)
}

func (app *BaseApp) SetEvmSysContractAddressHandler(handler sdk.EvmSysContractAddressHandler) {
	if app.sealed {
		panic("SetEvmSysContractAddressHandler() on sealed BaseApp")
	}
	app.EvmSysContractAddressHandler = handler
}
