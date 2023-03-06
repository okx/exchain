package keeper_test

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/erc20/types"
)

func (suite *KeeperTestSuite) TestQuerier() {

	testCases := []struct {
		msg      string
		path     []string
		malleate func()
		req      abci.RequestQuery
		expPass  bool
	}{
		{"unknown request", []string{"other"}, func() {}, abci.RequestQuery{}, false},
		{"parameters", []string{types.QueryParameters}, func() {}, abci.RequestQuery{}, true},
		{"all mapping", []string{types.QueryTokenMapping}, func() {
			denom1 := "testdenom1"
			denom2 := "testdenom2"

			autoContract := common.BigToAddress(big.NewInt(1))
			externalContract := common.BigToAddress(big.NewInt(2))
			suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, denom1, autoContract)
			suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, denom2, externalContract)
		}, abci.RequestQuery{}, true},
		{"contract by denom", []string{types.QueryContractByDenom}, func() {
			denom1 := "testdenom1"
			autoContract := common.BigToAddress(big.NewInt(1))
			suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, denom1, autoContract)
		}, abci.RequestQuery{Data: []byte(`{"denom":"testdenom1"}`)}, true},
		{"denom by contract", []string{types.QueryDenomByContract}, func() {
			denom1 := "testdenom1"
			autoContract := common.BigToAddress(big.NewInt(1))
			suite.app.Erc20Keeper.SetContractForDenom(suite.ctx, denom1, autoContract)
		}, abci.RequestQuery{Data: []byte(`{"contract":"0x01"}`)}, true},
		{"contract tem", []string{types.QueryContractTem}, func() {}, abci.RequestQuery{}, true},
	}

	for i, tc := range testCases {
		suite.Run("", func() {
			suite.SetupTest() // reset
			tc.malleate()

			bz, err := suite.querier(suite.ctx, tc.path, tc.req)
			if tc.expPass {
				suite.Require().NoError(err, "valid test %d failed: %s", i, tc.msg)
				suite.Require().NotZero(len(bz))
			} else {
				suite.Require().Error(err, "invalid test %d passed: %s", i, tc.msg)
			}
		})
	}
}
