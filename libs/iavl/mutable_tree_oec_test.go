package iavl

import (
	"fmt"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestSaveVersion(t *testing.T) {
	EnableAsyncCommit = true
	defer func() {
		EnableAsyncCommit = false
		treeMap.resetMap()
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

	tree := newTestTree(t, false, 100, "test")

	//_, _, err = tree.SaveVersion()
	//require.NoError(t, err)
	for k, v := range originData {
		tree.Set([]byte(k), []byte(v))
	}
	_, _, _, err := tree.SaveVersion(false)
	require.NoError(t, err)
	oldVersion := tree.version
	tree.Set([]byte(k1), []byte(modifiedValue))
	tree.Set([]byte(k2), []byte(modifiedValue))
	tree.Set([]byte(k3), []byte(modifiedValue))
	tree.Remove([]byte(k1))
	delete(modifiedData, k1)

	_, _, _, err = tree.SaveVersion(false)
	require.NoError(t, err)

	oldTree, err := tree.GetImmutable(oldVersion)
	require.NoError(t, err)

	newTree, err := tree.GetImmutable(tree.version)
	//require.Equal(t, oldTree.nodeSize(), newTree.nodeSize())
	testTree(originData, oldTree)
	testTree(modifiedData, newTree)

	for i := 0; i < 10; i++ {
		_, _, _, err = tree.SaveVersion(false)
		require.NoError(t, err)
	}
	for i := 0; i < 200; i++ {
		_, _, _, err = tree.SaveVersion(false)
		require.NoError(t, err)
		for j := 0; j < 8; j++ {
			tree, err := tree.GetImmutable(tree.version - int64(j))
			require.NoError(t, err)
			testTree(modifiedData, tree)
		}

	}

}

func TestSaveVersionCommitIntervalHeight(t *testing.T) {
	EnableAsyncCommit = true
	defer func() {
		EnableAsyncCommit = false
		treeMap.resetMap()
	}()
	tree := newTestTree(t, false, 10000, "test")

	_, _, _, err := tree.SaveVersion(false)
	require.NoError(t, err)
	keys, _ := initSetTree(tree)
	_, k2, _ := keys[0], keys[1], keys[2]

	_, _, _, err = tree.SaveVersion(false)
	require.NoError(t, err)
	tree.Set([]byte(k2), []byte("k22"))
	_, _, _, err = tree.SaveVersion(false)

	require.Equal(t, 5, len(tree.ndb.prePersistNodeCache)+len(tree.ndb.nodeCache))
	require.Equal(t, 3, len(tree.ndb.orphanNodeCache))

	_, _, _, err = tree.SaveVersion(false)
	require.NoError(t, err)
	require.NoError(t, err)
	for i := 0; i < 96; i++ {
		_, _, _, err = tree.SaveVersion(false)
		require.NoError(t, err)
	}

	_, _, _, err = tree.SaveVersion(false)
	require.NoError(t, err)
	require.Equal(t, 0, len(tree.ndb.prePersistNodeCache))
	require.Equal(t, 0, len(tree.ndb.orphanNodeCache))

	//require.Equal(t, 5, len(tree.ndb.nodeCache)+len(tree.ndb.tempPrePersistNodeCache))
	tree.Set([]byte("k5"), []byte("5555555555"))
	for i := 0; i < 98; i++ {
		_, _, _, err = tree.SaveVersion(false)
		require.NoError(t, err)
	}

	_, _, _, err = tree.SaveVersion(false)
	require.NoError(t, err)

}

func TestConcurrentGetNode(t *testing.T) {
	EnableAsyncCommit = true
	defer func() {
		EnableAsyncCommit = false
		treeMap.resetMap()
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

	tree := newTestTree(t, false, 10000, "test")

	//_, _, err = tree.SaveVersion()
	//require.NoError(t, err)
	for k, v := range originData {
		tree.Set([]byte(k), []byte(v))
	}
	_, _, _, err := tree.SaveVersion(false)
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
			tree.Set([]byte(key), []byte(value))

		}
		_, _, _, err = tree.SaveVersion(false)
		require.NoError(t, err)

	}
	wg.Wait()
}

func TestShareNode(t *testing.T) {
	EnableAsyncCommit = true
	CommitIntervalHeight = 10
	defer func() {
		EnableAsyncCommit = false
		CommitIntervalHeight = 100
		treeMap.resetMap()
	}()

	tree := newTestTree(t, false, 10000, "test")

	_, _, _, err := tree.SaveVersion(false)
	require.NoError(t, err)

	keys, _ := initSetTree(tree)
	_, k2, _ := keys[0], keys[1], keys[2]

	_, oldVersion, _, err := tree.SaveVersion(false)
	require.NoError(t, err)
	tree.Set([]byte(k2), []byte("k2new"))
	_, _, _, err = tree.SaveVersion(false)

	oldTree, err := tree.GetImmutable(oldVersion)
	require.NoError(t, err)

	oldK1Node := oldTree.root.getLeftNode(oldTree)
	newK1Node := tree.root.getLeftNode(tree.ImmutableTree)
	require.Equal(t, oldK1Node, newK1Node)
	nodeDBK1Node := tree.ndb.GetNode(oldK1Node.hash)
	require.Equal(t, oldK1Node, nodeDBK1Node)

	for i := 0; i < 10; i++ {
		_, _, _, err = tree.SaveVersion(false)
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

func TestPruningHistoryState(t *testing.T) {
	EnableAsyncCommit = true
	EnablePruningHistoryState = true
	defer func() {
		EnableAsyncCommit = false
		EnablePruningHistoryState = false
		treeMap.resetMap()
	}()
	tree := newTestTree(t, false, 10000, "test")
	keys, _ := initSetTree(tree)
	_, k2, _ := keys[0], keys[1], keys[2]

	_, _, _, err := tree.SaveVersion(false)
	require.NoError(t, err)

	batchSaveVersion(t, tree, int(CommitIntervalHeight))

	v2New := []byte("v22")
	tree.Set(k2, v2New)
	_, _, _, err = tree.SaveVersion(false)
	require.NoError(t, err)

	batchSaveVersion(t, tree, minHistoryStateNum*int(CommitIntervalHeight)-2)

	tree.commitCh <- commitEvent{-1, nil, nil, nil, nil, 0}

	iTree, err := tree.GetImmutable(CommitIntervalHeight * (minHistoryStateNum - 1))
	require.NoError(t, err)
	require.NotNil(t, iTree)
	_, v := iTree.Get(k2)
	require.Equal(t, v2New, v)

	iTree, err = tree.GetImmutable(CommitIntervalHeight * 1)
	require.Error(t, err)
	require.Nil(t, iTree)

	nodeCount := 0
	tree.ndb.traverseNodes(func(hash []byte, node *Node) {
		nodeCount++
	})
	require.Equal(t, 5, nodeCount)

	orphansCount := 0
	tree.ndb.traverseOrphans(func(k, v []byte) {
		orphansCount++
	})
	require.Equal(t, 0, orphansCount)
}

func batchSaveVersion(t *testing.T, tree *MutableTree, n int) {
	for i := 0; i < n; i++ {
		_, _, _, err := tree.SaveVersion(false)
		require.NoError(t, err)
	}
}

func openLog(moduleName string) {
	SetLogFunc(func(level int, format string, args ...interface{}) {
		if level == IavlInfo {
			fmt.Printf(format, args...)
			fmt.Printf("\n")
		}
	})
	OutputModules = make(map[string]int)
	OutputModules[moduleName] = 1
}

func TestPruningHistoryStateRandom(t *testing.T) {
	EnableAsyncCommit = true
	EnablePruningHistoryState = true
	defer func() {
		EnableAsyncCommit = false
		EnablePruningHistoryState = false
		treeMap.resetMap()
	}()
	tree := newTestTree(t, false, 10000, "test")
	keys, _ := initSetTree(tree)
	k1, k2, k3 := keys[0], keys[1], keys[2]

	_, _, _, err := tree.SaveVersion(false)
	require.NoError(t, err)

	for i := 0; i < 10000; i++ {
		tree.Set(k2, randBytes(i%64+1))
		_, _, _, err := tree.SaveVersion(false)
		require.NoError(t, err)
	}

	tree.commitCh <- commitEvent{-1, nil, nil, nil, nil, 0}

	nodeCount := 0
	tree.ndb.traverseNodes(func(hash []byte, node *Node) {
		nodeCount++
	})
	require.Equal(t, (minHistoryStateNum-1)*3+5, nodeCount)

	orphansCount := 0
	tree.ndb.traverseOrphans(func(k, v []byte) {
		orphansCount++
	})
	require.Equal(t, (minHistoryStateNum-1)*3, orphansCount)

	for i := 0; i < 10000; i++ {
		tree.Set(k1, randBytes(i%64+1))
		tree.Set(k2, randBytes(i%64+1))
		tree.Set(k3, randBytes(i%64+1))
		_, _, _, err := tree.SaveVersion(false)
		require.NoError(t, err)
	}

	nodeCount = 0
	tree.ndb.traverseNodes(func(hash []byte, node *Node) {
		nodeCount++
	})
	require.Equal(t, minHistoryStateNum*5, nodeCount)

	orphansCount = 0
	tree.ndb.traverseOrphans(func(k, v []byte) {
		orphansCount++
	})
	require.Equal(t, (minHistoryStateNum-1)*5, orphansCount)
}

func TestConcurrentQuery(t *testing.T) {
	EnableAsyncCommit = true
	EnablePruningHistoryState = true
	CommitIntervalHeight = 5
	defer func() {
		EnableAsyncCommit = false
		EnablePruningHistoryState = false
		CommitIntervalHeight = 100
		treeMap.resetMap()
	}()
	originData := make(map[string]string)
	var dataKey []string
	for i := 0; i < 100000; i++ {
		key := randstr(5)
		value := randstr(5)
		originData[key] = value
		dataKey = append(dataKey, key)
	}

	tree := newTestTree(t, false, 10000, "test")

	for k, v := range originData {
		tree.Set([]byte(k), []byte(v))
	}
	_, _, _, err := tree.SaveVersion(false)
	require.NoError(t, err)
	const num = 1000000
	queryEnd := false
	endCh := make(chan struct{})
	go func() {
		ch := make(chan struct{}, 20)
		wg := sync.WaitGroup{}
		wg.Add(num)
		for i := 0; i < num; i++ {
			ch <- struct{}{}
			go func() {
				queryVersion := tree.version
				//fmt.Println(time.Now().String(),"query version:", queryVersion)
				queryTree, newErr := tree.GetImmutable(queryVersion)
				require.Nil(t, newErr, "query:%d current:%d\n", queryVersion, tree.version)
				idx := rand.Int() % len(dataKey)
				_, value := queryTree.Get([]byte(dataKey[idx]))
				require.NotNil(t, value)
				require.NotEqual(t, []byte{}, value)
				wg.Done()
				<-ch
			}()
		}
		wg.Wait()
		queryEnd = true
	}()
	go func() {
		for i := 0; ; i++ {
			fmt.Println(time.Now().String(), "current version:", tree.version)
			for j := 0; j < 100; j++ {
				key := dataKey[rand.Intn(len(dataKey))]
				value := randstr(5)
				originData[key] = value
				tree.Set([]byte(key), []byte(value))
			}
			_, _, _, err = tree.SaveVersion(false)
			require.NoError(t, err)
			if queryEnd {
				break
			}
		}
		endCh <- struct{}{}
	}()
	<-endCh
}

func TestStopTree(t *testing.T) {
	EnableAsyncCommit = true
	EnablePruningHistoryState = true
	defer func() {
		EnableAsyncCommit = false
		EnablePruningHistoryState = false
	}()
	tree := newTestTree(t, false, 10000, "test")
	initSetTree(tree)

	_, _, _, err := tree.SaveVersion(false)
	require.NoError(t, err)
	tree.StopTree()
	require.Equal(t, 5, len(tree.ndb.nodeCache))
}

func TestLog(t *testing.T) {
	defer func() {
		treeMap.resetMap()
	}()
	tree := newTestTree(t, false, 10000, "test")
	dbRCount := tree.GetDBReadCount()
	dbWCount := tree.GetDBWriteCount()
	nodeRCount := tree.GetNodeReadCount()
	require.Zero(t, dbRCount)
	require.Zero(t, dbWCount)
	require.Zero(t, nodeRCount)
	tree.ResetCount()
}

func initSetTree(tree *MutableTree) ([][]byte, [][]byte) {
	keys := [][]byte{
		[]byte("k1"),
		[]byte("k2"),
		[]byte("k3"),
	}
	values := [][]byte{
		[]byte("v1"),
		[]byte("v2"),
		[]byte("v3"),
	}
	for i, key := range keys {
		tree.Set(key, values[i])
	}
	//    k1
	//  k1   k2
	//      k2 k3
	return keys, values
}

func newTestTree(t *testing.T, openLogFlag bool, cacheSize int, moduleName string) *MutableTree {
	if openLogFlag {
		openLog(moduleName)
	}

	memDB := db.NewPrefixDB(db.NewMemDB(), []byte(moduleName))
	tree, err := NewMutableTree(memDB, cacheSize)
	require.NoError(t, err)
	return tree
}

func TestCommitSchedule(t *testing.T) {
	EnableAsyncCommit = true
	EnablePruningHistoryState = true
	defer func() {
		EnableAsyncCommit = false
		EnablePruningHistoryState = false
		treeMap.resetMap()
	}()
	tree := newTestTree(t, false, 10000, "test")
	initSetTree(tree)

	for i := 0; i < int(CommitIntervalHeight); i++ {
		_, _, _, err := tree.SaveVersion(false)
		require.NoError(t, err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	versions := tree.deepCopyVersions()
	batch := tree.NewBatch()
	tree.commitCh <- commitEvent{CommitIntervalHeight, versions, batch, nil, nil, 0}

	tree.commitCh <- commitEvent{CommitIntervalHeight, versions, batch, nil, &wg, 0}
	wg.Wait()
}
