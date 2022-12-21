package tikv

import (
	"context"

	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/tikv/client-go/v2/rawkv"
)

type Batch struct {
	client *rawkv.Client
}

var _ dbm.Batch = (*Batch)(nil)

func newBatch(client *rawkv.Client) *Batch {
	return &Batch{
		client: client,
	}
}

func (b *Batch) Set(key, value []byte) {
	b.client.BatchPut(context.TODO(), [][]byte{key}, [][]byte{value})
}

func (b *Batch) Delete(key []byte) {
	b.client.BatchDelete(context.TODO(), [][]byte{key})
}

func (b *Batch) Write() error {
	return nil
}

func (b *Batch) WriteSync() error {
	return nil
}

func (b *Batch) Close() {
}
