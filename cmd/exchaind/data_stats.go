//go:build rocksdb
// +build rocksdb

package main

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/iavl"
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
			loop("application", fromDir)

			return nil
		},
	}

	return cmd
}

func loop(name, fromDir string) {
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
	total := 0
	latestPrefix := ""
	keys := make(map[string]int)

	orpanStatsTo := make(map[int64]int)
	orpanStatsFrom := make(map[int64]int)

	for ; iter.Valid(); iter.Next() {
		total++
		fk, format := getFormatKey(iter.Key())
		if format {
			if fk != latestPrefix {
				keys[latestPrefix] = counter + 1
				latestPrefix = fk
				counter = 0
			} else {
				counter++
			}
			if fk == "s/k:evm/o" {
				to, from := iavl.GetToFrom(iter.Key()[len(fk)-1:])
				orpanStatsTo[to]++
				orpanStatsFrom[from]++
			}
		} else {
			log.Println("unknown ", string(iter.Key()))
		}
	}
	keys[latestPrefix] = counter + 1

	tt := 0
	for _, v := range keys {
		tt += v
	}
	log.Println("-----print orpan from-----")
	for k, v := range orpanStatsFrom {
		log.Printf("%v,%v\n", k, v)
	}
	log.Println("-----print orpan to-----")
	for k, v := range orpanStatsTo {
		log.Printf("%v,%v\n", k, v)
	}
	log.Printf("total map %v done \n", tt)
	log.Printf("total %v done \n", total)
	log.Println(keys)
	iter.Close()
}

func getFormatKey(key []byte) (string, bool) {
	slashCount := 0
	slashIndex := 0
	for _, v := range key {
		if v == '/' {
			slashCount++
		}
		slashIndex++
		if slashCount == 2 && len(key) > slashIndex {
			return string(key[:slashIndex+1]), true
		}
	}

	return "", false
}
