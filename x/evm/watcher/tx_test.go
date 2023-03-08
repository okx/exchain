package watcher

import (
	"bytes"
	"math/big"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/suite"

	okexchaincodec "github.com/okex/exchain/app/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	authclient "github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	tm "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/global"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	etypes "github.com/okex/exchain/x/evm/types"
)

var (
	evmAmountZero = big.NewInt(0)
	evmGasLimit   = uint64(1000000)
	evmGasPrice   = big.NewInt(10000)
	evmChainID    = big.NewInt(3)

	nonce0 = uint64(0)
)

type TxTestSuite struct {
	suite.Suite
	Watcher   Watcher
	TxDecoder sdk.TxDecoder
	height    int64
}

// only used for comparing mockTx and ethTx in Case 2
func realTxBoolCompare(a sdk.Tx, b sdk.Tx) bool {
	// only Raw and Hash are compared, others are nil
	RawCmpResult := bytes.Compare(a.GetRaw(), b.GetRaw())
	HashCmpResult := bytes.Compare(a.TxHash(), b.TxHash())
	if RawCmpResult == 0 && HashCmpResult == 0 {
		return true
	}
	return false
}

func TestWatcherTx(t *testing.T) {
	suite.Run(t, new(TxTestSuite))
}

func (suite *TxTestSuite) TestGetRealTx() {
	//Decoder Settings
	codecProxy, _ := okexchaincodec.MakeCodecSuit(module.NewBasicManager())
	suite.TxDecoder = etypes.TxDecoder(codecProxy)
	suite.height = 10
	tmtypes.UnittestOnlySetMilestoneVenusHeight(1)
	global.SetGlobalHeight(suite.height)

	testCases := []struct {
		title   string
		buildTx func() (tm.TxEssentials, sdk.Tx)
	}{
		{
			title: "Tx directly asserted as realTx",
			buildTx: func() (tm.TxEssentials, sdk.Tx) {
				realTx := etypes.NewMsgEthereumTx(1, nil, big.NewInt(1), 1, nil, nil)
				return realTx, realTx
			},
		},
		{
			title: "Tx converted to realTx by txDecoder",
			buildTx: func() (tm.TxEssentials, sdk.Tx) {
				bytecode := ethcommon.FromHex("0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a36102c4806100dc6000396000f3fe608060405234801561001057600080fd5b5060043610610053576000357c010000000000000000000000000000000000000000000000000000000090048063893d20e814610058578063a6f9dae1146100a2575b600080fd5b6100606100e6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100e4600480360360208110156100b857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061010f565b005b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101d1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f43616c6c6572206973206e6f74206f776e65720000000000000000000000000081525060200191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505056fea265627a7a72315820f397f2733a89198bc7fed0764083694c5b828791f39ebcbc9e414bccef14b48064736f6c63430005100032")
				tx := etypes.NewMsgEthereumTx(nonce0, nil, evmAmountZero, evmGasLimit, evmGasPrice, bytecode)
				txBytes, err := authclient.GetTxEncoder(nil, authclient.WithEthereumTx())(tx)
				suite.Require().NoError(err)
				tx.SetRaw(txBytes)
				tx.SetTxHash(tmtypes.Tx(txBytes).Hash(suite.height))
				mockTx := tm.MockTx{txBytes, tx.TxHash(), tx.GetFrom(), tx.GetNonce(), tx.GetGasPrice()}
				return mockTx, tx
			},
		},
		{
			title: "Tx convertion error", //because tx bytes are empty
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
			resrTx, err := suite.Watcher.getRealTx(Tx, suite.TxDecoder)
			if err != nil {
				suite.Require().Nil(realTx)
			} else {
				suite.Require().True(realTxBoolCompare(resrTx, realTx), "%s error, convert Tx error", tc.title)
			}
		})
	}
}
