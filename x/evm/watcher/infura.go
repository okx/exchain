package watcher

type InfuraKeeper interface {
	OnSaveTransactionReceipt(TransactionReceipt)
}
