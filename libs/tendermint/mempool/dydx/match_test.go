package dydx

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var config = DydxConfig{
	PrivKeyHex:                 "e47a1fe74a7f9bfa44a362a3c6fbe96667242f62e6b8e138b3f61bd431c3215d",
	ChainID:                    "65",
	EthWsRpcUrl:                "wss://exchaintestws.okex.org:8443",
	PerpetualV1ContractAddress: "0xaC405bA85723d3E8d6D87B3B36Fd8D0D4e32D2c9",
	P1OrdersContractAddress:    "0xf1730217Bd65f86D2F008f1821D8Ca9A26d64619",
	P1MakerOracleAddress:       "0x4241DD684fbC5bCFCD2cA7B90b72885A79cf50B4",
}

type testTool struct {
	*testing.T
}

func TestMatch(t *testing.T) {
	tool := &testTool{T: t}

	book := NewDepthBook()
	me, err := NewMatchEngine(book, config, nil)
	require.NoError(t, err)

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(100, 100, true), nil))
	// buy
	// 100, 100
	// ------
	// sell

	// taker filled 100,10
	// maker filled
	// 100,10
	mr, err := me.Match(newTestOrder(100, 10, false), nil)
	// buy
	// 100, 90
	// ------
	// sell
	require.NoError(t, err)
	require.Equal(t, 1, len(mr.MatchedRecords))
	require.Equal(t, "100", mr.MatchedRecords[0].Fill.Price.String())
	require.Equal(t, "10", mr.MatchedRecords[0].Fill.Amount.String())

	require.Equal(t, "0", mr.TakerOrder.LeftAmount.String())
	require.Equal(t, "10", mr.TakerOrder.FrozenAmount.String())
	require.Equal(t, "90", mr.MatchedRecords[0].Maker.LeftAmount.String())

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(120, 5, false), nil))
	// buy
	// 100, 90
	// ------
	// 120, 5
	// sell

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(110, 6, false), nil))
	// buy
	// 100, 90
	// ------
	// 110, 6
	// 120, 5
	// sell

	// taker filled 121,10
	// maker filled
	// 110, 6
	// 120, 5
	mr, err = me.Match(newTestOrder(121, 10, true), nil)
	// buy
	// 100, 90
	// ------
	// 120, 1
	// sell
	require.NoError(t, err)
	require.Equal(t, 2, len(mr.MatchedRecords))
	require.Equal(t, "110", mr.MatchedRecords[0].Fill.Price.String())
	require.Equal(t, "6", mr.MatchedRecords[0].Fill.Amount.String())
	require.Equal(t, "0", mr.MatchedRecords[0].Maker.LeftAmount.String())

	require.Equal(t, "120", mr.MatchedRecords[1].Fill.Price.String())
	require.Equal(t, "4", mr.MatchedRecords[1].Fill.Amount.String())
	require.Equal(t, "1", mr.MatchedRecords[1].Maker.LeftAmount.String())

	require.Equal(t, "0", mr.TakerOrder.LeftAmount.String())
	require.Equal(t, "10", mr.TakerOrder.FrozenAmount.String())

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(99, 10, true), nil))
	// buy
	// 99, 10
	// 100, 90
	// ------
	// 120, 1
	// sell

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(97, 100, true), nil))
	// buy
	// 97, 100
	// 99, 10
	// 100, 90
	// ------
	// 120, 1
	// sell

	// no match
	tool.requireNoMatch(me.Match(newTestOrder(98, 15, true), nil))
	// buy
	// 97, 100
	// 98, 15
	// 99, 10
	// 100, 90
	// ------
	// 120, 1
	// sell

	// taker filled 98,115
	// maker filled
	// 100, 90
	// 99, 10
	// 98, 15
	mr, err = me.Match(newTestOrder(98, 130, false), nil)
	// buy
	// 97, 100
	// ------
	// 98, 15
	// 120, 1
	// sell
	require.NoError(t, err)
	require.Equal(t, 3, len(mr.MatchedRecords))
	require.Equal(t, "100", mr.MatchedRecords[0].Fill.Price.String())
	require.Equal(t, "90", mr.MatchedRecords[0].Fill.Amount.String())

	require.Equal(t, "99", mr.MatchedRecords[1].Fill.Price.String())
	require.Equal(t, "10", mr.MatchedRecords[1].Fill.Amount.String())

	require.Equal(t, "98", mr.MatchedRecords[2].Fill.Price.String())
	require.Equal(t, "15", mr.MatchedRecords[2].Fill.Amount.String())

	require.Equal(t, "15", mr.TakerOrder.LeftAmount.String())
	require.Equal(t, "115", mr.TakerOrder.FrozenAmount.String())
}

func (tool *testTool) requireNoMatch(mr *MatchResult, err error) {
	require.NoError(tool, err)
	require.Equal(tool, 0, len(mr.MatchedRecords))
}

func newTestOrder(price, amount uint64, isBuy bool) *WrapOrder {
	o := &WrapOrder{}
	o.LimitPrice = big.NewInt(0).SetUint64(price)
	o.Amount = big.NewInt(0).SetUint64(amount)
	o.LeftAmount = big.NewInt(0).SetUint64(amount)
	o.FrozenAmount = big.NewInt(0)
	o.TriggerPrice = big.NewInt(0)
	o.LimitFee = big.NewInt(0)
	// time.Now().Unix()*2 to avoid to be pruned
	// rand.Int63() to avoid repeated orderHash
	o.Expiration = big.NewInt(time.Now().Unix()*2 + rand.Int63())
	if !isBuy {
		o.Flags[31] = 1
	}
	return o
}
