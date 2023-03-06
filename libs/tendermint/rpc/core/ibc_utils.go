package core

import (
	"github.com/okx/okbchain/libs/tendermint/proto/version"
	coretypes "github.com/okx/okbchain/libs/tendermint/rpc/core/types"
	"github.com/okx/okbchain/libs/tendermint/types"
)

func ConvBlock2CM40Block(r *types.Block) *types.CM40Block {
	ret := &types.CM40Block{
		IBCHeader:  *ConvHeadersToIbcHeader(&r.Header),
		Data:       r.Data,
		Evidence:   r.Evidence,
		LastCommit: ConvCommitToIBCCommit(&r.Header, r.LastCommit),
	}
	return ret
}

func ConvBlockID2CM40BlockID(r types.BlockID) types.IBCBlockID {
	return types.IBCBlockID{
		Hash: r.Hash,
		PartSetHeader: types.IBCPartSetHeader{
			Total: uint32(r.PartsHeader.Total),
			Hash:  r.PartsHeader.Hash,
		},
	}
}

func ConvResultCommitTOIBC(r *coretypes.ResultCommit) *coretypes.IBCResultCommit {
	v := ConvSignheaderToIBCSignHeader(&r.SignedHeader)
	ret := &coretypes.IBCResultCommit{
		IBCSignedHeader: *v,
		CanonicalCommit: r.CanonicalCommit,
	}
	return ret
}

func ConvSignheaderToIBCSignHeader(h *types.SignedHeader) *types.IBCSignedHeader {
	ret := &types.IBCSignedHeader{
		IBCHeader: ConvHeadersToIbcHeader(h.Header),
		Commit:    ConvCommitToIBCCommit(h.Header, h.Commit),
	}

	return ret
}
func ConvHeadersToIbcHeader(h *types.Header) *types.IBCHeader {
	ret := &types.IBCHeader{
		Version: version.Consensus{
			Block: uint64(h.Version.Block),
			App:   uint64(h.Version.App),
		},
		ChainID: h.ChainID,
		Height:  h.Height,
		Time:    h.Time,
		LastBlockID: types.IBCBlockID{
			// TODO
			Hash: h.LastBlockID.Hash,
			PartSetHeader: types.IBCPartSetHeader{
				Total: uint32(h.LastBlockID.PartsHeader.Total),
				Hash:  h.LastBlockID.PartsHeader.Hash,
			},
		},
		LastCommitHash:     h.LastCommitHash,
		DataHash:           h.DataHash,
		ValidatorsHash:     h.ValidatorsHash,
		NextValidatorsHash: h.NextValidatorsHash,
		ConsensusHash:      h.ConsensusHash,
		AppHash:            h.AppHash,
		LastResultsHash:    h.LastResultsHash,
		EvidenceHash:       h.EvidenceHash,
		ProposerAddress:    h.ProposerAddress,
	}

	return ret
}

func ConvCommitToIBCCommit(hh *types.Header, h *types.Commit) *types.IBCCommit {
	ret := &types.IBCCommit{
		Height: h.Height,
		Round:  int32(h.Round),
		BlockID: types.IBCBlockID{
			Hash: h.BlockID.Hash,
			PartSetHeader: types.IBCPartSetHeader{
				Total: uint32(h.BlockID.PartsHeader.Total),
				Hash:  h.BlockID.PartsHeader.Hash,
			},
		},
		Signatures: h.Signatures,
	}

	return ret
}

func ConvEventBlock2CM40Event(block types.EventDataNewBlock) types.CM40EventDataNewBlock {
	ret := types.CM40EventDataNewBlock{}
	ret.Block = ConvBlock2CM40Block(block.Block)
	ret.ResultBeginBlock = block.ResultBeginBlock
	ret.ResultEndBlock = block.ResultEndBlock
	return ret
}
