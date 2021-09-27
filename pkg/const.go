package pkg

const (
	COMMIT_STATE_DB = "CommitStateDb"
	AVL             = "iavl"
	BLOCK_FORMAT    = "Block: Height<%d>, evmCost<%dns>, TxDetail: %s"
	TX_FORMAT = "Tx: %d,  evmCost<%dns>, "
	TX_DETAIL = "oper: %s, count<%d>, cost<%dns>, min<%dns>, max<%dns>, avg<%dns>"
)

var module = []string{COMMIT_STATE_DB, AVL}
