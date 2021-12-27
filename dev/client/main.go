package main

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

var privKey = "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17"

func main() {
	client, err := ethclient.Dial(RpcUrl)
	if err != nil {
		log.Fatalf("failed to initialize client: %+v", err)
	}
	send(client, "0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0")

}

func send(client *ethclient.Client, to string) {
	privateKey, senderAddress := initKey(privKey)
	toAddress := common.HexToAddress(to)

	// send 0.001okt
	transferOKT(client, senderAddress, toAddress, str2bigInt("0.001"), privateKey, 0)
}
