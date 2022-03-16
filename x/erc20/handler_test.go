package erc20_test

import (
	"testing"
	"time"

	"github.com/okex/exchain/app"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/erc20"
	"github.com/stretchr/testify/suite"
)

type CronosTestSuite struct {
	suite.Suite

	ctx     sdk.Context
	handler sdk.Handler
	app     *app.OKExChainApp
}

func TestCronosTestSuite(t *testing.T) {
	suite.Run(t, new(CronosTestSuite))
}

func (suite *CronosTestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-3", Time: time.Now().UTC()})
	suite.handler = erc20.NewHandler(suite.app.Erc20Keeper)
}
