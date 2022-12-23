package baseapp

import (
	"log"
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
	once             sync.Once
	GasUsedFactor    = 0.4
	regressionFactor = 0.01
	jobQueueLen      = 10
	cacheSize        = 10000

	historyGasUsedRecordDB HistoryGasUsedRecordDB
)

type gasKey struct {
	gas int64
	key string
}

type HistoryGasUsedRecordDB struct {
	latestGuMtx sync.Mutex
	latestGu    map[string][]int64
	cache       *lru.Cache
	guDB        db.DB

	jobQueue chan func()
}

func InstanceOfHistoryGasUsedRecordDB() *HistoryGasUsedRecordDB {
	once.Do(func() {
		cache, _ := lru.New(cacheSize)
		historyGasUsedRecordDB = HistoryGasUsedRecordDB{
			latestGu: make(map[string][]int64),
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
	h.latestGu[string(key)] = append(h.latestGu[string(key)], gasUsed)
	h.latestGuMtx.Unlock()
}

func (h *HistoryGasUsedRecordDB) GetHgu(key []byte) *HguRecord {
	hgu, cacheHit := h.getHgu(key)
	if hgu != nil && !cacheHit {
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
	for key, allGas := range h.latestGu {
		latestGasKeys = append(latestGasKeys, gasKey{
			gas: meanGas(allGas),
			key: key,
		})
		delete(h.latestGu, key)
	}
	h.jobQueue <- func() { h.flushHgu(latestGasKeys...) } // closure
}

func (h *HistoryGasUsedRecordDB) getHgu(key []byte) (hgu *HguRecord, fromCache bool) {
	v, ok := h.cache.Get(string(key))
	if ok {
		return v.(*HguRecord), true
	}

	data, err := h.guDB.Get(key)
	if err != nil || len(data) == 0 {
		return nil, false
	}

	var r HguRecord
	err = proto.Unmarshal(data, &r)
	if err != nil {
		return nil, false
	}
	return &r, false
}

func (h *HistoryGasUsedRecordDB) flushHgu(gks ...gasKey) {
	for _, gk := range gks {
		hgu, cacheHit := h.getHgu([]byte(gk.key))
		if hgu == nil {
			hgu = &HguRecord{
				MaxGas:           gk.gas,
				MinGas:           gk.gas,
				MovingAverageGas: gk.gas,
			}
		} else {
			// MovingAverageGas = 0.4 * newGas + 0.6 * oldMovingAverageGas
			hgu.MovingAverageGas = int64(GasUsedFactor*float64(gk.gas) + (1.0-GasUsedFactor)*float64(hgu.MovingAverageGas))
			// MaxGas = 0.01 * MovingAverageGas + 0.99 * oldMaxGas
			hgu.MaxGas = int64(regressionFactor*float64(hgu.MovingAverageGas) + (1.0-regressionFactor)*float64(hgu.MaxGas))
			// MinGas = 0.01 * MovingAverageGas + 0.99 * oldMinGas
			hgu.MinGas = int64(regressionFactor*float64(hgu.MovingAverageGas) + (1.0-regressionFactor)*float64(hgu.MinGas))
			if gk.gas > hgu.MaxGas {
				hgu.MaxGas = gk.gas
			} else if gk.gas < hgu.MinGas {
				hgu.MinGas = gk.gas
			}
			// add to cache if hit
			if cacheHit {
				h.cache.Add(gk.key, hgu)
			}
		}

		data, err := proto.Marshal(hgu)
		if err != nil {
			log.Println("flushHgu marshal error:", err)
			continue
		}

		h.guDB.Set([]byte(gk.key), data)
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

func meanGas(allGas []int64) int64 {
	var totalGas int64
	for _, gas := range allGas {
		totalGas += gas
	}
	return totalGas / int64(len(allGas))
}
