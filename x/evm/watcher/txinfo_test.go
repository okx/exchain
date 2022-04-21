package watcher

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
)

func (suite *WatcherTestSuite) TestWatcher_addTxsToBlock() {
	const txsCount = 10
	for i := txsCount; i > 0; i-- {
		suite.watcher.txInfoCollector = append(suite.watcher.txInfoCollector, &TxInfo{TxHash: ethcmn.Hash{byte(i)}, Index: uint64(i)})
	}
	suite.watcher.sortTxsAndUpdateCumulativeGas()
	suite.Equal(ethcmn.Hash{byte(1)}, suite.watcher.blockTxs[0])
	suite.Equal(ethcmn.Hash{byte(txsCount)}, suite.watcher.blockTxs[txsCount-1])
}
