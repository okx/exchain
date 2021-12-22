package kv

import (
	"bytes"
	"sort"

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

func MarshalPairToAmino(pair Pair) ([]byte, error) {
	var buf bytes.Buffer
	fieldKeysType := [2]byte{1<<3 | 2, 2<<3 | 2}
	var err error
	for pos := 1; pos <= 2; pos++ {
		switch pos {
		case 1:
			if len(pair.Key) == 0 {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, pair.Key)
			if err != nil {
				return nil, err
			}
		case 2:
			if len(pair.Value) == 0 {
				break
			}
			err = buf.WriteByte(fieldKeysType[pos-1])
			if err != nil {
				return nil, err
			}
			err = amino.EncodeByteSliceToBuffer(&buf, pair.Value)
			if err != nil {
				return nil, err
			}
		default:
			panic("unreachable")
		}
	}
	return buf.Bytes(), nil
}
