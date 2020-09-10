package stream

type CacheQueue struct {
	queue chan Context
}

type Context struct {
	blockHeight int64
	stream      *Stream
	taskData    *TaskWithData
}

func newCacheQueue(queueNum int) *CacheQueue {
	cacheQueue := &CacheQueue{
		queue: make(chan Context, queueNum),
	}
	return cacheQueue
}

func (cq *CacheQueue) Start() {
	for {
		streamContext := <-cq.queue
		execute(streamContext)
	}
}

func (cq *CacheQueue) Enqueue(sc Context) {
	cq.queue <- sc
}
