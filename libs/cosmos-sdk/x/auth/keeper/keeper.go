package keeper

import (
	"encoding/binary"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/prque"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/pkg/errors"

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

	// The prototypical Account constructor.
	proto func() exported.Account

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	paramSubspace subspace.Subspace

	observers []ObserverI

	trie ethstate.Trie
	db ethstate.Database

	accLRU   *lru.Cache
	deliverTxStore *types.CacheStore
	checkTxStore *types.CacheStore
	triegc *prque.Prque
}

// NewAccountKeeper returns a new sdk.AccountKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.Accounts.
// nolint
func NewAccountKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramstore subspace.Subspace, proto func() exported.Account,
) AccountKeeper {
	accLRU, e := lru.New(500000)
	if e != nil {
		panic(errors.New("Failed to init LRU Cause " + e.Error()))
	}

	ak := AccountKeeper{
		key:           key,
		proto:         proto,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
		db: types.InstanceOfEvmStore(),
		accLRU: accLRU,
		deliverTxStore: types.NewCacheStore(),
		checkTxStore: types.NewCacheStore(),
		triegc: prque.New(nil),
	}

	ak.OpenTrie()

	return ak
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
	store := ctx.KVStore(ak.key)
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

	return accNumber
}

// -----------------------------------------------------------------------------
// Misc.

func (ak AccountKeeper) decodeAccount(bz []byte) (acc exported.Account) {
	err := ak.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}

func (ak *AccountKeeper) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	ak.OpenTrie()
}

func (ak *AccountKeeper) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	ak.Commit(ctx)

	return []abci.ValidatorUpdate{}
}

var (
	KeyPrefixLatestHeight = []byte{0x01}
	KeyPrefixRootMptHash  = []byte{0x02}
)

// GetLatestBlockHeight get latest mpt storage height
func (ak *AccountKeeper) GetLatestBlockHeight() uint64 {
	rst, err := ak.db.TrieDB().DiskDB().Get(KeyPrefixLatestHeight)
	if err != nil || len(rst) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(rst)
}

// SetLatestBlockHeight sets the latest storage height
func (ak *AccountKeeper) SetLatestBlockHeight(height uint64) {
	hhash := sdk.Uint64ToBigEndian(height)
	ak.db.TrieDB().DiskDB().Put(KeyPrefixLatestHeight, hhash)
}

// GetRootMptHash gets root mpt hash from block height
func (ak *AccountKeeper) GetRootMptHash(height uint64) ethcmn.Hash {
	hhash := sdk.Uint64ToBigEndian(height)
	rst, err := ak.db.TrieDB().DiskDB().Get(append(KeyPrefixRootMptHash, hhash...))
	if err != nil || len(rst) == 0 {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(rst)
}

// SetRootMptHash sets the mapping from block height to root mpt hash
func (ak *AccountKeeper) SetRootMptHash(height uint64, hash ethcmn.Hash) {
	hhash := sdk.Uint64ToBigEndian(height)
	ak.db.TrieDB().DiskDB().Put(append(KeyPrefixRootMptHash, hhash...), hash.Bytes())
}

func (ak *AccountKeeper) OpenTrie() {
	latestHeight := ak.GetLatestBlockHeight()
	lastRootHash := ak.GetRootMptHash(latestHeight)

	tr, err := ak.db.OpenTrie(lastRootHash)
	if err != nil {
		panic("Fail to open root mpt: " + err.Error())
	}
	ak.trie = tr
}

func (ak *AccountKeeper) Commit(ctx sdk.Context) {
	// The onleaf func is called _serially_, so we can reuse the same account
	// for unmarshalling every time.
	var data []byte
	root, _ := ak.trie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
		if err := rlp.DecodeBytes(leaf, &data); err != nil {
			return nil
		}
		accStorageRoot := ak.decodeAccount(data).GetStorageRoot()

		if accStorageRoot != types2.EmptyRootHash && accStorageRoot != (ethcmn.Hash{}) {
			ak.db.TrieDB().Reference(accStorageRoot, parent)
		}

		return nil
	})

	latestHeight := uint64(ctx.BlockHeight())

	ak.SetRootMptHash(latestHeight, root)
	ak.SetLatestBlockHeight(latestHeight)
	ak.CleanCacheStore()

	ak.PushData2Database(ctx, root)
}

func (ak *AccountKeeper) OnStop() error {
	for !ak.triegc.Empty() {
		ak.db.TrieDB().Dereference(ak.triegc.PopItem().(ethcmn.Hash))
	}

	return nil
}