package baseapp

import (
	"fmt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var (
	cmHandles          = make(map[string]map[string]*CMHandle)
	evmResultConverter func(txHash, data []byte) ([]byte, error)
	evmConvertJudge    func(msg sdk.Msg) ([]byte, bool)
	evmParamParse      func(msg sdk.Msg) (*CMTxParam, error)
)

type CMTxParam struct {
	Module   string `json:"module"`
	Function string `json:"function"`
	Data     string `json:"data"`
}

type CMHandle struct {
	fn     func(data []byte, signers []sdk.AccAddress) (sdk.Msg, error)
	height int64
}

func NewCMHandle(fn func(data []byte, signers []sdk.AccAddress) (sdk.Msg, error), height int64) *CMHandle {
	return &CMHandle{
		fn:     fn,
		height: height,
	}
}

func RegisterCmHandle(moduleName, funcName string, create *CMHandle) {
	if create == nil {
		panic("Register CmHandle is nil")
	}
	v, ok := cmHandles[moduleName]
	if !ok {
		v = make(map[string]*CMHandle)
	}
	if _, dup := v[funcName]; dup {
		panic("Register CmHandle twice for same module and func " + moduleName + funcName)
	}
	v[funcName] = create
	cmHandles[moduleName] = v
}

func RegisterEvmResultConverter(create func(txHash, data []byte) ([]byte, error)) {
	if create == nil {
		panic("Register EvmResultConverter is nil")
	}
	evmResultConverter = create
}

func RegisterEvmParamParse(create func(msg sdk.Msg) (*CMTxParam, error)) {
	if create == nil {
		panic("Register EvmParamParse is nil")
	}
	evmParamParse = create
}

func RegisterEvmConvertJudge(create func(msg sdk.Msg) ([]byte, bool)) {
	if create == nil {
		panic("Register EvmConvertJudge is nil")
	}
	evmConvertJudge = create
}

func ConvertMsg(msg sdk.Msg, height int64) (sdk.Msg, error) {
	cmtx, err := evmParamParse(msg)
	if err != nil {
		return nil, err
	}
	if module, ok := cmHandles[cmtx.Module]; ok {
		cmh, ok := module[cmtx.Function]
		if ok && height >= cmh.height {
			return cmh.fn([]byte(cmtx.Data), msg.GetSigners())
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
