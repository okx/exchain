package iavl

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
	"math/rand"
	"runtime"
	"runtime/debug"
	"testing"
)

func TestGetNodeFromPool(t *testing.T) {
	n := GetNodeFromPool()
	require.Nil(t, n)
	node := NewNode([]byte("node"), []byte("this is node"), 1)
	SetNodeToPool(node)
	n = GetNodeFromPool()
	require.NotNil(t, n)
	require.Equal(t, node, n)
	n = GetNodeFromPool()
	require.Nil(t, n)
}

func TestSetNodeToPool(t *testing.T) {
	node := NewNode([]byte("node"), []byte("this is node"), 1)
	SetNodeToPool(node)
	n := GetNodeFromPool()
	require.NotNil(t, n)
	require.Equal(t, node, n)

	SetNodeToPool(node)
	SetNodeToPool(node)

	n = GetNodeFromPool()
	require.NotNil(t, n)
	require.Equal(t, node, n)
	n = GetNodeFromPool()
	require.NotNil(t, n)
	require.Equal(t, node, n)
	n = GetNodeFromPool()
	require.Nil(t, n)
}

func TestNode_Reset(t *testing.T) {
	leafNode1 := NewNode([]byte("leafNode1"), []byte("this is leafNode1"), 1)
	leafNode2 := NewNode([]byte("leafNode2"), []byte("this is leafNode2"), 1)
	leafNode1.Reset(leafNode2.key, leafNode2.value, nil, nil, nil, leafNode2.version, leafNode2.size, nil, nil, leafNode2.height, false, false)

	require.Equal(t, leafNode2, leafNode1)
}

func BenchmarkMakeNode_Alloc(b *testing.B) {
	leafNode1 := NewNode([]byte("leafNode1"), []byte("this is leafNode1"), 1)
	buffer := new(bytes.Buffer)
	leafNode1.writeBytes(buffer)
	data := buffer.Bytes()
	b.Run("MakeNode", func(b *testing.B) {
		b.ResetTimer()
		invalidOp := 0
		for i := 0; i < b.N; i++ {
			node := NewNode([]byte(fmt.Sprintf("leafNode%d", i)), []byte("this is leafNode1"), 1)
			node = GetNodeFromPool()
			if node != nil {
				invalidOp++
				continue
			}
			_, err := MakeNode(data)
			require.NoError(b, err)
		}
		b.ReportAllocs()
		if invalidOp > 10 {
			fmt.Println("invalid op", invalidOp)
		}
	})

	b.Run("MakeNodeGC", func(b *testing.B) {
		b.ResetTimer()
		invalidOp := 0
		for i := 0; i < b.N; i++ {
			node := NewNode([]byte(fmt.Sprintf("leafNode%d", i)), []byte("this is leafNode1"), 1)
			SetNodeToPool(node)
			node = GetNodeFromPool()
			if node == nil {
				invalidOp++
				continue
			}
			_, err := MakeNodeForGC(leafNode1, data)
			require.NoError(b, err)
		}
		b.ReportAllocs()
		if invalidOp > 10 {
			fmt.Println("invalid op", invalidOp)
		}
	})
}

func BenchmarkNode_Clone_GC_Compare_MemStats(b *testing.B) {
	h := []byte{1, 2, 3}
	node := &Node{key: []byte("test"), value: nBytes(1024 * 1024), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
	dataSize := 1024 * 1024
	data := make([]*Node, 0)

	for i := 0; i < dataSize; i++ {
		temp := &Node{key: []byte(fmt.Sprintf("innerNode%d", i)), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
		data = append(data, temp)
	}

	debug.SetGCPercent(-1)

	b.Run("nodeCloneGC", func(b *testing.B) {
		// reset
		for i := 0; i < dataSize; i++ {
			SetNodeToPool(data[i])
		}
		EnableNodePool = true
		nodeCloneTest(b, node)
	})

	b.Run("nodeClone", func(b *testing.B) {
		// reset
		for i := 0; i < dataSize; i++ {
			SetNodeToPool(data[i])
		}
		EnableNodePool = false
		nodeCloneTest(b, node)
	})
}

func nodeCloneTest(b *testing.B, node *Node) {
	before := getMemStats()
	oldGet, oldSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
	defer func() {
		b.ReportAllocs()
		after := getMemStats()
		newGet, newSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
		b.Logf("%s : Pool<%v> , GetNodePoolCounter<%d>, SetNodePoolCounter<%d>, Alloc<%dKB>, GC:%d", b.Name(), EnableNodePool, newGet-oldGet, newSet-oldSet, int64(after.Alloc/1024)-int64(before.Alloc/1024), after.NumGC-before.NumGC)
	}()
	temps1 := make([]*Node, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		temp := node.clone(1)
		temps1 = append(temps1, temp)
	}
}

func BenchmarkMakeNode_GC_Compare_MemStats(b *testing.B) {
	h := []byte{1, 2, 3}
	node := &Node{key: []byte("test"), value: nBytes(1024 * 1024), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
	buffer := new(bytes.Buffer)
	err := node.writeBytes(buffer)
	require.NoError(b, err)
	dataSize := 1024 * 1024
	data := make([]*Node, 0)

	for i := 0; i < dataSize; i++ {
		temp := &Node{key: []byte(fmt.Sprintf("innerNode%d", i)), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
		data = append(data, temp)
	}

	debug.SetGCPercent(-1)

	b.Run("MakeNodeGC", func(b *testing.B) {
		// reset
		for i := 0; i < dataSize; i++ {
			SetNodeToPool(data[i])
		}
		EnableNodePool = true
		makeCodeTest(b, buffer.Bytes())
	})

	b.Run("MakeNode", func(b *testing.B) {
		// reset
		for i := 0; i < dataSize; i++ {
			SetNodeToPool(data[i])
		}
		EnableNodePool = false
		makeCodeTest(b, buffer.Bytes())
	})
}

func makeCodeTest(b *testing.B, buff []byte) {
	before := getMemStats()
	oldGet, oldSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
	defer func() {
		b.ReportAllocs()
		after := getMemStats()
		newGet, newSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
		b.Logf("%s : Pool<%v> , GetNodePoolCounter<%d>, SetNodePoolCounter<%d>, Alloc<%dKB>, GC:%d", b.Name(), EnableNodePool, newGet-oldGet, newSet-oldSet, int64(after.Alloc/1024)-int64(before.Alloc/1024), after.NumGC-before.NumGC)
	}()
	temps1 := make([]*Node, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		temp, err := MakeNode(buff)
		require.NoError(b, err)
		temps1 = append(temps1, temp)
	}
}

func BenchmarkNewNode_GC_Compare_MemStats(b *testing.B) {
	h := []byte{1, 2, 3}
	node := &Node{key: []byte("test"), value: nBytes(1024 * 1024), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
	buffer := new(bytes.Buffer)
	err := node.writeBytes(buffer)
	require.NoError(b, err)
	dataSize := 1024 * 1024 * 32
	data := make([]*Node, 0)

	for i := 0; i < dataSize; i++ {
		temp := &Node{key: []byte(fmt.Sprintf("innerNode%d", i)), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
		data = append(data, temp)
	}

	debug.SetGCPercent(-1)

	b.Run("NewCode", func(b *testing.B) {
		// reset
		for i := 0; i < dataSize; i++ {
			SetNodeToPool(data[i])
		}
		EnableNodePool = false
		newCodeTest(b, node)
	})

	b.Run("NewCodeGC", func(b *testing.B) {
		// reset
		for i := 0; i < dataSize; i++ {
			SetNodeToPool(data[i])
		}
		EnableNodePool = true
		newCodeTest(b, node)
	})
}

func newCodeTest(b *testing.B, node *Node) {
	before := getMemStats()
	oldGet, oldSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
	defer func() {
		b.ReportAllocs()
		after := getMemStats()
		newGet, newSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
		b.Logf("%s : Pool<%v> , GetNodePoolCounter<%d>, SetNodePoolCounter<%d>, Alloc<%dKB>, GC:%d", b.Name(), EnableNodePool, newGet-oldGet, newSet-oldSet, int64(after.Alloc/1024)-int64(before.Alloc/1024), after.NumGC-before.NumGC)
	}()
	temps1 := make([]*Node, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		temp := NewNode(node.key, node.value, node.version)
		temps1 = append(temps1, temp)
	}
}

func BenchmarkNodeJsonToNode_GC_Compare_MemStats(b *testing.B) {
	h := []byte{1, 2, 3}
	node := &Node{key: []byte("test"), value: nBytes(1024 * 1024), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
	nodeJson := NodeToNodeJson(node)
	buffer := new(bytes.Buffer)
	err := node.writeBytes(buffer)
	require.NoError(b, err)
	dataSize := 1024 * 1024 * 32
	data := make([]*Node, 0)

	for i := 0; i < dataSize; i++ {
		temp := &Node{key: []byte(fmt.Sprintf("innerNode%d", i)), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
		data = append(data, temp)
	}

	debug.SetGCPercent(-1)

	b.Run("NodeJsonToNode", func(b *testing.B) {
		// reset
		for i := 0; i < dataSize; i++ {
			SetNodeToPool(data[i])
		}
		EnableNodePool = false
		memStatsTest(b, nodeJson)
	})
	b.Run("NodeJsonToNodeGC", func(b *testing.B) {
		// reset
		for i := 0; i < dataSize; i++ {
			SetNodeToPool(data[i])
		}
		EnableNodePool = true
		memStatsTest(b, nodeJson)
	})
}

func memStatsTest(b *testing.B, node *NodeJson) {
	before := getMemStats()
	oldGet, oldSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
	defer func() {
		b.ReportAllocs()
		after := getMemStats()
		newGet, newSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
		b.Logf("%s : Pool<%v> , GetNodePoolCounter<%d>, SetNodePoolCounter<%d>, Alloc<%dKB>, GC:%d", b.Name(), EnableNodePool, newGet-oldGet, newSet-oldSet, int64(after.Alloc/1024)-int64(before.Alloc/1024), after.NumGC-before.NumGC)
	}()
	temps1 := make([]*Node, 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		temp := NodeJsonToNode(node)
		temps1 = append(temps1, temp)
	}
}

func BenchmarkTreeSet_GC_Compare_MemStats(b *testing.B) {
	EnableAsyncCommit = true
	EnablePruningHistoryState = true
	CommitIntervalHeight = 1
	defer func() {
		EnableAsyncCommit = false
		EnablePruningHistoryState = false
		CommitIntervalHeight = 100
	}()
	testCases := []struct {
		name           string
		enableNodePool bool
		initTreeHeight int8
		setKVSize      int
	}{
		{
			name:           "test1",
			enableNodePool: true,
			initTreeHeight: 24,
			setKVSize:      10000,
		},
		{
			name:           "test2",
			enableNodePool: false,
			initTreeHeight: 24,
			setKVSize:      10000,
		},
	}

	for i, testCase := range testCases {
		fmt.Println("test case", i, ": ", testCase.name)
		benchmarkTreeSet(b, testCase.name, testCase.enableNodePool, testCase.initTreeHeight, testCase.setKVSize)
	}
}

func benchmarkTreeSet(b *testing.B, name string, enableNodePool bool, treeHeight int8, kvSize int) {
	b.Run(name, func(b *testing.B) {
		debug.SetGCPercent(-1)
		_, _, tree, err := mockGCSepcialHeightTree(treeHeight)
		require.NoError(b, err)
		EnableNodePool = true
		for i := 0; i < kvSize*1024; i++ {
			SetNodeToPool(&Node{key: []byte(fmt.Sprintf("innerNode%d", i)), version: 1, size: 1, height: 1, leftHash: []byte("test"), rightHash: []byte("test1")})
		}
		EnableNodePool = enableNodePool
		data := mockKVData(b.N)
		multiTreeSet(b, tree, data)

	})
}

func multiTreeSet(b *testing.B, tree *MutableTree, data []testKVData) {
	before := getMemStats()
	oldGet, oldSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
	defer func() {
		b.ReportAllocs()
		after := getMemStats()
		newGet, newSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
		b.Logf("after tree set: Pool<%v> , GetNodePoolCounter<%d>, SetNodePoolCounter<%d>, Alloc<%dKB>, GC:%d", EnableNodePool, newGet-oldGet, newSet-oldSet, int64(after.Alloc/1024)-int64(before.Alloc/1024), after.NumGC-before.NumGC)
	}()
	size := len(data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index := i % size
		tree.Set(data[index].key, data[index].value)
	}
}

type testKVData struct {
	key   []byte
	value []byte
}

func mockKVData(size int) []testKVData {
	data := make([]testKVData, 0)
	for i := 0; i < size; i++ {
		data = append(data, testKVData{key: nBytes(20), value: nBytes(20)})
	}
	return data
}

func mockGCSepcialHeightTree(height int8) (keySet []string, dataSet map[string]string, tree *MutableTree, err error) {
	tree, err = mockGCTree()
	if err != nil {
		return nil, nil, nil, err
	}
	dataSet = make(map[string]string)
	keySet = make([]string, 0)
	for i := 0; tree.root == nil || tree.root.height < height; i++ {
		key := randstr(20)
		value := randstr(2)
		dataSet[key] = string(value)
		keySet = append(keySet, key)
		tree.Set([]byte(key), []byte(value))
	}
	//PrintNode("ada", tree.ndb, tree.root)
	//tree.SaveVersion(false)
	//tree.commitCh <- commitEvent{-1, nil, nil, nil, nil, 0}
	return
}
func mockGCTree() (*MutableTree, error) {
	memDB := db.NewPrefixDB(db.NewMemDB(), []byte("mockMemTree"))
	return NewMutableTree(memDB, 0)
}

func nBytes(n int) []byte {
	buf := make([]byte, n)
	n, _ = rand.Read(buf)
	return buf[:n]
}

func getMemStats() (m runtime.MemStats) {
	runtime.ReadMemStats(&m)
	return m
}
