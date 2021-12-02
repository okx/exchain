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

func TestNode_Clone_GC_Compare_MemStats(t *testing.T) {
	h := []byte{1, 2, 3}
	node := &Node{key: []byte("test"), value: nBytes(1024 * 1024), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
	dataSize := 1024 * 1024
	data := make([]*Node, 0)

	for i := 0; i < dataSize; i++ {
		temp := &Node{key: []byte(fmt.Sprintf("innerNode%d", i)), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
		data = append(data, temp)
	}

	debug.SetGCPercent(-1)

	before := getMemStats()
	temps := make([]*Node, 0)
	for i := 0; i < dataSize; i++ {
		temp := node.clone(1)
		temps = append(temps, temp)
	}
	after := getMemStats()
	t.Logf("after node clone:gc<disable>, pool<disable>: %dMB,GC:%d", int64(after.Alloc/1024/1024)-int64(before.Alloc/1024/1024), after.NumGC-before.NumGC)

	// reset
	for i := 0; i < dataSize; i++ {
		SetNodeToPool(data[i])
	}

	before = getMemStats()
	temps1 := make([]*Node, 0)
	for i := 0; i < dataSize; i++ {
		temp := node.clone(1)
		temps1 = append(temps1, temp)
	}
	after = getMemStats()
	t.Logf("after node clone:gc<disable>, pool<enable> : %dMB,GC:%d", int64(after.Alloc/1024/1024)-int64(before.Alloc/1024/1024), after.NumGC-before.NumGC)
}

func TestMakeNode_GC_Compare_MemStats(t *testing.T) {
	h := []byte{1, 2, 3}
	node := &Node{key: []byte("test"), value: nBytes(1024 * 1024), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
	buffer := new(bytes.Buffer)
	err := node.writeBytes(buffer)
	require.NoError(t, err)
	dataSize := 1024 * 1024
	data := make([]*Node, 0)

	for i := 0; i < dataSize; i++ {
		temp := &Node{key: []byte(fmt.Sprintf("innerNode%d", i)), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
		data = append(data, temp)
	}

	debug.SetGCPercent(-1)

	before := getMemStats()
	temps := make([]*Node, 0)
	for i := 0; i < dataSize; i++ {
		temp, err := MakeNode(buffer.Bytes())
		require.NoError(t, err)
		temps = append(temps, temp)
	}
	after := getMemStats()
	t.Logf("after node MakeNode:gc<disable>, pool<disable>: %dMB,GC:%d", int64(after.Alloc/1024/1024)-int64(before.Alloc/1024/1024), after.NumGC-before.NumGC)

	// reset
	for i := 0; i < dataSize; i++ {
		SetNodeToPool(data[i])
	}

	before = getMemStats()
	temps1 := make([]*Node, 0)
	for i := 0; i < dataSize; i++ {
		temp, err := MakeNode(buffer.Bytes())
		require.NoError(t, err)
		temps1 = append(temps1, temp)
	}
	after = getMemStats()
	t.Logf("after node MakeNode:gc<disable>, pool<enable> : %dMB,GC:%d", int64(after.Alloc/1024/1024)-int64(before.Alloc/1024/1024), after.NumGC-before.NumGC)
}

func TestNewNode_GC_Compare_MemStats(t *testing.T) {
	h := []byte{1, 2, 3}
	node := &Node{key: []byte("test"), value: nBytes(1024 * 1024), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
	buffer := new(bytes.Buffer)
	err := node.writeBytes(buffer)
	require.NoError(t, err)
	dataSize := 1024 * 1024
	data := make([]*Node, 0)

	for i := 0; i < dataSize; i++ {
		temp := &Node{key: []byte(fmt.Sprintf("innerNode%d", i)), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
		data = append(data, temp)
	}

	debug.SetGCPercent(-1)

	before := getMemStats()
	temps := make([]*Node, 0)
	for i := 0; i < dataSize; i++ {
		temp := NewNode(node.key, node.value, node.version)
		temps = append(temps, temp)
	}
	after := getMemStats()
	t.Logf("after node NewNode:gc<disable>, pool<disable>: %dMB,GC:%d", int64(after.Alloc/1024/1024)-int64(before.Alloc/1024/1024), after.NumGC-before.NumGC)

	// reset
	for i := 0; i < dataSize; i++ {
		SetNodeToPool(data[i])
	}

	before = getMemStats()
	temps1 := make([]*Node, 0)
	for i := 0; i < dataSize; i++ {
		temp := NewNode(node.key, node.value, node.version)
		temps1 = append(temps1, temp)
	}
	after = getMemStats()
	t.Logf("after node NewNode:gc<disable>, pool<enable> : %dMB,GC:%d", int64(after.Alloc/1024/1024)-int64(before.Alloc/1024/1024), after.NumGC-before.NumGC)
}

func TestNodeJsonToNode_GC_Compare_MemStats(t *testing.T) {
	h := []byte{1, 2, 3}
	node := &Node{key: []byte("test"), value: nBytes(1024 * 1024), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
	nodeJson := NodeToNodeJson(node)
	buffer := new(bytes.Buffer)
	err := node.writeBytes(buffer)
	require.NoError(t, err)
	dataSize := 1024 * 1024
	data := make([]*Node, 0)

	for i := 0; i < dataSize; i++ {
		temp := &Node{key: []byte(fmt.Sprintf("innerNode%d", i)), version: 1, size: 1, height: 1, leftHash: h, rightHash: h}
		data = append(data, temp)
	}

	debug.SetGCPercent(-1)

	before := getMemStats()
	temps := make([]*Node, 0)
	for i := 0; i < dataSize; i++ {
		temp := NodeJsonToNode(nodeJson)
		temps = append(temps, temp)
	}
	after := getMemStats()
	t.Logf("after node NodeJsonToNode:gc<disable>, pool<disable>: %dMB,GC:%d", int64(after.Alloc/1024/1024)-int64(before.Alloc/1024/1024), after.NumGC-before.NumGC)

	// reset
	for i := 0; i < dataSize; i++ {
		SetNodeToPool(data[i])
	}

	before = getMemStats()
	temps1 := make([]*Node, 0)
	for i := 0; i < dataSize; i++ {
		temp := NodeJsonToNode(nodeJson)
		temps1 = append(temps1, temp)
	}
	after = getMemStats()
	t.Logf("after node NodeJsonToNode:gc<disable>, pool<enable> : %dMB,GC:%d", int64(after.Alloc/1024/1024)-int64(before.Alloc/1024/1024), after.NumGC-before.NumGC)
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
		benchmarkTreeSet(b, testCase.name, testCase.enableNodePool, testCase.initTreeHeight, 100)
	}
}

func benchmarkTreeSet(b *testing.B, name string, enableNodePool bool, treeHeight int8, kvSize int) {
	b.Run(name, func(b *testing.B) {
		_, _, tree, err := mockGCSepcialHeightTree(treeHeight)
		require.NoError(b, err)
		EnableNodePool = enableNodePool
		data := mockKVData(kvSize)
		for i := 0; i < kvSize*100; i++ {
			SetNodeToPool(&Node{key: []byte(fmt.Sprintf("innerNode%d", i)), version: 1, size: 1, height: 1, leftHash: []byte("test"), rightHash: []byte("test1")})
		}
		oldGet, oldSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
		b.ResetTimer()
		multiTreeSet(b, tree, data)
		b.ReportAllocs()
		newGet, newSet := GetNodeFromPoolCounter, SetNodeFromPoolCounter
		b.Logf("GetNodeFromPool:%d - SetNodeToPool:%d", newGet-oldGet, newSet-oldSet)
	})
}
func multiTreeSet(b *testing.B, tree *MutableTree, data []testKVData) {
	before := getMemStats()
	defer func() {
		after := getMemStats()
		b.Logf("after tree set:gc<disable>, pool<%v> : %dKB,GC:%d", EnableNodePool, int64(after.Alloc/1024)-int64(before.Alloc/1024), after.NumGC-before.NumGC)
	}()
	for i, _ := range data {
		tree.Set(data[i].key, data[i].value)
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
