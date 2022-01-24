package rootmulti

import (
	"bytes"
	"fmt"
	"sort"

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

	// map is unsorted, so when the data isn't changed, we
	// must order it as the same way.
	keys := make([]string, 0)
	for k, _ := range appliedData {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// encode a pair of data one by one
	for _, k := range keys {
		err := buf.WriteByte(fieldKeysType)
		if err != nil {
			return nil, err
		}
		// map must copy to new struct before it marshal
		td := new(iavl.TreeDelta)
		*td = appliedData[k]
		data, err := newAppliedDelta(k, td).MarshalToAmino()
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

// UnmarshalAppliedDeltaFromAmino decode bytes to map[string]*iavl.TreeDelta in amino format.
func UnmarshalAppliedDeltaFromAmino(data []byte) (map[string]*iavl.TreeDelta, error) {
	var dataLen uint64 = 0
	var subData []byte
	appliedData := map[string]*iavl.TreeDelta{}

	for {
		data = data[dataLen:]
		if len(data) == 0 {
			break
		}
		pos, pbType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return map[string]*iavl.TreeDelta{}, err
		}
		data = data[1:]

		if pbType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, _ = amino.DecodeUvarint(data)

			data = data[n:]
			if len(data) < int(dataLen) {
				return map[string]*iavl.TreeDelta{}, errors.New("not enough data")
			}
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			ad := new(appliedDelta)
			err := ad.UnmarshalFromAmino(subData)
			if err != nil {
				return map[string]*iavl.TreeDelta{}, err
			}
			// if tree is empty, it must be initialized
			if ad.appliedTree == nil {
				ad.appliedTree = &iavl.TreeDelta{
					NodesDelta:         map[string]*iavl.NodeJson{},
					OrphansDelta:       make([]*iavl.NodeJson, 0),
					CommitOrphansDelta: map[string]int64{},
				}
			}
			appliedData[ad.key] = ad.appliedTree

		default:
			return map[string]*iavl.TreeDelta{}, fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return appliedData, nil
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
			appliedData := &iavl.TreeDelta{
				NodesDelta:         map[string]*iavl.NodeJson{},
				OrphansDelta:       []*iavl.NodeJson{},
				CommitOrphansDelta: map[string]int64{},
			}
			err := appliedData.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
			ad.appliedTree = appliedData

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}
