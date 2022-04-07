package types

import (
	ce "github.com/okex/exchain/libs/tendermint/crypto/encoding"
	"github.com/okex/exchain/libs/tendermint/libs/bits"
	tmbytes "github.com/okex/exchain/libs/tendermint/libs/bytes"
	tmproto "github.com/okex/exchain/libs/tendermint/proto/types"
	tmversion "github.com/okex/exchain/libs/tendermint/proto/version"
	"time"
)

// SignedHeader is a header along with the commits that prove it.
type IBCSignedHeader struct {
	*IBCHeader `json:"header"`

	Commit *IBCCommit `json:"commit"`
}

type IBCPartSetHeader struct {
	Total uint32           `json:"total"`
	Hash  tmbytes.HexBytes `json:"hash"`
}

type IBCBlockID struct {
	Hash          tmbytes.HexBytes `json:"hash"`
	PartSetHeader IBCPartSetHeader `json:"parts"`
}

func (b IBCBlockID) ToBlockID() BlockID {
	return BlockID{
		Hash: b.Hash,
		PartsHeader: PartSetHeader{
			Total: int(b.PartSetHeader.Total),
			Hash:  b.PartSetHeader.Hash,
		},
	}
}

type IBCHeader struct {
	// basic block info
	Version tmversion.Consensus `json:"version"`
	ChainID string              `json:"chain_id"`
	Height  int64               `json:"height"`
	Time    time.Time           `json:"time"`

	// prev block info
	LastBlockID IBCBlockID `json:"last_block_id"`

	// hashes of block data
	LastCommitHash tmbytes.HexBytes `json:"last_commit_hash"` // commit from validators from the last block
	DataHash       tmbytes.HexBytes `json:"data_hash"`        // transactions

	// hashes from the app output from the prev block
	ValidatorsHash     tmbytes.HexBytes `json:"validators_hash"`      // validators for the current block
	NextValidatorsHash tmbytes.HexBytes `json:"next_validators_hash"` // validators for the next block
	ConsensusHash      tmbytes.HexBytes `json:"consensus_hash"`       // consensus params for current block
	AppHash            tmbytes.HexBytes `json:"app_hash"`             // state after txs from the previous block
	// root hash of all results from the txs from the previous block
	// see `deterministicResponseDeliverTx` to understand which parts of a tx is hashed into here
	LastResultsHash tmbytes.HexBytes `json:"last_results_hash"`

	// consensus info
	EvidenceHash    tmbytes.HexBytes `json:"evidence_hash"`    // evidence included in the block
	ProposerAddress Address          `json:"proposer_address"` // original proposer of the block
}

type IBCCommit struct {
	// NOTE: The signatures are in order of address to preserve the bonded
	// ValidatorSet order.
	// Any peer with a block can gossip signatures by index with a peer without
	// recalculating the active ValidatorSet.
	Height     int64       `json:"height"`
	Round      int32       `json:"round"`
	BlockID    IBCBlockID  `json:"block_id"`
	Signatures []CommitSig `json:"signatures"`

	// Memoized in first call to corresponding method.
	// NOTE: can't memoize in constructor because constructor isn't used for
	// unmarshaling.
	hash     tmbytes.HexBytes
	bitArray *bits.BitArray
}

func (c *IBCCommit) ToCommit() *Commit {
	return &Commit{
		Height:     c.Height,
		Round:      int(c.Round),
		BlockID:    c.BlockID.ToBlockID(),
		Signatures: c.Signatures,
		hash:       c.hash,
		bitArray:   c.bitArray,
	}
}

func (v *Validator) HeightBytes(h int64) []byte {
	if HigherThanVenus1(h) {
		pk, err := ce.PubKeyToProto(v.PubKey)
		if err != nil {
			panic(err)
		}

		pbv := tmproto.SimpleValidator{
			PubKey:      &pk,
			VotingPower: v.VotingPower,
		}

		bz, err := pbv.Marshal()
		if err != nil {
			panic(err)
		}
		return bz
	}
	return v.OriginBytes()
}
