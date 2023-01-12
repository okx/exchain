//go:build rocksdb
// +build rocksdb

package tikv

import (
	"fmt"
	"github.com/okex/exchain/libs/iavl"
	dbm "github.com/okex/exchain/libs/tm-db"
	"testing"

	"github.com/tikv/client-go/v2/rawkv"
)

func TestIterator_Valid(t *testing.T) {
	type fields struct {
		client    *rawkv.Client
		curKey    []byte
		curValue  []byte
		start     []byte
		end       []byte
		isReverse bool
		err       error
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"err is nil", fields{client: new(rawkv.Client)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				client:    tt.fields.client,
				curKey:    tt.fields.curKey,
				curValue:  tt.fields.curValue,
				start:     tt.fields.start,
				end:       tt.fields.end,
				isReverse: tt.fields.isReverse,
				err:       tt.fields.err,
			}
			if got := i.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadVersion(t *testing.T) {
	db, err := NewTiKV("", "127.0.0.1:2379")
	//db, err := dbm.NewRocksDB("application", "/Users/oker/workspace/exchain/dev/s0/data/")
	if err != nil {
		t.Fatal(err)
	}

	prefix := []byte(fmt.Sprintf("s/k:%s/", "evm"))

	prefixDB := dbm.NewPrefixDB(db, prefix)

	tree, err := iavl.NewMutableTree(prefixDB, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tree.LoadVersion(0)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(tree.IsUpgradeable())
}
