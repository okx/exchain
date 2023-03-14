package mpt

import (
	"bytes"
	"fmt"
	"log"

	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	apptypes "github.com/okx/okbchain/app/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/server"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/mpt"
	authexported "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"
	authtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"
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
	//TODO need re-build migrateAccFromIavlToMpt cmd
	nodes := trie.NewMergedNodeSet()
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
			pushData2Database(accMptDb, accTrie, committedHeight, false, nodes)
			log.Println(count)
		}

		// check if the account is a contract account
		if ethAcc, ok := account.(*apptypes.EthAccount); ok {
			if !bytes.Equal(ethAcc.CodeHash, mpt.EmptyCodeHashBytes) {
				contractCount++
				// update evm mpt. Key is the address of the contract; Value is the empty root hash
				panicError(evmTrie.TryUpdate(ethAcc.EthAddress().Bytes(), ethtypes.EmptyRootHash.Bytes()))
				if contractCount%100 == 0 {
					pushData2Database(evmMptDb, evmTrie, committedHeight, true, nodes)
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
	pushData2Database(accMptDb, accTrie, committedHeight, false, nodes)
	pushData2Database(evmMptDb, evmTrie, committedHeight, true, nodes)

	fmt.Println(fmt.Sprintf("Successfully migrate %d account (include %d contract account) at version %d", count, contractCount, committedHeight))
}
