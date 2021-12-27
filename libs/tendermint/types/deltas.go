package types

import (
	"encoding/hex"
	"fmt"
	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	"github.com/okex/exchain/libs/tendermint/libs/compress"
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
	FlagUploadDDS       = "upload-delta"
	FlagDDSCompressType = "compress-type"
	FlagDDSCompressFlag = "compress-flag"

	// redis
	FlagRedisUrl    = "delta-redis-url"
	FlagRedisAuth   = "delta-redis-auth"
	FlagRedisExpire = "delta-redis-expire"

	// fast-query
	FlagFastQuery = "fast-query"

	// delta version
	// when this DeltaVersion not equal with dds delta-version, can't use delta
	DeltaVersion = 2
)

var (
	fastQuery = false
	// fmt (http://ip:port/)
	centerUrl = "http://127.0.0.1:8030/"
	// fmt (ip:port)
	redisUrl  = "127.0.0.1:6379"
	redisAuth = "auth"
	// unit: second
	redisExpire = 300

	applyP2PDelta    = false
	broadcatP2PDelta = false
	downloadDelta    = false
	uploadDelta      = false

	onceFastQuery   sync.Once
	onceRedisUrl    sync.Once
	onceRedisAuth   sync.Once
	onceRedisExpire sync.Once

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

type DeltasMessage struct {
	Metadata     []byte `json:"metadata"`
	Height       int64  `json:"height"`
	Version      int    `json:"version"`
	CompressType int    `json:"compress_type"`
	PayloadHash  []byte `json:"payload_hash"`
}

type DeltaPayload struct {
	ABCIRsp     []byte
	DeltasBytes []byte
	WatchBytes  []byte
}

// Deltas defines the ABCIResponse and state delta
type Deltas struct {
	Height       int64
	Version      int
	Payload      DeltaPayload
	CompressType int
	CompressFlag int

	marshalElapsed  time.Duration
	compressElapsed time.Duration
	hashElapsed     time.Duration
}

// Size returns size of the deltas in bytes.
func (d *Deltas) Size() int {
	return len(d.ABCIRsp()) + len(d.DeltasBytes()) + len(d.WatchBytes())
}
func (d *Deltas) ABCIRsp() []byte {
	return d.Payload.ABCIRsp
}

func (d *Deltas) DeltasBytes() []byte {
	return d.Payload.DeltasBytes
}

func (d *Deltas) WatchBytes() []byte {
	return d.Payload.WatchBytes
}

func (d *Deltas) MarshalOrUnmarshalElapsed() time.Duration {
	return d.marshalElapsed
}
func (d *Deltas) CompressOrUncompressElapsed() time.Duration {
	return d.compressElapsed
}
func (d *Deltas) HashElapsed() time.Duration {
	return d.hashElapsed
}

// Marshal returns the amino encoding.
func (d *Deltas) Marshal() ([]byte, error) {
	t0 := time.Now()

	// marshal to payload bytes
	payload, err := cdc.MarshalBinaryBare(&d.Payload)
	if err != nil {
		return nil, err
	}

	t1 := time.Now()
	// calc payload hash
	payloadHash := tmhash.Sum(payload)

	// compress
	t2 := time.Now()
	payload, err = compress.Compress(d.CompressType, d.CompressFlag, payload)
	if err != nil {
		return nil, err
	}
	t3 := time.Now()

	dt := &DeltasMessage{
		Metadata:     payload,
		Height:       d.Height,
		Version:      d.Version,
		CompressType: d.CompressType,
		PayloadHash:  payloadHash,
	}

	// marshal to upload bytes
	res, err := cdc.MarshalBinaryBare(dt)
	t4 := time.Now()

	d.hashElapsed = t2.Sub(t1)
	d.compressElapsed = t3.Sub(t2)
	d.marshalElapsed = t4.Sub(t0) - d.compressElapsed - d.hashElapsed

	return res, err
}

// Unmarshal deserializes from amino encoded form.
func (d *Deltas) Unmarshal(bs []byte) error {
	t0 := time.Now()

	// unmarshal to DeltasMessage
	msg := &DeltasMessage{}
	err := cdc.UnmarshalBinaryBare(bs, msg)
	if err != nil {
		return err
	}

	t1 := time.Now()
	// calc payload hash
	payloadHash := tmhash.Sum(msg.Metadata)
	if hex.EncodeToString(payloadHash) != hex.EncodeToString(msg.PayloadHash) {
		return fmt.Errorf("metadata hash is different")
	}

	// uncompress
	d.CompressType = msg.CompressType
	t2 := time.Now()
	msg.Metadata, err = compress.UnCompress(d.CompressType, msg.Metadata)
	if err != nil {
		return err
	}
	t3 := time.Now()

	err = cdc.UnmarshalBinaryBare(msg.Metadata, &d.Payload)
	t4 := time.Now()

	d.Version = msg.Version
	d.Height = msg.Height

	d.hashElapsed = t2.Sub(t1)
	d.compressElapsed = t3.Sub(t2)
	d.marshalElapsed = t4.Sub(t0) - d.compressElapsed - d.hashElapsed
	return err
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
		len(dds.WatchBytes()) == 0 ||
		len(dds.ABCIRsp()) == 0 ||
		len(dds.DeltasBytes()) == 0 {
		return false
	}
	return true
}
