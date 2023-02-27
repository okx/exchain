package mpt

import (
	"bytes"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/cmd/exchaind/base"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/mpt"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/iavl"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strconv"
)

func cmpIavlMptCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cmpiavl <v2> <type> <iavl version>",
		Short: "cmpiavl <v2> <type> <iavl version>",
		Long:  "compare iavl v2 mpt data",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- cmp mpt start ---------")
			if len(args) < 3 {
				log.Println("need v2")
				return
			}
			v2, t := args[0], args[1]
			version, _ := strconv.Atoi(args[2])
			cmpIavl(v2, t, int64(version))
			log.Println("--------- cmp success ---------")
			log.Println("--------- cmp mpt end ---------")
		},
	}
	return cmd
}

func cmpIavl(v2, t string, version int64) {
	if t == "acc" {
		cmpIavlAcc(v2, version)
	} else {
		cmpIavlEvm(v2, version)
	}
}

func cmpIavlEvm(v2 string, version int64) {
	accMptDb := mpt.InstanceOfMptStore()
	heightBytes, err := accMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixAccLatestStoredHeight)
	panicError(err)
	rootHash, err := accMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixAccRootMptHash, heightBytes...))
	panicError(err)
	accTrie, err := accMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
	panicError(err)
	fmt.Println("accTrie root hash:", accTrie.Hash())

	tree, db := getTree(v2, "evm", version)
	defer db.Close()

	iavlIter := iavl.NewIterator(nil, nil, true, tree.ImmutableTree)
	defer iavlIter.Close()

	evmMptDb := getDB(v2)
	hhash, err := evmMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixAccLatestStoredHeight)
	panicError(err)
	v2rootHash, err := evmMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixEvmRootMptHash, hhash...))
	panicError(err)
	evmTrie, err := evmMptDb.OpenTrie(ethcmn.BytesToHash(v2rootHash))
	panicError(err)
	fmt.Println("evmTrie root hash:", evmTrie.Hash())

	v2itr := trie.NewIterator(evmTrie.NodeIterator(nil))

	var stateRoot ethcmn.Hash

	itr := trie.NewIterator(accTrie.NodeIterator(nil))
	for {
		acc := getEthAcc(itr.Value)
		if !acc.IsContract() {
			itr.Next()
			continue
		}
		addr := ethcmn.BytesToAddress(accTrie.GetKey(itr.Key))
		addrHash := ethcrypto.Keccak256Hash(addr[:])
		stateRoot.SetBytes(itr.Value)
		contractTrie := getStorageTrie(evmMptDb, addrHash, stateRoot)

		cItr := trie.NewIterator(contractTrie.NodeIterator(nil))

		for {
			if bytes.Compare(itr.Key, iavlIter.Key()) != 0 || bytes.Compare(itr.Value, iavlIter.Value()) != 0 {
				panic("evm not equal")
			}
			iavlIter.Next()
			if cItr.Next() {
				break
			}
		}
		if !iavlIter.Valid() {
			break
		}
	}

	for itr.Next() {
		acc := getEthAcc(itr.Value)
		if acc.IsContract() {
			panic("v1 still have valid constract acc")
		}
	}
}

func cmpIavlAcc(v2 string, version int64) {
	accMptDb := mpt.InstanceOfMptStore()
	heightBytes, err := accMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixAccLatestStoredHeight)
	panicError(err)
	rootHash, err := accMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixAccRootMptHash, heightBytes...))
	panicError(err)
	accTrie, err := accMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
	panicError(err)
	fmt.Println("accTrie root hash:", accTrie.Hash())

	itr := trie.NewIterator(accTrie.NodeIterator(nil))

	tree, db := getTree(v2, "acc", version)
	defer db.Close()

	iavlIter := iavl.NewIterator(nil, nil, true, tree.ImmutableTree)
	defer iavlIter.Close()

	for {
		acc := getEthAcc(itr.Value)
		if acc.IsContract() {
			itr.Next()
			continue
		}

		if !(bytes.Compare(itr.Key, iavlIter.Key()) == 0 && bytes.Compare(itr.Value, iavlIter.Value()) == 0) {
			panic("not equal")
		}
		if !iavlIter.Valid() {
			break
		}
		iavlIter.Next()
		if !itr.Next() {
			panic("v1 should not valid")
		}
	}
	for itr.Next() {
		acc := getEthAcc(itr.Value)
		if !acc.IsContract() {
			panic("v1 still have valid no constract acc")
		}
	}
}

func getTree(dataDir, module string, version int64) (*iavl.MutableTree, dbm.DB) {
	dbBackend := viper.GetString(sdk.FlagDBBackend)
	db, err := base.OpenDB(filepath.Join(dataDir, base.AppDBName), dbm.BackendType(dbBackend))
	if err != nil {
		panic(fmt.Errorf("error opening dir %v backend %v DB: %w", dataDir, dbBackend, err))
	}
	prefix := []byte(fmt.Sprintf("s/k:%s/", module))
	prefixDB := dbm.NewPrefixDB(db, prefix)
	log.Printf("Checking.... %v\n", module)

	mutableTree, err := iavl.NewMutableTree(prefixDB, 0)
	if err != nil {
		panic(err)
	}
	if _, err := mutableTree.LoadVersion(version); err != nil {
		panic(err)
	}

	return mutableTree, db
}
