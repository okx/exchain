package gasprice

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestRecommendGP(t *testing.T) {
	var testRecommendGP *big.Int
	config := NewGPOConfig(80, defaultPrice)
	gpo := NewOracle(config)

	coefficient := int64(200000)
	gpNum := 20000

	for blockNum := 1; blockNum <= 10; blockNum++ {
		for i := 0; i < gpNum; i++ {
			gp := big.NewInt(coefficient + params.GWei)
			gpo.CurrentBlockGPs.AddGP(gp)
			coefficient--
		}
		gpo.BlockGPQueue.Push(gpo.CurrentBlockGPs)
		gpo.CurrentBlockGPs.Clear()
		testRecommendGP = gpo.RecommendGP()
		require.NotNil(t, testRecommendGP)
	}
}

func BenchmarkRecommendGP(b *testing.B) {
	config := NewGPOConfig(80, defaultPrice)
	gpo := NewOracle(config)
	coefficient := int64(2000000)
	gpNum := 20000
	for blockNum := 1; blockNum <= 20; blockNum++ {
		for i := 0; i < gpNum; i++ {
			gp := big.NewInt(coefficient + params.GWei)
			gpo.CurrentBlockGPs.AddGP(gp)
			coefficient--
		}
		gpo.BlockGPQueue.Push(gpo.CurrentBlockGPs)
		gpo.CurrentBlockGPs.Clear()
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = gpo.RecommendGP()
	}
}
