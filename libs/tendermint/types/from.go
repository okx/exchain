package types

import (
	"github.com/okx/okbchain/libs/tendermint/proto/version"
)

func (e CM40EventDataNewBlock) From(block EventDataNewBlock) CM40EventDataNewBlock {
	e.Block = &CM40Block{}
	cm40Block := block.Block.To()
	e.Block = &cm40Block
	e.ResultEndBlock = block.ResultEndBlock
	e.ResultBeginBlock = block.ResultBeginBlock
	return e
}

func (b CM40Block) From(block Block) CM40Block {
	b.Data = block.Data
	b.Evidence = block.Evidence
	cmt := block.LastCommit.To()
	b.LastCommit = &cmt
	b.IBCHeader = b.IBCHeader.From(block.Header)
	return b
}

func (h IBCHeader) From(header Header) IBCHeader {
	h.Version = version.Consensus{
		Block: uint64(header.Version.Block),
		App:   uint64(header.Version.App),
	}
	h.ChainID = header.ChainID
	h.Height = header.Height
	h.Time = header.Time
	h.LastBlockID = h.LastBlockID.From(header.LastBlockID)
	h.LastCommitHash = header.LastCommitHash
	h.DataHash = header.DataHash
	h.ValidatorsHash = header.ValidatorsHash
	h.NextValidatorsHash = header.NextValidatorsHash
	h.ConsensusHash = header.ConsensusHash
	h.AppHash = header.AppHash
	h.LastResultsHash = header.LastResultsHash
	h.EvidenceHash = header.EvidenceHash
	h.ProposerAddress = header.ProposerAddress
	return h
}

func (p IBCPartSetHeader) From(header PartSetHeader) IBCPartSetHeader {
	p.Total = uint32(header.Total)
	p.Hash = header.Hash
	return p
}

func (b IBCBlockID) From(bb BlockID) IBCBlockID {
	b.PartSetHeader = b.PartSetHeader.From(bb.PartsHeader)
	b.Hash = bb.Hash
	return b
}

func (c IBCCommit) From(cc Commit) IBCCommit {
	c.BlockID = c.BlockID.From(cc.BlockID)
	c.Height = cc.Height
	c.Round = int32(cc.Round)
	c.Signatures = cc.Signatures
	return c
}
