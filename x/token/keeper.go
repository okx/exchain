package token

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/okex/okchain/x/params"
	"github.com/okex/okchain/x/token/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	bankKeeper       bank.Keeper
	supplyKeeper     SupplyKeeper
	feeCollectorName string // name of the FeeCollector ModuleAccount

	// The reference to the Paramstore to get and set gov specific params
	paramSpace    params.Subspace
	tokenStoreKey sdk.StoreKey // Unexposed key to access name store from sdk.Context
	lockStoreKey  sdk.StoreKey
	//TokenPairNewSignalChan chan types.TokenPair

	cdc *codec.Codec // The wire codec for binary encoding/decoding.

	enableBackend bool // whether open backend plugin

	// cache data in memory to avoid marshal/unmarshal too frequently
	// reset cache data in BeginBlock
	cache *Cache
}

// NewKeeper creates a new token keeper
func NewKeeper(bankKeeper bank.Keeper, paramSpace params.Subspace,
	feeCollectorName string, supplyKeeper SupplyKeeper, tokenStoreKey, lockStoreKey sdk.StoreKey, cdc *codec.Codec, enableBackend bool) Keeper {

	k := Keeper{
		bankKeeper:       bankKeeper,
		paramSpace:       paramSpace.WithKeyTable(types.ParamKeyTable()),
		feeCollectorName: feeCollectorName,
		supplyKeeper:     supplyKeeper,
		tokenStoreKey:    tokenStoreKey,
		lockStoreKey:     lockStoreKey,
		cdc:              cdc,
		enableBackend:    enableBackend,
		cache:            NewCache(),
	}
	return k
}

// nolint
func (k Keeper) ResetCache(ctx sdk.Context) {
	k.cache.reset()
}

// nolint
func (k Keeper) GetTokenInfo(ctx sdk.Context, symbol string) types.Token {
	var token types.Token
	store := ctx.KVStore(k.tokenStoreKey)
	bz := store.Get(types.GetTokenAddress(symbol))
	if bz == nil {
		return token
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &token)

	supply := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(token.Symbol)
	token.TotalSupply = supply

	return token
}

// TokenExist checks whether the token with symbol exist or not
func (k Keeper) TokenExist(ctx sdk.Context, symbol string) bool {
	store := ctx.KVStore(k.tokenStoreKey)
	bz := store.Get(types.GetTokenAddress(symbol))
	return bz != nil
}

// nolint
func (k Keeper) GetTokensInfo(ctx sdk.Context) (tokens []types.Token) {
	store := ctx.KVStore(k.tokenStoreKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenKey)
	defer iter.Close()
	for iter.Valid() {
		var token types.Token
		tokenBytes := iter.Value()
		k.cdc.MustUnmarshalBinaryBare(tokenBytes, &token)

		supply := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(token.Symbol)
		token.TotalSupply = supply

		tokens = append(tokens, token)
		iter.Next()
	}
	return tokens
}

// GetUserTokensInfo gets tokens info by owner address
func (k Keeper) GetUserTokensInfo(ctx sdk.Context, owner sdk.AccAddress) (tokens []types.Token) {
	userTokenPrefix := types.GetUserTokenPrefix(owner)
	userTokenPrefixLen := len(userTokenPrefix)
	store := ctx.KVStore(k.tokenStoreKey)
	iter := sdk.KVStorePrefixIterator(store, userTokenPrefix)
	defer iter.Close()
	for iter.Valid() {
		userTokenKey := iter.Key()
		symbol := string(userTokenKey[userTokenPrefixLen:])
		tokens = append(tokens, k.GetTokenInfo(ctx, symbol))

		iter.Next()
	}

	return tokens
}

// GetCurrenciesInfo returns all of the currencies info
func (k Keeper) GetCurrenciesInfo(ctx sdk.Context) (currencies []types.Currency) {
	store := ctx.KVStore(k.tokenStoreKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenKey)
	defer iter.Close()
	//iter := store.Iterator(nil, nil)
	for iter.Valid() {
		var token types.Token
		tokenBytes := iter.Value()
		k.cdc.MustUnmarshalBinaryBare(tokenBytes, &token)

		supply := k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(token.Symbol)
		token.TotalSupply = supply

		currencies = append(currencies,
			types.Currency{
				Description: token.Description,
				Symbol:      token.Symbol,
				TotalSupply: token.OriginalTotalSupply,
			})
		iter.Next()
	}
	return currencies
}

// DeleteUserToken deletes token by user address and symbol
func (k Keeper) DeleteUserToken(ctx sdk.Context, owner sdk.AccAddress, symbol string) {
	store := ctx.KVStore(k.tokenStoreKey)
	store.Delete(types.GetUserTokenKey(owner, symbol))
}

// nolint
func (k Keeper) NewToken(ctx sdk.Context, token types.Token) {
	// save token info
	store := ctx.KVStore(k.tokenStoreKey)
	store.Set(types.GetTokenAddress(token.Symbol), k.cdc.MustMarshalBinaryBare(token))
	store.Set(types.GetUserTokenKey(token.Owner, token.Symbol), []byte{})

	// update token number
	tokenNumber := k.getTokenNum(ctx)
	b := k.cdc.MustMarshalBinaryBare(tokenNumber + 1)
	store.Set(types.TokenNumberKey, b)
}

// send tokens(one or more coins) from one account to another
func (k Keeper) SendCoinsFromAccountToAccount(ctx sdk.Context, from, to sdk.AccAddress, amt sdk.DecCoins) error {
	return k.bankKeeper.SendCoins(ctx, from, to, amt)
}

// nolint
func (k Keeper) LockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins, lockCoinsType int) error {
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, addr, types.ModuleName, coins); err != nil {
		return err
	}
	// update lock coins
	return k.updateLockedCoins(ctx, addr, coins, true, lockCoinsType)
}

// nolint
func (k Keeper) updateLockedCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins, doAdd bool, lockCoinsType int) error {
	var key []byte
	switch lockCoinsType {
	case types.LockCoinsTypeQuantity:
		key = types.GetLockAddress(addr.Bytes())
	case types.LockCoinsTypeFee:
		key = types.GetLockFeeAddress(addr.Bytes())
	default:
		return fmt.Errorf("unrecognized lock coins type: %d", lockCoinsType)
	}

	var newCoins sdk.DecCoins
	var oldCoins sdk.DecCoins

	store := ctx.KVStore(k.lockStoreKey)
	coinsBytes := store.Get(key)

	if doAdd {
		// lock coins
		if coinsBytes == nil {
			newCoins = coins
		} else {
			k.cdc.MustUnmarshalBinaryBare(coinsBytes, &oldCoins)
			newCoins = oldCoins.Add(coins)
		}
	} else {
		// unlock coins
		if coinsBytes == nil {
			return fmt.Errorf("failed to unlock <%s>. Address <%s>, coins locked <0>", coins, addr)
		}
		k.cdc.MustUnmarshalBinaryBare(coinsBytes, &oldCoins)
		var isNegative bool
		newCoins, isNegative = oldCoins.SafeSub(coins)
		if isNegative {
			return fmt.Errorf("failed to unlock <%s>. Address <%s>, coins available <%s>", coins, addr, oldCoins)
		}
	}

	sort.Sort(newCoins)
	if len(newCoins) > 0 {
		store.Set(key, k.cdc.MustMarshalBinaryBare(newCoins))
	} else {
		store.Delete(key)
	}

	return nil
}

// nolint
func (k Keeper) UnlockCoins(ctx sdk.Context, addr sdk.AccAddress, coins sdk.DecCoins, lockCoinsType int) error {
	// update lock coins
	if err := k.updateLockedCoins(ctx, addr, coins, false, lockCoinsType); err != nil {
		return err
	}

	// update account
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

// GetLockCoins gets locked coins by address
func (k Keeper) GetLockedCoins(ctx sdk.Context, addr sdk.AccAddress) (coins sdk.DecCoins) {
	store := ctx.KVStore(k.lockStoreKey)
	coinsBytes := store.Get(types.GetLockAddress(addr.Bytes()))
	if coinsBytes == nil {
		return coins
	}
	k.cdc.MustUnmarshalBinaryBare(coinsBytes, &coins)
	return coins
}

// GetAllLockCoins iterates KVStore and gets all of the locked coins
func (k Keeper) GetAllLockedCoins(ctx sdk.Context) (locks []types.AccCoins) {
	store := ctx.KVStore(k.lockStoreKey)
	iter := sdk.KVStorePrefixIterator(store, types.LockKey)
	defer iter.Close()
	for iter.Valid() {
		var accCoins types.AccCoins
		accCoins.Acc = iter.Key()[len(types.LockKey):]
		coinsBytes := iter.Value()
		var coins sdk.DecCoins
		k.cdc.MustUnmarshalBinaryBare(coinsBytes, &coins)
		accCoins.Coins = coins
		locks = append(locks, accCoins)
		iter.Next()
	}
	return locks
}

// IterateAllDeposits iterates over the all the stored lock fee and performs a callback function
func (k Keeper) IterateLockedFees(ctx sdk.Context, cb func(acc sdk.AccAddress, coins sdk.DecCoins) (stop bool)) {
	store := ctx.KVStore(k.lockStoreKey)
	iter := sdk.KVStorePrefixIterator(store, types.LockedFeeKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		acc := iter.Key()[len(types.LockKey):]

		var coins sdk.DecCoins
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &coins)

		if cb(acc, coins) {
			break
		}
	}
}

// BalanceAccount is ONLY expected by the order module to settle an order where outputCoins
// is used to exchange inputCoins
func (k Keeper) BalanceAccount(ctx sdk.Context, addr sdk.AccAddress, outputCoins sdk.DecCoins,
	inputCoins sdk.DecCoins) (err error) {

	if !outputCoins.IsZero() {
		if err = k.updateLockedCoins(ctx, addr, outputCoins, false, types.LockCoinsTypeQuantity); err != nil {
			return err
		}
	}

	if !inputCoins.IsZero() {
		return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, inputCoins)
	}

	return nil
}

// nolint
func (k Keeper) GetCoins(ctx sdk.Context, addr sdk.AccAddress) sdk.DecCoins {
	return k.bankKeeper.GetCoins(ctx, addr)
}

// GetParams gets inflation params from the global param store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams set inflation params from the global param store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetCoinsInfo gets all of the coin info by addr
func (k Keeper) GetCoinsInfo(ctx sdk.Context, addr sdk.AccAddress) (coinsInfo types.CoinsInfo) {
	availableCoins := k.GetCoins(ctx, addr)
	lockedCoins := k.GetLockedCoins(ctx, addr)

	// merge coins
	coinsInfo = types.MergeCoinInfo(availableCoins, lockedCoins)
	return coinsInfo
}

// GetFeeDetailList gets fee detail list from cache
func (k Keeper) GetFeeDetailList() []*FeeDetail {
	return k.cache.getFeeDetailList()
}

// nolint
func (k Keeper) AddFeeDetail(ctx sdk.Context, from string, fee sdk.DecCoins, feeType string) {
	if k.enableBackend {
		feeDetail := &FeeDetail{
			Address:   from,
			Fee:       fee.String(),
			FeeType:   feeType,
			Timestamp: ctx.BlockHeader().Time.Unix(),
		}
		k.cache.addFeeDetail(feeDetail)
	}
}

func (k Keeper) getNumKeys(ctx sdk.Context) (tokenStoreKeyNum, lockStoreKeyNum int64) {
	{
		store := ctx.KVStore(k.tokenStoreKey)
		iter := store.Iterator(nil, nil)
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			tokenStoreKeyNum++
		}
	}
	{
		store := ctx.KVStore(k.lockStoreKey)
		iter := store.Iterator(nil, nil)
		defer iter.Close()
		for ; iter.Valid(); iter.Next() {
			lockStoreKeyNum++
		}
	}

	return
}

func (k Keeper) getTokenNum(ctx sdk.Context) (tokenNumber uint64) {
	store := ctx.KVStore(k.tokenStoreKey)
	b := store.Get(types.TokenNumberKey)
	if b != nil {
		k.cdc.MustUnmarshalBinaryBare(b, &tokenNumber)
	}
	return
}

// addTokenSuffix add token suffix
func addTokenSuffix(ctx sdk.Context, keeper Keeper, originalSymbol string) (name string, valid bool) {
	hash := fmt.Sprintf("%x", tmhash.Sum(ctx.TxBytes()))
	var i int
	for i = len(hash)/3 - 1; i >= 0; i-- {
		name = originalSymbol + "-" + strings.ToLower(hash[3*i:3*i+3])
		// check token name valid
		if sdk.ValidateDenom(name) != nil {
			return "", false
		}
		if !keeper.TokenExist(ctx, name) {
			break
		}
	}
	if i == -1 {
		return "", false
	}
	return name, true
}

func (k Keeper) SendCoinsFromAccountToModule(ctx sdk.Context, from sdk.AccAddress, amount sdk.Coins) sdk.Error {
	return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, from, types.ModuleName, amount)
}

func (k Keeper) SendCoinsFromModuleToAccount(ctx sdk.Context, to sdk.AccAddress, amount sdk.Coins) sdk.Error {
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, to, amount)
}

func (k Keeper) SendCoinsFromModuleToModule(ctx sdk.Context, recipientModule string, coins sdk.Coins) sdk.Error {
	return k.supplyKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, recipientModule, coins)
}
