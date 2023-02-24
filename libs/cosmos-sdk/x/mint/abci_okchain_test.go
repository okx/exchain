package mint_test

import (
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/simapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint"
	internaltypes "github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	BlocksPerYear  uint64 = (60 * 60 * 8766 / 3)
	DeflationEpoch uint64 = 3
	DeflationRate  string = "0.5"
	FarmProportion string = "0.5"
	Denom          string = "okt"
	FeeAccountName string = "fee_collector"
)

// returns context and an app with updated mint keeper
func createTestApp() (*simapp.SimApp, sdk.Context) {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.MintKeeper.SetParams(ctx, internaltypes.DefaultParams())
	app.MintKeeper.SetMinter(ctx, internaltypes.InitialMinterCustom())

	return app, ctx
}

type AbciOkchainSuite struct {
	suite.Suite
}

func TestAbciOkchainSuite(t *testing.T) {
	suite.Run(t, new(AbciOkchainSuite))
}

func (suite *AbciOkchainSuite) TestNormalBlockRewards() {
	testCases := []struct {
		title          string
		phase          uint64
		mintedPerBlock sdk.Dec
	}{
		{"phase 0", 0, sdk.MustNewDecFromStr("1.0")},
		{"phase 1", 1, sdk.MustNewDecFromStr("0.5")},
		{"phase 2", 2, sdk.MustNewDecFromStr("0.25")},
		{"phase 3", 3, sdk.MustNewDecFromStr("0.125")},
		{"phase 4", 4, sdk.MustNewDecFromStr("0.0625")},
		{"phase 5", 5, sdk.MustNewDecFromStr("0.03125")},
		{"phase 6", 6, sdk.MustNewDecFromStr("0.015625")},
		{"phase 7", 7, sdk.MustNewDecFromStr("0.0078125")},
		{"phase 8", 8, sdk.MustNewDecFromStr("0.00390625")},
		{"phase 9", 9, sdk.MustNewDecFromStr("0.001953125")},
		{"phase 10", 10, sdk.MustNewDecFromStr("0.0009765625")},
		{"phase 11", 11, sdk.MustNewDecFromStr("0.00048828125")},
		{"phase 12", 12, sdk.MustNewDecFromStr("0.000244140625")},
		{"phase 13", 13, sdk.MustNewDecFromStr("0.0001220703125")},
		{"phase 14", 14, sdk.MustNewDecFromStr("0.00006103515625")},
		{"phase 15", 15, sdk.MustNewDecFromStr("0.000030517578125")},
		{"phase 16", 16, sdk.MustNewDecFromStr("0.0000152587890625")},
		{"phase 17", 17, sdk.MustNewDecFromStr("0.00000762939453125")},
		{"phase 18", 18, sdk.MustNewDecFromStr("0.000003814697265625")},
		{"phase 19", 19, sdk.MustNewDecFromStr("0.000001907348632812")},
		{"phase 20", 20, sdk.MustNewDecFromStr("0.000000953674316406")},
		{"phase 21", 21, sdk.MustNewDecFromStr("0.000000476837158203")},
		{"phase 22", 22, sdk.MustNewDecFromStr("0.000000238418579102")},
		{"phase 23", 23, sdk.MustNewDecFromStr("0.000000119209289551")},
		{"phase 24", 24, sdk.MustNewDecFromStr("0.000000059604644776")},
		{"phase 25", 25, sdk.MustNewDecFromStr("0.000000029802322388")},
		{"phase 26", 26, sdk.MustNewDecFromStr("0.000000014901161194")},
		{"phase 27", 27, sdk.MustNewDecFromStr("0.000000007450580597")},
		{"phase 28", 28, sdk.MustNewDecFromStr("0.000000003725290298")},
		{"phase 29", 29, sdk.MustNewDecFromStr("0.000000001862645149")},
		{"phase 30", 30, sdk.MustNewDecFromStr("0.000000000931322574")},
		{"phase 31", 31, sdk.MustNewDecFromStr("0.000000000465661287")},
		{"phase 32", 32, sdk.MustNewDecFromStr("0.000000000232830644")},
		{"phase 33", 33, sdk.MustNewDecFromStr("0.000000000116415322")},
		{"phase 34", 34, sdk.MustNewDecFromStr("0.000000000058207661")},
		{"phase 35", 35, sdk.MustNewDecFromStr("0.000000000029103830")},
		{"phase 36", 36, sdk.MustNewDecFromStr("0.000000000014551915")},
		{"phase 37", 37, sdk.MustNewDecFromStr("0.000000000007275958")},
		{"phase 38", 38, sdk.MustNewDecFromStr("0.000000000003637979")},
		{"phase 39", 39, sdk.MustNewDecFromStr("0.000000000001818990")},
		{"phase 40", 40, sdk.MustNewDecFromStr("0.000000000000909495")},
		{"phase 41", 41, sdk.MustNewDecFromStr("0.000000000000454748")},
		{"phase 42", 42, sdk.MustNewDecFromStr("0.000000000000227374")},
		{"phase 43", 43, sdk.MustNewDecFromStr("0.000000000000113687")},
		{"phase 44", 44, sdk.MustNewDecFromStr("0.000000000000056844")},
		{"phase 45", 45, sdk.MustNewDecFromStr("0.000000000000028422")},
		{"phase 46", 46, sdk.MustNewDecFromStr("0.000000000000014211")},
		{"phase 47", 47, sdk.MustNewDecFromStr("0.000000000000007106")},
		{"phase 48", 48, sdk.MustNewDecFromStr("0.000000000000003553")},
		{"phase 49", 49, sdk.MustNewDecFromStr("0.000000000000001776")},
		{"phase 50", 50, sdk.MustNewDecFromStr("0.000000000000000888")},
		{"phase 51", 51, sdk.MustNewDecFromStr("0.000000000000000444")},
		{"phase 52", 52, sdk.MustNewDecFromStr("0.000000000000000222")},
		{"phase 53", 53, sdk.MustNewDecFromStr("0.000000000000000111")},
		{"phase 54", 54, sdk.MustNewDecFromStr("0.000000000000000056")},
		{"phase 55", 55, sdk.MustNewDecFromStr("0.000000000000000028")},
		{"phase 56", 56, sdk.MustNewDecFromStr("0.000000000000000014")},
		{"phase 57", 57, sdk.MustNewDecFromStr("0.000000000000000007")},
		{"phase 58", 58, sdk.MustNewDecFromStr("0.000000000000000004")},
		{"phase 59", 59, sdk.MustNewDecFromStr("0.000000000000000002")},
		{"phase 60", 60, sdk.MustNewDecFromStr("0.000000000000000001")},
		{"phase 61", 61, sdk.MustNewDecFromStr("0.000000000000000000")},
		{"phase 62", 62, sdk.MustNewDecFromStr("0.000000000000000000")},
	}

	simApp, ctx := createTestApp()

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			ctx.SetBlockHeight(int64(BlocksPerYear * DeflationEpoch * tc.phase))
			mint.BeginBlocker(ctx, simApp.MintKeeper)
			feeAccount := simApp.SupplyKeeper.GetModuleAccount(ctx, FeeAccountName)
			require.Equal(suite.T(), feeAccount.GetCoins().AmountOf(Denom),
				tc.mintedPerBlock.Sub(tc.mintedPerBlock.MulTruncate(sdk.MustNewDecFromStr(FarmProportion))))

			params := simApp.MintKeeper.GetParams(ctx)
			minter := simApp.MintKeeper.GetMinterCustom(ctx)
			require.Equal(suite.T(), params.MintDenom, Denom)
			require.Equal(suite.T(), params.BlocksPerYear, BlocksPerYear)
			require.Equal(suite.T(), params.DeflationRate, sdk.MustNewDecFromStr(DeflationRate))
			require.Equal(suite.T(), params.DeflationEpoch, DeflationEpoch)
			require.Equal(suite.T(), params.FarmProportion, sdk.MustNewDecFromStr(FarmProportion))

			require.Equal(suite.T(), minter.NextBlockToUpdate, BlocksPerYear*DeflationEpoch*(tc.phase+1))
			require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom), tc.mintedPerBlock)

			simApp.SupplyKeeper.SendCoinsFromModuleToModule(ctx, FeeAccountName, "bonded_tokens_pool", feeAccount.GetCoins())
		})
	}
}

const (
	CurrentBlock         int64  = 17601985 // 当前区块 17601985
	CurrentSupply        int64  = 19210060 // 当前奖励 19210060
	BlocksPerDay         uint64 = 22736    // 每天区块 22736 = 24 * 60 * 60 / 3.8
	DeflationEpochDay    uint64 = 273      // 需要天数 273 = 6207157 / 22736
	Target24DayBlock     uint64 = 555000   // 24天减半的时间点
	TargetDeflationBlock uint64 = 6207157  // 周期区块 6207157
	SupplyPhase0         int64  = 277500   // 第一阶段24天增发okt
)

func (suite *AbciOkchainSuite) initCurrentSupply(ctx *sdk.Context, simApp *simapp.SimApp, all_reward *sdk.Dec) {
	//init
	ctx.SetBlockHeight(CurrentBlock)
	coins := []sdk.Coin{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(CurrentSupply))}
	_ = simApp.SupplyKeeper.MintCoins(*ctx, mint.ModuleName, coins)
	_ = simApp.SupplyKeeper.SendCoinsFromModuleToModule(*ctx, mint.ModuleName, FeeAccountName, coins)

	mint.BeginBlocker(*ctx, simApp.MintKeeper)
	feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
	expect := feeAccount.GetCoins().AmountOf(Denom)
	minter := simApp.MintKeeper.GetMinterCustom(*ctx)
	require.Equal(suite.T(), minter.NextBlockToUpdate, BlocksPerYear*DeflationEpoch)
	require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom), sdk.MustNewDecFromStr("1.0"))
	reward := sdk.MustNewDecFromStr("1.0").Sub(sdk.MustNewDecFromStr("1.0").MulTruncate(sdk.MustNewDecFromStr(FarmProportion)))
	*all_reward = all_reward.Add(reward)
	require.Equal(suite.T(), expect, *all_reward)

	require.Equal(suite.T(), CurrentBlock, ctx.BlockHeight())
}

func (suite *AbciOkchainSuite) phase0(ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	//phase0
	var i int64
	for i = 1; i <= int64(Target24DayBlock-4000); i++ {
		ctx.SetBlockHeight(ctx.BlockHeight() + 1)
		mint.BeginBlocker(*ctx, simApp.MintKeeper)
		minter := simApp.MintKeeper.GetMinterCustom(*ctx)
		require.Equal(suite.T(), minter.NextBlockToUpdate, BlocksPerYear*DeflationEpoch)
		require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom), sdk.MustNewDecFromStr("1.0"))
		feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
		expect := feeAccount.GetCoins().AmountOf(Denom)
		reward := sdk.MustNewDecFromStr("1.0").Sub(sdk.MustNewDecFromStr("1.0").MulTruncate(sdk.MustNewDecFromStr(FarmProportion)))
		*allRewards = allRewards.Add(reward)
		require.Equal(suite.T(), expect, *allRewards)
	}
}

func (suite *AbciOkchainSuite) phase1(ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// The first proposal
	params := simApp.MintKeeper.GetParams(*ctx)
	params.BlocksPerYear = BlocksPerDay
	simApp.MintKeeper.SetParams(*ctx, params)

	for i := int64(1); i <= 1000; i++ {
		ctx.SetBlockHeight(ctx.BlockHeight() + 1)
		mint.BeginBlocker(*ctx, simApp.MintKeeper)
		minter := simApp.MintKeeper.GetMinterCustom(*ctx)
		require.Equal(suite.T(), minter.NextBlockToUpdate, BlocksPerYear*DeflationEpoch)
		require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom), sdk.MustNewDecFromStr("1.0"))
		feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
		expect := feeAccount.GetCoins().AmountOf(Denom)
		reward := sdk.MustNewDecFromStr("1.0").Sub(sdk.MustNewDecFromStr("1.0").MulTruncate(sdk.MustNewDecFromStr(FarmProportion)))
		*allRewards = allRewards.Add(reward)
		require.Equal(suite.T(), expect, *allRewards)

		params = simApp.MintKeeper.GetParams(*ctx)
		require.Equal(suite.T(), params.BlocksPerYear, BlocksPerDay)
		require.Equal(suite.T(), params.DeflationEpoch, DeflationEpoch)
	}
}

func (suite *AbciOkchainSuite) phase2(ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// The second proposal
	params := simApp.MintKeeper.GetParams(*ctx)
	params.DeflationEpoch = DeflationEpochDay
	simApp.MintKeeper.SetParams(*ctx, params)
	for i := int64(1); i <= 1000; i++ {
		ctx.SetBlockHeight(ctx.BlockHeight() + 1)
		mint.BeginBlocker(*ctx, simApp.MintKeeper)
		minter := simApp.MintKeeper.GetMinterCustom(*ctx)
		require.Equal(suite.T(), minter.NextBlockToUpdate, BlocksPerYear*DeflationEpoch)
		require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom), sdk.MustNewDecFromStr("1.0"))
		feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
		expect := feeAccount.GetCoins().AmountOf(Denom)
		reward := sdk.MustNewDecFromStr("1.0").Sub(sdk.MustNewDecFromStr("1.0").MulTruncate(sdk.MustNewDecFromStr(FarmProportion)))
		*allRewards = allRewards.Add(reward)
		require.Equal(suite.T(), expect, *allRewards)
		params = simApp.MintKeeper.GetParams(*ctx)
		require.Equal(suite.T(), params.BlocksPerYear, BlocksPerDay)
		require.Equal(suite.T(), params.DeflationEpoch, DeflationEpochDay)
	}

}

func (suite *AbciOkchainSuite) phase3(ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// The third proposal
	minter := simApp.MintKeeper.GetMinterCustom(*ctx)
	minter.NextBlockToUpdate = uint64(ctx.BlockHeight() + 1000)
	simApp.MintKeeper.SetMinterCustom(*ctx, minter)
	for i := int64(1); i <= 1000; i++ {
		ctx.SetBlockHeight(ctx.BlockHeight() + 1)
		mint.BeginBlocker(*ctx, simApp.MintKeeper)
		minter := simApp.MintKeeper.GetMinterCustom(*ctx)
		defaultMint := sdk.MustNewDecFromStr("1.0")
		if i == 1000 {
			defaultMint = sdk.MustNewDecFromStr("0.5")
		}
		require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom), defaultMint)
		feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
		expect := feeAccount.GetCoins().AmountOf(Denom)
		reward := defaultMint.Sub(defaultMint.MulTruncate(sdk.MustNewDecFromStr(FarmProportion)))
		*allRewards = allRewards.Add(reward)
		require.Equal(suite.T(), expect, *allRewards)
		params := simApp.MintKeeper.GetParams(*ctx)
		require.Equal(suite.T(), params.BlocksPerYear, BlocksPerDay)
		require.Equal(suite.T(), params.DeflationEpoch, DeflationEpochDay)
	}
}

func (suite *AbciOkchainSuite) phase4(ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// The fourth proposal
	minter := simApp.MintKeeper.GetMinterCustom(*ctx)
	minter.NextBlockToUpdate = uint64(ctx.BlockHeight() + 1000)
	simApp.MintKeeper.SetMinterCustom(*ctx, minter)
	for i := int64(1); i <= 1000; i++ {
		ctx.SetBlockHeight(ctx.BlockHeight() + 1)
		mint.BeginBlocker(*ctx, simApp.MintKeeper)
		minter := simApp.MintKeeper.GetMinterCustom(*ctx)
		defaultMint := sdk.MustNewDecFromStr("0.5")
		if i == 1000 {
			defaultMint = sdk.MustNewDecFromStr("0.25")
		}
		require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom), defaultMint)
		feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
		expect := feeAccount.GetCoins().AmountOf(Denom)
		reward := defaultMint.Sub(defaultMint.MulTruncate(sdk.MustNewDecFromStr(FarmProportion)))
		*allRewards = allRewards.Add(reward)
		require.Equal(suite.T(), expect, *allRewards)
		params := simApp.MintKeeper.GetParams(*ctx)
		require.Equal(suite.T(), params.BlocksPerYear, BlocksPerDay)
		require.Equal(suite.T(), params.DeflationEpoch, DeflationEpochDay)
	}
}

func (suite *AbciOkchainSuite) keeping(ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	testCases := []struct {
		title          string
		phase          uint64
		mintedPerBlock sdk.Dec
	}{
		{"phase 1", 1, sdk.MustNewDecFromStr("0.125")},
		{"phase 2", 2, sdk.MustNewDecFromStr("0.0625")},
		{"phase 3", 3, sdk.MustNewDecFromStr("0.03125")},
		{"phase 4", 4, sdk.MustNewDecFromStr("0.015625")},
		{"phase 5", 5, sdk.MustNewDecFromStr("0.0078125")},
		{"phase 6", 6, sdk.MustNewDecFromStr("0.00390625")},
		{"phase 7", 7, sdk.MustNewDecFromStr("0.001953125")},
		{"phase 8", 8, sdk.MustNewDecFromStr("0.0009765625")},
		{"phase 9", 9, sdk.MustNewDecFromStr("0.00048828125")},
		{"phase 10", 10, sdk.MustNewDecFromStr("0.000244140625")},
		{"phase 11", 11, sdk.MustNewDecFromStr("0.0001220703125")},
		{"phase 12", 12, sdk.MustNewDecFromStr("0.00006103515625")},
		{"phase 13", 13, sdk.MustNewDecFromStr("0.000030517578125")},
		{"phase 14", 14, sdk.MustNewDecFromStr("0.0000152587890625")},
		{"phase 15", 15, sdk.MustNewDecFromStr("0.00000762939453125")},
		{"phase 16", 16, sdk.MustNewDecFromStr("0.000003814697265625")},
		{"phase 17", 17, sdk.MustNewDecFromStr("0.000001907348632812")},
		{"phase 18", 18, sdk.MustNewDecFromStr("0.000000953674316406")},
		{"phase 19", 19, sdk.MustNewDecFromStr("0.000000476837158203")},
		{"phase 20", 20, sdk.MustNewDecFromStr("0.000000238418579102")},
		{"phase 21", 21, sdk.MustNewDecFromStr("0.000000119209289551")},
		{"phase 22", 22, sdk.MustNewDecFromStr("0.000000059604644776")},
		{"phase 23", 23, sdk.MustNewDecFromStr("0.000000029802322388")},
		{"phase 24", 24, sdk.MustNewDecFromStr("0.000000014901161194")},
		{"phase 25", 25, sdk.MustNewDecFromStr("0.000000007450580597")},
		{"phase 26", 26, sdk.MustNewDecFromStr("0.000000003725290298")},
		{"phase 27", 27, sdk.MustNewDecFromStr("0.000000001862645149")},
		{"phase 28", 28, sdk.MustNewDecFromStr("0.000000000931322574")},
		{"phase 29", 29, sdk.MustNewDecFromStr("0.000000000465661287")},
		{"phase 30", 30, sdk.MustNewDecFromStr("0.000000000232830644")},
		{"phase 31", 31, sdk.MustNewDecFromStr("0.000000000116415322")},
		{"phase 32", 32, sdk.MustNewDecFromStr("0.000000000058207661")},
		{"phase 33", 33, sdk.MustNewDecFromStr("0.000000000029103830")},
		{"phase 34", 34, sdk.MustNewDecFromStr("0.000000000014551915")},
		{"phase 35", 35, sdk.MustNewDecFromStr("0.000000000007275958")},
		{"phase 36", 36, sdk.MustNewDecFromStr("0.000000000003637979")},
		{"phase 37", 37, sdk.MustNewDecFromStr("0.000000000001818990")},
		{"phase 38", 38, sdk.MustNewDecFromStr("0.000000000000909495")},
		{"phase 39", 39, sdk.MustNewDecFromStr("0.000000000000454748")},
		{"phase 40", 40, sdk.MustNewDecFromStr("0.000000000000227374")},
		{"phase 41", 41, sdk.MustNewDecFromStr("0.000000000000113687")},
		{"phase 42", 42, sdk.MustNewDecFromStr("0.000000000000056844")},
		{"phase 43", 43, sdk.MustNewDecFromStr("0.000000000000028422")},
		{"phase 44", 44, sdk.MustNewDecFromStr("0.000000000000014211")},
		{"phase 45", 45, sdk.MustNewDecFromStr("0.000000000000007106")},
		{"phase 46", 46, sdk.MustNewDecFromStr("0.000000000000003553")},
		{"phase 47", 47, sdk.MustNewDecFromStr("0.000000000000001776")},
		{"phase 48", 48, sdk.MustNewDecFromStr("0.000000000000000888")},
		{"phase 49", 49, sdk.MustNewDecFromStr("0.000000000000000444")},
		{"phase 50", 50, sdk.MustNewDecFromStr("0.000000000000000222")},
		{"phase 51", 51, sdk.MustNewDecFromStr("0.000000000000000111")},
		{"phase 52", 52, sdk.MustNewDecFromStr("0.000000000000000056")},
		{"phase 53", 53, sdk.MustNewDecFromStr("0.000000000000000028")},
		{"phase 54", 54, sdk.MustNewDecFromStr("0.000000000000000014")},
		{"phase 55", 55, sdk.MustNewDecFromStr("0.000000000000000007")},
		{"phase 56", 56, sdk.MustNewDecFromStr("0.000000000000000004")},
		{"phase 57", 57, sdk.MustNewDecFromStr("0.000000000000000002")},
		{"phase 58", 58, sdk.MustNewDecFromStr("0.000000000000000001")},
		{"phase 59", 59, sdk.MustNewDecFromStr("0.000000000000000000")},
		{"phase 60", 60, sdk.MustNewDecFromStr("0.000000000000000000")},
	}

	for _, tc := range testCases {
		suite.Run(tc.title, func() {
			ctx.SetBlockHeight(CurrentBlock + int64(Target24DayBlock) + int64(BlocksPerDay*DeflationEpochDay*tc.phase))
			mint.BeginBlocker(*ctx, simApp.MintKeeper)
			feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
			expect := feeAccount.GetCoins().AmountOf(Denom)
			reward := tc.mintedPerBlock.Sub(tc.mintedPerBlock.MulTruncate(sdk.MustNewDecFromStr(FarmProportion)))
			*allRewards = allRewards.Add(reward)
			require.Equal(suite.T(), expect, *allRewards)

			params := simApp.MintKeeper.GetParams(*ctx)
			minter := simApp.MintKeeper.GetMinterCustom(*ctx)
			require.Equal(suite.T(), params.MintDenom, Denom)
			require.Equal(suite.T(), params.BlocksPerYear, BlocksPerDay)
			require.Equal(suite.T(), params.DeflationRate, sdk.MustNewDecFromStr(DeflationRate))
			require.Equal(suite.T(), params.DeflationEpoch, DeflationEpochDay)
			require.Equal(suite.T(), params.FarmProportion, sdk.MustNewDecFromStr(FarmProportion))

			require.Equal(suite.T(), minter.NextBlockToUpdate, uint64(CurrentBlock)+Target24DayBlock+BlocksPerDay*DeflationEpochDay*(tc.phase+1))
			require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom), tc.mintedPerBlock)
		})
	}
}

func (suite *AbciOkchainSuite) TestFakeUpdateNextBlock() {
	// init current supply block rewards
	simApp, ctx := createTestApp()
	allRewards := sdk.NewDec(CurrentSupply)

	suite.initCurrentSupply(&ctx, simApp, &allRewards)
	suite.phase0(&ctx, simApp, &allRewards)
	suite.phase1(&ctx, simApp, &allRewards)
	suite.phase2(&ctx, simApp, &allRewards)
	suite.phase3(&ctx, simApp, &allRewards)
	suite.phase4(&ctx, simApp, &allRewards)
	require.Equal(suite.T(), CurrentBlock+int64(Target24DayBlock), ctx.BlockHeight())
	suite.keeping(&ctx, simApp, &allRewards)
}
