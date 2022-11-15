package global

import "github.com/ethereum/go-ethereum/accounts/abi/bind"

var client bind.ContractBackend

func SetEthClient(backend bind.ContractBackend) {
	client = backend
}

func GetEthClient() bind.ContractBackend {
	return client
}
