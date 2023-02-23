package mpt

import (
	"bytes"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	apptypes "github.com/okex/exchain/app/types"
	ethermint "github.com/okex/exchain/app/types"
	sdkcodec "github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/mpt"
	"github.com/okex/exchain/libs/cosmos-sdk/store/mpt/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	auth "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

var cdc = sdkcodec.New()

func init() {
	sdk.RegisterCodec(cdc)
	ethsecp256k1.RegisterCodec(cdc)
	sdkcodec.RegisterCrypto(cdc)
	auth.RegisterCodec(cdc)
	ethermint.RegisterCodec(cdc)
}

func cmpMptCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cmp <v2> <type>",
		Short: "cmp <v2> <type>",
		Long:  "compare v2 mpt data",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- cmp mpt start ---------")
			if len(args) < 2 {
				log.Println("need v2")
				return
			}
			v2, t := args[0], args[1]
			cmp(v2, t)
			log.Println("--------- cmp success ---------")
			log.Println("--------- cmp mpt end ---------")
		},
	}
	return cmd
}

func cmp(v2, t string) {
	if t == "acc" {
		cmpAcc(v2)
	} else {
		cmpEvm(v2)
	}
}

func cmpEvm(v2 string) {
	accMptDb := mpt.InstanceOfMptStore()
	heightBytes, err := accMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixAccLatestStoredHeight)
	panicError(err)
	rootHash, err := accMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixAccRootMptHash, heightBytes...))
	panicError(err)
	accTrie, err := accMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
	panicError(err)
	fmt.Println("accTrie root hash:", accTrie.Hash())

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
	var v2stateRoot ethcmn.Hash

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

		v2addr := ethcmn.BytesToAddress(evmTrie.GetKey(itr.Key))
		v2addrHash := ethcrypto.Keccak256Hash(v2addr[:])
		v2stateRoot.SetBytes(itr.Value)
		v2contractTrie := getStorageTrie(evmMptDb, v2addrHash, v2stateRoot)
		v2cItr := trie.NewIterator(v2contractTrie.NodeIterator(nil))

		cItr := trie.NewIterator(contractTrie.NodeIterator(nil))

		for v2cItr.Next() && cItr.Next() {
			if bytes.Compare(itr.Key, v2itr.Key) != 0 || bytes.Compare(itr.Value, v2itr.Value) != 0 {
				panic("evm not equal")
			}
		}
		if valid := v2cItr.Next(); valid {
			panic("v2citer valid")
		}
		if valid := cItr.Next(); valid {
			panic("cItr valid")
		}
		if !v2itr.Next() {
			break
		}
		if !itr.Next() {
			panic("v1 should not valid")
		}
	}

	for itr.Next() {
		acc := getEthAcc(itr.Value)
		if acc.IsContract() {
			panic("v1 still have valid constract acc")
		}
	}
}

func cmpAcc(v2 string) {
	accMptDb := mpt.InstanceOfMptStore()
	heightBytes, err := accMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixAccLatestStoredHeight)
	panicError(err)
	rootHash, err := accMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixAccRootMptHash, heightBytes...))
	panicError(err)
	accTrie, err := accMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
	panicError(err)
	fmt.Println("accTrie root hash:", accTrie.Hash())

	v2AccMptDb := getDB(v2)

	v2heightBytes, err := v2AccMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixAccLatestStoredHeight)
	panicError(err)
	v2rootHash, err := v2AccMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixAccRootMptHash, v2heightBytes...))
	panicError(err)
	v2accTrie, err := v2AccMptDb.OpenTrie(ethcmn.BytesToHash(v2rootHash))
	panicError(err)

	v2itr := trie.NewIterator(v2accTrie.NodeIterator(nil))

	itr := trie.NewIterator(accTrie.NodeIterator(nil))
	for {
		acc := getEthAcc(itr.Value)
		if acc.IsContract() {
			itr.Next()
			continue
		}

		if !(bytes.Compare(itr.Key, v2itr.Key) == 0 && bytes.Compare(itr.Value, v2itr.Value) == 0) {
			panic("not equal")
		}
		if !v2itr.Next() {
			break
		}
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

const (
	mptDataDir = "data"
	mptSpace   = "mpt"
)

func getDB(homeDir string) ethstate.Database {
	path := filepath.Join(homeDir, mptDataDir)

	backend := viper.GetString(sdk.FlagDBBackend)
	if backend == "" {
		backend = string(types.GoLevelDBBackend)
	}

	kvstore, e := types.CreateKvDB(mptSpace, types.BackendType(backend), path)
	if e != nil {
		panic("fail to open database: " + e.Error())
	}
	db := rawdb.NewDatabase(kvstore)

	return ethstate.NewDatabaseWithConfig(db, &trie.Config{
		Cache:     int(2048),
		Journal:   "",
		Preimages: true,
	})
}

func getEthAcc(bz []byte) apptypes.EthAccount {
	account := getAcc(bz)

	acc, ok := account.(apptypes.EthAccount)
	if !ok {
		panic("not eth account")
	}

	return acc
}

func getAcc(bz []byte) exported.Account {
	val, err := cdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(bz, (*exported.Account)(nil))
	if err == nil {
		return val.(exported.Account)
	}
	var acc exported.Account
	err = cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return acc
}
