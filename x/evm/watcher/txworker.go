package watcher

import "runtime"

const (
	DefaultTxChanBuffer = 20
	DefaultTxWorkers    = 16
)

func (w *Watcher) setTxWorkerNums() {
	w.txWorkerNums = runtime.NumCPU()
}

func (w *Watcher) txRoutine() {
	if !w.Enabled() {
		return
	}

	w.txChan = make(chan func(), DefaultTxChanBuffer)
	for i := 0; i < w.txWorkerNums; i++ {
		go w.txWorker(w.txChan)
	}
}

func (w *Watcher) txWorker(jobs <-chan func()) {
	for job := range jobs {
		job()
	}
}

func (w *Watcher) dispatchTxJob(f func()) {
	// if jobRoutine were too slow to write data  to disk
	// we have to wait
	// why: something wrong happened: such as db panic(disk maybe is full)(it should be the only reason)
	//								  UseWatchData were executed every 4 seoncds(block schedual)
	select {
	case w.txChan <- f:
	default:
		w.log.Error("watch dispatch tx job too busy.")
		go func() {
			w.txChan <- f
		}()
	}
}
