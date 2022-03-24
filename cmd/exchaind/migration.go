package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/okex/exchain/libs/cosmos-sdk/store/rootmulti"
	dbm "github.com/okex/exchain/libs/tm-db"
	"log"
	"path/filepath"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/app"
	types2 "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/mpt"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/cobra"
)

var emptyCodeHash = ethcrypto.Keccak256(nil)

func migrateCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-state",
		Short: "migrate iavl state to mpt state (if use migrate mpt data, then you should set `--use-composite-key true` when you decide to use mpt to store the coming data)",
	}

	cmd.AddCommand(
		migrateAccountCmd(ctx),
		migrateContractCmd(ctx),
		cleanIavlStoreCmd(ctx),
		iteratorMptCmd(ctx),
		migrateMpt2IavlCmd(ctx),
	)

	return cmd
}

func migrateAccountCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-account",
		Short: "1. migrate iavl account to mpt account",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- migrate account start ---------")
			migrateAccount(ctx)
			log.Println("--------- migrate account end ---------")
		},
	}
	return cmd
}

func migrateContractCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-contract",
		Short: "2. migrate iavl contract state to mpt contract state",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- migrate contract state start ---------")
			migrateContract(ctx)
			log.Println("--------- migrate contract state end ---------")
		},
	}
	return cmd
}

func cleanIavlStoreCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean-IavlStore",
		Short: "3. clean up migrated iavl store",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- clean state start ---------")
			cleanIavlStore(ctx)
			log.Println("--------- clean state end ---------")
		},
	}
	return cmd
}

func iteratorMptCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iterate-mpt",
		Args:  cobra.ExactArgs(1),
		Short: "4. iterate mpt store (acc, evm)",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- iterate mpt data start ---------")
			name := args[0]
			iteratorMpt(ctx, name)
			log.Println("--------- iterate mpt data end ---------")
		},
	}
	return cmd
}

func migrateMpt2IavlCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-mpt2iavl",
		Args:  cobra.ExactArgs(1),
		Short: "5. migrate mpt data to iavl data",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- migrate mpt data to iavl data start ---------")
			name := args[0]
			migrateMpt2Iavl(ctx, name)
			log.Println("--------- migrate mpt data to iavl data end ---------")
		},
	}
	return cmd
}


func iteratorMpt(ctx *server.Context, name string) {
	switch name {
	case authtypes.StoreKey:
		accMptDb := mpt.InstanceOfAccStore()
		hhash, err := accMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixLatestStoredHeight)
		panicError(err)
		rootHash, err := accMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixRootMptHash, hhash...))
		panicError(err)
		accTrie, err := accMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
		panicError(err)
		fmt.Println("accTrie root hash:", accTrie.Hash())

		itr := trie.NewIterator(accTrie.NodeIterator(nil))
		for itr.Next() {
			fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(itr.Key), ethcmn.Bytes2Hex(itr.Value))
		}

	case evmtypes.StoreKey:
		evmMptDb := mpt.InstanceOfEvmStore()
		hhash, err := evmMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixLatestStoredHeight)
		panicError(err)
		rootHash, err := evmMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixRootMptHash, hhash...))
		panicError(err)
		evmTrie, err := evmMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
		panicError(err)
		fmt.Println("evmTrie root hash:", evmTrie.Hash())

		var stateRoot ethcmn.Hash
		itr := trie.NewIterator(evmTrie.NodeIterator(nil))
		for itr.Next() {
			addr := ethcmn.BytesToAddress(evmTrie.GetKey(itr.Key))
			addrHash := ethcrypto.Keccak256Hash(addr[:])
			stateRoot.SetBytes(itr.Value)

			contractTrie := getTrie(evmMptDb, addrHash, stateRoot)
			fmt.Println(addr.String(), contractTrie.Hash())

			cItr := trie.NewIterator(contractTrie.NodeIterator(nil))
			for cItr.Next() {
				fmt.Printf("%s: %s\n", cItr.Key, cItr.Value)
			}
		}

	}
}

func migrateMpt2Iavl(ctx *server.Context, name string) {
	switch name {
	case authtypes.StoreKey:
		iavl.SetIgnoreVersionCheck(true)
		accMptDb := mpt.InstanceOfAccStore()
		hhash, err := accMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixLatestStoredHeight)
		panicError(err)
		rootHash, err := accMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixRootMptHash, hhash...))
		panicError(err)
		accTrie, err := accMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
		panicError(err)
		fmt.Println("accTrie root hash:", accTrie.Hash())

		rootDir := ctx.Config.RootDir
		dataDir := filepath.Join(rootDir, "data")
		db, err := openDB(applicationDB, dataDir)
		if err != nil {
			panic("fail to open application db: " + err.Error())
		}
		db = dbm.NewPrefixDB(db, []byte(KeyAcc))
		defer db.Close()

		initialVersion := binary.BigEndian.Uint64(hhash)
		tree, err := iavl.NewMutableTreeWithOpts(db, DefaultCacheSize, &iavl.Options{InitialVersion: initialVersion})
		if err != nil {
			panic("fail to create iavl tree: " + err.Error())
		}

		itr := trie.NewIterator(accTrie.NodeIterator(nil))
		for itr.Next() {
			fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(itr.Key), ethcmn.Bytes2Hex(itr.Value))
			tree.Set(itr.Key, itr.Value)
		}
		_, _, _, err = tree.SaveVersion(false)
		if err != nil {
			fmt.Println("fail to migrate acc data to iavl: ", err)
		}

	case evmtypes.StoreKey:
		iavl.SetIgnoreVersionCheck(true)
		evmMptDb := mpt.InstanceOfEvmStore()
		hhash, err := evmMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixLatestStoredHeight)
		panicError(err)
		rootHash, err := evmMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixRootMptHash, hhash...))
		panicError(err)
		evmTrie, err := evmMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
		panicError(err)
		fmt.Println("evmTrie root hash:", evmTrie.Hash())

		rootDir := ctx.Config.RootDir
		dataDir := filepath.Join(rootDir, "data")
		db, err := openDB(applicationDB, dataDir)
		if err != nil {
			panic("fail to open application db: " + err.Error())
		}
		db = dbm.NewPrefixDB(db, []byte(KeyEvm))
		defer db.Close()

		initialVersion := binary.BigEndian.Uint64(hhash)
		tree, err := iavl.NewMutableTreeWithOpts(db, DefaultCacheSize, &iavl.Options{InitialVersion: initialVersion})
		if err != nil {
			panic("fail to create iavl tree: " + err.Error())
		}

		// 1.migrate bubble data to iavl
		bubbleDB, err := openDB("evmBuble", dataDir)
		if err != nil {
			panic("fail to open application db: " + err.Error())
		}
		bItr, err := bubbleDB.Iterator(nil, nil)
		if err != nil {
			panic("fail to create iavl tree: " + err.Error())
		}
		defer bItr.Close()
		fmt.Println("start to migrate bubble data to iavl")
		for ; bItr.Valid(); bItr.Next() {
			tree.Set(bItr.Key(), bItr.Value())
		}
		fmt.Println("finish migrate bubble data to iavl")

		// 2.migrate state data to iavl
		var stateRoot ethcmn.Hash
		itr := trie.NewIterator(evmTrie.NodeIterator(nil))
		for itr.Next() {
			addr := ethcmn.BytesToAddress(evmTrie.GetKey(itr.Key))
			addrHash := ethcrypto.Keccak256Hash(addr[:])
			stateRoot.SetBytes(itr.Value)

			contractTrie := getTrie(evmMptDb, addrHash, stateRoot)
			fmt.Println(addr.String(), contractTrie.Hash())

			cItr := trie.NewIterator(contractTrie.NodeIterator(nil))
			for cItr.Next() {
				fmt.Printf("%s: %s\n", cItr.Key, cItr.Value)
				tree.Set(cItr.Key, cItr.Value)
			}
		}
		_, _, _, err = tree.SaveVersion(false)
		if err != nil {
			fmt.Println("fail to migrate evm mpt data to iavl: ", err)
		}
	}
}

//----------------------------------------------------------------
func migrateAccount(ctx *server.Context) {
	migApp := newMigrationApp(ctx)

	ver, err := migApp.GetCommitVersion()
	panicError(err)

	// init deliver state
	migApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: ver + 1}})
	cmCtx := migApp.GetDeliverStateCtx()

	accMptDb := mpt.InstanceOfAccStore()
	accTrie, err := accMptDb.OpenTrie(ethcmn.Hash{})
	panicError(err)

	evmMptDb := mpt.InstanceOfEvmStore()
	evmTrie, err := evmMptDb.OpenTrie(ethcmn.Hash{})
	panicError(err)

	cnt := 0
	contractCnt := 0
	emptyRootHashByte := types.EmptyRootHash.Bytes()

	// update GlobalNumber
	accountNumber := migApp.AccountKeeper.GetNextAccountNumber(cmCtx)
	bz := migApp.Codec().MustMarshalBinaryLengthPrefixed(accountNumber)
	err = accTrie.TryUpdate(authtypes.GlobalAccountNumberKey, bz)
	panicError(err)

	// update every account
	migApp.AccountKeeper.MigrateAccounts(cmCtx, func(account authexported.Account, key, value []byte) (stop bool) {
		cnt += 1
		err := accTrie.TryUpdate(key, value)
		panicError(err)

		if cnt%100 == 0 {
			pushData2Database(accMptDb, accTrie, cmCtx.BlockHeight()-1)
			fmt.Println(cnt)
		}

		// contract account
		switch account.(type) {
		case *types2.EthAccount:
			ethAcc := account.(*types2.EthAccount)

			if !bytes.Equal(ethAcc.CodeHash, emptyCodeHash) {
				contractCnt += 1
				err = evmTrie.TryUpdate(ethAcc.EthAddress().Bytes(), emptyRootHashByte)
				panicError(err)

				cHash := ethcmn.BytesToHash(ethAcc.CodeHash)

				// migrate code
				codeWriter := evmMptDb.TrieDB().DiskDB().NewBatch()
				code := migApp.EvmKeeper.GetCodeByHash(cmCtx, cHash)
				rawdb.WriteCode(codeWriter, cHash, code)
				err = codeWriter.Write()
				panicError(err)

				if contractCnt%100 == 0 {
					pushData2Database(evmMptDb, evmTrie, cmCtx.BlockHeight()-1)
				}
			}
		default:
			//do nothing
		}

		return false
	})
	pushData2Database(accMptDb, accTrie, cmCtx.BlockHeight()-1)
	pushData2Database(evmMptDb, evmTrie, cmCtx.BlockHeight()-1)

	fmt.Println(fmt.Sprintf("Successfule migrate %d account (include %d contract account) at version %d", cnt, contractCnt, cmCtx.BlockHeight()-1))
}

func migrateContract(ctx *server.Context) {
	migApp := newMigrationApp(ctx)

	ver, err := migApp.GetCommitVersion()
	panicError(err)

	// init deliver state
	migApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: ver + 1}})
	cmCtx := migApp.GetDeliverStateCtx()

	evmMptDb := mpt.InstanceOfEvmStore()
	rootHash := migApp.EvmKeeper.GetMptRootHash(uint64(cmCtx.BlockHeight() - 1))
	evmTrie, err := evmMptDb.OpenTrie(rootHash)
	panicError(err)

	cnt := 0
	itr := trie.NewIterator(evmTrie.NodeIterator(nil))
	for itr.Next() {
		cnt += 1

		addr := ethcmn.BytesToAddress(evmTrie.GetKey(itr.Key))
		addrHash := ethcrypto.Keccak256Hash(addr[:])
		contractTrie := getTrie(evmMptDb, addrHash, ethcmn.Hash{})

		_ = migApp.EvmKeeper.ForEachStorage(cmCtx, addr, func(key, value ethcmn.Hash) bool {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(ethcmn.TrimLeftZeroes(value[:]))
			err := contractTrie.TryUpdate(key[:], v)
			panicError(err)

			return false
		})
		rootHash, err := contractTrie.Commit(nil)
		panicError(err)
		fmt.Println(addr.String(), rootHash.String())
		err = evmTrie.TryUpdate(addr[:], rootHash.Bytes())
		panicError(err)

		if cnt%100 == 0 {
			pushData2Database(evmMptDb, evmTrie, cmCtx.BlockHeight()-1)
			fmt.Println(cnt)
		}
	}
	pushData2Database(evmMptDb, evmTrie, cmCtx.BlockHeight()-1)

	fmt.Println(fmt.Sprintf("Successfule migrate %d contract stroage at version %d", cnt, cmCtx.BlockHeight()-1))
}

func cleanIavlStore(ctx *server.Context) {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	db, err := openDB(applicationDB, dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}

	rs := rootmulti.NewStore(db)
	latestVersion := rs.GetLatestVersion()

	// 1.clean account store
	fmt.Println("Start to clean account store")
	err = CleanIAVLStore(db, []byte(KeyAcc), latestVersion, DefaultCacheSize)
	if err != nil {
		fmt.Println("fail to clean iavl store: ", err)
	}

	// 2.migrate evm store's bubble data (which is not contract code and contract state) to a tmp key-val store
	bubbleDB, err := openDB("evmBuble", dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}

	fmt.Println("Start migrate evm store's bubble data")
	err = migrateEvmBubbleData(db, bubbleDB, []byte(KeyEvm), latestVersion, DefaultCacheSize)
	if err != nil {
		fmt.Println("fail to migrate evm bubble data: ", err)
	}

	// 3.clean evm store
	fmt.Println("Start to clean evm store")
	err = CleanIAVLStore(db, []byte(KeyEvm), latestVersion, DefaultCacheSize)
	if err != nil {
		fmt.Println("fail to clean iavl store: ", err)
	}
}

//----------------------------------------------------------------

func pushData2Database(db ethstate.Database, tr ethstate.Trie, height int64) {
	var storageRoot ethcmn.Hash
	root, err := tr.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
		storageRoot.SetBytes(leaf)
		if storageRoot != types.EmptyRootHash {
			db.TrieDB().Reference(storageRoot, parent)
		}
		return nil
	})
	panicError(err)

	err = db.TrieDB().Commit(root, false, nil)
	panicError(err)

	setAccMptRootHash(db, uint64(height), root)
}

func newMigrationApp(ctx *server.Context) *app.OKExChainApp {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	db, err := openDB(applicationDB, dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}

	return app.NewOKExChainApp(
		ctx.Logger,
		db,
		nil,
		true,
		map[int64]bool{},
		0,
	)
}

func getTrie(db ethstate.Database, addrHash ,  stateRoot ethcmn.Hash) ethstate.Trie {
	tr, _ := db.OpenStorageTrie(addrHash, stateRoot)
	return tr
}

// SetMptRootHash sets the mapping from block height to root mpt hash
func setAccMptRootHash(db ethstate.Database, height uint64, hash ethcmn.Hash) {
	hhash := sdk.Uint64ToBigEndian(height)
	db.TrieDB().DiskDB().Put(mpt.KeyPrefixLatestStoredHeight, hhash)
	db.TrieDB().DiskDB().Put(append(mpt.KeyPrefixRootMptHash, hhash...), hash.Bytes())
}

func CleanIAVLStore(db dbm.DB, prefix []byte, maxVersion int64, cacheSize int) error {
	if len(prefix) != 0 {
		db = dbm.NewPrefixDB(db, prefix)
	}

	tree, err := iavl.NewMutableTree(db, cacheSize)
	if err != nil {
		return err
	}

	// delete verion [from, to)
	return tree.DeleteVersionsRange(0, maxVersion + 1, true)
}

func migrateEvmBubbleData(db, bubbleDB dbm.DB, prefix []byte, maxVersion int64, cacheSize int) error{
	if len(prefix) != 0 {
		db = dbm.NewPrefixDB(db, prefix)
	}

	tree, err := iavl.NewMutableTree(db, cacheSize)
	if err != nil {
		return err
	}

	_, err = tree.LoadVersion(maxVersion)
	if err != nil {
		return err
	}

	tree.IterateRange(nil, nil, true, func(key []byte, value []byte) bool {
		saveEvmBubbleData(bubbleDB, key, value)
		return false
	})

	return nil
}

func saveEvmBubbleData(db dbm.DB, key []byte, value []byte) {
	switch key[0] {
	case evmtypes.KeyPrefixBlockHash[0]:
	case evmtypes.KeyPrefixBloom[0]:
	case evmtypes.KeyPrefixChainConfig[0]:
	case evmtypes.KeyPrefixHeightHash[0]:
	case evmtypes.KeyPrefixContractDeploymentWhitelist[0]:
	case evmtypes.KeyPrefixContractBlockedList[0]:
		db.Set(key, value)
	default:
	}
}