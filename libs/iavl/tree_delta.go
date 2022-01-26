package iavl

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
)

type TreeDeltaMap map[string]*TreeDelta

// MarshalToAmino encode to amino bytes
func (tdm TreeDeltaMap) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := byte(1<<3 | 2)

	if len(tdm) == 0 {
		return buf.Bytes(), nil
	}

	//// map is unsorted, so when the data isn't changed, we
	//// must order it as the same way.
	//keys := make([]string, 0)
	//for k, _ := range appliedData {
	//	keys = append(keys, k)
	//}
	//sort.Strings(keys)

	// encode a pair of data one by one
	for k, _ := range tdm {
		err := buf.WriteByte(fieldKeysType)
		if err != nil {
			return nil, err
		}

		// map must copy to new struct before it marshal
		data, err := newAppliedDelta(k, tdm[k]).MarshalToAmino()
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
func (tdm TreeDeltaMap) UnmarshalFromAmino(data []byte) error {
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
			ad := new(appliedDelta)
			err := ad.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
			tdm[ad.Key] = ad.AppliedTree

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

// appliedDelta convert map[string]*iavl.TreeDelta to struct
type appliedDelta struct {
	Key         string
	AppliedTree *TreeDelta
}

func newAppliedDelta(key string, treeValue *TreeDelta) *appliedDelta {
	return &appliedDelta{Key: key, AppliedTree: treeValue}
}
func (ad *appliedDelta) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(ad.Key) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}

			err = amino.EncodeStringToBuffer(&buf, ad.Key)
			if err != nil {
				return nil, err
			}
		case 2:
			data, err := ad.AppliedTree.MarshalToAmino()
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
			ad.Key = string(subData)
		case 2:
			appliedData := &TreeDelta{
				NodesDelta:         map[string]*NodeJson{},
				OrphansDelta:       []*NodeJson{},
				CommitOrphansDelta: map[string]int64{},
			}
			err := appliedData.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
			ad.AppliedTree = appliedData

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

func (td *TreeDelta) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [3]byte{1<<3 | 2, 2<<3 | 2, 3<<3 | 2}
	for pos := 1; pos <= 3; pos++ {
		switch pos {
		case 1:
			////sort keys
			//keys := []string{}
			//for k := range td.NodesDelta {
			//	keys = append(keys, k)
			//}
			//sort.Strings(keys)

			//encode data after it is sorted
			for k, _ := range td.NodesDelta {
				err := buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err := newNodesDelta(k, td.NodesDelta[k]).MarshalToAmino()
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
				data, err := v.MarshalToAmino()
				if err != nil {
					return nil, err
				}

				err = amino.EncodeByteSliceToBuffer(&buf, data)
				if err != nil {
					return nil, err
				}
			}
		case 3:
			////sort keys
			//keys := []string{}
			//for k := range td.NodesDelta {
			//	keys = append(keys, k)
			//}
			//sort.Strings(keys)

			for k, v := range td.CommitOrphansDelta {
				err := buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err := newCommitOrphansDelta(k, v).MarshalToAmino()
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
			nodeDelta := new(nodesDelta)
			err := nodeDelta.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
			td.NodesDelta[nodeDelta.Key] = nodeDelta.NodeData
		case 2:
			nodeData := new(NodeJson)
			err := nodeData.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
			td.OrphansDelta = append(td.OrphansDelta, nodeData)
		case 3:
			commitDelta := new(commitOrphansDelta)
			err := commitDelta.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
			td.CommitOrphansDelta[commitDelta.Key] = commitDelta.CommitValue

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

type nodesDelta struct {
	Key      string
	NodeData *NodeJson
}

// newNodesDelta convert map[string]*nodeJson to struct, and provide
// methods to marshal/unmarshal data in amino format.
func newNodesDelta(key string, nodeValue *NodeJson) *nodesDelta {
	return &nodesDelta{Key: key, NodeData: nodeValue}
}
func (nd *nodesDelta) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(nd.Key) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}

			err = amino.EncodeStringToBuffer(&buf, nd.Key)
			if err != nil {
				return nil, err
			}
		case 2:
			data, err := nd.NodeData.MarshalToAmino()
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
func (nd *nodesDelta) UnmarshalFromAmino(data []byte) error {
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
			nd.Key = string(subData)
		case 2:
			nodeData := new(NodeJson)
			err := nodeData.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
			nd.NodeData = nodeData
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

type commitOrphansDelta struct {
	Key         string
	CommitValue int64
}

// newNodesDelta convert map[string]int64 to struct, and provide
// methods to marshal/unmarshal data in amino format.
func newCommitOrphansDelta(key string, commitValue int64) *commitOrphansDelta {
	return &commitOrphansDelta{Key: key, CommitValue: commitValue}
}
func (cod *commitOrphansDelta) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	//key type list
	fieldKeysType := [2]byte{1<<3 | 2, 2 << 3}

	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(cod.Key) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}

			err = amino.EncodeStringToBuffer(&buf, cod.Key)
			if err != nil {
				return nil, err
			}
		case 2:
			if cod.CommitValue == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(cod.CommitValue))
			if err != nil {
				return nil, err
			}
		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}
func (cod *commitOrphansDelta) UnmarshalFromAmino(data []byte) error {
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
			cod.Key = string(subData)
		case 2:
			value, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			cod.CommitValue = int64(value)
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

func (nj *NodeJson) MarshalToAmino() ([]byte, error) {
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

func (nj *NodeJson) UnmarshalFromAmino(data []byte) error {
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
