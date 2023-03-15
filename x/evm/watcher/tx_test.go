package watcher

import (
	"bytes"
	"math/big"
	"os"
	"testing"

	ethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/suite"

	okexchaincodec "github.com/okex/exchain/app/codec"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	app "github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/module"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authclient "github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	tm "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	ctypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/distribution/keeper"
	etypes "github.com/okex/exchain/x/evm/types"
	//"github.com/okex/exchain/x/gov"
)

var (
	evmAmountZero = big.NewInt(0)
	evmGasLimit   = uint64(1000000)
	evmGasPrice   = big.NewInt(10000)
	evmChainID    = big.NewInt(3)

	cosmosChainId = "ethermint-3"

	nonce0 = uint64(0)
	nonce1 = uint64(1)
	//generate fees for stdTx
	Coin10   = sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)
	Coin1000 = sdk.NewInt64Coin(sdk.DefaultBondDenom, 1000)
	fees     = auth.NewStdFee(21000, sdk.NewCoins(Coin10))

	txCoin10   = sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)
	txFees     = auth.NewStdFee(21000, sdk.NewCoins(txCoin10))
	sysCoins10 = keeper.NewTestSysCoins(10, 0)
	sysCoins90 = keeper.NewTestSysCoins(90, 0)

	govProposalID1 = uint64(1)
	memo           = "hello, memo"

	blockHeight = int64(2)

	accountNum = uint64(0)
)

type TxTestSuite struct {
	suite.Suite
	Watcher            Watcher
	TxDecoder          sdk.TxDecoder
	height             int64
	evmSenderPrivKey   ethsecp256k1.PrivKey
	evmContractAddress ethcommon.Address

	stdSenderPrivKey    ethsecp256k1.PrivKey
	stdSenderAccAddress sdk.AccAddress
	app                 *app.BaseApp
}

//func (suite *TxTestSuite) Ctx() sdk.Context {
//	return suite.app..GetDeliverStateCtx()
//}

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

// For generating DeliverTxResponse with DeliverTx
func (suite *TxTestSuite) SetupTest() {
	//suite.app = app.Setup(false, app.WithChainId(cosmosChainId))
	//params := etypes.DefaultParams()
	//params.EnableCreate = true
	//params.EnableCall = true
	//suite.app.EvmKeeper.SetParams(suite.Ctx(), params)

	suite.evmSenderPrivKey, _ = ethsecp256k1.GenerateKey()
	codecProxy, _ := okexchaincodec.MakeCodecSuit(module.NewBasicManager())
	suite.TxDecoder = etypes.TxDecoder(codecProxy)

	suite.Watcher = *(NewWatcher(log.NewTMLogger(os.Stdout)))

	//streamMetrics := monitor.DefaultStreamMetrics(monitor.DefaultPrometheusConfig())
	//suite.app.InfuraKeeper = infura.NewKeeper(nil, log.NewTMLogger(os.Stdout), streamMetrics)
}

func (suite *TxTestSuite) beginFakeBlock() {
	suite.evmSenderPrivKey, _ = ethsecp256k1.GenerateKey()
	suite.evmContractAddress = ethcrypto.CreateAddress(ethcommon.HexToAddress(suite.evmSenderPrivKey.PubKey().Address().String()), 0)
	//accountEvm := suite.app.AccountKeeper.NewAccountWithAddress(suite.Ctx(), suite.evmSenderPrivKey.PubKey().Address().Bytes())
	//accountEvm.SetAccountNumber(accountNum)
	//accountEvm.SetCoins(sdk.NewCoins(Coin1000))
	//suite.app.AccountKeeper.SetAccount(suite.Ctx(), accountEvm)
	//
	//suite.stdSenderPrivKey, _ = ethsecp256k1.GenerateKey()
	//suite.stdSenderAccAddress = sdk.AccAddress(suite.stdSenderPrivKey.PubKey().Address())
	//accountStd := suite.app.AccountKeeper.NewAccountWithAddress(suite.Ctx(), suite.stdSenderAccAddress.Bytes())
	//accountStd.SetAccountNumber(accountNum)
	//accountStd.SetCoins(sdk.NewCoins(Coin1000))
	//suite.app.AccountKeeper.SetAccount(suite.Ctx(), accountStd)
	//err := suite.app.BankKeeper.SetCoins(suite.Ctx(), suite.stdSenderAccAddress, sdk.NewCoins(Coin1000))
	//suite.Require().NoError(err)

	tmtypes.UnittestOnlySetMilestoneVenusHeight(blockHeight - 1)
	global.SetGlobalHeight(blockHeight - 1)
	//suite.app.BeginBlocker(suite.Ctx(), tm.RequestBeginBlock{Header: tm.Header{Height: blockHeight}})
}

func (suite *TxTestSuite) endFakeBlock() {
	//suite.app.EndBlocker(suite.Ctx(), tm.RequestEndBlock{})
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

// UT for createWatchTx & extractEvmTx
// unextractable evmTx as a testCase is unnecessary
// because the input of extractEvmTx can only be EvmTxType sdk.Tx
// MsgEvmTx always return a none-nil msg with the Tx itself in it.

func (suite *TxTestSuite) TestCreateWatchTx() {
	var oldWIndex uint64

	testCases := []struct {
		title            string
		buildTx          func() (sdk.Tx, WatchTx)
		wIndexCompEnable bool //Test the increase of w.evmTxIndex if extract success
	}{
		{
			title: "extractable evmTx with correct result",
			buildTx: func() (sdk.Tx, WatchTx) {
				evmTx := etypes.NewMsgEthereumTx(1, nil, big.NewInt(1), 1, nil, nil)
				evmMsg := evmTx.GetMsgs()
				extEvmTx, ok := evmMsg[0].(*etypes.MsgEthereumTx)
				suite.Require().True(ok, "extract emv Tx from Msg error, type assertion in testCase error")
				txMsg := NewEvmTx(extEvmTx, ethcommon.BytesToHash(evmTx.TxHash()), suite.Watcher.blockHash, suite.Watcher.height, suite.Watcher.GetEvmTxIndex())
				return evmTx, txMsg
			},
			wIndexCompEnable: true,
		},
		{
			title: "stdTx, nil return",
			buildTx: func() (sdk.Tx, WatchTx) {
				stdTx := auth.NewStdTx([]sdk.Msg{}, fees, nil, "")
				return stdTx, nil
			},
			wIndexCompEnable: false,
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			evmTx, watchTx := tc.buildTx()
			if tc.wIndexCompEnable {
				oldWIndex = suite.Watcher.GetEvmTxIndex()
			}
			resWatchTx := suite.Watcher.createWatchTx(evmTx)
			if tc.wIndexCompEnable {
				suite.Require().Equal(oldWIndex+1, suite.Watcher.GetEvmTxIndex(), "%s evmTxIndex increase error", tc.title)
			}
			//how to compare WatchTx?
			suite.Require().Equal(resWatchTx, watchTx, "%s error", tc.title)
		})
	}
}

func (suite *TxTestSuite) testUpdateCumulativeGasWithArg(txIndex uint64, gasUsed uint64) uint64 {
	cumulativeGasMap := suite.Watcher.cumulativeGas
	suite.Require().True(len(cumulativeGasMap) > 0, "UpdateCumulativeGas failed, no append item")
	var cumulativeGas uint64
	if len(cumulativeGasMap) == 1 {
		cumulativeGas = cumulativeGasMap[txIndex]
		suite.Require().True(cumulativeGas == gasUsed, "UpdateCumulativeGas failed, append item error")
	} else {
		expectedValue := cumulativeGasMap[txIndex-1] + gasUsed
		cumulativeGas = cumulativeGasMap[txIndex]
		suite.Require().True(cumulativeGas == expectedValue, "UpdateCumulativeGas failed, append item error")
	}
	return cumulativeGas
}

//func (suite *TxTestSuite) testOnSaveTransactionReceiptWithArg(receipt *watcher.TransactionReceipt) {
//	expectedKeeper, ok := (suite.Watcher.InfuraKeeper).(infura.Keeper)
//	suite.Require().True(ok, "Get Infura Keeper error")
//	expectedAddedReceipt := expectedKeeper.GetCache().GetTransactionReceipts()
//	suite.Require().Equal(expectedAddedReceipt[len(expectedAddedReceipt)-1], *receipt, "On Save Transaction Receipt error")
//	return
//}

func (suite *TxTestSuite) testAppendMsgTransactionReceiptWithArg(receipt *TransactionReceipt, TxHash ethcommon.Hash, index int) {
	expectedMsg := NewMsgTransactionReceipt(*receipt, TxHash)
	getBatch := suite.Watcher.batch
	suite.Require().True(len(getBatch) >= index+1, "No MsgTransactionReceipt is appended to batch")
	getMsg := getBatch[index]
	getMsgTxReceipt, ok := getMsg.(*MsgTransactionReceipt)
	suite.Require().True(ok, "Convert WatchMessage to MsgTransactionReceipt Error")
	suite.Require().Equal(*(expectedMsg.TransactionReceipt), *(getMsgTxReceipt.TransactionReceipt), "Append MsgTransactionReceipt to batch Error")
	suite.Require().Equal(expectedMsg.txHash, getMsgTxReceipt.txHash, "Append MsgTransactionReceipt to batch Error")
	return
}

// only watcher != nil && watchTx != nil is used
// Used for comparing:
// 1. The Updated Cumulative Gas
// 2. The Transaction Receipt saved in Watcher.InfraKeeper
// 3. The Msg of Transaction Receipt saved in Watcher.Batch
// These 3 Tests are written in separate *WithArg functions above.

func (suite *TxTestSuite) testSaveFailedReceiptWithArg(watchTx WatchTx, gasUsed uint64, index int) {
	//watchTxIndex should minus 1 because the evmTxIndex increases by 1
	//after creating a new watchTx
	watchTxIndex := watchTx.GetIndex()
	watchTxHash := watchTx.GetTxHash()
	RespcumulativeGas := suite.testUpdateCumulativeGasWithArg(watchTxIndex, gasUsed)
	receipt := watchTx.GetFailedReceipts(RespcumulativeGas, gasUsed)
	//suite.testOnSaveTransactionReceiptWithArg(receipt)
	suite.testAppendMsgTransactionReceiptWithArg(receipt, watchTxHash, index)
	return
}

func (suite *TxTestSuite) testSaveTransactionReceiptWithArg(status uint32, msg *etypes.MsgEthereumTx, txHash ethcommon.Hash, txIndex uint64, data *etypes.ResultData, gasUsed uint64, index int) {
	_ = suite.testUpdateCumulativeGasWithArg(txIndex, gasUsed)
	cumulativeGas := suite.Watcher.cumulativeGas[txIndex]
	receipt := newTransactionReceipt(status, msg, txHash, suite.Watcher.blockHash, txIndex, suite.Watcher.height, data, cumulativeGas, gasUsed)
	//suite.testOnSaveTransactionReceiptWithArg(&receipt)
	suite.testAppendMsgTransactionReceiptWithArg(&receipt, txHash, index)
	return
}

func (suite *TxTestSuite) testSaveTxWithArg(watchTx WatchTx, batchIndex int) {

	suite.Require().NotNil(suite.Watcher, "Watcher is nil when testing SaveTx")
	suite.Require().NotNil(watchTx, "watchTx is nil when testing SaveTx")

	//if suite.Watcher.InfuraKeeper != nil {
	//	ethTx := watchTx.GetTransaction()
	//	if ethTx != nil {
	//		//Test OnSaveTransaction
	//		expectedKeeper, ok := (suite.Watcher.InfuraKeeper).(infura.Keeper)
	//		suite.Require().True(ok, "Get Infura Keeper error")
	//		expectedAddedTransaction := expectedKeeper.GetCache().GetTransactions()
	//		suite.Require().Equal(expectedAddedTransaction[len(expectedAddedTransaction)-1], *ethTx, "On Save Transaction Receipt error")
	//	}
	//}
	txWatchMessage := watchTx.GetTxWatchMessage()
	if txWatchMessage != nil {
		respBatch := suite.Watcher.batch
		suite.Require().True(len(respBatch) >= batchIndex+1, "Append txWatchMessage to batch Error : count error")
		respMsg := respBatch[batchIndex]
		//suite.Require().Equal(txWatchMessage, respMsg, "Append txWatchMessage to batch Error : mismatch message")
		expectEthMsg, ok := txWatchMessage.(*MsgEthTx)
		suite.Require().True(ok, "Convert to MsgEthTx error in testSaveTxWithArg")
		respEthMsg, ok := respMsg.(*MsgEthTx)
		suite.Require().True(ok, "Convert to MsgEthTx error in testSaveTxWithArg")
		expectTx := expectEthMsg.Transaction
		respTx := respEthMsg.Transaction
		suite.Require().Equal(expectTx.tx, respTx.tx, "Append txWatchMessage to batch Error : mismatch message")
		suite.Require().Equal(expectTx.Hash, respTx.Hash, "Append txWatchMessage to batch Error : mismatch message")
		suite.Require().Equal(expectEthMsg.GetKey(), respEthMsg.GetKey(), "Append txWatchMessage to batch Error : mismatch message")
	}
	respBlockTxs := suite.Watcher.blockTxs[len(suite.Watcher.blockTxs)-1]
	suite.Require().Equal(watchTx.GetTxHash(), respBlockTxs, "Append Tx Hash error in saveTx")
}

func (suite *TxTestSuite) TestRecordTxAndFailedReceipt() {
	testCases := []struct {
		title            string
		watcherEnabled   bool
		buildInput       func() (tm.TxEssentials, *tm.ResponseDeliverTx) // build the input of the tested function
		genStdTxResponse func(tm.TxEssentials, *tm.ResponseDeliverTx) *MsgStdTransactionResponse
		genWatchTx       func(tm.TxEssentials) (WatchTx, sdk.Tx)
		numBatch         int
	}{
		{
			title:          "evmTx success with none-nil ResponseDeliverTx",
			watcherEnabled: true,
			buildInput: func() (tm.TxEssentials, *tm.ResponseDeliverTx) {
				//Create evm contract - Owner.sol
				bytecode := ethcommon.FromHex("0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a36102c4806100dc6000396000f3fe608060405234801561001057600080fd5b5060043610610053576000357c010000000000000000000000000000000000000000000000000000000090048063893d20e814610058578063a6f9dae1146100a2575b600080fd5b6100606100e6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100e4600480360360208110156100b857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061010f565b005b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101d1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f43616c6c6572206973206e6f74206f776e65720000000000000000000000000081525060200191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505056fea265627a7a72315820f397f2733a89198bc7fed0764083694c5b828791f39ebcbc9e414bccef14b48064736f6c63430005100032")
				tx := etypes.NewMsgEthereumTx(nonce0, nil, evmAmountZero, evmGasLimit, evmGasPrice, bytecode)
				tx.Sign(evmChainID, suite.evmSenderPrivKey.ToECDSA())
				resp := suite.app.DeliverRealTx(tx)
				return tx, &resp
			},
			genStdTxResponse: func(tx tm.TxEssentials, resp *tm.ResponseDeliverTx) *MsgStdTransactionResponse {
				result := &ctypes.ResultTx{
					Hash:     tx.TxHash(),
					Height:   int64(suite.Watcher.height),
					TxResult: *resp,
					Tx:       tx.GetRaw(),
				}
				newMsg := NewStdTransactionResponse(result, suite.Watcher.header.Time, ethcommon.BytesToHash(result.Hash))
				return newMsg
			},
			genWatchTx: func(tx tm.TxEssentials) (WatchTx, sdk.Tx) {
				evmTx, ok := tx.(sdk.Tx)
				suite.Require().True(ok, "evmTx generate WatchTx error")
				watchTx := suite.Watcher.createWatchTx(evmTx)
				return watchTx, evmTx
			},
			numBatch: 2,
		},
		{
			title:          "evmTx fail with none-nil ResponseDeliverTx",
			watcherEnabled: true,
			buildInput: func() (tm.TxEssentials, *tm.ResponseDeliverTx) {
				//Create evm contract - Owner.sol
				bytecode := ethcommon.FromHex("0xa6f9dae10000000000000000000000006a82e4a67715c8412a9114fbd2cbaefbc8181424")
				tx := etypes.NewMsgEthereumTx(nonce1, nil, evmAmountZero, evmGasLimit, evmGasPrice, bytecode)
				tx.Sign(evmChainID, suite.evmSenderPrivKey.ToECDSA())
				resp := suite.app.DeliverRealTx(tx)
				return tx, &resp
			},
			genStdTxResponse: func(tx tm.TxEssentials, resp *tm.ResponseDeliverTx) *MsgStdTransactionResponse {
				result := &ctypes.ResultTx{
					Hash:     tx.TxHash(),
					Height:   int64(suite.Watcher.height),
					TxResult: *resp,
					Tx:       tx.GetRaw(),
				}
				newMsg := NewStdTransactionResponse(result, suite.Watcher.header.Time, ethcommon.BytesToHash(result.Hash))
				return newMsg
			},
			genWatchTx: func(tx tm.TxEssentials) (WatchTx, sdk.Tx) {
				var expectWatchTx WatchTx
				realTx, ok := tx.(sdk.Tx)
				suite.Require().True(ok, "evmTx generate WatchTx error")
				//Create the same Watcher as the tested function
				watchTx := suite.Watcher.createWatchTx(realTx)
				expectTx, ok := watchTx.(*evmTx)
				expectTx.index = expectTx.index - 1
				expectWatchTx = expectTx
				return expectWatchTx, realTx
			},
			numBatch: 3,
		},
		//{
		//	title:          "StdTx success with none-nil ResponseDeliverTx",
		//	watcherEnabled: true,
		//	buildInput: func() (tm.TxEssentials, *tm.ResponseDeliverTx) {
		//		//send std tx for gov, success
		//		content := gov.NewTextProposal("Test", "description")
		//		newProposalMsg := gov.NewMsgSubmitProposal(content, sysCoins10, suite.stdSenderAccAddress)
		//		depositMsg := gov.NewMsgDeposit(suite.stdSenderAccAddress, govProposalID1, sysCoins90)
		//		msgs := []sdk.Msg{newProposalMsg, depositMsg}
		//		tx := newTestStdTx(msgs, []crypto.PrivKey{suite.stdSenderPrivKey}, []uint64{accountNum}, []uint64{nonce0}, txFees, memo)
		//		resp := suite.app.DeliverRealTx(tx)
		//		return tx, &resp
		//	},
		//	genStdTxResponse: func(tx tm.TxEssentials, resp *tm.ResponseDeliverTx) *MsgStdTransactionResponse {
		//		result := &ctypes.ResultTx{
		//			Hash:     tx.TxHash(),
		//			Height:   int64(suite.Watcher.height),
		//			TxResult: *resp,
		//			Tx:       tx.GetRaw(),
		//		}
		//		newMsg := NewStdTransactionResponse(result, suite.Watcher.header.Time, ethcommon.BytesToHash(result.Hash))
		//		return newMsg
		//	},
		//	genWatchTx: func(tx tm.TxEssentials) (WatchTx, sdk.Tx) {
		//		stdTx, ok := tx.(sdk.Tx)
		//		suite.Require().True(ok, "StdTx generate WatchTx error")
		//		watchTx := suite.Watcher.createWatchTx(stdTx)
		//		return watchTx, stdTx
		//	},
		//	numBatch: 1,
		//},
		{
			title:          "Watcher Disabled",
			watcherEnabled: false,
			buildInput: func() (tm.TxEssentials, *tm.ResponseDeliverTx) {
				return nil, nil
			},
			genStdTxResponse: func(tx tm.TxEssentials, resp *tm.ResponseDeliverTx) *MsgStdTransactionResponse {
				return nil
			},
			genWatchTx: func(tx tm.TxEssentials) (WatchTx, sdk.Tx) {
				return nil, nil
			},
			numBatch: 0,
		},
	}

	suite.SetupTest()
	suite.beginFakeBlock()
	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			suite.Watcher.Enable(tc.watcherEnabled)
			tx, resp := tc.buildInput()
			//there is a mutex for batch, direct copy is inappropriate
			//suite.watcherBatch = suite.Watcher.batch
			oldlenbatch := len(suite.Watcher.batch)
			oldWatcher := suite.Watcher
			suite.Watcher.RecordTxAndFailedReceipt(tx, resp, suite.TxDecoder)
			newWatcher := suite.Watcher
			if !tc.watcherEnabled {
				suite.Require().Equal(oldWatcher, newWatcher, "Watcher Disabled error")
				return
			}
			//suite.watcherBatch = suite.Watcher.batch
			//suite.watcherBlockTxs = suite.Watcher.blockTxs
			//suite.watcherBlockStdTxs = suite.Watcher.blockStdTxs
			newlenbatch := len(suite.Watcher.batch)
			suite.Require().True(newlenbatch-oldlenbatch == tc.numBatch, "Add Batch Number Wrong")
			watchTx, realTx := tc.genWatchTx(tx)
			switch realTx.GetType() {
			case sdk.EvmTxType:
				if resp != nil && watchTx != nil {
					// Test record Result Tx
					respMsg := suite.Watcher.batch[oldlenbatch]
					expectMsg := tc.genStdTxResponse(tx, resp)
					suite.Require().Equal(expectMsg, respMsg, "Save txResult error")

					suite.testSaveTxWithArg(watchTx, oldlenbatch+1)
					//Test Save Receipt
					if resp.IsOK() && !suite.Watcher.IsRealEvmTx(resp) {
						msgs := realTx.GetMsgs()
						evmTx, ok := msgs[0].(*etypes.MsgEthereumTx)
						suite.Require().True(ok, "Eth tx get MsgEthereumTx error")
						suite.testSaveTransactionReceiptWithArg(TransactionSuccess, evmTx, watchTx.GetTxHash(), watchTx.GetIndex(), &etypes.ResultData{}, uint64(resp.GasUsed), oldlenbatch+2)
					} else if !resp.IsOK() {
						suite.testSaveFailedReceiptWithArg(watchTx, uint64(resp.GasUsed), oldlenbatch+2)
					}
				} else if resp == nil && watchTx != nil {
					//Only watchTx is saved
					suite.testSaveTxWithArg(watchTx, oldlenbatch)
				} else if resp != nil && watchTx == nil {
					//Only Result Tx is recorded
					respMsg := suite.Watcher.batch[oldlenbatch]
					expectMsg := tc.genStdTxResponse(tx, resp)
					suite.Require().Equal(expectMsg, respMsg, "Save txResult error")
				}
			case sdk.StdTxType:
				if resp != nil {
					expectMsgStdResponse := tc.genStdTxResponse(tx, resp)
					respMsgStdResponse := suite.Watcher.batch[len(suite.Watcher.batch)-1]
					suite.Require().Equal(expectMsgStdResponse, respMsgStdResponse, "StdTx save ResultTx error")
				}
				expectTxHash := ethcommon.BytesToHash(realTx.TxHash())
				respTxHash := suite.Watcher.blockStdTxs[len(suite.Watcher.blockStdTxs)-1]
				suite.Require().Equal(expectTxHash, respTxHash, "StdTx save blockStdTxs error")
			}
		})
	}
	suite.endFakeBlock()
}

func newTestStdTx(msgs []sdk.Msg, privs []crypto.PrivKey, accNums []uint64, seqs []uint64, fee auth.StdFee, memo string) sdk.Tx {
	sigs := make([]authtypes.StdSignature, len(privs))
	for i, priv := range privs {
		sig, err := priv.Sign(authtypes.StdSignBytes(cosmosChainId, accNums[i], seqs[i], fee, msgs, memo))
		if err != nil {
			panic(err)
		}
		sigs[i] = authtypes.StdSignature{PubKey: priv.PubKey(), Signature: sig}
	}

	tx := auth.NewStdTx(msgs, fee, sigs, memo)
	return tx
}
