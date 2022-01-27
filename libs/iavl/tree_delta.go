package iavl

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
)

type TreeDeltaMap map[string]*TreeDelta

// MarshalToAmino marshal to amino bytes
func (tdm TreeDeltaMap) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := byte(1<<3 | 2)

	if len(tdm) == 0 {
		return buf.Bytes(), nil
	}

	// encode a pair of data one by one
	for k := range tdm {
		err := buf.WriteByte(fieldKeysType)
		if err != nil {
			return nil, err
		}

		// map must convert to new struct before it marshal
		ti := &TreeDeltaMapImp{Key: k, TreeValue: tdm[k]}
		data, err := ti.MarshalToAmino(nil)
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

// UnmarshalFromAmino decode bytes from amino format.
func (tdm TreeDeltaMap) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			// a object to unmarshal data
			ad := &TreeDeltaMapImp{
				TreeValue: &TreeDelta{NodesDelta: map[string]*NodeJson{}, OrphansDelta: []*NodeJson{}, CommitOrphansDelta: map[string]int64{}},
			}
			err := ad.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}

			tdm[ad.Key] = ad.TreeValue

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

// TreeDeltaMapImp convert map[string]*TreeDelta to struct
type TreeDeltaMapImp struct {
	Key       string
	TreeValue *TreeDelta
}

//MarshalToAmino marshal data to amino bytes
func (ti *TreeDeltaMapImp) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(ti.Key) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}

			err = amino.EncodeStringToBuffer(&buf, ti.Key)
			if err != nil {
				return nil, err
			}
		case 2:
			data, err := ti.TreeValue.MarshalToAmino(nil)
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

// UnmarshalFromAmino unmarshal data from amino bytes.
func (ti *TreeDeltaMapImp) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			ti.Key = string(subData)

		case 2:
			// treevalue isn't empty,go on
			err := ti.TreeValue.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

// TreeDelta is the delta for applying on new version tree
type TreeDelta struct {
	NodesDelta         map[string]*NodeJson `json:"nodes_delta"`
	OrphansDelta       []*NodeJson          `json:"orphans_delta"`
	CommitOrphansDelta map[string]int64     `json:"commit_orphans_delta"`
}

func (td *TreeDelta) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [3]byte{1<<3 | 2, 2<<3 | 2, 3<<3 | 2}
	for pos := 1; pos <= 3; pos++ {
		switch pos {
		case 1:
			//encode data
			for k := range td.NodesDelta {
				err := buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				ni := &NodeJsonMapImp{Key: k, NodeValue: td.NodesDelta[k]}
				data, err := ni.MarshalToAmino(nil)
				if err != nil {
					return nil, err
				}
				err = amino.EncodeByteSliceToBuffer(&buf, data)
				if err != nil {
					return nil, err
				}
			}

		case 2:
			for _, v := range td.OrphansDelta {
				err := buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err := v.MarshalToAmino(nil)
				if err != nil {
					return nil, err
				}

				err = amino.EncodeByteSliceToBuffer(&buf, data)
				if err != nil {
					return nil, err
				}
			}
		case 3:
			for k := range td.CommitOrphansDelta {
				err := buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				cv := &CommitOrphansMapImp{Key: k, CommitValue: td.CommitOrphansDelta[k]}
				data, err := cv.MarshalToAmino(nil)
				if err != nil {
					return nil, err
				}
				err = amino.EncodeByteSliceToBuffer(&buf, data)
				if err != nil {
					return nil, err
				}
			}

		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}
func (td *TreeDelta) UnmarshalFromAmino(data []byte) error {
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
			ni := &NodeJsonMapImp{NodeValue: new(NodeJson)}
			err := ni.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}
			td.NodesDelta[ni.Key] = ni.NodeValue
		case 2:
			nodeData := new(NodeJson)
			err := nodeData.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}
			td.OrphansDelta = append(td.OrphansDelta, nodeData)
		case 3:
			ci := &CommitOrphansMapImp{}
			err := ci.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}
			// key and value must be existing together
			// CommitValue is initial in first.
			td.CommitOrphansDelta[ci.Key] = ci.CommitValue

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

type NodeJsonMapImp struct {
	Key       string
	NodeValue *NodeJson
}

// MarshalToAmino marshal data to amino bytes.
func (ni *NodeJsonMapImp) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(ni.Key) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}

			err = amino.EncodeStringToBuffer(&buf, ni.Key)
			if err != nil {
				return nil, err
			}
		case 2:
			data, err := ni.NodeValue.MarshalToAmino(nil)
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

// UnmarshalFromAmino unmarshal data from amino bytes.
func (ni *NodeJsonMapImp) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			ni.Key = string(subData)

		case 2:
			// NodeValue isn't empty, go on
			err := ni.NodeValue.UnmarshalFromAmino(nil, subData)
			if err != nil {
				return err
			}

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

type CommitOrphansMapImp struct {
	Key         string
	CommitValue int64
}

// MarshalToAmino marshal data to amino bytes.
func (ci *CommitOrphansMapImp) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	//key type list
	fieldKeysType := [2]byte{1<<3 | 2, 2 << 3}

	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(ci.Key) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}

			err = amino.EncodeStringToBuffer(&buf, ci.Key)
			if err != nil {
				return nil, err
			}

		case 2:
			if ci.CommitValue == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(ci.CommitValue))
			if err != nil {
				return nil, err
			}
		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalFromAmino unmarshal data from animo bytes.
func (ci *CommitOrphansMapImp) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			ci.Key = string(subData)

		case 2:
			value, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			ci.CommitValue = int64(value)

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

// NodeJson provide json Marshal of Node.
type NodeJson struct {
	Key          []byte `json:"key"`
	Value        []byte `json:"value"`
	Hash         []byte `json:"hash"`
	LeftHash     []byte `json:"left_hash"`
	RightHash    []byte `json:"right_hash"`
	Version      int64  `json:"version"`
	Size         int64  `json:"size"`
	Height       int8   `json:"height"`
	Persisted    bool   `json:"persisted"`
	PrePersisted bool   `json:"pre_persisted"`
}

// MarshalToAmino marshal data to amino bytes.
func (nj *NodeJson) MarshalToAmino(_ *amino.Codec) ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [10]byte{
		1<<3 | 2, 2<<3 | 2, 3<<3 | 2, 4<<3 | 2, 5<<3 | 2,
		6 << 3, 7 << 3, 8 << 3, 9 << 3, 10 << 3,
	}
	for pos := 1; pos <= 10; pos++ {
		switch pos {
		case 1:
			if len(nj.Key) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, nj.Key)
			if err != nil {
				return nil, err
			}
		case 2:
			if len(nj.Value) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, nj.Value)
			if err != nil {
				return nil, err
			}
		case 3:
			if len(nj.Hash) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, nj.Hash)
			if err != nil {
				return nil, err
			}
		case 4:
			if len(nj.LeftHash) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, nj.LeftHash)
			if err != nil {
				return nil, err
			}
		case 5:
			if len(nj.RightHash) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, nj.RightHash)
			if err != nil {
				return nil, err
			}
		case 6:
			if nj.Version == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(nj.Version))
			if err != nil {
				return nil, err
			}
		case 7:
			if nj.Size == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(nj.Size))
			if err != nil {
				return nil, err
			}
		case 8:
			if nj.Height == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(nj.Height))
			if err != nil {
				return nil, err
			}
		case 9:
			if !nj.Persisted {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = buf.WriteByte(1)
			if err != nil {
				return nil, err
			}
		case 10:
			if !nj.PrePersisted {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = buf.WriteByte(1)
			if err != nil {
				return nil, err
			}

		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}

// UnmarshalFromAmino unmarshal data from amino bytes.
func (nj *NodeJson) UnmarshalFromAmino(_ *amino.Codec, data []byte) error {
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
			nj.Key = make([]byte, len(subData))
			copy(nj.Key, subData)

		case 2:
			nj.Value = make([]byte, len(subData))
			copy(nj.Value, subData)

		case 3:
			nj.Hash = make([]byte, len(subData))
			copy(nj.Hash, subData)

		case 4:
			nj.LeftHash = make([]byte, len(subData))
			copy(nj.LeftHash, subData)

		case 5:
			nj.RightHash = make([]byte, len(subData))
			copy(nj.RightHash, subData)

		case 6:
			value, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			nj.Version = int64(value)

		case 7:
			value, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			nj.Size = int64(value)

		case 8:
			value, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			nj.Height = int8(value)

		case 9:
			if data[0] != 0 && data[0] != 1 {
				return fmt.Errorf("invalid Persisted")
			}
			nj.Persisted = data[0] == 1
			dataLen = 1

		case 10:
			if data[0] != 0 && data[0] != 1 {
				return fmt.Errorf("invalid prePersisted")
			}
			nj.PrePersisted = data[0] == 1
			dataLen = 1

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}
