package keeper

import (
	"encoding/binary"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/prque"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	types2 "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/wrap"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params/subspace"
	"github.com/okex/exchain/libs/tendermint/crypto"
	"github.com/okex/exchain/libs/tendermint/libs/log"
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
	triegc *prque.Prque

	accCommitStore *sdk.AccCommitStore
}

// NewAccountKeeper returns a new sdk.AccountKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.Accounts.
// nolint
func NewAccountKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramstore subspace.Subspace, proto func() exported.Account,
) AccountKeeper {
	ak := AccountKeeper{
		key:           key,
		proto:         proto,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),

		db: types.InstanceOfEvmStore(),
		triegc: prque.New(nil),

		accCommitStore: sdk.NewAccCommitStore(),
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

var (
	KeyPrefixLatestHeight = []byte{0x01}
	KeyPrefixRootMptHash  = []byte{0x02}
	KeyPrefixLatestStoredHeight  = []byte{0x03}
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

// GetLatestStoredBlockHeight get latest stored mpt storage height
func (ak *AccountKeeper) GetLatestStoredBlockHeight() uint64 {
	rst, err := ak.db.TrieDB().DiskDB().Get(KeyPrefixLatestStoredHeight)
	if err != nil || len(rst) == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(rst)
}

// SetLatestStoredBlockHeight sets the latest stored storage height
func (ak *AccountKeeper) SetLatestStoredBlockHeight(height uint64) {
	hhash := sdk.Uint64ToBigEndian(height)
	ak.db.TrieDB().DiskDB().Put(KeyPrefixLatestStoredHeight, hhash)
}

func (ak *AccountKeeper) OpenTrie() {
	//types3.GetStartBlockHeight() // start height of oec
	latestHeight := ak.GetLatestBlockHeight()
	latestRootHash := ak.GetRootMptHash(latestHeight)

	tr, err := ak.db.OpenTrie(latestRootHash)
	if err != nil {
		panic("Fail to open root mpt: " + err.Error())
	}
	ak.trie = tr

	ak.accCommitStore.SetMptTrie(tr)
}

func (ak *AccountKeeper) Commit(ctx sdk.Context) {
	ak.accCommitStore.Write() // cs.write()

	// The onleaf func is called _serially_, so we can reuse the same account
	// for unmarshalling every time.
	var wrapAcc wrap.WrapAccount
	root, _ := ak.trie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
		if err := rlp.DecodeBytes(leaf, &wrapAcc); err != nil {
			return nil
		}
		accStorageRoot := wrapAcc.RealAcc.GetStorageRoot()
		if accStorageRoot != types2.EmptyRootHash && accStorageRoot != (ethcmn.Hash{}) {
			ak.db.TrieDB().Reference(accStorageRoot, parent)
		}
		return nil
	})

	latestHeight := uint64(ctx.BlockHeight())
	ak.SetRootMptHash(latestHeight, root)
	ak.SetLatestBlockHeight(latestHeight)

	ak.PushData2Database(ctx, root)
}

func (ak *AccountKeeper) OnStop() error {
	for !ak.triegc.Empty() {
		ak.db.TrieDB().Dereference(ak.triegc.PopItem().(ethcmn.Hash))
	}

	return nil
}