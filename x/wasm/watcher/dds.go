package watcher

import (
	"bytes"
	"log"
	"sort"

	"github.com/golang/protobuf/proto"
	tmstate "github.com/okx/okbchain/libs/tendermint/state"
)

func SetWatchDataManager() {
	tmstate.SetWasmWatchDataManager(WatchDataManager{})
}

type WatchDataManager struct{}

func (w WatchDataManager) CreateWatchDataGenerator() func() ([]byte, error) {
	data := &WatchData{
		Messages: make([]*WatchMessage, 0, len(blockStateCache)),
	}
	for _, v := range blockStateCache {
		data.Messages = append(data.Messages, v)
	}
	sort.Sort(data)
	return func() ([]byte, error) {
		return proto.Marshal(data)
	}
}

func (w WatchDataManager) UnmarshalWatchData(b []byte) (interface{}, error) {
	if len(b) == 0 {
		return nil, nil
	}
	var data WatchData
	err := proto.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (w WatchDataManager) ApplyWatchData(v interface{}) {
	data, ok := v.(*WatchData)
	if !ok {
		return
	}
	task := func() {
		batch := db.NewBatch()
		for _, msg := range data.Messages {
			if msg.IsDelete {
				batch.Delete(msg.Key)
			} else {
				batch.Set(msg.Key, msg.Value)
			}
		}
		if err := batch.Write(); err != nil {
			log.Println("ApplyWatchData batch write error:" + err.Error())
		}
	}

	tasks <- task
}

func (d *WatchData) Len() int {
	return len(d.Messages)
}

func (d *WatchData) Less(i, j int) bool {
	return bytes.Compare(d.Messages[i].Key, d.Messages[j].Key) < 0
}

func (d *WatchData) Swap(i, j int) {
	d.Messages[i], d.Messages[j] = d.Messages[j], d.Messages[i]
}
