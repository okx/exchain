package dydx

import (
	"container/list"
	"errors"
	"sync"
)

type DepthBook struct {
	buyOrders  *OrderList
	sellOrders *OrderList
}

func NewDepthBook() *DepthBook {
	return &DepthBook{
		buyOrders:  NewOrderList(true),
		sellOrders: NewOrderList(false),
	}
}

func (d *DepthBook) Insert(order *WrapOrder) error {
	if order.Type() == SellOrderType {
		d.sellOrders.Insert(order)
	} else if order.Type() == BuyOrderType {
		d.buyOrders.Insert(order)
	} else {
		return errors.New("invalid order")
	}
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

func (d *DepthBook) Delete(key [KeySize]byte) *list.Element {
	if ele := d.buyOrders.Get(key); ele != nil {
		return d.buyOrders.Remove(ele)
	}
	if ele := d.sellOrders.Get(key); ele != nil {
		return d.sellOrders.Remove(ele)
	}
	return nil
}

type OrderList struct {
	sync.RWMutex

	reverse bool
	orders  *list.List
	index   map[[KeySize]byte]*list.Element
}

func NewOrderList(reverse bool) *OrderList {
	return &OrderList{
		reverse: reverse,
		orders:  list.New(),
		index:   make(map[[KeySize]byte]*list.Element),
	}
}

func (o *OrderList) Front() *list.Element {
	o.RLock()
	defer o.RUnlock()
	return o.orders.Front()
}

func (o *OrderList) Get(key [KeySize]byte) *list.Element {
	o.RLock()
	defer o.RUnlock()
	return o.index[key]
}

func (o *OrderList) Insert(order *WrapOrder) *list.Element {
	o.Lock()
	defer o.Unlock()

	ele := o.orders.Front()
	for ele != nil {
		cur := ele.Value.(*WrapOrder)
		if o.less(order, cur) {
			newEle := o.orders.InsertBefore(order, ele)
			o.index[order.Key()] = newEle
			return newEle
		}
		ele = ele.Next()
	}
	newEle := o.orders.PushBack(order)
	o.index[order.Key()] = newEle
	return newEle
}

func (o *OrderList) Pop() *list.Element {
	o.Lock()
	defer o.Unlock()

	front := o.orders.Front()
	o.orders.Remove(front)
	delete(o.index, front.Value.(*WrapOrder).Key())
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
	if _, ok = o.index[order.Key()]; ok {
		o.orders.Remove(ele)
		delete(o.index, order.Key())
		return ele
	}
	return nil
}

func (o *OrderList) less(order1, order2 *WrapOrder) bool {
	if o.reverse {
		return order1.Amount.Cmp(order2.Amount) > 0
	}
	return order1.Amount.Cmp(order2.Amount) < 0
}
