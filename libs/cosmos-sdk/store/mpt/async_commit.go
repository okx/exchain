package mpt

import (
	"container/list"
	"sync"
	"unsafe"

	"github.com/okex/exchain/libs/tendermint/libs/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)

type replayer interface {
	Replay(w ethdb.KeyValueWriter) error
	IsSingle() bool
}

type preCommitValue struct {
	value   []byte
	deleted bool
	ele     *list.Element
}

type singleOp struct {
	key    []byte
	value  []byte
	delete bool
}

func (op *singleOp) Replay(w ethdb.KeyValueWriter) error {
	if op.delete {
		return w.Delete(op.key)
	}
	return w.Put(op.key, op.value)
}

func (op *singleOp) IsSingle() bool {
	return true
}

type multiOp []singleOp

func (ops multiOp) Replay(w ethdb.KeyValueWriter) error {
	for _, op := range ops {
		if err := op.Replay(w); err != nil {
			return err
		}
	}
	return nil
}

func (ops multiOp) IsSingle() bool {
	return false
}

type actionOp struct {
	action func(w ethdb.KeyValueWriter)
}

func (op *actionOp) Replay(w ethdb.KeyValueWriter) error {
	if op.action != nil {
		op.action(w)
	}
	return nil
}

func (op *actionOp) IsSingle() bool {
	return true
}

type preCommitMap struct {
	data  map[string]preCommitValue
	store *AsyncKeyValueStore
}

func (w preCommitMap) Put(key []byte, value []byte) error {
	w.data[string(key)] = preCommitValue{
		value: value,
		ele:   w.store.preCommitTail,
	}
	return nil
}

func (w preCommitMap) Delete(key []byte) error {
	w.data[string(key)] = preCommitValue{
		deleted: true,
		ele:     w.store.preCommitTail,
	}
	return nil
}

type preCommitClearMap preCommitMap

type Element struct {
	next, prev *Element
	list       *list.List
}

func (w preCommitClearMap) Put(key []byte, _ []byte) error {
	if v, ok := w.data[string(key)]; ok {
		if elep := (*Element)((unsafe.Pointer)(v.ele)); elep.list == nil {
			delete(w.data, string(key))
		}
	}
	return nil
}

func (w preCommitClearMap) Delete(key []byte) error {
	if v, ok := w.data[string(key)]; ok {
		if elep := (*Element)((unsafe.Pointer)(v.ele)); elep.list == nil {
			delete(w.data, string(key))
		}
	}
	return nil
}

type commitTask struct {
	op replayer
}

type AsyncKeyValueStore struct {
	ethdb.KeyValueStore

	mtx           sync.RWMutex
	preCommit     preCommitMap
	preCommitList *list.List
	preCommitTail *list.Element
	listMtx       sync.Mutex

	enableAsyncCommit bool

	commitCh chan struct{}
	clearCh  chan []*commitTask
	closeWg  sync.WaitGroup

	logger log.Logger
}

func NewAsyncKeyValueStore(db ethdb.KeyValueStore) *AsyncKeyValueStore {
	store := &AsyncKeyValueStore{
		KeyValueStore: db,
		preCommit: preCommitMap{
			data: make(map[string]preCommitValue),
		},
		preCommitList: list.New(),
		commitCh:      make(chan struct{}, 10000*10),
		clearCh:       make(chan []*commitTask, 1024),
	}
	store.preCommit.store = store
	store.closeWg.Add(1)
	go store.commitRoutine()
	go store.clearRoutine()
	return store
}

func (store *AsyncKeyValueStore) SetLogger(logger log.Logger) {
	store.logger = logger
}

func (store *AsyncKeyValueStore) Has(key []byte) (bool, error) {
	store.mtx.RLock()
	defer store.mtx.RUnlock()

	if v, ok := store.preCommit.data[string(key)]; ok {
		return !v.deleted, nil
	}

	return store.KeyValueStore.Has(key)
}

func (store *AsyncKeyValueStore) Get(key []byte) ([]byte, error) {
	store.mtx.RLock()
	defer store.mtx.RUnlock()

	if v, ok := store.preCommit.data[string(key)]; ok {
		if v.deleted {
			return nil, nil
		}
		return v.value, nil
	}
	return store.KeyValueStore.Get(key)
}

func (store *AsyncKeyValueStore) Put(key []byte, value []byte) error {
	key, value = common.CopyBytes(key), common.CopyBytes(value)
	store.mtx.Lock()
	defer store.mtx.Unlock()

	store.listMtx.Lock()
	store.preCommitTail = store.preCommitList.PushBack(&commitTask{
		op: &singleOp{
			key:   key,
			value: value,
		},
	})
	store.listMtx.Unlock()

	_ = store.preCommit.Put(key, value)

	store.commitCh <- struct{}{}
	return nil
}

func (store *AsyncKeyValueStore) Delete(key []byte) error {
	key = common.CopyBytes(key)
	store.mtx.Lock()
	defer store.mtx.Unlock()

	store.listMtx.Lock()
	store.preCommitTail = store.preCommitList.PushBack(&commitTask{
		op: &singleOp{
			key:    key,
			delete: true,
		},
	})
	store.listMtx.Unlock()

	_ = store.preCommit.Delete(key)

	store.commitCh <- struct{}{}
	return nil
}

func (store *AsyncKeyValueStore) batchWrite(player replayer) error {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	store.listMtx.Lock()
	store.preCommitTail = store.preCommitList.PushBack(&commitTask{
		op: player,
	})
	store.listMtx.Unlock()

	if err := player.Replay(store.preCommit); err != nil {
		return err
	}

	store.commitCh <- struct{}{}
	return nil
}

func (store *AsyncKeyValueStore) LogInfoAfterWriteDone(msg string, args ...interface{}) {
	if store.logger == nil {
		return
	}

	store.ActionAfterWriteDone(func() {
		store.logger.Info(msg, args...)
	})
}

func (store *AsyncKeyValueStore) ActionAfterWriteDone(act func()) {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	store.listMtx.Lock()
	store.preCommitTail = store.preCommitList.PushBack(&commitTask{
		op: &actionOp{
			action: func(ethdb.KeyValueWriter) {
				act()
			},
		},
	})
	store.listMtx.Unlock()

	store.commitCh <- struct{}{}
}

func (store *AsyncKeyValueStore) NewBatch() ethdb.Batch {
	return newAsyncBatch(store)
}

func (store *AsyncKeyValueStore) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	// TODO
	return store.KeyValueStore.NewIterator(prefix, start)
}

//func (store *AsyncKeyValueStore) Stat(property string) (string, error) {
//	return store.KeyValueStore.Stat(property)
//}
//
//func (store *AsyncKeyValueStore) Compact(start []byte, limit []byte) error {
//	return store.KeyValueStore.Compact(start, limit)
//}

func (store *AsyncKeyValueStore) Close() error {
	store.mtx.Lock()
	defer store.mtx.Unlock()

	close(store.commitCh)
	store.closeWg.Wait()
	return store.KeyValueStore.Close()
}

func (store *AsyncKeyValueStore) commitRoutine() {
	defer func() {
		close(store.clearCh)
		store.closeWg.Done()
	}()

	doneList := make([]*commitTask, 0, 100)
	batchSize := 0

	for _ = range store.commitCh {
		store.listMtx.Lock()
		taskEle := store.preCommitList.Front()
		store.preCommitList.Remove(taskEle)
		store.listMtx.Unlock()

		task := taskEle.Value.(*commitTask)

		if task.op == nil {
			continue
		}
		var kvWriter ethdb.KeyValueWriter = store.KeyValueStore
		var batch ethdb.Batch

		if !task.op.IsSingle() {
			batch = store.KeyValueStore.NewBatch()
			kvWriter = batch
		}
		if err := task.op.Replay(kvWriter); err != nil {
			panic(err)
		}
		if batch != nil {
			if err := batch.Write(); err != nil {
				panic(err)
			}
			batchSize += batch.ValueSize()
		}

		doneList = append(doneList, task)
		if len(doneList) >= 100 || batchSize > 1_000_000 {
			store.clearCh <- doneList
			doneList = make([]*commitTask, 0, 100)
			batchSize = 0
		}
	}
}

func (store *AsyncKeyValueStore) clearRoutine() {
	for doneList := range store.clearCh {
		store.mtx.Lock()
		for _, task := range doneList {
			_ = task.op.Replay(preCommitClearMap(store.preCommit))
		}
		store.mtx.Unlock()
	}
}

type asyncBatch struct {
	ops   multiOp
	store *AsyncKeyValueStore
	size  int
}

func newAsyncBatch(store *AsyncKeyValueStore) *asyncBatch {
	return &asyncBatch{
		store: store,
	}
}

func (b *asyncBatch) Put(key []byte, value []byte) error {
	key, value = common.CopyBytes(key), common.CopyBytes(value)
	b.ops = append(b.ops, singleOp{
		key:   key,
		value: value,
	})
	b.size += len(value)
	return nil
}

func (b *asyncBatch) Delete(key []byte) error {
	key = common.CopyBytes(key)
	b.ops = append(b.ops, singleOp{
		key:    key,
		delete: true,
	})
	b.size += len(key)
	return nil
}

func (b *asyncBatch) ValueSize() int {
	return b.size
}

func (b *asyncBatch) Write() error {
	ops := b.ops
	b.ops = nil
	return b.store.batchWrite(ops)
}

func (b *asyncBatch) Reset() {
	b.ops = b.ops[:0]
	b.size = 0
}

func (b *asyncBatch) Replay(w ethdb.KeyValueWriter) error {
	return b.ops.Replay(w)
}

func (b *asyncBatch) IsSingle() bool {
	return false
}
