package iavl

import (
	"bytes"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
)

func TestDelete(t *testing.T) {
	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 0)
	require.NoError(t, err)

	tree.set([]byte("k1"), []byte("Fred"))
	hash, version, err := tree.SaveVersion()
	require.NoError(t, err)
	_, _, err = tree.SaveVersion()
	require.NoError(t, err)

	require.NoError(t, tree.DeleteVersion(version))

	k1Value, _, _ := tree.GetVersionedWithProof([]byte("k1"), version)
	require.Nil(t, k1Value)

	key := tree.ndb.rootKey(version)
	err = memDB.Set(key, hash)
	require.NoError(t, err)
	tree.versions.Set(version, true)

	k1Value, _, err = tree.GetVersionedWithProof([]byte("k1"), version)
	require.Nil(t, err)
	require.Equal(t, 0, bytes.Compare([]byte("Fred"), k1Value))
}

func TestTraverse(t *testing.T) {
	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 0)
	require.NoError(t, err)

	for i := 0; i < 6; i++ {
		tree.set([]byte(fmt.Sprintf("k%d", i)), []byte(fmt.Sprintf("v%d", i)))
	}

	require.Equal(t, 11, tree.nodeSize(), "Size of tree unexpected")
}

func TestMutableTree_DeleteVersions(t *testing.T) {
	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 0)
	require.NoError(t, err)

	type entry struct {
		key   []byte
		value []byte
	}

	versionEntries := make(map[int64][]entry)

	// create 10 tree versions, each with 1000 random key/value entries
	for i := 0; i < 10; i++ {
		entries := make([]entry, 1000)

		for j := 0; j < 1000; j++ {
			k := randBytes(10)
			v := randBytes(10)

			entries[j] = entry{k, v}
			_ = tree.Set(k, v)
		}

		_, v, err := tree.SaveVersion()
		require.NoError(t, err)

		versionEntries[v] = entries
	}

	// delete even versions
	versionsToDelete := []int64{2, 4, 6, 8}
	require.NoError(t, tree.DeleteVersions(versionsToDelete...))

	// ensure even versions have been deleted
	for _, v := range versionsToDelete {
		require.False(t, tree.versions.Get(v))

		_, err := tree.LazyLoadVersion(v)
		require.Error(t, err)
	}

	// ensure odd number versions exist and we can query for all set entries
	for _, v := range []int64{1, 3, 5, 7, 9, 10} {
		require.True(t, tree.versions.Get(v))

		_, err := tree.LazyLoadVersion(v)
		require.NoError(t, err)

		for _, e := range versionEntries[v] {
			_, val := tree.Get(e.key)
			require.Equal(t, e.value, val)
		}
	}
}

func TestMutableTree_DeleteVersionsRange(t *testing.T) {
	require := require.New(t)

	mdb := db.NewMemDB()
	tree, err := NewMutableTree(mdb, 0)
	require.NoError(err)

	const maxLength = 100
	const fromLength = 10

	versions := make([]int64, 0, maxLength)
	for count := 1; count <= maxLength; count++ {
		versions = append(versions, int64(count))
		countStr := strconv.Itoa(count)
		// Set kv pair and save version
		tree.Set([]byte("aaa"), []byte("bbb"))
		tree.Set([]byte("key"+countStr), []byte("value"+countStr))
		_, _, err = tree.SaveVersion()
		require.NoError(err, "SaveVersion should not fail")
	}

	tree, err = NewMutableTree(mdb, 0)
	require.NoError(err)
	targetVersion, err := tree.LoadVersion(int64(maxLength))
	require.NoError(err)
	require.Equal(targetVersion, int64(maxLength), "targetVersion shouldn't larger than the actual tree latest version")

	err = tree.DeleteVersionsRange(fromLength, int64(maxLength/2))
	require.NoError(err, "DeleteVersionsTo should not fail")

	for _, version := range versions[:fromLength-1] {
		require.True(tree.versions.Get(version), "versions %d no more than 10 should exist", version)

		v, err := tree.LazyLoadVersion(version)
		require.NoError(err, version)
		require.Equal(v, version)

		_, value := tree.Get([]byte("aaa"))
		require.Equal(string(value), "bbb")

		for _, count := range versions[:version] {
			countStr := strconv.Itoa(int(count))
			_, value := tree.Get([]byte("key" + countStr))
			require.Equal(string(value), "value"+countStr)
		}
	}

	for _, version := range versions[fromLength : int64(maxLength/2)-1] {
		require.False(tree.versions.Get(version), "versions %d more 10 and no more than 50 should have been deleted", version)

		_, err := tree.LazyLoadVersion(version)
		require.Error(err)
	}

	for _, version := range versions[int64(maxLength/2)-1:] {
		require.True(tree.versions.Get(version), "versions %d more than 50 should exist", version)

		v, err := tree.LazyLoadVersion(version)
		require.NoError(err)
		require.Equal(v, version)

		_, value := tree.Get([]byte("aaa"))
		require.Equal(string(value), "bbb")

		for _, count := range versions[:fromLength] {
			countStr := strconv.Itoa(int(count))
			_, value := tree.Get([]byte("key" + countStr))
			require.Equal(string(value), "value"+countStr)
		}
		for _, count := range versions[int64(maxLength/2)-1 : version] {
			countStr := strconv.Itoa(int(count))
			_, value := tree.Get([]byte("key" + countStr))
			require.Equal(string(value), "value"+countStr)
		}
	}
}

func TestMutableTree_InitialVersion(t *testing.T) {
	memDB := db.NewMemDB()
	tree, err := NewMutableTreeWithOpts(memDB, 0, &Options{InitialVersion: 9})
	require.NoError(t, err)

	tree.Set([]byte("a"), []byte{0x01})
	_, version, err := tree.SaveVersion()
	require.NoError(t, err)
	assert.EqualValues(t, 10, version)

	tree.Set([]byte("b"), []byte{0x02})
	_, version, err = tree.SaveVersion()
	require.NoError(t, err)
	assert.EqualValues(t, 11, version)

	// Reloading the tree with the same initial version is fine
	tree, err = NewMutableTreeWithOpts(memDB, 0, &Options{InitialVersion: 9})
	require.NoError(t, err)
	version, err = tree.Load()
	require.NoError(t, err)
	assert.EqualValues(t, 11, version)

	// Reloading the tree with an initial version beyond the lowest should error
	tree, err = NewMutableTreeWithOpts(memDB, 0, &Options{InitialVersion: 11})
	require.NoError(t, err)
	_, err = tree.Load()
	require.Error(t, err)

	// Reloading the tree with a lower initial version is fine, and new versions can be produced
	tree, err = NewMutableTreeWithOpts(memDB, 0, &Options{InitialVersion: 3})
	require.NoError(t, err)
	version, err = tree.Load()
	require.NoError(t, err)
	assert.EqualValues(t, 11, version)

	tree.Set([]byte("c"), []byte{0x03})
	_, version, err = tree.SaveVersion()
	require.NoError(t, err)
	assert.EqualValues(t, 12, version)
}

func TestMutableTree_SetInitialVersion(t *testing.T) {
	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 0)
	require.NoError(t, err)
	tree.SetInitialVersion(9)

	tree.Set([]byte("a"), []byte{0x01})
	_, version, err := tree.SaveVersion()
	require.NoError(t, err)
	assert.EqualValues(t, 10, version)
}

func BenchmarkMutableTree_Set(b *testing.B) {
	db := db.NewDB("test", db.MemDBBackend, "")
	t, err := NewMutableTree(db, 100000)
	require.NoError(b, err)
	for i := 0; i < 1000000; i++ {
		t.Set(randBytes(10), []byte{})
	}
	b.ReportAllocs()
	runtime.GC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		t.Set(randBytes(10), []byte{})
	}
}

func TodoTestSaveVersion(t *testing.T) {
	EnableAsyncCommit = true
	defer func() {
		EnableAsyncCommit = false
	}()
	originData := make(map[string]string)
	for i := 0; i < 100; i++ {
		key := randstr(5)
		value := randstr(5)
		originData[key] = value
	}

	k1 := "k1"
	k2 := "k2"
	k3 := "k3"
	value := "Fred"
	originData[k1] = value
	originData[k2] = value
	originData[k3] = value

	modifiedData := make(map[string]string)
	for k, v := range originData {
		modifiedData[k] = v
	}
	modifiedValue := "hhhhh"
	modifiedData[k1] = modifiedValue
	modifiedData[k2] = modifiedValue
	modifiedData[k3] = modifiedValue

	testTree := func(data map[string]string, tree *ImmutableTree) {
		for k, v := range data {
			_, value := tree.Get([]byte(k))
			require.Equal(t, value, []byte(v))
		}
	}

	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 100)
	require.NoError(t, err)

	//_, _, err = tree.SaveVersion()
	//require.NoError(t, err)
	for k, v := range originData {
		tree.set([]byte(k), []byte(v))
	}
	_, _, err = tree.SaveVersion()
	require.NoError(t, err)
	oldVersion := tree.version
	tree.Set([]byte(k1), []byte(modifiedValue))
	tree.Set([]byte(k2), []byte(modifiedValue))
	tree.Set([]byte(k3), []byte(modifiedValue))
	tree.Remove([]byte(k1))
	delete(modifiedData, k1)

	_, _, err = tree.SaveVersion()
	require.NoError(t, err)

	oldTree, err := tree.GetImmutable(oldVersion)
	require.NoError(t, err)

	newTree, err := tree.GetImmutable(tree.version)
	//require.Equal(t, oldTree.nodeSize(), newTree.nodeSize())
	testTree(originData, oldTree)
	testTree(modifiedData, newTree)

	for i := 0; i < 10; i++ {
		_, _, err = tree.SaveVersion()
		require.NoError(t, err)
	}
	for i := 0; i < 200; i++ {
		_, _, err = tree.SaveVersion()
		require.NoError(t, err)
		for j := 0; j < 8; j++ {
			tree, err := tree.GetImmutable(tree.version - int64(j))
			require.NoError(t, err)
			testTree(modifiedData, tree)
		}

	}

}

func TodoTestSaveVersionCommitIntervalHeight(t *testing.T) {
	EnableAsyncCommit = true
	defer func() {
		EnableAsyncCommit = false
	}()
	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 10000)
	require.NoError(t, err)
	_, _, err = tree.SaveVersion()
	require.NoError(t, err)
	k1 := "k1"
	k2 := "k2"
	k3 := "k3"
	tree.Set([]byte(k1), []byte("v1"))
	tree.Set([]byte(k2), []byte("v2"))
	tree.Set([]byte(k3), []byte("v3"))
	//    k1
	//  k1   k2
	//     k2  k3

	_, _, err = tree.SaveVersion()
	require.NoError(t, err)
	tree.Set([]byte(k2), []byte("k22"))
	_, _, err = tree.SaveVersion()

	require.Equal(t, 5, len(tree.ndb.prePersistNodeCache)+len(tree.ndb.nodeCache))
	require.Equal(t, 3, len(tree.ndb.orphanNodeCache))

	_, _, err = tree.SaveVersion()
	require.NoError(t, err)
	require.NoError(t, err)
	for i := 0; i < 96; i++ {
		_, _, err = tree.SaveVersion()
		require.NoError(t, err)
	}

	_, _, err = tree.SaveVersion()
	require.NoError(t, err)
	require.Equal(t, 0, len(tree.ndb.prePersistNodeCache))
	require.Equal(t, 0, len(tree.ndb.orphanNodeCache))

	//require.Equal(t, 5, len(tree.ndb.nodeCache)+len(tree.ndb.tempPrePersistNodeCache))
	tree.Set([]byte("k5"), []byte("5555555555"))
	for i := 0; i < 98; i++ {
		_, _, err = tree.SaveVersion()
		require.NoError(t, err)
	}

	_, _, err = tree.SaveVersion()
	require.NoError(t, err)

}

func TodoTestConcurrentGetNode(t *testing.T) {
	EnableAsyncCommit = true
	defer func() {
		EnableAsyncCommit = false
	}()
	originData := make(map[string]string)
	var dataKey []string
	var dataLock sync.RWMutex
	for i := 0; i < 10000; i++ {
		key := randstr(5)
		value := randstr(5)
		originData[key] = value
		dataKey = append(dataKey, key)
	}

	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 100)
	require.NoError(t, err)

	//_, _, err = tree.SaveVersion()
	//require.NoError(t, err)
	for k, v := range originData {
		tree.set([]byte(k), []byte(v))
	}
	_, _, err = tree.SaveVersion()
	require.NoError(t, err)
	wg := sync.WaitGroup{}
	const num = 50
	wg.Add(num)
	go func() {
		for i := 0; i < num; i++ {
			go func() {
				queryTree, newErr := tree.GetImmutable(tree.version)
				require.Nil(t, newErr)
				idx := rand.Int() % len(dataKey)
				_, value := queryTree.Get([]byte(dataKey[idx]))
				dataLock.RLock()
				if originData[string(dataKey[idx])] != string(value) {
					//fmt.Println("not equal", originData[string(dataKey[idx])], string(value))
					time.Sleep(time.Millisecond * 10)
				}
				dataLock.RUnlock()
				_, value = queryTree.Get([]byte(dataKey[idx]))
				dataLock.RLock()
				require.Equal(t, originData[string(dataKey[idx])], string(value))
				dataLock.RUnlock()
				wg.Done()
			}()
		}
	}()
	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond * 10)
		for j := 0; j < 100; j++ {
			key := randstr(5)
			value := randstr(5)
			dataLock.Lock()
			originData[key] = value
			dataLock.Unlock()
			tree.set([]byte(key), []byte(value))

		}
		_, _, err = tree.SaveVersion()
		require.NoError(t, err)

	}
	wg.Wait()
}

// todo
func TodoTestShareNode(t *testing.T) {
	EnableAsyncCommit = true
	defer func() {
		EnableAsyncCommit = false
	}()
	CommitIntervalHeight = 10
	memDB := db.NewMemDB()
	tree, err := NewMutableTree(memDB, 10000)
	require.NoError(t, err)
	_, _, err = tree.SaveVersion()
	require.NoError(t, err)

	k1 := "k1"
	k2 := "k2"
	k3 := "k3"
	tree.Set([]byte(k1), []byte("v1"))
	tree.Set([]byte(k2), []byte("v2"))
	tree.Set([]byte(k3), []byte("v3"))
	//    k1
	//  k1   k2
	//     k2  k3

	_, oldVersion, err := tree.SaveVersion()
	require.NoError(t, err)
	tree.Set([]byte(k2), []byte("k2new"))
	_, _, err = tree.SaveVersion()

	oldTree, err := tree.GetImmutable(oldVersion)
	require.NoError(t, err)

	oldK1Node := oldTree.root.getLeftNode(oldTree)
	newK1Node := tree.root.getLeftNode(tree.ImmutableTree)
	require.Equal(t, oldK1Node, newK1Node)
	nodeDBK1Node := tree.ndb.GetNode(oldK1Node.hash)
	require.Equal(t, oldK1Node, nodeDBK1Node)

	for i := 0; i < 10; i++ {
		_, _, err = tree.SaveVersion()
		require.NoError(t, err)
	}
	oldK1Node = oldTree.root.getLeftNode(oldTree)
	newK1Node = tree.root.getLeftNode(tree.ImmutableTree)
	require.Equal(t, oldK1Node, newK1Node)
	nodeDBK1Node = tree.ndb.GetNode(oldK1Node.hash)
	require.Equal(t, oldK1Node, nodeDBK1Node)
}

func TestParseDBName(t *testing.T) {
	str := "staking"
	memDB := db.NewMemDB()
	prefixDB := db.NewPrefixDB(memDB, []byte(str))

	result := ParseDBName(prefixDB)
	require.Equal(t, str, result)

	result2 := ParseDBName(memDB)
	require.Equal(t, "", result2)
}