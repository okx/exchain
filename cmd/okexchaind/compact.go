package main

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb/util"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/node"
)

//var wg sync.WaitGroup

// CompactCmd dumps app state to JSON.
func CompactCmd(ctx *server.Context) *cobra.Command {
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
			//go compactDB(appDB)
			wg.Wait()
			log.Println("--------- compact end ---------")
			return nil
		},
	}

	return cmd
}

func compactDB(db dbm.DB) {
	defer wg.Done()
	log.Println("--------- compact start ---------")
	err := db.(*dbm.GoLevelDB).DB().CompactRange(util.Range{})
	log.Println("--------- compact end ---------")
	panicError(err)
}
