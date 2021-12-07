package types

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"reflect"
	"sync"
	"unsafe"
)

var EthCryptoState = sync.Pool{
	New: func() interface{} {
		return ethcrypto.NewKeccakState()
	},
}

var HashPool = sync.Pool{
	New: func() interface{} {
		return &ethcmn.Hash{}
	},
}

func PutHashToPool(hash *ethcmn.Hash) {
	HashPool.Put(hash)
}

var BytesPool = sync.Pool{
	New: func() interface{} {
		buf := make([]byte, 52)
		return &buf
	},
}

// StrToByte is meant to make a zero allocation conversion
// from string -> []byte to speed up operations, it is not meant
// to be used generally, but for a specific pattern to check for available
// keys within a domain.
func StrToByte(s string) []byte {
	var b []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	hdr.Cap = len(s)
	hdr.Len = len(s)
	hdr.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
	return b
}

// ByteSliceToStr is meant to make a zero allocation conversion
// from []byte -> string to speed up operations, it is not meant
// to be used generally, but for a specific pattern to delete keys
// from a map.
func ByteSliceToStr(b []byte) string {
	hdr := (*reflect.StringHeader)(unsafe.Pointer(&b))
	return *(*string)(unsafe.Pointer(hdr))
}

func UnsafeToString(bytes []byte) *string {
	hdr := &reflect.StringHeader{
		Data: uintptr(unsafe.Pointer(&bytes[0])),
		Len:  len(bytes),
	}
	return (*string)(unsafe.Pointer(hdr))
}
