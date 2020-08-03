package dex

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/dex/types"
	ordertypes "github.com/okex/okchain/x/order/types"
	"github.com/stretchr/testify/require"
)

func TestInitGenesis(t *testing.T) {
	_, _, _, keeper, ctx := getMockTestCaseEvn(t)

	keeper.SetParams(ctx, *types.DefaultParams())
	params := keeper.GetParams(ctx)

	var tokenPairs []*types.TokenPair
	tokenPair := GetBuiltInTokenPair()
	tokenPairs = append(tokenPairs, tokenPair)

	var operators DEXOperators
	operators = append(operators, types.DEXOperator{
		Address:            tokenPair.Owner,
		HandlingFeeAddress: tokenPair.Owner,
		Website:            "http://www.okchain.com/operator.json",
		InitHeight:         100,
	})

	var withdrawInfos []WithdrawInfo
	now := time.Now()
	withdrawInfos = append(withdrawInfos, types.WithdrawInfo{
		Owner:        tokenPair.Owner,
		Deposits:     sdk.NewInt64DecCoin(tokenPair.BaseAssetSymbol, 1234),
		CompleteTime: now.Add(types.DefaultWithdrawPeriod),
	})

	lock := &ordertypes.ProductLock{
		BlockHeight:  666,
		Price:        sdk.NewDec(77),
		Quantity:     sdk.NewDec(88),
		BuyExecuted:  sdk.NewDec(99),
		SellExecuted: sdk.NewDec(66),
	}
	lockMap := ordertypes.NewProductLockMap()
	product := fmt.Sprintf("%s_%s", tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)
	lockMap.Data[product] = lock

	initGenesis := GenesisState{
		Params:         params,
		TokenPairs:     tokenPairs,
		WithdrawInfos:  withdrawInfos,
		ProductLocks:   *lockMap,
		Operators:      operators,
		MaxTokenPairID: 10,
	}
	require.NoError(t, ValidateGenesis(initGenesis))
	InitGenesis(ctx, keeper, initGenesis)
	require.Equal(t, initGenesis.Params, keeper.GetParams(ctx))
	require.Equal(t, initGenesis.TokenPairs, keeper.GetTokenPairs(ctx))
	require.Equal(t, initGenesis.MaxTokenPairID, keeper.GetMaxTokenPairID(ctx))
	require.Equal(t, initGenesis.ProductLocks, *keeper.LoadProductLocks(ctx))
	require.Equal(t, initGenesis.TokenPairs, keeper.GetUserTokenPairs(ctx, initGenesis.TokenPairs[0].Owner))

	var exportWithdrawInfos WithdrawInfos
	keeper.IterateWithdrawInfo(ctx, func(_ int64, withdrawInfo WithdrawInfo) (stop bool) {
		exportWithdrawInfos = append(exportWithdrawInfos, withdrawInfo)
		return false
	})
	require.True(t, initGenesis.WithdrawInfos.Equal(exportWithdrawInfos))

	var addr sdk.AccAddress
	keeper.IterateWithdrawAddress(ctx, initGenesis.WithdrawInfos[0].CompleteTime,
		func(_ int64, key []byte) (stop bool) {
			_, addr = types.SplitWithdrawTimeKey(key)
			return false
		})
	require.Equal(t, initGenesis.WithdrawInfos[0].Owner, addr)

	exportGenesis := ExportGenesis(ctx, keeper)
	require.Equal(t, initGenesis.Params, exportGenesis.Params)
	require.Equal(t, initGenesis.TokenPairs, exportGenesis.TokenPairs)
	require.True(t, initGenesis.WithdrawInfos.Equal(exportGenesis.WithdrawInfos))
	require.Equal(t, initGenesis.ProductLocks, exportGenesis.ProductLocks)
	require.Equal(t, initGenesis.MaxTokenPairID, exportGenesis.MaxTokenPairID)

	exportGenesis.Params.WithdrawPeriod = 55555
	exportGenesis.TokenPairs[0].ID = 66666
	exportGenesis.WithdrawInfos[0].CompleteTime = now.Add(2 * types.DefaultWithdrawPeriod)
	exportGenesis.ProductLocks.Data[product].BlockHeight = 123
	exportGenesis.MaxTokenPairID = 100

	_, _, _, newKeeper, newCtx := getMockTestCaseEvn(t)
	require.NoError(t, ValidateGenesis(exportGenesis))
	InitGenesis(newCtx, newKeeper, exportGenesis)
	require.Equal(t, exportGenesis.Params, newKeeper.GetParams(newCtx))
	require.Equal(t, exportGenesis.TokenPairs, newKeeper.GetTokenPairs(newCtx))
	require.Equal(t, exportGenesis.TokenPairs[0].ID, newKeeper.GetMaxTokenPairID(newCtx))
	require.Equal(t, exportGenesis.ProductLocks, *newKeeper.LoadProductLocks(newCtx))
	require.Equal(t, exportGenesis.TokenPairs, newKeeper.GetUserTokenPairs(newCtx, exportGenesis.TokenPairs[0].Owner))

	var exportWithdrawInfos1 WithdrawInfos
	newKeeper.IterateWithdrawInfo(newCtx, func(_ int64, withdrawInfo WithdrawInfo) (stop bool) {
		exportWithdrawInfos1 = append(exportWithdrawInfos1, withdrawInfo)
		return false
	})
	require.True(t, exportGenesis.WithdrawInfos.Equal(exportWithdrawInfos1))

	var addr1 sdk.AccAddress
	newKeeper.IterateWithdrawAddress(newCtx, exportGenesis.WithdrawInfos[0].CompleteTime,
		func(_ int64, key []byte) (stop bool) {
			_, addr1 = types.SplitWithdrawTimeKey(key)
			return false
		})
	require.Equal(t, exportGenesis.WithdrawInfos[0].Owner, addr1)

	newExportGenesis := ExportGenesis(newCtx, newKeeper)
	require.Equal(t, newExportGenesis.Params, newKeeper.GetParams(newCtx))
	require.Equal(t, newExportGenesis.TokenPairs, newKeeper.GetTokenPairs(newCtx))
	require.Equal(t, newExportGenesis.MaxTokenPairID, newKeeper.GetMaxTokenPairID(newCtx))
	var newExportWithdrawInfos WithdrawInfos
	newKeeper.IterateWithdrawInfo(newCtx, func(_ int64, withdrawInfo WithdrawInfo) (stop bool) {
		newExportWithdrawInfos = append(newExportWithdrawInfos, withdrawInfo)
		return false
	})
	require.True(t, newExportGenesis.WithdrawInfos.Equal(newExportWithdrawInfos))
	require.Equal(t, newExportGenesis.ProductLocks, *newKeeper.LoadProductLocks(newCtx))
}
