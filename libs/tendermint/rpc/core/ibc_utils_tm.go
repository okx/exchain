package core
//
//import (
//	"github.com/okex/exchain/libs/tendermint/proto/version"
//	coretypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
//	"github.com/okex/exchain/libs/tendermint/types"
//	tmbytes "github.com/tendermint/tendermint/libs/bytes"
//	tmversion "github.com/tendermint/tendermint/proto/tendermint/version"
//	coretypes2 "github.com/tendermint/tendermint/rpc/core/types"
//	types2 "github.com/tendermint/tendermint/types"
//)
//
////func ConvResultCommitTOIBC(r *coretypes.ResultCommit) *coretypes.IBCResultCommit {
////	v := ConvSignheaderToIBCSignHeader(&r.SignedHeader)
////	ret := &coretypes.IBCResultCommit{
////		IBCSignedHeader: *v,
////		CanonicalCommit: r.CanonicalCommit,
////	}
////
////	r.Commit.VoteSignBytes()
////	//for _,commit:=range v.Commit.Signatures{
////	//
////	//	voteSignBytes := commit.VoteSignBytes(chainID, int32(idx))
////	//}
////	if hhash, chash := r.Header.IBCHash(), ret.Commit.BlockID.Hash; !bytes.Equal(hhash, chash) {
////		panic("asd")
////	}
////
////	return ret
////}
//
////func VoteSignBytes(chainID string, vote *tmproto.Vote) []byte {
////	pb := CanonicalizeVote(chainID, vote)
////	bz, err := protoio.MarshalDelimited(&pb)
////	if err != nil {
////		panic(err)
////	}
////
////	return bz
////}
////func CanonicalizeVote(chainID string, vote *tmproto.Vote) tmproto.CanonicalVote {
////	return tmproto2.CanonicalVote{
////		Type:      tmproto2.SignedMsgType(vote.Type),
////		Height:    vote.Height,       // encoded as sfixed64
////		Round:     int64(vote.Round), // encoded as sfixed64
////		BlockID:   CanonicalizeBlockID(vote.BlockID),
////		Timestamp: vote.Timestamp,
////		ChainID:   chainID,
////	}
////}
////
////func CanonicalizeBlockID(bid tmproto.BlockID) *tmproto.CanonicalBlockID {
////	rbid, err := BlockIDFromProto(&bid)
////	if err != nil {
////		panic(err)
////	}
////	var cbid *tmproto.CanonicalBlockID
////	if rbid == nil || rbid.IsZero() {
////		cbid = nil
////	} else {
////		cbid = &tmproto.CanonicalBlockID{
////			Hash:          bid.Hash,
////			PartSetHeader: CanonicalizePartSetHeader(bid.PartSetHeader),
////		}
////	}
////
////	return cbid
////}
//
//func ConvResultCommitToTendermint(r *coretypes.ResultCommit) *coretypes2.ResultCommit {
//	ss := make([]types2.CommitSig, 0)
//	for _, v := range r.Commit.Signatures {
//		ss = append(ss, types2.CommitSig{
//			BlockIDFlag:      types2.BlockIDFlag(v.BlockIDFlag),
//			ValidatorAddress: types2.Address(v.ValidatorAddress),
//			Timestamp:        v.Timestamp,
//			Signature:        v.Signature,
//		})
//	}
//	rett := coretypes2.ResultCommit{
//		SignedHeader: types2.SignedHeader{
//			Header: &types2.Header{
//				Version: tmversion.Consensus{
//					Block: uint64(r.Header.Version.Block),
//					App:   uint64(r.Header.Version.App),
//				},
//				ChainID: r.SignedHeader.ChainID,
//				Height:  r.SignedHeader.Height,
//				Time:    r.SignedHeader.Time,
//				LastBlockID: types2.BlockID{
//					Hash: tmbytes.HexBytes(r.SignedHeader.LastBlockID.Hash),
//					PartSetHeader: types2.PartSetHeader{
//						Total: uint32(r.LastBlockID.PartsHeader.Total),
//						Hash:  tmbytes.HexBytes(r.LastBlockID.PartsHeader.Hash),
//					},
//				},
//				LastCommitHash:     tmbytes.HexBytes(r.Header.LastCommitHash),
//				DataHash:           tmbytes.HexBytes(r.Header.DataHash),
//				ValidatorsHash:     tmbytes.HexBytes(r.Header.ValidatorsHash),
//				NextValidatorsHash: tmbytes.HexBytes(r.Header.NextValidatorsHash),
//				ConsensusHash:      tmbytes.HexBytes(r.Header.ConsensusHash),
//				AppHash:            tmbytes.HexBytes(r.Header.AppHash),
//				LastResultsHash:    tmbytes.HexBytes(r.Header.LastResultsHash),
//				EvidenceHash:       tmbytes.HexBytes(r.Header.EvidenceHash),
//				ProposerAddress:    types2.Address(r.Header.ProposerAddress),
//			},
//			Commit: &types2.Commit{
//				Height: r.Commit.Height,
//				Round:  int32(r.Commit.Round),
//				BlockID: types2.BlockID{
//					Hash: tmbytes.HexBytes(r.Commit.BlockID.Hash),
//					PartSetHeader: types2.PartSetHeader{
//						Total: uint32(r.Commit.BlockID.PartsHeader.Total),
//						Hash:  tmbytes.HexBytes(r.Commit.BlockID.Hash),
//					},
//				},
//				Signatures: ss,
//			},
//		},
//		CanonicalCommit: r.CanonicalCommit,
//	}
//	rett.Commit.BlockID.Hash = rett.Header.Hash()
//	for _,sig:=range rett.Commit.Signatures{
//		for _,v:=range rett.Commit.Signatures{
//
//		}
//		rett.Commit.VoteSignBytes(rett.ChainID,1)
//	}
//	//if hhash, chash := rett.Header.Hash(), rett.Commit.BlockID.Hash; !bytes.Equal(hhash, chash) {
//	//	fmt.Println("asd")
//	//}
//	return &rett
//}
//
//func ConvSignheaderToIBCSignHeader(h *types.SignedHeader) *types.IBCSignedHeader {
//	ret := &types.IBCSignedHeader{
//		IBCHeader: ConvHeadersToIbcHeader(h.Header),
//		Commit:    ConvCommitToIBCCommit(h.Header, h.Commit),
//	}
//
//	return ret
//}
//func ConvHeadersToIbcHeader(h *types.Header) *types.IBCHeader {
//	ret := &types.IBCHeader{
//		Version: version.Consensus{
//			Block: uint64(h.Version.Block),
//			App:   uint64(h.Version.App),
//		},
//		ChainID: h.ChainID,
//		Height:  h.Height,
//		Time:    h.Time,
//		LastBlockID: types.IBCBlockID{
//			// TODO
//			Hash: h.LastBlockID.Hash,
//			PartSetHeader: types.IBCPartSetHeader{
//				Total: uint32(h.LastBlockID.PartsHeader.Total),
//				Hash:  h.LastBlockID.PartsHeader.Hash,
//			},
//		},
//		LastCommitHash:     h.LastCommitHash,
//		DataHash:           h.DataHash,
//		ValidatorsHash:     h.ValidatorsHash,
//		NextValidatorsHash: h.NextValidatorsHash,
//		ConsensusHash:      h.ConsensusHash,
//		AppHash:            h.AppHash,
//		LastResultsHash:    h.LastResultsHash,
//		EvidenceHash:       h.EvidenceHash,
//		ProposerAddress:    h.ProposerAddress,
//	}
//
//	return ret
//}
//
//func ConvCommitToIBCCommit(hh *types.Header, h *types.Commit) *types.IBCCommit {
//	ret := &types.IBCCommit{
//		Height: h.Height,
//		Round:  int32(h.Round),
//		BlockID: types.IBCBlockID{
//			Hash: hh.IBCHash(),
//			PartSetHeader: types.IBCPartSetHeader{
//				Total: uint32(h.BlockID.PartsHeader.Total),
//				Hash:  h.BlockID.PartsHeader.Hash,
//			},
//		},
//		Signatures: h.Signatures,
//	}
//
//	return ret
//}
