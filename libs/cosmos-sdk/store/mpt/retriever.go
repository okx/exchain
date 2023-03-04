package mpt

import (
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type StateRootRetriever interface {
	RetrieveStateRoot([]byte) ethcmn.Hash
}

type EmptyStateRootRetriever struct{}

func (e EmptyStateRootRetriever) RetrieveStateRoot([]byte) ethcmn.Hash {
	return ethtypes.EmptyRootHash
}
