package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gorpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

var (
	contextType      = reflect.TypeOf((*context.Context)(nil)).Elem()
	errorType        = reflect.TypeOf((*error)(nil)).Elem()
	subscriptionType = reflect.TypeOf(gorpc.Subscription{})
	stringType       = reflect.TypeOf("")
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

func GetAddress(netAddr string) ([]byte, error) {
	rpcRes, err := CallWithError(netAddr, "eth_accounts", []string{})
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

type callback struct {
	fn          reflect.Value  // the function
	rcvr        reflect.Value  // receiver object of method, set if fn is method
	argTypes    []reflect.Type // input argument types
	hasCtx      bool           // method's first argument is a context (not included in argTypes)
	errPos      int            // err return idx, of -1 when method cannot return error
	isSubscribe bool           // true if this is a subscription callback
}

func Call(t *testing.T, addr string, method string, params interface{}) *Response {
	req, err := json.Marshal(CreateRequest(method, params))
	require.NoError(t, err)

	var rpcRes *Response
	time.Sleep(1 * time.Second)
	/* #nosec */

	if addr == "" {
		addr = "http://localhost:8030"
	}
	res, err := http.Post(addr, "application/json", bytes.NewBuffer(req)) //nolint:gosec
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

func CallWithError(netAddr string, method string, params interface{}) (*Response, error) {
	req, err := json.Marshal(CreateRequest(method, params))
	if err != nil {
		return nil, err
	}

	var rpcRes *Response
	//time.Sleep(1 * time.Second)
	/* #nosec */

	if netAddr == "" {
		netAddr = "http://localhost:8030"
	}
	res, err := http.Post(netAddr, "application/json", bytes.NewBuffer(req)) //nolint:gosec
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

//nolint
func GetTransactionReceipt(t *testing.T, addr string, hash ethcmn.Hash) map[string]interface{} {
	param := []string{hash.Hex()}
	rpcRes := Call(t, addr, "eth_getTransactionReceipt", param)

	receipt := make(map[string]interface{})
	err := json.Unmarshal(rpcRes.Result, &receipt)
	require.NoError(t, err)

	return receipt
}

func WaitForReceipt(t *testing.T, netAddr string, hash ethcmn.Hash) map[string]interface{} {
	for i := 0; i < 12; i++ {
		receipt := GetTransactionReceipt(t, netAddr, hash)
		if receipt != nil {
			return receipt
		}

		time.Sleep(time.Second)
	}

	return nil
}

func GetNonce(t *testing.T, netAddr string, block string, addr string) hexutil.Uint64 {
	rpcRes := Call(t, netAddr, "eth_getTransactionCount", []interface{}{addr, block})

	var nonce hexutil.Uint64
	require.NoError(t, json.Unmarshal(rpcRes.Result, &nonce))
	return nonce
}

func UnlockAllAccounts(t *testing.T, netAddr string) {
	var accts []common.Address
	rpcRes := Call(t, netAddr, "eth_accounts", []map[string]string{})
	err := json.Unmarshal(rpcRes.Result, &accts)
	require.NoError(t, err)

	for _, acct := range accts {
		t.Logf("account: %v", acct)
		rpcRes = Call(t, netAddr, "personal_unlockAccount", []interface{}{acct, ""})
		var unlocked bool
		err = json.Unmarshal(rpcRes.Result, &unlocked)
		require.NoError(t, err)
		require.True(t, unlocked)
	}
}
