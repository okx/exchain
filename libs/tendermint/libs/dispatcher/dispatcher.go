package dispatcher

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

const (
	DEFAULT_TIMEOUT = 100 * time.Millisecond

	JOB_RUNNING int64 = iota
	JOB_DONE
)

type IdleDispatcher struct {
	chJobs chan func()
	isIdle bool
	cond   *sync.Cond

	chJobDone chan struct{}
	jobStatus int64
}

func NewIdleDispatcher() *IdleDispatcher {

	idp := &IdleDispatcher{
		isIdle:    false,
		cond:      sync.NewCond(&sync.Mutex{}),
		chJobs:    make(chan func()),
		chJobDone: make(chan struct{}, 1),
		jobStatus: JOB_DONE,
	}
	//start
	idp.chJobDone <- struct{}{}

	return idp
}

func (idp *IdleDispatcher) EnterCriticalState() {
	idp.cond.L.Lock()
	defer idp.cond.L.Unlock()
	idp.isIdle = false
	idp.cond.Broadcast()
}

func (idp *IdleDispatcher) LeaveCriticalState() {
	idp.cond.L.Lock()
	defer idp.cond.L.Unlock()
	idp.isIdle = true
	idp.cond.Broadcast()
}

func (idp *IdleDispatcher) AddJob(job func()) {
	idp.chJobs <- job
}

func (idp *IdleDispatcher) JobDoneChan() chan struct{} {
	return idp.chJobDone
}

func (idp *IdleDispatcher) JobChan() chan func() {
	return idp.chJobs
}

func (idp *IdleDispatcher) IdleDo() {
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
