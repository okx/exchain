package main

import (
	"fmt"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	ethstate "github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/okex/exchain/app"
	types2 "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/mpt"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	types3 "github.com/okex/exchain/libs/types"
	"github.com/spf13/cobra"
	"log"
	"path/filepath"
)

func migrateCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-state",
		Short: "migrate iavl state to mpt state",
	}

	cmd.AddCommand(
		migrateAccountCmd(ctx),
		migrateContractCmd(ctx),
		cleanRawDBCmd(ctx),
	)

	return cmd
}

func migrateAccountCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-account",
		Short: "1. migrate iavl account to mpt account",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- migrate account start ---------")
			migrateAccount(ctx)
			log.Println("--------- migrate account end ---------")
		},
	}
	cmd.Flags().String(FlagDisplayContractAddr, "", "target contract address to display")
	cmd.Flags().Int64(FlagDisplayVersion, 0, "target state version to display")
	return cmd
}

func migrateContractCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-contract",
		Short: "2. migrate iavl contract state to mpt contract state",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- display state start ---------")
			migrateContract(ctx)
			log.Println("--------- display state end ---------")
		},
	}
	cmd.Flags().String(FlagDisplayContractAddr, "", "target contract address to display")
	cmd.Flags().Int64(FlagDisplayVersion, 0, "target state version to display")
	return cmd
}

func cleanRawDBCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean-rawdb",
		Short: "3. clean up migrated iavl state",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- display state start ---------")
			//cleanRawDB(ctx)
			log.Println("--------- display state end ---------")
		},
	}
	cmd.Flags().String(FlagDisplayContractAddr, "", "target contract address to display")
	cmd.Flags().Int64(FlagDisplayVersion, 0, "target state version to display")
	return cmd
}

//----------------------------------------------------------------
func migrateAccount(ctx *server.Context) {
	migApp := newMigrationApp(ctx)

	ver, err := migApp.GetCommitVersion()
	panicError(err)

	// init deliver state
	migApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: ver + 1}})
	cmCtx := migApp.GetDeliverStateCtx()

	mptDb := mpt.InstanceOfMptStore()
	mptTrie, err := mptDb.OpenTrie(ethcmn.Hash{})
	panicError(err)

	cnt := 0
	contractCnt := 0

	migApp.AccountKeeper.MigrateAccounts(cmCtx, func(account authexported.Account, key, value []byte) (stop bool) {
		cnt += 1
		err := mptTrie.TryUpdate(key, value)
		panicError(err)

		if cnt % 100 == 0 {
			pushData2Database(mptDb, mptTrie, cmCtx.BlockHeight() - 1, migApp.AccountKeeper.RetrievalStorageRoot)
			fmt.Println(cnt)
		}

		// contract account, migrate contract code
		switch account.(type) {
		case *types2.EthAccount:
			ethAcc := account.(*types2.EthAccount)
			if len(ethAcc.CodeHash) > 0 {
				contractCnt += 1

				// migrate code
				cHash := ethcmn.BytesToHash(ethAcc.CodeHash)
				codeWriter := mptDb.TrieDB().DiskDB().NewBatch()
				code := migApp.EvmKeeper.GetCodeByHash(cmCtx, cHash)
				rawdb.WriteCode(codeWriter, cHash, code)
				err = codeWriter.Write()
				panicError(err)
			}
		default:
			//do nothing
		}
		return false
	})
	pushData2Database(mptDb, mptTrie, cmCtx.BlockHeight() - 1, migApp.AccountKeeper.RetrievalStorageRoot)

	fmt.Println(fmt.Sprintf("Successfule migrate %d account (include %d contract account) at version %d", cnt, contractCnt, cmCtx.BlockHeight() - 1))
}

func migrateContract(ctx *server.Context) {
	migApp := newMigrationApp(ctx)

	ver, err := migApp.GetCommitVersion()
	panicError(err)

	// init deliver state
	migApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: ver + 1}})
	cmCtx := migApp.GetDeliverStateCtx()

	mptDb := mpt.InstanceOfMptStore()
	rootHash := getMptRootHash(mptDb, uint64(cmCtx.BlockHeight() - 1))
	mptTrie, err := mptDb.OpenTrie(rootHash)
	panicError(err)

	cnt := 0
	migApp.AccountKeeper.MigrateAccounts(cmCtx, func(account authexported.Account, key, value []byte) (stop bool) {
		// contract account, migrate contract code
		switch account.(type) {
		case *types2.EthAccount:
			ethAcc := account.(*types2.EthAccount)

			if len(ethAcc.CodeHash) > 0 {
				cnt += 1

				addr := ethcmn.BytesToAddress(key)
				addrHash := ethcrypto.Keccak256Hash(addr[:])
				contractTrie, err := mptDb.OpenStorageTrie(addrHash, ethcmn.Hash{})
				panicError(err)

				_ = migApp.EvmKeeper.ForEachStorage(cmCtx, addr, func(key, value ethcmn.Hash) bool {
					// Encoding []byte cannot fail, ok to ignore the error.
					v, _ := rlp.EncodeToBytes(ethcmn.TrimLeftZeroes(value[:]))
					err := contractTrie.TryUpdate(key[:], v)
					panicError(err)

					return false
				})

				ethAcc.StateRoot, err = contractTrie.Commit(nil)
				panicError(err)

				bz, err := migApp.AccountKeeper.EncodeAccount(ethAcc)
				panicError(err)

				// use the FUCK key here, not the addr[:], for application logic will change the the addr with some odd prefix...
				err = mptTrie.TryUpdate(key, bz)
				panicError(err)

				if cnt % 100 == 0 {
					pushData2Database(mptDb, mptTrie, cmCtx.BlockHeight() - 1, migApp.AccountKeeper.RetrievalStorageRoot)
					fmt.Println(cnt)
				}

			}
		default:
			//do nothing
		}
		return false
	})
	pushData2Database(mptDb, mptTrie, cmCtx.BlockHeight() - 1, migApp.AccountKeeper.RetrievalStorageRoot)

	fmt.Println(fmt.Sprintf("Successfule migrate %d contract stroage at version %d", cnt, cmCtx.BlockHeight() - 1))
}

func cleanRawDB(ctx *server.Context) {
	fmt.Println("Not implement!!!")
}

//----------------------------------------------------------------

func pushData2Database(db ethstate.Database, tr ethstate.Trie, height int64, retrieval types3.StorageRootRetrieval) {
	var storageRoot ethcmn.Hash
	root, err := tr.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent ethcmn.Hash) error {
		storageRoot = retrieval(leaf)
		if storageRoot != types.EmptyRootHash && storageRoot != (ethcmn.Hash{}) {
			db.TrieDB().Reference(storageRoot, parent)
		}
		return nil
	})
	panicError(err)

	err = db.TrieDB().Commit(root, false, nil)
	panicError(err)

	setMptRootHash(db, uint64(height), root)
}

func newMigrationApp(ctx *server.Context) *app.OKExChainApp {
	rootDir := ctx.Config.RootDir
	dataDir := filepath.Join(rootDir, "data")
	db, err := openDB(applicationDB, dataDir)
	if err != nil {
		panic("fail to open application db: " + err.Error())
	}

	return app.NewOKExChainApp(
		ctx.Logger,
		db,
		nil,
		true,
		map[int64]bool{},
		0,
	)
}

// SetMptRootHash sets the mapping from block height to root mpt hash
func setMptRootHash(db ethstate.Database, height uint64, hash ethcmn.Hash) {
	hhash := sdk.Uint64ToBigEndian(height)
	db.TrieDB().DiskDB().Put(mpt.KeyPrefixLatestStoredHeight, hhash)
	db.TrieDB().DiskDB().Put(append(mpt.KeyPrefixRootMptHash, hhash...), hash.Bytes())
}

// getMptRootHash gets root mpt hash from block height
func getMptRootHash(db ethstate.Database, height uint64) ethcmn.Hash {
	hhash := sdk.Uint64ToBigEndian(height)
	rst, err := db.TrieDB().DiskDB().Get(append(mpt.KeyPrefixRootMptHash, hhash...))
	if err != nil || len(rst) == 0 {
		return ethcmn.Hash{}
	}

	return ethcmn.BytesToHash(rst)
}
