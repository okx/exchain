package types

import (
	"bytes"
	"fmt"
	"time"

	"github.com/okx/okbchain/libs/tendermint/crypto/tmhash"
	"github.com/okx/okbchain/libs/tendermint/libs/compress"
	"github.com/tendermint/go-amino"
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
	DeltaVersion = 10
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

func (m *DeltasMessage) AminoSize(_ *amino.Codec) int {
	var size int
	// field 1
	if len(m.Metadata) != 0 {
		size += 1 + amino.ByteSliceSize(m.Metadata)
	}
	// field 2
	if len(m.MetadataHash) != 0 {
		size += 1 + amino.ByteSliceSize(m.MetadataHash)
	}
	// field 3
	if m.Height != 0 {
		size += 1 + amino.UvarintSize(uint64(m.Height))
	}
	// field 4
	if m.CompressType != 0 {
		size += 1 + amino.UvarintSize(uint64(m.CompressType))
	}
	// field 5
	if m.From != "" {
		size += 1 + amino.EncodedStringSize(m.From)
	}
	return size
}

func (m *DeltasMessage) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(m.AminoSize(cdc))
	err := m.MarshalAminoTo(cdc, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (m *DeltasMessage) MarshalAminoTo(_ *amino.Codec, buf *bytes.Buffer) error {
	// field 1
	if len(m.Metadata) != 0 {
		const pbKey = 1<<3 | 2
		if err := amino.EncodeByteSliceWithKeyToBuffer(buf, m.Metadata, pbKey); err != nil {
			return err
		}
	}
	// field 2
	if len(m.MetadataHash) != 0 {
		const pbKey = 2<<3 | 2
		if err := amino.EncodeByteSliceWithKeyToBuffer(buf, m.MetadataHash, pbKey); err != nil {
			return err
		}
	}
	// field 3
	if m.Height != 0 {
		const pbKey = 3 << 3
		if err := amino.EncodeUvarintWithKeyToBuffer(buf, uint64(m.Height), pbKey); err != nil {
			return err
		}
	}
	// field 4
	if m.CompressType != 0 {
		const pbKey = 4 << 3
		if err := amino.EncodeUvarintWithKeyToBuffer(buf, uint64(m.CompressType), pbKey); err != nil {
			return err
		}
	}
	// field 5
	if m.From != "" {
		const pbKey = 5<<3 | 2
		if err := amino.EncodeStringWithKeyToBuffer(buf, m.From, pbKey); err != nil {
			return err
		}
	}
	return nil
}

func (m *DeltasMessage) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
	const fieldCount = 5
	var currentField int
	var currentType amino.Typ3
	var err error

	for cur := 1; cur <= fieldCount; cur++ {
		if len(data) != 0 && (currentField == 0 || currentField < cur) {
			var nextField int
			if nextField, currentType, err = amino.ParseProtoPosAndTypeMustOneByte(data[0]); err != nil {
				return err
			}
			if nextField < currentField {
				return fmt.Errorf("next field should greater than %d, got %d", currentField, nextField)
			} else {
				currentField = nextField
			}
		}

		if len(data) == 0 || currentField != cur {
			switch cur {
			case 1:
				m.Metadata = nil
			case 2:
				m.MetadataHash = nil
			case 3:
				m.Height = 0
			case 4:
				m.CompressType = 0
			case 5:
				m.From = ""
			default:
				return fmt.Errorf("unexpect feild num %d", cur)
			}
		} else {
			pbk := data[0]
			data = data[1:]
			var subData []byte
			if currentType == amino.Typ3_ByteLength {
				if subData, err = amino.DecodeByteSliceWithoutCopy(&data); err != nil {
					return err
				}
			}
			switch pbk {
			case 1<<3 | byte(amino.Typ3_ByteLength):
				amino.UpdateByteSlice(&m.Metadata, subData)
			case 2<<3 | byte(amino.Typ3_ByteLength):
				amino.UpdateByteSlice(&m.MetadataHash, subData)
			case 3<<3 | byte(amino.Typ3_Varint):
				if uvint, err := amino.DecodeUvarintUpdateBytes(&data); err != nil {
					return err
				} else {
					m.Height = int64(uvint)
				}
			case 4<<3 | byte(amino.Typ3_Varint):
				if m.CompressType, err = amino.DecodeIntUpdateBytes(&data); err != nil {
					return err
				}
			case 5<<3 | byte(amino.Typ3_ByteLength):
				m.From = string(subData)
			default:
				return fmt.Errorf("unexpect pb key %d", pbk)
			}
		}
	}

	if len(data) != 0 {
		return fmt.Errorf("unexpect data remain %X", data)
	}

	return nil
}

type DeltaPayload struct {
	ABCIRsp        []byte
	DeltasBytes    []byte
	WatchBytes     []byte
	WasmWatchBytes []byte
}

func (payload *DeltaPayload) AminoSize(_ *amino.Codec) int {
	var size int
	// field 1
	if len(payload.ABCIRsp) != 0 {
		size += 1 + amino.ByteSliceSize(payload.ABCIRsp)
	}
	// field 2
	if len(payload.DeltasBytes) != 0 {
		size += 1 + amino.ByteSliceSize(payload.DeltasBytes)
	}
	// field 3
	if len(payload.WatchBytes) != 0 {
		size += 1 + amino.ByteSliceSize(payload.WatchBytes)
	}
	// field 4
	if len(payload.WasmWatchBytes) != 0 {
		size += 1 + amino.ByteSliceSize(payload.WasmWatchBytes)
	}
	return size
}

func (payload *DeltaPayload) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(payload.AminoSize(cdc))
	err := payload.MarshalAminoTo(cdc, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (payload *DeltaPayload) MarshalAminoTo(_ *amino.Codec, buf *bytes.Buffer) error {
	// field 1
	if len(payload.ABCIRsp) != 0 {
		const pbKey = 1<<3 | 2
		if err := amino.EncodeByteSliceWithKeyToBuffer(buf, payload.ABCIRsp, pbKey); err != nil {
			return err
		}
	}
	// field 2
	if len(payload.DeltasBytes) != 0 {
		const pbKey = 2<<3 | 2
		if err := amino.EncodeByteSliceWithKeyToBuffer(buf, payload.DeltasBytes, pbKey); err != nil {
			return err
		}
	}
	// field 3
	if len(payload.WatchBytes) != 0 {
		const pbKey = 3<<3 | 2
		if err := amino.EncodeByteSliceWithKeyToBuffer(buf, payload.WatchBytes, pbKey); err != nil {
			return err
		}
	}
	// field 4
	if len(payload.WasmWatchBytes) != 0 {
		const pbKey = 4<<3 | 2
		if err := amino.EncodeByteSliceWithKeyToBuffer(buf, payload.WasmWatchBytes, pbKey); err != nil {
			return err
		}
	}
	return nil
}

func (payload *DeltaPayload) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
	const fieldCount = 4
	var currentField int
	var currentType amino.Typ3
	var err error

	for cur := 1; cur <= fieldCount; cur++ {
		if len(data) != 0 && (currentField == 0 || currentField < cur) {
			var nextField int
			if nextField, currentType, err = amino.ParseProtoPosAndTypeMustOneByte(data[0]); err != nil {
				return err
			}
			if nextField < currentField {
				return fmt.Errorf("next field should greater than %d, got %d", currentField, nextField)
			} else {
				currentField = nextField
			}
		}

		if len(data) == 0 || currentField != cur {
			switch cur {
			case 1:
				payload.ABCIRsp = nil
			case 2:
				payload.DeltasBytes = nil
			case 3:
				payload.WatchBytes = nil
			case 4:
				payload.WasmWatchBytes = nil
			default:
				return fmt.Errorf("unexpect feild num %d", cur)
			}
		} else {
			pbk := data[0]
			data = data[1:]
			var subData []byte
			if currentType == amino.Typ3_ByteLength {
				if subData, err = amino.DecodeByteSliceWithoutCopy(&data); err != nil {
					return err
				}
			}
			switch pbk {
			case 1<<3 | byte(amino.Typ3_ByteLength):
				if len(subData) != 0 {
					payload.ABCIRsp = make([]byte, len(subData))
					copy(payload.ABCIRsp, subData)
				} else {
					payload.ABCIRsp = nil
				}
			case 2<<3 | byte(amino.Typ3_ByteLength):
				if len(subData) != 0 {
					payload.DeltasBytes = make([]byte, len(subData))
					copy(payload.DeltasBytes, subData)
				} else {
					payload.DeltasBytes = nil
				}
			case 3<<3 | byte(amino.Typ3_ByteLength):
				if len(subData) != 0 {
					payload.WatchBytes = make([]byte, len(subData))
					copy(payload.WatchBytes, subData)
				} else {
					payload.WatchBytes = nil
				}
			case 4<<3 | byte(amino.Typ3_ByteLength):
				if len(subData) != 0 {
					payload.WasmWatchBytes = make([]byte, len(subData))
					copy(payload.WasmWatchBytes, subData)
				} else {
					payload.WasmWatchBytes = nil
				}
			default:
				return fmt.Errorf("unexpect pb key %d", pbk)
			}
		}
	}

	if len(data) != 0 {
		return fmt.Errorf("unexpect data remain %X", data)
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
	payload, err := d.Payload.MarshalToAmino(cdc)
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
	res, err := dt.MarshalToAmino(cdc)
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
	err := msg.UnmarshalFromAmino(cdc, bs)
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

	err = d.Payload.UnmarshalFromAmino(cdc, msg.Metadata)
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
		len(dds.ABCIRsp()) == 0 ||
		len(dds.DeltasBytes()) == 0 {
		return false
	}
	if FastQuery {
		if len(dds.WatchBytes()) == 0 {
			return false
		}
	}
	return true
}
