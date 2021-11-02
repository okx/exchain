package mempool

import (
	"github.com/okex/exchain/libs/tendermint/libs/clist"
	tmmath "github.com/okex/exchain/libs/tendermint/libs/math"
	"github.com/okex/exchain/libs/tendermint/types"
	"sync"
)

type AddressRecord struct {
	mtx    sync.RWMutex
	items map[string]map[string]*clist.CElement // Address -> (txHash -> *CElement)
}

func newAddressRecord() *AddressRecord {
	return &AddressRecord{
		items: make(map[string]map[string]*clist.CElement),
	}
}

func (ar *AddressRecord) AddItem(address string, txHash string, cElement *clist.CElement) {
	ar.mtx.Lock()
	defer ar.mtx.Unlock()
	if _, ok := ar.items[address]; !ok {
		ar.items[address] = make(map[string]*clist.CElement)
	}
	ar.items[address][txHash] = cElement
}

func (ar *AddressRecord) GetItem (address string) (item map[string]*clist.CElement, isExist bool) {
	ar.mtx.RLock()
	defer ar.mtx.RUnlock()
	item, isExist = ar.items[address]
	return
}

func (ar *AddressRecord) DeleteItem(e *clist.CElement) {
	ar.mtx.Lock()
	defer ar.mtx.Unlock()
	if userMap, ok := ar.items[e.Address]; ok {
		txHash := txID(e.Value.(*mempoolTx).tx)
		if _, ok = userMap[txHash]; ok {
			delete(userMap, txHash)
		}

		if len(userMap) == 0 {
			delete(ar.items, e.Address)
		}
	}
}

func (ar *AddressRecord) GetAddressList() []string {
	ar.mtx.RLock()
	defer ar.mtx.RUnlock()
	addressList := make([]string, 0, len(ar.items))
	for address, _ := range ar.items {
		addressList = append(addressList, address)
	}
	return addressList
}

func (ar *AddressRecord) GetAddressTxsCnt(address string) int {
	ar.mtx.RLock()
	defer ar.mtx.RUnlock()
	cnt := 0
	if userMap, ok := ar.items[address]; ok {
		cnt = len(userMap)
	}
	return cnt
}

func (ar *AddressRecord) GetAddressTxs(address string, txCount int, max int) types.Txs {
	ar.mtx.RLock()
	defer ar.mtx.RUnlock()
	userMap, ok := ar.items[address]
	if !ok || len(userMap) == 0 {
		return types.Txs{}
	}

	txNums := len(userMap)
	if max <= 0 || max > txNums {
		max = txNums
	}

	txs := make([]types.Tx, 0, tmmath.MinInt(txCount, max))

	for _, ele := range userMap {
		if len(txs) == max {
			break
		}

		txs = append(txs, ele.Value.(*mempoolTx).tx)
	}
	return txs
}
