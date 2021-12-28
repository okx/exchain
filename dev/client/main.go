package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

var privKey = "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17"

func main() {
	client, err := ethclient.Dial(RpcUrl)
	if err != nil {
		log.Fatalf("failed to initialize client: %+v", err)
	}
	toAddress := common.HexToAddress("0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0")

	send(client, toAddress)

	deployStandardOIP20ContractAndTransfer(client, toAddress, time.Second*1)
}

func send(client *ethclient.Client, toAddress common.Address) {
	privateKey, senderAddress := initKey(privKey)

	// send 0.001okt
	transferOKT(client, senderAddress, toAddress, str2bigInt("0.001"), privateKey, 0)
}

func deployStandardOIP20ContractAndTransfer(client *ethclient.Client, toAddress common.Address, waitDuration time.Duration) {
	privateKey, pubkey := initKey(privKey)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(ChainId))
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(auth.Context)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), pubkey)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = uint64(GasLimit) // in units
	auth.GasPrice = gasPrice
	auth.Context = context.Background()

	symbol := "OIP20"
	contractName := "OIP20 STD"
	decimals := 18

	_, oip20, err := deployStandardOIP20Contract(client, auth, symbol, contractName, uint8(decimals), str2bigInt("1000000000000000000000"), pubkey, waitDuration)

	if err != nil {
		log.Fatal(err)
	}

	balanceBefore, err := oip20.BalanceOf(&bind.CallOpts{Pending: false, Context: context.Background()}, toAddress)
	if err != nil {
		log.Fatal(err)
	}

	transferAmount := str2bigInt("10000000000000000000")
	nonce, err = client.PendingNonceAt(context.Background(), pubkey)
	if err != nil {
		log.Fatal(err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	_, err = oip20.Transfer(auth, toAddress, transferAmount)
	if err != nil {
		log.Fatal(err)
	}
	// waiting for token transfer successfully
	time.Sleep(waitDuration)

	if err != nil {
		log.Fatal(err)
	}

	balanceAfter, err := oip20.BalanceOf(&bind.CallOpts{Pending: false, Context: context.Background()}, toAddress)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"==================================================\n"+
			"Standard OIP20 transfer:\n"+
			"	contract name				: <%s>\n"+
			"	from					: <%s>\n"+
			"	to					: <%s>\n"+
			"	amount					: <%s>\n"+
			"	received balance before			: <%s>\n"+
			"	received balance after			: <%s>\n"+
			"==================================================\n",
		contractName,
		pubkey,
		toAddress,
		transferAmount,
		balanceBefore,
		balanceAfter,
	)
}
