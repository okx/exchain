package iavl

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
)

type TreeDeltaMap map[string]*TreeDelta

// MarshalTreeDeltaMapToAmino encode map[string]*iavl.TreeDelta to []byte in amino format
func MarshalTreeDeltaMapToAmino(treeDeltaList TreeDeltaMap) ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := byte(1<<3 | 2)
	if len(treeDeltaList) == 0 {
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
	for k, _ := range treeDeltaList {
		err := buf.WriteByte(fieldKeysType)
		if err != nil {
			return nil, err
		}

		// map must copy to new struct before it marshal
		data, err := newAppliedDelta(k, treeDeltaList[k]).MarshalToAmino()
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

// UnmarshalTreeDeltaMapFromAmino decode bytes to map[string]*iavl.TreeDelta in amino format.
func UnmarshalTreeDeltaMapFromAmino(data []byte) (TreeDeltaMap, error) {
	var dataLen uint64 = 0
	var subData []byte
	treeDeltaList := map[string]*TreeDelta{}

	for {
		data = data[dataLen:]
		if len(data) == 0 {
			break
		}
		pos, pbType, err := amino.ParseProtoPosAndTypeMustOneByte(data[0])
		if err != nil {
			return map[string]*TreeDelta{}, err
		}
		data = data[1:]

		if pbType == amino.Typ3_ByteLength {
			var n int
			dataLen, n, _ = amino.DecodeUvarint(data)

			data = data[n:]
			if len(data) < int(dataLen) {
				return map[string]*TreeDelta{}, errors.New("not enough data")
			}
			subData = data[:dataLen]
		}

		switch pos {
		case 1:
			ad := new(appliedDelta)
			err := ad.UnmarshalFromAmino(subData)
			if err != nil {
				return map[string]*TreeDelta{}, err
			}
			treeDeltaList[ad.key] = ad.appliedTree

		default:
			return map[string]*TreeDelta{}, fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return treeDeltaList, nil
}

// appliedDelta convert map[string]*iavl.TreeDelta to struct
type appliedDelta struct {
	key         string
	appliedTree *TreeDelta
}

func newAppliedDelta(key string, treeValue *TreeDelta) *appliedDelta {
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
			appliedData := &TreeDelta{
				NodesDelta:         map[string]*NodeJson{},
				OrphansDelta:       []*NodeJson{},
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
			td.NodesDelta[nodeDelta.key] = nodeDelta.nodeData
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
			td.CommitOrphansDelta[commitDelta.key] = commitDelta.commitValue

		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

type nodesDelta struct {
	key      string
	nodeData *NodeJson
}

// newNodesDelta convert map[string]*nodejson to struct, and privode
// methods to marshal/unmarshal data in amino format.
func newNodesDelta(key string, nodeValue *NodeJson) *nodesDelta {
	return &nodesDelta{key: key, nodeData: nodeValue}
}
func (nd *nodesDelta) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(nd.key) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}

			err = amino.EncodeStringToBuffer(&buf, nd.key)
			if err != nil {
				return nil, err
			}
		case 2:
			data, err := nd.nodeData.MarshalToAmino()
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
			nd.key = string(subData)
		case 2:
			nodeData := new(NodeJson)
			err := nodeData.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
			nd.nodeData = nodeData
		default:
			return fmt.Errorf("unexpect feild num %d", pos)
		}
	}
	return nil
}

type commitOrphansDelta struct {
	key         string
	commitValue int64
}

// newNodesDelta convert map[string]int64 to struct, and privode
// methods to marshal/unmarshal data in amino format.
func newCommitOrphansDelta(key string, commitValue int64) *commitOrphansDelta {
	return &commitOrphansDelta{key: key, commitValue: commitValue}
}
func (cod *commitOrphansDelta) MarshalToAmino() ([]byte, error) {
	var buf bytes.Buffer
	//key type list
	fieldKeysType := [2]byte{1<<3 | 2, 2 << 3}

	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(cod.key) == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}

			err = amino.EncodeStringToBuffer(&buf, cod.key)
			if err != nil {
				return nil, err
			}
		case 2:
			if cod.commitValue == 0 {
				break
			}
			err := buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeUvarintToBuffer(&buf, uint64(cod.commitValue))
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
			cod.key = string(subData)
		case 2:
			value, n, err := amino.DecodeUvarint(data)
			if err != nil {
				return err
			}
			dataLen = uint64(n)
			cod.commitValue = int64(value)
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
	prePersisted bool   `json:"pre_persisted"`
}
