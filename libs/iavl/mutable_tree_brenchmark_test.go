package iavl

import (
	"fmt"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
	"math/rand"
	"os"
	"testing"
	"time"
)


var dbDir = "testdata"
func prepareTree(b *testing.B, openLogFlag bool, dbName string, size int) (*MutableTree, []string, map[string]string) {
	moduleName := "test"
	dir := dbDir
	if openLogFlag {
		openLog(moduleName)
	}
	ldb, err := db.NewGoLevelDB(dbName, dir)
	memDB := db.NewPrefixDB(ldb, []byte(moduleName))
	tree, err := NewMutableTree(memDB, 0)
	require.NoError(b, err)

	fmt.Printf("init setting test %d data to MutableTree\n", size)
	dataSet := make(map[string]string)
	keySet := make([]string, 0, size)
	for i:=0;i<size;i++ {
		key := randstr(32)
		value := randstr(100)
		dataSet[key] = string(value)
		keySet = append(keySet, key)
		tree.Set([]byte(key), []byte(value))
	}
	//recursivePrint(tree.root, 0)

	tree.SaveVersion(false)
	tree.commitCh <- commitEvent{-1, nil,nil, nil, nil, 0}
	fmt.Println("init setting done")
	return tree, keySet, dataSet
}

func benchmarkTreeRead(b *testing.B, tree *MutableTree, keySet []string, readNum int) {
	fmt.Println("benchmark testing")
	t1 := time.Now()
	for i:=0;i<readNum;i++ {
		idx := rand.Int()%len(keySet)
		key := keySet[idx]
		_, v :=tree.Get([]byte(key))
		require.NotNil(b, v)
	}
	duration := time.Since(t1)
	fmt.Println("time:", duration.String())
}

func clearDB(dbName string) {
	path := dbDir + "/" + dbName + ".db"
	fmt.Println("clear db", path)
	err := os.RemoveAll(path)
	if err != nil {
		fmt.Println(err)
	}
	treeMap.mutableTreeList = nil
	treeMap.mutableTreeSavedMap = make(map[string]bool)
}


func BenchmarkMutableTree_Get(b *testing.B) {
	EnableAsyncCommit = true
	EnablePruningHistoryState = true
	CommitIntervalHeight = 1
	defer func() {
		EnableAsyncCommit = false
		EnablePruningHistoryState = false
		CommitIntervalHeight = 100
	}()
	testCases := []struct {
		dbName string
		openLog bool
		initDataSize int
		readNum int
	}{
		{
			dbName: "13-test",
			openLog: true,
			initDataSize: 130000,
			readNum: 100000,
		},
		{
			dbName: "10-test",
			openLog: true,
			initDataSize: 100000,
			readNum: 100000,
		},
		{
			dbName: "8-test",
			openLog: true,
			initDataSize: 80000,
			readNum: 100000,
		},
		{
			dbName: "5-test",
			openLog: true,
			initDataSize: 50000,
			readNum: 100000,
		},
		{
			dbName: "3-test",
			openLog: true,
			initDataSize: 30000,
			readNum: 100000,
		},
	}
	for i, testCase := range testCases {
		fmt.Println("test case", i, ": ", testCase.dbName)
		tree, keySet, _ := prepareTree(b, testCase.openLog, testCase.dbName, testCase.initDataSize)
		benchmarkTreeRead(b, tree, keySet, testCase.readNum)
		clearDB(testCase.dbName)
		fmt.Println()
	}
}


func BenchmarkMutableTree_Get2(b *testing.B) {
	EnableAsyncCommit = true
	EnablePruningHistoryState = true
	CommitIntervalHeight = 1
	defer func() {
		EnableAsyncCommit = false
		EnablePruningHistoryState = false
		CommitIntervalHeight = 100
	}()
	testCases := []struct {
		dbName string
		openLog bool
		initDataSize int
		readNum int
	}{
		{
			dbName: "16-test",
			openLog: true,
			initDataSize: 16,
			readNum: 100000,
		},
		{
			dbName: "256-test",
			openLog: true,
			initDataSize: 256,
			readNum: 100000,
		},
		{
			dbName: "4096-test",
			openLog: true,
			initDataSize: 4096,
			readNum: 100000,
		},
		{
			dbName: "65536-test",
			openLog: true,
			initDataSize: 65536,
			readNum: 100000,
		},
	}
	for i, testCase := range testCases {
		fmt.Println("test case", i, ": ", testCase.dbName)
		tree, keySet, _ := prepareTree(b, testCase.openLog, testCase.dbName, testCase.initDataSize)
		benchmarkTreeRead(b, tree, keySet, testCase.readNum)
		clearDB(testCase.dbName)
		fmt.Println()
	}
}


func BenchmarkMutableTree_Get3(b *testing.B) {
	EnableAsyncCommit = true
	EnablePruningHistoryState = true
	CommitIntervalHeight = 1
	defer func() {
		EnableAsyncCommit = false
		EnablePruningHistoryState = false
		CommitIntervalHeight = 100
	}()
	testCases := []struct {
		dbName string
		openLog bool
		initDataSize int
		readNum int
	}{
		{
			dbName: "16-test",
			openLog: true,
			initDataSize: 16,
			readNum: 100000,
		},
	}
	for i, testCase := range testCases {
		fmt.Println("test case", i, ": ", testCase.dbName)
		tree, keySet, _ := prepareTree(b, testCase.openLog, testCase.dbName, testCase.initDataSize)
		benchmarkTreeRead(b, tree, keySet, testCase.readNum)
		clearDB(testCase.dbName)
		fmt.Println()
	}
}

func recursivePrint(node *Node, n int) {
	if node == nil {
		return
	}
	for i:=0;i<n;i++ {
		fmt.Printf("   ")
	}
	fmt.Printf("height:%d key:%s\n", node.height, string(node.key))
	recursivePrint(node.leftNode, n+1)
	recursivePrint(node.rightNode, n+1)
}