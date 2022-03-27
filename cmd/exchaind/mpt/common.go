package mpt

import (
	"fmt"
	"path/filepath"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/okex/exchain/app"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/mpt"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmdb "github.com/okex/exchain/libs/tm-db"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

const (
	applicationDB = "application"

	accStoreKey = authtypes.StoreKey
	evmStoreKey = evmtypes.StoreKey

	accDbName = accStoreKey + ".db"
	evmDbName = evmStoreKey + ".db"
)

func panicError(err error) {
	if err != nil {
		panic(err)
	}
}

// checkValidKey checks if the key is equal to authtypes.StoreKey or evmtypes.StoreKey
func checkValidKey(key string) error {
	if key != accStoreKey && key != evmStoreKey {
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
	appDb, err := sdk.NewLevelDB(applicationDB, dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}
	return appDb
}

func getDeliverStateCtx(migrationApp *app.OKExChainApp) sdk.Context {
	committedHeight, err := migrationApp.GetCommitVersion()
	panicError(err)
	// init deliver state
	migrationApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: committedHeight + 1}})
	cmCtx := migrationApp.GetDeliverStateCtx()
	return cmCtx
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
func pushData2Database(db ethstate.Database, trie ethstate.Trie, height int64) {
	var storageRoot ethcmn.Hash
	root, err := trie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
		storageRoot.SetBytes(leaf)
		if storageRoot != mpt.EmptyRootHash {
			db.TrieDB().Reference(storageRoot, parent)
		}
		return nil
	})
	panicError(err)

	err = db.TrieDB().Commit(root, false, nil)
	panicError(err)

	setMptRootHash(db, uint64(height), root)
}

// setMptRootHash sets the mapping from block height to root mpt hash
func setMptRootHash(db ethstate.Database, height uint64, hash ethcmn.Hash) {
	heightBytes := sdk.Uint64ToBigEndian(height)
	db.TrieDB().DiskDB().Put(mpt.KeyPrefixLatestStoredHeight, heightBytes)
	db.TrieDB().DiskDB().Put(append(mpt.KeyPrefixRootMptHash, heightBytes...), hash.Bytes())
}

func writeDataToRawdb(batch ethdb.Batch) {
	if err := batch.Write(); err != nil {
		panic(err)
	}
	batch.Reset()
}
