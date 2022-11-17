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
	Blocks  int
}

func NewGPOConfig(weight int, checkBlocks int) GPOConfig {
	return GPOConfig{
		Weight:  weight,
		Default: defaultPrice,
		Blocks:  checkBlocks,
	}
}

func DefaultGPOConfig() GPOConfig {
	return NewGPOConfig(80, 5)
}

// Oracle recommends gas prices based on the content of recent blocks.
type Oracle struct {
	CurrentBlockGPs *types.SingleBlockGPs
	// hold the gas prices of the latest few blocks
	BlockGPQueue *types.BlockGPResults
	lastPrice    *big.Int
	weight       int
}

// NewOracle returns a new gasprice oracle which can recommend suitable
// gasprice for newly created transaction.
func NewOracle(params GPOConfig) *Oracle {
	cbgp := types.NewSingleBlockGPs()
	bgpq := types.NewBlockGPResults(params.Blocks)
	weight := params.Weight
	if weight < 0 {
		weight = 0
	}
	if weight > 100 {
		weight = 100
	}
	return &Oracle{
		CurrentBlockGPs: cbgp,
		BlockGPQueue:    bgpq,
		lastPrice:       params.Default,
		weight:          weight,
	}
}

func (gpo *Oracle) RecommendGP() *big.Int {
	maxGasUsed := appconfig.GetOecConfig().GetMaxGasUsedPerBlock()
	maxTxNum := appconfig.GetOecConfig().GetMaxTxNumPerBlock()
	allTxsLen := int64(len(gpo.CurrentBlockGPs.GetAll()))
	// If maxGasUsed is not negative and the current block's total gas consumption is more than 80% of it,
	// or the number of tx in the current block is equal to MaxTxNumPerBlock in mempool config,
	// then we consider the chain to be congested
	// and return recommend gas price.
	if (maxGasUsed > 0 && gpo.CurrentBlockGPs.GetGasUsed() >= uint64(maxGasUsed*80/100)) || allTxsLen == maxTxNum {
		txPrices := gpo.BlockGPQueue.ExecuteSamplingBy(gpo.lastPrice)

		price := gpo.lastPrice
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
	return defaultPrice
}
