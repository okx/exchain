package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/node"
	sm "github.com/tendermint/tendermint/state"
	"github.com/tendermint/tendermint/store"
	dbm "github.com/tendermint/tm-db"
)

var wg sync.WaitGroup

// PruningCmd dumps app state to JSON.
func PruningCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pruning",
		Short: "Pruning state and blocks and compact the leveldb",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))
			log.Println("--------- pruning start ---------")
			blockStoreDB, stateDB, _, err := initDBs(config, node.DefaultDBProvider)
			if err != nil {
				return err
			}

			blockStore := store.NewBlockStore(blockStoreDB)
			baseHeight := blockStore.Base()
			size := blockStore.Size()
			retainHeight := baseHeight + size - 2
			log.Printf("baseHeight:%d, size:%d, retainHeight:%d\n", baseHeight, size, retainHeight)

			//wg.Add(1)
			//go pruneBlocks(blockStore, stateDB, retainHeight)
			pruneBlocks(blockStore, stateDB, retainHeight)
			//
			//app, err := getApp(ctx.Logger, appDB, retainHeight)
			//if err != nil {
			//	return fmt.Errorf("error exporting state: %v", err)
			//}
			//cms := app.GetCMS()
			//var pruneHeights []int64
			//for i := int64(1121818); i < retainHeight; i++ {
			//	//if i%100 == 0 || i >= retainHeight-100 {
			//	pruneHeights = append(pruneHeights, i)
			//	//}
			//}
			//
			//wg.Add(1)
			//go pruneAppStates(cms.(*rootmulti.Store), pruneHeights)

			//wg.Wait()
			log.Println("--------- pruning end ---------")
			return nil
		},
	}

	return cmd
}

//
//func getApp(logger tmlog.Logger, db dbm.DB, height int64) (*app.OKExChainApp, error) {
//	var ethermintApp *app.OKExChainApp
//
//	ethermintApp = app.NewOKExChainApp(logger, db, nil, false, map[int64]bool{}, 0)
//	if err := ethermintApp.LoadHeight(height); err != nil {
//		panic(err)
//	} else {
//		ethermintApp = app.NewOKExChainApp(logger, db, nil, true, map[int64]bool{}, 0)
//	}
//
//	return ethermintApp, nil
//}

func initDBs(config *cfg.Config, dbProvider node.DBProvider) (blockStoreDB, stateDB, appDB dbm.DB, err error) {
	blockStoreDB, err = dbProvider(&node.DBContext{"blockstore", config})
	if err != nil {
		return
	}

	stateDB, err = dbProvider(&node.DBContext{"state", config})
	if err != nil {
		return
	}

	appDB, err = dbProvider(&node.DBContext{"application", config})
	if err != nil {
		return
	}

	return
}

//pruneStores will batch delete a list of heights from each mounted sub-store.
//Afterwards, pruneHeights is reset.
//func pruneAppStates(rs *rootmulti.Store, pruneHeights []int64) {
//	defer wg.Done()
//	log.Println("--------- pruning app start ---------")
//	if len(pruneHeights) == 0 {
//		return
//	}
//	fmt.Println(rs.GetStores())
//	//log.Println(pruneHeights)
//	for key, store := range rs.GetStores() {
//		if store.GetStoreType() == types.StoreTypeIAVL {
//			// If the store is wrapped with an inter-block cache, we must first unwrap
//			// it to get the underlying IAVL store.
//			store = rs.GetCommitKVStore(key)
//
//			if err := store.(*iavl.Store).DeleteVersions(pruneHeights...); err != nil {
//				fmt.Println(err)
//				if errCause := errors.Cause(err); errCause != nil && errCause != iavltree.ErrVersionDoesNotExist {
//					panic(err)
//				}
//			}
//		}
//	}
//	log.Println("--------- pruning app end ---------")
//}

func pruneBlocks(blockStore *store.BlockStore, stateDB dbm.DB, retainHeight int64) {
	//defer wg.Done()

	base := blockStore.Base()
	if retainHeight <= base {
		return
	}
	pruned, err := blockStore.PruneBlocks(retainHeight)
	if err != nil {
		panic(fmt.Errorf("failed to prune block store: %w", err))
	}
	err = sm.PruneStates(stateDB, base, retainHeight)
	if err != nil {
		panic(fmt.Errorf("failed to prune state database: %w", err))
	}

	log.Printf("pruned blocks: %d, retainHeight: %d\n", pruned, retainHeight)
	log.Printf("block store base: %d, block store size: %d\n", blockStore.Base(), blockStore.Size())
}
