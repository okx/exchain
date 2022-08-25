package watcher

import (
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
