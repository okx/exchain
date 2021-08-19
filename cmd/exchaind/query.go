package main

import (
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	sm "github.com/tendermint/tendermint/state"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/store"
)


func queryCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query",
		Short: "Query blocks and states in database",
	}

	queryBlock := &cobra.Command{
		Use:   "block",
		Short: "Query blocks info in database",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))
			blockStoreDB, _, _, err := initDBs(config, node.DefaultDBProvider)
			if err != nil {
				return err
			}

			blockStore := store.NewBlockStore(blockStoreDB)
            height:=blockStore.Height()
            if blockStore.Base()==0{
                return fmt.Errorf("base of blockStore cannot be zero, may be wrong path is used.")
            }

            list, err :=blockStore.GetValidBlocks(1, height+1)
            if err != nil {
				return err
			}
			log.Printf("Block Info: %v\n", list)

			return nil
		},
	}

	queryState := &cobra.Command{
		Use:   "state",
		Short: "Query states info in database",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))
			blockStoreDB, stateDB, _, err := initDBs(config, node.DefaultDBProvider)
			if err != nil {
				return err
			}
            // get height from blockstore, it also means the right end of search interval
            blockStore := store.NewBlockStore(blockStoreDB)
            height:=blockStore.Height()
            if blockStore.Base()==0{
                return fmt.Errorf("base of blockStore cannot be zero, may be wrong path is used.")
            }

            list, err := sm.GetValidStates(stateDB, 1, height+1/*including height*/)
            if err != nil {
				return err
			}
			log.Printf("State Info: start=%v\n", list)
			return nil
		},
	}

	cmd.AddCommand(queryBlock, queryState)

	return cmd
}
