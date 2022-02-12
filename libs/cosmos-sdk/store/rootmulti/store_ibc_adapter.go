package rootmulti

import abci "github.com/okex/exchain/libs/tendermint/abci/types"

func queryIbcProof( res *abci.ResponseQuery,info *commitInfo,storeName string){
	// Restore origin path and append proof op.
	res.Proof.Ops = append(res.Proof.Ops, info.ProofOp(storeName))
}