package coretypes

import (
	"github.com/okx/okbchain/libs/tendermint/types"
	"github.com/okx/okbchain/libs/tendermint/version"
)

type CM40ResultBlock struct {
	BlockID types.IBCBlockID `json:"block_id"`
	Block   *types.CM40Block `json:"block"`
}

func (c CM40ResultBlock) ToCM39ResultBlock() *ResultBlock {
	ret := ResultBlock{
		BlockID: c.BlockID.ToBlockID(),
		Block: &types.Block{
			Header:     c.Block.IBCHeader.ToCM39Header(),
			Data:       c.Block.Data,
			Evidence:   c.Block.Evidence,
			LastCommit: c.Block.LastCommit.ToCommit(),
		},
	}
	return &ret
}

// Commit and Header
type IBCResultCommit struct {
	types.IBCSignedHeader `json:"signed_header"`

	CanonicalCommit bool `json:"canonical"`
}

func NewIBCResultCommit(header *types.IBCHeader, commit *types.IBCCommit,
	canonical bool) *IBCResultCommit {

	return &IBCResultCommit{
		IBCSignedHeader: types.IBCSignedHeader{
			IBCHeader: header,
			Commit:    commit,
		},
		CanonicalCommit: canonical,
	}
}
func (c *IBCResultCommit) ToCommit() *ResultCommit {
	return &ResultCommit{
		SignedHeader: types.SignedHeader{
			Header: &types.Header{
				Version: version.Consensus{
					Block: version.Protocol(c.Version.Block),
					App:   version.Protocol(c.Version.App),
				},
				ChainID:            c.ChainID,
				Height:             c.Height,
				Time:               c.Time,
				LastBlockID:        c.LastBlockID.ToBlockID(),
				LastCommitHash:     c.LastCommitHash,
				DataHash:           c.DataHash,
				ValidatorsHash:     c.ValidatorsHash,
				NextValidatorsHash: c.NextValidatorsHash,
				ConsensusHash:      c.ConsensusHash,
				AppHash:            c.AppHash,
				LastResultsHash:    c.LastResultsHash,
				EvidenceHash:       c.EvidenceHash,
				ProposerAddress:    c.ProposerAddress,
			},
			Commit: c.Commit.ToCommit(),
		},
		CanonicalCommit: c.CanonicalCommit,
	}
}
