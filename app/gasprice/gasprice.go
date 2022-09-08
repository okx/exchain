package gasprice

import (
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/params"

	appconfig "github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/types"
)

var (
	maxPrice     = big.NewInt(500 * params.GWei)
	defaultPrice = big.NewInt(params.GWei / 10)
)

type GPOConfig struct {
	Weight  int
	Default *big.Int `toml:",omitempty"`
}

func NewGPOConfig(weight int, defaultPrice *big.Int) GPOConfig {
	return GPOConfig{
		Weight:  weight,
		Default: defaultPrice,
	}
}

func DefaultGPOConfig() GPOConfig {
	return NewGPOConfig(appconfig.GetOecConfig().GetDynamicGpWeight(), defaultPrice)
}

// Oracle recommends gas prices based on the content of recent blocks.
type Oracle struct {
	CurrentBlockGPs types.SingleBlockGPs
	// hold the gas prices of the latest few blocks
	BlockGPQueue types.BlockGPResults
	lastPrice    *big.Int
	weight       int
}

// NewOracle returns a new gasprice oracle which can recommend suitable
// gasprice for newly created transaction.
func NewOracle(params GPOConfig) *Oracle {
	weight := params.Weight
	if weight < 0 {
		weight = 0
	}
	if weight > 100 {
		weight = 100
	}
	return &Oracle{
		lastPrice: params.Default,
		weight:    weight,
	}
}

func (gpo *Oracle) RecommendGP() *big.Int {
	maxGasUsed := appconfig.GetOecConfig().GetMaxGasUsedPerBlock()
	if maxGasUsed > -1 && gpo.CurrentBlockGPs.GetGasUsed() < uint64(maxGasUsed) {
		return defaultPrice
	}
	// If the number of tx in the current block is less than the MaxTxNumPerBlock in mempool config,
	// the default gas price is returned.
	allGPsLen := int64(len(gpo.CurrentBlockGPs.GetAll()))
	maxTxNum := appconfig.GetOecConfig().GetMaxTxNumPerBlock()
	if allGPsLen < maxTxNum {
		return defaultPrice
	}

	lastPrice := gpo.lastPrice

	var txPrices []*big.Int
	for i := 0; i < len(gpo.BlockGPQueue); i++ {
		gpo.BlockGPQueue[i].SampleGP()

		// If block is empty, use the latest calculated price for sampling.
		if len(gpo.BlockGPQueue[i].GetSampled()) == 0 {
			gpo.BlockGPQueue[i].AddSampledGP(lastPrice)
		}

		txPrices = append(txPrices, gpo.BlockGPQueue[i].GetSampled()...)
	}
	price := lastPrice
	if len(txPrices) > 0 {
		sort.Sort(types.BigIntArray(txPrices))
		price = txPrices[(len(txPrices)-1)*gpo.weight/100]
	}
	if price.Cmp(maxPrice) > 0 {
		price = new(big.Int).Set(maxPrice)
	}
	gpo.lastPrice = price
	return price
}
