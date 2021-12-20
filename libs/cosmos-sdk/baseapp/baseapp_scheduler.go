package baseapp

type TaskScheduler struct {
	nextTaskId int
	taskNum int
	taskChan chan task
	sortChan chan task
	abciChan chan task
}

func newTaskScheduler(part1RoutineNum int) *TaskScheduler {
	ts := &TaskScheduler{}
	ts.taskChan = make(chan task, 10000)
	ts.sortChan = make(chan task, 10000)
	ts.abciChan = make(chan task, 10000)

	for i := 0; i < part1RoutineNum; i++ {
		go ts.part1Routine()
	}
	go ts.sortRoutine()
	go ts.abciRoutine()
	return ts
}

func (ts *TaskScheduler) part1Routine()  {
	for t := range ts.taskChan {
		t.part1()
		ts.sortChan <- t
	}
}

func (ts *TaskScheduler) sortRoutine()  {
	// run part2 after all part1 finished
	var taskMap = make(map[int]task)
	for t := range ts.sortChan {
		taskMap[t.id()] = t
		if len(taskMap) == ts.taskNum {
			for {
				if next, ok := taskMap[ts.nextTaskId]; ok {
					ts.abciChan <- next
					delete(taskMap, ts.nextTaskId)
					ts.nextTaskId++
				} else {
					break
				}
			}
		}
	}
}


func (ts *TaskScheduler) sortRoutine_concurrently()  {
	// run part2 and part1 concurrently
	var taskMap = make(map[int]task)
	for t := range ts.sortChan {
		if t.id() == ts.nextTaskId {
			ts.abciChan <- t
			ts.nextTaskId++
			for {
				if next, ok := taskMap[ts.nextTaskId]; ok {
					ts.abciChan <- next
					delete(taskMap, ts.nextTaskId)
					ts.nextTaskId++
				} else {
					break
				}
			}
		} else {
			taskMap[t.id()] = t
		}
	}
}

func (ts *TaskScheduler) abciRoutine()  {
	for t := range ts.abciChan {
		t.part2()
	}
}

func (ts *TaskScheduler) start(taskList []task)  {
	ts.nextTaskId = 0
	ts.taskNum = len(taskList)
	for _, task := range taskList {
		ts.taskChan <- task
	}
}

