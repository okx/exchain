package keeper

import (
	"testing"
	"time"

	"github.com/okex/exchain/x/common"

	"github.com/stretchr/testify/require"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/dex/types"
)

const TestProductNotExist = "product-not-exist"

func getTestTokenPair() *types.TokenPair {
	addr, err := sdk.AccAddressFromBech32(types.TestTokenPairOwner)
	if err != nil {
		panic(err)
	}
	return &types.TokenPair{
		BaseAssetSymbol:  "testToken",
		QuoteAssetSymbol: common.NativeToken,
		InitPrice:        sdk.MustNewDecFromStr("10.0"),
		MaxPriceDigit:    8,
		MaxQuantityDigit: 8,
		MinQuantity:      sdk.MustNewDecFromStr("0"),
		Owner:            addr,
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(0)),
	}
}

func TestSaveTokenPair(t *testing.T) {
	common.InitConfig()
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	keeper := testInput.DexKeeper
	tokenPair0 := getTestTokenPair()

	// SaveTokenPair
	err := keeper.SaveTokenPair(ctx, tokenPair0)
	require.Nil(t, err)

	maxTokenPairID := keeper.GetMaxTokenPairID(ctx)
	require.Equal(t, uint64(1), maxTokenPairID)

	// SaveTokenPair with id
	tokenPairId := uint64(100)
	tokenPair1 := getTestTokenPair()
	tokenPair1.ID = tokenPairId
	err = keeper.SaveTokenPair(ctx, tokenPair1)
	require.Nil(t, err)

	maxTokenPairID = keeper.GetMaxTokenPairID(ctx)
	require.Equal(t, tokenPairId, maxTokenPairID)

	// SaveTokenPair with smaller id
	tokenPair2 := getTestTokenPair()
	tokenPair2.ID = tokenPairId - 1
	err = keeper.SaveTokenPair(ctx, tokenPair2)
	require.Nil(t, err)

	maxTokenPairID = keeper.GetMaxTokenPairID(ctx)
	require.Equal(t, tokenPairId, maxTokenPairID)

}

func TestGetTokenPair(t *testing.T) {
	common.InitConfig()
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	keeper := testInput.DexKeeper
	tokenPair := getTestTokenPair()

	// SaveTokenPair
	err := keeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// GetTokenPair successful
	product := tokenPair.Name()
	getTokenPair := keeper.GetTokenPair(ctx, product)
	require.Equal(t, getTokenPair, tokenPair)

	userTokenPairs := keeper.GetUserTokenPairs(ctx, tokenPair.Owner)
	require.EqualValues(t, 1, len(userTokenPairs))

	// GetTokenPair failed
	getTokenPair = keeper.GetTokenPair(ctx, TestProductNotExist)
	require.Nil(t, getTokenPair)

	// GetTokenPairs from db
	getTokenPairs := keeper.GetTokenPairs(ctx)
	require.Equal(t, 1, len(getTokenPairs))

	// GetTokenPairs from cache111
	getTokenPairs = keeper.GetTokenPairs(ctx)
	require.Equal(t, 1, len(getTokenPairs))

	// GetTokenPairFromStore
	getTokenPair = keeper.GetTokenPairFromStore(ctx, product)
	require.Equal(t, getTokenPair, tokenPair)

}

func TestDeleteTokenPairByName(t *testing.T) {
	common.InitConfig()
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	keeper := testInput.DexKeeper
	tokenPair := getTestTokenPair()

	// SaveTokenPair
	err := keeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	userTokenPairs := keeper.GetUserTokenPairs(ctx, tokenPair.Owner)
	require.EqualValues(t, 1, len(userTokenPairs))

	// DeleteTokenPairByName successful
	product := tokenPair.Name()
	keeper.DeleteTokenPairByName(ctx, tokenPair.Owner, product)
	getTokenPair := keeper.GetTokenPair(ctx, product)
	require.Nil(t, getTokenPair)

	userTokenPairs = keeper.GetUserTokenPairs(ctx, tokenPair.Owner)
	require.EqualValues(t, 0, len(userTokenPairs))
}

func TestUpdateTokenPair(t *testing.T) {
	common.InitConfig()
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	keeper := testInput.DexKeeper
	tokenPair := getTestTokenPair()

	// SaveTokenPair
	err := keeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// UpdateTokenPair successful
	product := tokenPair.Name()
	blockHeight := tokenPair.BlockHeight
	updateTokenPair := tokenPair
	updateTokenPair.BlockHeight = blockHeight + 1
	keeper.UpdateTokenPair(ctx, product, updateTokenPair)
	getTokenPair := keeper.GetTokenPair(ctx, product)
	require.Equal(t, getTokenPair.BlockHeight, blockHeight+1)
}

func TestDeposit(t *testing.T) {
	common.InitConfig()
	testInput := createTestInputWithBalance(t, 2, 10000)
	ctx := testInput.Ctx
	keeper := testInput.DexKeeper
	tokenPair := getTestTokenPair()
	owner := testInput.TestAddrs[0]
	tokenPair.Owner = owner
	initDeposit := tokenPair.Deposits

	// SaveTokenPair
	err := keeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// Deposit successful
	product := tokenPair.Name()
	amount, err := sdk.ParseDecCoin("30" + sdk.DefaultBondDenom)
	require.Nil(t, err)
	err = keeper.Deposit(ctx, product, owner, amount)
	require.Nil(t, err)

	getTokenPair := keeper.GetTokenPair(ctx, product)
	require.Equal(t, getTokenPair.Deposits, initDeposit.Add(amount))

	// Deposit failed because of product not exist
	err = keeper.Deposit(ctx, TestProductNotExist, owner, amount)
	require.NotNil(t, err)

	// Deposit failed because of owner
	err = keeper.Deposit(ctx, product, testInput.TestAddrs[1], amount)
	require.NotNil(t, err)

	// Deposit failed because of invalid amount
	amountInvalid, err := sdk.ParseDecCoin("30" + common.TestToken)
	require.Nil(t, err)
	err = keeper.Deposit(ctx, product, owner, amountInvalid)
	require.NotNil(t, err)
}

func TestWithdraw(t *testing.T) {
	common.InitConfig()
	testInput := createTestInputWithBalance(t, 2, 10000)
	ctx := testInput.Ctx
	keeper := testInput.DexKeeper
	tokenPair := getTestTokenPair()
	owner := testInput.TestAddrs[0]
	tokenPair.Owner = owner
	initDeposit := tokenPair.Deposits
	keeper.SetParams(ctx, *types.DefaultParams())

	// SaveTokenPair
	err := keeper.SaveTokenPair(ctx, tokenPair)
	require.Nil(t, err)

	// Deposit successful
	product := tokenPair.Name()
	depositAmount, err := sdk.ParseDecCoin("30" + sdk.DefaultBondDenom)
	require.Nil(t, err)
	err = keeper.Deposit(ctx, product, owner, depositAmount)
	require.Nil(t, err)

	getTokenPair := keeper.GetTokenPair(ctx, product)
	require.Equal(t, getTokenPair.Deposits, initDeposit.Add(depositAmount))

	// Withdraw successful
	withdrawAmount, err := sdk.ParseDecCoin("10" + sdk.DefaultBondDenom)
	require.Nil(t, err)
	err = keeper.Withdraw(ctx, product, owner, withdrawAmount)
	require.Nil(t, err)
	getTokenPair = keeper.GetTokenPair(ctx, product)
	require.Equal(t, getTokenPair.Deposits, initDeposit.Add(depositAmount).Sub(withdrawAmount))

	// Withdraw failed because of product not exist
	err = keeper.Withdraw(ctx, TestProductNotExist, owner, withdrawAmount)
	require.NotNil(t, err)

	// Withdraw failed because of owner
	err = keeper.Withdraw(ctx, product, testInput.TestAddrs[1], withdrawAmount)
	require.NotNil(t, err)

	// Deposit failed because of invalid amount
	amountInvalid, err := sdk.ParseDecCoin("10" + common.TestToken)
	require.Nil(t, err)
	err = keeper.Withdraw(ctx, product, owner, amountInvalid)
	require.NotNil(t, err)

	// Withdraw failed because of deposits not enough
	err = keeper.Withdraw(ctx, product, owner, depositAmount)
	require.NotNil(t, err)
}

func TestGetTokenPairsOrdered(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	keeper := testInput.DexKeeper

	tokenPair0 := &types.TokenPair{
		BaseAssetSymbol:  "bToken0",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            testInput.TestAddrs[0],
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(50)),
		BlockHeight:      8,
	}

	tokenPair1 := &types.TokenPair{
		BaseAssetSymbol:  "bToken1",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            testInput.TestAddrs[0],
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100)),
		BlockHeight:      10,
	}

	tokenPair2 := &types.TokenPair{
		BaseAssetSymbol:  "bToken2",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            testInput.TestAddrs[0],
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(50)),
		BlockHeight:      9,
	}

	tokenPair3 := &types.TokenPair{
		BaseAssetSymbol:  "aToken0",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            testInput.TestAddrs[0],
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(50)),
		BlockHeight:      9,
	}

	// SaveTokenPair
	err := keeper.SaveTokenPair(ctx, tokenPair0)
	require.Nil(t, err)
	err = keeper.SaveTokenPair(ctx, tokenPair1)
	require.Nil(t, err)
	err = keeper.SaveTokenPair(ctx, tokenPair2)
	require.Nil(t, err)
	err = keeper.SaveTokenPair(ctx, tokenPair3)
	require.Nil(t, err)

	expectedSortedPairs := types.TokenPairs{tokenPair1, tokenPair0, tokenPair3, tokenPair2}

	// 1. compare deposit
	// 2. compare block height
	// 3. compare product name
	orderTokenPairs := keeper.GetTokenPairsOrdered(ctx)
	for idx, tokenPair := range orderTokenPairs {
		require.Equal(t, expectedSortedPairs[idx].Name(), tokenPair.Name())
	}
}

func TestSortProducts(t *testing.T) {
	testInput := createTestInputWithBalance(t, 1, 10000)
	ctx := testInput.Ctx
	keeper := testInput.DexKeeper

	tokenPair0 := &types.TokenPair{
		BaseAssetSymbol:  "bToken0",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            testInput.TestAddrs[0],
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(50)),
		BlockHeight:      8,
	}

	tokenPair1 := &types.TokenPair{
		BaseAssetSymbol:  "bToken1",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            testInput.TestAddrs[0],
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(100)),
		BlockHeight:      10,
	}

	tokenPair2 := &types.TokenPair{
		BaseAssetSymbol:  "bToken2",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            testInput.TestAddrs[0],
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(50)),
		BlockHeight:      9,
	}

	tokenPair3 := &types.TokenPair{
		BaseAssetSymbol:  "aToken0",
		QuoteAssetSymbol: common.NativeToken,
		Owner:            testInput.TestAddrs[0],
		Deposits:         sdk.NewDecCoin(sdk.DefaultBondDenom, sdk.NewInt(50)),
		BlockHeight:      9,
	}

	// SaveTokenPair
	err := keeper.SaveTokenPair(ctx, tokenPair0)
	require.Nil(t, err)
	err = keeper.SaveTokenPair(ctx, tokenPair1)
	require.Nil(t, err)
	err = keeper.SaveTokenPair(ctx, tokenPair2)
	require.Nil(t, err)
	err = keeper.SaveTokenPair(ctx, tokenPair3)
	require.Nil(t, err)

	unsoreProducts := []string{tokenPair0.Name(), tokenPair1.Name(), tokenPair2.Name(), tokenPair3.Name()}
	expectedSortedProducts := []string{tokenPair1.Name(), tokenPair0.Name(), tokenPair3.Name(), tokenPair2.Name()}

	// 1. compare deposit
	// 2. compare block height
	// 3. compare product name
	keeper.SortProducts(ctx, unsoreProducts)
	require.Equal(t, expectedSortedProducts, unsoreProducts)
}

func Test_IterateWithdrawInfo(t *testing.T) {
	testInput := createTestInputWithBalance(t, 2, 30)
	ctx := testInput.Ctx
	keeper := testInput.DexKeeper
	owner := testInput.TestAddrs[0]
	keeper.SetParams(ctx, *types.DefaultParams())

	withdrawAmount, err := sdk.ParseDecCoin("10" + sdk.DefaultBondDenom)
	require.Nil(t, err)
	withdrawInfo := types.WithdrawInfo{
		Owner:        owner,
		Deposits:     withdrawAmount,
		CompleteTime: time.Now().Add(types.DefaultWithdrawPeriod),
	}
	keeper.SetWithdrawInfo(ctx, withdrawInfo)
	expectWithdrawInfo, ok := keeper.GetWithdrawInfo(ctx, withdrawInfo.Owner)
	require.True(t, ok)
	require.True(t, withdrawInfo.Equal(expectWithdrawInfo))

	withdrawInfo.CompleteTime = withdrawInfo.CompleteTime.Add(types.DefaultWithdrawPeriod)
	keeper.SetWithdrawInfo(ctx, withdrawInfo)
	expectWithdrawInfo, ok = keeper.GetWithdrawInfo(ctx, withdrawInfo.Owner)
	require.True(t, ok)
	require.True(t, withdrawInfo.Equal(expectWithdrawInfo))

	var expectWithdrawInfos types.WithdrawInfos
	keeper.IterateWithdrawInfo(ctx, func(_ int64, withdrawInfo types.WithdrawInfo) (stop bool) {
		expectWithdrawInfos = append(expectWithdrawInfos, withdrawInfo)
		return false
	})
	require.Equal(t, 1, len(expectWithdrawInfos))
	require.True(t, expectWithdrawInfos[0].Equal(expectWithdrawInfo))
}

func TestKeeper_CheckTokenPairUnderDexDelist(t *testing.T) {
	testInput := createTestInputWithBalance(t, 2, 30)

	// fail case : the product is not exist
	isDelisting, err := testInput.DexKeeper.CheckTokenPairUnderDexDelist(testInput.Ctx, "no-product")
	require.Error(t, err)
	require.True(t, isDelisting)

	// save token pair
	tokenPair := getTestTokenPair()
	err = testInput.DexKeeper.SaveTokenPair(testInput.Ctx, tokenPair)
	require.Nil(t, err)

	// successful case
	isDelisting, err = testInput.DexKeeper.CheckTokenPairUnderDexDelist(testInput.Ctx, tokenPair.Name())
	require.Nil(t, err)
	require.Equal(t, isDelisting, tokenPair.Delisting)

}
