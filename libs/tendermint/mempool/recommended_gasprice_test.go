package mempool

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"

	"github.com/okx/okbchain/libs/tendermint/abci/example/kvstore"
	cfg "github.com/okx/okbchain/libs/tendermint/config"
	"github.com/okx/okbchain/libs/tendermint/proxy"
	"github.com/okx/okbchain/libs/tendermint/types"
)

// tx for recommended gas price test
type RGPTX struct {
	GU uint64
	GP *big.Int
}

func NewRGPTX(gu uint64, gp *big.Int) RGPTX {
	return RGPTX{
		GU: gu,
		GP: gp,
	}
}

func TestCListMempool_RecommendGP(t *testing.T) {
	testCases := []struct {
		title               string
		curBlockRGP         *big.Int
		isCurBlockCongested bool

		gpMaxTxNum   int64
		gpMaxGasUsed int64
		gpMode       int

		prepare func(int, int64, int64)
		// build txs for one block
		buildTxs func(int, int64, *int64, bool, bool) []RGPTX

		expectedTotalGU     []uint64
		expectedRecommendGp []string
		blocks              int

		needMultiple bool
		gpDecrease   bool
		tpb          []int
	}{
		{
			title:               "5/5 empty block, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{0, 0, 0, 0, 0},
			expectedRecommendGp: []string{"100000000", "100000000", "100000000", "100000000", "100000000"},
			blocks:              5,
			needMultiple:        false,
			gpDecrease:          false,
			tpb:                 []int{0, 0, 0, 0, 0},
		},
		{
			title:               "4/6 empty block, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 0, 0, 0, 0, 46329800},
			expectedRecommendGp: []string{"100200099", "100000000", "100000000", "100000000", "100000000", "100200099"},
			blocks:              6,
			needMultiple:        false,
			gpDecrease:          false,
			tpb:                 []int{200, 0, 0, 0, 0, 200},
		},

		{
			title:               "4/6 uncongested block, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 23164900, 23164900, 23164900, 23164900, 46329800},
			expectedRecommendGp: []string{"100200099", "100000000", "100000000", "100000000", "100000000", "100200500"},
			blocks:              6,
			needMultiple:        false,
			gpDecrease:          false,
			tpb:                 []int{200, 100, 100, 100, 100, 200},
		},
		{
			title:               "0/5 empty block, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 46329800, 46329800, 46329800, 46329800},
			expectedRecommendGp: []string{"100200099", "100200099", "100200299", "100200499", "100200699"},
			blocks:              5,
			needMultiple:        false,
			gpDecrease:          false,
			tpb:                 []int{200, 200, 200, 200, 200},
		},
		{
			title:               "2/5 empty block, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 0, 46329800, 0, 46329800},
			expectedRecommendGp: []string{"100200099", "100000000", "100200099", "100000000", "100200299"},
			blocks:              5,
			needMultiple:        false,
			gpDecrease:          false,
			tpb:                 []int{200, 0, 200, 0, 200},
		},
		{
			title:               "2/5 empty block, uncongestion, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        60000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 0, 46329800, 0, 46329800},
			expectedRecommendGp: []string{"100000000", "100000000", "100000000", "100000000", "100000000"},
			blocks:              5,
			needMultiple:        false,
			gpDecrease:          false,
			tpb:                 []int{200, 0, 200, 0, 200},
		},
		{
			title:               "0/5 empty block, minimal gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.MinimalGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 46329800, 46329800, 46329800, 46329800},
			expectedRecommendGp: []string{"0", "0", "0", "0", "0"},
			blocks:              5,
			needMultiple:        false,
			gpDecrease:          false,
			tpb:                 []int{200, 200, 200, 200, 200},
		},
		{
			title:               "0/5 empty block, gp decrease, normal gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.NormalGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 46329800, 46329800, 46329800, 46329800},
			expectedRecommendGp: []string{"100199802", "100199802", "100199801", "100199603", "100199603"},
			blocks:              5,
			needMultiple:        false,
			gpDecrease:          true,
			tpb:                 []int{200, 200, 200, 200, 200},
		},
		{
			title:               "0/5 empty block, gp decrease, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 46329800, 46329800, 46329800, 46329800},
			expectedRecommendGp: []string{"100199900", "100199700", "100199700", "100199700", "100199700"},
			blocks:              5,
			needMultiple:        false,
			gpDecrease:          true,
			tpb:                 []int{200, 200, 200, 200, 200},
		},
		{
			title:               "0/5 empty block, normal mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.NormalGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 46329800, 46329800, 46329800, 46329800},
			expectedRecommendGp: []string{"100200001", "100200201", "100200400", "100200402", "100200602"},
			blocks:              5,
			needMultiple:        false,
			gpDecrease:          false,
			tpb:                 []int{200, 200, 200, 200, 200},
		},
		{
			title:               "0/5 empty block, uncongestion, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        60000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 46329800, 46329800, 46329800, 46329800},
			expectedRecommendGp: []string{"100000000", "100000000", "100000000", "100000000", "100000000"},
			blocks:              5,
			needMultiple:        false,
			gpDecrease:          false,
			tpb:                 []int{200, 200, 200, 200, 200},
		},
		{
			title:               "0/5 empty block, uncongestion, gp multiple, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        60000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 46329800, 46329800, 46329800, 46329800},
			expectedRecommendGp: []string{"100000000", "100000000", "100000000", "100000000", "100000000"},
			blocks:              5,
			needMultiple:        true,
			gpDecrease:          false,
			tpb:                 []int{200, 200, 200, 200, 200},
		},
		{
			title:               "0/5 empty block, congestion, gp multiple, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 46329800, 46329800, 46329800, 46329800},
			expectedRecommendGp: []string{"5060109900", "5060109900", "5060120000", "5060130100", "5060140200"},
			blocks:              5,
			needMultiple:        true,
			gpDecrease:          false,
			tpb:                 []int{200, 200, 200, 200, 200},
		},
		{
			title:               "0/5 empty block, congestion, gp decrease, gp multiple, higher gp mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.CongestionHigherGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 46329800, 46329800, 46329800, 46329800},
			expectedRecommendGp: []string{"5060094901", "5060084801", "5060084801", "5060084801", "5060084801"},
			blocks:              5,
			needMultiple:        true,
			gpDecrease:          true,
			tpb:                 []int{200, 200, 200, 200, 200},
		},
		{
			title:               "0/5 empty block, congestion, gp multiple, normal mode",
			curBlockRGP:         big.NewInt(0),
			gpMaxTxNum:          300,
			gpMaxGasUsed:        40000000,
			gpMode:              types.NormalGpMode,
			prepare:             setMocConfig,
			buildTxs:            generateTxs,
			expectedTotalGU:     []uint64{46329800, 46329800, 46329800, 46329800, 46329800},
			expectedRecommendGp: []string{"100200001", "100200201", "100200400", "100200402", "100200602"},
			blocks:              5,
			needMultiple:        true,
			gpDecrease:          false,
			tpb:                 []int{200, 200, 200, 200, 200},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(tt *testing.T) {
			// init mempool
			app := kvstore.NewApplication()
			cc := proxy.NewLocalClientCreator(app)
			mempool, cleanup := newMempoolWithApp(cc)
			defer cleanup()

			tc.prepare(tc.gpMode, tc.gpMaxTxNum, tc.gpMaxGasUsed)
			gpOffset := int64(200000)
			baseGP := int64(params.GWei / 10)

			for i := 0; i < tc.blocks; i++ {
				totalGasUsed := uint64(0)
				txs := tc.buildTxs(tc.tpb[i], baseGP, &gpOffset, tc.gpDecrease, tc.needMultiple)
				for _, tx := range txs {
					if cfg.DynamicConfig.GetDynamicGpMode() != types.MinimalGpMode {
						mempool.gpo.CurrentBlockGPs.Update(tx.GP, tx.GU)
					}
					totalGasUsed += tx.GU
				}
				require.True(tt, totalGasUsed == tc.expectedTotalGU[i], "block gas expect %d, but get %d ", tc.expectedTotalGU[i], totalGasUsed)
				if cfg.DynamicConfig.GetDynamicGpMode() != types.MinimalGpMode {
					// calculate recommended GP
					currentBlockGPsCopy := mempool.gpo.CurrentBlockGPs.Copy()
					err := mempool.gpo.BlockGPQueue.Push(currentBlockGPsCopy)
					require.Nil(tt, err)
					tc.curBlockRGP, tc.isCurBlockCongested = mempool.gpo.RecommendGP()
					tc.curBlockRGP = postProcessGP(tc.curBlockRGP, tc.isCurBlockCongested)
					mempool.gpo.CurrentBlockGPs.Clear()
				}
				//fmt.Println("current recommended GP: ", tc.curBlockRGP)
				require.True(tt, tc.expectedRecommendGp[i] == tc.curBlockRGP.String(), "recommend gas price expect %s, but get %s ", tc.expectedRecommendGp[i], tc.curBlockRGP.String())
			}
		})
	}
}

func generateTxs(totalTxNum int, baseGP int64, gpOffset *int64, needDecreaseGP bool, needMultiple bool) []RGPTX {
	// guPerTx is gas used of evm contract - Owner.sol
	// bytecode := ethcommon.FromHex("0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a36102c4806100dc6000396000f3fe608060405234801561001057600080fd5b5060043610610053576000357c010000000000000000000000000000000000000000000000000000000090048063893d20e814610058578063a6f9dae1146100a2575b600080fd5b6100606100e6565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6100e4600480360360208110156100b857600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061010f565b005b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146101d1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260138152602001807f43616c6c6572206973206e6f74206f776e65720000000000000000000000000081525060200191505060405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f342827c97908e5e2f71151c08502a66d44b6f758e3ac2f1de95f02eb95f0a73560405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505056fea265627a7a72315820f397f2733a89198bc7fed0764083694c5b828791f39ebcbc9e414bccef14b48064736f6c63430005100032")
	guPerTx := uint64(231649)
	txs := make([]RGPTX, 0)
	for txCount := 0; txCount < totalTxNum/2; txCount++ {
		curTxGp := baseGP + *gpOffset
		if !needDecreaseGP {
			*gpOffset++
		} else {
			*gpOffset--
		}
		tx := NewRGPTX(guPerTx, big.NewInt(curTxGp))
		txs = append(txs, tx)
	}

	multiple := int64(1)
	if needMultiple {
		multiple = 100
	}
	for txCount := totalTxNum / 2; txCount < totalTxNum; txCount++ {
		curTxGp := (baseGP + *gpOffset) * multiple
		if !needDecreaseGP {
			*gpOffset++
		} else {
			*gpOffset--
		}
		tx := NewRGPTX(guPerTx, big.NewInt(curTxGp))

		txs = append(txs, tx)
	}
	return txs
}

func setMocConfig(gpMode int, gpMaxTxNum int64, gpMaxGasUsed int64) {
	moc := cfg.MockDynamicConfig{}
	moc.SetDynamicGpMode(gpMode)
	moc.SetDynamicGpMaxTxNum(gpMaxTxNum)
	moc.SetDynamicGpMaxGasUsed(gpMaxGasUsed)

	cfg.SetDynamicConfig(moc)
}

func postProcessGP(recommendedGP *big.Int, isCongested bool) *big.Int {
	// minGP for test is 0.1GWei
	minGP := big.NewInt(100000000)
	maxGP := new(big.Int).Mul(minGP, big.NewInt(5000))

	rgp := new(big.Int).Set(minGP)
	if cfg.DynamicConfig.GetDynamicGpMode() != types.MinimalGpMode {
		// If current block is not congested, rgp == minimal gas price.
		if isCongested {
			rgp.Set(recommendedGP)
		}

		if rgp.Cmp(minGP) == -1 {
			rgp.Set(minGP)
		}

		if rgp.Cmp(maxGP) == 1 {
			rgp.Set(maxGP)
		}
	}
	return rgp
}
