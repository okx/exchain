package watcher_test

import (
	"github.com/ethereum/go-ethereum/common"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBanchAsyncCommitWatcher(t *testing.T) {
	viper.Set(watcher.FlagWatchdbEnableAsyncCommit, true)
	watcher.SetCommitGapHeight(1)
	defer func() {
		watcher.SetAnableAsyncCommit(false)
		watcher.SetCommitGapHeight(100)
	}()
	TestDeployAndCallContract(t)
	TestDuplicateAddress(t)
	TestMsgEthereumTxByWatcher(t)
	TestHandleMsgEthereumTx(t)
}

func TestClassifyWatchMessageType(t *testing.T) {
	batch := []watcher.WatchMessage{
		&watcher.MsgCodeByHash{Key: []byte("0x01"), Code: "0x01"},
		&watcher.MsgCodeByHash{Key: []byte("0x02"), Code: "0x02"},
		&watcher.MsgCodeByHash{Key: []byte("hello1"), Code: "0x05"},
		watcher.NewMsgBlock(watcher.Block{}),
		&watcher.MsgEthTx{Key: []byte("11")},
		watcher.NewMsgTransactionReceipt(watcher.TransactionReceipt{}, common.Hash{}),
	}
	acbatch, noacbacth := watcher.ClassifyWatchMessageType(batch)
	require.Equal(t, 3, len(acbatch))
	require.Equal(t, 3, len(noacbacth))
	for i, bc := range batch {
		if i < 3 {
			require.Equal(t, bc.GetKey(), acbatch[i].GetKey())
		} else {
			require.Equal(t, bc.GetKey(), noacbacth[i-3].GetKey())
		}
	}
}

func TestClassifyWatchDataType(t *testing.T) {
	addr1 := sdk.AccAddress("0x01")
	addr2 := sdk.AccAddress("0x02")
	wd := watcher.WatchData{
		DirtyAccount: []*sdk.AccAddress{&addr1, &addr2},
		Batches: []*watcher.Batch{
			{Key: watcher.MsgCodeByHash{Key: []byte("0x01"), Code: "0x01"}.GetKey()},
			{Key: watcher.NewMsgBlock(watcher.Block{}).GetKey()},
		},
		DelayEraseKey: [][]byte{
			watcher.MsgEthTx{Key: []byte("11")}.GetKey(),
			[]byte("0x02")},
		BloomData: []*types.KV{{Key: []byte("0x01")}, {Value: []byte("0x01")}},
		DirtyList: [][]byte{
			watcher.NewMsgTransactionReceipt(watcher.TransactionReceipt{}, common.Hash{}).GetKey(),
			[]byte("0x02")},
	}
	acBatch, noACBatch := watcher.ClassifyWatchDataType(wd)
	require.Equal(t, 2, len(acBatch.DirtyAccount))
	require.Equal(t, wd.Batches[0].GetKey(), acBatch.Batches[0].GetKey())
	require.Equal(t, wd.DelayEraseKey[1], acBatch.DelayEraseKey[0])
	require.Nil(t, acBatch.BloomData)
	require.Equal(t, wd.DirtyList[1], acBatch.DirtyList[0])

	require.Nil(t, noACBatch.DirtyAccount)
	require.Equal(t, wd.Batches[1].GetKey(), noACBatch.Batches[0].GetKey())
	require.Equal(t, wd.DelayEraseKey[0], noACBatch.DelayEraseKey[0])
	require.Equal(t, 2, len(noACBatch.BloomData))
	require.Equal(t, wd.DirtyList[0], noACBatch.DirtyList[0])

}
