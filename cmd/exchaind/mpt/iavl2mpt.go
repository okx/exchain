package mpt

import (
	"bytes"
	"fmt"
	"log"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
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
		Use:   "iavl2mpt",
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
	cmCtx := getDeliverStateCtx(migrationApp)
	committedHeight := cmCtx.BlockHeight() - 1

	// 0.1 initialize database of acc mpt
	accMptDb := mpt.InstanceOfAccStore()
	accTrie, err := accMptDb.OpenTrie(mpt.NilHash)
	panicError(err)
	// 0.2 initialize database of evm mpt
	evmMptDb := mpt.InstanceOfEvmStore()
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
			pushData2Database(accMptDb, accTrie, committedHeight)
			log.Println(count)
		}

		// check if the account is a contract account
		if ethAcc, ok := account.(*apptypes.EthAccount); ok {
			if !bytes.Equal(ethAcc.CodeHash, mpt.EmptyCodeHashBytes) {
				contractCount++
				// update evm mpt. Key is the address of the contract; Value is the empty root hash
				panicError(evmTrie.TryUpdate(ethAcc.EthAddress().Bytes(), mpt.EmptyRootHashBytes))
				if contractCount%100 == 0 {
					pushData2Database(evmMptDb, evmTrie, committedHeight)
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
	pushData2Database(accMptDb, accTrie, committedHeight)
	pushData2Database(evmMptDb, evmTrie, committedHeight)

	fmt.Println(fmt.Sprintf("Successfully migrate %d account (include %d contract account) at version %d", count, contractCount, committedHeight))
}

// migrateEvmFromIavlToMpt migrate evm data from iavl to mpt
func migrateEvmFromIavlToMpt(ctx *server.Context) {
	// 0.1 initialize App and context
	migrationApp := newMigrationApp(ctx)
	cmCtx := getDeliverStateCtx(migrationApp)

	// 0.1 initialize database of evm mpt, and open trie based on the latest root hash
	evmMptDb := mpt.InstanceOfEvmStore()
	rootHash := migrationApp.EvmKeeper.GetMptRootHash(uint64(cmCtx.BlockHeight() - 1))
	evmTrie, err := evmMptDb.OpenTrie(rootHash)
	panicError(err)

	/* Here are prefix keys from evm module:
			KeyPrefixBlockHash                   = []byte{0x01}
			KeyPrefixBloom                       = []byte{0x02}
			KeyPrefixCode                        = []byte{0x04}
			KeyPrefixStorage                     = []byte{0x05}
			KeyPrefixChainConfig                 = []byte{0x06}
			KeyPrefixHeightHash                  = []byte{0x07}
			KeyPrefixContractDeploymentWhitelist = []byte{0x08}
			KeyPrefixContractBlockedList         = []byte{0x09}

	   So, here are data list about the migration process:
	   1. Accounts    -> evmTrie; Code、Storage -> a contractTire
	   2. ChainConfig              -> rawdb
	   3. BlockHash = HeightHash   -> rawdb
	   4. Bloom                    -> rawdb
	   5. ContractDeploymentWhitelist、ContractBlockedList -> evmTrie
	*/

	// 1. Migratess Accounts、Code、Storage
	migrateContractToMpt(migrationApp, cmCtx, evmMptDb, evmTrie)

	// 2. Migrates ChainConfig -> rawdb
	batch := evmMptDb.TrieDB().DiskDB().NewBatch()
	migrateChainConfigToDb(migrationApp, cmCtx, batch)

	// 3. Migrates BlockHash = HeightHash -> rawdb
	miragteBlockHashesToDb(migrationApp, cmCtx, batch)

	// 4. Migrates Bloom -> rawdb
	miragteBloomsToDb(migrationApp, cmCtx, batch)

	// 5. Migrate ContractDeploymentWhitelist、ContractBlockedList -> evmTrie
	migrateSpecialAddrsToMpt(migrationApp, cmCtx, evmMptDb, evmTrie)
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
			panicError(contractTrie.TryUpdate(key[:], v))
			return false
		})
		// 1.3 calculate rootHash of contract mpt
		rootHash, err := contractTrie.Commit(nil)
		panicError(err)
		// 1.4 set the rootHash of contract mpt into evm mpt
		panicError(evmTrie.TryUpdate(addr[:], rootHash.Bytes()))

		if count%1000 == 0 {
			pushData2Database(evmMptDb, evmTrie, committedHeight)
			log.Println(count)
		}
	}
	pushData2Database(evmMptDb, evmTrie, committedHeight)
	fmt.Printf("Successfully migrate %d contract stroage at version %d\n", count, committedHeight)
}

// 2. migrateChainConfigToDb Migrates chain config
func migrateChainConfigToDb(migrationApp *app.OKExChainApp, cmCtx sdk.Context, batch ethdb.Batch) {
	config, _ := migrationApp.EvmKeeper.GetChainConfig(cmCtx)
	panicError(batch.Put(evmtypes.KeyPrefixChainConfig, migrationApp.Codec().MustMarshalBinaryBare(config)))
	writeDataToRawdb(batch)
	fmt.Printf("Successfully migrate chain config\n")
}

// 3. miragteBlockHashesToDb Migrates BlockHash/HeightHash
func miragteBlockHashesToDb(migrationApp *app.OKExChainApp, cmCtx sdk.Context, batch ethdb.Batch) {
	count := 0
	migrationApp.EvmKeeper.IterateBlockHash(cmCtx, func(key []byte, value []byte) bool {
		count++
		panicError(batch.Put(key, value))
		panicError(batch.Put(append(evmtypes.KeyPrefixHeightHash, value...), key[1:]))

		if count%1000 == 0 {
			writeDataToRawdb(batch)
			log.Printf("write block hash between %d~%d\n", count-1000, count)
		}
		return false
	})
	writeDataToRawdb(batch)
	fmt.Printf("Successfully migrate %d block-hashes\n", count)
}

// 4. miragteBloomsToDb Migrates Bloom
func miragteBloomsToDb(migrationApp *app.OKExChainApp, cmCtx sdk.Context, batch ethdb.Batch) {
	count := 0
	migrationApp.EvmKeeper.IterateBlockBloom(cmCtx, func(key []byte, value []byte) bool {
		count++
		panicError(batch.Put(key, value))

		if count%1000 == 0 {
			writeDataToRawdb(batch)
			log.Printf("write bloom between %d~%d\n", count-1000, count)
		}
		return false
	})
	writeDataToRawdb(batch)
	fmt.Printf("Successfully migrate %d blooms\n", count)
}

// 5. migrateSpecialAddrsToMpt Migrates ContractDeploymentWhitelist、ContractBlockedList
func migrateSpecialAddrsToMpt(migrationApp *app.OKExChainApp, cmCtx sdk.Context, evmMptDb ethstate.Database, evmTrie ethstate.Trie) {
	csdb := evmtypes.CreateEmptyCommitStateDB(migrationApp.EvmKeeper.GenerateCSDBParams(), cmCtx)
	whiteList := csdb.GetContractDeploymentWhitelist()
	for i := 0; i < len(whiteList); i++ {
		panicError(evmTrie.TryUpdate(evmtypes.GetContractDeploymentWhitelistMemberKey(whiteList[i]), []byte(" ")))
	}

	blockedList := csdb.GetContractBlockedList()
	for i := 0; i < len(blockedList); i++ {
		panicError(evmTrie.TryUpdate(evmtypes.GetContractBlockedListMemberKey(blockedList[i]), []byte(" ")))
	}
	bcml := csdb.GetContractMethodBlockedList()
	for i := 0; i < len(bcml); i++ {
		evmtypes.SortContractMethods(bcml[i].BlockMethods)
		value := migrationApp.Codec().MustMarshalJSON(bcml[i].BlockMethods)
		sortedValue := sdk.MustSortJSON(value)
		panicError(evmTrie.TryUpdate(evmtypes.GetContractBlockedListMemberKey(bcml[i].Address), sortedValue))
	}

	committedHeight := cmCtx.BlockHeight() - 1
	pushData2Database(evmMptDb, evmTrie, committedHeight)
	fmt.Printf("Successfully migrate %d addresses in white list, %d addresses in blocked list, %d addresses in method block list at version %d\n",
		len(whiteList), len(blockedList), len(bcml), committedHeight)
}
