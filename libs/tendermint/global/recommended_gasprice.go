package global

import (
	"math/big"
)

var globalMinGasPrice *big.Int
var globalMaxGasPrice *big.Int

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
