package watcher

import (
	"github.com/sdming/goh"
	"github.com/sdming/goh/Hbase"
)

const (
	FlagFastQuery = "fast-query"
	TableName     = "infura-testnet"
	Column        = "Data:data"
)

type WatchStore struct {
	db *goh.HClient
}

var gWatchStore *WatchStore = nil

func InstanceOfWatchStore() *WatchStore {
	if gWatchStore == nil && IsWatcherEnabled() {
		db, e := initDb()
		if e == nil {
			gWatchStore = &WatchStore{db: db}
		}
	}
	return gWatchStore
}

func initDb() (*goh.HClient, error) {
	//todo getFrom config
	client, err := goh.NewTcpClient("127.0.0.1:9090", goh.TBinaryProtocol, false)
	if err != nil {
		panic(err)
	}
	if err = client.Open(); err != nil {
		panic(err)
	}
	return client, nil
}

func (w WatchStore) Set(key []byte, value []byte) {
	mutations := make([]*Hbase.Mutation, 1)
	mutations[0] = goh.NewMutation(Column, value)
	w.db.MutateRow(TableName, key, mutations, nil)
}

func (w WatchStore) Get(key []byte) ([]byte, error) {
	data, err := w.db.Get(TableName, key, Column, nil)
	if data == nil || len(data) == 0 {
		return nil, err
	}
	return data[0].Value, err
}
