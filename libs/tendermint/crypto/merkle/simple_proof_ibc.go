package merkle

import cryptomerkel "github.com/okex/exchain/libs/tendermint/proto/crypto/merkle"

func (sp *SimpleProof) ToProto() *cryptomerkel.SimpleProof {
	if sp == nil {
		return nil
	}
	pb := new(cryptomerkel.SimpleProof)
	pb.Total = int64(sp.Total)
	pb.Index = int64(sp.Index)
	pb.LeafHash = sp.LeafHash
	pb.Aunts = sp.Aunts
	return pb
}
