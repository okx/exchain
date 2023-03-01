package mpt

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/status-im/keycard-go/hexutils"
	"log"
	"math/big"

	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/mpt"
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
			case accStoreKey:
				iterateAccMpt(ctx)
			case evmStoreKey:
				iterateEvmMpt(ctx)
			}
			log.Printf("--------- iterate %s data end ---------\n", args[0])
		},
	}
	return cmd
}

func iterateAccMpt(ctx *server.Context) {
	accMptDb := mpt.InstanceOfMptStore()
	heightBytes, err := accMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixAccLatestStoredHeight)
	panicError(err)
	heightBytes = make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, 206)
	rootHash, err := accMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixAccRootMptHash, heightBytes...))
	panicError(err)
	accTrie, err := accMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
	panicError(err)
	rootHash2 := accTrie.Hash()
	fmt.Println("accTrie root hash:", rootHash2)

	var leafCount int
	var contractCount int
	var notOK int

	total := new(big.Int)
	itr := trie.NewIterator(accTrie.NodeIterator(nil))
	for itr.Next() {
		leafCount++
		acc := DecodeAccount(ethcmn.Bytes2Hex(itr.Key), itr.Value)
		if acc != nil {
			ethAcc, ok := acc.(*types.EthAccount)
			if !ok {
				notOK++
				continue
			}

			if hex.EncodeToString(ethAcc.CodeHash) != "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470" && len(acc.GetCoins()) != 0 {
				contractCount++
			}
			for _, coin := range acc.GetCoins() {
				total.Add(total, coin.Amount.Int)
			}
			//fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(itr.Key), acc.String())
		}
	}
	height := hex.EncodeToString(heightBytes)
	fmt.Println("accTrie root hash:", ethcmn.BytesToHash(rootHash), rootHash2, "leaf count:", leafCount, "height:", height, "total:", total.String(), "invalid contract:", contractCount, "moduleAccount", notOK)
}

func iterateEvmMpt(ctx *server.Context) {
	accMptDb := mpt.InstanceOfMptStore()
	heightBytes, err := accMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixAccLatestStoredHeight)
	panicError(err)
	rootHash, err := accMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixAccRootMptHash, heightBytes...))
	panicError(err)
	accTrie, err := accMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
	panicError(err)
	fmt.Println("accTrie root hash:", accTrie.Hash())

	var stateRoot ethcmn.Hash
	itr := trie.NewIterator(accTrie.NodeIterator(nil))
	for itr.Next() {
		addr := ethcmn.BytesToAddress(accTrie.GetKey(itr.Key))
		addrHash := ethcrypto.Keccak256Hash(addr[:])
		acc := DecodeAccount(addr.String(), itr.Value)
		if acc == nil {
			continue
		}
		stateRoot.SetBytes(acc.GetStateRoot().Bytes())

		contractTrie := getStorageTrie(accMptDb, addrHash, stateRoot)
		fmt.Println(addr.String(), contractTrie.Hash())

		cItr := trie.NewIterator(contractTrie.NodeIterator(nil))
		for cItr.Next() {
			fmt.Printf("%s: %s\n", ethcmn.Bytes2Hex(cItr.Key), ethcmn.Bytes2Hex(cItr.Value))
		}
	}
}

func DecodeAccount(key string, bz []byte) exported.Account {
	val, err := auth.ModuleCdc.UnmarshalBinaryBareWithRegisteredUnmarshaller(bz, (*exported.Account)(nil))
	if err == nil {
		return val.(exported.Account)
	}
	var acc exported.Account
	err = auth.ModuleCdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		fmt.Printf(" key(%s) value(%s) err(%s)\n", key, hexutils.BytesToHex(bz), err)
		panic(err)
	}
	return acc
}
