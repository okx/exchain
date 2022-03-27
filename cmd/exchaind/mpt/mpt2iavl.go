package mpt

import (
	"encoding/binary"
	"fmt"
	"log"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	iavlstore "github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/mpt"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/cobra"
)

func mpt2iavlCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mpt2iavl",
		Args:  cobra.ExactArgs(1),
		Short: "migrate mpt data to iavl data",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return checkValidKey(args[0])
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- migrate mpt data to iavl data start ---------")
			iavl.SetIgnoreVersionCheck(true)
			switch args[0] {
			case accStoreKey:
				migrateAccFroMptToIavl(ctx)
			case evmStoreKey:
				migrateEvmFroMptToIavl(ctx)
			}
			log.Println("--------- migrate mpt data to iavl data end ---------")
		},
	}
	return cmd
}

func migrateAccFroMptToIavl(ctx *server.Context) {
	accMptDb := mpt.InstanceOfAccStore()
	accTrie, height := openLatestTrie(accMptDb)
	fmt.Println("accTrie root hash:", accTrie.Hash())

	appDb := openApplicationDb(ctx.Config.RootDir)
	prefixDb := dbm.NewPrefixDB(appDb, []byte(accStoreKey))
	defer prefixDb.Close()

	tree, err := iavl.NewMutableTreeWithOpts(prefixDb, iavlstore.IavlCacheSize, &iavl.Options{InitialVersion: height})
	panicError(fmt.Errorf("fail to create iavl tree: " + err.Error()))

	itr := trie.NewIterator(accTrie.NodeIterator(nil))
	for itr.Next() {
		//fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(itr.Key), ethcmn.Bytes2Hex(itr.Value))
		tree.Set(itr.Key, itr.Value)
	}
	_, _, _, err = tree.SaveVersion(false)
	if err != nil {
		fmt.Println("fail to migrate acc data to iavl: ", err)
	}
}

func migrateEvmFroMptToIavl(ctx *server.Context) {
	evmMptDb := mpt.InstanceOfEvmStore()
	evmTrie, height := openLatestTrie(evmMptDb)
	fmt.Println("evmTrie root hash:", evmTrie.Hash())

	appDb := openApplicationDb(ctx.Config.RootDir)
	prefixDb := dbm.NewPrefixDB(appDb, []byte(evmStoreKey))
	defer prefixDb.Close()

	tree, err := iavl.NewMutableTreeWithOpts(prefixDb, iavlstore.IavlCacheSize, &iavl.Options{InitialVersion: height})
	if err != nil {
		panic("fail to create iavl tree: " + err.Error())
	}

	// 1.migrate bubble data to iavl

	// 2.migrate state data to iavl
	var stateRoot ethcmn.Hash
	itr := trie.NewIterator(evmTrie.NodeIterator(nil))
	for itr.Next() {
		addr := ethcmn.BytesToAddress(evmTrie.GetKey(itr.Key))
		addrHash := ethcrypto.Keccak256Hash(addr[:])
		stateRoot.SetBytes(itr.Value)

		contractTrie := getStorageTrie(evmMptDb, addrHash, stateRoot)
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

func openLatestTrie(db ethstate.Database) (ethstate.Trie, uint64) {
	heightBytes, err := db.TrieDB().DiskDB().Get(mpt.KeyPrefixLatestStoredHeight)
	panicError(err)
	rootHash, err := db.TrieDB().DiskDB().Get(append(mpt.KeyPrefixRootMptHash, heightBytes...))
	panicError(err)
	t, err := db.OpenTrie(ethcmn.BytesToHash(rootHash))
	panicError(err)
	return t, binary.BigEndian.Uint64(heightBytes)
}
