package baseapp

import (
	"path/filepath"
	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	hgutypes "github.com/okex/exchain/libs/tendermint/types/hgu"
	db "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/viper"
)

const (
	HistoryGasUsedDbDir  = "data"
	HistoryGasUsedDBName = "hgu"

	FlagGasUsedFactor = "gu_factor"

	jobQueueLen = 10
)

var (
	once          sync.Once
	GasUsedFactor = 0.4
)

type HistoryGasUsedRecordDB struct {
	guDB        db.DB
	latestGuMtx sync.Mutex
	latestGu    map[string][]int64
	jobQueue    chan func()
}

var historyGasUsedRecordDB HistoryGasUsedRecordDB

func InstanceOfHistoryGasUsedRecordDB() *HistoryGasUsedRecordDB {
	once.Do(func() {
		historyGasUsedRecordDB = HistoryGasUsedRecordDB{
			guDB:     initDb(),
			latestGu: make(map[string][]int64),
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

func (h *HistoryGasUsedRecordDB) GetHgu(key []byte) *hgutypes.HguRecord {
	var record hgutypes.HguRecord
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
	latestMeanGu := make(map[string]int64, len(h.latestGu))
	for key, gus := range h.latestGu {
		latestMeanGu[key] = meanInt64(gus...)
		delete(h.latestGu, key)
	}
	h.jobQueue <- func() {
		h.flushHgu(latestMeanGu)
	}

}

func (h *HistoryGasUsedRecordDB) flushHgu(lgu map[string]int64) {
	for key, meanGu := range lgu {
		hgu := h.GetHgu([]byte(key))
		if hgu == nil {
			hgu = &hgutypes.HguRecord{
				MaxGas:       meanGu,
				LastBlockGas: meanGu,
				HighGas:      meanGu * 2,
				StandardGas:  meanGu,
			}
		} else {
			hgu.LastBlockGas = meanGu
			if meanGu > hgu.MaxGas {
				hgu.MaxGas = meanGu
			}
			if meanGu >= hgu.StandardGas*2 {
				hgu.HighGas = int64(GasUsedFactor*float64(meanGu) + (1-GasUsedFactor)*float64(hgu.HighGas))
			} else if meanGu > hgu.StandardGas/2 {
				hgu.StandardGas = int64(GasUsedFactor*float64(meanGu) + (1-GasUsedFactor)*float64(hgu.StandardGas))
			} else {
				// meanGu <= hgu.StandardGas/2
				// in case that the meanGu is too big in the first init, so reinit StandardGas and HighGas of hgu
				hgu.StandardGas = meanGu
				hgu.HighGas = meanGu * 2
			}
		}
		data, err := proto.Marshal(hgu)
		if err != nil {
			return
		}
		h.guDB.Set([]byte(key), data)
	}
}

func (h *HistoryGasUsedRecordDB) updateRoutine() {
	for job := range h.jobQueue {
		job()
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
