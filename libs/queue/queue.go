package queue

import (
	"container/list"
	"sync"
)

type Queue interface {
	Take() interface{}
	Push(interface{}) (int, error)
}

var (
	_ Queue = (*linkedBlockingQueue)(nil)
	_ Queue = (*NonOpQueue)(nil)
)

type linkedBlockingQueue struct {
	sync.Mutex
	condition *sync.Cond
	list      *list.List
}

func NewLinkedBlockQueue() *linkedBlockingQueue {
	r := &linkedBlockingQueue{}
	r.condition = sync.NewCond(&r.Mutex)
	r.list = list.New()
	return r
}

func (l *linkedBlockingQueue) Take() interface{} {
	c := l.condition
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	for l.list.Len() == 0 {
		c.Wait()
	}
	if l.list.Len() > 0 {
		task := l.list.Front()
		if nil != task {
			l.list.Remove(task)
			return task.Value
		}
	}
	return nil
}
func (l *linkedBlockingQueue) Push(task interface{}) (int, error) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	l.list.PushBack(task)
	l.condition.Signal()
	return l.list.Len(), nil
}

type NonOpQueue struct {
}

func NewNonOpQueue() Queue {
	ret := &NonOpQueue{}
	return ret
}

func (n *NonOpQueue) Take() interface{} {
	panic("not supported")
}

func (n *NonOpQueue) Push(i interface{}) (int, error) {
	return 0, nil
}
