package stream

import (
	"fmt"
	"time"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/x/stream/analyservice"
	"github.com/okex/exchain/x/stream/common/kline"
	pushservicetypes "github.com/okex/exchain/x/stream/pushservice/types"
	"github.com/okex/exchain/x/stream/types"
	"github.com/okex/exchain/x/stream/websocket"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	k.stream.logger.Debug(fmt.Sprintf("stream endblock begin------%d", ctx.BlockHeight()))
	if k.stream.engines == nil {
		k.stream.logger.Debug("stream engine is not enable")
		return
	}

	// prepare task data
	sd := createStreamTaskWithData(ctx, k.stream)
	sc := Context{
		blockHeight: ctx.BlockHeight(),
		stream:      k.stream,
		taskData:    sd,
	}

	// cache111 queue
	if k.stream.cacheQueue != nil {
		// block if cache111 queue is full
		k.stream.logger.Debug(fmt.Sprintf("cache111 queue: len:%d, cap:%d, enqueue:%d", len(k.stream.cacheQueue.queue), cap(k.stream.cacheQueue.queue), sc.blockHeight))
		k.stream.cacheQueue.Enqueue(sc)
		k.metric.CacheSize.Set(float64(len(k.stream.cacheQueue.queue)))
		return
	}

	execute(sc)
	k.stream.logger.Debug("stream.Events", "size", len(ctx.EventManager().Events()))
	for i, e := range ctx.EventManager().ABCIEvents() {
		k.stream.logger.Debug("stream.Event", i, e.Type, "attrs", e.Attributes[0].String())
	}

}

func prepareStreamTask(blockHeight int64, s *Stream) (taskConst TaskConst, err error) {
	if s.distrLatestTask != nil && s.distrLatestTask.Height > blockHeight {
		return TaskPhase1NextActionJumpNextBlock, nil
	}

	// fetch distribute lock
	locked, err := s.scheduler.FetchDistLock(
		distributeLock, s.scheduler.GetLockerID(), atomTaskTimeout)

	if !locked || err != nil {
		return TaskPhase1NextActionRestart, err
	}

	tmpState, err := s.scheduler.GetDistState(latestTaskKey)
	if err != nil {
		return releaseLockWithStatus(s, TaskPhase1NextActionRestart, err)
	}

	if len(tmpState) > 0 {
		s.distrLatestTask, err = parseTaskFromJSON(tmpState)
		if err != nil {
			return releaseLockWithStatus(s, TaskPhase1NextActionRestart, err)
		}
		if s.distrLatestTask.Height > blockHeight {
			return releaseLockWithStatus(s, TaskPhase1NextActionJumpNextBlock, nil)
		}
		if s.distrLatestTask.Height == blockHeight {
			if s.distrLatestTask.GetStatus() == TaskStatusSuccess {
				return releaseLockWithStatus(s, TaskPhase1NextActionJumpNextBlock, nil)
			}
			return TaskPhase1NextActionReturnTask, nil
		}
		if s.distrLatestTask.Height+1 == blockHeight {
			return TaskPhase1NextActionNewTask, nil
		}
		return releaseLockWithStatus(s, TaskPhase1NextActionUnknown,
			fmt.Errorf("error: EndBlock-(%d) should never run into here, distrLatestBlock: %+v",
				blockHeight, s.distrLatestTask))
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

func createStreamTaskWithData(ctx sdk.Context, s *Stream) *TaskWithData {
	sd := TaskWithData{}
	sd.Task = NewTask(ctx.BlockHeight())
	sd.dataMap = make(map[Kind]types.IStreamData)

	for engineType := range s.engines {
		streamKind, ok := EngineKind2StreamKindMap[engineType]
		if ok {
			sd.Task.DoneMap[streamKind] = false
		}

		var data types.IStreamData
		switch engineType {
		case EngineAnalysisKind:
			adata := analyservice.NewDataAnalysis()
			adata.SetData(ctx, s.orderKeeper, s.tokenKeeper, s.Cache)
			data = adata
		case EngineNotifyKind:
			pBlock := pushservicetypes.NewRedisBlock()
			pBlock.SetData(ctx, s.orderKeeper, s.tokenKeeper, s.dexKeeper, s.swapKeeper, s.Cache)
			data = pBlock
		case EngineKlineKind:
			pData := kline.NewKlineData()
			pData.SetData(ctx, s.orderKeeper, s.Cache)
			// should init token pair map here
			kline.InitTokenPairMap(ctx, s.dexKeeper)
			data = pData
		case EngineWebSocketKind:
			websocket.InitialCache(ctx, s.orderKeeper, s.dexKeeper, s.logger)
			wsdata := websocket.NewPushData()
			wsdata.SetData(ctx, s.orderKeeper, s.tokenKeeper, s.dexKeeper, s.swapKeeper, s.Cache)
			data = wsdata
		}

		sd.dataMap[streamKind] = data
	}

	return &sd
}

func executeStreamTask(s *Stream, task *TaskWithData) (taskConst TaskConst, err error) {
	s.logger.Debug(fmt.Sprintf("executeStreamTask: task %+v, data: %+v", *task.Task, task.dataMap))
	s.taskChan <- task
	taskResult := <-s.resultChan
	s.logger.Debug(fmt.Sprintf("executeStreamTask: taskResult %+v", taskResult))

	stateStr := taskResult.toJSON()
	success, err := s.scheduler.UnlockDistLockWithState(
		distributeLock, s.scheduler.GetLockerID(), latestTaskKey, stateStr)

	if success && err == nil {
		s.distrLatestTask = &taskResult
		if s.distrLatestTask.GetStatus() != TaskStatusSuccess {
			return TaskPhase2NextActionRestart, nil
		}
		return TaskPhase2NextActionJumpNextBlock, nil
	}

	return TaskPhase2NextActionRestart, err

}

func execute(sc Context) {
	for {
		p1Status, p1err := prepareStreamTask(sc.blockHeight, sc.stream)
		if p1err != nil {
			sc.stream.logger.Error(p1err.Error())
		}
		sc.stream.logger.Debug(fmt.Sprintf("P1Status: %s", TaskConstDesc[p1Status]))
		switch p1Status {
		case TaskPhase1NextActionRestart:
			time.Sleep(1500 * time.Millisecond)
			continue
		case TaskPhase1NextActionUnknown:
			err := fmt.Errorf("stream unexpected exception, %+v", p1err)
			panic(err)
		case TaskPhase1NextActionJumpNextBlock:
			return
		default:
			if p1Status != TaskPhase1NextActionNewTask {
				sc.taskData.Task = sc.stream.distrLatestTask
			}
			p2Status, p2err := executeStreamTask(sc.stream, sc.taskData)
			if p2err != nil {
				sc.stream.logger.Error(p2err.Error())
			}

			sc.stream.logger.Debug(fmt.Sprintf("P2Status: %s", TaskConstDesc[p2Status]))

			switch p2Status {
			case TaskPhase2NextActionRestart:
				time.Sleep(5000 * time.Millisecond)
			case TaskPhase2NextActionJumpNextBlock:
				return
			}
		}
	}
}
