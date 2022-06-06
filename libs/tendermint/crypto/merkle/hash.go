package merkle

import (
	"bytes"
	"github.com/okex/exchain/libs/tendermint/crypto/tmhash"
	"sync"
)

// TODO: make these have a large predefined capacity
var (
	leafPrefix  = []byte{0}
	innerPrefix = []byte{1}

	hashBytesPool = &sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

// returns tmhash(0x00 || leaf)
func leafHash(leaf []byte) []byte {
	buf := hashBytesPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.Grow(len(leafPrefix) + len(leaf))
	buf.Write(leafPrefix)
	buf.Write(leaf)
	h := tmhash.Sum(buf.Bytes())
	hashBytesPool.Put(buf)
	return h
}

// returns tmhash(0x01 || left || right)
func innerHash(left []byte, right []byte) []byte {
	buf := hashBytesPool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.Grow(len(innerPrefix) + len(left) + len(right))
	buf.Write(innerPrefix)
	buf.Write(left)
	buf.Write(right)
	h := tmhash.Sum(buf.Bytes())
	hashBytesPool.Put(buf)
	return h
}
