package types

import (
	"bytes"

	ethcmn "github.com/ethereum/go-ethereum/common"
)

const (
	Uint64Length = 8
)

func AppendBlockHashKey(blockHash []byte) []byte {
	return append(KeyPrefixBlockHash, blockHash...)
}

func AppendBloomKey(height int64) []byte {
	return append(KeyPrefixBloom, BloomKey(height)...)
}

func AppendHeightHashKey(height uint64) []byte {
	return append(KeyPrefixHeightHash, HeightKey(height)...)
}

func AppendBlockByHeightKey(height uint64) []byte {
	return append(KeyPrefixEthBlockByHeight, HeightKey(height)...)
}

func AppendBlockByHashKey(blockHash []byte) []byte {
	return append(KeyPrefixEthBlockByHash, blockHash...)
}

/*
 * IsKey
 */
func IsBlockHashKey(key []byte) bool {
	return bytes.HasPrefix(key, KeyPrefixBlockHash) &&
		len(key) == (len(KeyPrefixBlockHash)+ethcmn.HashLength)
}

func IsBloomKey(key []byte) bool {
	return bytes.HasPrefix(key, KeyPrefixBloom) &&
		len(key) == (len(KeyPrefixBloom)+Uint64Length)
}
