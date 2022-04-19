package watcher

const (
	DefaultTxChanBuffer = 64
	DefaultTxWorkers    = 16
)

func (w *Watcher) txRoutine() {
	if !w.Enabled() {
		return
	}

	w.txChan = make(chan func(), DefaultTxChanBuffer)
	for i := 0; i < DefaultTxWorkers; i++ {
		go w.txWorker(w.txChan)
	}
}

func (w *Watcher) txWorker(jobs <-chan func()) {
	for job := range jobs {
		job()
	}
}

func (w *Watcher) dispatchTxJob(f func()) {
	select {
	case w.txChan <- f:
	default:
		w.log.Error("watch dispatch tx job too busy.")
		go func() {
			w.txChan <- f
		}()
	}
}
