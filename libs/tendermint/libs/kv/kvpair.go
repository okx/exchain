package kv

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/pkg/errors"

	"github.com/tendermint/go-amino"
)

//----------------------------------------
// KVPair

/*
Defined in types.proto

type Pair struct {
	Key   []byte
	Value []byte
}
*/

type Pairs []Pair

// Sorting
func (kvs Pairs) Len() int { return len(kvs) }
func (kvs Pairs) Less(i, j int) bool {
	switch bytes.Compare(kvs[i].Key, kvs[j].Key) {
	case -1:
		return true
	case 0:
		return bytes.Compare(kvs[i].Value, kvs[j].Value) < 0
	case 1:
		return false
	default:
		panic("invalid comparison result")
	}
}
func (kvs Pairs) Swap(i, j int) { kvs[i], kvs[j] = kvs[j], kvs[i] }
func (kvs Pairs) Sort()         { sort.Sort(kvs) }

func (pair *Pair) MarshalToAmino(cdc *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(pair.AminoSize(cdc))
	err := pair.MarshalAminoTo(cdc, &buf)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (pair *Pair) MarshalAminoTo(_ *amino.Codec, buf *bytes.Buffer) error {
	// field 1
	if len(pair.Key) != 0 {
		const pbKey = 1<<3 | 2
		err := amino.EncodeByteSliceWithKeyToBuffer(buf, pair.Key, pbKey)
		if err != nil {
			return err
		}
	}

	// field 2
	if len(pair.Value) != 0 {
		const pbKey = 2<<3 | 2
		err := amino.EncodeByteSliceWithKeyToBuffer(buf, pair.Value, pbKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (pair *Pair) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
	var dataLen uint64 = 0
	var subData []byte

	for {
		data = data[dataLen:]

		if len(data) == 0 {
			break
		}

		pos, aminoType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return err
		}
		if aminoType != amino.Typ3_ByteLength {
			return errors.New("invalid amino type")
		}
		data = data[1:]

		var n int
		dataLen, n, err = amino.DecodeUvarint(data)
		if err != nil {
			return err
		}

		data = data[n:]
		if len(data) < int(dataLen) {
			return errors.New("invalid data length")
		}
		subData = data[:dataLen]

		switch pos {
		case 1:
			pair.Key = make([]byte, dataLen)
			copy(pair.Key, subData)
		case 2:
			pair.Value = make([]byte, dataLen)
			copy(pair.Value, subData)
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

func (pair Pair) AminoSize(_ *amino.Codec) int {
	var size = 0
	if len(pair.Key) != 0 {
		size += 1 + amino.ByteSliceSize(pair.Key)
	}
	if len(pair.Value) != 0 {
		size += 1 + amino.ByteSliceSize(pair.Value)
	}
	return size
}

func (m *Pair) GetIndex() bool {
	if m != nil {
		return true
	}
	return false
}
