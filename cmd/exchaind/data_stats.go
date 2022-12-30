//go:build rocksdb
// +build rocksdb

package main

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/system"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

func dbStatsCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Statistics" + system.ChainName + " data",
		RunE: func(cmd *cobra.Command, args []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(flags.FlagHome))

			// {home}/data/*
			fromDir := ctx.Config.DBDir()
			R2TiKV("application", fromDir)

			return nil
		},
	}

	return cmd
}

func R2TiKV(name, fromDir string) {
	rdb, err := dbm.NewRocksDB(name, fromDir)
	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	iter, err := rdb.Iterator(nil, nil)
	if err != nil {
		panic(err)
	}

	counter := 0
	const commitGap = 50000

	for ; iter.Valid(); iter.Next() {
		if counter%commitGap == 0 {
			log.Printf("touching %v ...\n", counter)
		}
		counter++
		log.Println(iter.Key())
	}
	log.Printf("total %v done \n", counter)
	iter.Close()
}
