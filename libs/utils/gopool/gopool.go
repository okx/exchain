package gopool

import "github.com/panjf2000/ants/v2"

const (
	maxPoolSize = 20000
)

var (
	gPool *ants.Pool
)

func init() {
	var err error
	gPool, err = ants.NewPool(maxPoolSize)
	if err != nil {
		panic(err)
	}
}

func Run(f func()) error {
	return gPool.Submit(f)
}

func Close() {
	gPool.Release()
}
