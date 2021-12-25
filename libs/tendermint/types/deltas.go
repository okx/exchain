package types

import (
	"fmt"
	"github.com/spf13/viper"
	"sync"
	"time"
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
	FlagRedisUrl    = "delta-redis-url"
	FlagRedisAuth   = "delta-redis-auth"
	FlagRedisExpire = "delta-redis-expire"
	FlagRedisLocker = "delta-redis-locker"

	// data-center
	FlagDataCenter = "data-center-mode"
	DataCenterUrl  = "data-center-url"

	// fast-query
	FlagFastQuery = "fast-query"

	// delta version
	// when this DeltaVersion not equal with dds delta-version, can't use delta
	DeltaVersion = 1
)

var (
	fastQuery = false
	// fmt (http://ip:port/)
	centerUrl = "http://127.0.0.1:8030/"
	// fmt (ip:port)
	redisUrl  = "127.0.0.1:6379"
	redisAuth = "auth"
	redisLockerID = "locker"
	// unit: second
	redisExpire = 300

	applyP2PDelta    = false
	broadcatP2PDelta = false
	downloadDelta    = false
	uploadDelta      = false

	onceFastQuery   sync.Once
	onceCenterUrl   sync.Once
	onceRedisUrl    sync.Once
	onceRedisAuth   sync.Once
	onceRedisExpire sync.Once
	onceRedisLocker sync.Once

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
	})
	return uploadDelta
}

func RedisUrl() string {
	onceRedisUrl.Do(func() {
		redisUrl = viper.GetString(FlagRedisUrl)
	})
	return redisUrl
}

func RedisAuth() string {
	onceRedisAuth.Do(func() {
		redisAuth = viper.GetString(FlagRedisAuth)
	})
	return redisAuth
}

func RedisExpire() time.Duration {
	onceRedisExpire.Do(func() {
		redisExpire = viper.GetInt(FlagRedisExpire)
	})
	return time.Duration(redisExpire) * time.Second
}

func RedisLocker() string {
	onceRedisLocker.Do(func() {
		redisLockerID = viper.GetString(FlagRedisLocker)
	})
	return redisLockerID
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
	Version     int    `json:"version"`
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

func (d *Deltas) String() string {
	return fmt.Sprintf("height<%d>, version<%d>, size<%d>",
		d.Height,
		d.Version,
		d.Size(),
		)
}

func (dds *Deltas) Validate(height int64) bool {
	if DeltaVersion < dds.Version ||
		dds.Height != height ||
		len(dds.WatchBytes) == 0 ||
		len(dds.ABCIRsp) == 0 ||
		len(dds.DeltasBytes) == 0 {
		return false
	}
	return true
}
