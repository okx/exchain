package stream

import (
	"fmt"
	"time"

	"github.com/okex/okchain/x/stream/quoteslite"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/stream/analyservice"
	"github.com/okex/okchain/x/stream/pulsarclient"
	pushservicetypes "github.com/okex/okchain/x/stream/pushservice/types"
	"github.com/okex/okchain/x/stream/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) {
	k.stream.logger.Debug(fmt.Sprintf("stream endblock begin------%d", ctx.BlockHeight()))
	if k.stream.engines == nil {
		k.stream.logger.Debug("stream engine is not enable")
		return
	}

	// prepare task data
	sd := createStreamTaskWithData(ctx, k.stream)
	sc := StreamContext{
		blockHeight: ctx.BlockHeight(),
		stream:      k.stream,
		taskData:    sd,
	}

	// cache queue
	if k.stream.cacheQueue != nil {
		// block if cache queue is full
		k.stream.logger.Debug(fmt.Sprintf("cache queue: len:%d, cap:%d, enqueue:%d", len(k.stream.cacheQueue.queue), cap(k.stream.cacheQueue.queue), sc.blockHeight))
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
		return STREAM_TASK_PHRASE1_NEXT_ACTION_JUMP_NEXT_BLK, nil
	}

	// fetch distribute lock
	locked, err := s.scheduler.FetchDistLock(
		DISTRIBUTE_STREAM_LOCK, s.scheduler.GetLockerId(), atomTaskTimeout)

	if !locked || err != nil {
		return STREAM_TASK_PHRASE1_NEXT_ACTION_RESTART, err
	}

	tmpState, err := s.scheduler.GetDistState(LATEST_STREAM_TASK_KEY)
	if err != nil {
		return releaseLockWithStatus(s, STREAM_TASK_PHRASE1_NEXT_ACTION_RESTART, err)
	}

	if len(tmpState) > 0 {
		s.distrLatestTask, _ = parseTaskFromJsonStr(tmpState)
		if s.distrLatestTask.Height > blockHeight {
			return releaseLockWithStatus(s, STREAM_TASK_PHRASE1_NEXT_ACTION_JUMP_NEXT_BLK, nil)
		} else {
			if s.distrLatestTask.Height == blockHeight {
				if s.distrLatestTask.GetStatus() == STREAM_TASK_STATUS_SUCCESS {
					return releaseLockWithStatus(s, STREAM_TASK_PHRASE1_NEXT_ACTION_JUMP_NEXT_BLK, nil)
				} else {
					return STREAM_TASK_PHRASE1_NEXT_ACTION_RERUN_TASK, nil
				}
			} else {
				if s.distrLatestTask.Height+1 == blockHeight {
					return STREAM_TASK_PHRASE1_NEXT_ACTION_NEW_TASK, nil
				} else {
					return releaseLockWithStatus(s, STREAM_TASK_PHRASE1_NEXT_ACTION_UNKNOWN,
						fmt.Errorf("EndBlock-(%d) should never run into here, distrLatestBlock: %+v",
							blockHeight, s.distrLatestTask))
				}
			}
		}
	} else {
		return STREAM_TASK_PHRASE1_NEXT_ACTION_NEW_TASK, nil
	}

}

func releaseLockWithStatus(s *Stream, taskConst TaskConst, err error) (TaskConst, error) {
	rSuccess, rErr := s.scheduler.ReleaseDistLock(DISTRIBUTE_STREAM_LOCK, s.scheduler.GetLockerId())
	if rSuccess == false || rErr != nil {
		return STREAM_TASK_PHRASE1_NEXT_ACTION_RESTART, rErr
	} else {
		return taskConst, err
	}

}

func createStreamTaskWithData(ctx sdk.Context, s *Stream) *TaskWithData {

	sd := TaskWithData{}
	sd.Task = NewTask(ctx.BlockHeight())
	sd.dataMap = make(map[StreamKind]types.IStreamData)

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
			pBlock.SetData(ctx, s.orderKeeper, s.tokenKeeper, s.dexKeeper, s.Cache)
			data = pBlock
		case EngineKlineKind:
			pData := pulsarclient.NewPulsarData()
			pData.SetData(ctx, s.orderKeeper, s.Cache)
			// should init token pair map here
			pulsarclient.InitTokenPairMap(ctx, s.dexKeeper)
			data = pData
		case EngineWebSocketKind:
			quoteslite.InitialCache(ctx, s.orderKeeper, s.dexKeeper, s.logger)
			wsdata := quoteslite.NewWebSocketPushData()
			wsdata.SetData(ctx, s.orderKeeper, s.tokenKeeper, s.dexKeeper, s.Cache)
			data = wsdata
		}

		sd.dataMap[streamKind] = data
	}

	return &sd
}

func executeStreamTask(s *Stream, sd *TaskWithData) (taskConst TaskConst, err error) {

	s.logger.Debug(fmt.Sprintf("executeStreamTask: task %+v, data: %+v", *sd.Task, sd.dataMap))
	s.taskChan <- *sd
	taskResult := <-s.resultChan
	s.logger.Debug(fmt.Sprintf("executeStreamTask: taskResult %+v", taskResult))

	stateStr := taskResult.toJsonStr()
	success, err := s.scheduler.UnlockDistLockWithState(
		DISTRIBUTE_STREAM_LOCK, s.scheduler.GetLockerId(), LATEST_STREAM_TASK_KEY, stateStr)

	if success && err == nil {
		s.distrLatestTask = &taskResult
		if s.distrLatestTask.GetStatus() != STREAM_TASK_STATUS_SUCCESS {
			return STREAM_TASK_PHRASE2_NEXT_ACTION_RESTART, nil
		} else {
			return STREAM_TASK_PHRASE2_NEXT_ACTION_JUMP_NEXT_BLK, nil
		}
	}

	return STREAM_TASK_PHRASE2_NEXT_ACTION_RESTART, err

}

func execute(sc StreamContext) {
	for {
		p1Status, p1err := prepareStreamTask(sc.blockHeight, sc.stream)
		if p1err != nil {
			sc.stream.logger.Error(p1err.Error())
		}
		sc.stream.logger.Debug(fmt.Sprintf("P1Status: %s", StreamConstDesc[p1Status]))
		switch p1Status {
		case STREAM_TASK_PHRASE1_NEXT_ACTION_RESTART:
			time.Sleep(1500 * time.Millisecond)
			continue
		case STREAM_TASK_PHRASE1_NEXT_ACTION_UNKNOWN:
			err := fmt.Errorf("StreamPlugin's unexpected exception, %+v", p1err)
			panic(err)
		case STREAM_TASK_PHRASE1_NEXT_ACTION_JUMP_NEXT_BLK:
			return
		default:
			if p1Status != STREAM_TASK_PHRASE1_NEXT_ACTION_NEW_TASK {
				sc.taskData.Task = sc.stream.distrLatestTask
			}
			p2Status, p2err := executeStreamTask(sc.stream, sc.taskData)
			if p2err != nil {
				sc.stream.logger.Error(p2err.Error())
			}

			sc.stream.logger.Debug(fmt.Sprintf("P2Status: %s", StreamConstDesc[p2Status]))

			switch p2Status {
			case STREAM_TASK_PHRASE2_NEXT_ACTION_RESTART:
				time.Sleep(5000 * time.Millisecond)
			case STREAM_TASK_PHRASE2_NEXT_ACTION_JUMP_NEXT_BLK:
				return
			}
		}
	}
}
