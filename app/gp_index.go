package app

import (
	"math"
	"math/big"
	"sort"
)

type GasPriceIndex struct {
	Min *big.Int `json:"min"`
	Q1  *big.Int `json:"q1"`
	Q2  *big.Int `json:"q2"`
	Q3  *big.Int `json:"q3"`
	Max *big.Int `json:"max"`
}

func CalBlockGasPriceIndex(blockGasPrice []*big.Int) GasPriceIndex {
	num := len(blockGasPrice)
	if num == 0 {
		return GasPriceIndex{}
	}

	sort.SliceStable(blockGasPrice, func(i, j int) bool {
		return blockGasPrice[i].Cmp(blockGasPrice[j]) < 0
	})

	gpIndex := GasPriceIndex{
		Min: blockGasPrice[0],
		Q1:  blockGasPrice[0],
		Q2:  blockGasPrice[0],
		Q3:  blockGasPrice[0],
		Max: blockGasPrice[num-1],
	}

	if q1 := int(math.Round(float64(num) * 0.25)); q1 > 0 {
		gpIndex.Q1 = blockGasPrice[q1-1]
	}

	if q2 := int(math.Round(float64(num) * 0.5)); q2 > 0 {
		gpIndex.Q2 = blockGasPrice[q2-1]
	}

	if q3 := int(math.Round(float64(num) * 0.75)); q3 > 0 {
		gpIndex.Q3 = blockGasPrice[q3-1]
	}

	return gpIndex
}
