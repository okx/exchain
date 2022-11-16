package global

import (
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type LocalEthClient interface {
	bind.ContractBackend
	ethereum.ChainStateReader
}

var client LocalEthClient

func SetLocalEthClient(backend LocalEthClient) {
	client = backend
}

func GetLocalEthClient() LocalEthClient {
	return client
}
