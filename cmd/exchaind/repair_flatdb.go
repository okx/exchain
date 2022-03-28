package main

import (
	"log"
	"path/filepath"

	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/viper"

	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

const (
	latestVersionKey = "s/latest"
)

var cdc = codec.New()

func repairFlatDBCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair-flat-db",
		Short: "Repair flat kv db",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- repair flat.db start ---------")
			flatDB := newFlatKVDB()
			latestVersion := getLatestVersion(flatDB)
			log.Println("latest version before set:", latestVersion)
			setLatestVersion(flatDB, 0)
			latestVersion = getLatestVersion(flatDB)
			log.Println("latest version after set:", latestVersion)
			log.Println("--------- repair flat.db success ---------")
		},
	}
	return cmd
}

func newFlatKVDB() dbm.DB {
	rootDir := viper.GetString("home")
	dataDir := filepath.Join(rootDir, "data")
	var err error
	flatKVDB, err := sdk.NewLevelDB("flat", dataDir)
	if err != nil {
		panic(err)
	}
	return flatKVDB
}

func getLatestVersion(db dbm.DB) int64 {
	var latest int64
	latestBytes, err := db.Get([]byte(latestVersionKey))
	if err != nil {
		panic(err)
	} else if latestBytes == nil {
		return 0
	}

	err = cdc.UnmarshalBinaryLengthPrefixed(latestBytes, &latest)
	if err != nil {
		panic(err)
	}

	return latest
}

func setLatestVersion(db dbm.DB, version int64) {
	latestBytes := cdc.MustMarshalBinaryLengthPrefixed(version)
	db.SetSync([]byte(latestVersionKey), latestBytes)
}
