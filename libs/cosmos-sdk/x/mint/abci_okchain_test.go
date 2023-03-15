package mint_test

import (
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/simapp"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint"
	internaltypes "github.com/okex/exchain/libs/cosmos-sdk/x/mint/internal/types"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	BlocksPerYear  uint64 = 10519200        // Block per year, uint64(60 * 60 * 8766 / 3)
	DeflationEpoch uint64 = 3               // Default epoch, 3 year
	DeflationRate  string = "0.5"           // Default deflation rate 0.5
	FarmProportion string = "0.5"           // Default farm proportion 0.5
	Denom          string = "okt"           // OKT
	FeeAccountName string = "fee_collector" // Fee account

	InitStartBlock    int64  = 17601985 // Current mainnet block,  17601985
	InitStartSupply   int64  = 19210060 // Current mainnet supply, 19210060
	BlocksPerYearNew  uint64 = 8304636  // Reset new block per year, uint64(60 * 60 * 8766 / 3.8)
	DeflationEpochNew uint64 = 9        // Reset epoch, year to month, and 3 to 9
	Target24DayBlock  uint64 = 5000     // 24 day blocks must be 555000, but fake to 5000
)

// returns context and an app with updated mint keeper
func createTestApp() (*simapp.SimApp, sdk.Context) {
	isCheckTx := false
	app := simapp.Setup(isCheckTx)

	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.MintKeeper.SetParams(ctx, internaltypes.DefaultParams())
	app.MintKeeper.SetMinter(ctx, internaltypes.InitialMinterCustom())
	types.UnittestOnlySetMilestoneVenus5Height(0)

	return app, ctx
}

type AbciOkchainSuite struct {
	suite.Suite
}

func TestAbciOkchainSuite(t *testing.T) {
	suite.Run(t, new(AbciOkchainSuite))
}

func (suite *AbciOkchainSuite) TestNormalMint() {
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

func (suite *AbciOkchainSuite) TestDateAndSupply() {
	// TODO Check expected date and total supply
}

func (suite *AbciOkchainSuite) TestFakeUpdateNextBlock() {
	simApp, ctx := createTestApp()
	allRewards := sdk.NewDec(InitStartSupply)

	suite.step1(sdk.MustNewDecFromStr("0.5"), &ctx, simApp, &allRewards)
	suite.step2(sdk.MustNewDecFromStr("0.5"), &ctx, simApp, &allRewards)
	suite.step3(sdk.MustNewDecFromStr("0.5"), &ctx, simApp, &allRewards)
	suite.step4(sdk.MustNewDecFromStr("0.5"), &ctx, simApp, &allRewards)
	suite.step5(sdk.MustNewDecFromStr("0.5"), &ctx, simApp, &allRewards)
	suite.step6(sdk.MustNewDecFromStr("0.25"), &ctx, simApp, &allRewards)
	suite.step7(sdk.MustNewDecFromStr("0.25"), &ctx, simApp, &allRewards)
	suite.step8(sdk.MustNewDecFromStr("0.125"), &ctx, simApp, &allRewards)
	suite.LoopDeflation(&ctx, simApp, &allRewards)
}

func (suite *AbciOkchainSuite) step1(expectReward sdk.Dec, ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// Upgrade the main network code, wait N height to take effect.
	ctx.SetBlockHeight(InitStartBlock)
	coins := []sdk.Coin{sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(InitStartSupply))}
	_ = simApp.SupplyKeeper.MintCoins(*ctx, mint.ModuleName, coins)
	_ = simApp.SupplyKeeper.SendCoinsFromModuleToModule(*ctx, mint.ModuleName, FeeAccountName, coins)

	// Execution block.
	mint.BeginBlocker(*ctx, simApp.MintKeeper)
	feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
	minter := simApp.MintKeeper.GetMinterCustom(*ctx)
	*allRewards = allRewards.Add(expectReward)

	// Suit expect.
	require.Equal(suite.T(), minter.NextBlockToUpdate, BlocksPerYear*DeflationEpoch)
	require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom).MulTruncate(sdk.MustNewDecFromStr(FarmProportion)), expectReward)
	require.Equal(suite.T(), feeAccount.GetCoins().AmountOf(Denom), *allRewards)
	require.Equal(suite.T(), InitStartBlock, ctx.BlockHeight())

	// The target N height to take effect.
	types.UnittestOnlySetMilestoneVenus5Height(InitStartBlock + 1000)
}

func (suite *AbciOkchainSuite) step2(expectReward sdk.Dec, ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	//System: block height N+1
	for i := int64(1); i <= 1000+1; i++ {
		// Execution block.
		ctx.SetBlockHeight(ctx.BlockHeight() + 1)
		mint.BeginBlocker(*ctx, simApp.MintKeeper)
		minter := simApp.MintKeeper.GetMinterCustom(*ctx)
		feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
		*allRewards = allRewards.Add(expectReward)

		// Suit expect.
		require.Equal(suite.T(), minter.NextBlockToUpdate, BlocksPerYear*DeflationEpoch)
		require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom).MulTruncate(sdk.MustNewDecFromStr(FarmProportion)), expectReward)
		require.Equal(suite.T(), feeAccount.GetCoins().AmountOf(Denom), *allRewards)
	}

	require.Equal(suite.T(), InitStartBlock+1000, types.GetVenus5Height())
	require.Equal(suite.T(), InitStartBlock+1000+1, ctx.BlockHeight())
}

func (suite *AbciOkchainSuite) step3(expectReward sdk.Dec, ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// Send change BlocksPerYear proposal (effective immediately), first proposal.
	params := simApp.MintKeeper.GetParams(*ctx)
	params.BlocksPerYear = BlocksPerYearNew
	simApp.MintKeeper.SetParams(*ctx, params)

	for i := int64(1); i <= 1000; i++ {
		// Execution block.
		ctx.SetBlockHeight(ctx.BlockHeight() + 1)
		mint.BeginBlocker(*ctx, simApp.MintKeeper)
		minter := simApp.MintKeeper.GetMinterCustom(*ctx)
		params = simApp.MintKeeper.GetParams(*ctx)
		feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
		*allRewards = allRewards.Add(expectReward)

		// Suit expect.
		require.Equal(suite.T(), minter.NextBlockToUpdate, BlocksPerYear*DeflationEpoch)
		require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom).MulTruncate(sdk.MustNewDecFromStr(FarmProportion)), expectReward)
		require.Equal(suite.T(), feeAccount.GetCoins().AmountOf(Denom), *allRewards)
		require.Equal(suite.T(), params.BlocksPerYear, BlocksPerYearNew)
		require.Equal(suite.T(), params.DeflationEpoch, DeflationEpoch)
	}

	require.Equal(suite.T(), InitStartBlock+1000*2+1, ctx.BlockHeight())
}

func (suite *AbciOkchainSuite) step4(expectReward sdk.Dec, ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// Send change DeflationEpoch proposal (effective immediately) from 3 to 9, second proposal.
	params := simApp.MintKeeper.GetParams(*ctx)
	params.DeflationEpoch = DeflationEpochNew
	simApp.MintKeeper.SetParams(*ctx, params)

	for i := int64(1); i <= 1000; i++ {
		// Execution block.
		ctx.SetBlockHeight(ctx.BlockHeight() + 1)
		mint.BeginBlocker(*ctx, simApp.MintKeeper)
		minter := simApp.MintKeeper.GetMinterCustom(*ctx)
		feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
		*allRewards = allRewards.Add(expectReward)
		params = simApp.MintKeeper.GetParams(*ctx)

		// Suit expect.
		require.Equal(suite.T(), minter.NextBlockToUpdate, BlocksPerYear*DeflationEpoch)
		require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom).MulTruncate(sdk.MustNewDecFromStr(FarmProportion)), expectReward)
		require.Equal(suite.T(), feeAccount.GetCoins().AmountOf(Denom), *allRewards)
		require.Equal(suite.T(), params.BlocksPerYear, BlocksPerYearNew)
		require.Equal(suite.T(), params.DeflationEpoch, DeflationEpochNew)
	}
	require.Equal(suite.T(), InitStartBlock+1000*3+1, ctx.BlockHeight())
}

func (suite *AbciOkchainSuite) step5(expectReward sdk.Dec, ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// Send forced changes to the NextBlockUpdate proposal, third proposal.
	minter := simApp.MintKeeper.GetMinterCustom(*ctx)
	minter.NextBlockToUpdate = uint64(ctx.BlockHeight() + 1000)
	simApp.MintKeeper.SetMinterCustom(*ctx, minter)

	for i := int64(1); i < 1000; i++ {
		// Execution block.
		ctx.SetBlockHeight(ctx.BlockHeight() + 1)
		mint.BeginBlocker(*ctx, simApp.MintKeeper)
		minter := simApp.MintKeeper.GetMinterCustom(*ctx)
		feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
		*allRewards = allRewards.Add(expectReward)
		params := simApp.MintKeeper.GetParams(*ctx)

		// Suit expect.
		require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom).MulTruncate(sdk.MustNewDecFromStr(FarmProportion)), expectReward)
		require.Equal(suite.T(), feeAccount.GetCoins().AmountOf(Denom), *allRewards)
		require.Equal(suite.T(), params.BlocksPerYear, BlocksPerYearNew)
		require.Equal(suite.T(), params.DeflationEpoch, DeflationEpochNew)
	}
	require.Equal(suite.T(), InitStartBlock+1000*4, ctx.BlockHeight())
}

func (suite *AbciOkchainSuite) step6(expectReward sdk.Dec, ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// System code triggers halving: 0.5->0.25
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	mint.BeginBlocker(*ctx, simApp.MintKeeper)
	minter := simApp.MintKeeper.GetMinterCustom(*ctx)
	feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
	*allRewards = allRewards.Add(expectReward)
	params := simApp.MintKeeper.GetParams(*ctx)

	// Suit expect.
	require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom).MulTruncate(sdk.MustNewDecFromStr(FarmProportion)), expectReward)
	require.Equal(suite.T(), feeAccount.GetCoins().AmountOf(Denom), *allRewards)
	require.Equal(suite.T(), params.BlocksPerYear, BlocksPerYearNew)
	require.Equal(suite.T(), params.DeflationEpoch, DeflationEpochNew)
	require.Equal(suite.T(), InitStartBlock+1000*4+1, ctx.BlockHeight())
}

func (suite *AbciOkchainSuite) step7(expectReward sdk.Dec, ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// Send forced changes to the NextBlockUpdate proposal, fourth proposal.
	minter := simApp.MintKeeper.GetMinterCustom(*ctx)
	minter.NextBlockToUpdate = uint64(ctx.BlockHeight()) + Target24DayBlock - 4000
	simApp.MintKeeper.SetMinterCustom(*ctx, minter)

	for i := int64(1); i < int64(Target24DayBlock)-4000; i++ {
		// Execution block.
		ctx.SetBlockHeight(ctx.BlockHeight() + 1)
		mint.BeginBlocker(*ctx, simApp.MintKeeper)
		minter := simApp.MintKeeper.GetMinterCustom(*ctx)
		feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
		*allRewards = allRewards.Add(expectReward)
		params := simApp.MintKeeper.GetParams(*ctx)

		// Suit expect.
		require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom).MulTruncate(sdk.MustNewDecFromStr(FarmProportion)), expectReward)
		require.Equal(suite.T(), feeAccount.GetCoins().AmountOf(Denom), *allRewards)
		require.Equal(suite.T(), params.BlocksPerYear, BlocksPerYearNew)
		require.Equal(suite.T(), params.DeflationEpoch, DeflationEpochNew)
	}

	require.Equal(suite.T(), InitStartBlock+int64(Target24DayBlock), ctx.BlockHeight())
}

func (suite *AbciOkchainSuite) step8(expectReward sdk.Dec, ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
	// System code triggers halving: 0.25->0.125
	ctx.SetBlockHeight(ctx.BlockHeight() + 1)
	mint.BeginBlocker(*ctx, simApp.MintKeeper)
	minter := simApp.MintKeeper.GetMinterCustom(*ctx)
	feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
	*allRewards = allRewards.Add(expectReward)
	params := simApp.MintKeeper.GetParams(*ctx)

	// Suit expect.
	require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom).MulTruncate(sdk.MustNewDecFromStr(FarmProportion)), expectReward)
	require.Equal(suite.T(), feeAccount.GetCoins().AmountOf(Denom), *allRewards)
	require.Equal(suite.T(), params.BlocksPerYear, BlocksPerYearNew)
	require.Equal(suite.T(), params.DeflationEpoch, DeflationEpochNew)
	require.Equal(suite.T(), InitStartBlock+int64(Target24DayBlock)+1, ctx.BlockHeight())
}

func (suite *AbciOkchainSuite) LoopDeflation(ctx *sdk.Context, simApp *simapp.SimApp, allRewards *sdk.Dec) {
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
			ctx.SetBlockHeight(InitStartBlock + int64(Target24DayBlock) + int64(BlocksPerYearNew*DeflationEpochNew*tc.phase))
			mint.BeginBlocker(*ctx, simApp.MintKeeper)
			feeAccount := simApp.SupplyKeeper.GetModuleAccount(*ctx, FeeAccountName)
			expect := feeAccount.GetCoins().AmountOf(Denom)
			reward := tc.mintedPerBlock.Sub(tc.mintedPerBlock.MulTruncate(sdk.MustNewDecFromStr(FarmProportion)))
			*allRewards = allRewards.Add(reward)
			require.Equal(suite.T(), expect, *allRewards)

			params := simApp.MintKeeper.GetParams(*ctx)
			minter := simApp.MintKeeper.GetMinterCustom(*ctx)
			require.Equal(suite.T(), params.MintDenom, Denom)
			require.Equal(suite.T(), params.BlocksPerYear, BlocksPerYearNew)
			require.Equal(suite.T(), params.DeflationRate, sdk.MustNewDecFromStr(DeflationRate))
			require.Equal(suite.T(), params.DeflationEpoch, DeflationEpochNew)
			require.Equal(suite.T(), params.FarmProportion, sdk.MustNewDecFromStr(FarmProportion))

			require.Equal(suite.T(), minter.NextBlockToUpdate, uint64(InitStartBlock)+Target24DayBlock+BlocksPerYearNew/12*DeflationEpochNew*(tc.phase+1)+1)
			require.Equal(suite.T(), minter.MintedPerBlock.AmountOf(Denom), tc.mintedPerBlock)
		})
	}
}
