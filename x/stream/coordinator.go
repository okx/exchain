package stream

import (
	"errors"
	"fmt"
	"time"

	"github.com/okex/okchain/x/stream/types"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	atomTaskTimeout int = DISTRIBUTE_STREAM_LOCK_TIMEOUT * 0.98
)

type Coordinator struct {
	engineMap       map[EngineKind]types.IStreamEngine
	taskChan        chan TaskWithData
	resultChan      chan Task
	atomTaskTimeout int // In Million Second
	logger          log.Logger
}

func NewCoordinator(logger log.Logger, taskCh chan TaskWithData, resultCh chan Task, timeout int, engineMap map[EngineKind]types.IStreamEngine) *Coordinator {
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
	runners := []*AtomTaskRunner{}

	for streamType, done := range taskDesc.DoneMap {
		if !done {
			engineType := StreamKind2EngineKindMap[streamType]
			if engineType == EngineNilKind {
				err := errors.New(fmt.Sprintf("StreamKind: %+v not supported, no EngineKind found, Quite", streamType))
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
	for {
		select {

		// outer loop, block to get streamTask from taskChan
		case sTask := <-c.taskChan:
			func() {
				validTaskCnt := sTask.validAtomTaskCount()
				if validTaskCnt > 0 {
					notifyCh := make(chan AtomTaskResult, validTaskCnt)
					atomRunners := c.prepareAtomTasks(&sTask, notifyCh)
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
							sTask.DoneMap[taskResult.sType] = taskResult.successDone
							notifyCnt++
						}

						if validTaskCnt == notifyCnt {
							break
						}
					}

				}

				sTask.UpdatedAt = time.Now().Unix()
				c.resultChan <- *sTask.Task

			}()
		}
	}
}
