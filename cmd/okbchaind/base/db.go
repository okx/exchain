package base

import (
	"fmt"
	"strings"

	dbm "github.com/okx/okbchain/libs/tm-db"
)

const (
	AppDBName = "application.db"
)

func OpenDB(dir string, backend dbm.BackendType) (db dbm.DB, err error) {
	switch {
	case strings.HasSuffix(dir, ".db"):
		dir = dir[:len(dir)-3]
	case strings.HasSuffix(dir, ".db/"):
		dir = dir[:len(dir)-4]
	default:
		return nil, fmt.Errorf("database directory must end with .db")
	}
	//doesn't work on windows!
	cut := strings.LastIndex(dir, "/")
	if cut == -1 {
		return nil, fmt.Errorf("cannot cut paths on %s", dir)
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("couldn't create db: %v", r)
		}
	}()
	name := dir[cut+1:]
	db = dbm.NewDB(name, backend, dir[:cut])
	return db, nil
}
