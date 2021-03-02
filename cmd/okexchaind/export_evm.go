package main

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/okex/okexchain/app"
	"github.com/okex/okexchain/x/evm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const (
	flagHeight = "height"
)

var CodeKeyPrefix = []byte{0x01}
var StorageKeyPrefix = []byte{0x02}
var TxsLogKeyPrefix = []byte{0x03}

// ExportEVMCmd dumps app state to JSON.
func ExportEVMCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-evm",
		Short: "Export EVM to levelDB",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			dataDir := filepath.Join(ctx.Config.RootDir, "data")
			db, err := openDB(applicationDB, dataDir)
			if err != nil {
				return err
			}

			if isEmptyState(db) {
				return errors.New("State is not initialized.")
			}

			height := viper.GetInt64(flagHeight)
			err = exportEVM(ctx.Logger, db, height)
			if err != nil {
				return fmt.Errorf("error exporting evm: %v", err)
			}

			return nil
		},
	}

	cmd.Flags().Int64(flagHeight, -1, "Export state from a particular height (-1 means latest height)")
	return cmd
}

func isEmptyState(db dbm.DB) bool {
	return db.Stats()["leveldb.sstables"] == ""
}

func exportEVM(logger log.Logger, db dbm.DB, height int64) error {
	var ethermintApp *app.OKExChainApp
	if height != -1 {
		ethermintApp = app.NewOKExChainApp(logger, db, nil, false, map[int64]bool{}, 0)

		if err := ethermintApp.LoadHeight(height); err != nil {
			return err
		} else {
			ethermintApp = app.NewOKExChainApp(logger, db, nil, true, map[int64]bool{}, 0)
		}
	}

	// Creates context with current height and checks txs for ctx to be usable by start of next block
	ctx := ethermintApp.NewContext(true, abci.Header{Height: ethermintApp.LastBlockHeight()})

	// nolint: prealloc
	evmDB, err := createContractDB(".")
	if err != nil {
		panic(err)
	}
	defer evmDB.Close()

	//ethermintApp.AccountKeeper.IterateAccounts(ctx, func(account authexported.Account) bool {
	//	ethAccount, ok := account.(*ethermint.EthAccount)
	//	if !ok {
	//		// ignore non EthAccounts
	//		return false
	//	}
	//
	//	addr := ethAccount.EthAddress()
	//	fmt.Println(addr.String())
	//	codeKey := append(CodeKeyPrefix, addr.Bytes()...)
	//	if code := ethermintApp.EvmKeeper.GetCode(ctx, addr); code != nil {
	//		evmDB.Set(codeKey, code)
	//	}
	//
	//	//go exportStorage(ctx, *ethermintApp.EvmKeeper, addr, evmDB)
	//
	//	return false
	//})

	//txsLogs := ethermintApp.EvmKeeper.GetAllTxLogs(ctx)
	//fmt.Println(len(txsLogs))
	//for _, txsLog := range txsLogs {
	//	txLogKey := append(TxsLogKeyPrefix, txsLog.Hash.Bytes()...)
	//	logs, err := evmTypes.MarshalLogs(txsLog.Logs)
	//	if err != nil {
	//		panic(err)
	//	}
	//	evmDB.Set(txLogKey, logs)
	//}
	ethermintApp.EvmKeeper.IterateTxLogs(ctx, func(hash, logs []byte) bool {
		fmt.Println(common.BytesToHash(hash).String())
		evmDB.Set(hash, logs)

		return false
	})

	return nil
}

func createContractDB(rootDir string) (dbm.DB, error) {
	//dataDir := filepath.Join(rootDir, "data")
	dataDir := rootDir
	db, err := sdk.NewLevelDB("contract", dataDir)
	return db, err
}

func exportStorage(ctx sdk.Context, k evm.Keeper, addr ethcmn.Address, db dbm.DB) {
	storage, err := k.GetAccountStorage(ctx, addr)
	if err != nil {
		panic(err)
	}

	storageKey := append(StorageKeyPrefix, addr.Bytes()...)
	for _, state := range storage {
		db.Set(append(storageKey, state.Key[:]...), state.Value[:])
	}
}
