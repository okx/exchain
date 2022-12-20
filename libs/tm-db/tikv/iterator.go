package tikv

import (
	"context"

	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/tikv/client-go/v2/rawkv"
)

type Iterator struct {
	client    *rawkv.Client
	curKey    []byte
	curValue  []byte
	start     []byte
	end       []byte
	isReverse bool
	finish    bool
	err       error
}

var _ dbm.Iterator = (*Iterator)(nil)

func newIterator(start, end []byte, isReverse bool, client *rawkv.Client) *Iterator {
	return &Iterator{
		curKey:    start,
		start:     start,
		end:       end,
		isReverse: isReverse,
		client:    client,
	}
}

func (i *Iterator) Domain() (start []byte, end []byte) {
	return i.start, i.end
}

func (i *Iterator) Valid() bool {
	return i.client != nil && i.err == nil && !i.finish
}

func (i *Iterator) next() ([]byte, []byte, error) {
	keys, values, err := i.client.Scan(context.TODO(), i.curKey, i.end, 1)
	if err != nil {
		return nil, nil, err
	}
	if len(keys[0]) == 0 {
		i.finish = true
	}
	return keys[0], values[0], nil
}

func (i *Iterator) reverseNext() ([]byte, []byte, error) {
	keys, values, err := i.client.ReverseScan(context.TODO(), i.curKey, i.end, 1)
	if err != nil {
		return nil, nil, err
	}
	if len(keys[0]) == 0 {
		i.finish = true
	}
	return keys[0], values[0], nil
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
