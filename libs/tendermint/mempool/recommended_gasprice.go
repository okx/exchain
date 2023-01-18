package mempool

import (
	"math/big"
	"sort"

	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/libs/tendermint/types"
)

type GPOConfig struct {
	Weight  int
	Default *big.Int
	Blocks  int
}

func NewGPOConfig(weight int, checkBlocks int) GPOConfig {
	return GPOConfig{
		Weight: weight,
		// default gas price is 0.1GWei
		Default: big.NewInt(100000000),
		Blocks:  checkBlocks,
	}
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
		// Note: deep copy is necessary here
		lastPrice: new(big.Int).Set(params.Default),
		weight:    weight,
	}
}

func (gpo *Oracle) RecommendGP() *big.Int {

	minGP := global.GetGlobalMinGasPrice()
	maxGP := global.GetGlobalMaxGasPrice()

	maxGasUsed := cfg.DynamicConfig.GetDynamicGpMaxGasUsed()
	maxTxNum := cfg.DynamicConfig.GetDynamicGpMaxTxNum()

	allTxsLen := int64(len(gpo.CurrentBlockGPs.GetAll()))
	// If the current block's total gas consumption is more than maxGasUsed,
	// or the number of tx in the current block is more than maxTxNum,
	// then we consider the chain to be congested.
	isCongested := (int64(gpo.CurrentBlockGPs.GetGasUsed()) >= maxGasUsed) || (allTxsLen >= maxTxNum)

	// When the network is congested, increase the recommended gas price.
	adoptHigherGp := (cfg.DynamicConfig.GetDynamicGpMode() == types.CongestionHigherGpMode) && isCongested

	txPrices := gpo.BlockGPQueue.ExecuteSamplingBy(gpo.lastPrice, adoptHigherGp)

	price := new(big.Int).Set(gpo.lastPrice)

	// If the block is not congested, return the minimal GP
	if !isCongested {
		price.Set(minGP)
		gpo.lastPrice.Set(price)
		return price
	}

	if len(txPrices) > 0 {
		sort.Sort(types.BigIntArray(txPrices))
		price.Set(txPrices[(len(txPrices)-1)*gpo.weight/100])
	}

	if price.Cmp(maxGP) > 0 {
		price.Set(maxGP)
	}
	gpo.lastPrice.Set(price)
	return price
}
