package watcher

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/okex/exchain/x/evm/types"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestMessageCache(t *testing.T) {
	type wrapData struct {
		version int64
		msgCode []WatchMessage
		delAcc  *DelAccMsg
		delKeys [][]byte
		batchs  []*Batch // add keys
	}
	testcases := []struct {
		*wrapData
		fncheck func(cache *MessageCache, data *wrapData)
	}{
		{
			wrapData: &wrapData{
				version: 1,
				msgCode: []WatchMessage{&MsgCodeByHash{Key: []byte("0x01"), Code: "0x01"}, &MsgCodeByHash{Key: []byte("0x02"), Code: "0x02"}},
				delAcc:  &DelAccMsg{addr: []byte("0x01")},
				delKeys: [][]byte{[]byte("0x04"), []byte("0x05"), []byte("0x06")},
				batchs:  []*Batch{{Key: []byte("0x01"), Value: []byte("0x01"), TypeValue: 1}, {Key: []byte("0x02"), Value: []byte("0x02"), TypeValue: 2}},
			},
			fncheck: func(cache *MessageCache, data *wrapData) {
				for _, msg := range data.msgCode {
					res, ok := cache.Get(msg.GetKey())
					require.True(t, ok)
					require.Equal(t, res.GetValue(), msg.GetValue())
				}
				res, ok := cache.Get(data.delAcc.GetKey())
				require.True(t, ok)
				require.Equal(t, res.GetValue(), data.delAcc.GetValue())
				require.Equal(t, res.GetType(), data.delAcc.GetType())
				for _, k := range data.delKeys {
					res, ok := cache.Get(k)
					require.True(t, ok)
					require.Equal(t, TypeDelete, res.GetType())
				}
				for _, bs := range data.batchs {
					res, ok := cache.Get(bs.GetKey())
					require.True(t, ok)
					require.Equal(t, res.GetValue(), bs.GetValue())
					require.Equal(t, res.GetType(), bs.GetType())
				}
			},
		},
		{
			wrapData: &wrapData{
				version: 2,
				msgCode: []WatchMessage{&MsgCodeByHash{Key: []byte("0x01"), Code: "0x02"}, &MsgCodeByHash{Key: []byte("0x02"), Code: "0x04"}},
				batchs:  []*Batch{{Key: []byte("0x04"), Value: []byte("0x01"), TypeValue: 1}, {Key: []byte("0x05"), Value: []byte("0x02"), TypeValue: 2}},
			},
			fncheck: func(cache *MessageCache, data *wrapData) {
				for _, msg := range data.msgCode {
					res, ok := cache.Get(msg.GetKey())
					require.True(t, ok)
					require.Equal(t, res.GetValue(), msg.GetValue())
				}

				//pre version delkey check
				res, ok := cache.Get([]byte("0x06"))
				require.True(t, ok)
				require.Equal(t, TypeDelete, res.GetType())

				for _, bs := range data.batchs {
					res, ok := cache.Get(bs.GetKey())
					require.True(t, ok)
					require.Equal(t, res.GetValue(), bs.GetValue())
					require.Equal(t, res.GetType(), bs.GetType())
				}
			},
		},
	}

	msgCache := newMessageCache()
	for _, ts := range testcases {
		if ts.delAcc != nil {
			msgCache.Set(ts.delAcc)
		}
		msgCache.BatchSet(ts.msgCode)
		msgCache.BatchDel(ts.delKeys)
		msgCache.BatchSetEx(ts.batchs)
		ts.fncheck(msgCache, ts.wrapData)
	}
}

func TestCommitCachePushBackAndRemove(t *testing.T) {
	cc := newCommitCache()
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(j int) {
			cc.pushBack(int64(j), &MessageCacheEvent{version: int64(j)})
			wg.Done()
		}(i)
	}
	wg.Wait()
	require.Equal(t, 10, cc.size())

	for i := 0; i < 10; i++ {
		msge := cc.remove(int64(i))
		require.NotNil(t, msge)
		require.Equal(t, msge.version, int64(i))
	}
	require.Equal(t, 0, cc.size())
}

func TestCommitCacheGetTop(t *testing.T) {
	cc := newCommitCache()
	for i := 0; i < 10; i++ {
		cc.pushBack(int64(i), &MessageCacheEvent{version: int64(i)})
	}

	// check
	count := 0
	for {
		cmmiter, ok := cc.getTop()
		if !ok {
			break
		}
		require.Equal(t, cmmiter.version, int64(count))
		cc.remove(cmmiter.version)
		count++
	}

	require.Equal(t, count, 10)
}

func TestCommitCacheGetElementFromCache(t *testing.T) {
	type dataT struct {
		version int64
		msgs    []WatchMessage
	}

	testcases := []struct {
		td      []*dataT
		cc      *commitCache
		fnInit  func(cc *commitCache, datas []*dataT)
		fncheck func(cc *commitCache)
	}{
		{
			td: []*dataT{
				{
					version: 1,
					msgs: []WatchMessage{
						&MsgCodeByHash{Key: []byte("hello"), Code: "0x01"},
						&MsgCodeByHash{Key: []byte("hello1"), Code: "01"},
					},
				},
				{
					version: 3,
					msgs: []WatchMessage{
						&MsgCodeByHash{Key: []byte("hello2"), Code: "0x03"},
						&MsgCodeByHash{Key: []byte("hello3"), Code: "03"},
					},
				},
				{
					version: 5,
					msgs: []WatchMessage{
						&MsgCodeByHash{Key: []byte("hello1"), Code: "0x05"},
						&MsgCodeByHash{Key: []byte("hello3"), Code: "05"},
						&MsgCodeByHash{Key: []byte("hello4"), Code: "005"},
					},
				},
			},
			cc: newCommitCache(),
			fnInit: func(cc *commitCache, datas []*dataT) {
				for _, d := range datas {
					c := newMessageCache()
					c.BatchSet(d.msgs)
					cc.pushBack(d.version, &MessageCacheEvent{MessageCache: c, version: d.version})
				}
			},
			fncheck: func(cc *commitCache) {
				results := []struct {
					ok  bool
					msg *MsgCodeByHash
				}{
					{true, &MsgCodeByHash{Key: []byte("hello"), Code: "0x01"}},
					{true, &MsgCodeByHash{Key: []byte("hello1"), Code: "0x05"}},
					{true, &MsgCodeByHash{Key: []byte("hello2"), Code: "0x03"}},
					{true, &MsgCodeByHash{Key: []byte("hello3"), Code: "05"}},
					{true, &MsgCodeByHash{Key: []byte("hello4"), Code: "005"}},
				}

				for _, re := range results {
					r, ok := cc.getElementFromCache(re.msg.GetKey())
					require.Equal(t, re.ok, ok)
					require.Equal(t, re.msg.GetValue(), r.GetValue())
				}
				r, ok := cc.getElementFromCache([]byte("hello5"))
				require.False(t, ok)
				require.Nil(t, r)
			},
		},
	}

	for _, ts := range testcases {
		ts.fnInit(ts.cc, ts.td)
		ts.fncheck(ts.cc)
	}
}

func TestStatistic(t *testing.T) {
	mp := map[string]WatchMessage{
		hex.EncodeToString(append(prefixTx, []byte("hell1")...)):    &Batch{Key: append(prefixTx, []byte("hell1")...), Value: []byte("222")},
		hex.EncodeToString(append(prefixTx, []byte("hell2")...)):    &Batch{Key: append(prefixTx, []byte("hell2")...), Value: []byte("222")},
		hex.EncodeToString(append(prefixTx, []byte("hell3")...)):    &Batch{Key: append(prefixTx, []byte("hell3")...), Value: []byte("222")},
		hex.EncodeToString(append(prefixBlock, []byte("hell3")...)): &Batch{Key: append(prefixBlock, []byte("hell3")...), Value: []byte("222")},
		hex.EncodeToString(append(prefixBlock, []byte("hell4")...)): &Batch{Key: append(prefixBlock, []byte("hell4")...), Value: []byte("222")},
	}
	static := make(map[string]*Stat)
	for _, v := range mp {
		Statistic(v, static)
	}
	for k, v := range static {
		dbsize := float64(v.dbSize) / float64(1024*1024)
		structsize := float64(v.structSize) / float64(1024*1024)
		fmt.Printf("**** lyh ****** static %s, count %d, dbSize %.3f, structSize %.3f \n", k, v.count, dbsize, structsize)
	}
}

func TestTransactionReceiptSize(t *testing.T) {
	tx1 := &TransactionReceipt{
		TransactionHash: "11111111111111111111111111111111111111111",
		tx: &types.MsgEthereumTx{
			Data: types.TxData{Payload: []byte("222222222222222222222")},
		},
	}
	tx2 := &TransactionReceipt{
		TransactionHash: "11111111111111111111111111111111111111111",
		tx: &types.MsgEthereumTx{
			Data: types.TxData{Payload: []byte("222222222222222222222")},
		},
	}
	fmt.Println(getSize(tx1), getSize(tx2))

	//prefixAccount 21
	//PrefixState 53
	//prefixTx 33
	//prefixReceipt 33

	//**** lyh ****** static prefixReceipt, count 20000, dbSize 21.702, structSize 42.503
	//**** lyh ****** static prefixAccount, count 3473, dbSize 0.263, structSize 0.584
	//**** lyh ****** static prefixLatestHeight, count 1, dbSize 0.000, structSize 0.000
	//**** lyh ****** static prefixBlockInfo, count 1, dbSize 0.000, structSize 0.000
	//**** lyh ****** static prefixRpcDb, count 2559, dbSize 0.054, structSize 0.200
	//**** lyh ****** static prefixTx, count 20000, dbSize 8.052, structSize 22.866
	//**** lyh ****** static PrefixState, count 23084, dbSize 1.871, structSize 3.082
	//**** lyh ****** static prefixBlock, count 1, dbSize 1.317, structSize 1.317
	//**** lyh ****** static prefixParams, count 1, dbSize 0.000, structSize 0.000

	num := 10
	txcount := 20000 * num
	stcount := 23000 * num
	recount := 20000 * num
	accCount := 3500 * num
	rpcCount := 2500 * num
	n := newMessageCache()

	var batchs []*Batch
	for i := 0; i < (txcount + recount); i++ {
		batchs = append(batchs, &Batch{
			Key:       randBytes(33),
			Value:     nil,
			TypeValue: 0,
		})
	}

	for i := 0; i < stcount; i++ {
		batchs = append(batchs, &Batch{
			Key:       randBytes(53),
			Value:     nil,
			TypeValue: 0,
		})
	}

	for i := 0; i < accCount; i++ {
		batchs = append(batchs, &Batch{
			Key:       randBytes(21),
			Value:     nil,
			TypeValue: 0,
		})
	}

	for i := 0; i < rpcCount; i++ {
		batchs = append(batchs, &Batch{
			Key:       randBytes(53),
			Value:     nil,
			TypeValue: 0,
		})
	}
	printmem()
	n.BatchSetEx(batchs)
	printmem()

	// ******lyh****** Alloc 94.79 MB
	// ******lyh****** Alloc 229.93 MB
}

func randBytes(numBytes int) []byte {
	b := make([]byte, numBytes)
	_, _ = rand.Read(b)
	return b
}
