package keeper

import (
	"fmt"
	types2 "github.com/okex/exchain/libs/tendermint/types"
	types3 "github.com/okex/exchain/libs/types"

	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params/subspace"
)

// AccountKeeper encodes/decodes accounts using the go-amino (binary)
// encoding/decoding library.
type AccountKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	mptKey sdk.StoreKey

	// The prototypical Account constructor.
	proto func() exported.Account

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	paramSubspace subspace.Subspace

	observers []ObserverI
}

// NewAccountKeeper returns a new sdk.AccountKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.Accounts.
// nolint
func NewAccountKeeper(
	cdc *codec.Codec, key, keyMpt sdk.StoreKey, paramstore subspace.Subspace, proto func() exported.Account,
) AccountKeeper {

	return AccountKeeper{
		key:           key,
		mptKey:        keyMpt,
		proto:         proto,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
	}
}

// Logger returns a module-specific logger.
func (ak AccountKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetPubKey Returns the PubKey of the account at address
func (ak AccountKeeper) GetPubKey(ctx sdk.Context, addr sdk.AccAddress) (crypto.PubKey, error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
	}
	return acc.GetPubKey(), nil
}

// GetSequence Returns the Sequence of the account at address
func (ak AccountKeeper) GetSequence(ctx sdk.Context, addr sdk.AccAddress) (uint64, error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
	}
	return acc.GetSequence(), nil
}

// GetNextAccountNumber returns and increments the global account number counter.
// If the global account number is not set, it initializes it with value 0.
func (ak AccountKeeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	var accNumber uint64
	var store sdk.KVStore
	if types2.HigherThanMars(ctx.BlockHeight()) {
		store = ctx.KVStore(ak.mptKey)
	} else {
		store = ctx.KVStore(ak.key)
	}
	bz := store.Get(types.GlobalAccountNumberKey)
	if bz == nil {
		// initialize the account numbers
		accNumber = 0
	} else {
		err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &accNumber)
		if err != nil {
			panic(err)
		}
	}

	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(accNumber + 1)
	store.Set(types.GlobalAccountNumberKey, bz)
	if !types2.HigherThanMars(ctx.BlockHeight()) && types3.EnableDoubleWrite {
		ctx.MultiStore().GetKVStore(ak.mptKey).Set(types.GlobalAccountNumberKey, bz)
	}

	return accNumber
}

func (ak AccountKeeper) SetEvmRootHash(ctx sdk.Context, hash []byte) {
	var store sdk.KVStore
	if types2.HigherThanMars(ctx.BlockHeight()) {
		store = ctx.KVStore(ak.mptKey)
	} else {
		store = ctx.KVStore(ak.key)
	}

	store.Set(types.MptRootStoreKeyPrefix, hash)
}

// -----------------------------------------------------------------------------
// Misc.

func (ak AccountKeeper) decodeAccount(bz []byte) (acc exported.Account) {
	val, err := ak.cdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(bz, &acc)
	if err == nil {
		acc = val.(exported.Account)
		return
	}
	err = ak.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}
