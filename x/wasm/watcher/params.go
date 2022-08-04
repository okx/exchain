package watcher

import (
	"sync"

	"github.com/okex/exchain/x/wasm/types"
)

var params types.Params
var psMtx sync.RWMutex

func SetParams(para types.Params) {
	psMtx.Lock()
	params = para
	psMtx.Unlock()

}

func GetParams() types.Params {
	psMtx.RLock()
	defer psMtx.RUnlock()
	return params
}
