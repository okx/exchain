package types

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/store/iavl"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
	"path/filepath"
)

const storageDB = "storage"

var evmStorageDB dbm.DB

// InitEvmStorageDB inits storage db for contract account
func InitEvmStorageDB() {
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, "data")
	db, err := sdk.NewLevelDB(storageDB, dbPath)
	if err != nil {
		panic(err)
	}
	evmStorageDB = db
}

type StorageStore struct {
	store    sdk.CommitKVStore
	commitID sdk.CommitID
}

func (ss *StorageStore) GetCommitID() sdk.CommitID {
	return ss.commitID
}

func (ss *StorageStore) GetStore() sdk.CommitKVStore {
	return ss.store
}

type CacheStorageStores struct {
	stores map[ethcmn.Address]*StorageStore
}

func NewCacheStorageStore() *CacheStorageStores {
	return &CacheStorageStores{
		stores: make(map[ethcmn.Address]*StorageStore, 0),
	}
}

func (c *CacheStorageStores) Update(addr ethcmn.Address, store sdk.CommitKVStore) {
	c.stores[addr] = &StorageStore{
		store: store,
	}
}

func (c *CacheStorageStores) Reset() {
	c.stores = make(map[ethcmn.Address]*StorageStore, 0)
}

func (c *CacheStorageStores) Commit(ctx sdk.Context, storeKey sdk.StoreKey) error {
	for addr, store := range c.stores {
		commitID := store.store.Commit()
		store.commitID = commitID
		if err := SetStorageLatestCommitID(ctx, storeKey, addr, commitID); err != nil {
			return err
		}
	}

	// save the storage stores updated in the current height
	return SetCommitIDsByHeight(c, ctx, storeKey)
}

func getPruningHeights(ctx sdk.Context) []int64 {
	var pruneHeights []int64
	previousHeight := ctx.BlockHeight() - 1
	pruningOpts := storetypes.NewPruningOptionsFromString(viper.GetString("pruning"))
	if int64(pruningOpts.KeepRecent) < previousHeight {
		pruneHeight := previousHeight - int64(pruningOpts.KeepRecent)
		// We consider this height to be pruned iff:
		//
		// - KeepEvery is zero as that means that all heights should be pruned.
		// - KeepEvery % (height - KeepRecent) != 0 as that means the height is not
		// a 'snapshot' height.
		if pruningOpts.KeepEvery == 0 || pruneHeight%int64(pruningOpts.KeepEvery) != 0 {
			pruneHeights = append(pruneHeights, pruneHeight)
		}
	}
	return pruneHeights
}

func Pruning(ctx sdk.Context, storeKey sdk.StoreKey) error {
	pruHeights := getPruningHeights(ctx)
	for _, height := range pruHeights {
		commitIDs, err := GetCommitIDsByHeight(
			ctx, storeKey, height,
		)
		if err != nil {
			return err
		}
		for addr, commitID := range commitIDs {
			latestCommitID := GetStorageLatestCommitID(ctx, storeKey, addr)
			store, err := LoadAccountStorageStore(addr, latestCommitID)
			if err != nil {
				return err
			}
			if err = store.(*iavl.Store).DeleteVersions(commitID.Version); err != nil {
				return err
			}
		}
		DeleteCommitIDsByHeight(ctx, storeKey, height)
	}
	return nil
}

func SetCommitIDsByHeight(c *CacheStorageStores, ctx sdk.Context, storeKey sdk.StoreKey) error {
	if len(c.stores) == 0 {
		return nil
	}
	commitIDs := make(map[ethcmn.Address]sdk.CommitID, 0)
	for addr, store := range c.stores {
		commitIDs[addr] = store.commitID
	}
	store := ctx.KVStore(storeKey)
	bz, err := json.Marshal(commitIDs)
	if err != nil {
		return err
	}
	store.Set(HeightStoragesPrefix(ctx.BlockHeight()), bz)
	return nil
}

func GetCommitIDsByHeight(ctx sdk.Context, storeKey sdk.StoreKey, height int64) (map[ethcmn.Address]sdk.CommitID, error) {
	store := ctx.KVStore(storeKey)
	var c map[ethcmn.Address]sdk.CommitID
	bz := store.Get(HeightStoragesPrefix(height))
	if bz == nil {
		return make(map[ethcmn.Address]sdk.CommitID, 0), nil
	}
	err := json.Unmarshal(bz, &c)
	if err != nil {
		return make(map[ethcmn.Address]sdk.CommitID, 0), nil
	}

	return c, nil
}

func DeleteCommitIDsByHeight(ctx sdk.Context, storeKey sdk.StoreKey, height int64) {
	store := ctx.KVStore(storeKey)
	store.Delete(HeightStoragesPrefix(height))
}

func SetStorageLatestCommitID(ctx sdk.Context, storeKey sdk.StoreKey, addr ethcmn.Address, commitID sdk.CommitID) error {
	store := prefix.NewStore(ctx.KVStore(storeKey), AddressStoragePrefix(addr))
	bz, err := json.Marshal(commitID)
	if err != nil {
		return err
	}
	store.Set(addr.Bytes(), bz)
	return nil
}

func GetStorageLatestCommitID(ctx sdk.Context, storeKey sdk.StoreKey, addr ethcmn.Address) sdk.CommitID {
	store := prefix.NewStore(ctx.KVStore(storeKey), AddressStoragePrefix(addr))
	bz := store.Get(addr.Bytes())
	if bz == nil {
		return sdk.CommitID{}
	}
	var commitID sdk.CommitID
	if err := json.Unmarshal(bz, &commitID); err != nil {
		return sdk.CommitID{}
	}
	return commitID
}

func LoadAccountStorageStore(address ethcmn.Address, commitID sdk.CommitID) (sdk.CommitKVStore, error) {
	prefix := "s/k:" + address.String() + "/"
	db := dbm.NewPrefixDB(evmStorageDB, []byte(prefix))
	store, err := iavl.LoadStoreWithInitialVersion(db, commitID, false, uint64(tmtypes.GetStartBlockHeight()))
	if err != nil {
		return nil, err
	}
	return store, nil
}
