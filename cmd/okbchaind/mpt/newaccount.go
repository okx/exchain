package mpt

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/trie"
	apptypes "github.com/okx/okbchain/app/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/server"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"strconv"
)

type TempNewAccountPretty struct {
	Address       sdk.AccAddress    `json:"address" yaml:"address"`
	EthAddress    string            `json:"eth_address" yaml:"eth_address"`
	Coins         sdk.Coins         `json:"coins" yaml:"coins"`
	PubKey        string            `json:"public_key" yaml:"public_key"`
	AccountNumber uint64            `json:"account_number" yaml:"account_number"`
	Sequence      uint64            `json:"sequence" yaml:"sequence"`
	CodeHash      string            `json:"code_hash" yaml:"code_hash"`
	Storage       map[string]string `json:"storages" yaml:"storages"`
}

type TempModuleAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	EthAddress    string         `json:"eth_address" yaml:"eth_address"`
	Coins         sdk.Coins      `json:"coins" yaml:"coins"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Sequence      uint64         `json:"sequence" yaml:"sequence"`
	Name          string         `json:"name" yaml:"name"`               // name of the module
	Permissions   []string       `json:"permissions" yaml:"permissions"` // permissions of module account
}

func AccountGetCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [data] [height]",
		Args:  cobra.ExactArgs(2),
		Short: "get account all storage for diff mpt account",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("--------- iterate %s data start ---------\n", args[0])
			height, err := strconv.Atoi(args[1])
			panicError(err)
			mptAccount := getAccountFromMpt(uint64(height))
			buff, err := json.Marshal(mptAccount)
			if err != nil {
				fmt.Printf("Error:%s", err)
				return
			}
			if err := ioutil.WriteFile(args[0]+"mptaccount", buff, 0555); err != nil {
				fmt.Printf("Error:%s", err)
				return
			}
			fmt.Printf("--------- iterate %s data end ---------\n", args[0])
		},
	}
	return cmd
}

func getAccountFromMpt(height uint64) map[string]interface{} {
	result := make(map[string]interface{}, 0)
	accMptDb := mpt.InstanceOfMptStore()
	heightBytes, err := accMptDb.TrieDB().DiskDB().Get(mpt.KeyPrefixAccLatestStoredHeight)
	panicError(err)
	lastestHeight := binary.BigEndian.Uint64(heightBytes)
	if lastestHeight < height {
		panic(fmt.Errorf("height(%d) > lastestHeight(%d)", height, lastestHeight))
	}
	if height == 0 {
		height = lastestHeight
	}

	hhash := sdk.Uint64ToBigEndian(height)
	rootHash, err := accMptDb.TrieDB().DiskDB().Get(append(mpt.KeyPrefixAccRootMptHash, hhash...))
	panicError(err)
	accTrie, err := accMptDb.OpenTrie(ethcmn.BytesToHash(rootHash))
	panicError(err)

	var stateRoot ethcmn.Hash
	itr := trie.NewIterator(accTrie.NodeIterator(nil))
	for itr.Next() {
		addr := ethcmn.BytesToAddress(accTrie.GetKey(itr.Key))
		addrHash := ethcrypto.Keccak256Hash(addr[:])
		acc := DecodeAccount(addr.String(), itr.Value)
		buff, err := json.Marshal(acc)
		panicError(err)
		// check if the account is a contract account
		if ethAcc, ok := acc.(*apptypes.EthAccount); ok {
			var okbAcc = TempNewAccountPretty{Storage: make(map[string]string)}
			err = json.Unmarshal(buff, &okbAcc)
			panicError(err)

			if !bytes.Equal(ethAcc.CodeHash, mpt.EmptyCodeHashBytes) {
				stateRoot.SetBytes(acc.GetStateRoot().Bytes())

				contractTrie := getStorageTrie(accMptDb, addrHash, stateRoot)

				cItr := trie.NewIterator(contractTrie.NodeIterator(nil))
				for cItr.Next() {
					key := ethcmn.BytesToHash(contractTrie.GetKey(cItr.Key))
					prefixKey := GetStorageByAddressKey(ethAcc.EthAddress().Bytes(), key.Bytes())
					okbAcc.Storage[prefixKey.String()] = hex.EncodeToString(cItr.Value)
				}
			}

			result[acc.GetAddress().String()] = &okbAcc
		} else {
			result[acc.GetAddress().String()] = acc
		}
	}
	return result
}

// GetStorageByAddressKey returns a hash of the composite key for a state
// object's storage prefixed with it's address.
func GetStorageByAddressKey(prefix, key []byte) ethcmn.Hash {
	compositeKey := make([]byte, len(prefix)+len(key))

	copy(compositeKey, prefix)
	copy(compositeKey[len(prefix):], key)
	return keccak256HashWithSyncPool(compositeKey)
}

func keccak256HashWithSyncPool(data ...[]byte) (h ethcmn.Hash) {
	d := ethcrypto.NewKeccakState()
	d.Reset()
	for _, b := range data {
		d.Write(b)
	}
	d.Read(h[:])
	return h
}
