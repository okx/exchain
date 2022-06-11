package gopool

import "github.com/panjf2000/ants/v2"

type CustomPoolConfig struct {
	Size int
}

type CustomPool struct {
	config CustomPoolConfig
	pool   *ants.PoolWithFunc
}

func NewPool(config CustomPoolConfig, fn func(interface{})) (*CustomPool, error) {
	poolWithFunc, err := ants.NewPoolWithFunc(config.Size, fn)
	if err != nil {
		return nil, err
	}
	return &CustomPool{
		config: config,
		pool:   poolWithFunc,
	}, nil
}

func (p *CustomPool) Invoke(args interface{}) error {
	return p.pool.Invoke(args)
}

func (p *CustomPool) Release() {
	p.pool.Release()
}
