package keeper

import (
	"fmt"
	"sort"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okchain/x/margin/types"
	"github.com/okex/okchain/x/params"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the margin store
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           *codec.Codec
	paramSubspace params.Subspace

	dexKeeper    types.DexKeeper
	supplyKeeper types.SupplyKeeper
	tokenKeeper  types.TokenKeeper
}

// NewKeeper creates a margin keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSubspace types.ParamSubspace, dexKeeper types.DexKeeper, tokenKeeper types.TokenKeeper, supplyKeeper types.SupplyKeeper) Keeper {
	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSubspace: paramSubspace.WithKeyTable(types.ParamKeyTable()),

		dexKeeper:    dexKeeper,
		tokenKeeper:  tokenKeeper,
		supplyKeeper: supplyKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetCDC() *codec.Codec {
	return k.cdc
}

// GetSupplyKeeper returns supply Keeper
func (k Keeper) GetSupplyKeeper() types.SupplyKeeper {
	return k.supplyKeeper
}

// GetSupplyKeeper returns token Keeper
func (k Keeper) GetTokenKeeper() types.TokenKeeper {
	return k.tokenKeeper
}

// GetDexKeeper returns dex Keeper
func (k Keeper) GetDexKeeper() types.DexKeeper {
	return k.dexKeeper
}

// GetParamSubspace returns paramSubspace
func (k Keeper) GetParamSubspace() params.Subspace {
	return k.paramSubspace
}

// SetParams sets inflation params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.GetParamSubspace().SetParamSet(ctx, &params)
}

// GetParams gets inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.GetParamSubspace().GetParamSet(ctx, &params)
	return params
}

// GetTradePair returns  the trade pair by product
func (k Keeper) GetTradePair(ctx sdk.Context, product string) *types.TradePair {
	var tradePair types.TradePair
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetTradePairKey(product))
	if bytes == nil {
		return nil
	}

	if k.cdc.UnmarshalBinaryBare(bytes, &tradePair) != nil {
		ctx.Logger().Error("decoding of token pair is failed", product)
		return nil
	}
	return &tradePair
}

// SetTradePair saves the trade pair to db
func (k Keeper) SetTradePair(ctx sdk.Context, tradePair *types.TradePair) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetTradePairKey(tradePair.Name)
	store.Set(key, k.cdc.MustMarshalBinaryBare(tradePair))
}

// DexDeposit deposits amount of tokens for a product
func (k Keeper) DexDeposit(ctx sdk.Context, from sdk.AccAddress, product string, amount sdk.DecCoin) sdk.Error {
	if amount.Denom != sdk.DefaultBondDenom {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to deposit because deposits only support %s token", sdk.DefaultBondDenom))
	}
	tradePair := k.GetTradePair(ctx, product)
	if tradePair == nil {
		tradePair = &types.TradePair{
			Owner:       from,
			Name:        product,
			Deposit:     amount,
			BlockHeight: ctx.BlockHeight(),
		}
	} else {
		tradePair.Deposit = tradePair.Deposit.Add(amount)
	}

	err := k.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, from, types.ModuleName, amount.ToCoins())
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to deposits because  insufficient deposit coins(need %s)", amount.ToCoins().String()))
	}
	k.SetTradePair(ctx, tradePair)
	return nil
}

// DexWithdraw withdraws amount of tokens from a product
func (k Keeper) DexWithdraw(ctx sdk.Context, product string, to sdk.AccAddress, amount sdk.DecCoin) sdk.Error {
	tradePair := k.GetTradePair(ctx, product)
	if tradePair == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to withdraws because non-exist product: %s", product))
	}

	if !tradePair.Owner.Equals(to) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("failed to withdraws because %s is not the owner of product:%s", to.String(), product))
	}

	if tradePair.Deposit.IsLT(amount) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to withdraws because deposits:%s is less than withdraw:%s", tradePair.Deposit.String(), amount.String()))
	}

	completeTime := ctx.BlockHeader().Time.Add(k.GetParams(ctx).WithdrawPeriod)
	// add withdraw info to store
	withdrawInfo := k.GetDexWithdrawInfo(ctx, to)
	if withdrawInfo == nil {
		withdrawInfo = &types.DexWithdrawInfo{
			Owner:        to,
			Deposits:     amount,
			CompleteTime: completeTime,
		}
	} else {
		k.DeleteDexWithdrawCompleteTimeAddress(ctx, withdrawInfo.CompleteTime, to)
		withdrawInfo.Deposits = withdrawInfo.Deposits.Add(amount)
		withdrawInfo.CompleteTime = completeTime
	}
	k.SetDexWithdrawInfo(ctx, withdrawInfo)
	k.SetDexWithdrawCompleteTimeAddress(ctx, completeTime, to)

	// update token pair
	tradePair.Deposit = tradePair.Deposit.Sub(amount)
	k.SetTradePair(ctx, tradePair)
	return nil
}

// GetDexWithdrawInfo returns withdraw info binding the addr
func (k Keeper) GetDexWithdrawInfo(ctx sdk.Context, addr sdk.AccAddress) *types.DexWithdrawInfo {
	bytes := ctx.KVStore(k.storeKey).Get(types.GetDexWithdrawKey(addr))
	if bytes == nil {
		return nil
	}
	var withdrawInfo *types.DexWithdrawInfo
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &withdrawInfo)
	return withdrawInfo
}

// SetDexWithdrawInfo sets withdraw address key with withdraw info
func (k Keeper) SetDexWithdrawInfo(ctx sdk.Context, withdrawInfo *types.DexWithdrawInfo) {
	key := types.GetDexWithdrawKey(withdrawInfo.Owner)
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(withdrawInfo)
	ctx.KVStore(k.storeKey).Set(key, bytes)
}

// SetDexWithdrawCompleteTimeAddress sets withdraw time key with empty []byte{} value
func (k Keeper) SetDexWithdrawCompleteTimeAddress(ctx sdk.Context, completeTime time.Time, addr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetDexWithdrawTimeAddressKey(completeTime, addr), []byte{})
}

// DeleteDexWithdrawCompleteTimeAddress deletes withdraw time key
func (k Keeper) DeleteDexWithdrawCompleteTimeAddress(ctx sdk.Context, timestamp time.Time, delAddr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetDexWithdrawTimeAddressKey(timestamp, delAddr))
}
func (k Keeper) withdrawTimeKeyIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	key := types.GetWithdrawTimeKey(endTime)
	return store.Iterator(types.DexWithdrawTimeKeyPrefix, sdk.PrefixEndBytes(key))
}

// IterateDexWithdrawAddress itreates withdraw time keys, and returns address
func (k Keeper) IterateDexWithdrawAddress(ctx sdk.Context, currentTime time.Time,
	fn func(index int64, key []byte) (stop bool)) {
	// iterate for all keys of (time+delAddr) from time 0 until the current time
	timeKeyIterator := k.withdrawTimeKeyIterator(ctx, currentTime)
	defer timeKeyIterator.Close()

	for i := int64(0); timeKeyIterator.Valid(); timeKeyIterator.Next() {
		key := timeKeyIterator.Key()
		if stop := fn(i, key); stop {
			break
		}
		i++
	}
}

func (k Keeper) deleteWithdrawInfo(ctx sdk.Context, addr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetDexWithdrawKey(addr))
}

// CompleteWithdraw completes withdrawing of addr
func (k Keeper) CompleteWithdraw(ctx sdk.Context, addr sdk.AccAddress) error {
	withdrawInfo := k.GetDexWithdrawInfo(ctx, addr)
	if withdrawInfo == nil {
		return sdk.ErrInvalidAddress(fmt.Sprintf("there is no withdrawing for address%s", addr.String()))
	}
	withdrawCoins := withdrawInfo.Deposits.ToCoins()
	err := k.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawInfo.Owner, withdrawCoins)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("withdraw error: %s, insufficient deposit coins(need %s)",
			err.Error(), withdrawCoins.String()))
	}
	k.deleteWithdrawInfo(ctx, addr)
	return nil
}

// DexSet sets params for a margin product
func (k Keeper) DexSet(ctx sdk.Context, address sdk.AccAddress, product string, maxLeverage sdk.Dec, borrowRate sdk.Dec, maintenanceMarginRatio sdk.Dec) sdk.Error {
	tradePair := k.GetTradePair(ctx, product)
	if tradePair == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to set because non-exist product: %s", product))
	}

	if !tradePair.Owner.Equals(address) {
		return sdk.ErrInvalidAddress(fmt.Sprintf("failed to set because %s is not the owner of product:%s", address.String(), product))
	}
	if maxLeverage.IsPositive() {
		tradePair.MaxLeverage = maxLeverage
	}

	if borrowRate.IsPositive() {
		tradePair.BorrowRate = borrowRate
	}
	if maintenanceMarginRatio.IsPositive() {
		tradePair.MaintenanceMarginRatio = maintenanceMarginRatio
	}
	k.SetTradePair(ctx, tradePair)
	return nil
}

// DexSave saves amount of tokens for borrowing
func (k Keeper) DexSave(ctx sdk.Context, address sdk.AccAddress, product string, amount sdk.DecCoins) sdk.Error {
	saving := k.GetSaving(ctx, product)
	if saving == nil {
		saving = amount
	} else {
		saving = saving.Add(amount)
	}

	err := k.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, address, types.ModuleName, amount)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to deposits because  insufficient coins(need %s)", amount.String()))
	}
	k.SetSaving(ctx, product, saving)
	return nil
}

// DexReturn returns amount of tokens for borrowing
func (k Keeper) DexReturn(ctx sdk.Context, address sdk.AccAddress, product string, amount sdk.DecCoins) sdk.Error {
	saving := k.GetSaving(ctx, product)
	if saving == nil || saving.IsAllLT(amount) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to deposits because insufficient coins saved(need %s)", amount.String()))
	}
	err := k.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, amount)
	if err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to deposits because insufficient coins saved(need %s)", amount.String()))
	}
	saving = saving.Sub(amount)
	k.SetSaving(ctx, product, saving)
	return nil
}

// GetSaving returns  the saving of product
func (k Keeper) GetSaving(ctx sdk.Context, product string) sdk.DecCoins {
	var saving sdk.DecCoins
	store := ctx.KVStore(k.storeKey)
	bytes := store.Get(types.GetSavingKey(product))
	if bytes == nil {
		return nil
	}

	if k.cdc.UnmarshalBinaryBare(bytes, &saving) != nil {
		ctx.Logger().Error("decoding of saving is failed", product)
		return nil
	}
	return saving
}

// SetSaving saves the saving of product to db
func (k Keeper) SetSaving(ctx sdk.Context, product string, amount sdk.DecCoins) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetSavingKey(product)
	store.Set(key, k.cdc.MustMarshalBinaryBare(amount))
}

// GetAccount returns the account from db
func (k Keeper) GetAccount(ctx sdk.Context, address sdk.AccAddress, product string) *types.Account {
	var account *types.Account
	bytes := ctx.KVStore(k.storeKey).Get(types.GetAccountAddressProductKey(address, product))
	if bytes == nil {
		return nil
	}
	k.cdc.UnmarshalBinaryBare(bytes, &account)
	return account
}

func (k Keeper) SetAccount(ctx sdk.Context, address sdk.AccAddress, account *types.Account) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAccountAddressProductKey(address, account.Product)
	store.Set(key, k.cdc.MustMarshalBinaryBare(account))
}

// Deposit transfer token from spot account to margin account
func (k Keeper) Deposit(ctx sdk.Context, address sdk.AccAddress, product string, amount sdk.DecCoins) sdk.Error {
	account := k.GetAccount(ctx, address, product)
	// account info has exist
	if account != nil {
		account.Available = account.Available.Add(amount)
	} else {
		account = &types.Account{
			Product:   product,
			Available: amount,
			Locked:    sdk.DecCoins{},
			Borrowed:  sdk.DecCoins{},
			Interest:  sdk.DecCoins{},
		}
	}

	if err := k.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, address, types.ModuleName, amount); err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to deposits because  insufficient deposit coins(need %s)", amount.String()))
	}

	k.SetAccount(ctx, address, account)
	return nil
}

// Withdraw withdraws from margin account to address
func (k Keeper) Withdraw(ctx sdk.Context, address sdk.AccAddress, product string, amount sdk.DecCoins) sdk.Error {
	account := k.GetAccount(ctx, address, product)
	if account == nil {
		return types.ErrAccountNotExist(types.Codespace, fmt.Sprintf("failed to withdraw beacuse the margin account not exists "))
	}

	if !account.Borrowed.IsZero() {
		return types.ErrNotAllowed(types.Codespace, "should repay borrowed coins before withdraw")
	}

	if amount.IsAnyGT(account.Available) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to withdraw because insufficient coins saved(need %s)", amount.String()))
	}

	account.Available = account.Available.Sub(amount)

	if err := k.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, address, amount); err != nil {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to withdraw because insufficient coins saved(need %s)", amount.String()))
	}

	k.SetAccount(ctx, address, account)
	return nil
}

// Borrow record the loan information of an account under the margin trading pair
func (k Keeper) Borrow(ctx sdk.Context, address sdk.AccAddress, tradePair *types.TradePair, deposit sdk.DecCoin, leverage sdk.Dec) sdk.Error {
	account := k.GetAccount(ctx, address, tradePair.Name)
	if account == nil {
		return types.ErrAccountNotExist(types.Codespace, fmt.Sprintf("margin account not exists"))
	}

	borrowAmount := sdk.DecCoin{Denom: deposit.Denom, Amount: deposit.Amount.Mul(leverage.Sub(sdk.NewDec(1)))}
	maxCanBorrow := account.MaxCanBorrow(deposit.Denom, tradePair.MaxLeverage)
	if maxCanBorrow.IsLT(borrowAmount) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to borrow because insufficient coins, max can borrow: %s", maxCanBorrow.String()))
	}

	// sub saving
	saving := k.GetSaving(ctx, tradePair.Name)
	if saving == nil || !saving.IsAllGTE(sdk.NewCoins(borrowAmount)) {
		return sdk.ErrInsufficientCoins(fmt.Sprintf("failed to borrow because insufficient coins saved(need %s)", borrowAmount.String()))
	}
	saving = saving.Sub(sdk.NewCoins(borrowAmount))
	k.SetSaving(ctx, tradePair.Name, saving)

	// add borrow
	borrowInfo := k.GetBorrowInfo(ctx, address, tradePair.Name, uint64(ctx.BlockHeight()))
	if borrowInfo == nil {
		borrowInfo = &types.BorrowInfo{
			Address:      address,
			Product:      tradePair.Name,
			BorrowAmount: sdk.NewCoins(borrowAmount),
			BlockHeight:  ctx.BlockHeight(),
			Rate:         tradePair.BorrowRate,
			Leverage:     leverage,
		}
	} else {
		borrowInfo.BorrowAmount = borrowInfo.BorrowAmount.Add(sdk.NewCoins(borrowAmount))
	}
	k.SetBorrowInfo(ctx, borrowInfo)

	// add calculate interest key
	nextCalculateTime := ctx.BlockTime().Add(k.GetParams(ctx).InterestPeriod)
	k.SetCalculateInterestKey(ctx, nextCalculateTime, address, tradePair.Name, uint64(ctx.BlockHeight()))

	// update account
	interest := sdk.NewCoins(borrowAmount).MulDec(tradePair.BorrowRate)
	account.Borrowed = account.Borrowed.Add(sdk.NewCoins(borrowAmount))
	account.Available = account.Available.Add(sdk.NewCoins(borrowAmount))
	account.Interest = account.Interest.Add(interest)
	k.SetAccount(ctx, address, account)
	return nil
}

// nolint
func (k Keeper) GetBorrowInfo(ctx sdk.Context, address sdk.AccAddress, product string, blockHeight uint64) *types.BorrowInfo {
	key := types.GetBorrowInfoKey(address, product, blockHeight)
	return k.GetBorrowInfoByKey(ctx, key)
}

// nolint
func (k Keeper) GetBorrowInfoByKey(ctx sdk.Context, key []byte) *types.BorrowInfo {
	bytes := ctx.KVStore(k.storeKey).Get(key)
	if bytes == nil {
		return nil
	}
	var borrowInfo *types.BorrowInfo
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &borrowInfo)
	return borrowInfo
}

// SetBorrowInfo set or update the borrowInfo to db
func (k Keeper) SetBorrowInfo(ctx sdk.Context, borrowInfo *types.BorrowInfo) {
	key := types.GetBorrowInfoKey(borrowInfo.Address, borrowInfo.Product, uint64(ctx.BlockHeight()))
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(borrowInfo)
	ctx.KVStore(k.storeKey).Set(key, bytes)
}

func (k Keeper) deleteBorrowInfo(ctx sdk.Context, borrowInfo *types.BorrowInfo) {
	key := types.GetBorrowInfoKey(borrowInfo.Address, borrowInfo.Product, uint64(ctx.BlockHeight()))
	ctx.KVStore(k.storeKey).Delete(key)
}

// SetCalculateInterestKey use the interest calculation time as the key.
func (k Keeper) SetCalculateInterestKey(ctx sdk.Context, calculateTime time.Time, address sdk.AccAddress,
	product string, blockHeight uint64) {
	borrowInfoKey := types.GetBorrowInfoKey(address, product, blockHeight)
	ctx.KVStore(k.storeKey).Set(types.GetCalculateInterestKey(calculateTime, borrowInfoKey), []byte{})
}

// DeleteCalculateInterestKey delete the key when all the borrowings have been repaid
func (k Keeper) DeleteCalculateInterestKey(ctx sdk.Context, timestamp time.Time, borrowInfoKey []byte) {
	ctx.KVStore(k.storeKey).Delete(types.GetCalculateInterestKey(timestamp, borrowInfoKey))
}

// IterateCalculateInterest iterate through the borrowing information to calculate interest at EndBlock
func (k Keeper) IterateCalculateInterest(ctx sdk.Context, currentTime time.Time,
	fn func(index int64, key []byte) (stop bool)) {
	// iterate for all keys of (time+ borrowInfoKey) from time 0 until the current time
	timeKeyIterator := k.calculateTimeKeyIterator(ctx, currentTime)
	defer timeKeyIterator.Close()
	for i := int64(0); timeKeyIterator.Valid(); timeKeyIterator.Next() {
		key := timeKeyIterator.Key()
		if stop := fn(i, key); stop {
			break
		}
		i++
	}
}

//  calculateTimeKeyIterator traversal to get obtain loan key
func (k Keeper) calculateTimeKeyIterator(ctx sdk.Context, calculateTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	key := types.GetCalculateInterestTimeKey(calculateTime)
	return store.Iterator(types.CalculateInterestKeyPrefix, sdk.PrefixEndBytes(key))
}

// Repay repays the borrowing of product
// repay precedence: 1. return interest 2. repay borrowing which rate is greater 3. repay borrowing which borrowed earlier
func (k Keeper) Repay(ctx sdk.Context, address sdk.AccAddress, tradePair *types.TradePair, amount sdk.DecCoin) sdk.Error {
	account := k.GetAccount(ctx, address, tradePair.Name)
	if account == nil {
		return types.ErrAccountNotExist(types.Codespace, fmt.Sprintf("failed to repay"))
	}

	if account.Borrowed.AmountOf(amount.Denom).IsZero() {
		return sdk.ErrInvalidCoins(fmt.Sprintf("repay amount:%s mismatch borrowed coins:%s", amount.String(), account.Borrowed.String()))
	}

	denom := amount.Denom
	actualAmount := amount.Amount
	// when amount is greater than borrowed + interest
	if amount.Amount.GT(account.Borrowed.AmountOf(denom).Add(account.Interest.AmountOf(denom))) {
		actualAmount = account.Borrowed.AmountOf(denom).Add(account.Interest.AmountOf(denom))
	}
	// repay to saving, update saving
	saving := k.GetSaving(ctx, tradePair.Name)
	saving.AmountOf(denom) = saving.AmountOf(denom).Sub(actualAmount)
	k.SetSaving(ctx, tradePair.Name, saving)

	// only repay interest & update account
	if account.Interest.AmountOf(denom).GTE(actualAmount) {
		// update account
		account.Interest.AmountOf(denom) = account.Interest.AmountOf(denom).Sub(actualAmount)
		k.SetAccount(ctx, address, account)
		return nil
	}

	// update account
	remainAmount := actualAmount.Sub(account.Interest.AmountOf(denom))
	account.Interest.AmountOf(denom) = sdk.ZeroDec()
	account.Available.AmountOf(denom) = account.Available.AmountOf(denom).Sub(remainAmount)
	account.Borrowed.AmountOf(denom) = account.Borrowed.AmountOf(denom).Sub(remainAmount)
	k.SetAccount(ctx, address, account)

	// repay borrowing & update borrowInfo
	borrowInfoList := k.GetBorrowInfoList(ctx, address, tradePair.Name)
	sort.Sort(borrowInfoList)
	for _, borrowInfo := range borrowInfoList {
		if borrowInfo.BorrowAmount.AmountOf(denom).GT(remainAmount) {
			borrowInfo.BorrowAmount.AmountOf(denom) = borrowInfo.BorrowAmount.AmountOf(denom).Sub(remainAmount)
			k.SetBorrowInfo(ctx, borrowInfo)
			break
		}
		remainAmount = remainAmount.Sub(borrowInfo.BorrowAmount.AmountOf(denom))
		k.deleteBorrowInfo(ctx, borrowInfo)
	}
	return nil
}

// GetBorrowInfoList  returns all borrowInfos
func (k Keeper) GetBorrowInfoList(ctx sdk.Context, address sdk.AccAddress, product string) (borrowInfoList types.BorrowInfoList) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetBorrowInfoProductKey(address, product))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var borrowInfo types.BorrowInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &borrowInfo)
		borrowInfoList = append(borrowInfoList, &borrowInfo)
	}
	return
}

// GetAccounts return all margin accunts of address
func (k Keeper) GetAccounts(ctx sdk.Context, address sdk.AccAddress) (accounts []*types.Account) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetAccountAddressKey(address))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var account types.Account
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &account)
		accounts = append(accounts, &account)
	}
	return
}
