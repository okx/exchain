package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

type TestType string

const (
	abiFile = "./contracts/counter/counter.abi"
	binFile = "./contracts/counter/counter.bin"

	Oip20Test   = TestType("oip20")
	CounterTest = TestType("counter")
)

func main() {
	testTypeParam := flag.String("type", "oip20", "choose which test to run")
	privKey := []string{
		"8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17",
		"171786c73f805d257ceb07206d851eea30b3b41a2170ae55e1225e0ad516ef42",
		"b7700998b973a2cae0cb8e8a328171399c043e57289735aca5f2419bd622297a",
		"00dcf944648491b3a822d40bf212f359f699ed0dd5ce5a60f1da5e1142855949",
	}

	var testFunc func(privKey string, blockTime time.Duration)
	switch TestType(*testTypeParam) {
	case Oip20Test:
		testFunc = standardOip20Test
		break
	default:
		testFunc = writeRoutine
	}

	for _, key := range privKey {
		go testFunc(key, time.Millisecond*50)
	}
	<-make(chan struct{})
}

func writeRoutine(privKey string, blockTime time.Duration) {
	var (
		privateKey    *ecdsa.PrivateKey
		senderAddress common.Address
	)

	privateKey, senderAddress = initKey(privKey)
	counterContract := newContract("counter", "", abiFile, binFile)

	client, err := ethclient.Dial(RpcUrl)
	if err == nil {
		err = deployContract(client, senderAddress, privateKey, counterContract, 3)
	}

	for err == nil {
		err = writeContract(client, counterContract, senderAddress, privateKey, nil, blockTime, "add", big.NewInt(100))
		uint256Output(client, counterContract, "getCounter")
		err = writeContract(client, counterContract, senderAddress, privateKey, nil, blockTime, "subtract")
		uint256Output(client, counterContract, "getCounter")
	}

	sleep(3)
	log.Printf("recover writeRoutine")
	go writeRoutine(privKey, blockTime)
}

func standardOip20Test(privKey string, blockTime time.Duration) {
	privateKey, sender := initKey(privKey)

	client, err := ethclient.Dial(RpcUrl)
	if err != nil {
		log.Printf("failed to dial: %+v", err)
	}

	oip20, auth, err := deployOip(client, sender, privateKey)
	if err != nil {
		log.Printf("failed to deploy: %+v", err)
	}

	toAddress := common.HexToAddress("0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0")
	for err == nil {
		nonce, err := transferOip(client, oip20, sender, auth, toAddress)
		if err != nil {
			log.Printf("failed to transfer Oip: %+v", err)
			break
		}
		fmt.Printf(
			"==================================================\n"+
				"Standard OIP20 transfer:\n"+
				"	from					: <%s>\n"+
				"	nonce					: <%d>\n"+
				"	to					: <%s>\n",
			sender, nonce, toAddress,
		)
		time.Sleep(blockTime)
	}

	log.Printf("recover standardOip20Test")
	sleep(3)
	go standardOip20Test(privKey, blockTime)
}
