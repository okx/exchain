package watcher_test

import (
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
	"testing"
)

func TestBanchAsyncCommitWatcher(t *testing.T) {
	viper.Set(watcher.FlagWatchdbEnableAsyncCommit, true)
	viper.Set(watcher.FlagWatchdbCommitGapHeight, 1)
	defer func() {
		watcher.SetAnableAsyncCommit(false)
		watcher.SetCommitGapHeight(100)
	}()
	TestDeployAndCallContract(t)
	TestDuplicateAddress(t)
	TestMsgEthereumTxByWatcher(t)
	TestHandleMsgEthereumTx(t)
}
