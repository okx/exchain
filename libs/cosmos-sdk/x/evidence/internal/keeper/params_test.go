package keeper_test

import (
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/evidence/internal/types"
)

func (suite *KeeperTestSuite) TestParams() {
	ctx := suite.ctx
	ctx.SetRunTxMode(sdk.RunTxModeDeliver)
	suite.Equal(types.DefaultParams(), suite.keeper.GetParams(ctx))
	suite.Equal(types.DefaultMaxEvidenceAge, suite.keeper.MaxEvidenceAge(ctx))
}
