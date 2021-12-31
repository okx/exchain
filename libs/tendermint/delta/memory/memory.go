/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/31 3:37 下午
# @File : MemoryBroker.go
# @Description :
# @Attention :
*/
package memory

import (
	"fmt"
	"github.com/okex/exchain/libs/component/listener"
	"github.com/okex/exchain/libs/tendermint/delta"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/types"
	"sync"
	"sync/atomic"
	"time"
)

// TODO REPEATED CODE
var (
	_                   delta.DeltaBroker = (*MemoryBroker)(nil)
	lockerExpire                          = 4 * time.Second
	mostRecentHeightKey string
	deltaLockerKey      string
)

type MemoryBroker struct {
	mtx sync.RWMutex
	// semaphore
	lockC chan struct{}
	seq   int64
	log   log.Logger
	data  map[string]interface{}

	ttl time.Duration
	tw  *TimeWheel

	l listener.IListenerComponent
}

func NewMemoryBroker(lg log.Logger, ops ...MemoryBrokerOption) *MemoryBroker {
	ret := &MemoryBroker{
		mtx:   sync.RWMutex{},
		lockC: make(chan struct{}, 1),
		seq:   0,
		log:   lg.With("memoryBroker"),
		data:  make(map[string]interface{}),
	}
	for _, opt := range ops {
		opt(ret)
	}
	if nil == ret.l {
		ret.l = listener.DefaultNewListenerComponent(lg.With("module", "listener"), "listener")
	}
	// never mind
	ret.l.Start()
	return ret
}

func (m *MemoryBroker) ValidateBasic() error {
	const (
		mostRecentHeight = "MostRecentHeight"
		deltaLocker      = "DeltaLocker"
	)
	mostRecentHeightKey = fmt.Sprintf("dds:%d:%s", types.DeltaVersion, mostRecentHeight)
	deltaLockerKey = fmt.Sprintf("dds:%d:%s", types.DeltaVersion, deltaLocker)
	m.tw = New(time.Second, 60, func(i interface{}) {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		delete(m.data, i.(string))
	})
	m.tw.Start()
	m.log.Info("validate basic sucessfully")
	m.Start()
	return nil
}

func (m *MemoryBroker) GetLocker() bool {
	//to, cancel := context.WithTimeout(context.Background(), time.Second*3)
	//defer cancel()
	select {
	case m.lockC <- struct{}{}:
		// lock can be released by routine ,so we need atomic,but no need for compareAndSwap
		ret := atomic.AddInt64(&m.seq, 1)
		m.log.Info("lock successfully", "seq", ret)
		//return true
		//case <-to.Done():
		//	return false
		return true
	}
}

func (m *MemoryBroker) ReleaseLocker() {
	select {
	case <-m.lockC:
		m.log.Info("releaseLocker", "seq", atomic.LoadInt64(&m.seq))
	default:
		m.log.Info("lock HasBeen released automatically", "seq", atomic.LoadInt64(&m.seq))
	}
}

func (m *MemoryBroker) ResetMostRecentHeightAfterUpload(targetHeight int64, upload func(int64) bool) (bool, int64, error) {
	var (
		res bool
		mrh int64
	)
	m.mtx.RLock()
	v, ok := m.data[mostRecentHeightKey]
	m.mtx.RUnlock()
	if !ok {
		mrh = int64(0)
	} else {
		mrh = v.(int64)
	}

	if mrh < targetHeight && upload(mrh) {
		m.mtx.Lock()
		defer m.mtx.Unlock()
		m.data[mostRecentHeightKey] = targetHeight
		m.log.Info("Reset most recent height", "new-mrh", targetHeight, "old-mrh", mrh)
	} else {
	}
	return res, mrh, nil
}

func (m *MemoryBroker) Start() {
	go func() {
		tt := time.NewTicker(lockerExpire)
		for {
			select {
			case <-tt.C:
				select {
				case <-m.lockC:
					m.log.Info("release lock", "seq", atomic.LoadInt64(&m.seq))
				default:

				}
			}
		}
	}()
}

func (r *MemoryBroker) SetBlock(height int64, bytes []byte) error {
	if len(bytes) == 0 {
		return fmt.Errorf("block is empty")
	}
	key := genBlockKey(height)
	defer r.l.NotifyListener(height, key)
	r.mtx.Lock()
	r.data[key] = bytes
	r.mtx.Unlock()
	if r.ttl > 0 {
		r.tw.AddTimer(r.ttl, key, key)
	}
	r.log.Info("SetBlock", "height", height, "bytesLen", len(bytes))
	return nil
}

func (r *MemoryBroker) SetDeltas(height int64, bytes []byte) error {
	if len(bytes) == 0 {
		return fmt.Errorf("delta is empty")
	}
	key := genDeltaKey(height)
	defer r.l.NotifyListener(height, key)
	r.mtx.Lock()
	r.data[key] = bytes
	r.mtx.Unlock()
	if r.ttl > 0 {
		r.tw.AddTimer(r.ttl, key, key)
	}
	r.log.Info("setDeltas", "height", height, "bytesLen", len(bytes))
	return nil
}

func (r *MemoryBroker) GetBlock(height int64) ([]byte, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	v, exist := r.data[genBlockKey(height)]
	if !exist {
		r.log.Info("getBlock,not exist", "height", height)
		return nil, nil
	}
	r.log.Info("getBlock,", "height", height, "bytesLen", len(v.([]byte)))
	return v.([]byte), nil
}

func (r *MemoryBroker) GetDeltas(height int64) ([]byte, error, int64) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	mrh := r.getMostRecentHeight()
	key := genDeltaKey(height)
	v, exist := r.data[key]
	if !exist {
		r.log.Info("GetDeltas,not exist", "height", height)
		return nil, nil, mrh
	}
	r.log.Info("GetDeltas", "height", height, "bytesLen", len(v.([]byte)))
	return v.([]byte), nil, mrh
}

func (r *MemoryBroker) GetListener() listener.IListenerComponent {
	return r.l
}

func (r *MemoryBroker) getMostRecentHeight() (ret int64) {
	ret = -1
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	v, exist := r.data[mostRecentHeightKey]
	if !exist {
		ret = 0
	} else {
		ret = v.(int64)
	}
	return
}
func (r *MemoryBroker) RegisterBrokerDelta(h int64) <-chan interface{} {
	k := genDeltaKey(h)
	return r.l.RegisterListener(k)
}
func (r *MemoryBroker) RegisterBrokerBlock(h int64) <-chan interface{} {
	k := genBlockKey(h)
	return r.l.RegisterListener(k)
}
func genBlockKey(height int64) string {
	return fmt.Sprintf("BH:%d", height)
}

func genDeltaKey(height int64) string {
	return fmt.Sprintf("DH-%d:%d", types.DeltaVersion, height)
}
