package common

import (
	"errors"
	"sort"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
)

var (
	errDuplicateVersionSelector = errors.New("duplicate version ")
	errRepeatRegisterSelector   = errors.New("repeat register selectors ")
)

var (
	DefaultFactory = func(higherThan func(h int64) bool, version SelectVersion, internal interface{}) SelectorFactory {
		return func() Selector {
			return NewCommonHeightSelector(higherThan, version, internal)
		}
	}
)

// TODO,use genesis
type Selector interface {
	Version() SelectVersion
	Select(ctx sdk.Context) (interface{}, bool)
}

type SelectVersion float64

type SelectorFactory func() Selector

type Selectors []Selector

func (m Selectors) Len() int {
	return len(m)
}

func (m Selectors) Less(i, j int) bool {
	return m[i].Version() >= m[j].Version()
}

func (m Selectors) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

////////
var (
	_ Selector = (*CommonHeightSelector)(nil)
)

type CommonHeightSelector struct {
	higherThan func(h int64) bool
	version    SelectVersion
	internal   interface{}
}

func NewCommonHeightSelector(higherThan func(h int64) bool, version SelectVersion, internal interface{}) *CommonHeightSelector {
	return &CommonHeightSelector{higherThan: higherThan, version: version, internal: internal}
}

func (c *CommonHeightSelector) Version() SelectVersion {
	return c.version
}

func (c *CommonHeightSelector) Select(ctx sdk.Context) (interface{}, bool) {
	if c.higherThan(ctx.BlockHeight()) {
		return c.internal, true
	}
	return nil, false
}

///////

type SelectorStrategy struct {
	internal  interface{}
	selectors Selectors
	seal      bool
}

func (f *SelectorStrategy) RegisterSelectors(factories ...SelectorFactory) {
	if f.Seald() {
		panic(errDuplicateVersionSelector)
	}
	var selectors Selectors
	set := make(map[SelectVersion]struct{})

	for _, f := range factories {
		sel := f()
		selectors = append(selectors, sel)
		v := sel.Version()
		if _, contains := set[v]; contains {
			panic(errDuplicateVersionSelector)
		}
		set[v] = struct{}{}
	}
	sort.Sort(selectors)

	f.selectors = selectors
}

func NewSelectorStrategy(internal interface{}) *SelectorStrategy {
	return &SelectorStrategy{internal: internal}
}

func (f *SelectorStrategy) Seald() bool {
	return f.seal
}
func (f *SelectorStrategy) Seal() {
	f.seal = true
}

func (f *SelectorStrategy) GetProxy(ctx sdk.Context) interface{} {
	for _, s := range f.selectors {
		m, ok := s.Select(ctx)
		if ok {
			return m
		}
	}
	return f.internal
}

func (f *SelectorStrategy) GetInternal() interface{} {
	return f.internal
}
