package iavl

import (
	"fmt"
	"sync"
)

type CallBackTask interface {
	ExecTask()
	ParentTask() CallBackTask
	UpdateParentTask(interface{}, interface{})
}

type UpdateNodeTask struct {
	que         *TaskQueue
	node        *Node
	ndb         *nodeDB
	isStartFunc bool
	parent      CallBackTask
	isLeft      bool
	mtx         sync.Mutex
	done        chan struct{}
}

func (task *UpdateNodeTask) ExecTask() {
	if task.isStartFunc {
		fmt.Println("start task", task.node.height, string(task.node.key))
		task.StartFunc()
		fmt.Println("start task done", task.node.height, string(task.node.key))
	} else {
		fmt.Println("callback task", task.node.height, string(task.node.key))
		task.CallBackFunc()
		fmt.Println("callback task done", task.node.height, string(task.node.key))
	}
}

func (task *UpdateNodeTask) StartFunc() {
	task.mtx.Lock()
	defer task.mtx.Unlock()
	node := task.node
	if node.persisted || node.prePersisted {
		task.UpdateParentTask(node.hash, task.isLeft)
		return
	}
	if node.leftNode != nil || node.rightNode != nil {
		if node.leftNode != nil {
			isStartFunc := true
			if node.leftNode.isLeaf() {
				isStartFunc = false
			}
			leftTask := UpdateNodeTask{
				node:        node.leftNode,
				que:         task.que,
				ndb:         task.ndb,
				isStartFunc: isStartFunc,
				parent:      task,
				isLeft:      true,
			}
			task.que.que <- &leftTask
		}
		if node.rightNode != nil {
			isStartFunc := true
			if node.rightNode.isLeaf() {
				isStartFunc = false
			}
			rightTask := UpdateNodeTask{
				node:        node.rightNode,
				que:         task.que,
				ndb:         task.ndb,
				isStartFunc: isStartFunc,
				parent:      task,
				isLeft:      false,
			}
			task.que.que <- &rightTask
		}
	} else {
		panic("unexpected logic")
	}

}

func (task *UpdateNodeTask) CallBackFunc() {
	task.mtx.Lock()
	defer task.mtx.Unlock()
	node := task.node
	node._hash()
	task.ndb.saveNodeToPrePersistCache(node)

	node.leftNode = nil
	node.rightNode = nil
	task.UpdateParentTask(node.hash, task.isLeft)
}

func (task *UpdateNodeTask) SetStartFunc() {
	task.isStartFunc = true
}

func (task *UpdateNodeTask) SetCallBackFunc() {
	task.isStartFunc = false
}

func (task *UpdateNodeTask) ParentTask() CallBackTask {
	return task.parent
}

func (task *UpdateNodeTask) UpdateParentTask(hash interface{}, isLeft interface{}) {
	if task.parent == nil {
		task.done <- struct{}{}
		return
	}
	pTask := task.parent.(*UpdateNodeTask)
	hashBytes := hash.([]byte)
	isLeftBool := isLeft.(bool)

	pTask.mtx.Lock()
	defer pTask.mtx.Unlock()
	node := pTask.node
	if isLeftBool {
		node.leftHash = hashBytes
	} else {
		node.rightHash = hashBytes
	}
	if node.leftHash != nil && node.rightHash != nil {
		pTask.SetCallBackFunc()
		task.que.que <- pTask
	}
}
