package main

import (
	"fmt"
	"log"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
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

	queryBlockState := &cobra.Command{
		Use:   "block",
		Short: "Query blocks and states in database",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))
			blockStoreDB, _, _, err := initDBs(config, node.DefaultDBProvider)
			if err != nil {
				return err
			}

			blockStore := store.NewBlockStore(blockStoreDB)
			height := blockStore.Height()
			if blockStore.Base() == 0 {
				return fmt.Errorf("base of blockStore cannot be zero, may be wrong path is used.")
			}

			list, err := blockStore.GetValidBlocks(1, height+1)
			if err != nil {
				return err
			}
			log.Printf("Block Info: %v\n", list)

			return nil
		},
	}

	queryAppState := &cobra.Command{
		Use:   "state",
		Short: "Query application states info in database",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))
			_, _, appStateDB, err := initDBs(config, node.DefaultDBProvider)
			if err != nil {
				return err
			}

			rs := initAppStore(appStateDB)
			versions := rs.GetVersions()
			log.Printf("appState Info: %v\n", versions)
			return nil
		},
	}

	cmd.AddCommand(queryBlockState, queryAppState)

	return cmd
}
