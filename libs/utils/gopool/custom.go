package gopool

import "github.com/panjf2000/ants/v2"

type CustomPoolConfig struct {
	Size int
}

type CustomPool struct {
	config CustomPoolConfig
	pool   *ants.PoolWithFunc
}

type Option ants.Option

// WithNonblocking indicates that pool will return nil when there is no available workers.
func WithNonblocking(nonblocking bool) Option {
	return func(opts *ants.Options) {
		opts.Nonblocking = nonblocking
	}
}

func NewPool(config CustomPoolConfig, fn func(interface{}), opts ...Option) (*CustomPool, error) {
	var antsOpts []ants.Option
	for _, v := range opts {
		antsOpts = append(antsOpts, ants.Option(v))
	}
	poolWithFunc, err := ants.NewPoolWithFunc(config.Size, fn, antsOpts...)
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
