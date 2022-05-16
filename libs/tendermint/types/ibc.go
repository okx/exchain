package types

import (
	"github.com/okex/exchain/libs/tendermint/libs/protoio"
	tmproto "github.com/okex/exchain/libs/tendermint/proto/types"
)

// CanonicalizeVote transforms the given Proposal to a CanonicalProposal.
func IBCCanonicalizeProposal(chainID string, proposal *Proposal) tmproto.CanonicalProposal {
	return tmproto.CanonicalProposal{
		Type:      tmproto.ProposalType,
		Height:    proposal.Height,       // encoded as sfixed64
		Round:     int64(proposal.Round), // encoded as sfixed64
		POLRound:  int64(proposal.POLRound),
		BlockID:   IBCCanonicalizeBlockID(&proposal.BlockID),
		Timestamp: proposal.Timestamp,
		ChainID:   chainID,
	}
}

func IBCCanonicalizeBlockID(rbid *BlockID) *tmproto.CanonicalBlockID {
	var cbid *tmproto.CanonicalBlockID
	if rbid == nil || rbid.IsZero() {
		cbid = nil
	} else {
		cbid = &tmproto.CanonicalBlockID{
			Hash:          rbid.Hash,
			PartSetHeader: IBCCanonicalizePartSetHeader(rbid.PartsHeader),
		}
	}

	return cbid
}

// CanonicalizeVote transforms the given PartSetHeader to a CanonicalPartSetHeader.
func IBCCanonicalizePartSetHeader(psh PartSetHeader) tmproto.CanonicalPartSetHeader {
	pp := psh.ToIBCProto()
	return tmproto.CanonicalPartSetHeader{
		Total: uint32(pp.Total),
		Hash:  pp.Hash,
	}
}
func ProposalSignBytes(chainID string, p *Proposal) []byte {
	pb := IBCCanonicalizeProposal(chainID, p)
	bz, err := protoio.MarshalDelimited(&pb)
	if err != nil {
		panic(err)
	}

	return bz
}
func VoteSignBytes(chainID string, vote *Vote) []byte {
	pb := IBCCanonicalizeVote(chainID, vote)
	bz, err := protoio.MarshalDelimited(&pb)
	if err != nil {
		panic(err)
	}
	return bz
}
func IBCCanonicalizeVote(chainID string, vote *Vote) tmproto.CanonicalVote {
	return tmproto.CanonicalVote{
		Type:      tmproto.SignedMsgType(vote.Type),
		Height:    vote.Height,       // encoded as sfixed64
		Round:     int64(vote.Round), // encoded as sfixed64
		BlockID:   IBCCanonicalizeBlockID(&vote.BlockID),
		Timestamp: vote.Timestamp,
		ChainID:   chainID,
	}
}
