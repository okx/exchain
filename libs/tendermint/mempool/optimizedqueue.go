package mempool

func NewOptimizedTxQueue(txPriceBump int64) ITransactionQueue {
	return NewGasTxQueue(txPriceBump)
}
