package mpt

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
)

type SnapshotAccountIterator interface {
	snapshot.AccountIterator
}

func (ms *MptStore) SnapshotAccountIterator(root, seek ethcmn.Hash) (SnapshotAccountIterator, error) {
	if ms.snaps == nil {
		return nil, fmt.Errorf("create snap shot iterator error, snap shot is not available")
	}

	return ms.snaps.AccountIterator(root, seek)
}

type SnapshotStorageIterator interface {
	snapshot.StorageIterator
}

func (ms *MptStore) SnapshotStorageIterator(root, account, seek ethcmn.Hash) (SnapshotStorageIterator, error) {
	if ms.snaps == nil {
		return nil, fmt.Errorf("create snap shot iterator error, snap shot is not available")
	}

	return ms.snaps.StorageIterator(root, account, seek)
}
