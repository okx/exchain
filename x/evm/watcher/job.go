package watcher

type job interface{}

type watchDataCommitJob struct {
	watchData WatchData
}

type commitBatchJob struct {
	batch []WatchMessage
}