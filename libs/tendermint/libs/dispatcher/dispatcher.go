package dispatcher

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/viper"
)

const (
	DEFAULT_TIMEOUT = 100 * time.Millisecond

	JOB_RUNNING int64 = iota
	JOB_DONE
)

var (
	FlagIdleSimulateTx = "idle-simulate-tx"
)

type IdleDispatcher struct {
	enable bool
	chJobs chan func()
	isIdle bool
	cond   *sync.Cond

	chJobDone chan struct{}
	jobStatus int64
}

func NewIdleDispatcher() *IdleDispatcher {
	enableIdleSimulateTx := viper.GetBool(FlagIdleSimulateTx)
	idp := &IdleDispatcher{
		isIdle:    false,
		cond:      sync.NewCond(&sync.Mutex{}),
		chJobs:    make(chan func()),
		chJobDone: make(chan struct{}, 1),
		jobStatus: JOB_DONE,
		enable:    enableIdleSimulateTx,
	}
	//start
	idp.chJobDone <- struct{}{}

	return idp
}

func (idp *IdleDispatcher) EnterCriticalState() {
	if !idp.enable {
		return
	}
	idp.cond.L.Lock()
	defer idp.cond.L.Unlock()
	idp.isIdle = false
	idp.cond.Broadcast()

}

func (idp *IdleDispatcher) LeaveCriticalState() {
	if !idp.enable {
		return
	}
	idp.cond.L.Lock()
	defer idp.cond.L.Unlock()
	idp.isIdle = true
	idp.cond.Broadcast()

}

func (idp *IdleDispatcher) AddJob(job func()) {
	if !idp.enable {
		return
	}
	idp.chJobs <- job
}

func (idp *IdleDispatcher) JobDoneChan() chan struct{} {
	if !idp.enable {
		return nil
	}
	return idp.chJobDone
}

func (idp *IdleDispatcher) JobChan() chan func() {
	if !idp.enable {
		return nil
	}
	return idp.chJobs
}

func (idp *IdleDispatcher) IdleDo() {
	if !idp.enable {
		return
	}
	for {
		func() {
			idp.cond.L.Lock()
			defer idp.cond.L.Unlock()
			for !idp.isIdle {
				idp.cond.Wait()
			}

			if atomic.LoadInt64(&idp.jobStatus) == JOB_RUNNING {
				return
			}

			ctx, cancelFunc := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)

			select {
			case <-ctx.Done():
				go func() {
					defer cancelFunc()
				}()
				return
			case <-idp.JobDoneChan():
			}

			go func() {
				defer cancelFunc()

				job := <-idp.JobChan()
				atomic.StoreInt64(&idp.jobStatus, JOB_RUNNING)
				job()
				atomic.StoreInt64(&idp.jobStatus, JOB_DONE)
				idp.chJobDone <- struct{}{}
			}()
		}()
	}

}
