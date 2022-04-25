//This is a test utility for Ethermint's Web3 JSON-RPC services.
//
// To run these tests please first ensure you have the ethermintd running
// and have started the RPC service with `ethermintcli rest-server`.
//
// You can configure the desired HOST and MODE as well
package pending

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"testing"
	"time"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gorpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/rpc"
	"github.com/okex/exchain/app/rpc/backend"
	util "github.com/okex/exchain/app/rpc/tests"
	cosmos_context "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	cmserver "github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	apptesting "github.com/okex/exchain/libs/ibc-go/testing"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
	"github.com/okex/exchain/libs/tendermint/crypto/multisig"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/suite"
)

const (
	addrAStoreKey      = 0
	defaultMinGasPrice = "0.000000001okt"
	latestBlockNumber  = "latest"
	pendingBlockNumber = "pending"
)

var (
	receiverAddr    = ethcmn.BytesToAddress([]byte("receiver"))
	defaultGasPrice sdk.SysCoin
)

func init() {
	defaultGasPrice, _ = sdk.ParseDecCoin(defaultMinGasPrice)
}

type Request struct {
	Version string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	ID      int         `json:"id"`
}

type RPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type Response struct {
	Error  *RPCError       `json:"error"`
	ID     int             `json:"id"`
	Result json.RawMessage `json:"result,omitempty"`
}
type RPCPendingTestSuite struct {
	suite.Suite

	coordinator *apptesting.Coordinator

	// testing chains used for convenience and readability
	chain apptesting.TestChainI

	apiServer *gorpc.Server
	Mux       *http.ServeMux
	cliCtx    *cosmos_context.CLIContext
	addr      string
}

func init() {
	tmamino.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)
	multisig.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)
}

var (
	genesisAcc sdk.AccAddress
	senderAddr ethcmn.Address
)

func commitBlock(suite *RPCPendingTestSuite) {
	mck, ok := suite.cliCtx.Client.(*util.MockClient)
	suite.Require().True(ok)
	mck.CommitBlock()
}
func (suite *RPCPendingTestSuite) SetupTest() {
	chainId := apptesting.GetOKChainID(1)
	suite.coordinator = apptesting.NewEthCoordinator(suite.T(), 1)
	suite.chain = suite.coordinator.GetChain(chainId)
	suite.chain.App().SetOption(abci.RequestSetOption{
		Key:   "CheckChainID",
		Value: chainId,
	})

	cliCtx := cosmos_context.NewCLIContext().
		WithProxy(suite.chain.Codec()).
		WithTrustNode(true).
		WithChainID(chainId).
		WithClient(util.NewMockClient(chainId, suite.chain, suite.chain.App())).
		WithBroadcastMode(flags.BroadcastSync)

	suite.cliCtx = &cliCtx

	commitBlock(suite)

	suite.apiServer = gorpc.NewServer()

	viper.Set(rpc.FlagDisableAPI, "")
	viper.Set(backend.FlagApiBackendBlockLruCache, 100)
	viper.Set(backend.FlagApiBackendTxLruCache, 100)
	viper.Set(watcher.FlagFastQueryLru, 100)
	viper.Set("rpc.laddr", "127.0.0.1:0")
	viper.Set(flags.FlagKeyringBackend, "test")
	viper.Set(cmserver.FlagUlockKeyHome, fmt.Sprintf(".exchaincli/%s", time.Now().String()))

	viper.Set(rpc.FlagPersonalAPI, true)

	senderPv := suite.chain.SenderAccountPV()
	genesisAcc = suite.chain.SenderAccount().GetAddress()
	senderAddr = ethcmn.BytesToAddress(genesisAcc.Bytes())
	apis := rpc.GetAPIs(cliCtx, log.NewNopLogger(), []ethsecp256k1.PrivKey{ethsecp256k1.PrivKey(senderPv.Bytes())}...)
	for _, api := range apis {
		if err := suite.apiServer.RegisterName(api.Namespace, api.Service); err != nil {
			panic(err)
		}
	}
	suite.Mux = http.NewServeMux()
	suite.Mux.HandleFunc("/", suite.apiServer.ServeHTTP)
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	suite.addr = fmt.Sprintf("http://localhost:%d", listener.Addr().(*net.TCPAddr).Port)
	go func() {
		http.Serve(listener, suite.Mux)
	}()
}

func (suite *RPCPendingTestSuite) TestEth_Pending_GetBalance() {
	//waitForBlock(5)

	var resLatest, resPending hexutil.Big
	rpcResLatest := util.Call(suite.T(), suite.addr, "eth_getBalance", []interface{}{receiverAddr, latestBlockNumber})
	rpcResPending := util.Call(suite.T(), suite.addr, "eth_getBalance", []interface{}{receiverAddr, pendingBlockNumber})

	suite.Require().NoError(resLatest.UnmarshalJSON(rpcResLatest.Result))
	suite.Require().NoError(resPending.UnmarshalJSON(rpcResPending.Result))
	preTxLatestBalance := resLatest.ToInt()
	preTxPendingBalance := resPending.ToInt()
	//t.Logf("Got pending balance %s for %s pre tx\n", preTxPendingBalance, receiverAddr)
	//t.Logf("Got latest balance %s for %s pre tx\n", preTxLatestBalance, receiverAddr)

	// transfer
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = senderAddr.Hex()
	param[0]["to"] = receiverAddr.Hex()
	param[0]["value"] = "0xA"
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()

	_ = util.Call(suite.T(), suite.addr, "eth_sendTransaction", param)

	rpcResLatest = util.Call(suite.T(), suite.addr, "eth_getBalance", []interface{}{receiverAddr, latestBlockNumber})
	rpcResPending = util.Call(suite.T(), suite.addr, "eth_getBalance", []interface{}{receiverAddr, pendingBlockNumber})

	suite.Require().NoError(resPending.UnmarshalJSON(rpcResPending.Result))
	suite.Require().NoError(resLatest.UnmarshalJSON(rpcResLatest.Result))

	postTxPendingBalance := resPending.ToInt()
	postTxLatestBalance := resLatest.ToInt()
	//t.Logf("Got pending balance %s for %s post tx\n", postTxPendingBalance, receiverAddr)
	//t.Logf("Got latest balance %s for %s post tx\n", postTxLatestBalance, receiverAddr)

	suite.Require().Equal(preTxPendingBalance.Add(preTxPendingBalance, big.NewInt(10)), postTxPendingBalance)
	// preTxLatestBalance <= postTxLatestBalance
	suite.Require().True(preTxLatestBalance.Cmp(postTxLatestBalance) <= 0)
}

func (suite *RPCPendingTestSuite) TestEth_Pending_GetTransactionCount() {
	//waitForBlock(5)

	prePendingNonce := util.GetNonce(suite.T(), suite.addr, pendingBlockNumber, senderAddr.Hex())
	currentNonce := util.GetNonce(suite.T(), suite.addr, latestBlockNumber, senderAddr.Hex())
	//t.Logf("Pending nonce before tx is %d", prePendingNonce)
	//t.Logf("Current nonce is %d", currentNonce)

	suite.Require().True(prePendingNonce == currentNonce)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = senderAddr.Hex()
	param[0]["to"] = receiverAddr.Hex()
	param[0]["value"] = "0xA"
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()

	_ = util.Call(suite.T(), suite.addr, "eth_sendTransaction", param)

	pendingNonce := util.GetNonce(suite.T(), suite.addr, pendingBlockNumber, senderAddr.Hex())
	latestNonce := util.GetNonce(suite.T(), suite.addr, latestBlockNumber, senderAddr.Hex())
	//t.Logf("Latest nonce is %d", latestNonce)
	//t.Logf("Pending nonce is %d", pendingNonce)

	suite.Require().True(currentNonce <= latestNonce)
	suite.Require().True(latestNonce <= pendingNonce)
	suite.Require().True(prePendingNonce+1 == pendingNonce)
}

func (suite *RPCPendingTestSuite) TestEth_Pending_GetBlockTransactionCountByNumber() {
	//waitForBlock(5)

	rpcResLatest := util.Call(suite.T(), suite.addr, "eth_getBlockTransactionCountByNumber", []interface{}{latestBlockNumber})
	rpcResPending := util.Call(suite.T(), suite.addr, "eth_getBlockTransactionCountByNumber", []interface{}{pendingBlockNumber})
	var preTxPendingTxCount, preTxLatestTxCount hexutil.Uint
	suite.Require().NoError(json.Unmarshal(rpcResPending.Result, &preTxPendingTxCount))
	suite.Require().NoError(json.Unmarshal(rpcResLatest.Result, &preTxLatestTxCount))
	//t.Logf("Pre tx pending nonce is %d", preTxPendingTxCount)
	//t.Logf("Pre tx latest nonce is %d", preTxLatestTxCount)
	suite.Require().True(preTxPendingTxCount == preTxLatestTxCount)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = senderAddr.Hex()
	param[0]["to"] = receiverAddr.Hex()
	param[0]["value"] = "0xA"
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()

	txRes := util.Call(suite.T(), suite.addr, "eth_sendTransaction", param)
	suite.Require().Nil(txRes.Error)

	rpcResLatest = util.Call(suite.T(), suite.addr, "eth_getBlockTransactionCountByNumber", []interface{}{latestBlockNumber})
	rpcResPending = util.Call(suite.T(), suite.addr, "eth_getBlockTransactionCountByNumber", []interface{}{pendingBlockNumber})
	var postTxPendingTxCount, postTxLatestTxCount hexutil.Uint
	suite.Require().NoError(json.Unmarshal(rpcResPending.Result, &postTxPendingTxCount))
	suite.Require().NoError(json.Unmarshal(rpcResLatest.Result, &postTxLatestTxCount))
	//t.Logf("Post tx pending nonce is %d", postTxPendingTxCount)
	//t.Logf("Post tx latest nonce is %d", postTxLatestTxCount)

	suite.Require().True(preTxPendingTxCount+1 == postTxPendingTxCount)
	suite.Require().True((postTxPendingTxCount - preTxPendingTxCount) >= (postTxLatestTxCount - preTxLatestTxCount))
}

func (suite *RPCPendingTestSuite) TestEth_Pending_GetBlockByNumber() {
	//waitForBlock(5)

	rpcResLatest := util.Call(suite.T(), suite.addr, "eth_getBlockByNumber", []interface{}{latestBlockNumber, true})
	rpcResPending := util.Call(suite.T(), suite.addr, "eth_getBlockByNumber", []interface{}{pendingBlockNumber, true})

	var preTxLatestBlock, preTxPendingBlock map[string]interface{}
	suite.Require().NoError(json.Unmarshal(rpcResLatest.Result, &preTxLatestBlock))
	suite.Require().NoError(json.Unmarshal(rpcResPending.Result, &preTxPendingBlock))
	//preTxLatestTxs := len(preTxLatestBlock["transactions"].([]interface{}))
	//preTxPendingTxs := len(preTxPendingBlock["transactions"].([]interface{}))
	//t.Logf("Pre tx latest block tx number: %d\n", preTxLatestTxs)
	//t.Logf("Pre tx pending block tx number: %d\n", preTxPendingTxs)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = senderAddr.Hex()
	param[0]["to"] = receiverAddr.Hex()
	param[0]["value"] = "0xA"
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()

	_ = util.Call(suite.T(), suite.addr, "eth_sendTransaction", param)

	rpcResLatest = util.Call(suite.T(), suite.addr, "eth_getBlockByNumber", []interface{}{latestBlockNumber, true})
	rpcResPending = util.Call(suite.T(), suite.addr, "eth_getBlockByNumber", []interface{}{pendingBlockNumber, true})

	var postTxPendingBlock, postTxLatestBlock map[string]interface{}
	suite.Require().NoError(json.Unmarshal(rpcResPending.Result, &postTxPendingBlock))
	suite.Require().NoError(json.Unmarshal(rpcResLatest.Result, &postTxLatestBlock))
	//postTxPendingTxs := len(postTxPendingBlock["transactions"].([]interface{}))
	//postTxLatestTxs := len(postTxLatestBlock["transactions"].([]interface{}))
	//t.Logf("Post tx latest block tx number: %d\n", postTxLatestTxs)
	//t.Logf("Post tx pending block tx number: %d\n", postTxPendingTxs)

	//suite.Require().True(postTxPendingTxs >= preTxPendingTxs)
	//suite.Require().True(preTxLatestTxs == postTxLatestTxs)
}

//func (suite *RPCPendingTestSuite)TestEth_Pending_GetTransactionByBlockNumberAndIndex() {
//	var pendingTx []*rpctypes.Transaction
//	resPendingTxs := util.Call(suite.T(),"eth_pendingTransactions", []string{})
//	err := json.Unmarshal(resPendingTxs.Result, &pendingTx)
//	suite.Require().NoError( err)
//	pendingTxCount := len(pendingTx)
//
//	data := "0x608060405234801561001057600080fd5b5061011e806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063bc9c707d14602d575b600080fd5b603360ab565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101560715780820151818401526020810190506058565b50505050905090810190601f168015609d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040518060400160405280600681526020017f617261736b61000000000000000000000000000000000000000000000000000081525090509056fea2646970667358221220a31fa4c1ce0b3651fbf5401c511b483c43570c7de4735b5c3b0ad0db30d2573164736f6c63430007050033"
//	param := make([]map[string]string, 1)
//	param[0] = make(map[string]string)
//	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
//	param[0]["to"] = addrA
//	param[0]["value"] = "0xA"
//	param[0]["gasLimit"] = "0x5208"
//	param[0]["gasPrice"] = "0x1"
//	param[0]["data"] = data
//
//	txRes := util.Call(suite.T(),"eth_sendTransaction", param)
//	suite.Require().Nil( txRes.Error)
//
//	rpcRes := util.Call(suite.T(),"eth_getTransactionByBlockNumberAndIndex", []interface{}{"pending", "0x" + fmt.Sprintf("%X", pendingTxCount)})
//	var pendingBlockTx map[string]interface{}
//	err = json.Unmarshal(rpcRes.Result, &pendingBlockTx)
//	suite.Require().NoError( err)
//
//	// verify the pending tx has all the correct fields from the tx sent.
//	suite.Require().NotEmpty( pendingBlockTx["hash"])
//	suite.Require().Equal( pendingBlockTx["value"], "0xa")
//	suite.Require().Equal( data, pendingBlockTx["input"])
//
//	rpcRes = util.Call(suite.T(),"eth_getTransactionByBlockNumberAndIndex", []interface{}{"latest", "0x" + fmt.Sprintf("%X", pendingTxCount)})
//	var latestBlock map[string]interface{}
//	err = json.Unmarshal(rpcRes.Result, &latestBlock)
//	suite.Require().NoError( err)
//
//	// verify the pending trasnaction does not exist in the latest block info.
//	suite.Require().Empty( latestBlock)
//}
//
//func (suite *RPCPendingTestSuite)TestEth_Pending_GetTransactionByHash() {
//	data := "0x608060405234801561001057600080fd5b5061011e806100206000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806302eb691b14602d575b600080fd5b603360ab565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101560715780820151818401526020810190506058565b50505050905090810190601f168015609d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60606040518060400160405280600d81526020017f617261736b61776173686572650000000000000000000000000000000000000081525090509056fea264697066735822122060917c5c2fab8c058a17afa6d3c1d23a7883b918ea3c7157131ea5b396e1aa7564736f6c63430007050033"
//	param := make([]map[string]string, 1)
//	param[0] = make(map[string]string)
//	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
//	param[0]["to"] = addrA
//	param[0]["value"] = "0xA"
//	param[0]["gasLimit"] = "0x5208"
//	param[0]["gasPrice"] = "0x1"
//	param[0]["data"] = data
//
//	txRes := util.Call(suite.T(),"eth_sendTransaction", param)
//	var txHash common.Hash
//	err := txHash.UnmarshalJSON(txRes.Result)
//	suite.Require().NoError( err)
//
//	rpcRes := util.Call(suite.T(),"eth_getTransactionByHash", []interface{}{txHash})
//	var pendingBlockTx map[string]interface{}
//	err = json.Unmarshal(rpcRes.Result, &pendingBlockTx)
//	suite.Require().NoError( err)
//
//	// verify the pending tx has all the correct fields from the tx sent.
//	suite.Require().NotEmpty( pendingBlockTx)
//	suite.Require().NotEmpty( pendingBlockTx["hash"])
//	suite.Require().Equal( pendingBlockTx["value"], "0xa")
//	suite.Require().Equal( pendingBlockTx["input"], data)
//}
//
//func (suite *RPCPendingTestSuite)TestEth_Pending_SendTransaction_PendingNonce() {
//	currNonce := util.GetNonce( "latest")
//	param := make([]map[string]string, 1)
//	param[0] = make(map[string]string)
//	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
//	param[0]["to"] = addrA
//	param[0]["value"] = "0xA"
//	param[0]["gasLimit"] = "0x5208"
//	param[0]["gasPrice"] = "0x1"
//
//	// first transaction
//	txRes1 := util.Call(suite.T(),"eth_sendTransaction", param)
//	suite.Require().Nil( txRes1.Error)
//	pendingNonce1 := util.GetNonce( "pending")
//	suite.Require().Greater( uint64(pendingNonce1), uint64(currNonce))
//
//	// second transaction
//	param[0]["to"] = "0x7f0f463c4d57b1bd3e3b79051e6c5ab703e803d9"
//	txRes2 := util.Call(suite.T(),"eth_sendTransaction", param)
//	suite.Require().Nil( txRes2.Error)
//	pendingNonce2 := util.GetNonce( "pending")
//	suite.Require().Greater( uint64(pendingNonce2), uint64(currNonce))
//	suite.Require().Greater( uint64(pendingNonce2), uint64(pendingNonce1))
//
//	// third transaction
//	param[0]["to"] = "0x7fb24493808b3f10527e3e0870afeb8a953052d2"
//	txRes3 := util.Call(suite.T(),"eth_sendTransaction", param)
//	suite.Require().Nil( txRes3.Error)
//	pendingNonce3 := util.GetNonce( "pending")
//	suite.Require().Greater( uint64(pendingNonce3), uint64(currNonce))
//	suite.Require().Greater( uint64(pendingNonce3), uint64(pendingNonce2))
//}

/*func waitForBlock(second int64) {
	fmt.Printf("wait %ds for a clean slate of a new block\n", second)
	time.Sleep(time.Duration(second) * time.Second)
}*/

func TestRPCPendingTestSuite(t *testing.T) {
	suite.Run(t, new(RPCPendingTestSuite))
}
