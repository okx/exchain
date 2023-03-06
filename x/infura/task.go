package infura

import (
	"encoding/json"

	"github.com/okx/okbchain/x/infura/types"
)

type TaskConst int

const (
	TaskStatusInvalid TaskConst = 0 + iota
	TaskStatusSuccess
	TaskStatusStatusFail
)

const (
	// Phase 1 task status
	TaskPhase1NextActionRestart TaskConst = 100 + iota
	TaskPhase1NextActionJumpNextBlock
	TaskPhase1NextActionNewTask
	TaskPhase1NextActionReturnTask
	TaskPhase1NextActionUnknown

	// Phase 2 task status
	TaskPhase2NextActionRestart TaskConst = 200 + iota
	TaskPhase2NextActionJumpNextBlock
)

var TaskConstDesc = map[TaskConst]string{
	TaskStatusInvalid:                 "STREAM_TASK_STATUS_INVALID",
	TaskStatusSuccess:                 "STREAM_TASK_STATUS_SUCCESS",
	TaskStatusStatusFail:              "STREAM_TASK_STATUS_FAIL",
	TaskPhase1NextActionRestart:       "STREAM_TASK_PHRASE1_NEXT_ACTION_RESTART",
	TaskPhase1NextActionJumpNextBlock: "STREAM_TASK_PHRASE1_NEXT_ACTION_JUMP_NEXT_BLK",
	TaskPhase1NextActionNewTask:       "STREAM_TASK_PHRASE1_NEXT_ACTION_NEW_TASK",
	TaskPhase1NextActionReturnTask:    "STREAM_TASK_PHRASE1_NEXT_ACTION_RERUN_TASK",
	TaskPhase1NextActionUnknown:       "STREAM_TASK_PHRASE1_NEXT_ACTION_UNKNOWN",
	TaskPhase2NextActionRestart:       "STREAM_TASK_PHRASE2_NEXT_ACTION_RESTART",
	TaskPhase2NextActionJumpNextBlock: "STREAM_TASK_PHRASE2_NEXT_ACTION_JUMP_NEXT_BLK",
}

type Task struct {
	Height    int64             `json:"height"`
	Done      bool              `json:"done"`
	UpdatedAt int64             `json:"updatedAt"`
	Data      types.IStreamData `json:"-"`
}

func newTask(blockHeight int64, cache *Cache) *Task {
	return &Task{
		Height: blockHeight,
		Done:   false,
		Data:   getStreamData(cache),
	}
}

func getStreamData(cache *Cache) types.IStreamData {
	return types.StreamData{
		TransactionReceipts: cache.GetTransactionReceipts(),
		Block:               cache.GetBlock(),
		Transactions:        cache.GetTransactions(),
		ContractCodes:       cache.GetContractCodes(),
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

func (t *Task) GetStatus() TaskConst {
	if t.Done {
		return TaskStatusSuccess
	}
	return TaskStatusStatusFail
}
