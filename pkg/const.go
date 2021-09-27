package pkg

const (
	COMMIT_STATE_DB = "CommitStateDb"
	AVL             = "iavl"
	DEBUG_FORMAT    = "Block: Height<%d>, " +
		"evmCost<%dms>, " +
		"%s"
)

var module = []string{COMMIT_STATE_DB, AVL}
