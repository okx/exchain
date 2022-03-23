package watcher

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/go-amino"
)

var (
	watcherInitCdcOnce sync.Once
	watcherCdc         *amino.Codec
)

type LazyValueMarshaler interface {
	GetValue() string
}

type BaseMarshaler struct {
	origin interface{}
	value  string
}
type JsonMarshaler struct {
	BaseMarshaler
}

func newJsonMarshaller(o interface{}) *JsonMarshaler {
	return &JsonMarshaler{
		BaseMarshaler{
			origin: o,
			value:  "", //value will be set when GetValue() is called
		},
	}
}
func (b *JsonMarshaler) GetValue() string {
	if b.origin != nil {
		vs, err := json.Marshal(b.origin)
		if nil != err {
			panic(fmt.Sprintf("fail to marshaled by json, err : %s", err.Error()))
		}
		b.value = string(vs)
		b.origin = nil
	}
	return b.value
}

type AminoMarshaler struct {
	BaseMarshaler
}

func InitWatcherCdc(cdc *amino.Codec) {
	watcherCdc = cdc
	watcherCdc.RegisterInterface((*interface{})(nil), nil)
	watcherCdc.RegisterConcrete([]*Transaction{}, "watcher/Transaction", nil)
	watcherCdc.RegisterConcrete([]common.Hash{}, "common/hash", nil)
}

func WatcherCdc() *amino.Codec {
	return watcherCdc
}
func newAminoMarshaller(o interface{}) *AminoMarshaler {
	return &AminoMarshaler{
		BaseMarshaler{
			origin: o,
			value:  "", //value will be set when GetValue() is called
		},
	}
}
func (b *AminoMarshaler) GetValue() string {
	if b.origin != nil {
		vs, err := watcherCdc.MarshalBinaryBare(b.origin)
		if nil != err {
			panic(fmt.Sprintf("fail to marshaled by amino, origin : %s, err : %s", reflect.TypeOf(b.origin), err.Error()))
		}
		b.value = string(vs)
		b.origin = nil
	}
	return b.value
}
