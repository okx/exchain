package iavl

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"github.com/tendermint/go-amino"
)

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
			if len(td.NodesDelta) == 0 {
				break
			}
			for k, v := range td.NodesDelta {
				err := buf.WriteByte(fieldKeysType[pos-1])
				if err != nil {
					return nil, err
				}
				data, err := newNodesDelta(k, v).MarshalToAmino()
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
			err := nd.nodeData.UnmarshalFromAmino(subData)
			if err != nil {
				return err
			}
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
