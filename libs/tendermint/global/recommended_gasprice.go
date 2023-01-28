package global

import (
	"math/big"
)

// initialization
var globalMinGasPrice = big.NewInt(100000000)
var globalMaxGasPrice = big.NewInt(500000000000)

func SetGlobalMinAndMaxGasPrice(minGP *big.Int) {
	globalMinGasPrice = new(big.Int).Set(minGP)
	globalMaxGasPrice = new(big.Int).Mul(globalMinGasPrice, big.NewInt(5000))
}

func GetGlobalMinGasPrice() *big.Int {
	return globalMinGasPrice
}

func GetGlobalMaxGasPrice() *big.Int {
	return globalMaxGasPrice
}
