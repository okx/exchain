package stream

import (
	"encoding/json"
	"fmt"

	"github.com/okex/okchain/x/stream/types"
	"github.com/tendermint/tendermint/libs/log"
)

type TaskConst int

const (
	STREAM_TASK_STATUS_INVALID TaskConst = 0 + iota
	STREAM_TASK_STATUS_SUCCESS
	STREAM_TASK_STATUS_FAIL
	STREAM_TASK_STATUS_PARTITIAL_SUCCESS
)

const (
	// Phrase 1
	STREAM_TASK_PHRASE1_NEXT_ACTION_RESTART TaskConst = 100 + iota
	STREAM_TASK_PHRASE1_NEXT_ACTION_JUMP_NEXT_BLK
	STREAM_TASK_PHRASE1_NEXT_ACTION_NEW_TASK
	STREAM_TASK_PHRASE1_NEXT_ACTION_RERUN_TASK
	STREAM_TASK_PHRASE1_NEXT_ACTION_UNKNOWN

	// Phrase 2
	STREAM_TASK_PHRASE2_NEXT_ACTION_RESTART TaskConst = 200 + iota
	STREAM_TASK_PHRASE2_NEXT_ACTION_JUMP_NEXT_BLK
)

var StreamConstDesc = map[TaskConst]string{
	STREAM_TASK_STATUS_INVALID:                    "STREAM_TASK_STATUS_INVALID",
	STREAM_TASK_STATUS_SUCCESS:                    "STREAM_TASK_STATUS_SUCCESS",
	STREAM_TASK_STATUS_FAIL:                       "STREAM_TASK_STATUS_FAIL",
	STREAM_TASK_STATUS_PARTITIAL_SUCCESS:          "STREAM_TASK_STATUS_PARTITIAL_SUCCESS",
	STREAM_TASK_PHRASE1_NEXT_ACTION_RESTART:       "STREAM_TASK_PHRASE1_NEXT_ACTION_RESTART",
	STREAM_TASK_PHRASE1_NEXT_ACTION_JUMP_NEXT_BLK: "STREAM_TASK_PHRASE1_NEXT_ACTION_JUMP_NEXT_BLK",
	STREAM_TASK_PHRASE1_NEXT_ACTION_NEW_TASK:      "STREAM_TASK_PHRASE1_NEXT_ACTION_NEW_TASK",
	STREAM_TASK_PHRASE1_NEXT_ACTION_RERUN_TASK:    "STREAM_TASK_PHRASE1_NEXT_ACTION_RERUN_TASK",
	STREAM_TASK_PHRASE1_NEXT_ACTION_UNKNOWN:       "STREAM_TASK_PHRASE1_NEXT_ACTION_UNKNOWN",
	STREAM_TASK_PHRASE2_NEXT_ACTION_RESTART:       "STREAM_TASK_PHRASE2_NEXT_ACTION_RESTART",
	STREAM_TASK_PHRASE2_NEXT_ACTION_JUMP_NEXT_BLK: "STREAM_TASK_PHRASE2_NEXT_ACTION_JUMP_NEXT_BLK",
}

type Task struct {
	Height    int64               `json:"Height"`
	DoneMap   map[StreamKind]bool `json:"DoneMap"`
	UpdatedAt int64               `json:"UpdatedAt"`
}

func NewTask(blockHeight int64) *Task {
	doneMap := make(map[StreamKind]bool)
	//doneMap[StreamMysqlKind] = false
	//doneMap[StreamRedisKind] = false
	//doneMap[StreamPulsarKind] = false
	//doneMap[StreamWebSocketKind] = false

	return &Task{
		Height:  blockHeight,
		DoneMap: doneMap,
	}
}

func parseTaskFromJsonStr(s string) (*Task, error) {
	st := Task{}
	e := json.Unmarshal([]byte(s), &st)
	return &st, e
}

func (t *Task) toJsonStr() string {
	r, _ := json.Marshal(t)
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
		return STREAM_TASK_STATUS_SUCCESS
	}

	if 0 < doneCnt && doneCnt < allTaskCnt {
		return STREAM_TASK_STATUS_PARTITIAL_SUCCESS
	}

	if 0 == doneCnt {
		return STREAM_TASK_STATUS_FAIL
	}

	return STREAM_TASK_STATUS_INVALID
}

type TaskWithData struct {
	*Task
	dataMap map[StreamKind]types.IStreamData
}

type AtomTaskResult struct {
	sType       StreamKind
	successDone bool
}

type AtomTaskRunner struct {
	data   types.IStreamData
	engine types.IStreamEngine
	sType  StreamKind
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
