package global

import "github.com/ethereum/go-ethereum/accounts/abi/bind"

var client bind.ContractBackend

func SetEthClient(client bind.ContractBackend) {
	client = client
}

func GetEthClient() bind.ContractBackend {
	return client
}
