package watcher

import "encoding/json"

var (
	_ LazyValueMarshaler = (*baseLazyMarshal)(nil)
)

type LazyValueMarshaler interface {
	GetValue() string
}

type baseLazyMarshal struct {
	origin interface{}
	value  string
}

func newBaseLazyMarshal(o interface{}) *baseLazyMarshal {
	return &baseLazyMarshal{
		origin: o,
	}
}

func (b *baseLazyMarshal) GetValue() string {
	if b.origin != nil {
		vs, err := json.Marshal(b.origin)
		if err != nil {
			panic("cant happen")
		}
		b.value = string(vs)
		b.origin = nil
	}
	return b.value
}
