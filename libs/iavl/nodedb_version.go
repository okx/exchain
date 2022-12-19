package iavl

import dbm "github.com/okex/exchain/libs/tm-db"

// IsFastStorageStrategy check the db is FSS
func IsFastStorageStrategy(db dbm.DB) bool {
	ndb := &nodeDB{
		db: db,
	}
	if ndb.getLatestVersion() <= genesisVersion {
		return true
	}
	storeVersion, err := db.Get(metadataKeyFormat.Key([]byte(storageVersionKey)))
	if err != nil || storeVersion == nil {
		storeVersion = []byte(defaultStorageVersionValue)
	}
	ndb.storageVersion = string(storeVersion)

	return ndb.hasUpgradedToFastStorage() && !ndb.shouldForceFastStorageUpgrade()
}
