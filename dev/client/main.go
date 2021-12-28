package main

import (
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

var privKey = "8ff3ca2d9985c3a52b459e2f6e7822b23e1af845961e22128d5f372fb9aa5f17"


var (
	privateKey             *ecdsa.PrivateKey
	senderAddress          common.Address
)

const (
	abiFile = "./contracts/counter/counter.abi"
	binFile = "./contracts/counter/counter.bin"
)
func init() {
	privateKey, senderAddress = initKey(privKey)

}

func main() {
	client, err := ethclient.Dial(RpcUrl)
	if err != nil {
		log.Fatalf("failed to initialize client: %+v", err)
	}
	//send(client, "0x83D83497431C2D3FEab296a9fba4e5FaDD2f7eD0")

	contract := newContract("counter", "", abiFile, binFile)

	// 1. deploy contract
	deployContract(client, senderAddress, privateKey, contract)

	//contractAddr := common.HexToAddress("0x79BE5cc37B7e17594028BbF5d43875FDbed417db")

	//contract := NewContract("", "0x9a59ae3Fc0948717F94242fc170ac1d5dB3f0D5D", abiFile)

	// 2. call contract(write)
	uint256Output(client, contract, "getCounter")
	writeContract(client, contract, senderAddress, privateKey, nil, 3, "add", big.NewInt(100))
	uint256Output(client, contract, "getCounter")
	writeContract(client, contract, senderAddress, privateKey, nil, 3, "subtract",)
	uint256Output(client, contract, "getCounter")
}





func send(client *ethclient.Client, to string) {
	privateKey, senderAddress := initKey(privKey)
	toAddress := common.HexToAddress(to)

	// send 0.001okt
	transferOKT(client, senderAddress, toAddress, str2bigInt("0.001"), privateKey, 0)
}
