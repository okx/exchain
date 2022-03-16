package watcher

import (
	"encoding/json"

	"github.com/tendermint/go-amino"
)

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
		if nil != err {
			panic("cant happen")
		}
		b.value = string(vs)
		b.origin = nil
	}
	return b.value
}

type baseAminoMarshal struct {
	baseLazyMarshal
}

func newBaseAminoMarshal(o interface{}) *baseAminoMarshal {
	return &baseAminoMarshal{
		baseLazyMarshal{
			origin: o,
		},
	}
}
func (b *baseAminoMarshal) GetValue() string {
	if b.origin != nil {
		vs, err := amino.MarshalBinaryBare(b.origin)
		if nil != err {
			panic("cant happen")
		}
		b.value = string(vs)
		b.origin = nil
	}
	return b.value
}
