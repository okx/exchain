package keeper_test

import "github.com/okex/exchain/x/feesplit/types"

func (suite *KeeperTestSuite) TestParams() {
	params := suite.app.FeeSplitKeeper.GetParams(suite.ctx)
	params.EnableFeeSplit = true
	suite.Require().Equal(types.DefaultParams(), params)
	params.EnableFeeSplit = false
	suite.app.FeeSplitKeeper.SetParams(suite.ctx, params)
	newParams := suite.app.FeeSplitKeeper.GetParams(suite.ctx)
	suite.Require().Equal(newParams, params)
}
