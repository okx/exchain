package types

import (
	"errors"

	tmbytes "github.com/okex/exchain/libs/tendermint/libs/bytes"
	tmproto "github.com/okex/exchain/libs/tendermint/proto/types"
)

func (sh SignedHeader) ValidateBasicForIBC(chainID string) error {
	return sh.commonValidateBasic(chainID, true)
}

func (h *Header) PureIBCHash() tmbytes.HexBytes {
	if h == nil || len(h.ValidatorsHash) == 0 {
		return nil
	}
	return h.IBCHash()
}

func (commit *Commit) IBCVoteSignBytes(chainID string, valIdx int) []byte {
	return commit.GetVote(valIdx).ibcSignBytes(chainID)
}

// ToProto converts Block to protobuf
func (b *Block) ToProto() (*tmproto.Block, error) {
	if b == nil {
		return nil, errors.New("nil Block")
	}

	pb := new(tmproto.Block)

	pb.Header = *b.Header.ToProto()
	pb.LastCommit = b.LastCommit.ToProto()
	pb.Data = b.Data.ToProto()

	protoEvidence, err := b.Evidence.ToProto()
	if err != nil {
		return nil, err
	}
	pb.Evidence = tmproto.EvidenceData{
		Evidence: protoEvidence.Evidence,
	}

	return pb, nil
}

// ToProto converts Data to protobuf
func (data *Data) ToProto() tmproto.Data {
	tp := new(tmproto.Data)

	if len(data.Txs) > 0 {
		txBzs := make([][]byte, len(data.Txs))
		for i := range data.Txs {
			txBzs[i] = data.Txs[i]
		}
		tp.Txs = txBzs
	}

	return *tp
}

// ToProto converts EvidenceData to protobuf
func (data *EvidenceData) ToProto() (*tmproto.EvidenceList, error) {
	if data == nil {
		return nil, errors.New("nil evidence data")
	}

	evi := new(tmproto.EvidenceList)
	eviBzs := make([]tmproto.Evidence, len(data.Evidence))
	for i := range data.Evidence {
		protoEvi, err := EvidenceToProto(data.Evidence[i])
		if err != nil {
			return nil, err
		}
		eviBzs[i] = *protoEvi
	}
	evi.Evidence = eviBzs

	return evi, nil
}
