package gasprice

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"

	appconfig "github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/types"
)

func TestOracle_RecommendGP(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		appconfig.GetOecConfig().SetMaxTxNumPerBlock(300)
		appconfig.GetOecConfig().SetMaxGasUsedPerBlock(1000000)
		appconfig.GetOecConfig().SetDynamicGpCheckBlocks(5)
		config := NewGPOConfig(80, appconfig.GetOecConfig().GetDynamicGpCheckBlocks())
		var testRecommendGP *big.Int
		gpo := NewOracle(config)
		coefficient := int64(200000)
		gpNum := 20000

		for blockNum := 1; blockNum <= 10; blockNum++ {
			for i := 0; i < gpNum; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 35) // chain is uncongested
				//gpo.CurrentBlockGPs.Update(gp, 45) // chain is congested
				coefficient--
			}
			gpo.BlockGPQueue.Push(gpo.CurrentBlockGPs)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			//fmt.Println(testRecommendGP)
		}
	})
	t.Run("case 2", func(t *testing.T) {
		// Case 2 reproduces the problem of GP increase when the OKC's block height is 13527188
		appconfig.GetOecConfig().SetMaxTxNumPerBlock(300)
		appconfig.GetOecConfig().SetMaxGasUsedPerBlock(1000000)
		appconfig.GetOecConfig().SetDynamicGpCheckBlocks(5)
		config := NewGPOConfig(80, appconfig.GetOecConfig().GetDynamicGpCheckBlocks())
		var testRecommendGP *big.Int
		gpo := NewOracle(config)

		blockPGs13527188 := types.NewSingleBlockGPs()
		blockPGs13527188.Update(big.NewInt(10*params.GWei), 35)
		blockPGs13527188.Update(big.NewInt(10*params.GWei), 35)
		blockPGs13527188.Update(big.NewInt(params.GWei/10), 35)

		blockPGs13527187 := types.NewSingleBlockGPs()
		blockPGs13527187.Update(big.NewInt(params.GWei/10), 35)
		blockPGs13527187.Update(big.NewInt(params.GWei/10), 35)

		blockPGs13527186 := types.NewSingleBlockGPs()
		blockPGs13527186.Update(big.NewInt(params.GWei/10), 35)
		blockPGs13527186.Update(big.NewInt(params.GWei/10), 35)

		blockPGs13527185 := types.NewSingleBlockGPs()
		blockPGs13527185.Update(big.NewInt(params.GWei/10+params.GWei/1000), 35)
		blockPGs13527185.Update(big.NewInt(params.GWei/10), 35)
		blockPGs13527185.Update(big.NewInt(params.GWei/10), 35)

		blockPGs13527184 := types.NewSingleBlockGPs()
		blockPGs13527184.Update(big.NewInt(params.GWei*3/10), 35)
		blockPGs13527184.Update(big.NewInt(params.GWei*3/10), 35)
		blockPGs13527184.Update(big.NewInt(params.GWei*3/10), 35)
		blockPGs13527184.Update(big.NewInt(params.GWei/10), 35)
		blockPGs13527184.Update(big.NewInt(params.GWei/10), 35)
		blockPGs13527184.Update(big.NewInt(params.GWei/10), 35)

		gpo.BlockGPQueue.Push(blockPGs13527184)
		gpo.BlockGPQueue.Push(blockPGs13527185)
		gpo.BlockGPQueue.Push(blockPGs13527186)
		gpo.BlockGPQueue.Push(blockPGs13527187)
		gpo.BlockGPQueue.Push(blockPGs13527188)

		testRecommendGP = gpo.RecommendGP()
		require.NotNil(t, testRecommendGP)
		//fmt.Println(testRecommendGP)
		//testRecommendGP == 0.1GWei
	})
}

func BenchmarkRecommendGP(b *testing.B) {
	appconfig.GetOecConfig().SetMaxTxNumPerBlock(300)
	appconfig.GetOecConfig().SetMaxGasUsedPerBlock(1000000)
	appconfig.GetOecConfig().SetDynamicGpCheckBlocks(6)
	config := NewGPOConfig(80, appconfig.GetOecConfig().GetDynamicGpCheckBlocks())
	gpo := NewOracle(config)
	coefficient := int64(2000000)
	gpNum := 20000
	for blockNum := 1; blockNum <= 20; blockNum++ {
		for i := 0; i < gpNum; i++ {
			gp := big.NewInt(coefficient + params.GWei)
			gpo.CurrentBlockGPs.Update(gp, 35)
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
