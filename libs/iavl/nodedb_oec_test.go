package iavl

import (
	"container/list"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
)

func mockNodeDB() *nodeDB {
	memDB := db.NewMemDB()
	return newNodeDB(memDB, 10000, nil)
}

func mockNode(version int64) *Node {
	key, value := randBytes(32), randBytes(32)
	return NewNode(key, value, version)
}

func mockNodes(version int64, capacity int) (*Node, *list.List) {
	nodelist := list.New()

	root := mockNode(version)
	nodelist.PushBack(root)
	count := 1
	for count < capacity {
		lastNode := nodelist.Remove(nodelist.Front()).(*Node)
		if count < capacity {
			node := mockNode(version)
			count++
			lastNode.leftNode = node
			nodelist.PushBack(node)
		}

		if count < capacity {
			node := mockNode(version)
			count++
			lastNode.rightNode = node
			nodelist.PushBack(node)
		}
	}
	return root, nodelist
}

func Test_saveNodeToPrePersistCache(t *testing.T) {
	EnableAsyncCommit = true
	defer func() { EnableAsyncCommit = false }()

	cases := []struct {
		version  int64
		nodeNums int
	}{
		{100, 1},
		{200, 1000},
		{300, 100000},
	}

	ndb := mockNodeDB()
	for _, c := range cases {
		nodes := make([]*Node, c.nodeNums, c.nodeNums)
		for i := 0; i < c.nodeNums; i++ {
			node := mockNode(c.version)
			node.hash = node._hash()
			nodes[i] = node
			ndb.saveNodeToPrePersistCache(node)
		}

		for i := 0; i < c.nodeNums; i++ {
			require.True(t, nodes[i].prePersisted)
			require.NotNil(t, ndb.prePersistNodeCache[string(nodes[i].hash)])
		}
	}
}

func Test_batchSet(t *testing.T) {
	EnableAsyncCommit = true
	defer func() { EnableAsyncCommit = false }()

	cases := []struct {
		version  int64
		nodeNums int
	}{
		{100, 1},
		{200, 1000},
		{300, 100000},
	}
	judges := []struct {
		hashExisted  bool
		persisted    bool
		prePersisted bool
		panic        bool
	}{
		{true, true, true, false},
		{true, false, true, true},
		{false, false, true, true},
		{true, true, false, true},
		{false, true, false, true},
		{false, false, false, true},
	}

	ndb := mockNodeDB()
	for _, c := range cases {
		for _, g := range judges {

			batch := ndb.NewBatch()
			nodes := make([]*Node, c.nodeNums, c.nodeNums)
			for i := 0; i < c.nodeNums; i++ {
				nodes[i] = mockNode(c.version)
				if g.hashExisted {
					nodes[i].hash = nodes[i]._hash()
				}
				if g.persisted {
					nodes[i].persisted = true
				}
				if g.prePersisted {
					nodes[i].prePersisted = true
				}
				if g.panic {
					require.Panics(t, func() { ndb.batchSet(nodes[i], batch) })
				} else {
					require.NotPanics(t, func() { ndb.batchSet(nodes[i], batch) })
				}
			}

			if !g.panic {
				require.NoError(t, ndb.Commit(batch))

				for i := 0; i < c.nodeNums; i++ {
					require.NotNil(t, nodes[i].hash)
					get, err := ndb.dbGet(ndb.nodeKey(nodes[i].hash))
					require.NoError(t, err)
					require.NotEmpty(t, get)
				}
			}
		}
	}
}

func Test_updateBranch(t *testing.T) {
	EnableAsyncCommit = true
	defer func() { EnableAsyncCommit = false }()

	cases := []struct {
		version  int64
		nodeNums int
	}{
		{100, 1},
		{200, 10},
		{300, 100},
		{400, 1000},
		{500, 1000},
		{600, 10000},
		{700, 100000},
	}

	ndb := mockNodeDB()
	capacity := 0
	for _, c := range cases {
		capacity += c.nodeNums

		root, nodelist := mockNodes(c.version, c.nodeNums)
		ndb.updateBranch(root, map[string]*Node{})
		for elem := nodelist.Front(); elem != nil; elem = elem.Next() {
			node := elem.Value.(*Node)
			require.True(t, node.prePersisted)
			require.Nil(t, node.leftNode)
			require.Nil(t, node.rightNode)
		}
		require.Equal(t, len(ndb.prePersistNodeCache), capacity)
	}
}

func Test_updateBranchConcur(t *testing.T) {
	EnableAsyncCommit = true
	defer func() { EnableAsyncCommit = false }()

	cases := []struct {
		version  int64
		nodeNums int
	}{
		{100, 1},
		{200, 10},
		{300, 100},
		{400, 1000},
		{500, 1000},
		{600, 10000},
		{700, 100000},
	}

	ndb := mockNodeDB()
	capacity := 0
	for _, c := range cases {
		capacity += c.nodeNums

		root, nodelist := mockNodes(c.version, c.nodeNums)
		defaultConcurDepth := 4
		rootHash := make(chan []byte, 1)
		ndb.updateBranchConcur(root, map[string]*Node{}, defaultConcurDepth, rootHash)
		//fmt.Println("rootHash: ", hex.EncodeToString(<-rootHash))
		for elem := nodelist.Front(); elem != nil; elem = elem.Next() {
			node := elem.Value.(*Node)
			require.True(t, node.prePersisted)
			require.Nil(t, node.leftNode)
			require.Nil(t, node.rightNode)
		}
		require.Equal(t, len(ndb.prePersistNodeCache), capacity)
	}
}

func Test_saveCommitOrphans(t *testing.T) {
	EnableAsyncCommit = true
	defer func() { EnableAsyncCommit = false }()

	cases := []struct {
		version    int64
		orphansNum int
		exist      bool
	}{
		{100, 100, true},
		{200, 1000, true},
		{300, 10000, true},
	}

	ndb := mockNodeDB()
	for n, c := range cases {
		commitOrphans := make(map[string]int64)
		for i := 0; i < c.orphansNum; i++ {
			node := mockNode(c.version)
			commitOrphans[string(node._hash())] = rand.Int63n(100) + 100*int64(n)
		}

		batch1 := ndb.NewBatch()
		batch1.Set(ndb.rootKey(c.version), []byte("root"))
		require.NoError(t, ndb.Commit(batch1))

		batch2 := ndb.NewBatch()
		ndb.saveCommitOrphans(batch2, c.version+1, commitOrphans)
		require.NoError(t, ndb.Commit(batch2))

		for hash, fromVersion := range commitOrphans {
			key := ndb.orphanKey(fromVersion, c.version, []byte(hash))
			node, err := ndb.dbGet(key)
			require.NoError(t, err)
			require.Equal(t, []byte(hash), node)
		}
	}
}

func Test_getNodeInTpp(t *testing.T) {
	EnableAsyncCommit = true
	defer func() { EnableAsyncCommit = false }()

	cases := []struct {
		version    int64
		orphansNum int
		exist      bool
	}{
		{100, 100, true},
		{200, 1000, true},
		{300, 10000, true},
	}

	ndb := mockNodeDB()
	for _, c := range cases {
		for i := 0; i < c.orphansNum; i++ {
			node := mockNode(c.version)
			ndb.prePersistNodeCache[string(node._hash())] = node
		}

		tpp := ndb.prePersistNodeCache
		lItem := ndb.tppVersionList.PushBack(c.version)
		ndb.tppMap[c.version] = &tppItem{
			nodeMap:  tpp,
			listItem: lItem,
		}

		for hash, node := range tpp {
			getNode, found := ndb.getNodeInTpp([]byte(hash))
			require.True(t, found)
			require.EqualValues(t, node, getNode)
		}

		ndb.prePersistNodeCache = map[string]*Node{}
	}
}

func Test_getRootWithCache(t *testing.T) {
	EnableAsyncCommit = true
	defer func() { EnableAsyncCommit = false }()

	cases := []struct {
		version int64
		exist   bool
	}{
		{1, true},
		{2, true},
	}

	ndb := mockNodeDB()
	for _, c := range cases {
		rootHash := randBytes(32)
		ndb.heightOrphansMap[c.version] = &heightOrphansItem{c.version, rootHash, nil}

		actualHash, err := ndb.getRootWithCache(c.version)
		if c.exist {
			require.Equal(t, actualHash, rootHash)
		} else {
			require.Nil(t, actualHash)
		}
		require.NoError(t, err)

		actualHash, err = ndb.getRootWithCacheAndDB(c.version)
		if c.exist {
			require.Equal(t, actualHash, rootHash)
		} else {
			require.Nil(t, actualHash)
		}
		require.NoError(t, err)
	}
}

func Test_inVersionCacheMap(t *testing.T) {
	cases := []struct {
		version  int64
		expected bool
	}{
		{1, true},
		{2, true},
		{3, true},
		{4, true},
	}

	ndb := mockNodeDB()
	for _, c := range cases {
		rootHash := randBytes(32)
		orphanObj := &heightOrphansItem{version: c.version, rootHash: rootHash}
		ndb.heightOrphansMap[c.version] = orphanObj
		actualHash, existed := ndb.inVersionCacheMap(c.version)
		require.Equal(t, actualHash, rootHash)
		require.Equal(t, existed, c.expected)
	}
}
