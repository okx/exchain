package global

import (
	"sync"
)

var (
	signalMtx    sync.RWMutex
	commitSignal = make(chan struct{})
)

func init() {
	close(commitSignal)
}

func CommitLock() {
	signalMtx.Lock()
	commitSignal = make(chan struct{})
	signalMtx.Unlock()
}

func CommitUnlock() {
	close(commitSignal)
}

func WaitCommit() {
	signalMtx.RLock()
	<-commitSignal
	signalMtx.RUnlock()
}
