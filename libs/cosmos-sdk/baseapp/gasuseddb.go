package baseapp

import (
	"encoding/binary"
	"path/filepath"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	db "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/viper"
)

const (
	HistoryGasUsedDbDir  = "data"
	HistoryGasUsedDBName = "hgu"

	FlagGasUsedFactor = "gu_factor"
)

var (
	once          sync.Once
	GasUsedFactor = 0.4
	jobQueueLen   = 10
	cacheSize     = 10000

	historyGasUsedRecordDB HistoryGasUsedRecordDB
)

type gasKey struct {
	gas int64
	key string
}

type HistoryGasUsedRecordDB struct {
	latestGuMtx sync.Mutex
	latestGu    map[string]int64
	cache       *lru.Cache
	guDB        db.DB

	jobQueue chan func()
}

func InstanceOfHistoryGasUsedRecordDB() *HistoryGasUsedRecordDB {
	once.Do(func() {
		cache, _ := lru.New(cacheSize)
		historyGasUsedRecordDB = HistoryGasUsedRecordDB{
			latestGu: make(map[string]int64),
			cache:    cache,
			guDB:     initDb(),
			jobQueue: make(chan func(), jobQueueLen),
		}
		go historyGasUsedRecordDB.updateRoutine()
	})
	return &historyGasUsedRecordDB
}

func (h *HistoryGasUsedRecordDB) UpdateGasUsed(key []byte, gasUsed int64) {
	h.latestGuMtx.Lock()
	h.latestGu[string(key)] = gasUsed
	h.latestGuMtx.Unlock()
}

func (h *HistoryGasUsedRecordDB) GetHgu(key []byte) int64 {
	hgu, cacheHit := h.getHgu(key)
	if !cacheHit && hgu != -1 {
		// add to cache before returning hgu
		h.cache.Add(string(key), hgu)
	}
	return hgu
}

func (h *HistoryGasUsedRecordDB) FlushHgu() {
	if len(h.latestGu) == 0 {
		return
	}
	latestGasKeys := make([]gasKey, len(h.latestGu))
	for key, gas := range h.latestGu {
		latestGasKeys = append(latestGasKeys, gasKey{
			gas: gas,
			key: key,
		})
		delete(h.latestGu, key)
	}
	h.jobQueue <- func() { h.flushHgu(latestGasKeys...) } // closure
}

func (h *HistoryGasUsedRecordDB) getHgu(key []byte) (hgu int64, fromCache bool) {
	v, ok := h.cache.Get(string(key))
	if ok {
		return v.(int64), true
	}

	data, err := h.guDB.Get(key)
	if err != nil || len(data) == 0 {
		return -1, false
	}

	return bytesToInt64(data), false
}

func (h *HistoryGasUsedRecordDB) flushHgu(gks ...gasKey) {
	for _, gk := range gks {
		hgu, cacheHit := h.getHgu([]byte(gk.key))
		// avgGas = 0.4 * newGas + 0.6 * oldGas
		avgGas := int64(GasUsedFactor*float64(gk.gas) + (1.0-GasUsedFactor)*float64(hgu))
		// add to cache if hit
		if cacheHit {
			h.cache.Add(gk.key, avgGas)
		}
		h.guDB.Set([]byte(gk.key), int64ToBytes(avgGas))
	}
}

func (h *HistoryGasUsedRecordDB) updateRoutine() {
	for job := range h.jobQueue {
		job()
	}
}

func initDb() db.DB {
	homeDir := viper.GetString(flags.FlagHome)
	dbPath := filepath.Join(homeDir, HistoryGasUsedDbDir)

	db, err := sdk.NewDB(HistoryGasUsedDBName, dbPath)
	if err != nil {
		panic(err)
	}
	return db
}

func int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func bytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
