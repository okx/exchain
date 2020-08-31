package stream

type CacheQueue struct {
	queue chan StreamContext
}

type StreamContext struct {
	blockHeight int64
	stream      *Stream
	taskData    *TaskWithData
}

func newCacheQueue(queueNum int) *CacheQueue {
	cacheQueue := &CacheQueue{
		queue: make(chan StreamContext, queueNum),
	}
	return cacheQueue
}

func (cq *CacheQueue) Start() {
	for {
		streamContext := <-cq.queue
		execute(streamContext)
	}
}

func (cq *CacheQueue) Enqueue(sc StreamContext) {
	cq.queue <- sc
}
