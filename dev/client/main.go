package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type TestType string

const (
	abiFile = "./contracts/counter/counter.abi"
	binFile = "./contracts/counter/counter.bin"

	Oip20Test            = TestType("oip20")
	SingleEthereumTxTest = TestType("single-eth-tx")
	CounterTest          = TestType("counter")
)

func main() {
	testTypeParam := flag.String("type", "oip20", "choose which test to run")
	testRpcUrl := flag.String("type", "http://127.0.0.1:8545", "default test rpc url ")
	RpcUrl = *testRpcUrl
	flag.Parse()

	privKey := []string{
		"8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17",
		//"171786c73f805d257ceb07206d851eea30b3b41a2170ae55e1225e0ad516ef42",
		//"b7700998b973a2cae0cb8e8a328171399c043e57289735aca5f2419bd622297a",
		//"00dcf944648491b3a822d40bf212f359f699ed0dd5ce5a60f1da5e1142855949",
	}

	var testFunc func(privKey string, blockTime time.Duration) error
	switch TestType(*testTypeParam) {
	case Oip20Test:
		fmt.Printf("contract: %s\n", *testTypeParam)
		testFunc = standardOip20Test
		break
	case SingleEthereumTxTest:
		fmt.Printf("Single transfer tx")
		testFunc = standardSingleTransferTxTest
		break
	default:
		fmt.Printf("contract: %s\n", CounterTest)
		testFunc = counterTest
	}

	for _, k := range privKey {
		test := func(key string) {
			testFunc(key, time.Millisecond*5000)
		}
		go writeRoutine(test, k)
	}
	<-make(chan struct{})
}

func writeRoutine(test func(string), key string) {
	for {
		test(key)
		log.Printf("recover writeRoutine...")
		sleep(3)
	}
}

func counterTest(privKey string, blockTime time.Duration) error {
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
	return err
}

func standardSingleTransferTxTest(privKey string, blockTime time.Duration) error {
	privateKey, sender := initKey(privKey)
	client, err := ethclient.Dial(RpcUrl)
	if err != nil {
		log.Printf("failed to dial: %+v", err)
	}
	toAddress := common.HexToAddress("0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0")
	return transferOKT(client, sender, toAddress, str2bigInt("0.001"), privateKey, 0)
}

func standardOip20Test(privKey string, blockTime time.Duration) error {
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

	return err
}
