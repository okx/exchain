package types

import "sync"

var CommitStateDBPool = sync.Pool{
	New: func() interface{} {
		return new(CommitStateDB)
	},
}
