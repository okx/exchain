package types

import (
	"github.com/okex/exchain/libs/iavl"
	"github.com/spf13/viper"
	"sync"
)

const (
	// for getting flag of delta-mode
	FlagStateDelta = "state-sync-mode"

	// delta-mode
	// state-delta mode
	// na:not available. same as no FlagStateDelta
	// producer. product delta
	// consumer. consume delta; if get no delta, do as producer
	NoDelta      = "na"
	ProductDelta = "producer"
	ConsumeDelta = "consumer"

	// p2p-delta
	// the switch of delta to p2p
	// true: save into deltastore.db, and add delta into bcBlockResponseMessage
	FlagP2PDelta = "send-p2p-delta"

	// data-center
	FlagDataCenter = "data-center-mode"
	DataCenterUrl  = "data-center-url"
	DataCenterStr  = "dataCenter"

	// origin direct delta
	// get delta from dc/redis
	FlagOriginDDS = "origin-direct-delta"

	// send direct delta
	// send delta to dc/redis
	FlagSendDDS = "send-direct-delta"

	// fast-query
	FlagFastQuery = "fast-query"
)

var (
	deltaMode = NoDelta
	fastQuery = false
	p2pDelta = false
	originDDS = false
	sendDDS = false
	centerMode = false
	centerUrl = "127.0.0.1:8030"

	onceDeltaMode     sync.Once
	onceFastQuery	sync.Once
	onceP2PDelta	sync.Once
	onceOriginDDS	sync.Once
	onceSendDDS		sync.Once
	onceCenterMode	sync.Once
	onceCenterUrl	sync.Once
)

func GetDeltaMode() string {
	onceDeltaMode.Do(func() {
		deltaMode = viper.GetString(FlagStateDelta)
		iavl.SetDeltaMode(deltaMode)
	})
	return deltaMode
}

func IsFastQuery() bool {
	onceFastQuery.Do(func() {
		fastQuery = viper.GetBool(FlagFastQuery)
	})
	return fastQuery
}

func IsP2PDeltaEnabled() bool {
	onceP2PDelta.Do(func() {
		p2pDelta = viper.GetBool(FlagP2PDelta)
	})
	return p2pDelta
}

func EnableOriginDDS() bool {
	onceOriginDDS.Do(func() {
		originDDS = viper.GetBool(FlagOriginDDS)
	})
	return originDDS
}

func EnableSendDDS() bool {
	onceSendDDS.Do(func() {
		sendDDS = viper.GetBool(FlagSendDDS)
	})
	return sendDDS
}

func IsCenterEnabled() bool {
	onceCenterMode.Do(func() {
		centerMode = viper.GetBool(FlagDataCenter)
	})
	return centerMode
}

func GetCenterUrl() string {
	onceCenterUrl.Do(func() {
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
