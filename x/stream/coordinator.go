package stream

import (
	"fmt"
	"time"

	"github.com/okex/exchain/x/stream/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

const (
	atomTaskTimeout int = distributeLockTimeout * 0.98
)

type Coordinator struct {
	engineMap       map[EngineKind]types.IStreamEngine
	taskChan        chan *TaskWithData
	resultChan      chan Task
	atomTaskTimeout int // In Million Second
	logger          log.Logger
}

func NewCoordinator(logger log.Logger, taskCh chan *TaskWithData, resultCh chan Task, timeout int, engineMap map[EngineKind]types.IStreamEngine) *Coordinator {
	c := Coordinator{
		logger:          logger,
		taskChan:        taskCh,
		resultChan:      resultCh,
		atomTaskTimeout: timeout,
		engineMap:       engineMap,
	}
	return &c
}

func (c *Coordinator) prepareAtomTasks(taskDesc *TaskWithData, notify chan AtomTaskResult) []*AtomTaskRunner {
	var runners []*AtomTaskRunner
	for streamType, done := range taskDesc.DoneMap {
		if !done {
			engineType := StreamKind2EngineKindMap[streamType]
			if engineType == EngineNilKind {
				err := fmt.Errorf("stream kind: %+v not supported, no Kind found, Quite", streamType)
				panic(err)
			}

			r := AtomTaskRunner{
				data:   taskDesc.dataMap[streamType],
				engine: c.engineMap[engineType],
				result: notify,
				logger: c.logger,
				sType:  streamType,
			}

			runners = append(runners, &r)
		}
	}

	return runners
}

func (c *Coordinator) run() {
	for task := range c.taskChan {
		// outer loop, block to get streamTask from taskChan
		func() {
			validTaskCnt := task.validAtomTaskCount()
			if validTaskCnt > 0 {
				notifyCh := make(chan AtomTaskResult, validTaskCnt)
				atomRunners := c.prepareAtomTasks(task, notifyCh)
				c.logger.Debug(fmt.Sprintf("Coordinator loop: %d atomRunners prepared, %+v", len(atomRunners), atomRunners))

				for _, r := range atomRunners {
					go r.run()
				}

				timer := time.NewTimer(time.Duration(c.atomTaskTimeout * int(time.Millisecond)))
				notifyCnt := 0

				// inner loop, wait all jobs timeout or get all notified message
				for {
					select {
					case <-timer.C:
						c.logger.Error(fmt.Sprintf(
							"Coordinator: All atom tasks are forced stop becoz of %d millsecond timeout", c.atomTaskTimeout))
						notifyCnt = validTaskCnt
					case taskResult := <-notifyCh:
						task.DoneMap[taskResult.sType] = taskResult.successDone
						notifyCnt++
					}

					if validTaskCnt == notifyCnt {
						break
					}
				}

			}

			task.UpdatedAt = time.Now().Unix()
			c.resultChan <- *task.Task

		}()
	}
}
