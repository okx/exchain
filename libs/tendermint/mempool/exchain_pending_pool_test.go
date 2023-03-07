package mempool

import (
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"testing"
	"time"

	abci "github.com/okx/okbchain/libs/tendermint/abci/types"
	"github.com/stretchr/testify/assert"
)

func TestAddtx(t *testing.T) {
	pool := newPendingPool(100, 3, 10, 10)
	testCases := []struct {
		Tx *mempoolTx
	}{
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("1"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(3780)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("2"), from: "2", realTx: abci.MockTx{GasPrice: big.NewInt(5853)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("3"), from: "3", realTx: abci.MockTx{GasPrice: big.NewInt(8315)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("4"), from: "4", realTx: abci.MockTx{GasPrice: big.NewInt(9526)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("5"), from: "5", realTx: abci.MockTx{GasPrice: big.NewInt(9140)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("6"), from: "6", realTx: abci.MockTx{GasPrice: big.NewInt(9227)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("13"), from: "12", realTx: abci.MockTx{GasPrice: big.NewInt(2791), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("14"), from: "12", realTx: abci.MockTx{GasPrice: big.NewInt(2698), Nonce: 100}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("15"), from: "12", realTx: abci.MockTx{GasPrice: big.NewInt(2698), Nonce: 18446744073709551615}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("16"), from: "14", realTx: abci.MockTx{GasPrice: big.NewInt(6925), Nonce: 100}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("17"), from: "14", realTx: abci.MockTx{GasPrice: big.NewInt(2965), Nonce: 99}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("18"), from: "14", realTx: abci.MockTx{GasPrice: big.NewInt(2791), Nonce: 98}}},
	}

	for _, exInfo := range testCases {
		pool.addTx(exInfo.Tx)
	}
	assert.Equal(t, len(testCases), pool.Size(), fmt.Sprintf("Expected to txs length %v but got %v", len(testCases), pool.Size()))
	for _, exInfo := range testCases {
		tx := pool.addressTxsMap[exInfo.Tx.from][exInfo.Tx.realTx.GetNonce()]
		assert.Equal(t, tx, exInfo.Tx)
	}
}

func TestAddtxRandom(t *testing.T) {
	pool := newPendingPool(100000, 3, 10, 20000)
	txCount := 10000
	rand.Seed(time.Now().Unix())
	addrMap := map[int]string{
		0: "1234567",
		1: "0x333",
		2: "11111",
		3: "test",
	}

	for i := 0; i < txCount; i++ {
		nonce := rand.Intn(txCount)
		addrIndex := nonce % len(addrMap)
		tx := &mempoolTx{height: 1, gasWanted: 1, tx: []byte(strconv.Itoa(i)), from: addrMap[addrIndex], realTx: abci.MockTx{Nonce: uint64(nonce)}}
		pool.addTx(tx)
		txRes := pool.addressTxsMap[tx.from][tx.realTx.GetNonce()]
		assert.Equal(t, tx, txRes)
	}
	assert.Equal(t, txCount, pool.Size(), fmt.Sprintf("Expected to txs length %v but got %v", txCount, pool.Size()))
}

func TestRemovetx(t *testing.T) {
	pool := newPendingPool(100, 3, 10, 10)
	txs := []struct {
		Tx *mempoolTx
	}{
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("13"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(2791), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("14"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(2698), Nonce: 100}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("15"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(2698), Nonce: 18446744073709551615}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("16"), from: "14", realTx: abci.MockTx{GasPrice: big.NewInt(6925), Nonce: 100}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("17"), from: "14", realTx: abci.MockTx{GasPrice: big.NewInt(2965), Nonce: 99}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("18"), from: "14", realTx: abci.MockTx{GasPrice: big.NewInt(2791), Nonce: 98}}},
	}
	testCases := []struct {
		address string
		nonce   uint64
	}{
		{"18", 100},
		{"nonexist", 0},
		{"", 0},
		{"18", 1000},
		{"14", 98},
	}
	for _, exInfo := range txs {
		pool.addTx(exInfo.Tx)
		tx := pool.getTx(exInfo.Tx.from, exInfo.Tx.realTx.GetNonce())
		assert.Equal(t, tx, exInfo.Tx)
	}

	for _, tc := range testCases {
		pool.removeTx(tc.address, tc.nonce)
		res := pool.getTx(tc.address, tc.nonce)
		assert.Nil(t, res)
	}
}
func TestRemoveTxByHash(t *testing.T) {
	pool := newPendingPool(100, 3, 10, 10)
	txs := []struct {
		Tx *mempoolTx
	}{
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("13"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(2791), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("14"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(2698), Nonce: 100}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("15"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(2698), Nonce: 18446744073709551615}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("16"), from: "14", realTx: abci.MockTx{GasPrice: big.NewInt(6925), Nonce: 100}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("17"), from: "14", realTx: abci.MockTx{GasPrice: big.NewInt(2965), Nonce: 99}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("18"), from: "14", realTx: abci.MockTx{GasPrice: big.NewInt(2791), Nonce: 98}}},
	}

	for _, exInfo := range txs {
		pool.addTx(exInfo.Tx)
		tx := pool.getTx(exInfo.Tx.from, exInfo.Tx.realTx.GetNonce())
		assert.Equal(t, tx, exInfo.Tx)
	}

	for _, tc := range txs {
		pool.removeTxByHash(txID(tc.Tx.tx))
		res := pool.getTx(tc.Tx.from, tc.Tx.realTx.GetNonce())
		assert.Nil(t, res)
	}
}

func TestHandlePendingTx(t *testing.T) {
	pool := newPendingPool(100, 3, 10, 10)
	txs := []struct {
		Tx *mempoolTx
	}{
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("2"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(3780), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("3"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(5315), Nonce: 2}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("4"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(4526), Nonce: 3}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("5"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(2140), Nonce: 4}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("6"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(4227), Nonce: 5}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("7"), from: "2", realTx: abci.MockTx{GasPrice: big.NewInt(5315), Nonce: 2}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("8"), from: "2", realTx: abci.MockTx{GasPrice: big.NewInt(4526), Nonce: 3}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("9"), from: "2", realTx: abci.MockTx{GasPrice: big.NewInt(2140), Nonce: 5}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("10"), from: "2", realTx: abci.MockTx{GasPrice: big.NewInt(4227), Nonce: 6}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("11"), from: "3", realTx: abci.MockTx{GasPrice: big.NewInt(4526), Nonce: 3}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("12"), from: "3", realTx: abci.MockTx{GasPrice: big.NewInt(2140), Nonce: 4}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("13"), from: "3", realTx: abci.MockTx{GasPrice: big.NewInt(4227), Nonce: 5}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("14"), from: "4", realTx: abci.MockTx{GasPrice: big.NewInt(3780), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("15"), from: "4", realTx: abci.MockTx{GasPrice: big.NewInt(5315), Nonce: 2}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("16"), from: "4", realTx: abci.MockTx{GasPrice: big.NewInt(4526), Nonce: 3}}},
	}
	addressNonceTestCase := map[string]uint64{
		"1":         0,
		"2":         1,
		"3":         1,
		"4":         2,
		"non-exist": 0,
	}
	testCases := []struct {
		address       string
		nonceExpected uint64
	}{
		{"1", 1},
		{"2", 2},
		{"4", 3},
	}
	for _, exInfo := range txs {
		pool.addTx(exInfo.Tx)
	}
	assert.Equal(t, len(txs), pool.Size(), fmt.Sprintf("Expected to txs length %v but got %v", len(txs),
		pool.Size()))

	res := pool.handlePendingTx(addressNonceTestCase)
	for _, tc := range testCases {
		assert.Equal(t, tc.nonceExpected, res[tc.address], fmt.Sprintf("Expected tx nonce %v for  address %s, but got %d", tc.nonceExpected, tc.address,
			res[tc.address]))
	}
}

func TestHandlePeriodCounter(t *testing.T) {

	pool := newPendingPool(100, 3, 10, 10)
	txs := []struct {
		Tx *mempoolTx
	}{
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("1"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(3780), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("2"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(5315), Nonce: 2}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("3"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(4526), Nonce: 3}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("4"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(2140), Nonce: 4}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("5"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(4227), Nonce: 5}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("6"), from: "2", realTx: abci.MockTx{GasPrice: big.NewInt(5315), Nonce: 2}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("7"), from: "2", realTx: abci.MockTx{GasPrice: big.NewInt(4526), Nonce: 3}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("8"), from: "2", realTx: abci.MockTx{GasPrice: big.NewInt(2140), Nonce: 5}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("9"), from: "2", realTx: abci.MockTx{GasPrice: big.NewInt(4227), Nonce: 6}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("10"), from: "3", realTx: abci.MockTx{GasPrice: big.NewInt(4526), Nonce: 3}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("11"), from: "3", realTx: abci.MockTx{GasPrice: big.NewInt(2140), Nonce: 4}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("12"), from: "3", realTx: abci.MockTx{GasPrice: big.NewInt(4227), Nonce: 5}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("13"), from: "4", realTx: abci.MockTx{GasPrice: big.NewInt(3780), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("14"), from: "4", realTx: abci.MockTx{GasPrice: big.NewInt(5315), Nonce: 2}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("15"), from: "4", realTx: abci.MockTx{GasPrice: big.NewInt(4526), Nonce: 3}}},
	}
	for _, exInfo := range txs {
		pool.addTx(exInfo.Tx)
	}
	for i := 0; i < pool.reserveBlocks; i++ {
		pool.handlePeriodCounter()
	}
	assert.Equal(t, len(txs), pool.Size())
	pool.handlePeriodCounter()
	assert.Equal(t, 0, pool.Size())
}
