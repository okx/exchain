package tikv

import (
	"context"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/tikv/client-go/v2/rawkv"
)

type Iterator struct {
	client      *rawkv.Client
	curKey      []byte
	curKeyInner []byte
	curValue    []byte
	start       []byte
	end         []byte
	isReverse   bool
	finish      bool
	err         error

	keys   [][]byte
	values [][]byte
	cur    int
}

var _ dbm.Iterator = (*Iterator)(nil)

func newIterator(start, end []byte, isReverse bool, client *rawkv.Client) *Iterator {
	return &Iterator{
		curKey:      start,
		curKeyInner: start,
		start:       start,
		end:         end,
		isReverse:   isReverse,
		client:      client,
	}
}

func (i *Iterator) Domain() (start []byte, end []byte) {
	return i.start, i.end
}

func (i *Iterator) Valid() bool {
	return i.client != nil && i.err == nil && !i.finish
}

func (i *Iterator) _next() ([]byte, []byte, error) {
	var err error
	i.keys, i.values, err = i.client.Scan(context.TODO(), i.curKey, i.end, 1000)
	if err != nil {
		return nil, nil, err
	}

	if len(i.keys) == 0 {
		i.finish = true
		return nil, nil, nil
	}
	if i.cur >= len(i.keys) {
		i.finish = true
		return nil, nil, nil
	}
	key, value := i.keys[i.cur], i.values[i.cur]
	i.cur++
	return key, value, nil
}

func (i *Iterator) next() ([]byte, []byte, error) {
	var err error
	i.keys, i.values, err = i.client.Scan(context.TODO(), i.curKeyInner, i.end, 2)
	if err != nil {
		return nil, nil, err
	}

	if len(i.keys) == 0 || i.curKeyInner == nil {
		i.finish = true
		return nil, nil, nil
	}
	if len(i.keys) == 1 {
		i.curKeyInner = nil
	}
	if len(i.keys) == 2 {
		i.curKeyInner = make([]byte, len(i.keys[1]))
		copy(i.curKeyInner, i.keys[1])
	}

	key, value := i.keys[0], i.values[0]
	return key, value, nil
}

func (i *Iterator) reverseNext() ([]byte, []byte, error) {
	var err error
	i.keys, i.values, err = i.client.ReverseScan(context.TODO(), i.curKey, i.end, 2)
	if err != nil {
		return nil, nil, err
	}

	if len(i.keys) == 0 || i.curKeyInner == nil {
		i.finish = true
		return nil, nil, nil
	}
	if len(i.keys) == 1 {
		i.curKeyInner = nil
	}
	if len(i.keys) == 2 {
		i.curKeyInner = make([]byte, len(i.keys[1]))
		copy(i.curKeyInner, i.keys[1])
	}

	key, value := i.keys[0], i.values[0]
	return key, value, nil
}

func (i *Iterator) Next() {
	if i.isReverse {
		i.curKey, i.curValue, i.err = i.reverseNext()
		return
	}
	i.curKey, i.curValue, i.err = i.next()
}

func (i *Iterator) Key() (key []byte) {
	return i.curKey
}

func (i *Iterator) Value() (value []byte) {
	return i.curValue
}

func (i *Iterator) Error() error {
	return i.err
}

func (i *Iterator) Close() {}
