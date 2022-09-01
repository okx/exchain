package watcher

import (
	db "github.com/okex/exchain/libs/tm-db"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

var td = []*dataT{
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
}

var results = []*resultT{
	{true, &MsgCodeByHash{Key: []byte("hello"), Code: "0x01"}},
	{true, &MsgCodeByHash{Key: []byte("hello1"), Code: "0x05"}},
	{true, &MsgCodeByHash{Key: []byte("hello2"), Code: "0x03"}},
	{true, &MsgCodeByHash{Key: []byte("hello3"), Code: "05"}},
	{true, &MsgCodeByHash{Key: []byte("hello4"), Code: "005"}},
}

type dataT struct {
	version int64
	msgs    []WatchMessage
}

type resultT struct {
	ok  bool
	msg *MsgCodeByHash
}

func TestACProcessorMoveToCommitList(t *testing.T) {
	testcases := []struct {
		processor *ACProcessor
		results   []*resultT
		fnInit    func(processor *ACProcessor, datas []*dataT)
		fnCheck   func(processor *ACProcessor, results []*resultT)
	}{
		{
			results: results,
			fnInit: func(processor *ACProcessor, datas []*dataT) {
				for _, d := range datas {
					processor.BatchSet(d.msgs)
					processor.MoveToCommitList(d.version)
				}
			},
			fnCheck: func(processor *ACProcessor, results []*resultT) {
				for _, r := range results {
					res, ok := processor.Get(r.msg.GetKey())
					require.Equal(t, r.ok, ok)
					require.Equal(t, r.msg.GetValue(), res.GetValue())
				}
			},
		},
		{
			results: results,
			fnInit: func(processor *ACProcessor, datas []*dataT) {
				for i, d := range datas {
					processor.BatchSet(d.msgs)
					if i < len(datas)-1 {
						processor.MoveToCommitList(d.version)
					}
				}
			},
			fnCheck: func(processor *ACProcessor, results []*resultT) {
				wg := &sync.WaitGroup{}
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func(wg *sync.WaitGroup) {
						defer wg.Done()
						for _, r := range results {
							res, ok := processor.Get(r.msg.GetKey())
							require.Equal(t, r.ok, ok)
							require.Equal(t, r.msg.GetValue(), res.GetValue())
						}
					}(wg)
					wg.Wait()
				}
			},
		},
	}

	for _, ts := range testcases {
		processor := &ACProcessor{
			commitList:  newCommitCache(), // for support to querier
			curMsgCache: newMessageCache(),
		}
		ts.fnInit(processor, td)
		ts.fnCheck(processor, ts.results)
	}
}

func TestACProcessorPersistHander(t *testing.T) {
	type testcase struct {
		tdb       *testdb
		processor *ACProcessor
		results   []*resultT
		fnInit    func(ts *testcase)
		fnCheck   func(ts *testcase)
	}

	testcases := []*testcase{
		{
			tdb:       &testdb{db.NewMemDB()},
			processor: &ACProcessor{commitList: newCommitCache(), curMsgCache: newMessageCache()},
			results:   results,
			fnInit: func(ts *testcase) {
				for _, d := range td {
					ts.processor.BatchSet(d.msgs)
					ts.processor.MoveToCommitList(d.version)
					ts.processor.PersistHander(ts.tdb.commit)
				}
			},
			fnCheck: func(ts *testcase) {
				for _, r := range ts.results {
					res, err := ts.tdb.Get(r.msg.GetKey())
					require.Equal(t, r.ok, err == nil)
					require.Equal(t, []byte(r.msg.GetValue()), res)
				}
			},
		},
		{
			tdb:       &testdb{db.NewMemDB()},
			processor: &ACProcessor{commitList: newCommitCache(), curMsgCache: newMessageCache()},
			results:   results,
			fnInit: func(ts *testcase) { // the last version in the curMsgCache and other version in memdb
				for i, d := range td {
					ts.processor.BatchSet(d.msgs)
					if i < len(td)-1 {
						ts.processor.MoveToCommitList(d.version)
					}
				}
				ts.processor.PersistHander(ts.tdb.commit)
			},
			fnCheck: func(ts *testcase) {
				for _, r := range ts.results {
					res, ok := ts.processor.Get(r.msg.GetKey())
					if !ok {
						res, err := ts.tdb.Get(r.msg.GetKey())
						require.Equal(t, r.ok, err == nil)
						require.Equal(t, []byte(r.msg.GetValue()), res)
						continue
					}
					require.Equal(t, r.ok, ok)
					require.Equal(t, r.msg.GetValue(), res.GetValue())
				}
			},
		},
	}

	for _, ts := range testcases {
		ts.fnInit(ts)
		ts.fnCheck(ts)
	}
}

type testdb struct {
	db.DB
}

func (tdb *testdb) commit(epochCache *MessageCache) {
	dbBatch := tdb.NewBatch()
	defer dbBatch.Close()
	for key, b := range epochCache.mp {
		if b == nil {
			dbBatch.Delete([]byte(key))
			continue
		}
		key := b.GetKey()
		value := []byte(b.GetValue())
		typeValue := b.GetType()
		if typeValue == TypeDelete {
			dbBatch.Delete(key)
		} else {
			dbBatch.Set(key, value)
		}
	}
	dbBatch.Write()
}
