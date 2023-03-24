package types

import (
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Block represents a transaction returned to RPC clients.
type Block struct {
	Number           hexutil.Uint64 `json:"number"`
	Hash             common.Hash    `json:"hash"`
	ParentHash       common.Hash    `json:"parentHash"`
	Nonce            BlockNonce     `json:"nonce"`
	UncleHash        common.Hash    `json:"sha3Uncles"`
	LogsBloom        ethtypes.Bloom `json:"logsBloom"`
	TransactionsRoot common.Hash    `json:"transactionsRoot"`
	StateRoot        common.Hash    `json:"stateRoot"`
	Miner            common.Address `json:"miner"`
	MixHash          common.Hash    `json:"mixHash"`
	Difficulty       hexutil.Uint64 `json:"difficulty"`
	TotalDifficulty  hexutil.Uint64 `json:"totalDifficulty"`
	ExtraData        hexutil.Bytes  `json:"extraData"`
	Size             hexutil.Uint64 `json:"size"`
	GasLimit         hexutil.Uint64 `json:"gasLimit"`
	GasUsed          *hexutil.Big   `json:"gasUsed"`
	Timestamp        hexutil.Uint64 `json:"timestamp"`
	Uncles           []common.Hash  `json:"uncles"`
	ReceiptsRoot     common.Hash    `json:"receiptsRoot"`
	Transactions     interface{}    `json:"transactions"`
}

// EthHash returns block hash encode by rlp for being compatible with ethereum
func (b *Block) EthHash() common.Hash {
	var enc ethtypes.Header
	enc.ParentHash = b.ParentHash
	enc.UncleHash = b.UncleHash
	enc.Coinbase = b.Miner
	enc.Root = b.StateRoot
	enc.TxHash = b.TransactionsRoot
	enc.ReceiptHash = b.ReceiptsRoot
	enc.Bloom = b.LogsBloom
	enc.Difficulty = big.NewInt(int64(b.Difficulty))
	enc.Number = big.NewInt(int64(b.Number))
	enc.GasLimit = uint64(b.GasLimit)
	enc.GasUsed = b.GasUsed.ToInt().Uint64()
	enc.Time = uint64(b.Timestamp)
	enc.Extra = b.ExtraData
	enc.MixDigest = b.MixHash
	enc.Nonce = ethtypes.BlockNonce(b.Nonce)
	return rlpHash(&enc)
}

// A BlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type BlockNonce [8]byte

// EncodeNonce converts the given integer to a block nonce.
func EncodeNonce(i uint64) BlockNonce {
	var n BlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n BlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockNonce", input, n[:])
}
