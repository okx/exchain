package main

import (
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb/util"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/node"
)

var wg sync.WaitGroup

func compactCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compact",
		Short: "Compact the leveldb",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))
			log.Println("--------- compact start ---------")
			blockStoreDB, stateDB, appDB, err := initDBs(config, node.DefaultDBProvider)
			if err != nil {
				return err
			}

			if viper.GetBool(flagBlock) {
				wg.Add(2)
				go compactDB(blockStoreDB)
				go compactDB(stateDB)
			}
			if viper.GetBool(flagApp) {
				wg.Add(1)
				go compactDB(appDB)
			}
			wg.Wait()
			log.Println("--------- compact end ---------")
			return nil
		},
	}

	cmd.Flags().BoolP(flagBlock, "b", true, "Pruning block and state DB")
	cmd.Flags().BoolP(flagApp, "a", true, "Pruning application DB")
	return cmd
}

func compactDB(db dbm.DB) {
	defer wg.Done()
	err := db.(*dbm.GoLevelDB).DB().CompactRange(util.Range{})
	panicError(err)
}
