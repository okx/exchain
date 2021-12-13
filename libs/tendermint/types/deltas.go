package types

import (
	"github.com/okex/exchain/libs/iavl"
	"github.com/spf13/viper"
	"sync"
)

const (
	// use delta from bcBlockResponseMessage or not
	FlagApplyP2PDelta = "apply-p2p-delta"
	// save into deltastore.db, and add delta into bcBlockResponseMessage
	FlagBroadcastP2PDelta = "broadcast-delta"
	// get delta from dc/redis
	FlagDownloadDDS = "download-delta"
	// send delta to dc/redis
	FlagUploadDDS = "upload-delta"

	// redis
	FlagRedisUrl = "redis-url"

	// data-center
	FlagDataCenter = "data-center-mode"
	DataCenterUrl  = "data-center-url"

	// fast-query
	FlagFastQuery = "fast-query"
)

var (
	fastQuery  = false
	centerMode = false
	// fmt (http://ip:port/)
	centerUrl  = "http://127.0.0.1:8030/"
	// fmt (ip:port)
	redisUrl   = "127.0.0.1:6379"

	applyP2PDelta    = false
	broadcatP2PDelta = false
	downloadDelta    = false
	uploadDelta      = false

	onceFastQuery  sync.Once
	onceCenterMode sync.Once
	onceCenterUrl  sync.Once
	onceRedisUrl   sync.Once

	onceApplyP2P     sync.Once
	onceBroadcastP2P sync.Once
	onceDownload     sync.Once
	onceUpload       sync.Once
)

func IsFastQuery() bool {
	onceFastQuery.Do(func() {
		fastQuery = viper.GetBool(FlagFastQuery)
	})
	return fastQuery
}

func EnableApplyP2PDelta() bool {
	onceApplyP2P.Do(func() {
		applyP2PDelta = viper.GetBool(FlagApplyP2PDelta)
	})
	return applyP2PDelta
}

func EnableBroadcastP2PDelta() bool {
	onceBroadcastP2P.Do(func() {
		broadcatP2PDelta = viper.GetBool(FlagBroadcastP2PDelta)
		iavl.SetProduceDelta(broadcatP2PDelta)
	})
	return broadcatP2PDelta
}

func EnableDownloadDelta() bool {
	onceDownload.Do(func() {
		downloadDelta = viper.GetBool(FlagDownloadDDS)
	})
	return downloadDelta
}

func EnableUploadDelta() bool {
	onceUpload.Do(func() {
		uploadDelta = viper.GetBool(FlagUploadDDS)
		iavl.SetProduceDelta(uploadDelta)
	})
	return uploadDelta
}

func RedisUrl() string {
	onceRedisUrl.Do(func() {
		redisUrl = viper.GetString(FlagRedisUrl)
	})
	return redisUrl
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
	WatchBytes  []byte `json:"watch_bytes"`
	Height      int64  `json:"height"`
}

// Size returns size of the deltas in bytes.
func (d *Deltas) Size() int {
	return len(d.ABCIRsp) + len(d.DeltasBytes) + len(d.WatchBytes)
}

// Marshal returns the amino encoding.
func (d *Deltas) Marshal() ([]byte, error) {
	return cdc.MarshalBinaryBare(d)
}

// Unmarshal deserializes from amino encoded form.
func (d *Deltas) Unmarshal(bs []byte) error {
	return cdc.UnmarshalBinaryBare(bs, d)
}
