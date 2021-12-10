package iavl

type Task interface {
	ExecTask()
}

type TaskQueue struct {
	que           chan Task
	consumerCount int
}

const (
	maxConsumers           = 10000
	defaultQueueBufferSize = 100000
	defaultConsumers       = 32
)

func NewDefaultTaskQueue() *TaskQueue {
	return NewTaskQueue(defaultQueueBufferSize, defaultConsumers)
}

func NewTaskQueue(queueBufferSize int, consumer int) *TaskQueue {
	if consumer > maxConsumers {
		consumer = maxConsumers
	}
	taskQueue := &TaskQueue{
		que:           make(chan Task, queueBufferSize),
		consumerCount: consumer,
	}
	for i := 0; i < consumer; i++ {
		go taskQueue.consumeTask()
	}
	return taskQueue
}

func (taskQueue *TaskQueue) consumeTask() {
	for task := range taskQueue.que {
		task.ExecTask()
	}
}

func (taskQueue *TaskQueue) SendTask(task Task) {
	taskQueue.que <- task
}
