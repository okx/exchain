package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"

	dbm "github.com/okex/exchain/libs/tm-db"
)

const (
	flagCheckDBName = "db"
	flagCheckPath1  = "path1"
	flagCheckPath2  = "path2"
)

func dbCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "db-check",
		Short: "Check all keys and values in db",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- check db start ---------")
			dbName := viper.GetString(flagCheckDBName)
			path1 := viper.GetString(flagCheckPath1)
			path2 := viper.GetString(flagCheckPath2)
			err := diffData(dbName, path1, path2)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("--------- check db end ---------")
		},
	}
	cmd.Flags().String(flagCheckDBName, "watch.db", "")
	cmd.Flags().String(flagCheckPath1, "", "")
	cmd.Flags().String(flagCheckPath2, "", "")
	return cmd
}

func diffData(dbname, path1, path2 string) error {
	db1 := dbm.NewDB(dbname, dbm.RocksDBBackend, path1)
	db2 := dbm.NewDB(dbname, dbm.RocksDBBackend, path2)
	defer func() {
		db1.Close()
		db2.Close()
	}()
	//KeySlice hold all keys of db
	keySlice := make([][]byte, 0)
	itr, err := db1.Iterator(nil, nil)
	if err != nil {
		return err
	}
	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		key := itr.Key()
		keySlice = append(keySlice, key)
	}
	log.Printf("%d keys is in %s/%s\n", len(keySlice), path1, dbname)
	for _, k := range keySlice {
		v1, err := db1.Get(k)
		if err != nil {
			return err
		}

		v2, err := db2.Get(k)
		if err != nil {
			return err
		}
		if bytes.Compare(v1, v2) != 0 {
			return fmt.Errorf("values not match k=%v v1=%v v2=%v", k, v1, v2)
		}
	}

	log.Println(dbname, "'s data matched correctly")
	return nil
}
