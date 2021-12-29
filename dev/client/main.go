package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
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

	switch TestType(*testTypeParam) {

	case Oip20Test:
		for _, key := range privKey {
			go standardOip20Test(key, time.Millisecond*50)
		}
		break
	case CounterTest:
		for _, key := range privKey {
			go writeRoutine(key, time.Millisecond*50)
		}
		break
	}

	<-make(chan struct{})
}

func writeRoutine(privKey string, blockTime time.Duration) {
	var (
		privateKey    *ecdsa.PrivateKey
		senderAddress common.Address
	)

	defer func() {
		if r := recover(); r != nil {
			sleep(3)
			go writeRoutine(privKey, blockTime)
		}
	}()
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
	panic(err)
}

func standardOip20Test(privKey string, blockTime time.Duration) {
	toAddress := common.HexToAddress("0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0")
	privateKey, pubkey := initKey(privKey)

	defer func() {
		if r := recover(); r != nil {
			sleep(3)
			go standardOip20Test(privKey, blockTime)
		}
	}()

	client, err := ethclient.Dial(RpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), pubkey)
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(ChainId))
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(auth.Context)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = GasLimit   // in units
	auth.GasPrice = gasPrice
	auth.Context = context.Background()

	symbol := "OIP20"
	contractName := "OIP20 STD"
	decimals := 18

	var (
		oip20 *Oip20
	)

	if err == nil {
		_, oip20, err = deployStandardOIP20Contract(client, auth, symbol, contractName, uint8(decimals), str2bigInt("100000000000000000000000"), pubkey, blockTime)
	}

	for err == nil {

		nonce++
		auth.Nonce = big.NewInt(int64(nonce))

		transferAmount := str2bigInt("100000000000000000")
		_, err = oip20.Transfer(auth, toAddress, transferAmount)

		fmt.Printf(
			"==================================================\n"+
				"Standard OIP20 transfer:\n"+
				"	contract name				: <%s>\n"+
				"	from					: <%s>\n"+
				"	to					: <%s>\n"+
				"	amount					: <%s>\n"+
				"==================================================\n",
			contractName,
			pubkey,
			toAddress,
			transferAmount,
		)
	}
}
