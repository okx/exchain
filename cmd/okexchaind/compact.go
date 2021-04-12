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
			blockStoreDB, stateDB, _, err := initDBs(config, node.DefaultDBProvider)
			if err != nil {
				return err
			}

			wg.Add(2)
			go compactDB(blockStoreDB)
			go compactDB(stateDB)
			wg.Wait()
			log.Println("--------- compact end ---------")
			return nil
		},
	}

	return cmd
}

func compactDB(db dbm.DB) {
	defer wg.Done()
	err := db.(*dbm.GoLevelDB).DB().CompactRange(util.Range{})
	panicError(err)
}
