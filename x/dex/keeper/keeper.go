package keeper

import (
	"fmt"
	"sort"
	"time"

	"github.com/okex/okexchain/x/stream/exported"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/dex/types"
	"github.com/okex/okexchain/x/params"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	supplyKeeper      SupplyKeeper
	feeCollectorName  string // name of the FeeCollector ModuleAccount
	tokenKeeper       TokenKeeper
	stakingKeeper     StakingKeeper         // The reference to the staking keeper to check whether proposer is  validator
	bankKeeper        BankKeeper            // The reference to the bank keeper to check whether proposer can afford  proposal deposit
	govKeeper         GovKeeper             // The reference to the gov keeper to handle proposal
	observerKeeper    exported.StreamKeeper // The reference to the stream keeper
	storeKey          sdk.StoreKey
	tokenPairStoreKey sdk.StoreKey
	paramSubspace     params.Subspace // The reference to the Paramstore to get and set gov modifiable params
	cdc               *codec.Codec    // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the token Keeper
func NewKeeper(feeCollectorName string, supplyKeeper SupplyKeeper, dexParamsSubspace params.Subspace, tokenKeeper TokenKeeper,
	stakingKeeper StakingKeeper, bankKeeper BankKeeper, storeKey, tokenPairStoreKey sdk.StoreKey, cdc *codec.Codec) Keeper {

	k := Keeper{
		tokenKeeper:       tokenKeeper,
		feeCollectorName:  feeCollectorName,
		supplyKeeper:      supplyKeeper,
		stakingKeeper:     stakingKeeper,
		bankKeeper:        bankKeeper,
		paramSubspace:     dexParamsSubspace.WithKeyTable(types.ParamKeyTable()),
		storeKey:          storeKey,
		tokenPairStoreKey: tokenPairStoreKey,
		cdc:               cdc,
	}

	return k
}

// GetSupplyKeeper returns supply Keeper
func (k Keeper) GetSupplyKeeper() SupplyKeeper {
	return k.supplyKeeper
}

// GetBankKeeper returns bank Keeper
func (k Keeper) GetBankKeeper() BankKeeper {
	return k.bankKeeper
}

// GetFeeCollector returns feeCollectorName
func (k Keeper) GetFeeCollector() string {
	return k.feeCollectorName
}

// GetTokenKeeper returns token Keeper
func (k Keeper) GetTokenKeeper() TokenKeeper {
	return k.tokenKeeper
}

func (k Keeper) deleteUserTokenPair(ctx sdk.Context, owner sdk.AccAddress, pair string) {
	store := ctx.KVStore(k.tokenPairStoreKey)
	store.Delete(types.GetUserTokenPairAddress(owner, pair))
}

// SaveTokenPair saves the token pair to db
// key is base:quote
func (k Keeper) SaveTokenPair(ctx sdk.Context, tokenPair *types.TokenPair) error {
	store := ctx.KVStore(k.tokenPairStoreKey)

	maxTokenPairID := k.GetMaxTokenPairID(ctx)
	// list new tokenPair
	if tokenPair.ID == 0 {
		tokenPair.ID = maxTokenPairID + 1
	}

	// update maxTokenPairID to db
	// to load exported data from genesis file.
	if tokenPair.ID > maxTokenPairID {
		k.SetMaxTokenPairID(ctx, tokenPair.ID)
	}

	keyPair := tokenPair.BaseAssetSymbol + "_" + tokenPair.QuoteAssetSymbol
	store.Set(types.GetTokenPairAddress(keyPair), k.cdc.MustMarshalBinaryBare(tokenPair))
	store.Set(types.GetUserTokenPairAddress(tokenPair.Owner, keyPair), []byte{})

	if k.observerKeeper != nil {
		k.observerKeeper.OnAddNewTokenPair(ctx, tokenPair)
	}

	return nil
}

// GetTokenPair gets the token pair by product
func (k Keeper) GetTokenPair(ctx sdk.Context, product string) *types.TokenPair {
	var tokenPair *types.TokenPair

	store := ctx.KVStore(k.tokenPairStoreKey)
	bytes := store.Get(types.GetTokenPairAddress(product))
	if bytes == nil {
		return nil
	}

	if k.cdc.UnmarshalBinaryBare(bytes, &tokenPair) != nil {
		ctx.Logger().Info("decoding of token pair is failed", product)
		return nil
	}

	return tokenPair
}

// GetTokenPairFromStore returns token pair from store without cache
func (k Keeper) GetTokenPairFromStore(ctx sdk.Context, product string) *types.TokenPair {
	var tokenPair types.TokenPair
	store := ctx.KVStore(k.tokenPairStoreKey)
	bytes := store.Get(types.GetTokenPairAddress(product))
	if bytes == nil {
		return nil
	}
	if k.cdc.UnmarshalBinaryBare(bytes, &tokenPair) != nil {
		ctx.Logger().Info("decoding of token pair is failed", product)
		return nil
	}

	return &tokenPair
}

// GetTokenPairs returns all token pairs from store without cache
func (k Keeper) GetTokenPairs(ctx sdk.Context) (tokenPairs []*types.TokenPair) {
	store := ctx.KVStore(k.tokenPairStoreKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenPairKey)
	defer iter.Close()
	for iter.Valid() {
		var tokenPair types.TokenPair
		tokenPairBytes := iter.Value()
		k.cdc.MustUnmarshalBinaryBare(tokenPairBytes, &tokenPair)
		tokenPairs = append(tokenPairs, &tokenPair)
		iter.Next()
	}

	return tokenPairs
}

// GetUserTokenPairs returns all token pairs belong to an account from store
func (k Keeper) GetUserTokenPairs(ctx sdk.Context, owner sdk.AccAddress) (tokenPairs []*types.TokenPair) {
	store := ctx.KVStore(k.tokenPairStoreKey)
	userTokenPairPrefix := types.GetUserTokenPairAddressPrefix(owner)
	prefixLen := len(userTokenPairPrefix)

	iter := sdk.KVStorePrefixIterator(store, userTokenPairPrefix)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		tokenPairName := string(key[prefixLen:])

		tokenPair := k.GetTokenPairFromStore(ctx, tokenPairName)
		if tokenPair != nil {
			tokenPairs = append(tokenPairs, tokenPair)
		}
	}

	return tokenPairs
}

// DeleteTokenPairByName deletes the token pair by name
func (k Keeper) DeleteTokenPairByName(ctx sdk.Context, owner sdk.AccAddress, product string) {
	// get store
	store := ctx.KVStore(k.tokenPairStoreKey)
	// delete the token pair from the store
	store.Delete(types.GetTokenPairAddress(product))
	// remove the user-tokenpair relationship
	k.deleteUserTokenPair(ctx, owner, product)

	if k.observerKeeper != nil {
		k.observerKeeper.OnTokenPairUpdated(ctx)
	}
}

func (k Keeper) updateUserTokenPair(ctx sdk.Context, product string, owner, to sdk.AccAddress) {
	store := ctx.KVStore(k.tokenPairStoreKey)
	store.Delete(types.GetUserTokenPairAddress(owner, product))
	store.Set(types.GetUserTokenPairAddress(to, product), []byte{})
}

// UpdateTokenPair updates token pair in the store and the cache
func (k Keeper) UpdateTokenPair(ctx sdk.Context, product string, tokenPair *types.TokenPair) {
	store := ctx.KVStore(k.tokenPairStoreKey)
	store.Set(types.GetTokenPairAddress(product), k.cdc.MustMarshalBinaryBare(*tokenPair))

	if k.observerKeeper != nil {
		k.observerKeeper.OnTokenPairUpdated(ctx)
	}
}

// CheckTokenPairUnderDexDelist checks if token pair is under delist. for x/order: It's not allowed to place an order about the tokenpair under dex delist
func (k Keeper) CheckTokenPairUnderDexDelist(ctx sdk.Context, product string) (isDelisting bool, err error) {
	tp := k.GetTokenPair(ctx, product)
	if tp != nil {
		isDelisting = tp.Delisting
	} else {
		isDelisting = true
		msg := fmt.Sprintf("product %s doesn't exist", product)
		err = types.ErrTokenPairNotFound(msg)
	}
	return isDelisting, err
}

// Deposit deposits amount of tokens for a product
func (k Keeper) Deposit(ctx sdk.Context, product string, from sdk.AccAddress, amount sdk.DecCoin) sdk.Error {
	tokenPair := k.GetTokenPair(ctx, product)
	if tokenPair == nil {
		return types.ErrInvalidTokenPair(product)
	}

	if !tokenPair.Owner.Equals(from) {
		return types.ErrMustTokenPairOwner(from.String(), product)
	}

	if amount.Denom != sdk.DefaultBondDenom {
		return types.ErrDepositOnlySupportDefaultBondDenom(sdk.DefaultBondDenom)
	}

	depositCoins := amount.ToCoins()
	err := k.GetSupplyKeeper().SendCoinsFromAccountToModule(ctx, from, types.ModuleName, depositCoins)
	if err != nil {
		return types.ErrInsufficientDepositCoins(err.Error(), depositCoins.String())
	}

	tokenPair.Deposits = tokenPair.Deposits.Add(amount)
	k.UpdateTokenPair(ctx, product, tokenPair)
	return nil
}

// Withdraw withdraws amount of tokens from a product
func (k Keeper) Withdraw(ctx sdk.Context, product string, to sdk.AccAddress, amount sdk.DecCoin) sdk.Error {
	tokenPair := k.GetTokenPair(ctx, product)
	if tokenPair == nil {
		return types.ErrInvalidTokenPair(product)
	}

	if !tokenPair.Owner.Equals(to) {
		return types.ErrMustTokenPairOwner(to.String(), product)
	}

	if amount.Denom != sdk.DefaultBondDenom {
		return types.ErrWithdrawOnlySupportDefaultBondDenom(sdk.DefaultBondDenom)
	}

	if tokenPair.Deposits.IsLT(amount) {
		return types.ErrInsufficientWithdrawCoins(tokenPair.Deposits.String(), amount.String())
	}

	completeTime := ctx.BlockHeader().Time.Add(k.GetParams(ctx).WithdrawPeriod)
	// add withdraw info to store
	withdrawInfo, ok := k.GetWithdrawInfo(ctx, to)
	if !ok {
		withdrawInfo = types.WithdrawInfo{
			Owner:        to,
			Deposits:     amount,
			CompleteTime: completeTime,
		}
	} else {
		k.DeleteWithdrawCompleteTimeAddress(ctx, withdrawInfo.CompleteTime, to)
		withdrawInfo.Deposits = withdrawInfo.Deposits.Add(amount)
		withdrawInfo.CompleteTime = completeTime
	}
	k.SetWithdrawInfo(ctx, withdrawInfo)
	k.SetWithdrawCompleteTimeAddress(ctx, completeTime, to)

	// update token pair
	tokenPair.Deposits = tokenPair.Deposits.Sub(amount)
	k.UpdateTokenPair(ctx, product, tokenPair)
	return nil
}

// GetTokenPairsOrdered returns token pairs ordered by product
func (k Keeper) GetTokenPairsOrdered(ctx sdk.Context) types.TokenPairs {
	var result types.TokenPairs
	tokenPairs := k.GetTokenPairs(ctx)
	for _, tp := range tokenPairs {
		result = append(result, tp)
	}
	sort.Sort(result)
	return result
}

// SortProducts sorts products
func (k Keeper) SortProducts(ctx sdk.Context, products []string) {
	tokenPairs := make(types.TokenPairs, 0, len(products))
	for _, product := range products {
		tokenPair := k.GetTokenPair(ctx, product)
		if tokenPair != nil {
			tokenPairs = append(tokenPairs, tokenPair)
		}
	}
	sort.Sort(tokenPairs)

	for i, tokenPair := range tokenPairs {
		products[i] = fmt.Sprintf("%s_%s", tokenPair.BaseAssetSymbol, tokenPair.QuoteAssetSymbol)
	}
}

// GetParams gets inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.GetParamSubspace().GetParamSet(ctx, &params)
	return params
}

// SetParams sets inflation params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.GetParamSubspace().SetParamSet(ctx, &params)
}

// GetParamSubspace returns paramSubspace
func (k Keeper) GetParamSubspace() params.Subspace {
	return k.paramSubspace
}

// TransferOwnership transfers ownership of product
func (k Keeper) TransferOwnership(ctx sdk.Context, product string, from sdk.AccAddress, to sdk.AccAddress) sdk.Error {
	tokenPair := k.GetTokenPair(ctx, product)
	if tokenPair == nil {
		return types.ErrTokenPairNotFound(product)
	}

	if !tokenPair.Owner.Equals(from) {
		return types.ErrMustTokenPairOwner(from.String(), product)
	}

	// Withdraw
	if tokenPair.Deposits.IsPositive() {
		if err := k.Withdraw(ctx, product, from, tokenPair.Deposits); err != nil {
			return types.ErrWithdrawDepositsError(tokenPair.Deposits.String(), err.Error())
		}
	}

	// transfer ownership
	tokenPair.Owner = to
	tokenPair.Deposits = types.DefaultTokenPairDeposit
	k.UpdateTokenPair(ctx, product, tokenPair)
	k.updateUserTokenPair(ctx, product, from, to)

	return nil
}

// GetWithdrawInfo returns withdraw info binding the addr
func (k Keeper) GetWithdrawInfo(ctx sdk.Context, addr sdk.AccAddress) (withdrawInfo types.WithdrawInfo, ok bool) {
	bytes := ctx.KVStore(k.storeKey).Get(types.GetWithdrawAddressKey(addr))
	if bytes == nil {
		return
	}

	k.cdc.MustUnmarshalBinaryLengthPrefixed(bytes, &withdrawInfo)
	return withdrawInfo, true
}

// SetWithdrawInfo sets withdraw address key with withdraw info
func (k Keeper) SetWithdrawInfo(ctx sdk.Context, withdrawInfo types.WithdrawInfo) {
	key := types.GetWithdrawAddressKey(withdrawInfo.Owner)
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(withdrawInfo)
	ctx.KVStore(k.storeKey).Set(key, bytes)
}

func (k Keeper) deleteWithdrawInfo(ctx sdk.Context, addr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetWithdrawAddressKey(addr))
}

func (k Keeper) withdrawTimeKeyIterator(ctx sdk.Context, endTime time.Time) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	key := types.GetWithdrawTimeKey(endTime)
	return store.Iterator(types.WithdrawTimeKeyPrefix, sdk.PrefixEndBytes(key))
}

// SetWithdrawCompleteTimeAddress sets withdraw time key with empty []byte{} value
func (k Keeper) SetWithdrawCompleteTimeAddress(ctx sdk.Context, completeTime time.Time, addr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(types.GetWithdrawTimeAddressKey(completeTime, addr), []byte{})
}

// DeleteWithdrawCompleteTimeAddress deletes withdraw time key
func (k Keeper) DeleteWithdrawCompleteTimeAddress(ctx sdk.Context, timestamp time.Time, delAddr sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Delete(types.GetWithdrawTimeAddressKey(timestamp, delAddr))
}

// IterateWithdrawInfo iterates withdraw address key， and returns withdraw info
func (k Keeper) IterateWithdrawInfo(ctx sdk.Context, fn func(index int64, withdrawInfo types.WithdrawInfo) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WithdrawAddressKeyPrefix)
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		var withdrawInfo types.WithdrawInfo
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &withdrawInfo)
		if stop := fn(i, withdrawInfo); stop {
			break
		}
		i++
	}
}

// IterateWithdrawAddress itreates withdraw time keys, and returns address
func (k Keeper) IterateWithdrawAddress(ctx sdk.Context, currentTime time.Time,
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

// CompleteWithdraw completes withdrawing of addr
func (k Keeper) CompleteWithdraw(ctx sdk.Context, addr sdk.AccAddress) error {
	withdrawInfo, ok := k.GetWithdrawInfo(ctx, addr)
	if !ok {
		return sdk.ErrInvalidAddress(fmt.Sprintf("there is no withdrawing for address %s", addr.String()))
	}
	withdrawCoins := withdrawInfo.Deposits.ToCoins()
	err := k.GetSupplyKeeper().SendCoinsFromModuleToAccount(ctx, types.ModuleName, withdrawInfo.Owner, withdrawCoins)
	if err != nil {
		return types.ErrInsufficientDepositCoins(err.Error(), withdrawCoins.String())
	}
	k.deleteWithdrawInfo(ctx, addr)
	return nil
}

// SetGovKeeper sets keeper of gov
func (k *Keeper) SetGovKeeper(gk GovKeeper) {
	k.govKeeper = gk
}

// GetMaxTokenPairID returns the max ID of token pair
func (k Keeper) GetMaxTokenPairID(ctx sdk.Context) (tokenPairMaxID uint64) {
	store := ctx.KVStore(k.tokenPairStoreKey)
	b := store.Get(types.MaxTokenPairIDKey)
	if b != nil {
		k.cdc.MustUnmarshalBinaryBare(b, &tokenPairMaxID)
	}
	return
}

// GetOperator gets the DEXOperator and checks whether the operator with address exist or not
func (k Keeper) GetOperator(ctx sdk.Context, addr sdk.AccAddress) (operator types.DEXOperator, isExist bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOperatorAddressKey(addr))
	if bz == nil {
		return operator, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &operator)
	return operator, true
}

// IterateOperators iterates over the all the operators and performs a callback function
func (k Keeper) IterateOperators(ctx sdk.Context, cb func(operator types.DEXOperator) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.DEXOperatorKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var operator types.DEXOperator
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &operator)

		if cb(operator) {
			break
		}
	}
}

// SetOperator save the operator information
func (k Keeper) SetOperator(ctx sdk.Context, operator types.DEXOperator) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetOperatorAddressKey(operator.Address)
	bytes := k.cdc.MustMarshalBinaryLengthPrefixed(operator)
	store.Set(key, bytes)
}

// SetMaxTokenPairID sets the max ID of token pair
func (k Keeper) SetMaxTokenPairID(ctx sdk.Context, MaxtokenPairID uint64) {
	store := ctx.KVStore(k.tokenPairStoreKey)
	b := k.cdc.MustMarshalBinaryBare(MaxtokenPairID)
	store.Set(types.MaxTokenPairIDKey, b)
}

func (k *Keeper) SetObserverKeeper(sk exported.StreamKeeper) {
	k.observerKeeper = sk
}
