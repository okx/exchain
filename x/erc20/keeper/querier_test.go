package keeper_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/erc20/types"
)

func (suite *KeeperTestSuite) TestQuerier() {

	testCases := []struct {
		msg      string
		path     []string
		malleate func()
		expPass  bool
	}{
		{"unknown request", []string{"other"}, func() {}, false},
		{"parameters", []string{types.QueryParameters}, func() {}, true},
		{"all mapping", []string{types.QueryAllMapping}, func() {
			denom1 := "testdenom1"
			denom2 := "testdenom2"

			autoContract := common.BigToAddress(big.NewInt(1))
			externalContract := common.BigToAddress(big.NewInt(2))
			suite.app.Erc20Keeper.SetExternalContractForDenom(suite.ctx, denom1, autoContract)
			suite.app.Erc20Keeper.SetExternalContractForDenom(suite.ctx, denom2, externalContract)
		}, true},
	}

	for i, tc := range testCases {
		suite.Run("", func() {
			suite.SetupTest() // reset
			tc.malleate()

			bz, err := suite.querier(suite.ctx, tc.path, abci.RequestQuery{})
			if tc.expPass {
				suite.Require().NoError(err, "valid test %d failed: %s", i, tc.msg)
				suite.Require().NotZero(len(bz))
			} else {
				suite.Require().Error(err, "invalid test %d passed: %s", i, tc.msg)
			}
		})
	}
}
