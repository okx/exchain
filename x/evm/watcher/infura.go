package watcher

type InfuraKeeper interface {
	OnSaveTransactionReceipt(TransactionReceipt)
	OnSaveBlock(Block)
	OnSaveTransaction(Transaction)
}
