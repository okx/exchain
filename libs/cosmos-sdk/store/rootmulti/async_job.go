package rootmulti

import "sync"

const (
	MaxAsyncJob = 100
)

func (rs *Store) EnableAsyncJob() {
	rs.enableAsyncJob = true
	rs.cacheMetadata()
	rs.lazySetupAsyncJob()
	go rs.jobRoutine()
}

func (rs *Store) dispatchJob(job func()) {
	rs.jobChan <- job
}

func (rs *Store) lazySetupAsyncJob() {
	rs.jobChan = make(chan func(), MaxAsyncJob)
	rs.jobDone = new(sync.WaitGroup)
	rs.jobDone.Add(1)
}

func (rs *Store) jobRoutine() {
	for job := range rs.jobChan {
		job()
	}
	rs.jobDone.Done()
}

func (rs *Store) stopJob() {
	if !rs.enableAsyncJob {
		return
	}
	close(rs.jobChan)
	rs.jobDone.Wait()
}
