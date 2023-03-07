package erc20_test

import (
	"testing"
	"time"

	"github.com/okx/exchain/app"
	sdk "github.com/okx/exchain/libs/cosmos-sdk/types"
	abci "github.com/okx/exchain/libs/tendermint/abci/types"
	"github.com/okx/exchain/x/erc20"
	"github.com/okx/exchain/x/erc20/types"
	"github.com/stretchr/testify/suite"
)

type Erc20TestSuite struct {
	suite.Suite

	ctx     sdk.Context
	handler sdk.Handler
	app     *app.OKExChainApp
}

func TestErc20TestSuite(t *testing.T) {
	suite.Run(t, new(Erc20TestSuite))
}

func (suite *Erc20TestSuite) SetupTest() {
	checkTx := false

	suite.app = app.Setup(false)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: "ethermint-3", Time: time.Now().UTC()})
	suite.handler = erc20.NewHandler(suite.app.Erc20Keeper)
	suite.app.Erc20Keeper.SetParams(suite.ctx, types.DefaultParams())
}
