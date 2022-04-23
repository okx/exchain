package iavl

import (
	"container/list"
	"sync"
)

type syncList struct {
	mtx sync.RWMutex
	*list.List
}

func newSyncList() *syncList {
	sl := &syncList{
		List: list.New(),
	}
	return sl
}

func (sl *syncList) MoveToBack(e *list.Element) {
	sl.mtx.Lock()
	sl.List.MoveToBack(e)
	sl.mtx.Unlock()
}

func (sl *syncList) Len() (l int) {
	sl.mtx.RLock()
	l = sl.List.Len()
	sl.mtx.RUnlock()
	return
}

func (sl *syncList) Front() (front *list.Element) {
	sl.mtx.RLock()
	front = sl.List.Front()
	sl.mtx.RUnlock()
	return
}

func (sl *syncList) RemoveFront() interface{} {
	sl.mtx.Lock()
	oldest := sl.List.Front()
	ret := sl.List.Remove(oldest)
	sl.mtx.Unlock()
	return ret
}

func (sl *syncList) PushBack(e interface{}) *list.Element {
	sl.mtx.Lock()
	ret := sl.List.PushBack(e)
	sl.mtx.Unlock()
	return ret
}

func (sl *syncList) Remove(e *list.Element) (removed interface{}) {
	sl.mtx.Lock()
	removed = sl.List.Remove(e)
	sl.mtx.Unlock()
	return
}
