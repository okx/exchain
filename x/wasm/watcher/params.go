package watcher

import (
	"github.com/golang/protobuf/proto"
	"github.com/okex/exchain/x/wasm/types"
)

var (
	paramsKey = []byte("wasm-parameters")
)

func SetParams(para types.Params) {
	b, err := proto.Marshal(&para)
	if err != nil {
		panic("wasm watchDB SetParams marshal error:" + err.Error())
	}
	if err = db.Set(paramsKey, b); err != nil {
		panic("wasm watchDB SetParams set error:" + err.Error())
	}
}

func GetParams() types.Params {
	b, err := db.Get(paramsKey)
	if err != nil {
		panic("wasm watchDB GetParams get error:" + err.Error())
	}
	var p types.Params
	if err = proto.Unmarshal(b, &p); err != nil {
		panic("wasm watchDB GetParams unmarshal error:" + err.Error())
	}
	return p

}
