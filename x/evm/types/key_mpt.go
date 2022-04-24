package types

import (
	"bytes"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
)

const (
	Uint64Length = 8
)

// Below are the keys which are different from the key in iavl
var (
	UpgradedKeyPrefixCode = rawdb.CodePrefix // Old: KeyPrefixCode = []byte{0x04}
)

/*
 * KeyPrefixBlockHash         = []byte{0x01}
 * KeyPrefixBloom             = []byte{0x02}
 * UpgradedKeyPrefixCode      = []byte{"c"}
 * KeyPrefixStorage           not stored in db directly
 * KeyPrefixChainConfig       = []byte{0x06}
 * KeyPrefixHeightHash        = []byte{0x07}
 *
 * Below are functions used for setting in DiskDB
 */
/*
 * Append
 */
func AppendBlockHashKey(blockHash []byte) []byte {
	return append(KeyPrefixBlockHash, blockHash...)
}

func AppendBloomKey(height int64) []byte {
	return append(KeyPrefixBloom, BloomKey(height)...)
}

func AppendUpgradedCodeKey(codeHash []byte) []byte {
	return append(UpgradedKeyPrefixCode, codeHash...)
}

func AppendHeightHashKey(height uint64) []byte {
	return append(KeyPrefixHeightHash, HeightHashKey(height)...)
}

/*
 * Split
 */
func SplitUpgradedCodeHashKey(key []byte) []byte {
	return key[len(UpgradedKeyPrefixCode):]
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

func IsCodeHashKey(key []byte) bool {
	return bytes.HasPrefix(key, UpgradedKeyPrefixCode) &&
		len(key) == (len(UpgradedKeyPrefixCode)+ethcmn.HashLength)
}

func IsHeightHashKey(key []byte) bool {
	return bytes.HasPrefix(key, KeyPrefixHeightHash) &&
		len(key) == (len(KeyPrefixHeightHash)+Uint64Length)
}
