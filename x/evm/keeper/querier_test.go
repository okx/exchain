package keeper_test

import (
	"fmt"
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/okexchain/x/evm/types"

	abci "github.com/tendermint/tendermint/abci/types"
)

func (suite *KeeperTestSuite) TestQuerier() {

	testCases := []struct {
		msg      string
		path     []string
		malleate func()
		expPass  bool
	}{
		{"balance", []string{types.QueryBalance, addrHex}, func() {
			suite.app.EvmKeeper.SetBalance(suite.ctx, suite.address, big.NewInt(5))
		}, true},
		//{"balance fail", []string{types.QueryBalance, "0x01232"}, func() {}, false},
		{"block number", []string{types.QueryBlockNumber}, func() {}, true},
		{"storage", []string{types.QueryStorage, "0x0", "0x0"}, func() {}, true},
		{"code", []string{types.QueryCode, "0x0"}, func() {}, true},
		{"hash to height", []string{types.QueryHashToHeight, hex}, func() {
			suite.app.EvmKeeper.SetBlockHash(suite.ctx, hash, 8)
		}, true},
		{"fail hash to height", []string{types.QueryHashToHeight, "0x00"}, func() {
			suite.app.EvmKeeper.SetBlockHash(suite.ctx, hash, 8)
		}, false},
		{"tx logs", []string{types.QueryTransactionLogs, "0x0"}, func() {}, true},
		{"bloom", []string{types.QueryBloom, "4"}, func() {
			testBloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
			suite.app.EvmKeeper.SetBlockBloom(suite.ctx, 4, testBloom)
		}, true},
		{"fail bloom height", []string{types.QueryBloom, ""}, func() {
			testBloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
			suite.app.EvmKeeper.SetBlockBloom(suite.ctx, 4, testBloom)
		}, false},
		{"fail bloom number", []string{types.QueryBloom, "4"}, func() {
			testBloom := ethtypes.BytesToBloom([]byte{0x1, 0x3})
			suite.app.EvmKeeper.SetBlockBloom(suite.ctx, 3, testBloom)
		}, false},
		{"logs", []string{types.QueryLogs, "0x0"}, func() {}, true},
		{"account", []string{types.QueryAccount, "0x0"}, func() {}, true},
		{"exportAccount", []string{types.QueryExportAccount, suite.address.String()}, func() {
			for i := 0; i < 5; i++ {
				suite.app.EvmKeeper.SetState(suite.ctx, suite.address, ethcmn.BytesToHash([]byte(fmt.Sprintf("key%d", i))), ethcmn.BytesToHash([]byte(fmt.Sprintf("value%d", i))))
			}
			suite.app.EvmKeeper.Finalise(suite.ctx, false)
		}, true},
		{"unknown request", []string{"other"}, func() {}, false},
		{"parameters", []string{types.QueryParameters}, func() {}, true},
	}

	for i, tc := range testCases {
		suite.Run("", func() {
			//nolint
			tc := tc
			suite.SetupTest() // reset
			//nolint
			tc.malleate()

			bz, err := suite.querier(suite.ctx, tc.path, abci.RequestQuery{})

			//nolint
			if tc.expPass {
				//nolint
				suite.Require().NoError(err, "valid test %d failed: %s", i, tc.msg)
				suite.Require().NotZero(len(bz))
			} else {
				//nolint
				suite.Require().Error(err, "invalid test %d passed: %s", i, tc.msg)
			}
		})
	}
}
