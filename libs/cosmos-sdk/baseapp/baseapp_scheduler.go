package baseapp

type TaskScheduler struct {
	expectedId int
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
	var taskMap = make(map[int]task)
	for t := range ts.sortChan {
		if t.id() == ts.expectedId {
			ts.abciChan <- t
			ts.expectedId++
			for {
				if next, ok := taskMap[ts.expectedId]; ok {
					ts.abciChan <- next
					delete(taskMap, ts.expectedId)
					ts.expectedId++
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
	ts.expectedId = 0
	for _, task := range taskList {
		ts.taskChan <- task
	}
}

