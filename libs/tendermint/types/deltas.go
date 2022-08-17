package types

import (
	"bytes"
	"fmt"
	"github.com/tendermint/go-amino"
	"time"

	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	"github.com/okex/exchain/libs/tendermint/libs/compress"
)

const (
	FlagDownloadDDS     = "download-delta"
	FlagUploadDDS       = "upload-delta"
	FlagAppendPid       = "append-pid"
	FlagBufferSize      = "delta-buffer-size"
	FlagDDSCompressType = "compress-type"
	FlagDDSCompressFlag = "compress-flag"

	// redis
	// url fmt (ip:port)
	FlagRedisUrl  = "delta-redis-url"
	FlagRedisAuth = "delta-redis-auth"
	// expire unit: second
	FlagRedisExpire = "delta-redis-expire"
	FlagRedisDB     = "delta-redis-db"
	FlagFastQuery   = "fast-query"

	// FlagDeltaVersion specify the DeltaVersion
	FlagDeltaVersion = "delta-version"
)

var (
	// DeltaVersion do not apply delta if version does not match
	// if user specify the flag 'FlagDeltaVersion'(--delta-version) use user's setting,
	// otherwise use the default value
	DeltaVersion = 9
)

var (
	FastQuery     = false
	DownloadDelta = false
	UploadDelta   = false
)

type DeltasMessage struct {
	Metadata     []byte `json:"metadata"`
	MetadataHash []byte `json:"metadata_hash"`
	Height       int64  `json:"height"`
	CompressType int    `json:"compress_type"`
	From         string `json:"from"`
}

type DeltaPayload struct {
	ABCIRsp        []byte
	DeltasBytes    []byte
	WatchBytes     []byte
	WasmWatchBytes []byte
}

func (p DeltaPayload) AminoSize(_ *amino.Codec) int {
	var size int
	if len(p.ABCIRsp) != 0 {
		size += 1 + amino.ByteSliceSize(p.ABCIRsp)
	}
	if len(p.DeltasBytes) != 0 {
		size += 1 + amino.ByteSliceSize(p.DeltasBytes)
	}
	if len(p.WatchBytes) != 0 {
		size += 1 + amino.ByteSliceSize(p.WatchBytes)
	}
	return size
}

func (p DeltaPayload) MarshalAminoTo(_ *amino.Codec, buf *bytes.Buffer) error {
	var err error
	// field 1
	if len(p.ABCIRsp) != 0 {
		const pbKey = byte(1<<3 | amino.Typ3_ByteLength)
		err = amino.EncodeByteSliceWithKeyToBuffer(buf, p.ABCIRsp, pbKey)
		if err != nil {
			return err
		}
	}
	// field 2
	if len(p.DeltasBytes) != 0 {
		const pbKey = byte(2<<3 | amino.Typ3_ByteLength)
		err = amino.EncodeByteSliceWithKeyToBuffer(buf, p.DeltasBytes, pbKey)
		if err != nil {
			return err
		}
	}
	// field 3
	if len(p.WatchBytes) != 0 {
		const pbKey = byte(3<<3 | amino.Typ3_ByteLength)
		err = amino.EncodeByteSliceWithKeyToBuffer(buf, p.WatchBytes, pbKey)
		if err != nil {
			return err
		}
	}
	return nil
}

// Deltas defines the ABCIResponse and state delta
type Deltas struct {
	Height       int64
	Payload      DeltaPayload
	CompressType int
	CompressFlag int
	From         string

	marshalElapsed  time.Duration
	compressElapsed time.Duration
	hashElapsed     time.Duration
}

func (d *Deltas) AminoSize(cdc *amino.Codec) int {
	var size int
	if d.Height != 0 {
		size += 1 + amino.UvarintSize(uint64(d.Height))
	}
	payloadSize := d.Payload.AminoSize(cdc)
	if payloadSize > 0 {
		size += 1 + amino.UvarintSize(uint64(payloadSize)) + payloadSize
	}
	if d.CompressType != 0 {
		size += 1 + amino.UvarintSize(uint64(d.CompressType))
	}
	if d.CompressFlag != 0 {
		size += 1 + amino.UvarintSize(uint64(d.CompressFlag))
	}
	if d.From != "" {
		size += 1 + amino.EncodedStringSize(d.From)
	}
	return size
}

func (d *Deltas) MarshalAminoTo(cdc *amino.Codec, buf *bytes.Buffer) error {
	var err error
	// field 1
	if d.Height != 0 {
		const pbKey = byte(1<<3 | amino.Typ3_Varint)
		err = amino.EncodeUvarintWithKeyToBuffer(buf, uint64(d.Height), pbKey)
		if err != nil {
			return err
		}
	}
	// field 2
	payLoadSize := d.Payload.AminoSize(cdc)
	if payLoadSize > 0 {
		const pbKey = byte(2<<3 | amino.Typ3_ByteLength)
		buf.WriteByte(pbKey)
		err = amino.EncodeUvarintToBuffer(buf, uint64(payLoadSize))
		if err != nil {
			return err
		}
		lenBeforeData := buf.Len()
		err = d.Payload.MarshalAminoTo(cdc, buf)
		if err != nil {
			return err
		}
		if buf.Len()-lenBeforeData != payLoadSize {
			return amino.NewSizerError(payLoadSize, buf.Len()-lenBeforeData, payLoadSize)
		}
	}
	// field 3
	if d.CompressType != 0 {
		const pbKey = byte(3<<3 | amino.Typ3_Varint)
		err = amino.EncodeUvarintWithKeyToBuffer(buf, uint64(d.CompressType), pbKey)
		if err != nil {
			return err
		}
	}
	// field 4
	if d.CompressFlag != 0 {
		const pbKey = byte(4<<3 | amino.Typ3_Varint)
		err = amino.EncodeUvarintWithKeyToBuffer(buf, uint64(d.CompressFlag), pbKey)
		if err != nil {
			return err
		}
	}
	// field 5
	if d.From != "" {
		const pbKey = byte(5<<3 | amino.Typ3_ByteLength)
		err = amino.EncodeStringWithKeyToBuffer(buf, d.From, pbKey)
		if err != nil {
			return err
		}
	}
	return nil
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

func (d *Deltas) WasmWatchBytes() []byte {
	return d.Payload.WasmWatchBytes
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
		CompressType: d.CompressType,
		MetadataHash: payloadHash,
		From:         d.From,
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
	// uncompress
	d.CompressType = msg.CompressType
	msg.Metadata, err = compress.UnCompress(d.CompressType, msg.Metadata)
	if err != nil {
		return err
	}

	t2 := time.Now()
	// calc payload hash
	payloadHash := tmhash.Sum(msg.Metadata)
	if bytes.Compare(payloadHash, msg.MetadataHash) != 0 {
		return fmt.Errorf("metadata hash is different")
	}
	t3 := time.Now()

	err = cdc.UnmarshalBinaryBare(msg.Metadata, &d.Payload)
	t4 := time.Now()

	d.Height = msg.Height
	d.From = msg.From

	d.compressElapsed = t2.Sub(t1)
	d.hashElapsed = t3.Sub(t2)
	d.marshalElapsed = t4.Sub(t0) - d.compressElapsed - d.hashElapsed
	return err
}

func (d *Deltas) String() string {
	return fmt.Sprintf("height<%d>, size<%d>, from<%s>",
		d.Height,
		d.Size(),
		d.From,
	)
}

func (dds *Deltas) Validate(height int64) bool {
	if dds.Height != height ||
		len(dds.WatchBytes()) == 0 ||
		len(dds.ABCIRsp()) == 0 ||
		len(dds.DeltasBytes()) == 0 {
		return false
	}
	return true
}
