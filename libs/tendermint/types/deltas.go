package types

import (
	"fmt"
	"github.com/okex/exchain/libs/tendermint/libs/compress"
	"github.com/spf13/viper"
	"sync"
	"time"
)

const (
	FlagDownloadDDS = "download-delta"
	FlagUploadDDS = "upload-delta"
	FlagAppendPid = "append-pid"
	FlagDDSCompressType = "compress-type"
	FlagDDSCompressFlag = "compress-flag"

	// redis
	// url fmt (ip:port)
	FlagRedisUrl    = "delta-redis-url"
	FlagRedisAuth   = "delta-redis-auth"
	// expire unit: second
	FlagRedisExpire = "delta-redis-expire"

	FlagFastQuery = "fast-query"

	// when this DeltaVersion not equal with dds delta-version, can't use delta
	DeltaVersion = 2
)

var (
	fastQuery = false
	downloadDelta    = false
	uploadDelta      = false

	onceFastQuery   sync.Once
	onceDownload     sync.Once
	onceUpload       sync.Once
)

func IsFastQuery() bool {
	onceFastQuery.Do(func() {
		fastQuery = viper.GetBool(FlagFastQuery)
	})
	return fastQuery
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

type DeltasMessage struct {
	Metadata         []byte `json:"metadata"`
	Height           int64  `json:"height"`
	Version          int    `json:"version"`
	CompressType     int    `json:"compress_type"`
	From			 string `json:"from"`
}

type DeltaPayload struct {
	ABCIRsp     []byte
	DeltasBytes []byte
	WatchBytes  []byte
}

// Deltas defines the ABCIResponse and state delta
type Deltas struct {
	Height           int64
	Version          int
	Payload          DeltaPayload
	CompressType     int
	CompressFlag     int
	From			 string

	marshalElapsed    time.Duration
	compressElapsed   time.Duration
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


// Marshal returns the amino encoding.
func (d *Deltas) Marshal() ([]byte, error) {
	t0 := time.Now()

	payload, err := cdc.MarshalBinaryBare(&d.Payload)
	if err != nil {
		return nil, err
	}

	t1 := time.Now()
	payload, err = compress.Compress(d.CompressType, d.CompressFlag, payload)
	if err != nil {
		return nil, err
	}
	t2 := time.Now()

	dt := &DeltasMessage{
		Metadata: payload,
		Height: d.Height,
		Version: d.Version,
		CompressType: d.CompressType,
		From: d.From,
	}

	res, err := cdc.MarshalBinaryBare(dt)
	t3 := time.Now()

	d.compressElapsed = t2.Sub(t1)
	d.marshalElapsed = t3.Sub(t0) - d.compressElapsed

	return res, err
}

// Unmarshal deserializes from amino encoded form.
func (d *Deltas) Unmarshal(bs []byte) error {
	t0 := time.Now()

	msg := &DeltasMessage{}
	err := cdc.UnmarshalBinaryBare(bs, msg)
	if err != nil {
		return err
	}
	d.CompressType = msg.CompressType

	t1 := time.Now()
	msg.Metadata, err = compress.UnCompress(d.CompressType, msg.Metadata)
	if err != nil {
		return err
	}
	t2 := time.Now()

	err = cdc.UnmarshalBinaryBare(msg.Metadata, &d.Payload)
	t3 := time.Now()

	d.Version = msg.Version
	d.Height = msg.Height
	d.From = msg.From

	d.compressElapsed = t2.Sub(t1)
	d.marshalElapsed = t3.Sub(t0) - d.compressElapsed
	return err
}

func (d *Deltas) String() string {
	return fmt.Sprintf("height<%d>, version<%d>, size<%d>, from<%s>",
		d.Height,
		d.Version,
		d.Size(),
		d.From,
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
