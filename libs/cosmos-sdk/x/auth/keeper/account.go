package keeper

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/mpt"
	mpttypes "github.com/okex/exchain/libs/mpt/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/tendermint/go-amino"
)

// NewAccountWithAddress implements sdk.AccountKeeper.
func (ak AccountKeeper) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) exported.Account {
	acc := ak.proto()
	err := acc.SetAddress(addr)
	if err != nil {
		panic(err)
	}
	return ak.NewAccount(ctx, acc)
}

// NewAccount sets the next account number to a given account interface
func (ak AccountKeeper) NewAccount(ctx sdk.Context, acc exported.Account) exported.Account {
	if err := acc.SetAccountNumber(ak.GetNextAccountNumber(ctx)); err != nil {
		panic(err)
	}
	return acc
}

// GetAccount implements sdk.AccountKeeper.
func (ak AccountKeeper) GetAccount(ctx sdk.Context, addr sdk.AccAddress) exported.Account {
	if data, gas, ok := ctx.Cache().GetAccount(ethcmn.BytesToAddress(addr)); ok {
		ctx.GasMeter().ConsumeGas(gas, "x/auth/keeper/account.go/GetAccount")
		if data == nil {
			return nil
		}

		return data.Copy().(exported.Account)
	}

	var store sdk.KVStore
	if tmtypes.HigherThanMars(ctx.BlockHeight()) {
		store = ctx.KVStore(ak.mptKey)
	} else {
		store = ctx.KVStore(ak.key)
	}

	bz := store.Get(types.AddressStoreKey(addr))
	if bz == nil {
		ctx.Cache().UpdateAccount(addr, nil, len(bz), false)
		return nil
	}
	acc := ak.decodeAccount(bz)
	ctx.Cache().UpdateAccount(addr, acc, len(bz), false)
	return acc
}

// GetAllAccounts returns all accounts in the accountKeeper.
func (ak AccountKeeper) GetAllAccounts(ctx sdk.Context) (accounts []exported.Account) {
	ak.IterateAccounts(ctx,
		func(acc exported.Account) (stop bool) {
			accounts = append(accounts, acc)
			return false
		})
	return accounts
}

// SetAccount implements sdk.AccountKeeper.
func (ak AccountKeeper) SetAccount(ctx sdk.Context, acc exported.Account) {
	addr := acc.GetAddress()

	var store sdk.KVStore
	if tmtypes.HigherThanMars(ctx.BlockHeight()) {
		store = ctx.KVStore(ak.mptKey)
	} else {
		store = ctx.KVStore(ak.key)
	}

	bz, err := ak.cdc.MarshalBinaryBareWithRegisteredMarshaller(acc)
	if err != nil {
		bz, err = ak.cdc.MarshalBinaryBare(acc)
	}
	if err != nil {
		panic(err)
	}

	storeAccKey := types.AddressStoreKey(addr)
	store.Set(storeAccKey, bz)
	if !tmtypes.HigherThanMars(ctx.BlockHeight()) && mpttypes.EnableDoubleWrite {
		ctx.MultiStore().GetKVStore(ak.mptKey).Set(storeAccKey, bz)
	}
	ctx.Cache().UpdateAccount(addr, acc, len(bz), true)

	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		mpt.GAccToPrefetchChannel <- [][]byte{storeAccKey}

		if ak.observers != nil {
			for _, observer := range ak.observers {
				if observer != nil {
					observer.OnAccountUpdated(acc)
				}
			}
		}
	}
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (ak AccountKeeper) RemoveAccount(ctx sdk.Context, acc exported.Account) {
	addr := acc.GetAddress()
	var store sdk.KVStore
	if tmtypes.HigherThanMars(ctx.BlockHeight()) {
		store = ctx.KVStore(ak.mptKey)
	} else {
		store = ctx.KVStore(ak.key)
	}

	storeAccKey := types.AddressStoreKey(addr)
	store.Delete(storeAccKey)
	if !tmtypes.HigherThanMars(ctx.BlockHeight()) && mpttypes.EnableDoubleWrite {
		ctx.MultiStore().GetKVStore(ak.mptKey).Delete(storeAccKey)
	}

	if !ctx.IsCheckTx() && !ctx.IsReCheckTx() {
		mpt.GAccToPrefetchChannel <- [][]byte{storeAccKey}
	}

	ctx.Cache().UpdateAccount(addr, nil, 0, true)
}

// IterateAccounts iterates over all the stored accounts and performs a callback function
func (ak AccountKeeper) IterateAccounts(ctx sdk.Context, cb func(account exported.Account) (stop bool)) {
	var store sdk.KVStore
	if tmtypes.HigherThanMars(ctx.BlockHeight()) {
		store = ctx.KVStore(ak.mptKey)
	} else {
		store = ctx.KVStore(ak.key)
	}
	iterator := sdk.KVStorePrefixIterator(store, types.AddressStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		account := ak.decodeAccount(iterator.Value())

		if cb(account) {
			break
		}
	}
}

// IterateAccounts iterates over all the stored accounts and performs a callback function
func (ak AccountKeeper) MigrateAccounts(ctx sdk.Context, cb func(account exported.Account, key, value []byte) (stop bool)) {
	var store sdk.KVStore
	if tmtypes.HigherThanMars(ctx.BlockHeight()) {
		store = ctx.KVStore(ak.mptKey)
	} else {
		store = ctx.KVStore(ak.key)
	}
	iterator := sdk.KVStorePrefixIterator(store, types.AddressStoreKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		account := ak.decodeAccount(iterator.Value())

		if cb(account, iterator.Key(), iterator.Value()) {
			break
		}
	}
}

func (ak AccountKeeper) GetEncodedAccountSize(acc exported.Account) int {
	if sizer, ok := acc.(amino.Sizer); ok {
		// typeprefix + amino bytes
		return 4 + sizer.AminoSize(ak.cdc)
	} else {
		return len(ak.cdc.MustMarshalBinaryBare(acc))
	}
}
