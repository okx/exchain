package baseapp

import (
	"math"
	"path/filepath"
	"sync"

	"github.com/gogo/protobuf/proto"
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
	jobQueueLen   = 10
	cacheSize     = 10000
	GasUsedFactor = 0.4
)

type gasKey struct {
	gas int64
	key []byte
}

type HistoryGasUsedRecordDB struct {
	guDB        db.DB
	latestGuMtx sync.Mutex
	latestGu    map[string][]int64
	cache       *lru.Cache
	jobQueue    chan func()
}

var historyGasUsedRecordDB HistoryGasUsedRecordDB

func InstanceOfHistoryGasUsedRecordDB() *HistoryGasUsedRecordDB {
	once.Do(func() {
		cache, _ := lru.New(cacheSize)
		historyGasUsedRecordDB = HistoryGasUsedRecordDB{
			guDB:     initDb(),
			latestGu: make(map[string][]int64),
			cache:    cache,
			jobQueue: make(chan func(), jobQueueLen),
		}
		go historyGasUsedRecordDB.updateRoutine()
	})
	return &historyGasUsedRecordDB
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

func (h *HistoryGasUsedRecordDB) UpdateGasUsed(key []byte, gasUsed int64) {
	h.latestGuMtx.Lock()
	h.latestGu[string(key)] = append(h.latestGu[string(key)], gasUsed)
	h.latestGuMtx.Unlock()
}

func (h *HistoryGasUsedRecordDB) GetHgu(key []byte) *HguRecord {
	v, ok := h.cache.Get(string(key))
	if ok {
		return v.(*HguRecord)
	}
	var record HguRecord
	data, err := h.guDB.Get(key)
	if err != nil || len(data) == 0 {
		return nil
	}
	err = proto.Unmarshal(data, &record)
	if err != nil {
		return nil
	}
	return &record
}

func (h *HistoryGasUsedRecordDB) FlushHgu() {
	if len(h.latestGu) == 0 {
		return
	}
	latestMeanGu := make([]gasKey, len(h.latestGu))
	for key, gus := range h.latestGu {
		latestMeanGu = append(latestMeanGu, gasKey{
			gas: meanInt64(gus...),
			key: []byte(key),
		})
		delete(h.latestGu, key)
	}
	h.jobQueue <- func() { h.flushHgu(latestMeanGu...) } // closure function
}

func (h *HistoryGasUsedRecordDB) flushHgu(gks ...gasKey) {
	for _, gk := range gks {
		meanGu := gk.gas
		key := gk.key
		hgu := h.GetHgu(key)
		if hgu == nil {
			hgu = &HguRecord{
				MaxGas:       meanGu,
				LastBlockGas: meanGu,
				HighGas:      meanGu,
				StandardGas:  meanGu,
			}
		} else {
			hgu.LastBlockGas = meanGu
			if meanGu > hgu.MaxGas {
				hgu.MaxGas = meanGu
			}
			if meanGu >= hgu.StandardGas*2 {
				// abnormal cases which meanGu of last block is much higher than StandardGas
				hgu.HighGas = int64(GasUsedFactor*float64(meanGu) + (1-GasUsedFactor)*float64(hgu.HighGas))
			} else if meanGu > hgu.StandardGas/2 {
				// StandardGas/2 < meanGu < StandardGas*2
				// normal cases
				hgu.StandardGas = int64(GasUsedFactor*float64(meanGu) + (1-GasUsedFactor)*float64(hgu.StandardGas))
				hgu.HighGas = maxInt64(hgu.HighGas, hgu.StandardGas)
			} else {
				// meanGu <= hgu.StandardGas/2
				// in case that the meanGu is too big in the first init, so reinit StandardGas and HighGas of hgu
				hgu.StandardGas = meanGu
			}
		}
		h.setHgu(key, hgu)
	}
}

func (h *HistoryGasUsedRecordDB) setHgu(key []byte, hgu *HguRecord) {
	h.cache.Add(string(key), hgu)
	data, err := proto.Marshal(hgu)
	if err != nil {
		return
	}
	_ = h.guDB.Set(key, data)
}

func (h *HistoryGasUsedRecordDB) updateRoutine() {
	for job := range h.jobQueue {
		job()
	}
}

// for unit test
func (h *HistoryGasUsedRecordDB) close() {
	h.guDB.Close()
	// for recreate the instance
	once = sync.Once{}
}

func estimateGas(gasLimit int64, hgu *HguRecord) int64 {
	if hgu == nil {
		return gasLimit
	}
	switch {
	case hgu.LastBlockGas >= hgu.HighGas:
		return minInt64(gasLimit, hgu.MaxGas)
	case hgu.LastBlockGas >= hgu.StandardGas*2:
		return minInt64(gasLimit, hgu.MaxGas) * 3 / 4 // 75%
	case hgu.LastBlockGas >= hgu.StandardGas:
		return minInt64(gasLimit, (hgu.LastBlockGas+hgu.StandardGas)/2)
	default:
		return minInt64(gasLimit, hgu.StandardGas)
	}
}

func meanInt64(numbers ...int64) int64 {
	if len(numbers) == 0 {
		return 0
	}
	var total int64
	for _, num := range numbers {
		total += num
	}
	return total / int64(len(numbers))
}

func minInt64(numbers ...int64) int64 {
	res := int64(math.MaxInt64)
	for _, num := range numbers {
		if num < res {
			res = num
		}
	}
	return res
}

func maxInt64(numbers ...int64) int64 {
	res := int64(math.MinInt64)
	for _, num := range numbers {
		if num > res {
			res = num
		}
	}
	return res
}
