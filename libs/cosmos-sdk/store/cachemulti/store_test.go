package cachemulti

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/store/cachekv"
	"github.com/okex/exchain/libs/cosmos-sdk/store/dbadapter"
	dbm "github.com/okex/exchain/libs/tm-db"
	"strconv"
	"testing"

	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/stretchr/testify/require"
)

var keys []*types.KVStoreKey

func TestStoreGetKVStore(t *testing.T) {
	require := require.New(t)

	s := Store{stores: map[types.StoreKey]types.CacheWrap{}}
	key := types.NewKVStoreKey("abc")
	errMsg := fmt.Sprintf("kv store with key %v has not been registered in stores", key)

	require.PanicsWithValue(errMsg,
		func() { s.GetStore(key) })

	require.PanicsWithValue(errMsg,
		func() { s.GetKVStore(key) })
}

func TestStore_WriteGetMultiSnapShotWSet(t *testing.T) {
	store := setupCacheMulti()
	for i := 0; i < 10; i++ {
		if i%3 == 0 {
			//update key
			store.GetKVStore(keys[i]).Set([]byte("1"), []byte("2"))
		} else if i%3 == 1 {
			//insert key
			store.GetKVStore(keys[i]).Set([]byte("2"), []byte("2"))
		} else if i%3 == 2 {
			//delete key
			store.GetKVStore(keys[i]).Delete([]byte("1"))
		}
	}

	snapshot := store.WriteGetMultiSnapshotWSet()
	for i := 0; i < 10; i++ {
		if i%3 == 0 {
			//update key

			iter := store.GetKVStore(keys[i]).Iterator(nil, nil)
			require.True(t, iter.Valid())
			iter.Next()
			require.False(t, iter.Valid())
			value := store.GetKVStore(keys[i]).Get([]byte("1"))
			require.Equal(t, []byte("2"), value)
		} else if i%3 == 1 {
			//insert key
			value := store.GetKVStore(keys[i]).Get([]byte("2"))
			require.Equal(t, []byte("2"), value)
			value = store.GetKVStore(keys[i]).Get([]byte("1"))
			require.Equal(t, []byte("1"), value)

			iter := store.GetKVStore(keys[i]).Iterator(nil, nil)
			require.True(t, iter.Valid())
			iter.Next()
			require.True(t, iter.Valid())
			iter.Next()
			require.False(t, iter.Valid())

		} else if i%3 == 2 {
			//delete key
			value := store.GetKVStore(keys[i]).Get([]byte("1"))
			require.Nil(t, value)
			iter := store.GetKVStore(keys[i]).Iterator(nil, nil)
			require.False(t, iter.Valid())
		}
	}

	require.Equal(t, 10, len(snapshot.Stores))
	for i := 0; i < 10; i++ {
		if i%3 == 0 {
			//update key

			require.Equal(t, 1, len(snapshot.Stores[keys[i]].Write))
			require.Equal(t, snapshot.Stores[keys[i]].Write["1"].PrevValue, []byte("1"))
		} else if i%3 == 1 {
			//insert key

			require.Equal(t, 1, len(snapshot.Stores[keys[i]].Write))
			require.Nil(t, snapshot.Stores[keys[i]].Write["2"].PrevValue)
		} else if i%3 == 2 {
			//delete key

			require.Equal(t, 1, len(snapshot.Stores[keys[i]].Write))
			require.Equal(t, snapshot.Stores[keys[i]].Write["1"].PrevValue, []byte("1"))
		}
	}
}

func TestStore_RevertDBWithMultiSnapShotRWSet(t *testing.T) {
	store := setupCacheMulti()
	for i := 0; i < 10; i++ {
		if i%3 == 0 {
			//update key
			store.GetKVStore(keys[i]).Set([]byte("1"), []byte("2"))
		} else if i%3 == 1 {
			//insert key
			store.GetKVStore(keys[i]).Set([]byte("2"), []byte("2"))
		} else if i%3 == 2 {
			//delete key
			store.GetKVStore(keys[i]).Delete([]byte("1"))
		}
	}

	snapshot := store.WriteGetMultiSnapshotWSet()
	for i := 0; i < 10; i++ {
		if i%3 == 0 {
			//update key

			iter := store.GetKVStore(keys[i]).Iterator(nil, nil)
			require.True(t, iter.Valid())
			iter.Next()
			require.False(t, iter.Valid())
			value := store.GetKVStore(keys[i]).Get([]byte("1"))
			require.Equal(t, []byte("2"), value)
		} else if i%3 == 1 {
			//insert key
			value := store.GetKVStore(keys[i]).Get([]byte("2"))
			require.Equal(t, []byte("2"), value)
			value = store.GetKVStore(keys[i]).Get([]byte("1"))
			require.Equal(t, []byte("1"), value)

			iter := store.GetKVStore(keys[i]).Iterator(nil, nil)
			require.True(t, iter.Valid())
			iter.Next()
			require.True(t, iter.Valid())
			iter.Next()
			require.False(t, iter.Valid())

		} else if i%3 == 2 {
			//delete key
			value := store.GetKVStore(keys[i]).Get([]byte("1"))
			require.Nil(t, value)
			iter := store.GetKVStore(keys[i]).Iterator(nil, nil)
			require.False(t, iter.Valid())
		}
	}

	require.Equal(t, 10, len(snapshot.Stores))
	for i := 0; i < 10; i++ {
		if i%3 == 0 {
			//update key

			require.Equal(t, 1, len(snapshot.Stores[keys[i]].Write))
			require.Equal(t, snapshot.Stores[keys[i]].Write["1"].PrevValue, []byte("1"))
		} else if i%3 == 1 {
			//insert key

			require.Equal(t, 1, len(snapshot.Stores[keys[i]].Write))
			require.Nil(t, snapshot.Stores[keys[i]].Write["2"].PrevValue)
		} else if i%3 == 2 {
			//delete key

			require.Equal(t, 1, len(snapshot.Stores[keys[i]].Write))
			require.Equal(t, snapshot.Stores[keys[i]].Write["1"].PrevValue, []byte("1"))
		}
	}
	store.RevertDBWithMultiSnapshotRWSet(snapshot)

	for i := 0; i < 10; i++ {
		iter := store.GetKVStore(keys[i]).Iterator(nil, nil)
		require.True(t, iter.Valid())
		iter.Next()
		require.False(t, iter.Valid())

		value := store.GetKVStore(keys[i]).Get([]byte("1"))
		require.Equal(t, []byte("1"), value)
	}
}

func setupCacheMulti() Store {
	keys = make([]*types.KVStoreKey, 0)
	s := Store{stores: map[types.StoreKey]types.CacheWrap{}}
	for i := 0; i < 10; i++ {
		key := types.NewKVStoreKey(strconv.Itoa(i))
		keys = append(keys, key)
		mem := dbadapter.Store{DB: dbm.NewMemDB()}
		s.stores[key] = cachekv.NewStore(mem)
		s.GetKVStore(key).Set([]byte("1"), []byte("1"))
	}
	mem := dbadapter.Store{DB: dbm.NewMemDB()}
	s.db = cachekv.NewStore(mem)
	s.Write()
	return s
}
