package mpt

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	"github.com/tendermint/go-amino"
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
	once   bool
}

func (op *actionOp) Replay(w ethdb.KeyValueWriter) error {
	if op.action != nil {
		op.action(w)
		if op.once {
			op.action = nil
		}
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
	w.data[amino.BytesToStr(key)] = preCommitValue{
		value: value,
		ele:   w.store.preCommitTail,
	}
	return nil
}

func (w preCommitMap) Delete(key []byte) error {
	w.data[amino.BytesToStr(key)] = preCommitValue{
		deleted: true,
		ele:     w.store.preCommitTail,
	}
	return nil
}

func (w preCommitMap) Len() int {
	return len(w.data)
}

type preCommitClearMap preCommitMap

func (w preCommitClearMap) Put(key []byte, _ []byte) error {
	if v, ok := w.data[string(key)]; ok {
		if v.ele == w.store.waitClearPtr {
			delete(w.data, amino.BytesToStr(key))
			atomic.AddInt64(&w.store.clearNum, 1)
		}
	}
	return nil
}

func (w preCommitClearMap) Delete(key []byte) error {
	if v, ok := w.data[string(key)]; ok {
		if v.ele == w.store.waitClearPtr {
			delete(w.data, amino.BytesToStr(key))
			atomic.AddInt64(&w.store.clearNum, 1)
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
	preCommitPtr  *list.Element
	waitClearPtr  *list.Element

	enableCommit     bool
	disableAutoPrune bool

	commitCh chan struct{}
	clearCh  chan struct{}
	closeWg  sync.WaitGroup

	logger log.Logger

	waitClear  int64
	waitCommit int64

	clearNum int64
}

func NewAsyncKeyValueStore(db ethdb.KeyValueStore, autoClearOff bool) *AsyncKeyValueStore {
	store := &AsyncKeyValueStore{
		KeyValueStore: db,
		preCommit: preCommitMap{
			data: make(map[string]preCommitValue, 200_0000),
		},
		preCommitList:    list.New(),
		commitCh:         make(chan struct{}, 10000*10),
		clearCh:          make(chan struct{}, 10000*10),
		logger:           log.NewNopLogger(),
		disableAutoPrune: autoClearOff,
	}
	store.preCommit.store = store
	store.closeWg.Add(1)
	go store.commitRoutine()
	if !store.disableAutoPrune {
		go store.pruneRoutine()
	}
	store.preCommitPtr = store.preCommitList.PushBack(nil)
	store.waitClearPtr = store.preCommitPtr
	return store
}

func (store *AsyncKeyValueStore) SetLogger(logger log.Logger) {
	if store != nil {
		store.logger = logger
	}
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
		return common.CopyBytes(v.value), nil
	}
	return store.KeyValueStore.Get(key)
}

func (store *AsyncKeyValueStore) Put(key []byte, value []byte) error {
	key, value = common.CopyBytes(key), common.CopyBytes(value)
	task := &commitTask{
		op: &singleOp{
			key:   key,
			value: value,
		},
	}
	store.mtx.Lock()
	store.preCommitTail = store.preCommitList.PushBack(task)
	_ = store.preCommit.Put(key, value)
	store.mtx.Unlock()
	atomic.AddInt64(&store.waitCommit, 1)

	store.commitCh <- struct{}{}
	return nil
}

func (store *AsyncKeyValueStore) Delete(key []byte) error {
	key = common.CopyBytes(key)
	task := &commitTask{
		op: &singleOp{
			key:    key,
			delete: true,
		},
	}
	store.mtx.Lock()
	store.preCommitTail = store.preCommitList.PushBack(task)
	_ = store.preCommit.Delete(key)
	store.mtx.Unlock()
	atomic.AddInt64(&store.waitCommit, 1)

	store.commitCh <- struct{}{}
	return nil
}

func (store *AsyncKeyValueStore) batchWrite(player replayer) error {
	task := &commitTask{
		op: player,
	}
	store.mtx.Lock()
	store.preCommitTail = store.preCommitList.PushBack(task)
	if err := player.Replay(store.preCommit); err != nil {
		return err
	}
	store.mtx.Unlock()
	atomic.AddInt64(&store.waitCommit, 1)

	store.commitCh <- struct{}{}
	return nil
}

func (store *AsyncKeyValueStore) LogInfoAfterWriteDone(msg string, args ...interface{}) {
	if store.logger == nil {
		return
	}

	store.ActionAfterWriteDone(func() {
		store.logger.Info(msg, args...)
	}, true)
}

func (store *AsyncKeyValueStore) ActionAfterWriteDone(act func(), once bool) {
	task := &commitTask{
		op: &actionOp{
			action: func(ethdb.KeyValueWriter) {
				act()
			},
			once: once,
		},
	}
	store.mtx.Lock()
	store.preCommitTail = store.preCommitList.PushBack(task)
	store.mtx.Unlock()
	atomic.AddInt64(&store.waitCommit, 1)

	store.commitCh <- struct{}{}
}

func (store *AsyncKeyValueStore) NewBatch() ethdb.Batch {
	return newAsyncBatch(store)
}

func (store *AsyncKeyValueStore) Stat(property string) (string, error) {
	return store.KeyValueStore.Stat(property)
}

func (store *AsyncKeyValueStore) Compact(start []byte, limit []byte) error {
	return store.KeyValueStore.Compact(start, limit)
}

func (store *AsyncKeyValueStore) Close() error {
	if store == nil {
		return nil
	}
	store.mtx.Lock()
	defer store.mtx.Unlock()

	close(store.commitCh)
	store.closeWg.Wait()
	return store.KeyValueStore.Close()
}

func (store *AsyncKeyValueStore) LogStats() {
	if store == nil || store.logger == nil {
		return
	}

	store.logger.Info("AsyncKeyValueStore stats",
		"waitCommitOp", atomic.LoadInt64(&store.waitCommit),
		"waitClearOp", atomic.LoadInt64(&store.waitClear),
		"preCommitMapSize", store.preCommit.Len(),
		"clearInMap", atomic.LoadInt64(&store.clearNum),
	)
}

func (store *AsyncKeyValueStore) commitRoutine() {
	defer func() {
		close(store.clearCh)
		store.closeWg.Done()
	}()

	batchSize := 0

	for _ = range store.commitCh {
		taskEle := store.preCommitPtr.Next()

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
		atomic.AddInt64(&store.waitCommit, -1)

		store.setPreCommitPtr(taskEle)
		waitClear := atomic.AddInt64(&store.waitClear, 1)

		if waitClear >= 100 || batchSize > 1_000_000 {
			store.clearCh <- struct{}{}
			batchSize = 0
		}
	}
}

func (store *AsyncKeyValueStore) pruneRoutine() {
	for _ = range store.clearCh {
		preCommitPtr := store.getPreCommitPtr()
		for store.waitClearPtr != preCommitPtr {
			needRemove := store.waitClearPtr
			needClear := store.waitClearPtr.Next()
			commitedTask := needClear.Value.(*commitTask)
			for {
				if store.mtx.TryLock() {
					store.preCommitList.Remove(needRemove)
					store.waitClearPtr = needClear
					_ = commitedTask.op.Replay(preCommitClearMap(store.preCommit))
					store.mtx.Unlock()
					atomic.AddInt64(&store.waitClear, -1)
					break
				} else {
					time.Sleep(1 * time.Millisecond)
				}
			}
		}
	}
}

func (store *AsyncKeyValueStore) Prune() {
	if !store.disableAutoPrune {
		return
	}
	store.mtx.Lock()
	defer store.mtx.Unlock()

	preCommitPtr := store.getPreCommitPtr()
	for store.waitClearPtr != preCommitPtr {
		needRemove := store.waitClearPtr
		needClear := store.waitClearPtr.Next()
		commitedTask := needClear.Value.(*commitTask)
		store.preCommitList.Remove(needRemove)
		store.waitClearPtr = needClear
		_ = commitedTask.op.Replay(preCommitClearMap(store.preCommit))
		atomic.AddInt64(&store.waitClear, -1)
	}
}

func (store *AsyncKeyValueStore) getPreCommitPtr() *list.Element {
	return (*list.Element)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&store.preCommitPtr))))
}

func (store *AsyncKeyValueStore) setPreCommitPtr(ptr *list.Element) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&store.preCommitPtr)), unsafe.Pointer(ptr))
}

type asyncBatch struct {
	ops       multiOp
	store     *AsyncKeyValueStore
	valueSize int
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
	b.valueSize += len(value)
	return nil
}

func (b *asyncBatch) Delete(key []byte) error {
	key = common.CopyBytes(key)
	b.ops = append(b.ops, singleOp{
		key:    key,
		delete: true,
	})
	b.valueSize += len(key)
	return nil
}

func (b *asyncBatch) ValueSize() int {
	return b.valueSize
}

func (b *asyncBatch) Write() error {
	return b.store.batchWrite(b.ops)
}

func (b *asyncBatch) Reset() {
	b.ops = make(multiOp, 0, len(b.ops))
	b.valueSize = 0
}

func (b *asyncBatch) Replay(w ethdb.KeyValueWriter) error {
	return b.ops.Replay(w)
}

func (b *asyncBatch) IsSingle() bool {
	return false
}

var _ ethdb.Batch = (*asyncBatch)(nil)
var _ replayer = (*asyncBatch)(nil)
var _ replayer = (*multiOp)(nil)
var _ replayer = (*singleOp)(nil)
var _ replayer = (*actionOp)(nil)
