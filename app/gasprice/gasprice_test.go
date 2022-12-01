package gasprice

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"

	appconfig "github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/types"
)

func TestOracle_RecommendGP(t *testing.T) {
	t.Run("case 1: mainnet case", func(t *testing.T) {
		// Case 1 reproduces the problem of GP increase when the OKC's block height is 13527188
		appconfig.GetOecConfig().SetDynamicGpMaxTxNum(300)
		appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(1000000)
		appconfig.GetOecConfig().SetDynamicGpMode(1)
		config := NewGPOConfig(80, 5)
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
		fmt.Println(testRecommendGP)
		//testRecommendGP == 0.1GWei
	})
	t.Run("case 2: not full tx, not full gasUsed, maxGasUsed configured, mode 1", func(t *testing.T) {
		appconfig.GetOecConfig().SetDynamicGpMaxTxNum(300)
		appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(1000000)
		appconfig.GetOecConfig().SetDynamicGpMode(1)
		config := NewGPOConfig(80, 5)
		var testRecommendGP *big.Int
		gpo := NewOracle(config)
		delta := int64(200000)
		gpNum := 200

		for blockNum := 1; blockNum <= 10; blockNum++ {
			for i := 0; i < gpNum; i++ {
				gp := big.NewInt(delta + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 3500)
				delta--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
	})
	t.Run("case 3: not full tx, not full gasUsed, maxGasUsed unconfigured, mode 0", func(t *testing.T) {
		appconfig.GetOecConfig().SetDynamicGpMaxTxNum(300)
		appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(-1)
		appconfig.GetOecConfig().SetDynamicGpMode(0)
		config := NewGPOConfig(80, 5)
		var testRecommendGP *big.Int
		gpo := NewOracle(config)
		coefficient := int64(200000)
		gpNum := 200

		for blockNum := 1; blockNum <= 10; blockNum++ {
			for i := 0; i < gpNum; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 3500) // chain is uncongested
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
	})
	t.Run("case 4: not full tx, not full gasUsed, maxGasUsed unconfigured, mode 1", func(t *testing.T) {
		appconfig.GetOecConfig().SetDynamicGpMaxTxNum(300)
		appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(-1)
		appconfig.GetOecConfig().SetDynamicGpMode(1)
		config := NewGPOConfig(80, 5)
		var testRecommendGP *big.Int
		gpo := NewOracle(config)
		coefficient := int64(200000)
		gpNum := 200

		for blockNum := 1; blockNum <= 10; blockNum++ {
			for i := 0; i < gpNum; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 3500)
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
	})

	t.Run("case 5: not full tx, full gasUsed, gp surge, mode 0", func(t *testing.T) {
		appconfig.GetOecConfig().SetDynamicGpMaxTxNum(200)
		appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(1000000)
		appconfig.GetOecConfig().SetDynamicGpMode(0)
		config := NewGPOConfig(80, 5)
		var testRecommendGP *big.Int
		gpo := NewOracle(config)
		coefficient := int64(200000)
		gpNum := 180

		for blockNum := 1; blockNum <= 5; blockNum++ {
			for i := 0; i < gpNum; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 6000) // chain is congested
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
		for blockNum := 5; blockNum <= 10; blockNum++ {
			for i := 0; i < gpNum/2; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 6000) // chain is congested
				coefficient--
			}
			for i := gpNum / 2; i < gpNum; i++ {
				gp := big.NewInt((coefficient + params.GWei) * 100)
				gpo.CurrentBlockGPs.Update(gp, 6000)
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
	})
	t.Run("case 6: not full tx, full gasUsed, gp surge, mode 1", func(t *testing.T) {
		appconfig.GetOecConfig().SetDynamicGpMaxTxNum(300)
		appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(1000000)
		appconfig.GetOecConfig().SetDynamicGpMode(1)
		config := NewGPOConfig(80, 5)
		var testRecommendGP *big.Int
		gpo := NewOracle(config)
		coefficient := int64(200000)
		gpNum := 200

		for blockNum := 1; blockNum <= 5; blockNum++ {
			for i := 0; i < gpNum; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 4500) // chain is congested
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
		for blockNum := 5; blockNum <= 10; blockNum++ {
			for i := 0; i < gpNum/2; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 4500) // chain is congested
				coefficient--
			}
			for i := gpNum / 2; i < gpNum; i++ {
				gp := big.NewInt((coefficient + params.GWei) * 100)
				gpo.CurrentBlockGPs.Update(gp, 4500) // chain is congested
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
	})
	t.Run("case 7: full tx, not full gasUsed, gp surge, mode 0", func(t *testing.T) {
		appconfig.GetOecConfig().SetDynamicGpMaxTxNum(300)
		appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(1000000)
		appconfig.GetOecConfig().SetDynamicGpMode(0)
		config := NewGPOConfig(80, 5)
		var testRecommendGP *big.Int
		gpo := NewOracle(config)
		coefficient := int64(200000)
		gpNum := 300
		for blockNum := 1; blockNum <= 5; blockNum++ {
			for i := 0; i < gpNum; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 450) // chain is congested
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
		for blockNum := 5; blockNum <= 10; blockNum++ {
			for i := 0; i < gpNum/2; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 450) // chain is congested
				coefficient--
			}
			for i := gpNum / 2; i < gpNum; i++ {
				gp := big.NewInt((coefficient + params.GWei) * 100)
				gpo.CurrentBlockGPs.Update(gp, 450) // chain is congested
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
	})
	t.Run("case 8: full tx, not full gasUsed, gp surge, mode 1", func(t *testing.T) {
		appconfig.GetOecConfig().SetDynamicGpMaxTxNum(300)
		appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(1000000)
		appconfig.GetOecConfig().SetDynamicGpMode(1)
		config := NewGPOConfig(80, 5)
		var testRecommendGP *big.Int
		gpo := NewOracle(config)
		coefficient := int64(200000)
		gpNum := 300
		for blockNum := 1; blockNum <= 5; blockNum++ {
			for i := 0; i < gpNum; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 450) // chain is congested
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
		for blockNum := 5; blockNum <= 10; blockNum++ {
			for i := 0; i < gpNum/2; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 450)
				coefficient--
			}
			for i := gpNum / 2; i < gpNum; i++ {
				gp := big.NewInt((coefficient + params.GWei) * 100)
				gpo.CurrentBlockGPs.Update(gp, 450)
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
	})
	t.Run("case 9: not full tx, full gasUsed, gp decrease, mode 0", func(t *testing.T) {
		appconfig.GetOecConfig().SetDynamicGpMaxTxNum(300)
		appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(1000000)
		appconfig.GetOecConfig().SetDynamicGpMode(0)
		config := NewGPOConfig(80, 5)
		var testRecommendGP *big.Int
		gpo := NewOracle(config)
		coefficient := int64(200000)
		gpNum := 200

		for blockNum := 1; blockNum <= 5; blockNum++ {
			for i := 0; i < gpNum/2; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 6000)
				coefficient--
			}
			for i := gpNum / 2; i < gpNum; i++ {
				gp := big.NewInt((coefficient + params.GWei) * 100)
				gpo.CurrentBlockGPs.Update(gp, 6000)
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
		for blockNum := 5; blockNum <= 10; blockNum++ {
			for i := 0; i < gpNum; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 6000)
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
	})
	t.Run("case 10: not full tx, full gasUsed, gp decrease, mode 1", func(t *testing.T) {
		appconfig.GetOecConfig().SetDynamicGpMaxTxNum(300)
		appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(1000000)
		appconfig.GetOecConfig().SetDynamicGpMode(1)
		config := NewGPOConfig(80, 5)
		var testRecommendGP *big.Int
		gpo := NewOracle(config)
		coefficient := int64(200000)
		gpNum := 200

		for blockNum := 1; blockNum <= 5; blockNum++ {
			for i := 0; i < gpNum/2; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 6000)
				coefficient--
			}
			for i := gpNum / 2; i < gpNum; i++ {
				gp := big.NewInt((coefficient + params.GWei) * 100)
				gpo.CurrentBlockGPs.Update(gp, 6000)
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
		for blockNum := 5; blockNum <= 10; blockNum++ {
			for i := 0; i < gpNum; i++ {
				gp := big.NewInt(coefficient + params.GWei)
				gpo.CurrentBlockGPs.Update(gp, 6000)
				coefficient--
			}
			cp := gpo.CurrentBlockGPs.Copy()
			gpo.BlockGPQueue.Push(cp)
			testRecommendGP = gpo.RecommendGP()
			gpo.CurrentBlockGPs.Clear()
			require.NotNil(t, testRecommendGP)
			fmt.Println(testRecommendGP)
		}
	})
}

func BenchmarkRecommendGP(b *testing.B) {
	appconfig.GetOecConfig().SetDynamicGpMaxTxNum(300)
	appconfig.GetOecConfig().SetDynamicGpMaxGasUsed(1000000)
	appconfig.GetOecConfig().SetDynamicGpMode(0)
	config := NewGPOConfig(80, 6)
	gpo := NewOracle(config)
	coefficient := int64(2000000)
	gpNum := 300
	for blockNum := 1; blockNum <= 20; blockNum++ {
		for i := 0; i < gpNum; i++ {
			gp := big.NewInt(coefficient + params.GWei)
			gpo.CurrentBlockGPs.Update(gp, 6000)
			coefficient--
		}
		cp := gpo.CurrentBlockGPs.Copy()
		gpo.BlockGPQueue.Push(cp)
		gpo.CurrentBlockGPs.Clear()
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = gpo.RecommendGP()
	}
}
