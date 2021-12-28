package main

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

var privKey = "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17"

func main() {
	client, err := ethclient.Dial(RpcUrl)
	if err != nil {
		log.Fatalf("failed to initialize client: %+v", err)
	}
	toAddress := common.HexToAddress("0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0")

	send(client, toAddress)

	deployStandardOIP20ContractAndTransfer(client, toAddress)
}

func send(client *ethclient.Client, toAddress common.Address) {
	privateKey, senderAddress := initKey(privKey)

	// send 0.001okt
	transferOKT(client, senderAddress, toAddress, str2bigInt("0.001"), privateKey, 0)
}

func deployStandardOIP20ContractAndTransfer(client *ethclient.Client, toAddress common.Address) {
	privateKey, pubkey := initKey(privKey)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(int64(ChainId)))
	auth.Nonce = big.NewInt(int64(0))
	auth.Value = big.NewInt(0)       // in wei
	auth.GasLimit = uint64(GasLimit) // in units
	auth.GasPrice = big.NewInt(int64(GasPrice))

	standardContractAddress := deployStandardOIP20Contract(client, auth, "oip", "oip std", 18, big.NewInt(int64(100000000000000000000)), pubkey)

	oip20, err := NewOIP20(standardContractAddress, client)
	if err != nil {
		return
	}

	_, err = oip20.Transfer(auth, toAddress, big.NewInt(int64(10000000000000000000)))
	if err != nil {
		log.Fatal(err)
	}
}
