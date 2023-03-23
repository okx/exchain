package mpt

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state/snapshot"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	mpttypes "github.com/okx/okbchain/libs/cosmos-sdk/store/mpt/types"
)

var (
	gDisableSnapshot = false
	gSnapshotRebuild = false
)

const (
	// snapshotMemoryLayerCount snapshot memory layer count
	// as we dont rollback transactions so we only keep 1 memory layer
	snapshotMemoryLayerCount = 1
)

func DisableSnapshot() {
	gDisableSnapshot = true
}

func SetSnapshotRebuild(rebuild bool) {
	gSnapshotRebuild = rebuild
}

func (ms *MptStore) openSnapshot() error {
	if ms == nil || ms.db == nil || ms.trie == nil || ms.db.TrieDB().DiskDB() == nil || gDisableSnapshot {
		return fmt.Errorf("mpt store is nil or mpt trie is nil")
	}
	// If the chain was rewound past the snapshot persistent layer (causing
	// a recovery block number to be persisted to disk), check if we're still
	// in recovery mode and in that case, don't invalidate the snapshot on a
	// head mismatch.
	var recovery bool

	version := ms.CurrentVersion()
	if layer := rawdb.ReadSnapshotRecoveryNumber(ms.db.TrieDB().DiskDB()); layer != nil && *layer > uint64(version) {
		ms.logger.Error("Enabling snapshot recovery", "chainhead", version, "diskbase", *layer)
		recovery = true
	}
	var err error
	ms.snaps, err = snapshot.NewCustom(ms.db.TrieDB().DiskDB(), ms.db.TrieDB(), 256, ms.originalRoot, false, gSnapshotRebuild, recovery, ms.retriever)
	if err != nil {
		ms.logger.Error("open snapshot error", "chainhead", version, "error", err)
		return fmt.Errorf("open snapshot error %v", err)
	}

	ms.prepareSnap(ms.originalRoot)

	return nil
}

func (ms *MptStore) updateDestructs(addr []byte) {
	if ms.snap == nil {
		return
	}
	addrHash := mpttypes.Keccak256HashWithSyncPool(addr[:])
	ms.snapRWLock.Lock()
	ms.snapDestructs[addrHash] = struct{}{}
	ms.snapRWLock.Unlock()
}

func (ms *MptStore) prepareSnap(root common.Hash) {
	if ms.snaps == nil {
		return
	}
	if ms.snap = ms.snaps.Snapshot(root); ms.snap != nil {
		ms.snapDestructs = make(map[common.Hash]struct{})
		ms.snapAccounts = make(map[common.Hash][]byte)
		ms.snapStorage = make(map[common.Hash]map[common.Hash][]byte)
	}
}

func (ms *MptStore) commitSnap(root common.Hash) {
	// If snapshotting is enabled, update the snapshot tree with this new version
	if ms.snap == nil {
		return
	}
	// Only update if there's a state transition
	if parent := ms.snap.Root(); parent != root {
		if err := ms.snaps.Update(root, parent, ms.snapDestructs, ms.snapAccounts, ms.snapStorage); err != nil {
			ms.logger.Error("Failed to update snapshot tree", "from", parent, "to", root, "err", err)
		}
		// Keep snapshotMemoryLayerCount diff layers in the memory,
		// persistent layer is snapshotMemoryLayerCount+1 th.
		// - head layer is paired with HEAD state
		// - head-1 layer is paired with HEAD-1 state
		if err := ms.snaps.Cap(root, snapshotMemoryLayerCount); err != nil {
			ms.logger.Error("Failed to cap snapshot tree", "root", root, "layers", snapshotMemoryLayerCount, "err", err)
		}
	}
	ms.snap, ms.snapDestructs, ms.snapAccounts, ms.snapStorage = nil, nil, nil, nil

	ms.prepareSnap(root)
}

func (ms *MptStore) updateSnapAccounts(addr, bz []byte) {
	if ms.snap == nil {
		return
	}

	// If state snapshotting is active, cache the data til commit
	addrHash := ethcrypto.Keccak256Hash(addr[:])
	ms.snapRWLock.Lock()
	ms.snapAccounts[addrHash] = bz
	ms.snapRWLock.Unlock()
}

func (ms *MptStore) updateSnapStorages(addr common.Address, key, bz []byte) {
	if ms.snap == nil {
		return
	}

	addrHash := mpttypes.Keccak256HashWithSyncPool(AddressStoreKey(addr.Bytes()))
	// The snapshot storage map for the object
	var storage map[common.Hash][]byte
	// If state snapshotting is active, cache the data til commit
	ms.snapRWLock.Lock()
	// Retrieve the old storage map, if available, create a new one otherwise
	if storage = ms.snapStorage[addrHash]; storage == nil {
		storage = make(map[common.Hash][]byte)
		ms.snapStorage[addrHash] = storage
	}
	storageHash := mpttypes.Keccak256HashWithSyncPool(key[:])
	storage[storageHash] = bz // v will be nil if value is 0x00
	ms.snapRWLock.Unlock()
}

func (ms *MptStore) getSnapAccount(addr []byte) ([]byte, error) {
	if ms.snap == nil {
		return nil, fmt.Errorf("snap is unavaliable")
	}

	addrHash := mpttypes.Keccak256HashWithSyncPool(addr)

	ms.snapRWLock.RLock()
	defer ms.snapRWLock.RUnlock()
	if _, ok := ms.snapDestructs[addrHash]; ok {
		return nil, nil
	}
	if v, ok := ms.snapAccounts[addrHash]; ok {
		return v, nil
	}

	data, err := ms.snap.AccountRLP(addrHash)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 { // can be both nil and []byte{}
		return nil, nil
	}

	return data, err
}

func (ms *MptStore) getSnapStorage(addr common.Address, key []byte) ([]byte, error) {
	if ms.snap == nil {
		return nil, fmt.Errorf("snap is unavaliable")
	}
	addrHash := mpttypes.Keccak256HashWithSyncPool(AddressStoreKey(addr.Bytes()))
	storageHash := mpttypes.Keccak256HashWithSyncPool(key[:])
	ms.snapRWLock.RLock()
	defer ms.snapRWLock.RUnlock()
	if storage, ok := ms.snapStorage[addrHash]; ok {
		if v, ok := storage[storageHash]; ok {
			return v, nil
		}
	}

	data, err := ms.snap.Storage(addrHash, storageHash)

	if err != nil {
		return nil, err
	}
	if len(data) == 0 { // can be both nil and []byte{}
		return nil, nil
	}

	return data, err
}

func (ms *MptStore) flattenPersistSnapshot() error {
	if ms.snap == nil {
		return fmt.Errorf("snap is unavaliable")
	}
	ms.snapRWLock.RLock()
	defer ms.snapRWLock.RUnlock()
	latestStoreVersion := ms.GetLatestStoredBlockHeight()
	root := ms.GetMptRootHash(latestStoreVersion)
	if parent := ms.snap.Root(); parent != root {
		if err := ms.snaps.Update(root, parent, ms.snapDestructs, ms.snapAccounts, ms.snapStorage); err != nil {
			ms.logger.Error("Failed to update snapshot tree", "from", parent, "to", root, "err", err)
		}
	}
	return ms.snaps.Cap(root, 0)
}
