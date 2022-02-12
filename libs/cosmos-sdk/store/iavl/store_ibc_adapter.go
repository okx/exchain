package iavl

import (
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/types"
	sdkerrors "github.com/okex/exchain/libs/cosmos-sdk/types/errors"
	//types2 "github.com/okex/exchain/libs/cosmos-sdk/x/ibc/core/23-commitment/types"
	"github.com/okex/exchain/libs/iavl"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
)

func (st *Store) queryKeyForIBC(req abci.RequestQuery) (res abci.ResponseQuery) {

	if len(req.Data) == 0 {
		return sdkerrors.QueryResult(sdkerrors.Wrap(sdkerrors.ErrTxDecode, "query cannot be zero length"))
	}

	tree := st.tree

	// store the height we chose in the response, with 0 being changed to the
	// latest height
	res.Height = getHeight(tree, req)

	switch req.Path {
	case "/key": // get by key
		key := req.Data // data holds the key bytes

		res.Key = key
		if !st.VersionExists(res.Height) {
			res.Log = iavl.ErrVersionDoesNotExist.Error()
			break
		}

		_, res.Value = tree.GetVersioned(key, res.Height)
		if !req.Prove {
			break
		}

		// Continue to prove existence/absence of value
		// Must convert store.Tree to iavl.MutableTree with given version to use in CreateProof
		iTree, err := tree.GetImmutable(res.Height)
		if err != nil {
			// sanity check: If value for given version was retrieved, immutable tree must also be retrievable
			panic(fmt.Sprintf("version exists in store but could not retrieve corresponding versioned tree in store, %s", err.Error()))
		}
		mtree := &iavl.MutableTree{
			ImmutableTree: iTree,
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
