package keeper

import (
	"encoding/json"
	"math"
	"sort"
	"time"

	"github.com/okex/okexchain/x/ammswap"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/backend/types"
	"github.com/okex/okexchain/x/common"
	farm "github.com/okex/okexchain/x/farm/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// queryFarmPools returns pools of farm
func queryFarmPools(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QueryFarmPoolsParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}

	offset, limit := common.GetPage(queryParams.Page, queryParams.PerPage)
	if offset < 0 || limit < 0 {
		return nil, common.ErrInvalidPaginateParam(queryParams.Page, queryParams.PerPage)
	}

	// all farm pools
	allFarmPools := keeper.farmKeeper.GetFarmPools(ctx)
	// whitelist
	whitelist := keeper.farmKeeper.GetWhitelist(ctx)
	whitelistMap := make(map[string]bool, len(whitelist))
	for _, name := range whitelist {
		whitelistMap[name] = true
	}
	// farm pools
	var farmPools []farm.FarmPool
	switch queryParams.PoolType {
	case types.WhitelistFarmPool:
		for _, farmPool := range allFarmPools {
			if whitelistMap[farmPool.Name] {
				farmPools = append(farmPools, farmPool)
			}
		}
	case types.NormalFarmPool:
		for _, farmPool := range allFarmPools {
			if !whitelistMap[farmPool.Name] {
				farmPools = append(farmPools, farmPool)
			}
		}
	}

	allPoolStaked := sdk.ZeroDec()
	// response
	responseList := make(types.FarmResponseList, len(farmPools))
	for i, farmPool := range farmPools {
		// calculate total staked in dollars
		totalStakedDollars := keeper.farmKeeper.GetPoolLockedValue(ctx, farmPool)
		// calculate start at and finish at
		startAt := calculateFarmPoolStartAt(ctx, farmPool)
		finishAt := calculateFarmPoolFinishAt(ctx, keeper, farmPool, startAt)
		// calculate pool rate and farm apy
		yieldedInDay := farmPool.YieldedTokenInfos[0].AmountYieldedPerBlock.MulInt64(int64(types.BlocksPerDay))
		poolRate := sdk.NewDecCoinsFromDec(farmPool.YieldedTokenInfos[0].RemainingAmount.Denom, yieldedInDay)
		yieldedDollarsInDay := calculateAmountToDollars(ctx, keeper,
			sdk.NewDecCoinFromDec(farmPool.YieldedTokenInfos[0].RemainingAmount.Denom, yieldedInDay))
		apy := sdk.ZeroDec()
		if !totalStakedDollars.IsZero() {
			apy = yieldedDollarsInDay.Quo(totalStakedDollars).MulInt64(types.DaysInYear)
		}
		farmApy := sdk.NewDecCoinsFromDec(farmPool.YieldedTokenInfos[0].RemainingAmount.Denom, apy)
		status := getFarmPoolStatus(startAt, finishAt, farmPool)
		responseList[i] = types.FarmPoolResponse{
			PoolName:    farmPool.Name,
			LockSymbol:  farmPool.MinLockAmount.Denom,
			YieldSymbol: farmPool.YieldedTokenInfos[0].RemainingAmount.Denom,
			TotalStaked: totalStakedDollars,
			StartAt:     startAt,
			FinishAt:    finishAt,
			PoolRate:    poolRate,
			FarmApy:     farmApy,
			InWhitelist: whitelistMap[farmPool.Name],
			Status:      status,
		}

		// update allPoolStaked
		allPoolStaked = allPoolStaked.Add(totalStakedDollars)
	}

	// calculate pool rate and apy in whitelist
	if queryParams.PoolType == types.WhitelistFarmPool && allPoolStaked.IsPositive() && keeper.farmKeeper.GetParams(ctx).YieldNativeToken {
		moduleAcc := keeper.farmKeeper.SupplyKeeper().GetModuleAccount(ctx, farm.MintFarmingAccount)
		yieldedNativeTokenPerBlock := moduleAcc.GetCoins().AmountOf(sdk.DefaultBondDenom)
		yieldedNativeTokenPerDay := yieldedNativeTokenPerBlock.MulInt64(types.BlocksPerDay)
		for _, poolResponse := range responseList {
			nativeTokenRate := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, yieldedNativeTokenPerDay.Mul(poolResponse.TotalStaked.Quo(allPoolStaked)))
			poolResponse.PoolRate = poolResponse.PoolRate.Add(nativeTokenRate)
			nativeTokenToDollarsPerDay := calculateAmountToDollars(ctx, keeper, sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, nativeTokenRate.Amount))
			if !poolResponse.TotalStaked.IsZero() {
				nativeTokenApy := nativeTokenToDollarsPerDay.Quo(poolResponse.TotalStaked).MulInt64(types.DaysInYear)
				poolResponse.FarmApy = poolResponse.FarmApy.Add(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, nativeTokenApy))
			} else {
				poolResponse.FarmApy = poolResponse.FarmApy.Add(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.ZeroDec()))
			}
		}
	}
	// sort
	sort.Sort(responseList)

	// paginate
	total := len(responseList)
	switch {
	case total < offset:
		responseList = responseList[0:0]
	case total < offset+limit:
		responseList = responseList[offset:]
	default:
		responseList = responseList[offset : offset+limit]
	}

	// response
	var response *common.ListResponse
	if len(responseList) > 0 {
		response = common.GetListResponse(total, queryParams.Page, queryParams.PerPage, responseList)
	} else {
		response = common.GetEmptyListResponse(total, queryParams.Page, queryParams.PerPage)
	}

	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

// queryFarmDashboard returns dashboard of farm
func queryFarmDashboard(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QueryFarmDashboardParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}

	offset, limit := common.GetPage(queryParams.Page, queryParams.PerPage)
	if offset < 0 || limit < 0 {
		return nil, common.ErrInvalidPaginateParam(queryParams.Page, queryParams.PerPage)
	}

	address, err := sdk.AccAddressFromBech32(queryParams.Address)
	if err != nil {
		return nil, common.ErrCreateAddrFromBech32Failed(queryParams.Address, err.Error())
	}
	// staked pools
	stakedPools := keeper.farmKeeper.GetFarmPoolNamesForAccount(ctx, address)
	// whitelist
	whitelist := keeper.farmKeeper.GetWhitelist(ctx)
	whitelistMap := make(map[string]bool, len(whitelist))
	for _, name := range whitelist {
		whitelistMap[name] = true
	}
	// response
	var responseList types.FarmResponseList
	hasWhiteList := false
	for _, poolName := range stakedPools {
		farmPool, found := keeper.farmKeeper.GetFarmPool(ctx, poolName)
		if !found {
			continue
		}
		if whitelistMap[poolName] {
			hasWhiteList = true
		}
		// calculate staked in dollars and pool ratio
		poolRatio := sdk.ZeroDec()
		userStaked := sdk.ZeroDec()
		totalStakedDollars := keeper.farmKeeper.GetPoolLockedValue(ctx, farmPool)
		if lockInfo, found := keeper.farmKeeper.GetLockInfo(ctx, address, poolName); found {
			if !farmPool.TotalValueLocked.Amount.IsZero() {
				poolRatio = lockInfo.Amount.Amount.Quo(farmPool.TotalValueLocked.Amount)
				userStaked = poolRatio.Mul(totalStakedDollars)
			}
		}

		// calculate start at and finish at
		startAt := calculateFarmPoolStartAt(ctx, farmPool)
		finishAt := calculateFarmPoolFinishAt(ctx, keeper, farmPool, startAt)
		// calculate pool rate and farm apy
		yieldedInDay := farmPool.YieldedTokenInfos[0].AmountYieldedPerBlock.MulInt64(int64(types.BlocksPerDay))
		poolRate := sdk.NewDecCoinsFromDec(farmPool.YieldedTokenInfos[0].RemainingAmount.Denom, yieldedInDay)
		yieldedDollarsInDay := calculateAmountToDollars(ctx, keeper,
			sdk.NewDecCoinFromDec(farmPool.YieldedTokenInfos[0].RemainingAmount.Denom, yieldedInDay))
		apy := sdk.ZeroDec()
		if !totalStakedDollars.IsZero() {
			apy = yieldedDollarsInDay.Quo(totalStakedDollars).MulInt64(types.DaysInYear)
		}
		farmApy := sdk.NewDecCoinsFromDec(farmPool.YieldedTokenInfos[0].RemainingAmount.Denom, apy)

		// calculate total farmed and claim infos
		var unclaimed sdk.SysCoins
		var unclaimedInDollars sdk.SysCoins
		var claimed sdk.SysCoins
		var claimedInDollars sdk.SysCoins
		earning, err := keeper.farmKeeper.GetEarnings(ctx, farmPool.Name, address)
		if err == nil {
			unclaimed = earning.AmountYielded
			unclaimedInDollars = calculateSysCoinsInDollars(ctx, keeper, unclaimed)
		}
		farmDetails := generateFarmDetails(claimed, earning.AmountYielded)
		totalFarmed := calculateTotalFarmed(claimedInDollars, unclaimedInDollars)

		status := getFarmPoolStatus(startAt, finishAt, farmPool)
		responseList = append(responseList, types.FarmPoolResponse{
			PoolName:      farmPool.Name,
			LockSymbol:    farmPool.MinLockAmount.Denom,
			YieldSymbol:   farmPool.YieldedTokenInfos[0].RemainingAmount.Denom,
			TotalStaked:   userStaked,
			PoolRatio:     poolRatio,
			StartAt:       startAt,
			FinishAt:      finishAt,
			PoolRate:      poolRate,
			FarmApy:       farmApy,
			InWhitelist:   whitelistMap[poolName],
			FarmedDetails: farmDetails,
			TotalFarmed:   totalFarmed,
			Status:        status,
		})
	}

	// calculate whitelist apy
	if hasWhiteList && keeper.farmKeeper.GetParams(ctx).YieldNativeToken {
		moduleAcc := keeper.farmKeeper.SupplyKeeper().GetModuleAccount(ctx, farm.MintFarmingAccount)
		yieldedNativeTokenPerBlock := moduleAcc.GetCoins().AmountOf(sdk.DefaultBondDenom)
		yieldedNativeTokenPerDay := yieldedNativeTokenPerBlock.MulInt64(types.BlocksPerDay)
		whitelistTotalStaked := calculateWhitelistTotalStaked(ctx, keeper, whitelist)
		for _, poolResponse := range responseList {
			if !whitelistMap[poolResponse.PoolName] {
				continue
			}
			nativeTokenRate := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, yieldedNativeTokenPerDay.Mul(poolResponse.TotalStaked.Quo(whitelistTotalStaked)))
			poolResponse.PoolRate = poolResponse.PoolRate.Add(nativeTokenRate)
			nativeTokenToDollarsPerDay := calculateAmountToDollars(ctx, keeper, sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, nativeTokenRate.Amount))
			if !poolResponse.TotalStaked.IsZero() {
				nativeTokenApy := nativeTokenToDollarsPerDay.Quo(poolResponse.TotalStaked).MulInt64(types.DaysInYear)
				poolResponse.FarmApy = poolResponse.FarmApy.Add(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, nativeTokenApy))
			} else {
				poolResponse.FarmApy = poolResponse.FarmApy.Add(sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.ZeroDec()))
			}
		}
	}

	// sort
	sort.Sort(responseList)

	// paginate
	total := len(responseList)
	switch {
	case total < offset:
		responseList = responseList[0:0]
	case total < offset+limit:
		responseList = responseList[offset:]
	default:
		responseList = responseList[offset : offset+limit]
	}

	// response
	var response *common.ListResponse
	if len(responseList) > 0 {
		response = common.GetListResponse(total, queryParams.Page, queryParams.PerPage, responseList)
	} else {
		response = common.GetEmptyListResponse(total, queryParams.Page, queryParams.PerPage)
	}

	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

// queryFarmMaxApy returns max apy of farm pools
func queryFarmMaxApy(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	// whitelist
	whitelist := keeper.farmKeeper.GetWhitelist(ctx)
	apyMap := make(map[string]sdk.Dec, len(whitelist))
	allPoolStaked := sdk.ZeroDec()
	var responseList types.FarmResponseList
	for _, poolName := range whitelist {
		pool, found := keeper.farmKeeper.GetFarmPool(ctx, poolName)
		if !found {
			continue
		}
		totalStakedDollars := keeper.farmKeeper.GetPoolLockedValue(ctx, pool)
		yieldedInDay := pool.YieldedTokenInfos[0].AmountYieldedPerBlock.MulInt64(int64(types.BlocksPerDay))
		yieldedDollarsInDay := calculateAmountToDollars(ctx, keeper,
			sdk.NewDecCoinFromDec(pool.YieldedTokenInfos[0].RemainingAmount.Denom, yieldedInDay))
		apy := sdk.ZeroDec()
		if !totalStakedDollars.IsZero() {
			apy = yieldedDollarsInDay.Quo(totalStakedDollars).MulInt64(types.DaysInYear)
		}
		apyMap[poolName] = apy
		allPoolStaked = allPoolStaked.Add(totalStakedDollars)
		responseList = append(responseList, types.FarmPoolResponse{
			PoolName:    poolName,
			TotalStaked: totalStakedDollars,
		})
	}

	// calculate native token farmed apy
	if allPoolStaked.IsPositive() && keeper.farmKeeper.GetParams(ctx).YieldNativeToken {
		moduleAcc := keeper.farmKeeper.SupplyKeeper().GetModuleAccount(ctx, farm.MintFarmingAccount)
		yieldedNativeTokenPerBlock := moduleAcc.GetCoins().AmountOf(sdk.DefaultBondDenom)
		yieldedNativeTokenPerDay := yieldedNativeTokenPerBlock.MulInt64(types.BlocksPerDay)
		for _, poolResponse := range responseList {
			nativeTokenRate := sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, yieldedNativeTokenPerDay.Mul(poolResponse.TotalStaked.Quo(allPoolStaked)))
			nativeTokenToDollarsPerDay := calculateAmountToDollars(ctx, keeper, sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, nativeTokenRate.Amount))
			if !poolResponse.TotalStaked.IsZero() {
				nativeTokenApy := nativeTokenToDollarsPerDay.Quo(poolResponse.TotalStaked).MulInt64(types.DaysInYear)
				apyMap[poolResponse.PoolName] = apyMap[poolResponse.PoolName].Add(nativeTokenApy)
			}
		}
	}

	// max apy
	maxApy := sdk.ZeroDec()
	for _, apy := range apyMap {
		if apy.GT(maxApy) {
			maxApy = apy
		}
	}

	// response
	response := common.GetBaseResponse(maxApy)
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

// queryFarmStakedInfo returns farm staked info of the account
func queryFarmStakedInfo(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var queryParams types.QueryFarmStakedInfoParams
	err := keeper.cdc.UnmarshalJSON(req.Data, &queryParams)
	if err != nil {
		return nil, common.ErrUnMarshalJSONFailed(err.Error())
	}
	// validate params
	if queryParams.Address == "" {
		return nil, types.ErrAddressIsRequired()
	}
	address, err := sdk.AccAddressFromBech32(queryParams.Address)
	if err != nil {
		return nil, common.ErrCreateAddrFromBech32Failed(queryParams.Address, err.Error())
	}

	// query farm pool
	farmPool, found := keeper.farmKeeper.GetFarmPool(ctx, queryParams.PoolName)
	if !found {
		return nil, farm.ErrNoFarmPoolFound(queryParams.PoolName)
	}

	// query balance
	accountCoins := keeper.TokenKeeper.GetCoins(ctx, address)
	balance := accountCoins.AmountOf(farmPool.MinLockAmount.Denom)

	// locked info
	accountStaked := sdk.ZeroDec()
	if lockedInfo, found := keeper.farmKeeper.GetLockInfo(ctx, address, farmPool.Name); found {
		accountStaked = lockedInfo.Amount.Amount
	}

	// pool ratio
	poolRatio := sdk.ZeroDec()
	if !farmPool.TotalValueLocked.IsZero() {
		poolRatio = accountStaked.Quo(farmPool.TotalValueLocked.Amount)
	}

	// staked info
	stakedInfo := types.FarmStakedInfo{
		PoolName:        farmPool.Name,
		Balance:         balance,
		AccountStaked:   accountStaked,
		PoolTotalStaked: farmPool.TotalValueLocked.Amount,
		PoolRatio:       poolRatio,
	}
	// response
	response := common.GetBaseResponse(stakedInfo)
	bz, err := json.Marshal(response)
	if err != nil {
		return nil, common.ErrMarshalJSONFailed(err.Error())
	}
	return bz, nil
}

func generateFarmDetails(claimed sdk.SysCoins, unClaimed sdk.SysCoins) []types.FarmInfo {
	var farmDetails []types.FarmInfo
	claimedMap := make(map[string]sdk.Dec, len(claimed))
	for _, coin := range claimed {
		claimedMap[coin.Denom] = coin.Amount
	}
	for _, coin := range unClaimed {
		claimedAmount := claimedMap[coin.Denom]
		farmDetails = append(farmDetails, types.FarmInfo{
			Symbol:    coin.Denom,
			UnClaimed: coin.Amount,
			Claimed:   claimedAmount,
		})
	}
	return farmDetails
}

func calculateSysCoinsInDollars(ctx sdk.Context, keeper Keeper, coins sdk.SysCoins) sdk.SysCoins {
	result := sdk.SysCoins{}
	for _, coin := range coins {
		amountInDollars := calculateAmountToDollars(ctx, keeper, coin)
		result = append(result, sdk.NewDecCoinFromDec(coin.Denom, amountInDollars))
	}
	return result
}

// calculates totalLockedValue in dollar by usdk
func calculateAmountToDollars(ctx sdk.Context, keeper Keeper, amount sdk.SysCoin) sdk.Dec {
	if amount.Amount.IsZero() {
		return sdk.ZeroDec()
	}
	dollarAmount := sdk.ZeroDec()
	dollarQuoteToken := keeper.farmKeeper.GetParams(ctx).QuoteSymbol
	if amount.Denom == dollarQuoteToken {
		dollarAmount = amount.Amount
	} else {
		tokenPairName := ammswap.GetSwapTokenPairName(amount.Denom, dollarQuoteToken)
		if tokenPair, err := keeper.swapKeeper.GetSwapTokenPair(ctx, tokenPairName); err == nil {
			if tokenPair.BasePooledCoin.Denom == dollarQuoteToken && tokenPair.QuotePooledCoin.Amount.IsPositive() {
				dollarAmount = common.MulAndQuo(tokenPair.BasePooledCoin.Amount, amount.Amount, tokenPair.QuotePooledCoin.Amount)
			} else if tokenPair.BasePooledCoin.Amount.IsPositive() {
				dollarAmount = common.MulAndQuo(tokenPair.QuotePooledCoin.Amount, amount.Amount, tokenPair.BasePooledCoin.Amount)
			}
		}
	}
	return dollarAmount
}

func calculateFarmPoolStartAt(ctx sdk.Context, farmPool farm.FarmPool) int64 {
	if farmPool.YieldedTokenInfos[0].StartBlockHeightToYield == 0 {
		return 0
	}
	return time.Now().Unix() + (farmPool.YieldedTokenInfos[0].StartBlockHeightToYield-ctx.BlockHeight())*types.BlockInterval
}

func calculateFarmPoolFinishAt(ctx sdk.Context, keeper Keeper, farmPool farm.FarmPool, startAt int64) int64 {
	var finishAt int64
	updatedPool, _ := keeper.farmKeeper.CalculateAmountYieldedBetween(ctx, farmPool)
	if updatedPool.YieldedTokenInfos[0].RemainingAmount.Amount.IsPositive() && updatedPool.YieldedTokenInfos[0].AmountYieldedPerBlock.IsPositive() {
		if startAt > time.Now().Unix() {
			finishAt = startAt + updatedPool.YieldedTokenInfos[0].RemainingAmount.Amount.Quo(
				updatedPool.YieldedTokenInfos[0].AmountYieldedPerBlock).TruncateInt64()/int64(math.Pow10(sdk.Precision))*types.BlockInterval

		} else {
			finishAt = time.Now().Unix() + updatedPool.YieldedTokenInfos[0].RemainingAmount.Amount.Quo(
				updatedPool.YieldedTokenInfos[0].AmountYieldedPerBlock).TruncateInt64()/int64(math.Pow10(sdk.Precision))*types.BlockInterval
		}
	}
	return finishAt
}

func calculateWhitelistTotalStaked(ctx sdk.Context, keeper Keeper, whitelist []string) sdk.Dec {
	totalStaked := sdk.ZeroDec()
	for _, poolName := range whitelist {
		pool, found := keeper.farmKeeper.GetFarmPool(ctx, poolName)
		if !found {
			continue
		}
		poolValue := keeper.farmKeeper.GetPoolLockedValue(ctx, pool)
		totalStaked = totalStaked.Add(poolValue)
	}
	return totalStaked
}

func calculateTotalFarmed(claimed sdk.SysCoins, uncalimed sdk.SysCoins) sdk.Dec {
	sum := sdk.ZeroDec()
	for _, coin := range claimed {
		sum = sum.Add(coin.Amount)
	}
	for _, coin := range uncalimed {
		sum = sum.Add(coin.Amount)
	}
	return sum
}

func getFarmPoolStatus(startAt int64, finishAt int64, farmPool farm.FarmPool) types.FarmPoolStatus {
	if startAt == 0 {
		return types.FarmPoolCreated
	}
	if startAt > time.Now().Unix() && farmPool.YieldedTokenInfos[0].RemainingAmount.IsPositive() {
		return types.FarmPoolProvided
	}
	if time.Now().Unix() > startAt && time.Now().Unix() < finishAt {
		return types.FarmPoolYielded
	}
	return types.FarmPoolFinished
}
