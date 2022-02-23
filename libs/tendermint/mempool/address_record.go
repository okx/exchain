package mempool

import (
	"sync"

	"github.com/okex/exchain/libs/tendermint/libs/clist"
	tmmath "github.com/okex/exchain/libs/tendermint/libs/math"
	"github.com/okex/exchain/libs/tendermint/types"
)

type AddressRecord struct {
	addrTxs sync.Map // address -> map[txHash]*addrMap
}

type addrMap struct {
	sync.RWMutex

	items map[string]*clist.CElement  // txHash -> *mempoolTx
	maxNonce uint64
}

func newAddressRecord() *AddressRecord {
	return &AddressRecord{}
}

func (ar *AddressRecord) AddItem(address string, txHash string, cElement *clist.CElement) {
	v, ok := ar.addrTxs.Load(address)
	if !ok {
		// LoadOrStore to prevent double storing
		v, ok = ar.addrTxs.LoadOrStore(address, &addrMap{items:make(map[string]*clist.CElement)})
	}
	am := v.(*addrMap)
	am.Lock()
	defer am.Unlock()
	am.items[txHash] = cElement
	if cElement.Nonce > am.maxNonce {
		am.maxNonce = cElement.Nonce
	}
}

func (ar *AddressRecord) checkRepeatedElement(info ExTxInfo) *clist.CElement {
	v, ok := ar.addrTxs.Load(info.Sender)
	if !ok {
		return nil
	}
	am := v.(*addrMap)
	am.RLock()
	defer am.RUnlock()
	if info.Nonce > am.maxNonce {
		return nil
	}
	for _, ele := range am.items {
		if ele.Nonce == info.Nonce {
			return ele
		}
	}
	return nil
}

func (ar *AddressRecord) CleanItems(address string, nonce uint64) []*clist.CElement {
	v, ok := ar.addrTxs.Load(address)
	if !ok {
		return nil
	}
	am := v.(*addrMap)
	var l []*clist.CElement
	am.Lock()
	defer am.Unlock()
	for k, v := range am.items {
		if v.Nonce <= nonce {
			l = append(l, v)
			delete(am.items, k)
		}
	}
	if len(am.items) == 0 {
		ar.addrTxs.Delete(address)
	}
	return l
}

func (ar *AddressRecord) GetItems(address string) []*clist.CElement {
	v, ok := ar.addrTxs.Load(address)
	if !ok {
		return nil
	}
	am := v.(*addrMap)
	var l []*clist.CElement
	am.RLock()
	defer am.RUnlock()
	for _, v := range am.items {
		l = append(l, v)
	}
	return l
}

func (ar *AddressRecord) DeleteItem(e *clist.CElement) {
	if v, ok := ar.addrTxs.Load(e.Address); ok {
		am := v.(*addrMap)
		memTx := e.Value.(*mempoolTx)
		txHash := txID(memTx.tx, memTx.height)
		am.Lock()
		defer am.Unlock()
		delete(am.items, txHash)
		if len(am.items) == 0 {
			ar.addrTxs.Delete(e.Address)
		}
	}
}

func (ar *AddressRecord) GetAddressList() []string {
	var addrList []string
	ar.addrTxs.Range(func(k, v interface{}) bool {
		addrList = append(addrList, k.(string))
		return true
	})
	return addrList
}

func (ar *AddressRecord) GetAddressTxsCnt(address string) int {
	v, ok := ar.addrTxs.Load(address)
	if !ok {
		return 0
	}
	am := v.(*addrMap)
	am.RLock()
	defer am.RUnlock()
	return len(am.items)
}

func (ar *AddressRecord) GetAddressNonce(address string) uint64 {
	v, ok := ar.addrTxs.Load(address)
	if !ok {
		return 0
	}
	am := v.(*addrMap)
	am.RLock()
	defer am.RUnlock()
	var nonce uint64
	for _, e := range am.items {
		if e.Nonce > nonce {
			nonce = e.Nonce
		}
	}
	return nonce
}

func (ar *AddressRecord) GetAddressTxs(address string, txCount int, max int) types.Txs {
	v, ok := ar.addrTxs.Load(address)
	if !ok {
		return nil
	}
	am := v.(*addrMap)
	am.RLock()
	defer am.RUnlock()
	if max <= 0 || max > len(am.items) {
		max = len(am.items)
	}
	txs := make([]types.Tx, 0, tmmath.MinInt(txCount, max))
	for _, e := range am.items {
		if len(txs) == cap(txs) {
			break
		}
		txs = append(txs, e.Value.(*mempoolTx).tx)
	}
	return txs
}

