package iavl

import (
	"container/list"
	"sync"
)

type syncList struct {
	mtx sync.Mutex
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
	defer sl.mtx.Unlock()
	sl.List.MoveToBack(e)
}

//func (sl *syncList) Len() int {
//	sl.mtx.RLock()
//	defer sl.mtx.RUnlock()
//	return sl.List.Len()
//}
//func (sl *syncList) Front() *list.Element {
//	sl.mtx.RLock()
//	defer sl.mtx.RUnlock()
//	return sl.List.Front()
//}
//func (sl *syncList) PushBack(e interface{}) *list.Element {
//	sl.mtx.Lock()
//	defer sl.mtx.Unlock()
//	return sl.List.PushBack(e)
//}
//func (sl *syncList) Remove(e *list.Element) interface{} {
//	sl.mtx.Lock()
//	defer sl.mtx.Unlock()
//	return sl.List.Remove(e)
//}
