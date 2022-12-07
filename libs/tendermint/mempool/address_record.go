package mempool

import (
	"sync"

	"github.com/okex/exchain/libs/tendermint/libs/clist"
	"github.com/okex/exchain/libs/tendermint/types"
)

type elementManager interface {
	removeElement(*clist.CElement)
	removeElementByKey(key [32]byte) *clist.CElement
	reorganizeElements([]*clist.CElement)
}

type AddressRecord struct {
	addrTxs sync.Map // address -> *addrMap

	elementManager
}

type addrMap struct {
	sync.RWMutex

	items    map[uint64]*clist.CElement // nonce -> *mempoolTx
	maxNonce uint64
}

func newAddressRecord(em elementManager) *AddressRecord {
	return &AddressRecord{elementManager: em}
}

func (ar *AddressRecord) AddItem(address string, cElement *clist.CElement) {
	v, ok := ar.addrTxs.Load(address)
	if !ok {
		// LoadOrStore to prevent double storing
		v, ok = ar.addrTxs.LoadOrStore(address, &addrMap{items: make(map[uint64]*clist.CElement)})
	}
	am := v.(*addrMap)
	am.Lock()
	defer am.Unlock()
	am.items[cElement.Nonce] = cElement
	if cElement.Nonce > am.maxNonce {
		am.maxNonce = cElement.Nonce
	}
}

func (ar *AddressRecord) checkRepeatedAndAddItem(memTx *mempoolTx, txPriceBump int64, cb func(*clist.CElement) *clist.CElement) *clist.CElement {
	gasPrice := memTx.realTx.GetGasPrice()
	nonce := memTx.realTx.GetNonce()
	newElement := clist.NewCElement(memTx, memTx.from, gasPrice, nonce)

	v, ok := ar.addrTxs.Load(memTx.from)
	if !ok {
		v, ok = ar.addrTxs.LoadOrStore(memTx.from, &addrMap{items: make(map[uint64]*clist.CElement)})
	}
	am := v.(*addrMap)
	am.Lock()
	defer am.Unlock()
	// do not need to check element nonce
	if newElement.Nonce > am.maxNonce {
		cb(newElement)
		am.maxNonce = newElement.Nonce
		am.items[newElement.Nonce] = newElement
		return newElement
	}

	for _, e := range am.items {
		if e.Nonce == nonce {
			// only replace tx for bigger gas price
			expectedGasPrice := MultiPriceBump(e.GasPrice, txPriceBump)
			if gasPrice.Cmp(expectedGasPrice) <= 0 {
				return nil
			}

			// delete the old element and reorganize the elements whose nonce is greater the the new element
			ar.removeElement(e)
			items := []*clist.CElement{newElement}
			for _, item := range am.items {
				if item.Nonce > nonce {
					items = append(items, item)
				}
			}
			ar.reorganizeElements(items)
			am.items[newElement.Nonce] = newElement
			return newElement
		}
	}

	cb(newElement)
	am.items[newElement.Nonce] = newElement

	return newElement
}

func (ar *AddressRecord) CleanItems(address string, nonce uint64, cb func(element *clist.CElement)) {
	v, ok := ar.addrTxs.Load(address)
	if !ok {
		return
	}
	am := v.(*addrMap)
	am.Lock()
	defer am.Unlock()
	for k, v := range am.items {
		if v.Nonce <= nonce {
			cb(v)
			delete(am.items, k)
		}
	}
	if len(am.items) == 0 {
		ar.addrTxs.Delete(address)
	}
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
		am.Lock()
		defer am.Unlock()
		delete(am.items, e.Nonce)
		//calculate max Nonce
		am.maxNonce = calculateMaxNonce(am)
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

func (ar *AddressRecord) GetAddressNonce(address string) (uint64, bool) {
	v, ok := ar.addrTxs.Load(address)
	if !ok {
		return 0, false
	}
	am := v.(*addrMap)
	am.RLock()
	defer am.RUnlock()
	return am.maxNonce, true
}

func (ar *AddressRecord) GetAddressTxs(address string, max int) types.Txs {
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
	txs := make([]types.Tx, 0, max)
	for _, e := range am.items {
		if len(txs) == max {
			break
		}
		txs = append(txs, e.Value.(*mempoolTx).tx)
	}
	return txs
}

func calculateMaxNonce(data *addrMap) uint64 {
	maxNonce := uint64(0)
	for k, _ := range data.items {
		if k > maxNonce {
			maxNonce = k
		}
	}
	return maxNonce
}

type AddressNonce struct {
	addr  string
	nonce uint64
}

var addressNoncePool = sync.Pool{
	New: func() interface{} {
		return &AddressNonce{}
	},
}
