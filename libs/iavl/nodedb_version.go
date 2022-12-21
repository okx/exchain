package iavl

import (
	"fmt"
	"strconv"
	"strings"

	dbm "github.com/okex/exchain/libs/tm-db"
)

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

func GetFastStorageVersion(db dbm.DB) (int64, error) {
	storeVersion, err := db.Get(metadataKeyFormat.Key([]byte(storageVersionKey)))
	if err != nil {
		return 0, err
	}
	versions := strings.Split(string(storeVersion), fastStorageVersionDelimiter)

	if len(versions) != 2 {
		return 0, fmt.Errorf("error fast storage version format")
	}

	return strconv.ParseInt(versions[1], 10, 64)
}
