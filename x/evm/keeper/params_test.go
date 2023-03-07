package keeper_test

import (
	"github.com/okx/okbchain/x/evm/types"
)

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.EvmKeeper.GetParams(suite.ctx)
	suite.Require().Equal(types.DefaultParams(), params)
	suite.app.EvmKeeper.SetParams(suite.ctx, params)
	newParams := suite.app.EvmKeeper.GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)
	newParams = suite.app.EvmKeeper.GetParams(*suite.ctx.SetDeliver())
	suite.Require().Equal(newParams, params)
	types.GetEvmParamsCache().UpdateParams(params, false)
	newParams = suite.app.EvmKeeper.GetParams(*suite.ctx.SetDeliver())
	suite.Require().Equal(newParams, params)
}
