package iavl

import (
	"fmt"

	ics23 "github.com/confio/ics23/go"
	storetyeps "github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	"github.com/okex/exchain/libs/iavl"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
)

func (st *Store) queryKeyForIBC(req abci.RequestQuery) (res abci.ResponseQuery) {

	if len(req.Data) == 0 {
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrTxDecode, "query cannot be zero length"))
	}

	// store the height we chose in the response, with 0 being changed to the
	// latest height
	res.Height = getHeight(st.tree, req)
	tree, err := st.tree.GetImmutable(req.Height)
	if err != nil {
		return sdkerrors.QueryResult(sdkerrors.Wrapf(iavl.ErrVersionDoesNotExist, "request height %d", req.Height))
	}
	switch req.Path {
	case "/key": // get by key
		key := req.Data // data holds the key bytes
		res.Key = key

		_, res.Value = tree.Get(key)
		if !req.Prove {
			break
		}

		mtree := &iavl.MutableTree{
			ImmutableTree: tree,
		}

		// get proof from tree and convert to merkle.Proof before adding to result
		res.Proof = getProofFromTree(mtree, req.Data, res.Value != nil)
	case "/subspace":
		var KVs []types.KVPair

		subspace := req.Data
		res.Key = subspace

		iterator := types.KVStorePrefixIterator(st, subspace)
		for ; iterator.Valid(); iterator.Next() {
			KVs = append(KVs, types.KVPair{Key: iterator.Key(), Value: iterator.Value()})
		}

		iterator.Close()
		res.Value = cdc.MustMarshalBinaryLengthPrefixed(KVs)
	default:
		return sdkerrors.QueryResult(sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unexpected query path: %v", req.Path))
	}

	return res
}

func getProofFromTree(tree *iavl.MutableTree, key []byte, exists bool) *merkle.Proof {

	var (
		commitmentProof *ics23.CommitmentProof
		err             error
	)
	//tmcrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"
	if exists {
		// value was found
		commitmentProof, err = tree.GetMembershipProof(key)
		if err != nil {
			// sanity check: If value was found, membership proof must be creatable
			panic(fmt.Sprintf("unexpected value for empty proof: %s", err.Error()))
		}
	} else {
		// value wasn't found
		commitmentProof, err = tree.GetNonMembershipProof(key)
		if err != nil {
			// sanity check: If value wasn't found, nonmembership proof must be creatable
			panic(fmt.Sprintf("unexpected error for nonexistence proof: %s", err.Error()))
		}
	}

	op := storetyeps.NewIavlCommitmentOp(key, commitmentProof)

	//&merkle.Proof{Ops: []merkle.ProofOp{iavl.NewValueOp(key, proof).ProofOp()}}
	opp := op.ProofOp()
	return &merkle.Proof{
		Ops:                  []merkle.ProofOp{opp},
		XXX_NoUnkeyedLiteral: struct{}{},
		XXX_unrecognized:     nil,
		XXX_sizecache:        0,
	}
}
