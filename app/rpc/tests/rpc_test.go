// This is a test utility for Ethermint's Web3 JSON-RPC services.
//
// To run these tests please first ensure you have the ethermintd running
// and have started the RPC service with `ethermintcli rest-server`.
//
// You can configure the desired HOST and MODE as well
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/okex/okexchain/app/rpc/types"
	"github.com/okex/okexchain/app/rpc/websockets"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

const (
	addrAStoreKey          = 0
	defaultProtocolVersion = 65
	defaultChainID         = 65
	defaultMinGasPrice     = "0.000000001okt"
	latestBlockNumber      = "latest"
	pendingBlockNumber     = "pending"
)

var (
	receiverAddr   = ethcmn.BytesToAddress([]byte("receiver"))
	inexistentAddr = ethcmn.BytesToAddress([]byte{0})
	inexistentHash = ethcmn.BytesToHash([]byte("inexistent hash"))
	MODE           = os.Getenv("MODE")
	from           = []byte{1}
	zeroString     = "0x0"
)

func TestMain(m *testing.M) {
	var err error
	from, err = GetAddress()
	if err != nil {
		fmt.Printf("failed to get account: %s\n", err)
		os.Exit(1)
	}

	// Start all tests
	code := m.Run()
	os.Exit(code)
}

func TestEth_Accounts(t *testing.T) {
	// all unlocked addresses
	rpcRes, err := CallWithError("eth_accounts", nil)
	require.NoError(t, err)
	require.Equal(t, 1, rpcRes.ID)

	var addrsUnlocked []ethcmn.Address
	require.NoError(t, json.Unmarshal(rpcRes.Result, &addrsUnlocked))
	require.Equal(t, addrCounter, len(addrsUnlocked))
	require.True(t, addrsUnlocked[0] == hexAddr1)
	require.True(t, addrsUnlocked[1] == hexAddr2)
}

func TestEth_ProtocolVersion(t *testing.T) {
	rpcRes, err := CallWithError("eth_protocolVersion", nil)
	require.NoError(t, err)

	var version hexutil.Uint
	require.NoError(t, json.Unmarshal(rpcRes.Result, &version))
	require.Equal(t, version, hexutil.Uint(defaultProtocolVersion))
}

func TestEth_ChainId(t *testing.T) {
	rpcRes, err := CallWithError("eth_chainId", nil)
	require.NoError(t, err)

	var chainID hexutil.Uint
	require.NoError(t, json.Unmarshal(rpcRes.Result, &chainID))
	require.Equal(t, chainID, hexutil.Uint(defaultChainID))
}

func TestEth_Syncing(t *testing.T) {
	rpcRes, err := CallWithError("eth_syncing", nil)
	require.NoError(t, err)

	// single node for test.sh -> always leading without syncing
	var catchingUp bool
	require.NoError(t, json.Unmarshal(rpcRes.Result, &catchingUp))
	require.False(t, catchingUp)

	// TODO: set an evn in multi-nodes testnet to test the sycing status of a lagging node
}

func TestEth_Coinbase(t *testing.T) {
	// single node -> always the same addr for coinbase
	rpcRes, err := CallWithError("eth_coinbase", nil)
	require.NoError(t, err)

	var coinbaseAddr1 ethcmn.Address
	require.NoError(t, json.Unmarshal(rpcRes.Result, &coinbaseAddr1))

	// wait for 5s as an block interval
	time.Sleep(5 * time.Second)

	// query again
	rpcRes, err = CallWithError("eth_coinbase", nil)
	require.NoError(t, err)

	var coinbaseAddr2 ethcmn.Address
	require.NoError(t, json.Unmarshal(rpcRes.Result, &coinbaseAddr2))

	require.Equal(t, coinbaseAddr1, coinbaseAddr2)
}

func TestEth_PowAttribute(t *testing.T) {
	// eth_mining -> always false
	rpcRes, err := CallWithError("eth_mining", nil)
	require.NoError(t, err)

	var mining bool
	require.NoError(t, json.Unmarshal(rpcRes.Result, &mining))
	require.False(t, mining)

	// eth_hashrate -> always 0
	rpcRes, err = CallWithError("eth_hashrate", nil)
	require.NoError(t, err)

	var hashrate hexutil.Uint64
	require.NoError(t, json.Unmarshal(rpcRes.Result, &hashrate))
	require.True(t, hashrate == 0)

	// eth_getUncleCountByBlockHash -> 0 for any hash
	rpcRes, err = CallWithError("eth_getUncleCountByBlockHash", []interface{}{inexistentHash})
	require.NoError(t, err)

	var uncleCount hexutil.Uint
	require.NoError(t, json.Unmarshal(rpcRes.Result, &uncleCount))
	require.True(t, uncleCount == 0)

	// eth_getUncleCountByBlockNumber -> 0 for any block number
	rpcRes, err = CallWithError("eth_getUncleCountByBlockNumber", []interface{}{latestBlockNumber})
	require.NoError(t, err)

	require.NoError(t, json.Unmarshal(rpcRes.Result, &uncleCount))
	require.True(t, uncleCount == 0)

	// eth_getUncleByBlockHashAndIndex -> always "null"
	rand.Seed(time.Now().UnixNano())
	luckyNum := int64(rand.Int())
	randomBlockHash := ethcmn.BigToHash(big.NewInt(luckyNum))
	randomIndex := hexutil.Uint(luckyNum)
	rpcRes, err = CallWithError("eth_getUncleByBlockHashAndIndex", []interface{}{randomBlockHash, randomIndex})
	require.NoError(t, err)
	assertNullFromJSONResponse(t, rpcRes.Result)

	// error check
	// miss argument
	_, err = CallWithError("eth_getUncleByBlockHashAndIndex", []interface{}{randomBlockHash})
	require.Error(t, err)

	_, err = CallWithError("eth_getUncleByBlockHashAndIndex", nil)
	require.Error(t, err)

	// eth_getUncleByBlockNumberAndIndex -> always "null"
	luckyNum = int64(rand.Int())
	randomBlockHeight := hexutil.Uint(luckyNum)
	randomIndex = hexutil.Uint(luckyNum)
	rpcRes, err = CallWithError("eth_getUncleByBlockNumberAndIndex", []interface{}{randomBlockHeight, randomIndex})
	require.NoError(t, err)
	assertNullFromJSONResponse(t, rpcRes.Result)

	// error check
	// miss argument
	_, err = CallWithError("eth_getUncleByBlockNumberAndIndex", []interface{}{randomBlockHeight})
	require.Error(t, err)

	_, err = CallWithError("eth_getUncleByBlockNumberAndIndex", nil)
	require.Error(t, err)
}

func TestEth_GasPrice(t *testing.T) {
	rpcRes, err := CallWithError("eth_gasPrice", nil)
	require.NoError(t, err)

	var gasPrice hexutil.Big
	require.NoError(t, json.Unmarshal(rpcRes.Result, &gasPrice))

	// min gas price in test.sh is "0.000000001okt"
	mgp, err := sdk.ParseDecCoin(defaultMinGasPrice)
	require.NoError(t, err)

	require.True(t, mgp.Amount.BigInt().Cmp(gasPrice.ToInt()) == 0)
}

func TestEth_BlockNumber(t *testing.T) {
	rpcRes := Call(t, "eth_blockNumber", nil)
	var blockNumber1 hexutil.Uint64
	require.NoError(t, json.Unmarshal(rpcRes.Result, &blockNumber1))

	// wait for 5s as an block interval
	time.Sleep(5 * time.Second)

	rpcRes = Call(t, "eth_blockNumber", nil)
	var blockNumber2 hexutil.Uint64
	require.NoError(t, json.Unmarshal(rpcRes.Result, &blockNumber2))

	require.True(t, blockNumber2 > blockNumber1)
}

func TestEth_GetBalance(t *testing.T) {
	// initial balance of hexAddr2 is 1000000000okt in test.sh
	initialBalance, err := sdk.ParseDecCoin("1000000000okt")
	require.NoError(t, err)

	rpcRes, err := CallWithError("eth_getBalance", []interface{}{hexAddr2, latestBlockNumber})
	require.NoError(t, err)

	var balance hexutil.Big
	require.NoError(t, json.Unmarshal(rpcRes.Result, &balance))
	require.True(t, initialBalance.Amount.Int.Cmp(balance.ToInt()) == 0)

	// query on certain block height (2)
	rpcRes, err = CallWithError("eth_getBalance", []interface{}{hexAddr2, hexutil.EncodeUint64(2)})
	require.NoError(t, json.Unmarshal(rpcRes.Result, &balance))
	require.NoError(t, err)
	require.True(t, initialBalance.Amount.Int.Cmp(balance.ToInt()) == 0)

	// query with pending -> no tx in mempool
	rpcRes, err = CallWithError("eth_getBalance", []interface{}{hexAddr2, pendingBlockNumber})
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &balance))
	require.True(t, initialBalance.Amount.Int.Cmp(balance.ToInt()) == 0)

	// inexistent addr -> zero balance
	rpcRes, err = CallWithError("eth_getBalance", []interface{}{inexistentAddr, latestBlockNumber})
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &balance))
	require.True(t, sdk.ZeroDec().Int.Cmp(balance.ToInt()) == 0)

	// error check
	// empty hex string
	_, err = CallWithError("eth_getBalance", []interface{}{hexAddr2, ""})
	require.Error(t, err)

	// missing argument
	_, err = CallWithError("eth_getBalance", []interface{}{hexAddr2})
	require.Error(t, err)
}

func TestEth_SendTransaction_Transfer(t *testing.T) {
	value := sdk.NewDec(1024)
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = hexAddr1.Hex()
	param[0]["to"] = receiverAddr.Hex()
	param[0]["value"] = (*hexutil.Big)(value.BigInt()).String()
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()

	rpcRes := Call(t, "eth_sendTransaction", param)

	var hash ethcmn.Hash
	require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))
	receipt := WaitForReceipt(t, hash)
	require.NotNil(t, receipt)
	require.Equal(t, "0x1", receipt["status"].(string))
	t.Logf("%s transfers %sokt to %s successfully\n", hexAddr1.Hex(), value.String(), receiverAddr.Hex())

	// TODO: logic bug, fix it later
	// ignore gas price -> default 'ethermint.DefaultGasPrice' on node -> successfully
	//delete(param[0], "gasPrice")
	//rpcRes = Call(t, "eth_sendTransaction", param)
	//
	//require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))
	//receipt = WaitForReceipt(t, hash)
	//require.NotNil(t, receipt)
	//require.Equal(t, "0x1", receipt["status"].(string))
	//t.Logf("%s transfers %sokt to %s successfully with nil gas price \n", hexAddr1.Hex(), value.String(), receiverAddr.Hex())

	// error check
	// sender is not unlocked on the node
	param[0]["from"] = receiverAddr.Hex()
	param[0]["to"] = hexAddr1.Hex()
	rpcRes, err := CallWithError("eth_sendTransaction", param)
	require.Error(t, err)

	// data.Data and data.Input are not same
	param[0]["from"], param[0]["to"] = param[0]["to"], param[0]["from"]
	param[0]["data"] = "0x1234567890abcdef"
	param[0]["input"] = param[0]["data"][:len(param[0]["data"])-2]
	rpcRes, err = CallWithError("eth_sendTransaction", param)
	require.Error(t, err)

	// input and toAddr are all empty
	delete(param[0], "to")
	delete(param[0], "input")
	delete(param[0], "data")

	rpcRes, err = CallWithError("eth_sendTransaction", param)
	require.Error(t, err)

	// 0 gas price
	param[0]["to"] = receiverAddr.Hex()
	param[0]["gasPrice"] = (*hexutil.Big)(sdk.ZeroDec().BigInt()).String()
	rpcRes, err = CallWithError("eth_sendTransaction", param)
	require.Error(t, err)
}

func TestEth_SendTransaction_ContractDeploy(t *testing.T) {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = hexAddr1.Hex()
	param[0]["data"] = "0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029"
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()
	rpcRes := Call(t, "eth_sendTransaction", param)

	var hash ethcmn.Hash
	require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))
	receipt := WaitForReceipt(t, hash)
	require.NotNil(t, receipt)
	require.Equal(t, "0x1", receipt["status"].(string))
	t.Logf("%s deploys contract (filled \"data\") successfully with tx hash %s\n", hexAddr1.Hex(), hash.String())

	// TODO: logic bug, fix it later
	// ignore gas price -> default 'ethermint.DefaultGasPrice' on node -> successfully
	//delete(param[0], "gasPrice")
	//rpcRes = Call(t, "eth_sendTransaction", param)
	//
	//require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))
	//receipt = WaitForReceipt(t, hash)
	//require.NotNil(t, receipt)
	//require.Equal(t, "0x1", receipt["status"].(string))
	//t.Logf("%s deploys contract successfully with tx hash %s and nil gas price\n", hexAddr1.Hex(), hash.String())

	// same payload filled in both 'input' and 'data' -> ok
	param[0]["input"] = param[0]["data"]
	rpcRes = Call(t, "eth_sendTransaction", param)

	require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))
	receipt = WaitForReceipt(t, hash)
	require.NotNil(t, receipt)
	require.Equal(t, "0x1", receipt["status"].(string))
	t.Logf("%s deploys contract (filled \"input\" and \"data\") successfully with tx hash %s\n", hexAddr1.Hex(), hash.String())

	// TODO: logic bug, fix it later
	// filled in 'input' -> ok
	//delete(param[0], "data")
	//rpcRes = Call(t, "eth_sendTransaction", param)
	//
	//require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))
	//receipt = WaitForReceipt(t, hash)
	//require.NotNil(t, receipt)
	//require.Equal(t, "0x1", receipt["status"].(string))
	//t.Logf("%s deploys contract (filled \"input\") successfully with tx hash %s\n", hexAddr1.Hex(), hash.String())

	// error check
	// sender is not unlocked on the node
	param[0]["from"] = receiverAddr.Hex()
	rpcRes, err := CallWithError("eth_sendTransaction", param)
	require.Error(t, err)

	// data.Data and data.Input are not same
	param[0]["from"] = hexAddr1.Hex()
	param[0]["input"] = param[0]["data"][:len(param[0]["data"])-2]
	rpcRes, err = CallWithError("eth_sendTransaction", param)
	require.Error(t, err)

	// 0 gas price
	delete(param[0], "input")
	param[0]["gasPrice"] = (*hexutil.Big)(sdk.ZeroDec().BigInt()).String()
	rpcRes, err = CallWithError("eth_sendTransaction", param)
	require.Error(t, err)

	// no payload of contract deployment
	delete(param[0], "data")

	rpcRes, err = CallWithError("eth_sendTransaction", param)
	require.Error(t, err)
}

func TestEth_GetStorageAt(t *testing.T) {
	expectedRes := hexutil.Bytes{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	rpcRes := Call(t, "eth_getStorageAt", []string{hexAddr1.Hex(), fmt.Sprint(addrAStoreKey), latestBlockNumber})

	var storage hexutil.Bytes
	require.NoError(t, storage.UnmarshalJSON(rpcRes.Result))

	t.Logf("Got value [%X] for %s with key %X\n", storage, hexAddr1.Hex(), addrAStoreKey)

	require.True(t, bytes.Equal(storage, expectedRes), "expected: %d (%d bytes) got: %d (%d bytes)", expectedRes, len(expectedRes), storage, len(storage))

	// error check
	// miss argument
	_, err := CallWithError("eth_getStorageAt", []string{hexAddr1.Hex(), fmt.Sprint(addrAStoreKey)})
	require.Error(t, err)

	_, err = CallWithError("eth_getStorageAt", []string{hexAddr1.Hex()})
	require.Error(t, err)
}

func TestEth_GetTransactionByHash(t *testing.T) {
	hash := sendTestTransaction(t, hexAddr1, receiverAddr, 1024)

	rpcRes := Call(t, "eth_getTransactionByHash", []interface{}{hash})

	var transaction types.Transaction
	require.NoError(t, json.Unmarshal(rpcRes.Result, &transaction))
	require.True(t, hexAddr1 == transaction.From)
	require.True(t, receiverAddr == *transaction.To)
	require.True(t, hash == transaction.Hash)
	require.True(t, transaction.Value.ToInt().Cmp(big.NewInt(1024)) == 0)
	require.True(t, transaction.GasPrice.ToInt().Cmp(defaultGasPrice.Amount.BigInt()) == 0)
	// no input for a transfer tx
	require.Equal(t, 0, len(transaction.Input))

	// hash not found -> rpcRes.Result -> "null"
	rpcRes, err := CallWithError("eth_getTransactionByHash", []interface{}{inexistentHash})
	require.NoError(t, err)
	assertNullFromJSONResponse(t, rpcRes.Result)
	require.Nil(t, rpcRes.Error)
}

func TestEth_GetTransactionCount(t *testing.T) {
	hash := sendTestTransaction(t, hexAddr1, receiverAddr, 1024)

	// sleep for a while
	time.Sleep(3 * time.Second)
	height := getBlockHeightFromTxHash(t, hash)

	rpcRes := Call(t, "eth_getTransactionCount", []interface{}{hexAddr1, height.String()})

	var nonce, preNonce hexutil.Uint64
	require.NoError(t, json.Unmarshal(rpcRes.Result, &nonce))

	// query height - 1
	rpcRes = Call(t, "eth_getTransactionCount", []interface{}{hexAddr1, (height - 1).String()})
	require.NoError(t, json.Unmarshal(rpcRes.Result, &preNonce))

	require.True(t, nonce-preNonce == 1)

	// latestBlock query
	rpcRes = Call(t, "eth_getTransactionCount", []interface{}{hexAddr1, latestBlockNumber})
	require.NoError(t, json.Unmarshal(rpcRes.Result, &preNonce))
	require.Equal(t, nonce, preNonce)

	// pendingBlock query
	rpcRes = Call(t, "eth_getTransactionCount", []interface{}{hexAddr1, pendingBlockNumber})
	require.NoError(t, json.Unmarshal(rpcRes.Result, &nonce))
	require.Equal(t, preNonce, nonce)

	// error check
	// miss argument
	_, err := CallWithError("eth_getTransactionCount", []interface{}{hexAddr1})
	require.Error(t, err)

	_, err = CallWithError("eth_getTransactionCount", nil)
	require.Error(t, err)
}

func TestEth_GetBlockTransactionCountByHash(t *testing.T) {
	hash := sendTestTransaction(t, hexAddr1, receiverAddr, 1024)

	// sleep for a while
	time.Sleep(3 * time.Second)
	blockHash := getBlockHashFromTxHash(t, hash)
	require.NotNil(t, blockHash)

	rpcRes := Call(t, "eth_getBlockTransactionCountByHash", []interface{}{*blockHash})

	var txCount hexutil.Uint
	require.NoError(t, json.Unmarshal(rpcRes.Result, &txCount))
	// only 1 tx on that height in this single node testnet
	require.True(t, txCount == 1)

	// inexistent hash -> return nil
	rpcRes = Call(t, "eth_getBlockTransactionCountByHash", []interface{}{inexistentHash})
	assertNullFromJSONResponse(t, rpcRes.Result)

	// error check
	// miss argument
	_, err := CallWithError("eth_getBlockTransactionCountByHash", nil)
	require.Error(t, err)
}

func TestEth_GetBlockTransactionCountByNumber(t *testing.T) {
	hash := sendTestTransaction(t, hexAddr1, receiverAddr, 1024)

	// sleep for a while
	time.Sleep(3 * time.Second)
	height := getBlockHeightFromTxHash(t, hash)
	require.True(t, height != 0)

	rpcRes := Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{height.String()})

	var txCount hexutil.Uint
	require.NoError(t, json.Unmarshal(rpcRes.Result, &txCount))
	// only 1 tx on that height in this single node testnet
	require.True(t, txCount == 1)

	// latestBlock query
	rpcRes = Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{latestBlockNumber})
	require.NoError(t, json.Unmarshal(rpcRes.Result, &txCount))
	// there is no tx on latest block
	require.True(t, txCount == 0)

	// pendingBlock query
	rpcRes = Call(t, "eth_getBlockTransactionCountByNumber", []interface{}{pendingBlockNumber})
	require.NoError(t, json.Unmarshal(rpcRes.Result, &txCount))
	// there is no tx on latest block and mempool
	require.True(t, txCount == 0)

	// error check
	// miss argument
	_, err := CallWithError("eth_getBlockTransactionCountByNumber", nil)
	require.Error(t, err)
	fmt.Println(err)
}

func TestEth_GetCode(t *testing.T) {
	// TODO: logic bug, fix it later
	// erc20 contract
	//hash, receipet := deployTestContract(t, hexAddr1, erc20ContractKind)
	//height := getBlockHeightFromTxHash(t, hash)
	//require.True(t, height != 0)
	//
	//rpcRes := Call(t, "eth_getCode", []interface{}{receipet["contractAddress"], height.String()})
	//var code hexutil.Bytes
	//require.NoError(t, json.Unmarshal(rpcRes.Result, &code))
	//require.True(t, strings.EqualFold(erc20ContractByteCode, code.String()))

	// test contract
	// TODO: logic bug, fix it later
	//hash, receipet := deployTestContract(t, hexAddr1, testContractKind)
	//height := getBlockHeightFromTxHash(t, hash)
	//require.True(t, height != 0)
	//
	//rpcRes := Call(t, "eth_getCode", []interface{}{receipet["contractAddress"], height.String()})
	//var code hexutil.Bytes
	//require.NoError(t, json.Unmarshal(rpcRes.Result, &code))
	//fmt.Println(testContractByteCode)
	//fmt.Println(code.String())
	//require.True(t, strings.EqualFold(testContractByteCode, code.String()))

	// error check
	// miss argument
	// TODO: use a valid contract address as the first argument in params
	_, err := CallWithError("eth_getCode", []interface{}{hexAddr1})
	require.Error(t, err)

	_, err = CallWithError("eth_getCode", nil)
	require.Error(t, err)
}

func TestEth_GetTransactionLogs(t *testing.T) {
	hash := sendTestTransaction(t, hexAddr1, receiverAddr, 1024)

	// sleep for a while
	time.Sleep(3 * time.Second)

	rpcRes := Call(t, "eth_getTransactionLogs", []interface{}{hash})
	var transactionLogs []ethtypes.Log
	require.NoError(t, json.Unmarshal(rpcRes.Result, &transactionLogs))
	// no transaction log for an evm transfer
	assertNullFromJSONResponse(t, rpcRes.Result)

	// test contract that emits an event in its constructor
	hash, receipt := deployTestContract(t, hexAddr1, testContractKind)

	rpcRes = Call(t, "eth_getTransactionLogs", []interface{}{hash})
	require.NoError(t, json.Unmarshal(rpcRes.Result, &transactionLogs))
	require.Equal(t, 1, len(transactionLogs))
	require.True(t, ethcmn.HexToAddress(receipt["contractAddress"].(string)) == transactionLogs[0].Address)
	require.True(t, hash == transactionLogs[0].TxHash)
	// event in test contract constructor keeps the value: 1024
	require.True(t, transactionLogs[0].Topics[1].Big().Cmp(big.NewInt(1024)) == 0)

	// inexistent tx hash
	rpcRes, err := CallWithError("eth_getTransactionLogs", []interface{}{inexistentHash})
	require.NoError(t, err)
	assertNullFromJSONResponse(t, rpcRes.Result)

	// error check
	// miss argument
	_, err = CallWithError("eth_getTransactionLogs", nil)
	require.Error(t, err)
}

func TestEth_Sign(t *testing.T) {
	data := []byte("context to sign")
	expectedSignature, err := signWithAccNameAndPasswd("alice", defaultPassWd, data)
	require.NoError(t, err)

	rpcRes := Call(t, "eth_sign", []interface{}{hexAddr1, hexutil.Bytes(data)})
	var sig hexutil.Bytes
	require.NoError(t, json.Unmarshal(rpcRes.Result, &sig))

	require.True(t, bytes.Equal(expectedSignature, sig))

	// error check
	// inexistent signer
	_, err = CallWithError("eth_sign", []interface{}{receiverAddr, hexutil.Bytes(data)})
	require.Error(t, err)

	// miss argument
	_, err = CallWithError("eth_sign", []interface{}{receiverAddr})
	require.Error(t, err)

	_, err = CallWithError("eth_sign", nil)
	require.Error(t, err)
}

func TestEth_Call(t *testing.T) {
	// simulate evm transfer
	callArgs := make(map[string]string)
	callArgs["from"] = hexAddr1.Hex()
	callArgs["to"] = receiverAddr.Hex()
	callArgs["value"] = hexutil.Uint(1024).String()
	callArgs["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()
	_, err := CallWithError("eth_call", []interface{}{callArgs, latestBlockNumber})
	require.NoError(t, err)

	// simulate contract deployment
	delete(callArgs, "to")
	delete(callArgs, "value")
	callArgs["data"] = erc20ContractDeployedByteCode
	_, err = CallWithError("eth_call", []interface{}{callArgs, latestBlockNumber})
	require.NoError(t, err)

	// error check
	// miss argument
	_, err = CallWithError("eth_call", []interface{}{callArgs})
	require.Error(t, err)

	_, err = CallWithError("eth_call", nil)
	require.Error(t, err)
}

func TestEth_EstimateGas_WithoutArgs(t *testing.T) {
	// error check
	// miss argument
	res, err := CallWithError("eth_estimateGas", nil)
	require.Error(t, err)
	require.Nil(t, res)
}

func TestEth_EstimateGas_Transfer(t *testing.T) {
	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["to"] = "0x1122334455667788990011223344556677889900"
	param[0]["value"] = "0x1"
	param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()
	rpcRes := Call(t, "eth_estimateGas", param)
	require.NotNil(t, rpcRes)
	require.NotEmpty(t, rpcRes.Result)

	var gas string
	err := json.Unmarshal(rpcRes.Result, &gas)
	require.NoError(t, err, string(rpcRes.Result))

	require.Equal(t, "0x100bb", gas)
}

func TestEth_EstimateGas_ContractDeployment(t *testing.T) {
	bytecode := "0x608060405234801561001057600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a260d08061004d6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063eb8ac92114602d575b600080fd5b606060048036036040811015604157600080fd5b8101908080359060200190929190803590602001909291905050506062565b005b8160008190555080827ff3ca124a697ba07e8c5e80bebcfcc48991fc16a63170e8a9206e30508960d00360405160405180910390a3505056fea265627a7a723158201d94d2187aaf3a6790527b615fcc40970febf0385fa6d72a2344848ebd0df3e964736f6c63430005110032"

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
	param[0]["data"] = bytecode

	rpcRes := Call(t, "eth_estimateGas", param)
	require.NotNil(t, rpcRes)
	require.NotEmpty(t, rpcRes.Result)

	var gas hexutil.Uint64
	err := json.Unmarshal(rpcRes.Result, &gas)
	require.NoError(t, err, string(rpcRes.Result))

	require.Equal(t, "0x1b243", gas.String())
}

func TestEth_GetBlockByHash(t *testing.T) {
	hash := sendTestTransaction(t, hexAddr1, receiverAddr, 1024)
	time.Sleep(3 * time.Second)
	expectedBlockHash := getBlockHashFromTxHash(t, hash)

	// TODO: OKExChain only supports the block query with txs' hash inside no matter what the second bool argument is.
	// 		eth rpc: 	false -> txs' hash inside
	//				  	true  -> txs full content

	// TODO: block hash bug , wait for pr merge
	//rpcRes := Call(t, "eth_getBlockByHash", []interface{}{expectedBlockHash, false})
	//var res map[string]interface{}
	//require.NoError(t, json.Unmarshal(rpcRes.Result, &res))
	//require.True(t, strings.EqualFold(expectedBlockHash, res["hash"].(string)))
	//
	//rpcRes = Call(t, "eth_getBlockByHash", []interface{}{expectedBlockHash, true})
	//require.NoError(t, json.Unmarshal(rpcRes.Result, &res))
	//require.True(t, strings.EqualFold(expectedBlockHash, res["hash"].(string)))

	// inexistent hash
	//rpcRes, err := CallWithError("eth_getBlockByHash", []interface{}{inexistentHash, false})

	// error check
	// miss argument
	_, err := CallWithError("eth_getBlockByHash", []interface{}{expectedBlockHash})
	require.Error(t, err)

	_, err = CallWithError("eth_getBlockByHash", nil)
	require.Error(t, err)
}

func TestEth_GetBlockByNumber(t *testing.T) {
	hash := sendTestTransaction(t, hexAddr1, receiverAddr, 1024)

	// sleep for a while
	time.Sleep(3 * time.Second)
	expectedHeight := getBlockHeightFromTxHash(t, hash)

	// TODO: OKExChain only supports the block query with txs' hash inside no matter what the second bool argument is.
	// 		eth rpc: 	false -> txs' hash inside
	rpcRes := Call(t, "eth_getBlockByNumber", []interface{}{expectedHeight, false})
	var res map[string]interface{}
	require.NoError(t, json.Unmarshal(rpcRes.Result, &res))
	require.True(t, strings.EqualFold(expectedHeight.String(), res["number"].(string)))

	rpcRes = Call(t, "eth_getBlockByNumber", []interface{}{expectedHeight, true})
	require.NoError(t, json.Unmarshal(rpcRes.Result, &res))
	require.True(t, strings.EqualFold(expectedHeight.String(), res["number"].(string)))

	// error check
	// future block height -> return nil without error
	rpcRes = Call(t, "eth_blockNumber", nil)
	var currentBlockHeight hexutil.Uint64
	require.NoError(t, json.Unmarshal(rpcRes.Result, &currentBlockHeight))

	rpcRes, err := CallWithError("eth_getBlockByNumber", []interface{}{currentBlockHeight + 100, false})
	require.NoError(t, err)
	assertNullFromJSONResponse(t, rpcRes.Result)

	// miss argument
	_, err = CallWithError("eth_getBlockByNumber", []interface{}{currentBlockHeight})
	require.Error(t, err)

	_, err = CallWithError("eth_getBlockByNumber", nil)
	require.Error(t, err)
}

func TestEth_GetTransactionByBlockHashAndIndex(t *testing.T) {
	hash := sendTestTransaction(t, hexAddr1, receiverAddr, 1024)

	// sleep for a while
	time.Sleep(5 * time.Second)
	blockHash, index := getBlockHashFromTxHash(t, hash), hexutil.Uint(0)
	rpcRes := Call(t, "eth_getTransactionByBlockHashAndIndex", []interface{}{blockHash, index})
	var transaction types.Transaction
	require.NoError(t, json.Unmarshal(rpcRes.Result, &transaction))
	require.True(t, hash == transaction.Hash)
	require.True(t, *blockHash == *transaction.BlockHash)
	require.True(t, hexutil.Uint64(index) == *transaction.TransactionIndex)

	// inexistent block hash
	// TODO: error:{"code":1,"log":"internal","height":1497,"codespace":"undefined"}, fix it later
	//rpcRes, err := CallWithError("eth_getTransactionByBlockHashAndIndex", []interface{}{inexistentHash, index})
	//fmt.Println(err)

	// inexistent transaction index -> nil
	rpcRes, err := CallWithError("eth_getTransactionByBlockHashAndIndex", []interface{}{blockHash, index + 100})
	require.NoError(t, err)
	assertNullFromJSONResponse(t, rpcRes.Result)

	// error check
	// miss argument
	rpcRes, err = CallWithError("eth_getTransactionByBlockHashAndIndex", []interface{}{blockHash})
	require.Error(t, err)

	rpcRes, err = CallWithError("eth_getTransactionByBlockHashAndIndex", nil)
	require.Error(t, err)
}

func TestEth_GetTransactionReceipt(t *testing.T) {
	hash := sendTestTransaction(t, hexAddr1, receiverAddr, 1024)

	// sleep for a while
	time.Sleep(3 * time.Second)
	rpcRes := Call(t, "eth_getTransactionReceipt", []interface{}{hash})

	var receipt map[string]interface{}
	require.NoError(t, json.Unmarshal(rpcRes.Result, &receipt))
	require.True(t, strings.EqualFold(hexAddr1.Hex(), receipt["from"].(string)))
	require.True(t, strings.EqualFold(receiverAddr.Hex(), receipt["to"].(string)))
	require.True(t, strings.EqualFold(hexutil.Uint(1).String(), receipt["status"].(string)))
	require.True(t, strings.EqualFold(hash.Hex(), receipt["transactionHash"].(string)))

	// contract deployment
	hash, receipt = deployTestContract(t, hexAddr1, erc20ContractKind)
	require.True(t, strings.EqualFold(hexAddr1.Hex(), receipt["from"].(string)))
	require.True(t, strings.EqualFold(hexutil.Uint(1).String(), receipt["status"].(string)))
	require.True(t, strings.EqualFold(hash.Hex(), receipt["transactionHash"].(string)))

	// inexistent hash -> nil without error
	rpcRes, err := CallWithError("eth_getTransactionReceipt", []interface{}{inexistentHash})
	require.NoError(t, err)
	assertNullFromJSONResponse(t, rpcRes.Result)

	// error check
	// miss argument
	_, err = CallWithError("eth_getTransactionReceipt", nil)
	require.Error(t, err)
}

func TestEth_PendingTransactions(t *testing.T) {
	// there will be no pending tx in mempool because of the quick grab of block building
	rpcRes := Call(t, "eth_pendingTransactions", nil)
	var transactions []types.Transaction
	require.NoError(t, json.Unmarshal(rpcRes.Result, &transactions))
	require.Zero(t, len(transactions))
}

func TestBlockBloom(t *testing.T) {
	hash, receipt := deployTestContract(t, hexAddr1, testContractKind)

	rpcRes := Call(t, "eth_getBlockByNumber", []interface{}{receipt["blockNumber"].(string), false})
	var blockInfo map[string]interface{}
	require.NoError(t, json.Unmarshal(rpcRes.Result, &blockInfo))
	logsBloom := hexToBloom(t, blockInfo["logsBloom"].(string))

	// get the transaction log with tx hash
	rpcRes = Call(t, "eth_getTransactionLogs", []interface{}{hash})
	var transactionLogs []ethtypes.Log
	require.NoError(t, json.Unmarshal(rpcRes.Result, &transactionLogs))
	require.Equal(t, 1, len(transactionLogs))

	// all the topics in the transactionLogs should be included in the logs bloom of the block
	require.True(t, logsBloom.Test(transactionLogs[0].Topics[0].Bytes()))
	require.True(t, logsBloom.Test(transactionLogs[0].Topics[1].Bytes()))
	// check the consistency of tx hash
	require.True(t, strings.EqualFold(hash.Hex(), blockInfo["transactions"].([]interface{})[0].(string)))
}

func TestEth_GetLogs_NoLogs(t *testing.T) {
	param := make([]map[string][]string, 1)
	param[0] = make(map[string][]string)
	// inexistent topics
	inexistentTopicsHash := ethcmn.BytesToHash([]byte("inexistent topics")).Hex()
	param[0]["topics"] = []string{inexistentTopicsHash}
	rpcRes, err := CallWithError("eth_getLogs", param)
	require.NoError(t, err)

	var logs []ethtypes.Log
	require.NoError(t, json.Unmarshal(rpcRes.Result, &logs))
	require.Zero(t, len(logs))

	// error check
	_, err = CallWithError("eth_getLogs", nil)
	require.Error(t, err)
}

func TestEth_GetLogs_GetTopicsFromHistory(t *testing.T) {
	_, receipt := deployTestContract(t, hexAddr1, testContractKind)
	param := make([]map[string]interface{}, 1)
	param[0] = make(map[string]interface{})
	param[0]["topics"] = []string{helloTopic, worldTopic}
	param[0]["fromBlock"] = receipt["blockNumber"].(string)

	rpcRes := Call(t, "eth_getLogs", param)

	var logs []ethtypes.Log
	require.NoError(t, json.Unmarshal(rpcRes.Result, &logs))
	require.Equal(t, 1, len(logs))
	require.Equal(t, 2, len(logs[0].Topics))
	require.True(t, logs[0].Topics[0].Hex() == helloTopic)
	require.True(t, logs[0].Topics[1].Hex() == worldTopic)

	// get block number from receipt
	blockNumber, err := hexutil.DecodeUint64(receipt["blockNumber"].(string))
	require.NoError(t, err)

	// get current block height -> there is no logs from that height
	param[0]["fromBlock"] = hexutil.Uint64(blockNumber + 1).String()

	rpcRes, err = CallWithError("eth_getLogs", param)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &logs))
	require.Zero(t, len(logs))
}

func TestEth_GetProof(t *testing.T) {
	// initial balance of hexAddr2 is 1000000000okt in test.sh
	initialBalance, err := sdk.ParseDecCoin("1000000000okt")
	require.NoError(t, err)

	rpcRes := Call(t, "eth_getProof", []interface{}{hexAddr2, []string{fmt.Sprint(addrAStoreKey)}, "latest"})
	require.NotNil(t, rpcRes)

	var accRes types.AccountResult
	require.NoError(t, json.Unmarshal(rpcRes.Result, &accRes))
	require.True(t, accRes.Address == hexAddr2)
	require.True(t, initialBalance.Amount.Int.Cmp(accRes.Balance.ToInt()) == 0)
	require.NotEmpty(t, accRes.AccountProof)
	require.NotEmpty(t, accRes.StorageProof)

	// inexistentAddr -> zero value account result
	rpcRes, err = CallWithError("eth_getProof", []interface{}{inexistentAddr, []string{fmt.Sprint(addrAStoreKey)}, "latest"})
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &accRes))
	require.True(t, accRes.Address == inexistentAddr)
	require.True(t, sdk.ZeroDec().Int.Cmp(accRes.Balance.ToInt()) == 0)

	// error check
	// miss argument
	_, err = CallWithError("eth_getProof", []interface{}{hexAddr2, []string{fmt.Sprint(addrAStoreKey)}})
	require.Error(t, err)

	_, err = CallWithError("eth_getProof", []interface{}{hexAddr2})
	require.Error(t, err)

	_, err = CallWithError("eth_getProof", nil)
	require.Error(t, err)
}

func TestEth_NewFilter(t *testing.T) {
	param := make([]map[string]interface{}, 1)
	param[0] = make(map[string]interface{})
	// random topics
	param[0]["topics"] = []ethcmn.Hash{ethcmn.BytesToHash([]byte("random topics"))}
	rpcRes := Call(t, "eth_newFilter", param)

	var ID string
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	require.NotZero(t, ID)

	// fromBlock: latest, toBlock: latest -> no error
	delete(param[0], "topics")
	param[0]["fromBlock"] = latestBlockNumber
	param[0]["toBlock"] = latestBlockNumber
	rpcRes, err := CallWithError("eth_newFilter", param)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	require.NotZero(t, ID)

	// fromBlock: nil, toBlock: latest -> no error
	delete(param[0], "fromBlock")
	rpcRes, err = CallWithError("eth_newFilter", param)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	require.NotZero(t, ID)

	// fromBlock: latest, toBlock: nil -> no error
	delete(param[0], "toBlock")
	param[0]["fromBlock"] = latestBlockNumber
	rpcRes, err = CallWithError("eth_newFilter", param)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	require.NotZero(t, ID)

	// fromBlock: pending, toBlock: pending -> no error
	param[0]["fromBlock"] = pendingBlockNumber
	param[0]["toBlock"] = pendingBlockNumber
	rpcRes, err = CallWithError("eth_newFilter", param)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	require.NotZero(t, ID)

	// fromBlock: latest, toBlock: pending -> no error
	param[0]["fromBlock"] = latestBlockNumber
	rpcRes, err = CallWithError("eth_newFilter", param)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	require.NotZero(t, ID)

	// toBlock > fromBlock -> no error
	param[0]["fromBlock"] = (*hexutil.Big)(big.NewInt(2)).String()
	param[0]["toBlock"] = (*hexutil.Big)(big.NewInt(3)).String()
	rpcRes, err = CallWithError("eth_newFilter", param)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	require.NotZero(t, ID)

	// error check
	// miss argument
	_, err = CallWithError("eth_newFilter", nil)
	require.Error(t, err)

	// fromBlock > toBlock -> error: invalid from and to block combination: from > to
	param[0]["fromBlock"] = (*hexutil.Big)(big.NewInt(3)).String()
	param[0]["toBlock"] = (*hexutil.Big)(big.NewInt(2)).String()
	rpcRes, err = CallWithError("eth_newFilter", param)
	require.Error(t, err)

	// fromBlock: pending, toBlock: latest
	param[0]["fromBlock"] = pendingBlockNumber
	param[0]["toBlock"] = latestBlockNumber
	rpcRes, err = CallWithError("eth_newFilter", param)
	require.Error(t, err)
}

func TestEth_NewBlockFilter(t *testing.T) {
	rpcRes := Call(t, "eth_newBlockFilter", nil)

	var ID string
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	require.NotZero(t, ID)
}

func TestEth_GetFilterChanges_BlockFilter(t *testing.T) {
	rpcRes := Call(t, "eth_newBlockFilter", nil)

	var ID string
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))

	// wait for block generation
	time.Sleep(5 * time.Second)

	changesRes := Call(t, "eth_getFilterChanges", []interface{}{ID})
	var hashes []ethcmn.Hash
	require.NoError(t, json.Unmarshal(changesRes.Result, &hashes))
	require.GreaterOrEqual(t, len(hashes), 1)

	// error check
	// miss argument
	_, err := CallWithError("eth_getFilterChanges", nil)
	require.Error(t, err)
}

func TestEth_GetFilterChanges_NoLogs(t *testing.T) {
	param := make([]map[string]interface{}, 1)
	param[0] = make(map[string]interface{})
	param[0]["topics"] = []ethcmn.Hash{ethcmn.BytesToHash([]byte("random topics"))}

	rpcRes := Call(t, "eth_newFilter", param)

	var ID string
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))

	changesRes := Call(t, "eth_getFilterChanges", []interface{}{ID})

	var logs []ethtypes.Log
	require.NoError(t, json.Unmarshal(changesRes.Result, &logs))
	// no logs
	require.Empty(t, logs)
}

func TestEth_GetFilterChanges_WrongID(t *testing.T) {
	// ID's length is 16
	inexistentID := "0x1234567890abcdef"
	_, err := CallWithError("eth_getFilterChanges", []interface{}{inexistentID})
	require.Error(t, err)
}

func TestEth_GetFilterChanges_NoTopics(t *testing.T) {
	// create a new filter with no topics and latest block height for "fromBlock"
	param := make([]map[string]interface{}, 1)
	param[0] = make(map[string]interface{})
	param[0]["fromBlock"] = latestBlockNumber

	rpcRes := Call(t, "eth_newFilter", param)
	require.Nil(t, rpcRes.Error)
	var ID string
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	t.Logf("create filter successfully with ID %s\n", ID)

	// deploy contract with emitting events
	_, _ = deployTestContract(t, hexAddr1, testContractKind)

	// get filter changes
	changesRes := Call(t, "eth_getFilterChanges", []string{ID})

	var logs []ethtypes.Log
	require.NoError(t, json.Unmarshal(changesRes.Result, &logs))
	require.Equal(t, 1, len(logs))
}

func TestEth_GetFilterChanges_Addresses(t *testing.T) {
	// TODO: logic bug, fix it later
	//// deploy contract with emitting events
	//_, receipt := deployTestContract(t, hexAddr1, testContractKind)
	//contractAddrHex := receipt["contractAddress"].(string)
	//blockHeight := receipt["blockNumber"].(string)
	//// create a filter
	//param := make([]map[string]interface{}, 1)
	//param[0] = make(map[string]interface{})
	//// focus on the contract by its address
	//param[0]["addresses"] = []string{contractAddrHex}
	//param[0]["topics"] = []string{helloTopic, worldTopic}
	//param[0]["fromBlock"] = blockHeight
	//rpcRes := Call(t, "eth_newFilter", param)
	//
	//var ID string
	//require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	//t.Logf("create filter focusing on contract %s successfully with ID %s\n", contractAddrHex, ID)
	//
	//// get filter changes
	//changesRes := Call(t, "eth_getFilterChanges", []string{ID})
	//
	//var logs []ethtypes.Log
	//require.NoError(t, json.Unmarshal(changesRes.Result, &logs))
	//require.Equal(t, 1, len(logs))
}

func TestEth_GetFilterChanges_BlockHash(t *testing.T) {
	// TODO: logic bug, fix it later
	//// deploy contract with emitting events
	//_, receipt := deployTestContract(t, hexAddr1, testContractKind)
	//blockHash := receipt["blockHash"].(string)
	//contractAddrHex := receipt["contractAddress"].(string)
	//// create a filter
	//param := make([]map[string]interface{}, 1)
	//param[0] = make(map[string]interface{})
	//// focus on the contract by its address
	//param[0]["blockHash"] = blockHash
	//param[0]["addresses"] = []string{contractAddrHex}
	//param[0]["topics"] = []string{helloTopic, worldTopic}
	//rpcRes := Call(t, "eth_newFilter", param)
	//
	//var ID string
	//require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	//t.Logf("create filter focusing on contract %s in the block with block hash %s successfully with ID %s\n", contractAddrHex, blockHash, ID)
	//
	//// get filter changes
	//changesRes := Call(t, "eth_getFilterChanges", []string{ID})
	//
	//var logs []ethtypes.Log
	//require.NoError(t, json.Unmarshal(changesRes.Result, &logs))
	//require.Equal(t, 1, len(logs))
}

// Tests topics case where there are topics in first two positions
func TestEth_GetFilterChanges_Topics_AB(t *testing.T) {
	param := make([]map[string]interface{}, 1)
	param[0] = make(map[string]interface{})
	// set topics in filter with A && B
	param[0]["topics"] = []string{helloTopic, worldTopic}
	param[0]["fromBlock"] = latestBlockNumber

	// create new filter
	rpcRes := Call(t, "eth_newFilter", param)

	var ID string
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	t.Logf("create filter successfully with ID %s\n", ID)

	// deploy contract with emitting events
	_, _ = deployTestContract(t, hexAddr1, testContractKind)

	// get filter changes
	changesRes := Call(t, "eth_getFilterChanges", []string{ID})

	var logs []ethtypes.Log
	require.NoError(t, json.Unmarshal(changesRes.Result, &logs))
	require.Equal(t, 1, len(logs))
}

func TestEth_GetFilterChanges_Topics_XB(t *testing.T) {
	param := make([]map[string]interface{}, 1)
	param[0] = make(map[string]interface{})
	// set topics in filter with X && B
	param[0]["topics"] = []interface{}{nil, worldTopic}
	param[0]["fromBlock"] = latestBlockNumber

	// create new filter
	rpcRes := Call(t, "eth_newFilter", param)

	var ID string
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	t.Logf("create filter successfully with ID %s\n", ID)

	// deploy contract with emitting events
	_, _ = deployTestContract(t, hexAddr1, testContractKind)

	// get filter changes
	changesRes := Call(t, "eth_getFilterChanges", []string{ID})

	var logs []ethtypes.Log
	require.NoError(t, json.Unmarshal(changesRes.Result, &logs))
	require.Equal(t, 1, len(logs))
}

//func TestEth_GetFilterChanges_Topics_XXC(t *testing.T) {
//	t.Skip()
//	// TODO: call test function, need tx receipts to determine contract address
//}

func TestEth_PendingTransactionFilter(t *testing.T) {
	rpcRes := Call(t, "eth_newPendingTransactionFilter", nil)

	var ID string
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))

	for i := 0; i < 5; i++ {
		_, _ = deployTestContract(t, hexAddr1, erc20ContractKind)
	}

	time.Sleep(10 * time.Second)

	// get filter changes
	changesRes := Call(t, "eth_getFilterChanges", []string{ID})
	require.NotNil(t, changesRes)

	var txs []hexutil.Bytes
	require.NoError(t, json.Unmarshal(changesRes.Result, &txs))

	require.True(t, len(txs) >= 2, "could not get any txs", "changesRes.Result", string(changesRes.Result))
}

func TestEth_UninstallFilter(t *testing.T) {
	// create a new filter, get id
	rpcRes := Call(t, "eth_newBlockFilter", nil)
	var ID string
	require.NoError(t, json.Unmarshal(rpcRes.Result, &ID))
	require.NotZero(t, ID)

	// based on id, uninstall filter
	rpcRes = Call(t, "eth_uninstallFilter", []string{ID})
	require.NotNil(t, rpcRes)
	var status bool
	require.NoError(t, json.Unmarshal(rpcRes.Result, &status))
	require.Equal(t, true, status)

	// uninstall a non-existent filter
	rpcRes = Call(t, "eth_uninstallFilter", []string{ID})
	require.NotNil(t, rpcRes)
	require.NoError(t, json.Unmarshal(rpcRes.Result, &status))
	require.Equal(t, false, status)

}

func TestEth_Subscribe_And_UnSubscribe(t *testing.T) {
	// create websocket
	origin, url := "http://127.0.0.1:8546/", "ws://127.0.0.1:8546"
	ws, err := websocket.Dial(url, "", origin)
	require.NoError(t, err)
	defer func() {
		// close websocket
		err = ws.Close()
		require.NoError(t, err)
	}()

	// send valid message
	validMessage := []byte(`{"id": 2, "method": "eth_subscribe", "params": ["newHeads"]}`)
	excuteValidMessage(t, ws, validMessage)

	// send invalid message
	invalidMessage := []byte(`{"id": 2, "method": "eth_subscribe", "params": ["non-existent method"]}`)
	excuteInvalidMessage(t, ws, invalidMessage)

	invalidMessage = []byte(`{"id": 2, "method": "eth_subscribe", "params": [""]}`)
	excuteInvalidMessage(t, ws, invalidMessage)
}

func excuteValidMessage(t *testing.T, ws  *websocket.Conn, message []byte) {
	fmt.Println("Send:", string(message))
	_, err := ws.Write(message)
	require.NoError(t, err)

	msg := make([]byte, 10240)
	// receive subscription id
	n, err := ws.Read(msg)
	require.NoError(t, err)
	var res Response
	require.NoError(t, json.Unmarshal(msg[:n], &res))
	subscriptionId := string(res.Result)

	// receive message three times
	for i := 0; i < 3; i++ {
		n, err = ws.Read(msg)
		require.NoError(t, err)
		fmt.Println("Receive:", string(msg[:n]))
	}

	// cancel the subscription
	cancelMsg := fmt.Sprintf(`{"id": 2, "method": "eth_unsubscribe", "params": [%s]}`, subscriptionId)
	fmt.Println("Send:", cancelMsg)
	_, err = ws.Write([]byte(cancelMsg))
	require.NoError(t, err)

	// receive the result of eth_unsubscribe
	n, err = ws.Read(msg)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(msg[:n], &res))
	require.Equal(t, "true", string(res.Result))
}

func excuteInvalidMessage(t *testing.T, ws  *websocket.Conn, message []byte) {
	fmt.Println("Send:", string(message))
	_, err := ws.Write(message)
	require.NoError(t, err)

	msg := make([]byte, 10240)
	// receive error msg
	n, err := ws.Read(msg)
	require.NoError(t, err)

	var res Response
	require.NoError(t, json.Unmarshal(msg[:n], &res))
	require.Equal(t, -32600, res.Error.Code)
	require.Equal(t, 0, res.ID)
}

func TestWebsocket_PendingTransaction(t *testing.T) {
	// create websocket
	origin, url := "http://127.0.0.1:8546/", "ws://127.0.0.1:8546"
	ws, err := websocket.Dial(url, "", origin)
	require.NoError(t, err)
	defer func() {
		// close websocket
		err = ws.Close()
		require.NoError(t, err)
	}()

	// send message to call newPendingTransactions ws api
	_, err = ws.Write([]byte(`{"id": 2, "method": "eth_subscribe", "params": ["newPendingTransactions"]}`))
	require.NoError(t, err)

	msg := make([]byte, 10240)
	// receive subscription id
	n, err := ws.Read(msg)
	require.NoError(t, err)
	var res Response
	require.NoError(t, json.Unmarshal(msg[:n], &res))
	subscriptionId := string(res.Result)

	// send transactions
	var expectedHashList [3]ethcmn.Hash
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			param := make([]map[string]string, 1)
			param[0] = make(map[string]string)
			param[0]["from"] = hexAddr1.Hex()
			param[0]["data"] = "0x6080604052348015600f57600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a2603580604b6000396000f3fe6080604052600080fdfea165627a7a723058206cab665f0f557620554bb45adf266708d2bd349b8a4314bdff205ee8440e3c240029"
			param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()
			rpcRes := Call(t, "eth_sendTransaction", param)

			var hash ethcmn.Hash
			require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))
			expectedHashList[i] = hash
		}
	}()
	var actualHashList [3]ethcmn.Hash
	// receive message three times
	for i := 0; i < 3; i++ {
		n, err = ws.Read(msg)
		require.NoError(t, err)
		var notification websockets.SubscriptionNotification
		require.NoError(t, json.Unmarshal(msg[:n], &notification))
		actualHashList[i] = ethcmn.HexToHash(notification.Params.Result.(string))
	}
	wg.Wait()
	require.EqualValues(t, expectedHashList, actualHashList)

	// cancel the subscription
	cancelMsg := fmt.Sprintf(`{"id": 2, "method": "eth_unsubscribe", "params": [%s]}`, subscriptionId)
	_, err = ws.Write([]byte(cancelMsg))
	require.NoError(t, err)
}

//{} or nil          matches any topic list
//{A}                matches topic A in first position
//{{}, {B}}          matches any topic in first position AND B in second position
//{{A}, {B}}         matches topic A in first position AND B in second position
//{{A, B}, {C, D}}   matches topic (A OR B) in first position AND (C OR D) in second position
func TestWebsocket_Logs(t *testing.T) {
	contractAddr1, contractAddr2, contractAddr3 := deployTestTokenContract(t), deployTestTokenContract(t), deployTestTokenContract(t)

	// init test cases
	tests := []struct {
		addressList   string   // input
		topicsList    string   // input
		expected      int // expected result
	}{
		{"","",21},                                    // matches any address & any topics
		{fmt.Sprintf(`"address":"%s"`, contractAddr1),"",7},                // matches one address & any topics
		{fmt.Sprintf(`"address":["%s","%s"]`, contractAddr1, contractAddr2),"",14}, // matches two addressses & any topics
		//{fmt.Sprintf(`"address":["%s","%s"]`, contractAddr1, contractAddr2),fmt.Sprintf(`"topics":"%s"`, approveFuncHash),true},
		//{fmt.Sprintf(`"address":["%s","%s"]`, contractAddr1, contractAddr2),fmt.Sprintf(`"topics":[null, null, ["%s"]]`, recvAddr1Hash),true},
		//{fmt.Sprintf(`"address":["%s","%s"]`, contractAddr1, contractAddr2),fmt.Sprintf(`"topics":[["%s"], null, ["%s"]]`, approveFuncHash, recvAddr1Hash),true},
		//{fmt.Sprintf(`"address":["%s","%s"]`, contractAddr1, contractAddr2),fmt.Sprintf(`"topics":[["%s","%s"], null, ["%s","%s"]]`, approveFuncHash,transferFuncHash, recvAddr1Hash, recvAddr2Hash),true},
	}

	// create websocket
	origin, url := "http://127.0.0.1:8546/", "ws://127.0.0.1:8546"
	ws, err := websocket.Dial(url, "", origin)
	require.NoError(t, err)
	defer func() {
		// close websocket
		err := ws.Close()
		require.NoError(t, err)
	}()

	msg := make([]byte, 10240)
	for _, test := range tests {
		// fulfill parameters
		param := assembleParameters(test.addressList, test.topicsList)
		_, err = ws.Write([]byte(param))
		require.NoError(t, err)

		// receive subscription id
		n, err := ws.Read(msg)
		var res Response
		require.NoError(t, err)
		require.NoError(t, json.Unmarshal(msg[:n], &res))
		subscriptionId := string(res.Result)

		// send txs
		// fetch logs
		var wg sync.WaitGroup
		wg.Add(1)
		go sendTxs(t, &wg, contractAddr1, contractAddr2, contractAddr3)
		for i := 0; i < test.expected; i++ {
			n, err = ws.Read(msg)
			require.NoError(t, err)
			var notification websockets.SubscriptionNotification
			require.NoError(t, json.Unmarshal(msg[:n], &notification))
		}
		wg.Wait()

		// cancel the subscription
		cancelMsg := fmt.Sprintf(`{"id": 2, "method": "eth_unsubscribe", "params": [%s]}`, subscriptionId)
		_, err = ws.Write([]byte(cancelMsg))
		require.NoError(t, err)
	}
}

func deployTestTokenContract(t *testing.T) string {
	param := make([]map[string]string, 1)
	param[0] = map[string]string{
		"from": hexAddr1.Hex(),
		"data": ttokenContractByteCode,
		"gasPrice": (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String(),
	}
	rpcRes := Call(t, "eth_sendTransaction", param)
	var hash ethcmn.Hash
	require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))
	receipt := WaitForReceipt(t, hash)
	require.NotNil(t, receipt)
	contractAddr, ok := receipt["contractAddress"].(string)
	require.True(t, ok)
	return contractAddr
}

func verifyWebSocketRecvNum(t *testing.T, addressList, topicsList string, expected int) {
	// create websocket
	origin, url := "http://127.0.0.1:8546/", "ws://127.0.0.1:8546"
	ws, err := websocket.Dial(url, "", origin)
	require.NoError(t, err)
	defer func() {
		// close websocket
		err := ws.Close()
		require.NoError(t, err)
	}()

	// fulfill parameters
	param := assembleParameters(addressList, topicsList)
	_, err = ws.Write([]byte(param))
	require.NoError(t, err)

	msg := make([]byte, 10240)
	// receive subscription id
	n, err := ws.Read(msg)
	var res Response
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(msg[:n], &res))
	subscriptionId := string(res.Result)

	for i := 0; i < expected; i++ {
		n, err = ws.Read(msg)
		require.NoError(t, err)
		var notification websockets.SubscriptionNotification
		require.NoError(t, json.Unmarshal(msg[:n], &notification))
	}

	// cancel the subscription
	cancelMsg := fmt.Sprintf(`{"id": 2, "method": "eth_unsubscribe", "params": [%s]}`, subscriptionId)
	_, err = ws.Write([]byte(cancelMsg))
	require.NoError(t, err)
}

func assembleParameters(addressList string, topicsList string) string {
	var param string
	if addressList == "" {
		param = topicsList
	}
	if topicsList == "" {
		param = addressList
	}
	if addressList != "" && topicsList != "" {
		param = addressList+","+topicsList
	}
	return fmt.Sprintf(`{"id": 2, "method": "eth_subscribe", "params": ["logs",{%s}]}`, param)
}

func sendTxs(t *testing.T, wg *sync.WaitGroup, contractAddrs ...string) {
	dataList := []string{
		// 0. mint  4294967295coin -> 0x2cf4ea7df75b513509d95946b43062e26bd88035
		"0x40c10f190000000000000000000000002cf4ea7df75b513509d95946b43062e26bd8803500000000000000000000000000000000000000000000000000000000ffffffff",
		// 1. approve 12345678coin -> 0x9ad84c8630e0282f78e5479b46e64e17779e3cfb
		"0x095ea7b30000000000000000000000009ad84c8630e0282f78e5479b46e64e17779e3cfb0000000000000000000000000000000000000000000000000000000000bc614e",
		// 2. approve 12345678coin -> 0xc9c9b43322f5e1dc401252076fa4e699c9122cd6
		"0x095ea7b3000000000000000000000000c9c9b43322f5e1dc401252076fa4e699c9122cd60000000000000000000000000000000000000000000000000000000000bc614e",
		// 3. approve 12345678coin -> 0x2B5Cf24AeBcE90f0B8f80Bc42603157b27cFbf47
		"0x095ea7b30000000000000000000000002b5cf24aebce90f0b8f80bc42603157b27cfbf470000000000000000000000000000000000000000000000000000000000bc614e",
		// 4. transfer 1234coin    -> 0x9ad84c8630e0282f78e5479b46e64e17779e3cfb
		"0xa9059cbb0000000000000000000000009ad84c8630e0282f78e5479b46e64e17779e3cfb00000000000000000000000000000000000000000000000000000000000004d2",
		// 5. transfer 1234coin    -> 0xc9c9b43322f5e1dc401252076fa4e699c9122cd6
		"0xa9059cbb000000000000000000000000c9c9b43322f5e1dc401252076fa4e699c9122cd600000000000000000000000000000000000000000000000000000000000004d2",
		// 6. transfer 1234coin    -> 0x2B5Cf24AeBcE90f0B8f80Bc42603157b27cFbf47
		"0xa9059cbb0000000000000000000000002b5cf24aebce90f0b8f80bc42603157b27cfbf4700000000000000000000000000000000000000000000000000000000000004d2",
	}
	defer wg.Done()
	for _, contractAddr := range contractAddrs{
		for i := 0; i < 7; i++{
			param := make([]map[string]string, 1)
			param[0] = make(map[string]string)
			param[0]["from"] = hexAddr1.Hex()
			param[0]["to"] = contractAddr
			param[0]["data"] = dataList[i]
			param[0]["gasPrice"] = (*hexutil.Big)(defaultGasPrice.Amount.BigInt()).String()
			rpcRes := Call(t, "eth_sendTransaction", param)
			var hash ethcmn.Hash
			require.NoError(t, json.Unmarshal(rpcRes.Result, &hash))

			time.Sleep(time.Second*1)
		}
	}
}