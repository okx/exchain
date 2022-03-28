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
	evmtypes "github.com/okex/exchain/x/evm/types"
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
			log.Printf("--------- migrate %s data to iavl data start ---------\n", args[0])
			iavl.SetIgnoreVersionCheck(true)
			switch args[0] {
			case accStoreKey:
				migrateAccFroMptToIavl(ctx)
			case evmStoreKey:
				migrateEvmFroMptToIavl(ctx)
			}
			log.Printf("--------- migrate %s data to iavl data end ---------\n", args[0])
		},
	}
	return cmd
}

func migrateAccFroMptToIavl(ctx *server.Context) {
	accMptDb := mpt.InstanceOfAccStore()
	accTrie, height := openLatestTrie(accMptDb)
	fmt.Println("accTrie root hash:", accTrie.Hash(), ", height:", height)

	appDb := openApplicationDb(ctx.Config.RootDir)
	prefixDb := dbm.NewPrefixDB(appDb, []byte(iavlAccKey))
	defer prefixDb.Close()

	tree, err := iavl.NewMutableTreeWithOpts(prefixDb, iavlstore.IavlCacheSize, &iavl.Options{InitialVersion: height - 1})
	if err != nil {
		panic("fail to create iavl tree: " + err.Error())
	}

	itr := trie.NewIterator(accTrie.NodeIterator(nil))
	for itr.Next() {
		originKey := accTrie.GetKey(itr.Key)
		fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(originKey), ethcmn.Bytes2Hex(itr.Value))
		tree.Set(originKey, itr.Value)
	}
	_, _, _, err = tree.SaveVersion(false)
	if err != nil {
		fmt.Println("fail to migrate acc data to iavl: ", err)
	}
}

func migrateEvmFroMptToIavl(ctx *server.Context) {
	evmMptDb := mpt.InstanceOfEvmStore()
	evmTrie, height := openLatestTrie(evmMptDb)
	fmt.Println("evmTrie root hash:", evmTrie.Hash(), ", height:", height)

	appDb := openApplicationDb(ctx.Config.RootDir)
	prefixDb := dbm.NewPrefixDB(appDb, []byte(iavlEvmKey))
	defer prefixDb.Close()

	tree, err := iavl.NewMutableTreeWithOpts(prefixDb, iavlstore.IavlCacheSize, &iavl.Options{InitialVersion: height})
	if err != nil {
		panic("fail to create iavl tree: " + err.Error())
	}

	// 1.migrate rawdb's data to iavl
	/*  ChainConfig              -> rawdb
	 *  BlockHash = HeightHash   -> rawdb
	 *  Bloom                    -> rawdb
	 */
	diskdb := evmMptDb.TrieDB().DiskDB()
	// 1.1 set ChainConfig back to iavl
	configValue, err := diskdb.Get(evmtypes.KeyPrefixChainConfig)
	panicError(err)
	tree.Set(evmtypes.KeyPrefixChainConfig, configValue)
	// 1.2 set BlockHash/HeightHash back to iavl
	dIter := diskdb.NewIterator(evmtypes.KeyPrefixBlockHash, nil)
	for dIter.Next() {
		fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(dIter.Key()), ethcmn.Bytes2Hex(dIter.Value()))
		tree.Set(dIter.Key(), dIter.Value())
	}
	dIter = diskdb.NewIterator(evmtypes.KeyPrefixHeightHash, nil)
	for dIter.Next() {
		fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(dIter.Key()), ethcmn.Bytes2Hex(dIter.Value()))
		tree.Set(dIter.Key(), dIter.Value())
	}
	// 1.3 set Bloom back to iavl
	dIter = diskdb.NewIterator(evmtypes.KeyPrefixBloom, nil)
	for dIter.Next() {
		fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(dIter.Key()), ethcmn.Bytes2Hex(dIter.Value()))
		tree.Set(dIter.Key(), dIter.Value())
	}

	// 2.migrate state data to iavl
	var originKey []byte
	var stateRoot ethcmn.Hash
	itr := trie.NewIterator(evmTrie.NodeIterator(nil))
	for itr.Next() {
		addr := ethcmn.BytesToAddress(evmTrie.GetKey(itr.Key))
		stateRoot.SetBytes(itr.Value)

		contractTrie := getStorageTrie(evmMptDb, ethcrypto.Keccak256Hash(addr[:]), stateRoot)
		fmt.Println(addr.String(), contractTrie.Hash())

		cItr := trie.NewIterator(contractTrie.NodeIterator(nil))
		for cItr.Next() {
			originKey = contractTrie.GetKey(cItr.Key)
			fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(originKey), ethcmn.Bytes2Hex(cItr.Value))
			tree.Set(originKey, cItr.Value)
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
