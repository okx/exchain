package main

import (
	"fmt"

	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


type dbContext struct {
	DataDir   string
	Prefix    string
	DbBackend dbm.BackendType
}

func dbCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db",
	}
	dbCtx := &dbContext{}

	cmd.AddCommand(
		dbReadCmd(dbCtx),
		dbWriteCmd(dbCtx),
	)
	dbCtx.DbBackend = dbm.BackendType(*cmd.PersistentFlags().String(flagDBBackend, "", "Database backend: goleveldb | rocksdb"))
	dbCtx.DataDir = *cmd.PersistentFlags().String(flags.FlagHome, "", "")
	return cmd
}

func dbReadCmd(ctx *dbContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read --home app.db --key key",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return dbReadData()
		},
	}
	cmd.PersistentFlags().String("key", "", "")
	return cmd
}

// dbReadData reads key-value from leveldb
func dbReadData() error {
	dataDir := viper.GetString(flags.FlagHome)
	dbBackend := dbm.BackendType(viper.GetString(flagDBBackend))
	db, err := OpenDB(dataDir, dbBackend)
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}
	defer db.Close()

	key := viper.GetString("key")
	value, err := db.Get([]byte(key))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	fmt.Println(value)
	return nil
}

func dbWriteCmd(ctx *dbContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write --home app.db --key key --value value",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return dbWriteData(ctx)
		},
	}
	cmd.PersistentFlags().String("key", "", "")
	cmd.PersistentFlags().String("value", "", "")
	return cmd
}

// dbWriteData reads key-value from leveldb
func dbWriteData(ctx *dbContext) error {
	dataDir := viper.GetString(flags.FlagHome)
	dbBackend := dbm.BackendType(viper.GetString(flagDBBackend))
	db, err := OpenDB(dataDir, dbBackend)
	if err != nil {
		return fmt.Errorf("error opening DB: %w", err)
	}
	defer db.Close()

	key := viper.GetString("key")
	value := viper.GetString("value")
	err = db.SetSync([]byte(key), []byte(value))
	if err != nil {
		fmt.Println(err)
	}
	return nil
}