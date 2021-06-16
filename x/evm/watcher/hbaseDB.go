package watcher

import (
	"context"
	"fmt"
	"github.com/silenceper/pool"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"time"
)

type HbaseDB struct {
	db pool.Pool
}

const (
	TableName = "infura-testnet"
)

func initHbaseDB(dbUrl string) *HbaseDB {
	factory := func() (interface{}, error) { return gohbase.NewClient(dbUrl), nil }
	close := func(v interface{}) error { v.(gohbase.Client).Close(); return nil }

	//pool config
	poolConfig := &pool.Config{
		InitialCap:  1000,
		MaxIdle:     1000,
		MaxCap:      1000,
		Factory:     factory,
		Close:       close,
		IdleTimeout: 5 * time.Second,
	}
	pool, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		panic(err)
	}
	return &HbaseDB{db: pool}
}

func (db *HbaseDB) Set(key []byte, value []byte) {
	v, err := db.db.Get()
	if v == nil || err != nil {
		fmt.Println("db get error:", err)
	}
	client, _ := v.(gohbase.Client)
	defer db.db.Put(client)

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

func (db *HbaseDB) Get(key []byte) ([]byte, error) {
	v, err := db.db.Get()
	if v == nil || err != nil {
		fmt.Println("db get error:", err)
		return []byte{}, err
	}
	client, _ := v.(gohbase.Client)
	defer db.db.Put(client)

	getRequest, _ := hrpc.NewGet(context.Background(), []byte(TableName), key)
	//todo del
	fromTime := time.Now()
	result, err := client.Get(getRequest)
	fmt.Println("HbaseDB get spend time ", time.Since(fromTime))
	if err != nil || len(result.Cells) == 0 {
		return []byte{}, err
	}
	return result.Cells[0].Value, err
}

func (db *HbaseDB) Delete(key []byte) {
	v, err := db.db.Get()
	if v == nil || err != nil {
		fmt.Println("db get error:", err)
	}
	client, _ := v.(gohbase.Client)
	defer db.db.Put(client)

	mutate, err := hrpc.NewDel(context.Background(), []byte(TableName), key, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = client.Delete(mutate)
	if err != nil {
		fmt.Println(err)
	}
}

func (db *HbaseDB) Has(key []byte) bool {
	data, err := db.Get(key)
	return len(data) > 0 && err == nil
}
