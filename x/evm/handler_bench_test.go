package evm_test

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/x/evm"
	"github.com/okex/exchain/x/evm/keeper"
	"github.com/okex/exchain/x/evm/types"
	govtypes "github.com/okex/exchain/x/gov/types"
	"math/big"
	"testing"
	"time"
)

var (
	//counterCode       = "0x60806040526001600055600060015534801561001a57600080fd5b506102748061002a6000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c80631ae34bc81461005c57806343114db81461009c578063434aefd3146100dc578063a68fd6ac1461011c578063b8b085f21461013a575b600080fd5b61009a6004803603602081101561007257600080fd5b8101908080356fffffffffffffffffffffffffffffffff169060200190929190505050610158565b005b6100da600480360360208110156100b257600080fd5b8101908080356fffffffffffffffffffffffffffffffff16906020019092919050505061018f565b005b61011a600480360360208110156100f257600080fd5b8101908080356fffffffffffffffffffffffffffffffff1690602001909291905050506101d9565b005b610124610233565b6040518082815260200191505060405180910390f35b610142610239565b6040518082815260200191505060405180910390f35b600080600090505b826fffffffffffffffffffffffffffffffff1681101561018a576000549150806001019050610160565b505050565b600080600090505b826fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff1610156101d457600182019150806001019050610197565b505050565b60008090505b816fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff16101561022f57806fffffffffffffffffffffffffffffffff166001819055508060010190506101df565b5050565b60015481565b6000548156fea265627a7a723158206f5270b1d04421e6f7e326a2dffdf5dc66f5d427524fb366bf5a2fa2663825e764736f6c63430005110032"
	//counterABI        = `[{"constant":false,"inputs":[{"internalType":"uint128","name":"times","type":"uint128"}],"name":"add","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint128","name":"times","type":"uint128"}],"name":"read","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"readCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"uint128","name":"times","type":"uint128"}],"name":"write","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"writeCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`

	//counterCode       = "0x608060405260016000556000600155600060025534801561001f57600080fd5b506102ed8061002f6000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063434aefd31161005b578063434aefd314610120578063a68fd6ac14610160578063b8b085f21461017e578063d14d42ba1461019c5761007d565b806308615e55146100825780631ae34bc8146100a057806343114db8146100e0575b600080fd5b61008a6101ba565b6040518082815260200191505060405180910390f35b6100de600480360360208110156100b657600080fd5b8101908080356fffffffffffffffffffffffffffffffff1690602001909291905050506101c0565b005b61011e600480360360208110156100f657600080fd5b8101908080356fffffffffffffffffffffffffffffffff1690602001909291905050506101f7565b005b61015e6004803603602081101561013657600080fd5b8101908080356fffffffffffffffffffffffffffffffff169060200190929190505050610248565b005b6101686102a2565b6040518082815260200191505060405180910390f35b6101866102a8565b6040518082815260200191505060405180910390f35b6101a46102ae565b6040518082815260200191505060405180910390f35b60025481565b600080600090505b826fffffffffffffffffffffffffffffffff168110156101f25760005491508060010190506101c8565b505050565b600080600090505b826fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff16101561023c576001820191508060010190506101ff565b50806002819055505050565b60008090505b816fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff16101561029e57806fffffffffffffffffffffffffffffffff1660018190555080600101905061024e565b5050565b60015481565b60005481565b600060025490509056fea265627a7a72315820655d1631108e29b72fe10b4dab90e13570c914dbd7f7fb32b8d35225a20ed39264736f6c63430005110032"
	//counterABI        = `[{"constant":false,"inputs":[{"internalType":"uint128","name":"times","type":"uint128"}],"name":"add","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"addCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getAdd","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"uint128","name":"times","type":"uint128"}],"name":"read","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"readCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"uint128","name":"times","type":"uint128"}],"name":"write","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"writeCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`
	counterCode       = "0x608060405260016000556000600155600060025534801561001f57600080fd5b5061042b8061002f6000396000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c8063434aefd31161005b578063434aefd3146100d8578063a68fd6ac146100f4578063b8b085f214610112578063d14d42ba146101305761007d565b806308615e55146100825780631ae34bc8146100a057806343114db8146100bc575b600080fd5b61008a61014e565b604051610097919061027d565b60405180910390f35b6100ba60048036038101906100b591906102e5565b610154565b005b6100d660048036038101906100d191906102e5565b610191565b005b6100f260048036038101906100ed91906102e5565b6101f1565b005b6100fc61024e565b604051610109919061027d565b60405180910390f35b61011a610254565b604051610127919061027d565b60405180910390f35b61013861025a565b604051610145919061027d565b60405180910390f35b60025481565b600080600090505b826fffffffffffffffffffffffffffffffff1681101561018c5760005491508061018590610341565b905061015c565b505050565b600080600090505b826fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff1610156101e5576001826101d29190610389565b9150806101de906103bd565b9050610199565b50806002819055505050565b60005b816fffffffffffffffffffffffffffffffff16816fffffffffffffffffffffffffffffffff16101561024a57806fffffffffffffffffffffffffffffffff1660018190555080610243906103bd565b90506101f4565b5050565b60015481565b60005481565b6000600254905090565b6000819050919050565b61027781610264565b82525050565b6000602082019050610292600083018461026e565b92915050565b600080fd5b60006fffffffffffffffffffffffffffffffff82169050919050565b6102c28161029d565b81146102cd57600080fd5b50565b6000813590506102df816102b9565b92915050565b6000602082840312156102fb576102fa610298565b5b6000610309848285016102d0565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600061034c82610264565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff820361037e5761037d610312565b5b600182019050919050565b600061039482610264565b915061039f83610264565b92508282019050808211156103b7576103b6610312565b5b92915050565b60006103c88261029d565b91506fffffffffffffffffffffffffffffffff82036103ea576103e9610312565b5b60018201905091905056fea26469706673582212201d31f72251150d467fe1ac46963efc23b17905daa45b992321a79594a560386f64736f6c63430008110033"
	counterABI        = `[{"inputs":[{"internalType":"uint128","name":"times","type":"uint128"}],"name":"add","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"addCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"getAdd","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint128","name":"times","type":"uint128"}],"name":"read","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"readCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint128","name":"times","type":"uint128"}],"name":"write","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"writeCounter","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
	addParam    int64 = 3500000
	readParam   int64 = 1
	writeParam  int64 = 1
)

func GenAddParam(count int64, name string) []byte {
	abis, err := types.NewABI(counterABI)
	if err != nil {
		panic(err)
	}
	re, err := abis.Pack(name, big.NewInt(count))
	if err != nil {
		panic(err)
	}
	return re
}

func GenReadAddParam() []byte {
	abis, err := types.NewABI(counterABI)
	if err != nil {
		panic(err)
	}
	re, err := abis.Pack("getAdd")
	if err != nil {
		panic(err)
	}
	return re
}

func TestGenAddParam(t *testing.T) {
	fmt.Println(hexutil.Encode(GenAddParam(1, "add")))
	fmt.Println(hexutil.Encode(GenAddParam(1, "read")))
	fmt.Println(hexutil.Encode(GenAddParam(1, "write")))
}

func UnpackRetValue(methodName string, data []byte) ([]interface{}, error) {
	abis, err := types.NewABI(counterABI)
	if err != nil {
		panic(err)
	}
	method, ok := abis.Methods[methodName]
	if !ok {
		return nil, fmt.Errorf("method %s is not exist in abi", methodName)
	}
	return method.Outputs.UnpackValues(data)
}

//func (suite *EvmTestSuite) TestBatchCounter() {
//
//	gasLimit := uint64(100000000)
//	gasPrice := big.NewInt(10000)
//
//	priv, err := ethsecp256k1.GenerateKey()
//	suite.Require().NoError(err, "failed to create key")
//
//	bytecode := common.FromHex(counterCode)
//	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
//	tx.Sign(big.NewInt(3), priv.ToECDSA())
//	suite.Require().NoError(err)
//
//	result, err := suite.handler(suite.ctx, tx)
//	suite.Require().NoError(err, "failed to handle eth tx msg")
//
//	resultData, err := types.DecodeResultData(result.Data)
//	suite.Require().NoError(err, "failed to decode result data")
//
//	// store - changeOwner
//	gasLimit = uint64(100000000000)
//	gasPrice = big.NewInt(100)
//	receiver := common.HexToAddress(resultData.ContractAddress.String())
//
//	bytecode = GenAddParam(addParam, "add")
//	tx = types.NewMsgEthereumTx(2, &receiver, big.NewInt(0), gasLimit, gasPrice, bytecode)
//	tx.Sign(big.NewInt(3), priv.ToECDSA())
//	suite.Require().NoError(err)
//
//	st := time.Now()
//	result, err = suite.handler(suite.ctx, tx)
//	end := time.Now()
//
//	suite.Require().NoError(err, "failed to handle eth tx msg")
//
//	resultData, err = types.DecodeResultData(result.Data)
//	suite.Require().NoError(err, "failed to decode result data")
//
//	fmt.Println("the res gas used ", "cost time", end.Sub(st))
//
//	// query - getOwner
//	bytecode = GenReadAddParam()
//	tx = types.NewMsgEthereumTx(2, &receiver, big.NewInt(0), gasLimit, gasPrice, bytecode)
//	tx.Sign(big.NewInt(3), priv.ToECDSA())
//	suite.Require().NoError(err)
//
//	result, err = suite.handler(suite.ctx, tx)
//	suite.Require().NoError(err, "failed to handle eth tx msg")
//
//	resultData, err = types.DecodeResultData(result.Data)
//	suite.Require().NoError(err, "failed to decode result data")
//	r, err := UnpackRetValue("getAdd", resultData.Ret)
//	if len(r) == 1 {
//		if rr, ok := r[0].(*big.Int); ok {
//			fmt.Println("the res is ", rr.Int64())
//		}
//	}
//}

type EvmBatchTest struct {
	ctx        sdk.Context
	handler    sdk.Handler
	govHandler govtypes.Handler
	querier    sdk.Querier
	app        *app.OKExChainApp
	stateDB    *types.CommitStateDB
	codec      *codec.Codec

	//
	receiver common.Address
}

func InitEvmBatchTest() *EvmBatchTest {
	ebt := &EvmBatchTest{}
	checkTx := false
	chain_id := "ethermint-3"

	ebt.app = app.Setup(checkTx)
	ebt.ctx = ebt.app.BaseApp.NewContext(checkTx, abci.Header{Height: 1, ChainID: chain_id, Time: time.Now().UTC()})
	ebt.ctx.SetDeliver()
	ebt.stateDB = types.CreateEmptyCommitStateDB(ebt.app.EvmKeeper.GenerateCSDBParams(), ebt.ctx)
	ebt.handler = evm.NewHandler(ebt.app.EvmKeeper)
	ebt.querier = keeper.NewQuerier(*ebt.app.EvmKeeper)
	ebt.codec = codec.New()

	err := ethermint.SetChainId(chain_id)
	if err != nil {
		panic(err)
	}

	params := types.DefaultParams()
	params.EnableCreate = true
	params.EnableCall = true
	ebt.app.EvmKeeper.SetParams(ebt.ctx, params)
	return ebt
}

func (e *EvmBatchTest) Deploy() {
	gasLimit := uint64(100000000)
	gasPrice := big.NewInt(10000)

	priv, err := ethsecp256k1.GenerateKey()
	if err != nil {
		panic(err)
	}

	bytecode := common.FromHex(counterCode)
	tx := types.NewMsgEthereumTx(1, nil, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())

	result, err := e.handler(e.ctx, tx)
	if err != nil {
		panic(err)
	}

	resultData, err := types.DecodeResultData(result.Data)
	if err != nil {
		panic(err)
	}
	e.receiver = common.HexToAddress(resultData.ContractAddress.String())
}

func (e *EvmBatchTest) Run(name string) {
	gasLimit := uint64(100000000000000)
	gasPrice := big.NewInt(10000)
	bytecode := GenAddParam(addParam, name)
	priv, err := ethsecp256k1.GenerateKey()
	tx := types.NewMsgEthereumTx(2, &e.receiver, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())

	st := time.Now()
	result, err := e.handler(e.ctx, tx)
	if err != nil {
		panic(err)
	}
	end := time.Now()

	resultData, err := types.DecodeResultData(result.Data)
	if err != nil {
		panic(err)
	}
	fmt.Println("opt", name, "count", addParam, "cost time", end.Sub(st), len(resultData.Ret))
}

func (e *EvmBatchTest) Query() {
	gasLimit := uint64(100000000)
	gasPrice := big.NewInt(10000)
	priv, err := ethsecp256k1.GenerateKey()
	bytecode := GenReadAddParam()
	tx := types.NewMsgEthereumTx(2, &e.receiver, big.NewInt(0), gasLimit, gasPrice, bytecode)
	tx.Sign(big.NewInt(3), priv.ToECDSA())

	result, err := e.handler(e.ctx, tx)
	if err != nil {
		panic(err)
	}

	resultData, err := types.DecodeResultData(result.Data)
	if err != nil {
		panic(err)
	}
	r, err := UnpackRetValue("getAdd", resultData.Ret)
	if len(r) == 1 {
		if rr, ok := r[0].(*big.Int); ok {
			fmt.Println("the res is ", rr.Int64())
		}
	}
}

func TestEvmBatchAdd(t *testing.T) {
	e := InitEvmBatchTest()
	e.Deploy()
	e.Run("add")
	e.Query()
}

func TestEvmBatchRead(t *testing.T) {
	e := InitEvmBatchTest()
	e.Deploy()
	e.Run("read")
	e.Query()
}

func TestEvmBatchWrite(t *testing.T) {
	e := InitEvmBatchTest()
	e.Deploy()
	e.Run("write")
	e.Query()
}

func TestEvmBatch(t *testing.T) {
	counts := []int64{1000000, 1250000, 1500000, 1750000, 2000000, 2250000, 2500000, 2750000, 3000000, 3250000, 3500000, 3750000, 4000000, 4250000, 4500000, 4750000, 5000000}
	e := InitEvmBatchTest()
	e.Deploy()
	for _, c := range counts {
		addParam = c
		e.Run("add")
		e.Run("read")
		e.Run("write")
	}
}
