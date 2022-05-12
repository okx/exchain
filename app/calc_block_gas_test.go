package app

import (
	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	sm "github.com/okex/exchain/libs/tendermint/state"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"math/big"
	"testing"
	"time"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authclient "github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/secp256k1"
	tender_types "github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/x/evm/types"
	evm_types "github.com/okex/exchain/x/evm/types"
	"github.com/stretchr/testify/suite"
)

var (
	tx_coin10      = sdk.NewInt64Coin(sdk.DefaultBondDenom, 10)
	tx_coin90      = sdk.NewInt64Coin(sdk.DefaultBondDenom, 90)
	tx_coin100     = sdk.NewInt64Coin(sdk.DefaultBondDenom, 100)
	tx_fees        = auth.NewStdFee(21000, sdk.NewCoins(tx_coin10))
	expectGas      = [...]int64{60602, 61938, 63274, 63274}
	expectBlockGas = sdk.Gas(63274)
)

type BlockTxTestSuite struct {
	suite.Suite
	ctx     sdk.Context
	app     *OKExChainApp
	stateDB *evm_types.CommitStateDB
	codec   *codec.Codec
	block   *tender_types.Block
	txs     []tender_types.Tx

	contractDeloyerPrivKey ethsecp256k1.PrivKey
	contractAddress        ethcmn.Address
}

func (suite *BlockTxTestSuite) SetupTest() {
	checkTx := false
	chain_id := "ethermint-3"

	suite.app = Setup(checkTx)
	suite.ctx = suite.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: chain_id, Time: time.Now().UTC()})
	suite.ctx.SetDeliver()
	suite.stateDB = evm_types.CreateEmptyCommitStateDB(suite.app.EvmKeeper.GenerateCSDBParams(), suite.ctx)
	suite.codec = codec.New()

	err := ethermint.SetChainId(chain_id)
	suite.Nil(err)

	params := evm_types.DefaultParams()
	params.EnableCreate = true
	params.EnableCall = true
	suite.app.EvmKeeper.SetParams(suite.ctx, params)
}

func TestBlcokInnerTxTestSuite(t *testing.T) {
	suite.Run(t, new(BlockTxTestSuite))
}

func (suite *BlockTxTestSuite) makeBlock(height int64, state sm.State, lastCommit *tender_types.Commit) {
	//suite.block, _ = state.MakeBlock(height, suite.makeTxs(height), lastCommit, nil, nil)
}

func (suite *BlockTxTestSuite) makeTxs(height int64) (txs []tender_types.Tx) {
	txs = append(txs, suite.buildTxEvmDeploy(height))
	txs = append(txs, suite.buildTxEvmDeployError(height))
	txs = append(txs, suite.buildTxEvmCallStore(height))
	txs = append(txs, suite.buildTxEvmCallStoreError(height))
	txs = append(txs, suite.buildTxEvmCallQuery(height))
	txs = append(txs, suite.buildTxEvmCallQueryError(height))
	txs = append(txs, suite.buildStdTxSendMsgBank(height))
	txs = append(txs, suite.buildStdTxSendMsgBankError(height))
	return txs
}

func (suite *BlockTxTestSuite) TestDeployAndCallContract() {

	startBlock := func() {
		newHeader := suite.ctx.BlockHeader()
		newHeader.Time = suite.ctx.BlockHeader().Time.Add(time.Duration(1) * time.Second)
		suite.ctx.SetBlockHeader(newHeader)
		suite.app.BeginBlocker(suite.ctx, abci.RequestBeginBlock{Header: abci.Header{Height: 1}})
		suite.makeTxs(12)
	}

	deliverTxs := func() {
		for i := 0; i < len(suite.txs); i++ {
			txReal := suite.app.PreDeliverRealTx(suite.txs[i])
			resp := suite.app.DeliverRealTx(txReal)
			actualGas := resp.GasUsed
			suite.Require().True(expectGas[i] == actualGas, "expect gas %d, not equal actual gas %d ", expectGas[i], actualGas)
		}
	}

	endBlock := func() {
		suite.app.Commit(abci.RequestCommit{})
		suite.app.EndBlocker(suite.ctx, abci.RequestEndBlock{})
		blockActualGas := suite.ctx.BlockGasMeter().GasConsumed()
		suite.Require().True(expectBlockGas == blockActualGas, "expect gas %d, not equal actual gas %d ", expectBlockGas, blockActualGas)
	}

	testCases := []struct {
		msg     string
		prepare func()
	}{
		{
			"process tx",
			func() {
				//Process default
			},
		},
		{
			"parallel tx",
			func() {
				//Parallel setting, TODO
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.msg, func() {
			suite.SetupTest() // reset
			startBlock()
			tc.prepare()
			deliverTxs()
			endBlock()
		})
	}
}

func (suite *BlockTxTestSuite) buildTxEvmDeploy(height int64) []byte {
	//Create evm contract - Owner.sol
	gasLimit := uint64(100000000)
	gasPrice := big.NewInt(10000)

	priv, err := ethsecp256k1.GenerateKey()
	suite.Require().NoError(err, "failed to create private key")

	sender := ethcmn.HexToAddress(priv.PubKey().Address().String())
	suite.app.EvmKeeper.SetBalance(suite.ctx, sender, big.NewInt(100))

	bytecode := common.FromHex("0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a36102c4806100dc6000396000f3fe608060405234801561001057600080fd5b5060043610610053576000357c010000000000000000000000000000000000000000000000000000000090048063893d20e814610058578063a6f9dae1146100a2575b600080fd5b6100606100e6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100e4600480360360208110156100b857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061010f565b005b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101d1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f43616c6c6572206973206e6f74206f776e65720000000000000000000000000081525060200191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505056fea265627a7a72315820f397f2733a89198bc7fed0764083694c5b828791f39ebcbc9e414bccef14b48064736f6c63430005100032")
	tx := types.NewMsgEthereumTx(1, &sender, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())

	var txEncoder sdk.TxEncoder
	if tmtypes.HigherThanVenus(int64(height)) {
		txEncoder = authclient.GetTxEncoder(nil, authclient.WithEthereumTx())
	} else {
		txEncoder = authclient.GetTxEncoder(suite.codec)
	}
	// Encode transaction by RLP encoder
	txBytes, err := txEncoder(tx)

	return txBytes
}

func (suite *BlockTxTestSuite) buildTxEvmDeployError(height int64) []byte {
	//TODO
	return nil
}

func (suite *BlockTxTestSuite) buildTxEvmCallStore(height int64) []byte {
	// Execute evm contract with function changeOwner, for saving data.
	gasLimit := uint64(100000000000)
	gasPrice := big.NewInt(100)

	storeAddr := "0xa6f9dae10000000000000000000000006a82e4a67715c8412a9114fbd2cbaefbc8181424"
	bytecode := common.FromHex(storeAddr)
	tx := types.NewMsgEthereumTx(2, &suite.contractAddress, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), suite.contractDeloyerPrivKey.ToECDSA())

	var txEncoder sdk.TxEncoder
	if tmtypes.HigherThanVenus(int64(height)) {
		txEncoder = authclient.GetTxEncoder(nil, authclient.WithEthereumTx())
	} else {
		txEncoder = authclient.GetTxEncoder(suite.codec)
	}
	// Encode transaction by RLP encoder
	txBytes, _ := txEncoder(tx)

	return txBytes

	//suite.Require().NoError(err, "failed to handle eth tx msg")
	//_, err = types.DecodeResultData(result.Data)
	//suite.Require().NoError(err, "failed to decode result data")
	//actualGas := suite.ctx.GasMeter().GasConsumed()
	//suite.Require().True(expectEVMCallSaveGas == actualGas, "expect gas %d, not equal actual gas %d ", expectEVMCallSaveGas, actualGas)
}

func (suite *BlockTxTestSuite) buildTxEvmCallStoreError(height int64) []byte {
	return nil
}

func (suite *BlockTxTestSuite) buildTxEvmCallQuery(height int64) []byte {
	//Execute evm contract with function getOwner, for querying data.
	bytecode := common.FromHex("0x893d20e8")
	gasLimit := uint64(100000000000)
	gasPrice := big.NewInt(100)
	tx := types.NewMsgEthereumTx(2, &suite.contractAddress, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), suite.contractDeloyerPrivKey.ToECDSA())
	var txEncoder sdk.TxEncoder
	if tmtypes.HigherThanVenus(int64(height)) {
		txEncoder = authclient.GetTxEncoder(nil, authclient.WithEthereumTx())
	} else {
		txEncoder = authclient.GetTxEncoder(suite.codec)
	}
	// Encode transaction by RLP encoder
	txBytes, _ := txEncoder(tx)
	return txBytes
	//suite.Require().NoError(err, "failed to handle eth tx msg")
	//resultData, _ := types.DecodeResultData(result.Data)
	//suite.Require().NoError(err, "failed to decode result data")
	//
	//storeAddr := "0xa6f9dae10000000000000000000000006a82e4a67715c8412a9114fbd2cbaefbc8181424"
	//bytecode = common.FromHex(storeAddr)
	//getAddr := strings.ToLower(hexutils.BytesToHex(resultData.Ret))
	//suite.Require().True(strings.HasSuffix(storeAddr, getAddr), "Fail to query the address")
	//actualGas := suite.ctx.GasMeter().GasConsumed()
	//suite.Require().True(expectEVMCallQueryGas == actualGas, "expect gas %d, not equal actual gas %d ", expectEVMCallQueryGas, actualGas)
}

func (suite *BlockTxTestSuite) buildTxEvmCallQueryError(height int64) []byte {
	//TODO
	return nil
}

func (suite *BlockTxTestSuite) buildStdTxSendMsgBank(height int64) []byte {
	var (
		tx          sdk.Tx
		privFrom, _ = ethsecp256k1.GenerateKey()
		cmFrom      = sdk.AccAddress(privFrom.PubKey().Address())
		privTo      = secp256k1.GenPrivKeySecp256k1([]byte("private key to"))
		cmTo        = sdk.AccAddress(privTo.PubKey().Address())
	)
	normal := func() {
		err := suite.app.BankKeeper.SetCoins(suite.ctx, cmFrom, sdk.NewCoins(tx_coin100))
		suite.Require().NoError(err)
	}

	suite.SetupTest() // reset
	normal()

	msg := bank.NewMsgSend(cmFrom, cmTo, sdk.NewCoins(tx_coin10))
	tx = auth.NewStdTx([]sdk.Msg{msg}, tx_fees, nil, "")

	var txEncoder sdk.TxEncoder
	if tmtypes.HigherThanVenus(int64(height)) {
		txEncoder = authclient.GetTxEncoder(nil, authclient.WithEthereumTx())
	} else {
		txEncoder = authclient.GetTxEncoder(suite.codec)
	}
	// Encode transaction by RLP encoder
	txBytes, _ := txEncoder(tx)
	return txBytes

	//suite.ctx.SetGasMeter(sdk.NewInfiniteGasMeter())
	//
	//msgs := tx.GetMsgs()
	//for _, msg := range msgs {
	//	suite.Require().NoError(err)
	//}
	//
	//fromBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmFrom).GetCoins()
	//suite.Require().True(fromBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(tx_coin90))))
	//
	//toBalance := suite.app.AccountKeeper.GetAccount(suite.ctx, cmTo).GetCoins()
	//suite.Require().True(toBalance.IsEqual(sdk.NewDecCoins(sdk.NewDecCoinFromCoin(tx_coin10))))
	//
	//actualGas := suite.ctx.GasMeter().GasConsumed() //TODO is it cumulative?
	//suite.Require().True(expectStdTxGas == actualGas, "expect gas %d, not equal actual gas %d ", expectStdTxGas, actualGas)
}

func (suite *BlockTxTestSuite) buildStdTxSendMsgBankError(height int64) []byte {
	//TODO
	return nil
}
