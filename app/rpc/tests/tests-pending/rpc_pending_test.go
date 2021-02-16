// This is a test utility for Ethermint's Web3 JSON-RPC services.
//
// To run these tests please first ensure you have the ethermintd running
// and have started the RPC service with `ethermintcli rest-server`.
//
// You can configure the desired HOST and MODE as well
package pending

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	util "github.com/okex/okexchain/app/rpc/tests"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"testing"
	"time"
)

const (
	senderAddrHex      = "0x2CF4ea7dF75b513509d95946B43062E26bD88035"
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

func TestMain(m *testing.M) {
	// Start all tests
	code := m.Run()
	os.Exit(code)
}

func TestEth_Pending_GetBalance(t *testing.T) {
	waitForBlock(5)

	var resLatest, resPending hexutil.Big
	rpcResLatest := util.Call(t, "eth_getBalance", []interface{}{receiverAddr, latestBlockNumber})
	rpcResPending := util.Call(t, "eth_getBalance", []interface{}{receiverAddr, pendingBlockNumber})

	require.NoError(t, resLatest.UnmarshalJSON(rpcResLatest.Result))
	require.NoError(t, resPending.UnmarshalJSON(rpcResPending.Result))
	preTxLatestBalance := resLatest.ToInt()
	preTxPendingBalance := resPending.ToInt()
	t.Logf("Got pending balance %s for %s pre tx\n", preTxPendingBalance, receiverAddr)
	t.Logf("Got latest balance %s for %s pre tx\n", preTxLatestBalance, receiverAddr)

	// transfer
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = senderAddrHex
	param[0]["to"] = receiverAddr.Hex()
	param[0]["value"] = "0xA"
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()

	_ = util.Call(t, "eth_sendTransaction", param)

	rpcResLatest = util.Call(t, "eth_getBalance", []interface{}{receiverAddr, latestBlockNumber})
	rpcResPending = util.Call(t, "eth_getBalance", []interface{}{receiverAddr, pendingBlockNumber})

	require.NoError(t, resPending.UnmarshalJSON(rpcResPending.Result))
	require.NoError(t, resLatest.UnmarshalJSON(rpcResLatest.Result))

	postTxPendingBalance := resPending.ToInt()
	postTxLatestBalance := resLatest.ToInt()
	t.Logf("Got pending balance %s for %s post tx\n", postTxPendingBalance, receiverAddr)
	t.Logf("Got latest balance %s for %s post tx\n", postTxLatestBalance, receiverAddr)

	require.Equal(t, preTxPendingBalance.Add(preTxPendingBalance, big.NewInt(10)), postTxPendingBalance)
	// preTxLatestBalance <= postTxLatestBalance
	require.True(t, preTxLatestBalance.Cmp(postTxLatestBalance) <= 0)
}

func TestEth_Pending_GetTransactionCount(t *testing.T) {
	waitForBlock(5)

	prePendingNonce := util.GetNonce(t, pendingBlockNumber, senderAddrHex)
	currentNonce := util.GetNonce(t, latestBlockNumber, senderAddrHex)
	t.Logf("Pending nonce before tx is %d", prePendingNonce)
	t.Logf("Current nonce is %d", currentNonce)

	require.True(t, prePendingNonce == currentNonce)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = senderAddrHex
	param[0]["to"] = receiverAddr.Hex()
	param[0]["value"] = "0xA"
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()

	_ = util.Call(t, "eth_sendTransaction", param)

	pendingNonce := util.GetNonce(t, pendingBlockNumber, senderAddrHex)
	latestNonce := util.GetNonce(t, latestBlockNumber, senderAddrHex)
	t.Logf("Latest nonce is %d", latestNonce)
	t.Logf("Pending nonce is %d", pendingNonce)

	require.True(t, currentNonce <= latestNonce)
	require.True(t, latestNonce <= pendingNonce)
	require.True(t, prePendingNonce+1 == pendingNonce)
}

func TestEth_Pending_GetBlockTransactionCountByNumber(t *testing.T) {
	waitForBlock(5)

	rpcResLatest := util.Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{latestBlockNumber})
	rpcResPending := util.Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{pendingBlockNumber})
	var preTxPendingTxCount, preTxLatestTxCount hexutil.Uint
	require.NoError(t, json.Unmarshal(rpcResPending.Result, &preTxPendingTxCount))
	require.NoError(t, json.Unmarshal(rpcResLatest.Result, &preTxLatestTxCount))
	t.Logf("Pre tx pending nonce is %d", preTxPendingTxCount)
	t.Logf("Pre tx latest nonce is %d", preTxLatestTxCount)
	require.True(t, preTxPendingTxCount == preTxLatestTxCount)

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = senderAddrHex
	param[0]["to"] = receiverAddr.Hex()
	param[0]["value"] = "0xA"
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()

	txRes := util.Call(t, "eth_sendTransaction", param)
	require.Nil(t, txRes.Error)

	rpcResLatest = util.Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{latestBlockNumber})
	rpcResPending = util.Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{pendingBlockNumber})
	var postTxPendingTxCount, postTxLatestTxCount hexutil.Uint
	require.NoError(t, json.Unmarshal(rpcResPending.Result, &postTxPendingTxCount))
	require.NoError(t, json.Unmarshal(rpcResLatest.Result, &postTxLatestTxCount))
	t.Logf("Post tx pending nonce is %d", postTxPendingTxCount)
	t.Logf("Post tx latest nonce is %d", postTxLatestTxCount)

	require.True(t, preTxPendingTxCount+1 == postTxPendingTxCount)
	require.True(t, (postTxPendingTxCount-preTxPendingTxCount) >= (postTxLatestTxCount-preTxLatestTxCount))
}

func TestEth_Pending_GetBlockByNumber(t *testing.T) {
	waitForBlock(5)

	rpcResPending := util.Call(t, "eth_getBlockByNumber", []interface{}{pendingBlockNumber, true})
	rpcResLatest := util.Call(t, "eth_getBlockByNumber", []interface{}{latestBlockNumber, true})

	var preTxLatestBlock, preTxPendingBlock map[string]interface{}
	require.NoError(t, json.Unmarshal(rpcResLatest.Result, &preTxLatestBlock))
	require.NoError(t, json.Unmarshal(rpcResPending.Result, &preTxPendingBlock))
	preTxLatestTxs := len(preTxLatestBlock["transactions"].([]interface{}))
	preTxPendingTxs := len(preTxPendingBlock["transactions"].([]interface{}))

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = senderAddrHex
	param[0]["to"] = receiverAddr.Hex()
	param[0]["value"] = "0xA"
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()

	_ = util.Call(t, "eth_sendTransaction", param)

	rpcResPending = util.Call(t, "eth_getBlockByNumber", []interface{}{pendingBlockNumber, true})
	rpcResLatest = util.Call(t, "eth_getBlockByNumber", []interface{}{latestBlockNumber, true})

	var postTxPendingBlock, postTxLatestBlock map[string]interface{}
	require.NoError(t, json.Unmarshal(rpcResPending.Result, &postTxPendingBlock))
	require.NoError(t, json.Unmarshal(rpcResLatest.Result, &postTxLatestBlock))
	postTxPendingTxs := len(postTxPendingBlock["transactions"].([]interface{}))
	postTxLatestTxs := len(postTxLatestBlock["transactions"].([]interface{}))

	require.True(t, postTxPendingTxs > preTxPendingTxs)
	require.True(t, preTxLatestTxs == postTxLatestTxs)
	require.True(t, postTxPendingTxs > preTxPendingTxs)
}

//func TestEth_Pending_GetTransactionByBlockNumberAndIndex(t *testing.T) {
//	var pendingTx []*rpctypes.Transaction
//	resPendingTxs := util.Call(t, "eth_pendingTransactions", []string{})
//	err := json.Unmarshal(resPendingTxs.Result, &pendingTx)
//	require.NoError(t, err)
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
//	txRes := util.Call(t, "eth_sendTransaction", param)
//	require.Nil(t, txRes.Error)
//
//	rpcRes := util.Call(t, "eth_getTransactionByBlockNumberAndIndex", []interface{}{"pending", "0x" + fmt.Sprintf("%X", pendingTxCount)})
//	var pendingBlockTx map[string]interface{}
//	err = json.Unmarshal(rpcRes.Result, &pendingBlockTx)
//	require.NoError(t, err)
//
//	// verify the pending tx has all the correct fields from the tx sent.
//	require.NotEmpty(t, pendingBlockTx["hash"])
//	require.Equal(t, pendingBlockTx["value"], "0xa")
//	require.Equal(t, data, pendingBlockTx["input"])
//
//	rpcRes = util.Call(t, "eth_getTransactionByBlockNumberAndIndex", []interface{}{"latest", "0x" + fmt.Sprintf("%X", pendingTxCount)})
//	var latestBlock map[string]interface{}
//	err = json.Unmarshal(rpcRes.Result, &latestBlock)
//	require.NoError(t, err)
//
//	// verify the pending trasnaction does not exist in the latest block info.
//	require.Empty(t, latestBlock)
//}
//
//func TestEth_Pending_GetTransactionByHash(t *testing.T) {
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
//	txRes := util.Call(t, "eth_sendTransaction", param)
//	var txHash common.Hash
//	err := txHash.UnmarshalJSON(txRes.Result)
//	require.NoError(t, err)
//
//	rpcRes := util.Call(t, "eth_getTransactionByHash", []interface{}{txHash})
//	var pendingBlockTx map[string]interface{}
//	err = json.Unmarshal(rpcRes.Result, &pendingBlockTx)
//	require.NoError(t, err)
//
//	// verify the pending tx has all the correct fields from the tx sent.
//	require.NotEmpty(t, pendingBlockTx)
//	require.NotEmpty(t, pendingBlockTx["hash"])
//	require.Equal(t, pendingBlockTx["value"], "0xa")
//	require.Equal(t, pendingBlockTx["input"], data)
//}
//
//func TestEth_Pending_SendTransaction_PendingNonce(t *testing.T) {
//	currNonce := util.GetNonce(t, "latest")
//	param := make([]map[string]string, 1)
//	param[0] = make(map[string]string)
//	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
//	param[0]["to"] = addrA
//	param[0]["value"] = "0xA"
//	param[0]["gasLimit"] = "0x5208"
//	param[0]["gasPrice"] = "0x1"
//
//	// first transaction
//	txRes1 := util.Call(t, "eth_sendTransaction", param)
//	require.Nil(t, txRes1.Error)
//	pendingNonce1 := util.GetNonce(t, "pending")
//	require.Greater(t, uint64(pendingNonce1), uint64(currNonce))
//
//	// second transaction
//	param[0]["to"] = "0x7f0f463c4d57b1bd3e3b79051e6c5ab703e803d9"
//	txRes2 := util.Call(t, "eth_sendTransaction", param)
//	require.Nil(t, txRes2.Error)
//	pendingNonce2 := util.GetNonce(t, "pending")
//	require.Greater(t, uint64(pendingNonce2), uint64(currNonce))
//	require.Greater(t, uint64(pendingNonce2), uint64(pendingNonce1))
//
//	// third transaction
//	param[0]["to"] = "0x7fb24493808b3f10527e3e0870afeb8a953052d2"
//	txRes3 := util.Call(t, "eth_sendTransaction", param)
//	require.Nil(t, txRes3.Error)
//	pendingNonce3 := util.GetNonce(t, "pending")
//	require.Greater(t, uint64(pendingNonce3), uint64(currNonce))
//	require.Greater(t, uint64(pendingNonce3), uint64(pendingNonce2))
//}

func waitForBlock(second int64) {
	fmt.Printf("wait %ds for a clean slate of a new block\n", second)
	time.Sleep(time.Duration(second) * time.Second)
}
