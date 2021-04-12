package main

import (
	"fmt"
	"log"

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

func pruningCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pruning",
		Short: "Pruning blocks",
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

			pruneBlocks(blockStore, stateDB, retainHeight)

			log.Println("--------- pruning end ---------")
			return nil
		},
	}

	return cmd
}

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

func pruneBlocks(blockStore *store.BlockStore, stateDB dbm.DB, retainHeight int64) {
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
