package watcher

import (
	"container/list"
	"sync"
)

type workQueue struct {
	mtx   sync.Mutex
	queue *list.List
}

func NewWorkQueue() *workQueue {
	return &workQueue{queue: list.New()}
}

func (q *workQueue) PushBack(data *WatchMessage) {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	q.queue.PushBack(data)
}

func (q *workQueue) PopFront() *WatchMessage {
	q.mtx.Lock()
	defer q.mtx.Unlock()
	if q.queue.Len() > 1 {
		q.queue.Remove(q.queue.Front())
	}
	return nil
}
