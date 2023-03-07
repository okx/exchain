package watcher_test

import (
	okexchaincodec "github.com/okex/exchain/app/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	tm "github.com/okex/exchain/libs/tendermint/abci/types"
	etypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/stretchr/testify/suite"
	"math/big"
	"testing"
)

type TxTestSuite struct {
	suite.Suite
	Watcher   watcher.Watcher
	TxDecoder sdk.TxDecoder
}

// 要在最开始通过函数建立tesing.T和suite的关系
func TestWatcherTx(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}

func (suite *TxTestSuite) TestGetRealTx() {
	codecProxy, _ := okexchaincodec.MakeCodecSuit(module.NewBasicManager())
	suite.TxDecoder = etypes.TxDecoder(codecProxy)

	testCases := []struct {
		title   string
		buildTx func() (tm.TxEssentials, sdk.Tx)
	}{
		{
			//直接生成一个sdk.Tx类型的交易
			title: "Tx directly assert as realTx",
			buildTx: func() (tm.TxEssentials, sdk.Tx) {
				realTx := etypes.NewMsgEthereumTx(1, nil, big.NewInt(1), 1, nil, nil)
				return realTx, realTx
			},
		},
		{
			//生成一个可以被转化的交易
			title: "Tx convert to realTx by txDecoder",
			buildTx: func() (tm.TxEssentials, sdk.Tx) {

				mockTx := tm.MockTx{}
				realTx := etypes.NewMsgEthereumTx(1, nil, big.NewInt(1), 1, nil, nil)
				return mockTx, realTx
			},
		},
		{
			title: "Tx convert error", //tx bytes are empty
			buildTx: func() (tm.TxEssentials, sdk.Tx) {
				mockTx := tm.MockTx{}
				return mockTx, nil
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			Tx, realTx := tc.buildTx()
			suite.Require().NotNil(Tx)
			resrTx, err := suite.Watcher.GetRealTx(Tx, suite.TxDecoder)
			if err != nil {
				suite.Require().Nil(realTx)
			} else {
				suite.Require().True(resrTx == realTx, "%s error, convert Tx error", tc.title)
			}
		})
	}
}
