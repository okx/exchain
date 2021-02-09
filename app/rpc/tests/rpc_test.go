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
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/okex/okexchain/app/rpc/types"
	"github.com/stretchr/testify/require"
	"math/big"
	"os"
	"testing"
	"time"
)

const (
	addrA                  = "0xc94770007dda54cF92009BFF0dE90c06F603a09f"
	addrAStoreKey          = 0
	defaultProtocolVersion = 65
	defaultChainID         = 65
	defaultMinGasPrice     = "0.000000001okt"
	latestBlockNumber      = "latest"
	pendingBlockNumber     = "pending"
)

var (
	receiverAddr   = ethcmn.BytesToAddress([]byte("receiver"))
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
	require.Equal(t, hexutil.Uint64(0), hashrate)
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
	inexistentAddr := ethcmn.BytesToAddress([]byte{0})
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

	var hash hexutil.Bytes
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

	var hash hexutil.Bytes
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
	require.True(t, bytes.Equal(rpcRes.Result, []byte("null")))
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

	require.Equal(t, hexutil.Uint64(0x1), nonce-preNonce)
}

//
//func TestBlockBloom(t *testing.T) {
//	hash := DeployTestContractWithFunction(t, from)
//	receipt := WaitForReceipt(t, hash)
//
//	number := receipt["blockNumber"].(string)
//	param := []interface{}{number, false}
//	rpcRes := Call(t, "eth_getBlockByNumber", param)
//
//	block := make(map[string]interface{})
//	err := json.Unmarshal(rpcRes.Result, &block)
//	require.NoError(t, err)
//
//	lb := HexToBigInt(t, block["logsBloom"].(string))
//	require.NotEqual(t, big.NewInt(0), lb)
//	require.Equal(t, hash.String(), block["transactions"].([]interface{})[0])
//}
//
//func TestEth_GetLogs_NoLogs(t *testing.T) {
//	param := make([]map[string][]string, 1)
//	param[0] = make(map[string][]string)
//	param[0]["topics"] = []string{}
//	rpcRes := Call(t, "eth_getLogs", param)
//	require.NotNil(t, rpcRes)
//	require.Nil(t, rpcRes.Error)
//
//	var logs []*ethtypes.Log
//	err := json.Unmarshal(rpcRes.Result, &logs)
//	require.NoError(t, err)
//	require.Empty(t, logs)
//}
//
//func TestEth_GetLogs_Topics_AB(t *testing.T) {
//	// TODO: this test passes on when run on its own, but fails when run with the other tests
//	if testing.Short() {
//		t.Skip("skipping TestEth_GetLogs_Topics_AB")
//	}
//
//	rpcRes := Call(t, "eth_blockNumber", []string{})
//
//	var res hexutil.Uint64
//	err := res.UnmarshalJSON(rpcRes.Result)
//	require.NoError(t, err)
//
//	param := make([]map[string]interface{}, 1)
//	param[0] = make(map[string]interface{})
//	param[0]["topics"] = []string{helloTopic, worldTopic}
//	param[0]["fromBlock"] = res.String()
//
//	hash := DeployTestContractWithFunction(t, from)
//	WaitForReceipt(t, hash)
//
//	rpcRes = Call(t, "eth_getLogs", param)
//
//	var logs []*ethtypes.Log
//	err = json.Unmarshal(rpcRes.Result, &logs)
//	require.NoError(t, err)
//
//	require.Equal(t, 1, len(logs))
//}
//
//func TestEth_GetTransactionCount(t *testing.T) {
//	// TODO: this test passes on when run on its own, but fails when run with the other tests
//	if testing.Short() {
//		t.Skip("skipping TestEth_GetTransactionCount")
//	}
//
//	prev := GetNonce(t, "latest")
//	SendTestTransaction(t, from)
//	post := GetNonce(t, "latest")
//	require.Equal(t, prev, post-1)
//}
//
//func TestEth_GetTransactionLogs(t *testing.T) {
//	// TODO: this test passes on when run on its own, but fails when run with the other tests
//	if testing.Short() {
//		t.Skip("skipping TestEth_GetTransactionLogs")
//	}
//
//	hash, _ := DeployTestContract(t, from)
//
//	param := []string{hash.String()}
//	rpcRes := Call(t, "eth_getTransactionLogs", param)
//
//	logs := new([]*ethtypes.Log)
//	err := json.Unmarshal(rpcRes.Result, logs)
//	require.NoError(t, err)
//	require.Equal(t, 1, len(*logs))
//}
//

//func TestEth_GetProof(t *testing.T) {
//	params := make([]interface{}, 3)
//	params[0] = addrA
//	params[1] = []string{fmt.Sprint(addrAStoreKey)}
//	params[2] = "latest"
//	rpcRes := Call(t, "eth_getProof", params)
//	require.NotNil(t, rpcRes)
//
//	var accRes rpctypes.AccountResult
//	err := json.Unmarshal(rpcRes.Result, &accRes)
//	require.NoError(t, err)
//	require.NotEmpty(t, accRes.AccountProof)
//	require.NotEmpty(t, accRes.StorageProof)
//
//	t.Logf("Got AccountResult %s", rpcRes.Result)
//}
//
//func TestEth_GetCode(t *testing.T) {
//	expectedRes := hexutil.Bytes{}
//	rpcRes := Call(t, "eth_getCode", []string{addrA, zeroString})
//
//	var code hexutil.Bytes
//	err := code.UnmarshalJSON(rpcRes.Result)
//
//	require.NoError(t, err)
//
//	t.Logf("Got code [%X] for %s\n", code, addrA)
//	require.True(t, bytes.Equal(expectedRes, code), "expected: %X got: %X", expectedRes, code)
//}
//

//func TestEth_NewFilter(t *testing.T) {
//	param := make([]map[string][]string, 1)
//	param[0] = make(map[string][]string)
//	param[0]["topics"] = []string{"0x0000000000000000000000000000000000000000000000000000000012341234"}
//	rpcRes := Call(t, "eth_newFilter", param)
//
//	var ID string
//	err := json.Unmarshal(rpcRes.Result, &ID)
//	require.NoError(t, err)
//}
//
//func TestEth_NewBlockFilter(t *testing.T) {
//	rpcRes := Call(t, "eth_newBlockFilter", []string{})
//
//	var ID string
//	err := json.Unmarshal(rpcRes.Result, &ID)
//	require.NoError(t, err)
//}
//
//func TestEth_GetFilterChanges_BlockFilter(t *testing.T) {
//	rpcRes := Call(t, "eth_newBlockFilter", []string{})
//
//	var ID string
//	err := json.Unmarshal(rpcRes.Result, &ID)
//	require.NoError(t, err)
//
//	time.Sleep(5 * time.Second)
//
//	changesRes := Call(t, "eth_getFilterChanges", []string{ID})
//	var hashes []ethcmn.Hash
//	err = json.Unmarshal(changesRes.Result, &hashes)
//	require.NoError(t, err)
//	require.GreaterOrEqual(t, len(hashes), 1)
//}
//
//func TestEth_GetFilterChanges_NoLogs(t *testing.T) {
//	param := make([]map[string][]string, 1)
//	param[0] = make(map[string][]string)
//	param[0]["topics"] = []string{}
//	rpcRes := Call(t, "eth_newFilter", param)
//
//	var ID string
//	err := json.Unmarshal(rpcRes.Result, &ID)
//	require.NoError(t, err)
//
//	changesRes := Call(t, "eth_getFilterChanges", []string{ID})
//
//	var logs []*ethtypes.Log
//	err = json.Unmarshal(changesRes.Result, &logs)
//	require.NoError(t, err)
//}
//
//func TestEth_GetFilterChanges_WrongID(t *testing.T) {
//	req, err := json.Marshal(CreateRequest("eth_getFilterChanges", []string{"0x1122334400000077"}))
//	require.NoError(t, err)
//
//	var rpcRes *Response
//	time.Sleep(1 * time.Second)
//	/* #nosec */
//	res, err := http.Post(HOST, "application/json", bytes.NewBuffer(req))
//	require.NoError(t, err)
//
//	decoder := json.NewDecoder(res.Body)
//	rpcRes = new(Response)
//	err = decoder.Decode(&rpcRes)
//	require.NoError(t, err)
//
//	err = res.Body.Close()
//	require.NoError(t, err)
//	require.NotNil(t, "invalid filter ID", rpcRes.Error.Message)
//}
//
//func TestEth_GetTransactionReceipt(t *testing.T) {
//	hash := SendTestTransaction(t, from)
//
//	time.Sleep(time.Second * 5)
//
//	param := []string{hash.String()}
//	rpcRes := Call(t, "eth_getTransactionReceipt", param)
//	require.Nil(t, rpcRes.Error)
//
//	receipt := make(map[string]interface{})
//	err := json.Unmarshal(rpcRes.Result, &receipt)
//	require.NoError(t, err)
//	require.NotEmpty(t, receipt)
//	require.Equal(t, "0x1", receipt["status"].(string))
//	require.Equal(t, []interface{}{}, receipt["logs"].([]interface{}))
//}
//
//func TestEth_GetTransactionReceipt_ContractDeployment(t *testing.T) {
//	hash, _ := DeployTestContract(t, from)
//
//	time.Sleep(time.Second * 5)
//
//	param := []string{hash.String()}
//	rpcRes := Call(t, "eth_getTransactionReceipt", param)
//
//	receipt := make(map[string]interface{})
//	err := json.Unmarshal(rpcRes.Result, &receipt)
//	require.NoError(t, err)
//	require.Equal(t, "0x1", receipt["status"].(string))
//
//	require.NotEqual(t, ethcmn.Address{}.String(), receipt["contractAddress"].(string))
//	require.NotNil(t, receipt["logs"])
//
//}
//
//func TestEth_GetFilterChanges_NoTopics(t *testing.T) {
//	rpcRes := Call(t, "eth_blockNumber", []string{})
//
//	var res hexutil.Uint64
//	err := res.UnmarshalJSON(rpcRes.Result)
//	require.NoError(t, err)
//
//	param := make([]map[string]interface{}, 1)
//	param[0] = make(map[string]interface{})
//	param[0]["topics"] = []string{}
//	param[0]["fromBlock"] = res.String()
//
//	// instantiate new filter
//	rpcRes = Call(t, "eth_newFilter", param)
//	require.Nil(t, rpcRes.Error)
//	var ID string
//	err = json.Unmarshal(rpcRes.Result, &ID)
//	require.NoError(t, err)
//
//	// deploy contract, emitting some event
//	DeployTestContract(t, from)
//
//	// get filter changes
//	changesRes := Call(t, "eth_getFilterChanges", []string{ID})
//
//	var logs []*ethtypes.Log
//	err = json.Unmarshal(changesRes.Result, &logs)
//	require.NoError(t, err)
//	require.Equal(t, 1, len(logs))
//}
//
//func TestEth_GetFilterChanges_Addresses(t *testing.T) {
//	t.Skip()
//	// TODO: need transaction receipts to determine contract deployment address
//}
//
//func TestEth_GetFilterChanges_BlockHash(t *testing.T) {
//	t.Skip()
//	// TODO: need transaction receipts to determine tx block
//}
//
//// hash of Hello event
//var helloTopic = "0x775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd738898"
//
//// world parameter in Hello event
//var worldTopic = "0x0000000000000000000000000000000000000000000000000000000000000011"
//
//// Tests topics case where there are topics in first two positions
//func TestEth_GetFilterChanges_Topics_AB(t *testing.T) {
//	time.Sleep(time.Second)
//
//	rpcRes := Call(t, "eth_blockNumber", []string{})
//
//	var res hexutil.Uint64
//	err := res.UnmarshalJSON(rpcRes.Result)
//	require.NoError(t, err)
//
//	param := make([]map[string]interface{}, 1)
//	param[0] = make(map[string]interface{})
//	param[0]["topics"] = []string{helloTopic, worldTopic}
//	param[0]["fromBlock"] = res.String()
//
//	// instantiate new filter
//	rpcRes = Call(t, "eth_newFilter", param)
//	var ID string
//	err = json.Unmarshal(rpcRes.Result, &ID)
//	require.NoError(t, err, string(rpcRes.Result))
//
//	DeployTestContractWithFunction(t, from)
//
//	// get filter changes
//	changesRes := Call(t, "eth_getFilterChanges", []string{ID})
//
//	var logs []*ethtypes.Log
//	err = json.Unmarshal(changesRes.Result, &logs)
//	require.NoError(t, err)
//
//	require.Equal(t, 1, len(logs))
//}
//
//func TestEth_GetFilterChanges_Topics_XB(t *testing.T) {
//	rpcRes := Call(t, "eth_blockNumber", []string{})
//
//	var res hexutil.Uint64
//	err := res.UnmarshalJSON(rpcRes.Result)
//	require.NoError(t, err)
//
//	param := make([]map[string]interface{}, 1)
//	param[0] = make(map[string]interface{})
//	param[0]["topics"] = []interface{}{nil, worldTopic}
//	param[0]["fromBlock"] = res.String()
//
//	// instantiate new filter
//	rpcRes = Call(t, "eth_newFilter", param)
//	var ID string
//	err = json.Unmarshal(rpcRes.Result, &ID)
//	require.NoError(t, err)
//
//	DeployTestContractWithFunction(t, from)
//
//	// get filter changes
//	changesRes := Call(t, "eth_getFilterChanges", []string{ID})
//
//	var logs []*ethtypes.Log
//	err = json.Unmarshal(changesRes.Result, &logs)
//	require.NoError(t, err)
//
//	require.Equal(t, 1, len(logs))
//}
//
//func TestEth_GetFilterChanges_Topics_XXC(t *testing.T) {
//	t.Skip()
//	// TODO: call test function, need tx receipts to determine contract address
//}
//
//func TestEth_PendingTransactionFilter(t *testing.T) {
//	rpcRes := Call(t, "eth_newPendingTransactionFilter", []string{})
//
//	var ID string
//	err := json.Unmarshal(rpcRes.Result, &ID)
//	require.NoError(t, err)
//
//	for i := 0; i < 5; i++ {
//		DeployTestContractWithFunction(t, from)
//	}
//
//	time.Sleep(10 * time.Second)
//
//	// get filter changes
//	changesRes := Call(t, "eth_getFilterChanges", []string{ID})
//	require.NotNil(t, changesRes)
//
//	var txs []*hexutil.Bytes
//	err = json.Unmarshal(changesRes.Result, &txs)
//	require.NoError(t, err, string(changesRes.Result))
//
//	require.True(t, len(txs) >= 2, "could not get any txs", "changesRes.Result", string(changesRes.Result))
//}
//
//func TestEth_EstimateGas(t *testing.T) {
//	param := make([]map[string]string, 1)
//	param[0] = make(map[string]string)
//	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
//	param[0]["to"] = "0x1122334455667788990011223344556677889900"
//	param[0]["value"] = "0x1"
//	rpcRes := Call(t, "eth_estimateGas", param)
//	require.NotNil(t, rpcRes)
//	require.NotEmpty(t, rpcRes.Result)
//
//	var gas string
//	err := json.Unmarshal(rpcRes.Result, &gas)
//	require.NoError(t, err, string(rpcRes.Result))
//
//	require.Equal(t, "0x1006b", gas)
//}
//
//func TestEth_EstimateGas_ContractDeployment(t *testing.T) {
//	bytecode := "0x608060405234801561001057600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a260d08061004d6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063eb8ac92114602d575b600080fd5b606060048036036040811015604157600080fd5b8101908080359060200190929190803590602001909291905050506062565b005b8160008190555080827ff3ca124a697ba07e8c5e80bebcfcc48991fc16a63170e8a9206e30508960d00360405160405180910390a3505056fea265627a7a723158201d94d2187aaf3a6790527b615fcc40970febf0385fa6d72a2344848ebd0df3e964736f6c63430005110032"
//
//	param := make([]map[string]string, 1)
//	param[0] = make(map[string]string)
//	param[0]["from"] = "0x" + fmt.Sprintf("%x", from)
//	param[0]["data"] = bytecode
//
//	rpcRes := Call(t, "eth_estimateGas", param)
//	require.NotNil(t, rpcRes)
//	require.NotEmpty(t, rpcRes.Result)
//
//	var gas hexutil.Uint64
//	err := json.Unmarshal(rpcRes.Result, &gas)
//	require.NoError(t, err, string(rpcRes.Result))
//
//	require.Equal(t, "0x1b243", gas.String())
//}
//
//func TestEth_GetBlockByNumber(t *testing.T) {
//	param := []interface{}{"0x1", false}
//	rpcRes := Call(t, "eth_getBlockByNumber", param)
//
//	block := make(map[string]interface{})
//	err := json.Unmarshal(rpcRes.Result, &block)
//	require.NoError(t, err)
//	require.Equal(t, "0x0", block["extraData"].(string))
//	require.Equal(t, []interface{}{}, block["uncles"].([]interface{}))
//}
