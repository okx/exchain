package stream

import (
	"encoding/json"
	"fmt"

	"github.com/okex/exchain/x/stream/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

type TaskConst int

const (
	TaskStatusInvalid TaskConst = 0 + iota
	TaskStatusSuccess
	TaskStatusStatusFail
	TaskStatusPartialSuccess
)

const (
	// Phase 1
	TaskPhase1NextActionRestart TaskConst = 100 + iota
	TaskPhase1NextActionJumpNextBlock
	TaskPhase1NextActionNewTask
	TaskPhase1NextActionReturnTask
	TaskPhase1NextActionUnknown

	// Phase 2
	TaskPhase2NextActionRestart TaskConst = 200 + iota
	TaskPhase2NextActionJumpNextBlock
)

var TaskConstDesc = map[TaskConst]string{
	TaskStatusInvalid:                 "STREAM_TASK_STATUS_INVALID",
	TaskStatusSuccess:                 "STREAM_TASK_STATUS_SUCCESS",
	TaskStatusStatusFail:              "STREAM_TASK_STATUS_FAIL",
	TaskStatusPartialSuccess:          "STREAM_TASK_STATUS_PARTITIAL_SUCCESS",
	TaskPhase1NextActionRestart:       "STREAM_TASK_PHRASE1_NEXT_ACTION_RESTART",
	TaskPhase1NextActionJumpNextBlock: "STREAM_TASK_PHRASE1_NEXT_ACTION_JUMP_NEXT_BLK",
	TaskPhase1NextActionNewTask:       "STREAM_TASK_PHRASE1_NEXT_ACTION_NEW_TASK",
	TaskPhase1NextActionReturnTask:    "STREAM_TASK_PHRASE1_NEXT_ACTION_RERUN_TASK",
	TaskPhase1NextActionUnknown:       "STREAM_TASK_PHRASE1_NEXT_ACTION_UNKNOWN",
	TaskPhase2NextActionRestart:       "STREAM_TASK_PHRASE2_NEXT_ACTION_RESTART",
	TaskPhase2NextActionJumpNextBlock: "STREAM_TASK_PHRASE2_NEXT_ACTION_JUMP_NEXT_BLK",
}

type Task struct {
	Height    int64         `json:"Height"`
	DoneMap   map[Kind]bool `json:"DoneMap"`
	UpdatedAt int64         `json:"UpdatedAt"`
}

func NewTask(blockHeight int64) *Task {
	doneMap := make(map[Kind]bool)
	return &Task{
		Height:  blockHeight,
		DoneMap: doneMap,
	}
}

func parseTaskFromJSON(s string) (*Task, error) {
	st := Task{}
	e := json.Unmarshal([]byte(s), &st)
	return &st, e
}

func (t *Task) toJSON() string {
	r, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(r)
}

func (t *Task) validAtomTaskCount() int {
	cnt := 0
	for _, done := range t.DoneMap {
		if !done {
			cnt++
		}
	}
	return cnt
}

func (t *Task) GetStatus() TaskConst {
	allTaskCnt := len(t.DoneMap)
	unDoneCnt := t.validAtomTaskCount()
	doneCnt := allTaskCnt - unDoneCnt

	if doneCnt == allTaskCnt {
		return TaskStatusSuccess
	}

	if doneCnt > 0 && doneCnt < allTaskCnt {
		return TaskStatusPartialSuccess
	}

	if doneCnt == 0 {
		return TaskStatusStatusFail
	}

	return TaskStatusInvalid
}

type TaskWithData struct {
	*Task
	dataMap map[Kind]types.IStreamData
}

type AtomTaskResult struct {
	sType       Kind
	successDone bool
}

type AtomTaskRunner struct {
	data   types.IStreamData
	engine types.IStreamEngine
	sType  Kind
	result chan AtomTaskResult
	logger log.Logger
}

func (r *AtomTaskRunner) run() {
	taskSuccess := false

	defer func() {
		if e := recover(); e != nil {
			r.logger.Error(fmt.Sprintf("AtomTaskRunner Panic: %+v", e))
			r.result <- AtomTaskResult{sType: r.sType, successDone: false}
		} else {
			r.result <- AtomTaskResult{sType: r.sType, successDone: taskSuccess}
		}
	}()

	if r.data == nil {
		taskSuccess = true
	} else {
		r.engine.Write(r.data, &taskSuccess)
	}
}
