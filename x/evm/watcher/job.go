package watcher

type job func()

type watchDataCommitJob struct {
	watchData WatchData
}

type commitBatchJob struct {
	batch []WatchMessage
}