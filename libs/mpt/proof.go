package mpt

import (
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
)

type ProofList [][]byte

func (n *ProofList) Put(key []byte, value []byte) error {
	*n = append(*n, value)
	return nil
}

func (n *ProofList) Delete(key []byte) error {
	panic("not supported")
}

const ProofOpMptValue = "mpt:v"
const ProofOpMptAbsence = "mpt:a"

func newProofOpMptValue(key []byte, proof ProofList) merkle.ProofOp {
	bz := cdc.MustMarshalBinaryLengthPrefixed(proof)
	return merkle.ProofOp{
		Type: ProofOpMptValue,
		Key:  key,
		Data: bz,
	}
}

func newProofOpMptAbsence(key []byte, proof ProofList) merkle.ProofOp {
	bz := cdc.MustMarshalBinaryLengthPrefixed(proof)
	return merkle.ProofOp{
		Type: ProofOpMptAbsence,
		Key:  key,
		Data: bz,
	}
}
