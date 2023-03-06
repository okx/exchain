package watcher

import (
	"log"

	"github.com/golang/protobuf/proto"
	"github.com/okx/okbchain/x/wasm/types"
)

var (
	paramsKey      = []byte("wasm-parameters")
	sendEnabledKey = []byte("send-enabled")
)

func SetParams(para types.Params) {
	if !Enable() {
		return
	}
	b, err := proto.Marshal(&para)
	if err != nil {
		panic("wasm watchDB SetParams marshal error:" + err.Error())
	}
	if err = db.Set(paramsKey, b); err != nil {
		panic("wasm watchDB SetParams set error:" + err.Error())
	}
}

func GetParams() types.Params {
	ensureChecked()
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

type ParamsManager struct{}

func (p ParamsManager) SetSendEnabled(enable bool) {
	if !Enable() {
		return
	}
	var ok byte
	if enable {
		ok = 1
	}
	if err := db.Set(sendEnabledKey, []byte{ok}); err != nil {
		log.Println("SetSendEnabled error:", err)
	}
}

func (p ParamsManager) GetSendEnabled() bool {
	ensureChecked()
	v, err := db.Get(sendEnabledKey)
	if err != nil {
		log.Println("SetSendEnabled error:", err)
		return false
	}
	if len(v) == 0 || v[0] == 0 {
		return false
	}
	return true
}
