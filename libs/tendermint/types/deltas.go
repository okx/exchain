package types

import (
	"github.com/okex/exchain/libs/iavl"
	"github.com/spf13/viper"
	"sync"
)

// state-delta mode
// 0 same as no state-delta
// 1 product delta and save into deltastore.db
// 2 consume delta and save into deltastore.db; if get no delta, do as 1
const (
	// for getting flag of delta-mode
	FlagStateDelta = "state-sync-mode"

	// delta-mode
	NoDelta      = "na"
	ProductDelta = "producer"
	ConsumeDelta = "consumer"

	// data-center
	FlagDataCenter = "data-center-mode"
	DataCenterUrl  = "data-center-url"
	DataCenterStr  = "dataCenter"

	// fast-query
	FlagFastQuery = "fast-query"
)

var (
	deltaMode = NoDelta
	fastQuery = false
	centerMode = false
	centerUrl = "127.0.0.1:8030"
	onceEnable     sync.Once
)

func GetDeltaMode() string {
	onceEnable.Do(func() {
		deltaMode = viper.GetString(FlagStateDelta)
		iavl.SetDeltaMode(deltaMode)
	})
	return deltaMode
}

func IsFastQuery() bool {
	onceEnable.Do(func() {
		fastQuery = viper.GetBool(FlagFastQuery)
	})
	return fastQuery
}

func IsCenterEnabled() bool {
	onceEnable.Do(func() {
		centerMode = viper.GetBool(FlagDataCenter)
	})
	return centerMode
}

func GetCenterUrl() string {
	onceEnable.Do(func() {
		centerUrl = viper.GetString(DataCenterUrl)
	})
	return centerUrl
}

// Deltas defines the ABCIResponse and state delta
type Deltas struct {
	ABCIRsp     []byte `json:"abci_rsp"`
	DeltasBytes []byte `json:"deltas_bytes"`
	Height      int64  `json:"height"`
}

// Size returns size of the deltas in bytes.
func (d *Deltas) Size() int {
	if d == nil {
		return -1
	}
	return len(d.ABCIRsp) + len(d.DeltasBytes)
}

// Marshal returns the amino encoding.
func (d *Deltas) Marshal() ([]byte, error) {
	return cdc.MarshalBinaryBare(d)
}

// Unmarshal deserializes from amino encoded form.
func (d *Deltas) Unmarshal(bs []byte) error {
	return cdc.UnmarshalBinaryBare(bs, d)
}

// WatchData defines the batch in watchDB and accounts for delete
type WatchData struct {
	WatchDataByte []byte `json:"watch_data_byte"`
	Height        int64  `json:"height"`
}

// Size returns size of the deltas in bytes.
func (wd *WatchData) Size() int {
	if wd == nil {
		return -1
	}
	return len(wd.WatchDataByte)
}
