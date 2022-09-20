package types

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

var (
	DefaultFactory = func(higherThan func(h int64) bool, version int, internal Middleware) SelectorFactory {
		return func() MiddlewareSelector {
			return NewDefaultKeepSelector(higherThan, version, internal)
		}
	}
)

// TODO,抽到公共中取
type MiddlewareSelector interface {
	Version() int
	Select(ctx sdk.Context) (Middleware, bool)
}

type MiddlewareSelectors []MiddlewareSelector

func (m MiddlewareSelectors) Len() int {
	return len(m)
}

func (m MiddlewareSelectors) Less(i, j int) bool {
	return m[i].Version() >= m[j].Version()
}

func (m MiddlewareSelectors) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type SelectorFactory func() MiddlewareSelector

var (
	_ MiddlewareSelector = (*DefaultKeepSelector)(nil)
)

type DefaultKeepSelector struct {
	higherThan func(h int64) bool
	version    int
	internal   Middleware
}

func (d *DefaultKeepSelector) Version() int {
	return d.version
}

func NewDefaultKeepSelector(higherThan func(h int64) bool, version int, internal Middleware) *DefaultKeepSelector {
	ret := &DefaultKeepSelector{
		higherThan: higherThan,
		version:    version,
		internal:   internal,
	}
	return ret
}

func (d *DefaultKeepSelector) Select(ctx sdk.Context) (Middleware, bool) {
	if d.higherThan(ctx.BlockHeight()) {
		return d.internal, true
	}
	return nil, false
}
