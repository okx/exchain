package dydx

import (
	"container/list"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type DepthBook struct {
	buyOrders  *OrderList
	sellOrders *OrderList

	addrOrders map[common.Address][]*WrapOrder
	addrMtx    sync.RWMutex
}

func NewDepthBook() *DepthBook {
	return &DepthBook{
		buyOrders:  NewOrderList(true),
		sellOrders: NewOrderList(false),
		addrOrders: make(map[common.Address][]*WrapOrder),
	}
}

func (d *DepthBook) Insert(order *WrapOrder) error {
	var ele *list.Element
	if order.isBuy() {
		ele = d.buyOrders.Insert(order)
	} else {
		ele = d.sellOrders.Insert(order)
	}
	if ele == nil {
		return ErrRepeatedOrder
	}
	d.addrMtx.Lock()
	d.addrOrders[order.Maker] = append(d.addrOrders[order.Maker], order)
	d.addrMtx.Unlock()

	return nil
}

func (d *DepthBook) SellFront() *WrapOrder {
	ele := d.sellOrders.Front()
	wodr, ok := ele.Value.(*WrapOrder)
	if !ok {
		//TODO log err
		return nil
	}
	return wodr
}

func (d *DepthBook) Delete(order *WrapOrder) *list.Element {
	var ele *list.Element
	if order.Flags[31] == 1 {
		if ele = d.buyOrders.Get(order.Hash()); ele != nil {
			d.buyOrders.Remove(ele)
		}
	} else {
		if ele = d.sellOrders.Get(order.Hash()); ele != nil {
			d.sellOrders.Remove(ele)
		}
	}
	if ele != nil {
		maker := ele.Value.(*WrapOrder).Maker
		hash := ele.Value.(*WrapOrder).Hash()
		orders := d.addrOrders[maker]
		for i, order := range orders {
			if order.Hash() == hash {
				d.addrOrders[maker] = append(orders[:i], orders[i+1:]...)
			}
		}
	}
	return ele

}

func (d *DepthBook) DeleteByHash(hash common.Hash) *list.Element {
	var ele *list.Element
	if ele == nil {
		if ele = d.buyOrders.Get(hash); ele != nil {
			d.buyOrders.Remove(ele)
		}
	}
	if ele == nil {
		if ele = d.sellOrders.Get(hash); ele != nil {
			d.sellOrders.Remove(ele)
		}
	}
	if ele != nil {
		maker := ele.Value.(*WrapOrder).Maker
		orders := d.addrOrders[maker]
		for i, order := range orders {
			if order.Hash() == hash {
				d.addrOrders[maker] = append(orders[:i], orders[i+1:]...)
			}
		}
	}
	return ele
}

type OrderList struct {
	sync.RWMutex

	reverse bool
	orders  *list.List
	index   map[common.Hash]*list.Element
}

func NewOrderList(reverse bool) *OrderList {
	ol := &OrderList{
		reverse: reverse,
		orders:  list.New(),
		index:   make(map[common.Hash]*list.Element),
	}
	go ol.prune()
	return ol
}

func (o *OrderList) Front() *list.Element {
	o.RLock()
	defer o.RUnlock()
	return o.orders.Front()
}

func (o *OrderList) Get(hash common.Hash) *list.Element {
	o.RLock()
	defer o.RUnlock()
	return o.index[hash]
}

func (o *OrderList) Insert(order *WrapOrder) *list.Element {
	o.Lock()
	defer o.Unlock()

	if _, ok := o.index[order.Hash()]; ok {
		return nil
	}

	ele := o.orders.Front()
	for ele != nil {
		cur := ele.Value.(*WrapOrder)
		if o.less(order, cur) {
			newEle := o.orders.InsertBefore(order, ele)
			o.index[order.Hash()] = newEle
			return newEle
		}
		ele = ele.Next()
	}
	newEle := o.orders.PushBack(order)
	o.index[order.Hash()] = newEle
	return newEle
}

func (o *OrderList) Pop() *list.Element {
	o.Lock()
	defer o.Unlock()

	front := o.orders.Front()
	o.orders.Remove(front)
	delete(o.index, front.Value.(*WrapOrder).Hash())
	return front
}

func (o *OrderList) Remove(ele *list.Element) *list.Element {
	o.Lock()
	defer o.Unlock()
	order, ok := ele.Value.(*WrapOrder)
	if !ok {
		//TODO: log error
		return nil
	}
	if _, ok = o.index[order.Hash()]; ok {
		o.orders.Remove(ele)
		delete(o.index, order.Hash())
		return ele
	}
	return nil
}

func (o *OrderList) List() []*WrapOrder {
	o.Lock()
	defer o.Unlock()
	var orders []*WrapOrder
	for ele := o.orders.Front(); ele != nil; ele = ele.Next() {
		orders = append(orders, ele.Value.(*WrapOrder))
	}
	return orders
}

func (o *OrderList) Len() int {
	o.RLock()
	defer o.RUnlock()
	return len(o.index)
}

// TODO, use block.timestamp?
func (o *OrderList) prune() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			o.Lock()
			for ele := o.orders.Front(); ele != nil; ele = ele.Next() {
				if ele.Value.(*WrapOrder).Expiration.Uint64() < uint64(time.Now().Unix()) {
					o.orders.Remove(ele)
				}
			}
			o.Unlock()
		}
	}
}

func (o *OrderList) less(order1, order2 *WrapOrder) bool {
	if o.reverse {
		return order1.Price().Cmp(order2.Price()) > 0
	}
	return order1.Price().Cmp(order2.Price()) < 0
}
