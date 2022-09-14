package types

import (
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/params"
)

const (
	checkBlocks  = 5
	sampleNumber = 3 // Number of transactions sampled in a block
)

var ignorePrice = big.NewInt(2 * params.Wei)

// SingleBlockGPs holds the gas price of all transactions in a block
// and will sample the lower few gas prices according to sampleNumber.
type SingleBlockGPs struct {
	// gas price of all transactions
	all []*big.Int
	// gas price of transactions sampled in a block
	sampled []*big.Int
	// total gas of all tx in the block
	gasUsed uint64
}

func NewSingleBlockGPs() SingleBlockGPs {
	return SingleBlockGPs{
		all:     make([]*big.Int, 0),
		sampled: make([]*big.Int, 0),
		gasUsed: 0,
	}
}

func (bgp *SingleBlockGPs) GetAll() []*big.Int {
	return bgp.all
}

func (bgp *SingleBlockGPs) GetSampled() []*big.Int {
	return bgp.sampled
}

func (bgp *SingleBlockGPs) GetGasUsed() uint64 {
	return bgp.gasUsed
}

func (bgp *SingleBlockGPs) AddGP(gp *big.Int) {
	bgp.all = append(bgp.all, gp)
}

func (bgp *SingleBlockGPs) AddSampledGP(gp *big.Int) {
	bgp.sampled = append(bgp.sampled, gp)
}

func (bgp *SingleBlockGPs) AddGas(gas uint64) {
	bgp.gasUsed += gas
}

func (bgp *SingleBlockGPs) Clear() {
	bgp.all = bgp.all[:0]
	bgp.sampled = bgp.sampled[:0]
	bgp.gasUsed = 0
}

func (bgp *SingleBlockGPs) SampleGP() {
	// "len(bgp.sampled) != 0" means it has been sampled
	if len(bgp.all) == 0 && len(bgp.sampled) != 0 {
		return
	}

	txGPs := make([]*big.Int, len(bgp.all))
	copy(txGPs, bgp.all)
	sort.Sort(BigIntArray(txGPs))

	for _, gp := range txGPs {
		// If a GP is too cheap, discard it.
		if gp.Cmp(ignorePrice) == -1 {
			continue
		}
		bgp.AddSampledGP(gp)
		if len(bgp.sampled) >= sampleNumber {
			break
		}
	}
}

// BlockGPResults holds the gas prices of the latest few blocks
type BlockGPResults []SingleBlockGPs

func NewBlockGPResults() BlockGPResults {
	return BlockGPResults{}
}

func (rs *BlockGPResults) Push(gp SingleBlockGPs) {
	if rs.isFull() {
		_, _ = rs.Pop()
	}
	*rs = append(*rs, gp)
}

func (rs *BlockGPResults) Pop() (SingleBlockGPs, error) {
	if rs.isEmpty() {
		return SingleBlockGPs{}, errors.New("pop failed: slice is empty")
	}
	res := (*rs)[0]
	*rs = (*rs)[1:]
	return res, nil
}

func (rs BlockGPResults) isEmpty() bool {
	return len(rs) == 0
}

func (rs BlockGPResults) isFull() bool {
	return len(rs) == checkBlocks
}

type BigIntArray []*big.Int

func (s BigIntArray) Len() int           { return len(s) }
func (s BigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s BigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
