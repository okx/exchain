package proxy

import (
	abcicli "github.com/okex/exchain/libs/tendermint/abci/client"
	"github.com/okex/exchain/libs/tendermint/abci/types"
)

//----------------------------------------

type AppConnTest interface {
	EchoAsync(string) *abcicli.ReqRes
	FlushSync() error
	InfoSync(types.RequestInfo) (*types.ResponseInfo, error)
}

type appConnTest struct {
	appConn abcicli.Client
}

func NewAppConnTest(appConn abcicli.Client) AppConnTest {
	return &appConnTest{appConn}
}

func (app *appConnTest) EchoAsync(msg string) *abcicli.ReqRes {
	return app.appConn.EchoAsync(msg)
}

func (app *appConnTest) FlushSync() error {
	return app.appConn.FlushSync()
}

func (app *appConnTest) InfoSync(req types.RequestInfo) (*types.ResponseInfo, error) {
	return app.appConn.InfoSync(req)
}
