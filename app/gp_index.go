package app

import (
	"math"
	"math/big"
	"sort"
)

type GasPriceIndex struct {
	RecommendGp *big.Int `json:"recommend-gp"`
}

func CalBlockGasPriceIndex(blockGasPrice []*big.Int, weight int) GasPriceIndex {
	num := len(blockGasPrice)
	if num == 0 {
		return GasPriceIndex{}
	}

	sort.SliceStable(blockGasPrice, func(i, j int) bool {
		return blockGasPrice[i].Cmp(blockGasPrice[j]) < 0
	})

	idx := int(math.Round(float64(weight) / 100.0 * float64(num)))
	if idx > 0 {
		idx -= 1
	}

	return GasPriceIndex{
		RecommendGp: blockGasPrice[idx],
	}
}
