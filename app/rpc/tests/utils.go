package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

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

var (
	HOST = os.Getenv("HOST")
)

func GetAddress() ([]byte, error) {
	rpcRes, err := CallWithError("eth_accounts", []string{})
	if err != nil {
		return nil, err
	}

	var res []hexutil.Bytes
	err = json.Unmarshal(rpcRes.Result, &res)
	if err != nil {
		return nil, err
	}

	return res[0], nil
}

func CreateRequest(method string, params interface{}) Request {
	return Request{
		Version: "2.0",
		Method:  method,
		Params:  params,
		ID:      1,
	}
}

func Call(t *testing.T, method string, params interface{}) *Response {
	req, err := json.Marshal(CreateRequest(method, params))
	require.NoError(t, err)

	var rpcRes *Response
	time.Sleep(1 * time.Second)
	/* #nosec */

	if HOST == "" {
		HOST = "http://localhost:8545"
	}
	res, err := http.Post(HOST, "application/json", bytes.NewBuffer(req)) //nolint:gosec
	require.NoError(t, err)

	decoder := json.NewDecoder(res.Body)
	rpcRes = new(Response)
	err = decoder.Decode(&rpcRes)
	require.NoError(t, err)

	err = res.Body.Close()
	require.NoError(t, err)
	require.Nil(t, rpcRes.Error)

	return rpcRes
}

func CallWithError(method string, params interface{}) (*Response, error) {
	req, err := json.Marshal(CreateRequest(method, params))
	if err != nil {
		return nil, err
	}

	var rpcRes *Response
	time.Sleep(1 * time.Second)
	/* #nosec */

	if HOST == "" {
		HOST = "http://localhost:8545"
	}
	res, err := http.Post(HOST, "application/json", bytes.NewBuffer(req)) //nolint:gosec
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(res.Body)
	rpcRes = new(Response)
	err = decoder.Decode(&rpcRes)
	if err != nil {
		return nil, err
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	if rpcRes.Error != nil {
		return nil, fmt.Errorf(rpcRes.Error.Message)
	}

	return rpcRes, nil
}

func DeployTestContractWithFunction(t *testing.T, addr []byte) ethcmn.Hash {
	// pragma solidity ^0.5.1;

	// contract Test {
	//     event Hello(uint256 indexed world);
	//     event TestEvent(uint256 indexed a, uint256 indexed b);

	//     uint256 myStorage;

	//     constructor() public {
	//         emit Hello(17);
	//     }

	//     function test(uint256 a, uint256 b) public {
	//         myStorage = a;
	//         emit TestEvent(a, b);
	//     }
	// }

	bytecode := "0x608060405234801561001057600080fd5b5060117f775a94827b8fd9b519d36cd827093c664f93347070a554f65e4a6f56cd73889860405160405180910390a260d08061004d6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c8063eb8ac92114602d575b600080fd5b606060048036036040811015604157600080fd5b8101908080359060200190929190803590602001909291905050506062565b005b8160008190555080827ff3ca124a697ba07e8c5e80bebcfcc48991fc16a63170e8a9206e30508960d00360405160405180910390a3505056fea265627a7a723158201d94d2187aaf3a6790527b615fcc40970febf0385fa6d72a2344848ebd0df3e964736f6c63430005110032"

	param := make([]map[string]string, 1)
	param[0] = make(map[string]string)
	param[0]["from"] = "0x" + fmt.Sprintf("%x", addr)
	param[0]["data"] = bytecode
	param[0]["gaslimit"] = "0x2000000"
	param[0]["gasprice"] = "0x2000000000"

	rpcRes := Call(t, "eth_sendTransaction", param)

	var hash ethcmn.Hash
	err := json.Unmarshal(rpcRes.Result, &hash)
	require.NoError(t, err)

	receipt := WaitForReceipt(t, hash)
	require.NotNil(t, receipt, "transaction failed")
	require.Equal(t, "0x1", receipt["status"].(string))

	return hash
}

//nolint
func GetTransactionReceipt(t *testing.T, hash ethcmn.Hash) map[string]interface{} {
	param := []string{hash.Hex()}
	rpcRes := Call(t, "eth_getTransactionReceipt", param)

	receipt := make(map[string]interface{})
	err := json.Unmarshal(rpcRes.Result, &receipt)
	require.NoError(t, err)

	return receipt
}

func WaitForReceipt(t *testing.T, hash ethcmn.Hash) map[string]interface{} {
	for i := 0; i < 12; i++ {
		receipt := GetTransactionReceipt(t, hash)
		if receipt != nil {
			return receipt
		}

		time.Sleep(time.Second)
	}

	return nil
}

func GetNonce(t *testing.T, block string) hexutil.Uint64 {
	from, err := GetAddress()
	require.NoError(t, err)

	param := []interface{}{hexutil.Bytes(from), block}
	rpcRes := Call(t, "eth_getTransactionCount", param)

	var nonce hexutil.Uint64
	err = json.Unmarshal(rpcRes.Result, &nonce)
	require.NoError(t, err)
	return nonce
}

func UnlockAllAccounts(t *testing.T) {
	var accts []common.Address
	rpcRes := Call(t, "eth_accounts", []map[string]string{})
	err := json.Unmarshal(rpcRes.Result, &accts)
	require.NoError(t, err)

	for _, acct := range accts {
		t.Logf("account: %v", acct)
		rpcRes = Call(t, "personal_unlockAccount", []interface{}{acct, ""})
		var unlocked bool
		err = json.Unmarshal(rpcRes.Result, &unlocked)
		require.NoError(t, err)
		require.True(t, unlocked)
	}
}
