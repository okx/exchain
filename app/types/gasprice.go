package types

import (
	"errors"
	"math/big"
	"sort"
)

const (
	sampleNumber = 3 // Number of transactions sampled in a block

	CongestionHigherGpMode = 0
	NormalGpMode           = 1
	MinimalGpMode          = 2

	NoGasUsedCap = -1
)

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

func NewSingleBlockGPs() *SingleBlockGPs {
	return &SingleBlockGPs{
		all:     make([]*big.Int, 0),
		sampled: make([]*big.Int, 0),
		gasUsed: 0,
	}
}

func (bgp SingleBlockGPs) GetAll() []*big.Int {
	return bgp.all
}

func (bgp SingleBlockGPs) GetSampled() []*big.Int {
	return bgp.sampled
}

func (bgp SingleBlockGPs) GetGasUsed() uint64 {
	return bgp.gasUsed
}

func (bgp *SingleBlockGPs) AddSampledGP(gp *big.Int) {
	gpCopy := new(big.Int).Set(gp)
	bgp.sampled = append(bgp.sampled, gpCopy)
}

func (bgp *SingleBlockGPs) Update(gp *big.Int, gas uint64) {
	gpCopy := new(big.Int).Set(gp)
	bgp.all = append(bgp.all, gpCopy)
	bgp.gasUsed += gas
}

func (bgp *SingleBlockGPs) Clear() {
	bgp.all = make([]*big.Int, 0)
	bgp.sampled = make([]*big.Int, 0)
	bgp.gasUsed = 0
}

func (bgp *SingleBlockGPs) Copy() *SingleBlockGPs {
	return &SingleBlockGPs{
		all:     bgp.all,
		sampled: bgp.sampled,
		gasUsed: bgp.gasUsed,
	}
}

func (bgp *SingleBlockGPs) SampleGP(adoptHigherGp bool) {
	// "len(bgp.sampled) != 0" means it has been sampled
	if len(bgp.sampled) != 0 {
		return
	}

	txGPs := make([]*big.Int, len(bgp.all))
	copy(txGPs, bgp.all)
	sort.Sort(BigIntArray(txGPs))

	if adoptHigherGp {

		rowSampledGPs := make([]*big.Int, 0)

		// Addition of sampleNumber lower-priced gp
		for i := 0; i < len(txGPs); i++ {
			if i >= sampleNumber {
				break
			}
			rowSampledGPs = append(rowSampledGPs, new(big.Int).Set(txGPs[i]))
		}

		// Addition of sampleNumber higher-priced gp
		for i := len(txGPs) - 1; i >= 0; i-- {
			if len(txGPs)-1-i >= sampleNumber {
				break
			}
			rowSampledGPs = append(rowSampledGPs, new(big.Int).Set(txGPs[i]))
		}

		if len(rowSampledGPs) != 0 {
			sampledGPLen := big.NewInt(int64(len(rowSampledGPs)))
			sum := big.NewInt(0)
			for _, gp := range rowSampledGPs {
				sum.Add(sum, gp)
			}

			avgGP := new(big.Int).Quo(sum, sampledGPLen)
			bgp.AddSampledGP(avgGP)
		}
	} else {
		for _, gp := range txGPs {
			bgp.AddSampledGP(gp)
			if len(bgp.sampled) >= sampleNumber {
				break
			}
		}
	}
}

// BlockGPResults is a circular queue of SingleBlockGPs
type BlockGPResults struct {
	items    []*SingleBlockGPs
	front    int
	rear     int
	capacity int
}

func NewBlockGPResults(checkBlocksNum int) *BlockGPResults {
	circularQueue := &BlockGPResults{
		items:    make([]*SingleBlockGPs, checkBlocksNum, checkBlocksNum),
		front:    -1,
		rear:     -1,
		capacity: checkBlocksNum,
	}
	return circularQueue
}

func (rs BlockGPResults) IsFull() bool {
	if rs.front == 0 && rs.rear == rs.capacity-1 {
		return true
	}
	if rs.front == rs.rear+1 {
		return true
	}
	return false
}

func (rs BlockGPResults) IsEmpty() bool {
	return rs.front == -1
}

func (rs BlockGPResults) Front() int {
	return rs.front
}

func (rs BlockGPResults) Rear() int {
	return rs.rear
}

func (rs BlockGPResults) Cap() int {
	return rs.capacity
}

func (rs *BlockGPResults) Push(gp *SingleBlockGPs) error {
	if rs.IsFull() {
		_, err := rs.Pop()
		if err != nil {
			return err
		}
	}
	if rs.front == -1 {
		rs.front = 0
	}
	rs.rear = (rs.rear + 1) % rs.capacity
	rs.items[rs.rear] = gp
	return nil
}

func (rs *BlockGPResults) Pop() (*SingleBlockGPs, error) {
	if rs.IsEmpty() {
		return nil, errors.New("pop failed: BlockGPResults is empty")
	}
	element := rs.items[rs.front]
	if rs.front == rs.rear {
		// rs has only one element,
		// so we reset the queue after deleting it
		rs.front = -1
		rs.rear = -1
	} else {
		rs.front = (rs.front + 1) % rs.capacity
	}
	return element, nil
}

func (rs *BlockGPResults) ExecuteSamplingBy(lastPrice *big.Int, adoptHigherGp bool) []*big.Int {
	var txPrices []*big.Int
	if !rs.IsEmpty() {
		// traverse the circular queue
		for i := rs.front; i != rs.rear; i = (i + 1) % rs.capacity {
			rs.items[i].SampleGP(adoptHigherGp)
			// If block is empty, use the latest gas price for sampling.
			if len(rs.items[i].sampled) == 0 {
				rs.items[i].AddSampledGP(lastPrice)
			}
			txPrices = append(txPrices, rs.items[i].sampled...)
		}
		rs.items[rs.rear].SampleGP(adoptHigherGp)
		if len(rs.items[rs.rear].sampled) == 0 {
			rs.items[rs.rear].AddSampledGP(lastPrice)
		}
		txPrices = append(txPrices, rs.items[rs.rear].sampled...)
	}
	return txPrices
}

type BigIntArray []*big.Int

func (s BigIntArray) Len() int           { return len(s) }
func (s BigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s BigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
