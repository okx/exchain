package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/okexchain/app"
	ethermint "github.com/okex/okexchain/app/types"
	"github.com/okex/okexchain/x/evm"
	evmtypes "github.com/okex/okexchain/x/evm/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

const (
	flagHeight = "height"
	flagDBPath = "db_path"
)

var defaultHome, _ = os.Getwd()
var wg sync.WaitGroup

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
	cmd.PersistentFlags().StringP(flagDBPath, "", defaultHome, "directory for config and data")
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
	evmByteCodeDB, evmStateDB, err := createEVMDB(viper.GetString(flagDBPath))
	if err != nil {
		panic(err)
	}
	//defer evmDB.Close()

	ethermintApp.AccountKeeper.IterateAccounts(ctx, func(account authexported.Account) bool {
		ethAccount, ok := account.(*ethermint.EthAccount)
		if !ok {
			// ignore non EthAccounts
			return false
		}

		addr := ethAccount.EthAddress()
		if code := ethermintApp.EvmKeeper.GetCode(ctx, addr); code != nil {
			evmByteCodeDB.Set(append(evmtypes.KeyPrefixCode, ethcrypto.Keccak256Hash(code).Bytes()...), code)
		}

		wg.Add(1)
		go exportStorage(ctx, *ethermintApp.EvmKeeper, addr, evmStateDB)

		return false
	})
	wg.Wait()
	return nil
}

func createEVMDB(path string) (evmByteCodeDB, evmStateDB dbm.DB, err error) {
	evmByteCodeDB, err = sdk.NewLevelDB("evm_bytecode", path)
	if err != nil {
		return
	}
	evmStateDB, err = sdk.NewLevelDB("evm_state", path)
	return
}

func exportStorage(ctx sdk.Context, k evm.Keeper, addr ethcmn.Address, db dbm.DB) {
	defer wg.Done()
	k.IterateStorage(ctx, addr, func(hash, storage []byte) bool {
		prefix := evmtypes.AddressStoragePrefix(addr)
		db.Set(append(prefix, hash...), storage)
		return false
	})
}
