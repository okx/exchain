package main

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"math/big"
	"time"
)

const (
	abiFile = "./contracts/counter/counter.abi"
	binFile = "./contracts/counter/counter.bin"
)

func main() {
	privKey := []string {
		"8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17",
		"171786c73f805d257ceb07206d851eea30b3b41a2170ae55e1225e0ad516ef42",
		"b7700998b973a2cae0cb8e8a328171399c043e57289735aca5f2419bd622297a",
		"00dcf944648491b3a822d40bf212f359f699ed0dd5ce5a60f1da5e1142855949",
	}

	for _, key := range privKey {
		go writeRoutine(key, time.Millisecond*50)
	}

	<-make(chan struct{})
}

func writeRoutine(privKey string, blockTime time.Duration) {
	var (
		privateKey             *ecdsa.PrivateKey
		senderAddress          common.Address
	)
	privateKey, senderAddress = initKey(privKey)

	defer func() {
		if r := recover(); r != nil {
			sleep(3)
			go writeRoutine(privKey, blockTime)
		}
	}()

	client, err := ethclient.Dial(RpcUrl)

	contract := newContract("counter", "", abiFile, binFile)
	err = deployContract(client, senderAddress, privateKey, contract, 3)

	for err == nil {
		err = writeContract(client, contract, senderAddress, privateKey, nil, blockTime, "add", big.NewInt(100))
		uint256Output(client, contract, "getCounter")
		err = writeContract(client, contract, senderAddress, privateKey, nil, blockTime, "subtract",)
		uint256Output(client, contract, "getCounter")
	}
	panic(err)
}
