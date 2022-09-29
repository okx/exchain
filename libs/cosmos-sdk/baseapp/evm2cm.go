package baseapp

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var (
	cmHandles          = make(map[string]map[string]func(data []byte, signers []sdk.AccAddress) (sdk.Msg, error))
	evmResultConverter func(data []byte) ([]byte, error)
	evmConvertJudge    func(msg sdk.Msg) (*CMTxParam, []byte, bool)
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

func RegisterEvmResultConverter(create func(data []byte) ([]byte, error)) {
	if create == nil {
		panic("Execute: Register EvmResultConverter is nil")
	}
	evmResultConverter = create
}

func RegisterEvmConvertJudge(create func(msg sdk.Msg) (*CMTxParam, []byte, bool)) {
	if create == nil {
		panic("Execute: Register EvmConvertJudge is nil")
	}
	evmConvertJudge = create
}

func ConvertMsg(cmtx *CMTxParam, msg sdk.Msg) (sdk.Msg, error) {
	if module, ok := cmHandles[cmtx.Module]; ok {
		if fn, ok := module[cmtx.Function]; ok {
			return fn([]byte(cmtx.Data), msg.GetSigners())
		}
	}
	return nil, fmt.Errorf("not find handle")
}

func EvmResultConvert(data []byte) ([]byte, error) {
	return evmResultConverter(data)
}

func (app *BaseApp) JudgeEvmConvert(ctx sdk.Context, msg sdk.Msg) (*CMTxParam, bool) {
	cmtp, addr, ok := evmConvertJudge(msg)
	if !ok {
		return nil, false
	}
	if app.EvmSysContractAddressHandler == nil {
		return nil, false
	}

	ok = app.EvmSysContractAddressHandler(ctx, addr)
	return cmtp, ok
}

func (app *BaseApp) SetEvmSysContractAddressHandler(handler sdk.EvmSysContractAddressHandler) {
	if app.sealed {
		panic("SetEvmSysContractAddressHandler() on sealed BaseApp")
	}
	app.EvmSysContractAddressHandler = handler
}
