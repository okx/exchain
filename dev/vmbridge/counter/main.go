package main

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TestType string

const (
	abiFilePath = "./evmContract/evmContract.abi"
	binFilePath = "./evmContract/evmContract.bin"

	Deploy  = TestType("deploy")
	Execute = TestType("execute")
	Query   = TestType("query")

	GasPrice int64  = 1000000000 // 1 gwei
	GasLimit uint64 = 30000000
)

var (
	client                                         *ethclient.Client
	evmContract, WasmContract, deltaValue, privKey *string
	chainID                                        *big.Int
)

func main() {
	actionTypeParam := flag.String("action", "deploy", "deploy/execute/call")
	evmContract = flag.String("contract", "", "counter contract address")
	WasmContract = flag.String("wasmContract", "", "wasm contract address")
	deltaValue = flag.String("delta", "1", "wasm contract address")
	privKey = flag.String("key", "", "private key")

	flag.Parse()

	var err error
	client, err = ethclient.Dial("http://127.0.0.1:8545")
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*5000)
	defer cancel()
	chainID, err = client.ChainID(ctx)
	if err != nil {
		panic(err)
	}

	// privKey := []string{
	// 	"8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17",
	// }

	// for _, k := range privKey {
	test := func(key string) {

		switch TestType(*actionTypeParam) {
		case Deploy:
			counterDeploy(*privKey)
		case Execute:
			counterExecute(*privKey, common.HexToAddress(*evmContract), common.HexToAddress(*WasmContract))
		case Query:
			counterQuery(*privKey, common.HexToAddress(*evmContract))
		default:
			panic("action not found")
		}
	}
	test(*privKey)
	// }
}

func counterDeploy(privKey string) (err error) {
	var (
		privateKey    *ecdsa.PrivateKey
		senderAddress common.Address
	)

	privateKey, senderAddress = initKey(privKey)

	counterContract := newContract("counter", "", abiFilePath, binFilePath)

	return deployContract(client, senderAddress, privateKey, counterContract, time.Second*1)

}

func counterExecute(privKey string, evmContract common.Address, wasmContract common.Address) error {
	privateKey, _ := initKey(privKey)

	abiFile, err := ioutil.ReadFile(abiFilePath)
	if err != nil {
		log.Printf("Failed to read ABI file: %v", err)
		return err
	}

	contractAbi, err := abi.JSON(bytes.NewReader(abiFile))
	if err != nil {
		log.Printf("Failed to parse ABI: %v", err)
		return err
	}

	txHash, err := sendTransaction(client, privateKey, evmContract, contractAbi, "addCounterForWasm", []interface{}{wasmContract.Hex(), &deltaValue})
	if err != nil {
		log.Printf("error: %+v", err)
		return err
	}

	time.Sleep(3 * time.Second)
	_, err = getReceipt(client, txHash)
	if err != nil {
		log.Printf("error: %+v", err)
		return err
	}

	fmt.Printf("hash is: %v\n", txHash.Hex())

	return nil
}

func counterQuery(privKey string, evmContract common.Address) (*big.Int, error) {
	var result *big.Int

	abiFile, err := ioutil.ReadFile(abiFilePath)
	if err != nil {
		log.Printf("Failed to read ABI file: %v", err)
		return nil, err
	}
	contractAbi, err := abi.JSON(bytes.NewReader(abiFile))
	if err != nil {
		log.Printf("Failed to parse ABI: %v", err)
		return nil, err
	}

	err = callContractStatic(client, &result, evmContract, &contractAbi, "count", []interface{}{}, nil)
	if err != nil {
		log.Printf("error: %+v", err)
	}

	fmt.Println(result)

	return result, err
}

func initKey(key string) (*ecdsa.PrivateKey, common.Address) {
	privateKey, err := crypto.HexToECDSA(key)
	if err != nil {
		panic("failed to switch unencrypted private key -> secp256k1 private key: " + err.Error())
	}
	pubkey := privateKey.Public()
	pubkeyECDSA, ok := pubkey.(*ecdsa.PublicKey)
	if !ok {
		panic("failed to switch secp256k1 private key -> pubkey")
	}
	senderAddress := crypto.PubkeyToAddress(*pubkeyECDSA)

	return privateKey, senderAddress
}
