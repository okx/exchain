package queue

import (
	"github.com/emirpasic/gods/lists/arraylist"
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
	list      *arraylist.List
}

func NewLinkedBlockQueue() *linkedBlockingQueue {
	r := &linkedBlockingQueue{}
	r.condition = sync.NewCond(&r.Mutex)
	r.list = arraylist.New()
	return r
}

func (l *linkedBlockingQueue) Take() interface{} {
	c := l.condition
	l.Mutex.Lock()
	defer l.Mutex.Unlock()

	for l.list.Size() == 0 {
		c.Wait()
	}
	if l.list.Size() > 0 {
		task, _ := l.list.Get(0)
		l.list.Remove(0)
		return task
	}
	return nil
}
func (l *linkedBlockingQueue) Push(task interface{}) (int, error) {
	l.Mutex.Lock()
	defer l.Mutex.Unlock()
	l.list.Add(task)
	l.condition.Signal()
	return l.list.Size(), nil
}

type NonOpQueue struct {
}

func NewNonOpQueue()Queue{
	ret:=&NonOpQueue{}
	return ret
}

func (n *NonOpQueue) Take() interface{} {
	panic("not supported")
}

func (n *NonOpQueue) Push(i interface{}) (int, error) {
	return 0,nil
}

