package rootmulti

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/libs/iavl"
	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
)

// MarshalAppliedDeltaToAmino encode map[string]iavl.TreeDelta to []byte in amino format
func MarshalAppliedDeltaToAmino(appliedData map[string]iavl.TreeDelta) ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := byte(1<<3 | 2)
	if len(appliedData) == 0 {
		return buf.Bytes(), nil
	}

	// encode a pair of data one by one
	for k, v := range appliedData {
		err := buf.WriteByte(fieldKeysType)
		if err != nil {
			return nil, err
		}
		// map must convert to struct before it marshal
		data, err := newAppliedDelta(k, &v).MarshalToAmino()
		if err != nil {
			return nil, err
		}
		// write marshal result to buffer
		err = amino.EncodeByteSliceToBuffer(&buf, data)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}
func UnmarshalAppliedDeltaFromAmino(data []byte) (map[string]*iavl.TreeDelta, error) {
	var dataLen uint64 = 0
	var subData []byte
	appliedList := make(map[string]*iavl.TreeDelta)
	for {
		data = data[dataLen:]
		if len(data) == 0 {
			break
		}
		pos, pbType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return appliedList, err
		}
		data = data[1:]

		if pbType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, _ = amino.DecodeUvarint(data)

			data = data[n:]
			if len(data) < int(dataLen) {
				return appliedList, errors.New("not enough data")
			}
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			appliedData := new(appliedDelta)
			err := appliedData.UnmarshalFromAmino(subData)
			if err != nil {
				return nil, err
			}
			appliedList[appliedData.key] = appliedData.appliedTree

		default:
			return nil, fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return appliedList, nil
}

// appliedDelta convert map[string]*iavl.TreeDelta to struct
type appliedDelta struct {
	key         string
	appliedTree *iavl.TreeDelta
}

func newAppliedDelta(key string, treeValue *iavl.TreeDelta) *appliedDelta {
	return &appliedDelta{key: key, appliedTree: treeValue}
}
func (ad *appliedDelta) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(ad.key) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}

			err = amino.EncodeStringToBuffer(&buf, ad.key)
			if err != nil {
				return nil, err
			}
		case 2:
			data, err := ad.appliedTree.MarshalToAmino()
			if err != nil {
				return nil, err
			}
			if len(data) == 0 {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, data)
			if err != nil {
				return nil, err
			}

		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}
func (ad *appliedDelta) UnmarshalFromAmino(data []byte) error {
	var dataLen uint64 = 0
	var subData []byte

	for {
		data = data[dataLen:]
		if len(data) == 0 {
			break
		}
		pos, pbType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return err
		}
		data = data[1:]

		if pbType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, _ = amino.DecodeUvarint(data)

			data = data[n:]
			if len(data) < int(dataLen) {
				return errors.New("not enough data")
			}
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			ad.key = string(subData)
		case 2:
			err := ad.appliedTree.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}
