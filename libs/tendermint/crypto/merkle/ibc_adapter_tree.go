package merkle

import "github.com/okex/exchain/libs/tendermint/crypto/tmhash"

func HashFromByteSlices(items [][]byte) []byte {
	switch len(items) {
	case 0:
		return emptyHash()
	case 1:
		return leafHash(items[0])
	default:
		k := getSplitPoint(len(items))
		left := HashFromByteSlices(items[:k])
		right := HashFromByteSlices(items[k:])
		return innerHash(left, right)
	}
}

// returns tmhash(<empty>)
func emptyHash() []byte {
	return tmhash.Sum([]byte{})
}

// ProofsFromByteSlices computes inclusion proof for given items.
// proofs[0] is the proof for items[0].
func ProofsFromByteSlices(items [][]byte) (rootHash []byte, proofs []*SimpleProof) {
	trails, rootSPN := trailsFromByteSlices(items)
	rootHash = rootSPN.Hash
	proofs = make([]*SimpleProof, len(items))
	for i, trail := range trails {
		proofs[i] = &SimpleProof{
			Total:    len(items),
			Index:    i,
			LeafHash: trail.Hash,
			Aunts:    trail.FlattenAunts(),
		}
	}
	return
}
