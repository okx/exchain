package mpt

import (
	"bytes"
	"fmt"
	"github.com/okex/exchain/libs/iavl"
	"log"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/okex/exchain/app"
	apptypes "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	authtypes "github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/mpt"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/cobra"
)

func iavl2mptCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Args:  cobra.ExactArgs(1),
		Use:   "iavl2mpt acc/evm",
		Short: "migrate data from iavl to mpt",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return checkValidKey(args[0])
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("--------- migrate %s start ---------\n", args[0])
			switch args[0] {
			case accStoreKey:
				migrateAccFromIavlToMpt(ctx)
			case evmStoreKey:
				migrateEvmFromIavlToMpt(ctx)
			}
			log.Printf("--------- migrate %s end ---------\n", args[0])
		},
	}
	return cmd
}

// migrateAccFromIavlToMpt migrate acc data from iavl to mpt
func migrateAccFromIavlToMpt(ctx *server.Context) {
	// 0.1 initialize App and context
	migrationApp := newMigrationApp(ctx)
	cmCtx := migrationApp.MockContext()
	committedHeight := cmCtx.BlockHeight() - 1

	// 0.1 initialize database of acc mpt
	accMptDb := mpt.InstanceOfMptStore()
	accTrie, err := accMptDb.OpenTrie(mpt.NilHash)
	panicError(err)

	// 0.2 initialize database of evm mpt
	evmMptDb := mpt.InstanceOfMptStore()
	evmTrie, err := evmMptDb.OpenTrie(mpt.NilHash)
	panicError(err)

	// 1.1 update GlobalNumber to mpt
	accountNumber := migrationApp.AccountKeeper.GetNextAccountNumber(cmCtx)
	err = accTrie.TryUpdate(authtypes.GlobalAccountNumberKey, migrationApp.Codec().MustMarshalBinaryLengthPrefixed(accountNumber))
	panicError(err)
	fmt.Println("GlobalNumber", accountNumber)

	// 1.2 update every account to mpt
	count, contractCount := 0, 0
	batch := evmMptDb.TrieDB().DiskDB().NewBatch()
	migrationApp.AccountKeeper.MigrateAccounts(cmCtx, func(account authexported.Account, key, value []byte) (stop bool) {
		count++
		if len(value) == 0 {
			log.Printf("[warning] %s has nil value\n", account.GetAddress().String())
		}

		// update acc mpt for every account
		panicError(accTrie.TryUpdate(key, value))
		if count%100 == 0 {
			pushData2Database(accMptDb, accTrie, committedHeight, false)
			log.Println(count)
		}

		// check if the account is a contract account
		if ethAcc, ok := account.(*apptypes.EthAccount); ok {
			if !bytes.Equal(ethAcc.CodeHash, mpt.EmptyCodeHashBytes) {
				contractCount++
				// update evm mpt. Key is the address of the contract; Value is the empty root hash
				panicError(evmTrie.TryUpdate(ethAcc.EthAddress().Bytes(), mpt.EmptyRootHashBytes))
				if contractCount%100 == 0 {
					pushData2Database(evmMptDb, evmTrie, committedHeight, true)
				}

				// write code to evm.db in direct
				codeHash := ethcmn.BytesToHash(ethAcc.CodeHash)
				rawdb.WriteCode(batch, codeHash, migrationApp.EvmKeeper.GetCodeByHash(cmCtx, codeHash))
				writeDataToRawdb(batch)
			}
		}

		return false
	})

	// 1.3 make sure the last data is committed to the database
	pushData2Database(accMptDb, accTrie, committedHeight, false)
	pushData2Database(evmMptDb, evmTrie, committedHeight, true)

	fmt.Println(fmt.Sprintf("Successfully migrate %d account (include %d contract account) at version %d", count, contractCount, committedHeight))
}

// migrateEvmFromIavlToMpt migrate evm data from iavl to mpt
func migrateEvmFromIavlToMpt(ctx *server.Context) {
	// 0.1 initialize App and context
	migrationApp := newMigrationApp(ctx)
	cmCtx := migrationApp.MockContext()

	// 0.1 initialize database of evm mpt, and open trie based on the latest root hash
	evmMptDb := mpt.InstanceOfMptStore()
	rootHash := migrationApp.EvmKeeper.GetMptRootHash(uint64(cmCtx.BlockHeight() - 1))
	evmTrie, err := evmMptDb.OpenTrie(rootHash)
	panicError(err)

	/* Here are prefix keys from evm module:
			KeyPrefixBlockHash
			KeyPrefixBloom
			KeyPrefixCode
			KeyPrefixStorage
			KeyPrefixChainConfig
			KeyPrefixHeightHash
			KeyPrefixContractDeploymentWhitelist
			KeyPrefixContractBlockedList

	   So, here are data list about the migration process:
	   1. Accounts    -> evmTrie
	      Code        -> rawdb   (note: done in iavl2mpt acc)
	      Storage     -> a contractTire

	   2. ChainConfig              -> iavl
	   3. BlockHash = HeightHash   -> iavl
	   4. Bloom                    -> iavl
	   5. ContractDeploymentWhitelist、ContractBlockedList -> iavl
	*/

	// 1. Migratess Accounts、Storage
	migrateContractToMpt(migrationApp, cmCtx, evmMptDb, evmTrie)

	evm2Tree := getIavlTree(migrationApp.GetDB())
	// 2. Migrates ChainConfig -> rawdb
	migrateChainConfigToIavl(evm2Tree, migrationApp, cmCtx)

	// 3. Migrates BlockHash = HeightHash -> rawdb
	miragteBlockHashesToIavl(evm2Tree, migrationApp, cmCtx)

	// 4. Migrates Bloom -> rawdb
	miragteBloomsToIavl(evm2Tree, migrationApp, cmCtx)

	// 5. Migrate ContractDeploymentWhitelist、ContractBlockedList -> rawdb
	migrateSpecialAddrsToIavl(evm2Tree, migrationApp, cmCtx)

	evm2Tree.SaveVersion(false)
}

// 1. migrateContractToMpt Migrates Accounts、Code、Storage
func migrateContractToMpt(migrationApp *app.OKExChainApp, cmCtx sdk.Context, evmMptDb ethstate.Database, evmTrie ethstate.Trie) {
	committedHeight := cmCtx.BlockHeight() - 1
	count := 0
	itr := trie.NewIterator(evmTrie.NodeIterator(nil))
	for itr.Next() {
		count++

		addr := ethcmn.BytesToAddress(evmTrie.GetKey(itr.Key))
		// 1.1 get solo contract mpt
		contractTrie := getStorageTrie(evmMptDb, ethcrypto.Keccak256Hash(addr[:]), mpt.NilHash)

		_ = migrationApp.EvmKeeper.ForEachStorage(cmCtx, addr, func(key, value ethcmn.Hash) bool {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ := rlp.EncodeToBytes(ethcmn.TrimLeftZeroes(value[:]))
			if len(v) == 0 {
				log.Printf("[warning] %s in %s has nil value\n", addr.String(), key.String())
			}
			// 1.2 set every storage into solo
			panicError(contractTrie.TryUpdate(key.Bytes(), v))
			return false
		})
		// 1.3 calculate rootHash of contract mpt
		rootHash, err := contractTrie.Commit(nil)
		panicError(err)
		// 1.4 set the rootHash of contract mpt into evm mpt
		panicError(evmTrie.TryUpdate(addr[:], rootHash.Bytes()))

		if count%100 == 0 {
			pushData2Database(evmMptDb, evmTrie, committedHeight, true)
			log.Println(count)
		}
	}
	pushData2Database(evmMptDb, evmTrie, committedHeight, true)
	fmt.Printf("Successfully migrate %d contract stroage at version %d\n", count, committedHeight)
}

// 2. migrateChainConfigToIavl Migrates chain config
func migrateChainConfigToIavl(tree *iavl.MutableTree, migrationApp *app.OKExChainApp, cmCtx sdk.Context) {
	config, _ := migrationApp.EvmKeeper.GetChainConfig(cmCtx)
	tree.Set(evmtypes.KeyPrefixChainConfig, migrationApp.Codec().MustMarshalBinaryBare(config))
	fmt.Printf("Successfully migrate chain config\n")
}

// 3. miragteBlockHashesToIavl Migrates BlockHash/HeightHash
func miragteBlockHashesToIavl(tree *iavl.MutableTree, migrationApp *app.OKExChainApp, cmCtx sdk.Context) {
	count := 0
	migrationApp.EvmKeeper.IterateBlockHash(cmCtx, func(key []byte, value []byte) bool {
		count++

		tree.Set(key, value)
		tree.Set(append(evmtypes.KeyPrefixHeightHash, value...), key[1:])

		return false
	})
	fmt.Printf("Successfully migrate %d block-hashes\n", count)
}

// 4. miragteBloomsToIavl Migrates Bloom
func miragteBloomsToIavl(tree *iavl.MutableTree, migrationApp *app.OKExChainApp, cmCtx sdk.Context) {
	count := 0
	migrationApp.EvmKeeper.IterateBlockBloom(cmCtx, func(key []byte, value []byte) bool {
		count++

		tree.Set(key, value)

		return false
	})
	fmt.Printf("Successfully migrate %d blooms\n", count)
}

// 5. migrateSpecialAddrsToIavl Migrates ContractDeploymentWhitelist、ContractBlockedList
func migrateSpecialAddrsToIavl(tree *iavl.MutableTree,migrationApp *app.OKExChainApp, cmCtx sdk.Context) {
	csdb := evmtypes.CreateEmptyCommitStateDB(migrationApp.EvmKeeper.GenerateCSDBParams(), cmCtx)

	// 5.1、deploy white list
	whiteList := csdb.GetContractDeploymentWhitelist()
	for i := 0; i < len(whiteList); i++ {
		tree.Set(append(evmtypes.KeyPrefixContractDeploymentWhitelist, whiteList[i]...), []byte(""))
	}

	// 5.2、deploy blocked list
	blockedList := csdb.GetContractBlockedList()
	for i := 0; i < len(blockedList); i++ {
		tree.Set(append(evmtypes.KeyPrefixContractBlockedList, blockedList[i]...), []byte(""))
	}

	// 5.3、deploy blocked method list
	bcml := csdb.GetContractMethodBlockedList()
	for i := 0; i < len(bcml); i++ {
		if !bcml[i].IsAllMethodBlocked() {
			evmtypes.SortContractMethods(bcml[i].BlockMethods)
			value := migrationApp.Codec().MustMarshalJSON(bcml[i].BlockMethods)
			sortedValue := sdk.MustSortJSON(value)
			tree.Set(append(evmtypes.KeyPrefixContractBlockedList, bcml[i].Address...), sortedValue)
		}
	}

	fmt.Printf("Successfully migrate %d addresses in white list, %d addresses in blocked list, %d addresses in method block list\n",
		len(whiteList), len(blockedList), len(bcml))
}
