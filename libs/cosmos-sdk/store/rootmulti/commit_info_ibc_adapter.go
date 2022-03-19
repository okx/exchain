package rootmulti

import (
	"fmt"
	ics23 "github.com/confio/ics23/go"
	sdkmaps "github.com/okex/exchain/libs/cosmos-sdk/store/internal/maps"
	sdkproofs "github.com/okex/exchain/libs/cosmos-sdk/store/internal/proofs"
	"github.com/okex/exchain/libs/cosmos-sdk/store/types"
	"github.com/okex/exchain/libs/tendermint/crypto/merkle"
)

func (ci commitInfo) ProofOp(storeName string) merkle.ProofOp {
	cmap := ci.toMap()
	_, proofs, _ := sdkmaps.ProofsFromMap(cmap)

	proof := proofs[storeName]
	if proof == nil {
		panic(fmt.Sprintf("ProofOp for %s but not registered store name", storeName))
	}

	// convert merkle.SimpleProof to CommitmentProof
	existProof, err := sdkproofs.ConvertExistenceProof(proof, []byte(storeName), cmap[storeName])
	if err != nil {
		panic(fmt.Errorf("could not convert simple proof to existence proof: %w", err))
	}

	commitmentProof := &ics23.CommitmentProof{
		Proof: &ics23.CommitmentProof_Exist{
			Exist: existProof,
		},
	}
	return types.NewSimpleMerkleCommitmentOp([]byte(storeName), commitmentProof).ProofOp()
}

func (ci commitInfo) toMap() map[string][]byte {
	m := make(map[string][]byte, len(ci.StoreInfos))
	for _, storeInfo := range ci.StoreInfos {
		m[storeInfo.Name] = storeInfo.GetHash()
	}

	return m
}
func (si storeInfo) GetHash() []byte {
	return si.Core.CommitID.Hash
}
