package mempool

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/okx/okbchain/libs/system/trace"
	cfg "github.com/okx/okbchain/libs/tendermint/config"
	"github.com/okx/okbchain/libs/tendermint/types"
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

func (gpo *Oracle) RecommendGP() (*big.Int, bool) {

	gasUsedThreshold := cfg.DynamicConfig.GetDynamicGpMaxGasUsed()
	txNumThreshold := cfg.DynamicConfig.GetDynamicGpMaxTxNum()

	allTxsLen := int64(len(gpo.CurrentBlockGPs.GetAll()))
	// If the current block's total gas consumption is more than gasUsedThreshold,
	// or the number of tx in the current block is more than txNumThreshold,
	// then we consider the chain to be congested.
	isCongested := (int64(gpo.CurrentBlockGPs.GetGasUsed()) >= gasUsedThreshold) || (allTxsLen >= txNumThreshold)
	trace.GetElapsedInfo().AddInfo(trace.IsCongested, fmt.Sprintf("%t", isCongested))

	// If GP mode is CongestionHigherGpMode, increase the recommended gas price when the network is congested.
	adoptHigherGp := (cfg.DynamicConfig.GetDynamicGpMode() == types.CongestionHigherGpMode) && isCongested

	txPrices := gpo.BlockGPQueue.ExecuteSamplingBy(gpo.lastPrice, adoptHigherGp)

	price := new(big.Int).Set(gpo.lastPrice)

	if len(txPrices) > 0 {
		sort.Sort(types.BigIntArray(txPrices))
		price.Set(txPrices[(len(txPrices)-1)*gpo.weight/100])
	}
	gpo.lastPrice.Set(price)

	trace.GetElapsedInfo().AddInfo(trace.RecommendedGP, fmt.Sprintf("%sWei", price.String()))
	return price, isCongested
}
