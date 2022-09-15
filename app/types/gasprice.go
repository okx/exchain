package types

import (
	"errors"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/params"
)

const sampleNumber = 3 // Number of transactions sampled in a block

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

func (bgp *SingleBlockGPs) AddSampledGP(gp *big.Int) {
	bgp.sampled = append(bgp.sampled, gp)
}

func (bgp *SingleBlockGPs) Update(gp *big.Int, gas uint64) {
	bgp.all = append(bgp.all, gp)
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

// BlockGPResults is a circular queue of SingleBlockGPs
type BlockGPResults struct {
	Items    []SingleBlockGPs
	front    int
	rear     int
	capacity int
}

func NewBlockGPResults(checkBlocksNum int) BlockGPResults {
	circularQueue := BlockGPResults{
		Items:    make([]SingleBlockGPs, checkBlocksNum, checkBlocksNum),
		front:    -1,
		rear:     -1,
		capacity: checkBlocksNum,
	}
	return circularQueue
}

func (rs *BlockGPResults) IsFull() bool {
	if rs.front == 0 && rs.rear == rs.capacity-1 {
		return true
	}
	if rs.front == rs.rear+1 {
		return true
	}
	return false
}

func (rs *BlockGPResults) IsEmpty() bool {
	return rs.front == -1
}

func (rs *BlockGPResults) Front() int {
	return rs.front
}

func (rs *BlockGPResults) Rear() int {
	return rs.rear
}

func (rs *BlockGPResults) Cap() int {
	return rs.capacity
}

func (rs *BlockGPResults) Push(gp SingleBlockGPs) error {
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
	rs.Items[rs.rear] = gp
	return nil
}

func (rs *BlockGPResults) Pop() (*SingleBlockGPs, error) {
	if rs.IsEmpty() {
		return nil, errors.New("pop failed: BlockGPResults is empty")
	}
	element := rs.Items[rs.front]
	if rs.front == rs.rear {
		// rs has only one element,
		// so we reset the queue after deleting it
		rs.front = -1
		rs.rear = -1
	} else {
		rs.front = (rs.front + 1) % rs.capacity
	}
	return &element, nil
}

type BigIntArray []*big.Int

func (s BigIntArray) Len() int           { return len(s) }
func (s BigIntArray) Less(i, j int) bool { return s[i].Cmp(s[j]) < 0 }
func (s BigIntArray) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
