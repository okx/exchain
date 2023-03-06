package mpt

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func genKVStore() ethdb.KeyValueStore {
	kvstore, _ := types.NewMptMemDB("", "")
	return kvstore
}

type kv struct {
	key   []byte
	value []byte
}

func genTestKVs(num int) []*kv {
	var kvs []*kv
	for i := 0; i < num; i++ {
		kvs = append(kvs, &kv{randBytes(32), randBytes(32)})
	}
	return kvs
}

func TestStatKeyValueStoreInterface(t *testing.T) {
	kvstore := genKVStore()
	nkvstore := NewStatKeyValueStore(kvstore, nil)
	tkv := genTestKVs(3)
	for _, kv := range tkv {
		err := nkvstore.Put(kv.key, kv.value)
		assert.NoError(t, err)
		v, err := nkvstore.Get(kv.key)
		assert.NoError(t, err)
		assert.Equal(t, kv.value, v)
		b, err := nkvstore.Has(kv.key)
		assert.NoError(t, err)
		assert.True(t, b)
		err = nkvstore.Delete(kv.key)
		assert.NoError(t, err)
	}
	batch := nkvstore.NewBatch()
	for _, kv := range tkv {
		batch.Put(kv.key, kv.value)
	}
	batch.Write()
	v, err := nkvstore.Get(tkv[0].key)
	assert.NoError(t, err)
	assert.Equal(t, tkv[0].value, v)

	it := nkvstore.NewIterator(nil, nil)
	for it.Next() {
		fmt.Println("key is", it.Key(), "value is", it.Value())
	}
	it.Release()

}

func TestStatKeyValueStore(t *testing.T) {
	kvstore := genKVStore()
	stat := NewRuntimeState()
	nkvstore := NewStatKeyValueStore(kvstore, stat)
	tkv := genTestKVs(10)
	dbCount := 0
	statTime := 0
	for _, kv := range tkv {
		nkvstore.Put(kv.key, kv.value)
		nkvstore.Get(kv.key)
		dbCount++
		assert.Equal(t, stat.getDBReadCount(), dbCount)
		assert.Equal(t, stat.getDBWriteCount(), dbCount)
		assert.NotEqual(t, stat.getDBReadTime(), statTime)
	}
	stat.resetCount()

	batch := nkvstore.NewBatch()
	for _, kv := range tkv {
		batch.Put(kv.key, kv.value)
	}
	batch.Write()
	assert.Equal(t, stat.getDBWriteCount(), len(tkv))
}
