package baseapp

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	db "github.com/okex/exchain/libs/tm-db"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

type testMsg struct {
	route string
}

func (msg testMsg) Route() string                { return msg.route }
func (msg testMsg) Type() string                 { return "testMsg" }
func (msg testMsg) GetSigners() []sdk.AccAddress { return nil }
func (msg testMsg) GetSignBytes() []byte         { return nil }
func (msg testMsg) ValidateBasic() error         { return nil }

func TestRegisterCmHandle_ConvertMsg(t *testing.T) {
	testcases := []struct {
		module   string
		funcName string
	}{
		{
			module:   "module1",
			funcName: "test1",
		},
		{
			module:   "module1",
			funcName: "test2",
		},
		{
			module:   "module3",
			funcName: "test1",
		},
		{
			module:   "module4",
			funcName: "test1",
		},
		{
			module:   "module4",
			funcName: "test2",
		},
		{
			module:   "module4",
			funcName: "test3",
		},
	}
	for _, ts := range testcases {
		RegisterCmHandle(ts.module, ts.funcName, func(data []byte, signers []sdk.AccAddress) (sdk.Msg, error) {
			return nil, nil
		})
	}

	// check
	for _, ts := range testcases {
		RegisterEvmParamParse(func(msg sdk.Msg) (*CMTxParam, error) {
			return &CMTxParam{Module: ts.module, Function: ts.funcName}, nil
		})
		_, err := ConvertMsg(testMsg{})
		require.NoError(t, err)
	}
}

func TestJudgeEvmConvert(t *testing.T) {
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).
		With("module", "mock")
	testcases := []struct {
		app     *BaseApp
		fnInit  func(app *BaseApp)
		fnCheck func(is bool)
	}{
		{
			app:    NewBaseApp("test", logger, db.NewMemDB(), nil),
			fnInit: func(app *BaseApp) {},
			fnCheck: func(is bool) {
				require.False(t, is)
			},
		},
		{
			app: NewBaseApp("test", logger, db.NewMemDB(), nil),
			fnInit: func(app *BaseApp) {
				app.SetEvmSysContractAddressHandler(func(ctx sdk.Context, addr sdk.AccAddress) bool {
					return true
				})
				RegisterEvmParamParse(func(msg sdk.Msg) (*CMTxParam, error) {
					return nil, nil
				})
				RegisterEvmResultConverter(func(txHash, data []byte) ([]byte, error) {
					return nil, nil
				})
				RegisterEvmConvertJudge(func(msg sdk.Msg) ([]byte, bool) {
					return []byte{1, 2, 3}, true
				})
			},
			fnCheck: func(is bool) {
				require.True(t, is)
			},
		},
		{
			app:    NewBaseApp("test", logger, db.NewMemDB(), nil),
			fnInit: func(app *BaseApp) {},
			fnCheck: func(is bool) {
				require.False(t, is)
			},
		},
	}

	for _, ts := range testcases {
		ts.fnInit(ts.app)
		re := ts.app.JudgeEvmConvert(sdk.Context{}, nil)
		ts.fnCheck(re)
	}
}
