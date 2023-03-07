package infura

import (
	"fmt"
	"time"

	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	k.stream.logger.Debug("infura EndBlocker begin")
	if !k.stream.enable {
		k.stream.logger.Debug("infura engine is not enable")
		return
	}
	// prepare task data
	sc := StreamContext{
		blockHeight: ctx.BlockHeight(),
		stream:      k.stream,
		task:        newTask(ctx.BlockHeight(), k.stream.cache),
	}

	// cache queue
	if k.stream.cacheQueue != nil {
		k.stream.logger.Debug(fmt.Sprintf("cache queue: len:%d, cap:%d, enqueue:%d",
			len(k.stream.cacheQueue.queue), cap(k.stream.cacheQueue.queue), sc.blockHeight))
		// block if cache queue is full
		k.stream.cacheQueue.Enqueue(sc)
		k.metric.CacheSize.Set(float64(len(k.stream.cacheQueue.queue)))
		return
	}

	execute(sc)
	k.stream.logger.Debug("infura EndBlocker end")
}

func prepareStreamTask(blockHeight int64, ctx StreamContext) (taskConst TaskConst, err error) {
	if ctx.task != nil && ctx.task.Height > blockHeight {
		return TaskPhase1NextActionJumpNextBlock, nil
	}

	// fetch distribute lock
	locked, err := ctx.stream.scheduler.FetchDistLock(
		distributeLock, ctx.stream.scheduler.GetLockerID(), taskTimeout)

	if !locked || err != nil {
		return TaskPhase1NextActionRestart, err
	}

	state := ctx.stream.scheduler.GetDistState(latestTaskKey)
	if len(state) > 0 {
		ctx.task, err = parseTaskFromJSON(state)
		if err != nil {
			return releaseLockWithStatus(ctx.stream, TaskPhase1NextActionRestart, err)
		}
		if ctx.task.Height > blockHeight {
			return releaseLockWithStatus(ctx.stream, TaskPhase1NextActionJumpNextBlock, nil)
		}
		if ctx.task.Height == blockHeight {
			if ctx.task.GetStatus() == TaskStatusSuccess {
				return releaseLockWithStatus(ctx.stream, TaskPhase1NextActionJumpNextBlock, nil)
			}
			return TaskPhase1NextActionReturnTask, nil
		}
		if ctx.task.Height+1 == blockHeight {
			return TaskPhase1NextActionNewTask, nil
		}
		return releaseLockWithStatus(ctx.stream, TaskPhase1NextActionUnknown,
			fmt.Errorf("error: EndBlock-(%d) should never run into here, distrLatestBlock: %+v",
				blockHeight, ctx.task))
	}
	return TaskPhase1NextActionNewTask, nil

}

func releaseLockWithStatus(s *Stream, taskConst TaskConst, err error) (TaskConst, error) {
	rSuccess, rErr := s.scheduler.ReleaseDistLock(distributeLock, s.scheduler.GetLockerID())
	if !rSuccess || rErr != nil {
		return TaskPhase1NextActionRestart, rErr
	}
	return taskConst, err
}

func executeStreamTask(s *Stream, task *Task) (taskConst TaskConst, err error) {
	done := s.engine.Write(task.Data)
	task.Done = done
	stateStr := task.toJSON()
	success, err := s.scheduler.UnlockDistLockWithState(
		distributeLock, s.scheduler.GetLockerID(), latestTaskKey, stateStr)
	if success && err == nil {
		if task.GetStatus() != TaskStatusSuccess {
			return TaskPhase2NextActionRestart, nil
		}
		return TaskPhase2NextActionJumpNextBlock, nil
	}

	return TaskPhase2NextActionRestart, err

}

func execute(ctx StreamContext) {
	for {
		p1Status, p1err := prepareStreamTask(ctx.blockHeight, ctx)
		if p1err != nil {
			ctx.stream.logger.Error(p1err.Error())
		}
		ctx.stream.logger.Debug(fmt.Sprintf("P1Status: %s", TaskConstDesc[p1Status]))
		switch p1Status {
		case TaskPhase1NextActionRestart:
			time.Sleep(1500 * time.Millisecond)
			continue
		case TaskPhase1NextActionUnknown:
			err := fmt.Errorf("infura unexpected exception, %+v", p1err)
			panic(err)
		case TaskPhase1NextActionJumpNextBlock:
			return
		default:
			p2Status, p2err := executeStreamTask(ctx.stream, ctx.task)
			if p2err != nil {
				ctx.stream.logger.Error(p2err.Error())
			}

			ctx.stream.logger.Debug(fmt.Sprintf("P2Status: %s", TaskConstDesc[p2Status]))

			switch p2Status {
			case TaskPhase2NextActionRestart:
				time.Sleep(5000 * time.Millisecond)
			case TaskPhase2NextActionJumpNextBlock:
				return
			}
		}
	}
}
