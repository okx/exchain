package client

import (
	"fmt"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	dbm "github.com/tendermint/tm-db"
)

const (
	flagRetainHeight = "retain-height"
)

// PruningCmd dumps app state to JSON.
func PruningCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pruning",
		Short: "Export state to JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			blockStore, stateDB, err := initDBs(config, node.DefaultDBProvider)
			if err != nil {
				return err
			}
			fmt.Println(blockStore.Base(), blockStore.Size())
			retainHeight := viper.GetInt64(flagRetainHeight)
			pruned, err := pruneBlocks(blockStore, stateDB, retainHeight)
			if err != nil {
				return err
			}
			fmt.Println(blockStore.Base(), blockStore.Size())
			fmt.Printf("Pruned blocks, pruned: %d, retainHeight: %d\n", pruned, retainHeight)

			//db, err := openDB(config.RootDir)
			//if err != nil {
			//	return err
			//}
			//cms := csstore.NewCommitMultiStore(db)
			//
			//if isEmptyState(db) {
			//	if _, err := fmt.Fprintln(os.Stderr, "WARNING: State is not initialized. Returning genesis file."); err != nil {
			//		return err
			//	}
			//	return nil
			//}
			//
			//appState, validators, err := exportAppStateAndTMValidators(ctx.Logger, db, height, true, nil)
			//if err != nil {
			//	return fmt.Errorf("error exporting state: %v", err)
			//}

			fmt.Println("done")
			return nil
		},
	}

	cmd.Flags().Int64(flagRetainHeight, 1, "Export state from a particular height (-1 means latest height)")
	return cmd
}

func isEmptyState(db dbm.DB) bool {
	return db.Stats()["leveldb.sstables"] == ""
}

//
//func exportAppStateAndTMValidators(
//	logger log.Logger, db dbm.DB, height int64, forZeroHeight bool, jailWhiteList []string,
//) (json.RawMessage, []tmtypes.GenesisValidator, error) {
//	var ethermintApp *app.OKExChainApp
//	if height != -1 {
//		ethermintApp = app.NewOKExChainApp(logger, db, nil, false, map[int64]bool{}, 0)
//
//		if err := ethermintApp.LoadHeight(height); err != nil {
//			return nil, nil, err
//		} else {
//			ethermintApp = app.NewOKExChainApp(logger, db, nil, true, map[int64]bool{}, 0)
//		}
//	}
//
//	return ethermintApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
//}

func initDBs(config *cfg.Config, dbProvider node.DBProvider) (blockStore *store.BlockStore, stateDB dbm.DB, err error) {
	var blockStoreDB dbm.DB
	blockStoreDB, err = dbProvider(&node.DBContext{"blockstore", config})
	if err != nil {
		return
	}
	blockStore = store.NewBlockStore(blockStoreDB)

	stateDB, err = dbProvider(&node.DBContext{"state", config})
	if err != nil {
		return
	}

	return
}

func openDB(rootDir string) (dbm.DB, error) {
	dataDir := filepath.Join(rootDir, "data")
	db, err := sdk.NewLevelDB("application", dataDir)
	return db, err
}

// pruneStores will batch delete a list of heights from each mounted sub-store.
// Afterwards, pruneHeights is reset.
//func pruneStores(rs *rootmulti.Store, pruneHeights []int64) {
//	if len(pruneHeights) == 0 {
//		return
//	}
//
//	for key, store := range rs.stores {
//		if store.GetStoreType() == types.StoreTypeIAVL {
//			// If the store is wrapped with an inter-block cache, we must first unwrap
//			// it to get the underlying IAVL store.
//			store = rs.GetCommitKVStore(key)
//
//			if err := store.(*iavl.Store).DeleteVersions(pruneHeights...); err != nil {
//				if errCause := errors.Cause(err); errCause != nil && errCause != iavltree.ErrVersionDoesNotExist {
//					panic(err)
//				}
//			}
//		}
//	}
//}

func pruneBlocks(blockStore *store.BlockStore, stateDB dbm.DB, retainHeight int64) (uint64, error) {
	base := blockStore.Base()
	if retainHeight <= base {
		return 0, nil
	}
	pruned, err := blockStore.PruneBlocks(retainHeight)
	if err != nil {
		return 0, fmt.Errorf("failed to prune block store: %w", err)
	}
	err = sm.PruneStates(stateDB, base, retainHeight)
	if err != nil {
		return 0, fmt.Errorf("failed to prune state database: %w", err)
	}
	return pruned, nil
}
