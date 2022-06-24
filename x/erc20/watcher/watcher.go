package watcher

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	clientcontext "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/x/erc20/types"
	"github.com/spf13/viper"
	"strings"
	"sync"
	"sync/atomic"
)

type Watcher struct {
	store         *WatchStore
	sw            bool
	isLoadHistory int32
}

var (
	watcherEnable = false
	onceEnable    sync.Once
	isLoadHistory = int32(0)
)

func IsWatcherEnabled() bool {
	onceEnable.Do(func() {
		watcherEnable = viper.GetBool(FlagFastQuery)
	})
	return watcherEnable
}

func NewWatcher() *Watcher {
	return &Watcher{
		store: InstanceOfWatchStore(),
		sw:    IsWatcherEnabled(),
	}
}

func (w *Watcher) LoadItem(denom, contract string) bool {
	if strings.HasPrefix(contract, "0x") {
		contractBytes, err := hex.DecodeString(contract[2:])
		if err != nil {
			panic(fmt.Errorf("erc20  LoadTokenMapping error", err))
		}
		w.SetContractToDenom(contractBytes, []byte(denom))
	}
	return false
}

func (w *Watcher) LoadTokenMapping(f func(func(denom, contract string) bool)) {
	if atomic.CompareAndSwapInt32(&isLoadHistory, 0, 1) {
		f(w.LoadItem)
	}
}

func (w *Watcher) LoadHistoryTokenMapping(ctx clientcontext.CLIContext) error {
	if atomic.CompareAndSwapInt32(&isLoadHistory, 0, 1) {
		queryPath := fmt.Sprintf("custom/erc20/token-mapping")
		res, _, err := ctx.QueryWithData(queryPath, []byte{})
		if err != nil {
			return err
		}
		var mapping []types.TokenMapping
		err = json.Unmarshal(res, &mapping)
		if err != nil {
			return err
		}
		for _, v := range mapping {
			//rm 0x prefix
			contract, err := hex.DecodeString(v.Contract[2:])
			if err != nil {
				return err
			}
			w.SetContractToDenom(contract, []byte(v.Denom))
		}
	}
	return nil
}

func (w *Watcher) SetContractToDenom(key, value []byte) {
	if w.Enabled() {
		w.store.Set(types.ContractToDenomKey(key), value)
	}
}

func (w *Watcher) GetDenomByContract(key []byte) []byte {

	if w.Enabled() {
		r, err := w.store.Get(types.ContractToDenomKey(key))
		if err != nil {
			return nil
		}
		return r
	}
	return nil
}

func (w *Watcher) Delete(key []byte) {
	if w.Enabled() {
		w.store.Delete(types.ContractToDenomKey(key))
	}
}

func (w *Watcher) Enabled() bool {
	return w.sw
}
