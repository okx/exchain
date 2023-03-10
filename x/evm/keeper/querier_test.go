package keeper_test

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/okx/okbchain/x/evm/types"
	"math/big"
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
			suite.app.EvmKeeper.SetBlockHeight(suite.ctx, hash, 8)
		}, true},
		{"fail hash to height", []string{types.QueryHashToHeight, "0x00"}, func() {
			suite.app.EvmKeeper.SetBlockHeight(suite.ctx, hash, 8)
		}, false},
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
		}, true},
		{"account", []string{types.QueryAccount, "0x0"}, func() {}, true},
		{"exportAccount", []string{types.QueryExportAccount, suite.address.String()}, func() {
			for i := 0; i < 5; i++ {
				suite.stateDB.WithContext(suite.ctx).SetState(suite.address, ethcmn.BytesToHash([]byte(fmt.Sprintf("key%d", i))), ethcmn.BytesToHash([]byte(fmt.Sprintf("value%d", i))))
			}
			suite.stateDB.WithContext(suite.ctx).Finalise(false)
		}, true},
		{"unknown request", []string{"other"}, func() {}, false},
		{"parameters", []string{types.QueryParameters}, func() {}, true},
		{"storage by key", []string{types.QueryStorageByKey, "0xE3Db5e3cfDbBa56FfdDED5792DaAB8A2DC9c52c4", "key"}, func() {}, true},
		{"storage height to hash", []string{types.QueryHeightToHash, "1"}, func() {}, true},
		//{"storage section", []string{types.QuerySection, "1"}, func() {}, true},
		{"contract deploy white list", []string{types.QueryContractDeploymentWhitelist}, func() {}, true},
		{"contract blocked list", []string{types.QueryContractBlockedList}, func() {}, true},
		{"contract method blocked list", []string{types.QueryContractMethodBlockedList}, func() {}, true},
	}

	for i, tc := range testCases {
		suite.Run("", func() {
			//nolint
			tc := tc
			suite.SetupTest() // reset
			//nolint
			tc.malleate()
			fmt.Println(tc.msg)
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
