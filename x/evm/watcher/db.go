package watcher

import (
	"context"
	"fmt"
	"github.com/silenceper/pool"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"time"
)

const (
	FlagFastQuery = "fast-query"
	TableName     = "infura-testnet"
)

type WatchStore struct {
	db pool.Pool
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

func initDb() (pool.Pool, error) {
	//todo getFrom config
	factory := func() (interface{}, error) { return gohbase.NewClient("18.167.164.175:21811"), nil }
	close := func(v interface{}) error { v.(gohbase.Client).Close(); return nil }

	//pool config
	poolConfig := &pool.Config{
		InitialCap:  5,
		MaxIdle:     200,
		MaxCap:      2000,
		Factory:     factory,
		Close:       close,
		IdleTimeout: 15 * time.Second,
	}
	return pool.NewChannelPool(poolConfig)
}

func (w WatchStore) Set(key []byte, value []byte) {
	v, err := w.db.Get()
	if v == nil || err != nil {
		fmt.Println("db get error:", err)
	}
	client, _ := v.(gohbase.Client)
	defer w.db.Put(client)

	values := map[string]map[string][]byte{
		"Data": map[string][]byte{
			"data": value,
		}}

	mutate, err := hrpc.NewPut(context.Background(), []byte(TableName), key, values)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = client.Put(mutate)
	if err != nil {
		fmt.Println(err)
	}
}

func (w WatchStore) Get(key []byte) ([]byte, error) {
	v, err := w.db.Get()
	if v == nil || err != nil {
		fmt.Println("db get error:", err)
		return []byte{}, err
	}
	client, _ := v.(gohbase.Client)
	defer w.db.Put(client)

	getRequest, _ := hrpc.NewGet(context.Background(), []byte(TableName), key)
	result, err := client.Get(getRequest)
	if err != nil || len(result.Cells) == 0 {
		return []byte{}, err
	}
	return result.Cells[0].Value, err
}
