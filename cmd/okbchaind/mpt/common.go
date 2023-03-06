package mpt

import (
	"fmt"
	"path/filepath"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/okx/okbchain/app"
	"github.com/okx/okbchain/libs/cosmos-sdk/server"
	iavlstore "github.com/okx/okbchain/libs/cosmos-sdk/store/iavl"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/rootmulti"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	authtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okx/okbchain/libs/iavl"
	dbm "github.com/okx/okbchain/libs/tm-db"
	tmdb "github.com/okx/okbchain/libs/tm-db"
	evmtypes "github.com/okx/okbchain/x/evm/types"
)

const (
	applicationDB = "application"

	accStoreKey    = authtypes.StoreKey
	evmStoreKey    = evmtypes.StoreKey
	legacyStoreKey = "evmlegacy"

	iavlAccKey       = "s/k:acc/"
	iavlEvmKey       = "s/k:evm/"
	iavlEvmLegacyKey = "s/k:evmlegacy/"
	KeyParams        = "s/k:params/"
)

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

// checkValidKey checks if the key is equal to authtypes.StoreKey or evmtypes.StoreKey
func checkValidKey(key string) error {
	if key != accStoreKey && key != evmStoreKey && key != legacyStoreKey {
		return fmt.Errorf("invalid key %s", key)
	}
	return nil
}

/*
 * Common functions about cosmos-sdk
 */
// newMigrationApp generates a new app with the given key and application.db
func newMigrationApp(ctx *server.Context) *app.OKExChainApp {
	appDb := openApplicationDb(ctx.Config.RootDir)
	return app.NewOKExChainApp(
		ctx.Logger,
		appDb,
		nil,
		true,
		map[int64]bool{},
		0,
	)
}

func openApplicationDb(rootdir string) tmdb.DB {
	dataDir := filepath.Join(rootdir, "data")
	appDb, err := sdk.NewDB(applicationDB, dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}
	return appDb
}

/*
 * Common functions about mpt
 */
// getStorageTrie returns the trie of the given address and stateRoot
func getStorageTrie(db ethstate.Database, addrHash, stateRoot ethcmn.Hash) ethstate.Trie {
	tr, err := db.OpenStorageTrie(addrHash, stateRoot)
	panicError(err)
	return tr
}

// pushData2Database commit the data to the database
func pushData2Database(db ethstate.Database, trie ethstate.Trie, height int64, isEvm bool) {
	var storageRoot ethcmn.Hash
	root, err := trie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
		storageRoot.SetBytes(leaf)
		if storageRoot != ethtypes.EmptyRootHash {
			db.TrieDB().Reference(storageRoot, parent)
		}
		return nil
	})
	panicError(err)

	err = db.TrieDB().Commit(root, false, nil)
	panicError(err)

	setMptRootHash(db, uint64(height), root, isEvm)
}

// setMptRootHash sets the mapping from block height to root mpt hash
func setMptRootHash(db ethstate.Database, height uint64, hash ethcmn.Hash, isEvm bool) {
	heightBytes := sdk.Uint64ToBigEndian(height)
	db.TrieDB().DiskDB().Put(mpt.KeyPrefixAccLatestStoredHeight, heightBytes)
	db.TrieDB().DiskDB().Put(append(mpt.KeyPrefixAccRootMptHash, heightBytes...), hash.Bytes())
}

func writeDataToRawdb(batch ethdb.Batch) {
	if err := batch.Write(); err != nil {
		panic(err)
	}
	batch.Reset()
}

func getUpgradedTree(db dbm.DB, prefix []byte, usePreLatest bool) *iavl.MutableTree {
	rs := rootmulti.NewStore(db)
	latestVersion := rs.GetLatestVersion()
	if latestVersion == 0 {
		return nil
	}

	db = dbm.NewPrefixDB(db, prefix)

	tree, err := iavl.NewMutableTree(db, iavlstore.IavlCacheSize)
	if err != nil {
		panic("Fail to get tree: " + err.Error())
	}

	if usePreLatest {
		latestVersion -= 1
	}

	if latestVersion <= 0 {
		panic(fmt.Sprintf("invalid version to load: %d", latestVersion))
	}

	_, err = tree.LoadVersion(latestVersion)
	if err != nil {
		panic("fail to load target version tree: " + err.Error())
	}

	return tree
}
