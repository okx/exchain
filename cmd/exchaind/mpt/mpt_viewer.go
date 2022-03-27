package mpt

import (
	"fmt"
	"log"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/mpt"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/cobra"
)

func mptViewerCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mptviewer",
		Args:  cobra.ExactArgs(1),
		Short: "iterate mpt store (acc, evm)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return checkValidKey(args[0])
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("--------- iterate %s data start ---------\n", args[0])
			switch args[0] {
			case authtypes.StoreKey:
				iterateAccMpt(ctx)
			case evmtypes.StoreKey:
				iterateEvmMpt(ctx)
			}
			log.Printf("--------- iterate %s data end ---------\n", args[0])
		},
	}
	return cmd
}

func iterateAccMpt(ctx *server.Context) {
	accMptDb := mpt.InstanceOfAccStore()
	heightBytes, err := accMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixLatestStoredHeight)
	panicError(err)
	rootHash, err := accMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixRootMptHash, heightBytes...))
	panicError(err)
	accTrie, err := accMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
	panicError(err)
	fmt.Println("accTrie root hash:", accTrie.Hash())

	itr := trie.NewIterator(accTrie.NodeIterator(nil))
	for itr.Next() {
		fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(itr.Key), ethcmn.Bytes2Hex(itr.Value))
	}
}

func iterateEvmMpt(ctx *server.Context) {
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

		contractTrie := getStorageTrie(evmMptDb, addrHash, stateRoot)
		fmt.Println(addr.String(), contractTrie.Hash())

		cItr := trie.NewIterator(contractTrie.NodeIterator(nil))
		for cItr.Next() {
			fmt.Printf("%s: %s\n", cItr.Key, cItr.Value)
		}
	}
}
