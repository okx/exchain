package watcher

import (
	"encoding/json"
	"fmt"
	"sync"

	cryptocodec "github.com/okex/exchain/app/crypto/ethsecp256k1"
	apptypes "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/tendermint/go-amino"
)

var (
	watcherInitCdcOnce sync.Once
	watcherCdc         *amino.Codec
)

type LazyValueMarshaler interface {
	GetOrigin() interface{}
	GetValue() string
}

type BaseMarshaler struct {
	origin interface{}
	value  *string
}

func (b *BaseMarshaler) GetOrigin() interface{} {
	return b.origin
}

type JsonMarshaler struct {
	BaseMarshaler
}

func newJsonMarshaller(o interface{}) *JsonMarshaler {
	return &JsonMarshaler{
		BaseMarshaler{
			origin: o,
			value:  nil, //value will be set when GetValue() is called
		},
	}
}
func (b *JsonMarshaler) GetValue() (ret string) {
	if b.value != nil {
		vs, err := json.Marshal(b.origin)
		if nil != err {
			panic(fmt.Sprintf("fail to marshaled by json, err : %s", err.Error()))
		}
		ret = string(vs)
		b.value = &ret
	}
	return
}

type AminoMarshaler struct {
	BaseMarshaler
}

func InitWatcherCdc() {
	watcherInitCdcOnce.Do(func() {
		watcherCdc = codec.New()
		apptypes.RegisterCodec(watcherCdc)
		cryptocodec.RegisterCodec(watcherCdc)
		codec.RegisterCrypto(watcherCdc)
	})
}
func WatcherCdc() *amino.Codec {
	InitWatcherCdc()
	return watcherCdc
}
func newAminoMarshaller(o interface{}) *AminoMarshaler {
	InitWatcherCdc()
	return &AminoMarshaler{
		BaseMarshaler{
			origin: o,
			value:  nil, //value will be set when GetValue() is called
		},
	}
}
func (b *AminoMarshaler) GetValue() (ret string) {
	if b.value != nil {
		vs, err := watcherCdc.MarshalBinaryBare(b.origin)
		if nil != err {
			panic(fmt.Sprintf("fail to marshaled by amino, err : %s", err.Error()))
		}
		ret = string(vs)
		b.value = &ret
	}
	return
}
