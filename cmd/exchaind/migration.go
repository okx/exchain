package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"path/filepath"

	"github.com/okex/exchain/libs/cosmos-sdk/store/rootmulti"
	dbm "github.com/okex/exchain/libs/tm-db"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/mpt"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/cobra"
)

func migrateCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-state",
		Short: "migrate iavl state to mpt state (if use migrate mpt data, then you should set `--use-composite-key true` when you decide to use mpt to store the coming data)",
	}

	cmd.AddCommand(
		cleanIavlStoreCmd(ctx),
		iteratorMptCmd(ctx),
		migrateMpt2IavlCmd(ctx),
	)

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

func getTrie(db ethstate.Database, addrHash, stateRoot ethcmn.Hash) ethstate.Trie {
	tr, _ := db.OpenStorageTrie(addrHash, stateRoot)
	return tr
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
	return tree.DeleteVersionsRange(0, maxVersion+1, true)
}

func migrateEvmBubbleData(db, bubbleDB dbm.DB, prefix []byte, maxVersion int64, cacheSize int) error {
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
