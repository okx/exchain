package mempool

import (
	"fmt"
	"math/big"
	"sync"
	"testing"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/stretchr/testify/suite"
)

type PendingPoolTestSuite struct {
	suite.Suite

	Pool *PendingPool
}

func (suite *PendingPoolTestSuite) SetupTest() {
	suite.Pool = newPendingPool(100, 3, 100, 10)
}
func (suite *PendingPoolTestSuite) TestAddtx(t *testing.T) {
	testCases := []struct {
		Tx *mempoolTx
	}{
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("1"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(3780)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("2"), from: "6", realTx: abci.MockTx{GasPrice: big.NewInt(5853)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("3"), from: "7", realTx: abci.MockTx{GasPrice: big.NewInt(8315)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("4"), from: "10", realTx: abci.MockTx{GasPrice: big.NewInt(9526)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("5"), from: "15", realTx: abci.MockTx{GasPrice: big.NewInt(9140)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("6"), from: "9", realTx: abci.MockTx{GasPrice: big.NewInt(9227)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("7"), from: "3", realTx: abci.MockTx{GasPrice: big.NewInt(761)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("8"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(9740)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("9"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(6574)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("10"), from: "8", realTx: abci.MockTx{GasPrice: big.NewInt(9656)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("11"), from: "12", realTx: abci.MockTx{GasPrice: big.NewInt(6554)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("12"), from: "16", realTx: abci.MockTx{GasPrice: big.NewInt(5609)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("13"), from: "6", realTx: abci.MockTx{GasPrice: big.NewInt(2791), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("14"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(2698), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("15"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(6925), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("16"), from: "3", realTx: abci.MockTx{GasPrice: big.NewInt(3171)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("17"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(2965), Nonce: 2}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("18"), from: "19", realTx: abci.MockTx{GasPrice: big.NewInt(2484)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("19"), from: "13", realTx: abci.MockTx{GasPrice: big.NewInt(9722)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("20"), from: "7", realTx: abci.MockTx{GasPrice: big.NewInt(4236), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("21"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(1780)}}},
	}

	for _, exInfo := range testCases {
		suite.Pool.addTx(exInfo.Tx)
	}
	suite.Require().Equal(t, 18, suite.Pool.Size(), fmt.Sprintf("Expected to txs length %v but got %v", 18,
		suite.Pool.Size()))
	for _, exInfo := range testCases {
		t := suite.Pool.getTx(exInfo.Tx.from, exInfo.Tx.senderNonce)
		suite.Require().NotNil(t)
	}
}

func (suite *PendingPoolTestSuite) TestAddtxConcurrency(t *testing.T) {
	type Case struct {
		Tx *mempoolTx
	}

	testCases := []Case{
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("1"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(3780)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("2"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(3780), Nonce: 1}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("3"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(5315), Nonce: 2}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("4"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(4526), Nonce: 3}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("5"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(2140), Nonce: 4}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("6"), from: "1", realTx: abci.MockTx{GasPrice: big.NewInt(4227), Nonce: 5}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("7"), from: "2", realTx: abci.MockTx{GasPrice: big.NewInt(2161)}}},
	}

	var wait sync.WaitGroup
	for _, exInfo := range testCases {
		wait.Add(1)
		go func(p Case) {
			suite.Pool.addTx(p.Tx)
			wait.Done()
		}(exInfo)
	}

	wait.Wait()
}

func (suite *PendingPoolTestSuite) TestRemovetx(t *testing.T) {
	testCases := []struct {
		Tx *mempoolTx
	}{
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("1"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(3780)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("2"), from: "6", realTx: abci.MockTx{GasPrice: big.NewInt(5853)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("3"), from: "7", realTx: abci.MockTx{GasPrice: big.NewInt(8315)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("4"), from: "10", realTx: abci.MockTx{GasPrice: big.NewInt(9526)}}},
	}

	for _, exInfo := range testCases {
		suite.Pool.addTx(exInfo.Tx)
	}

	for _, exInfo := range testCases {
		t := suite.Pool.getTx(exInfo.Tx.from, exInfo.Tx.senderNonce)
		suite.Require().NotNil(t)
	}
	suite.Pool.removeTx("18", 0)
	res := suite.Pool.getTx("18", 0)
	suite.Require().Nil(res)

	suite.Pool.removeTx("6", 0)
	res = suite.Pool.getTx("6", 0)
	suite.Require().Nil(res)

	suite.Pool.removeTx("7", 0)
	res = suite.Pool.getTx("7", 0)
	suite.Require().Nil(res)

	suite.Pool.removeTx("10", 0)
	res = suite.Pool.getTx("10", 0)
	suite.Require().Nil(res)
}

func (suite *PendingPoolTestSuite) TestRemoveTxByHash(t *testing.T) {
	testCases := []struct {
		Tx *mempoolTx
	}{
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("1"), from: "18", realTx: abci.MockTx{GasPrice: big.NewInt(3780)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("2"), from: "6", realTx: abci.MockTx{GasPrice: big.NewInt(5853)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("3"), from: "7", realTx: abci.MockTx{GasPrice: big.NewInt(8315)}}},
		{&mempoolTx{height: 1, gasWanted: 1, tx: []byte("4"), from: "10", realTx: abci.MockTx{GasPrice: big.NewInt(9526)}}},
	}

	for _, exInfo := range testCases {
		suite.Pool.addTx(exInfo.Tx)
	}

	for _, exInfo := range testCases {
		t := suite.Pool.getTx(exInfo.Tx.from, exInfo.Tx.senderNonce)
		suite.Require().NotNil(t)
	}
	suite.Pool.removeTxByHash("0x000")
	res := suite.Pool.getTx("18", 0)
	suite.Require().NotNil(res)
}

func (suite *PendingPoolTestSuite) TestHandlePendingTx(t *testing.T) {

}

func (suite *PendingPoolTestSuite) TestHandlePeriodCounter(t *testing.T) {

}
