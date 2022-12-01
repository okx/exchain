package fss

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"

	"github.com/okex/exchain/app/utils/appstatus"
	"github.com/okex/exchain/cmd/exchaind/base"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/iavl"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fssCmd.AddCommand(checkCmd)
}

// checkCmd represents the create command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check with the fast index with IAVL original nodes",
	Long: `Check fast index with IAVL original nodes:
This command is a tool to check the IAVL fast index.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storeKeys := appstatus.GetAllStoreKeys()
		return check(storeKeys)
	},
}

func check(storeKeys []string) error {
	if !appstatus.IsFastStorageStrategy() {
		return fmt.Errorf("db haven't upgraded to fast IAVL")
	}

	dataDir := viper.GetString(flagDataDir)
	dbBackend := viper.GetString(sdk.FlagDBBackend)
	db, err := base.OpenDB(filepath.Join(dataDir, base.AppDBName), dbm.BackendType(dbBackend))
	if err != nil {
		return fmt.Errorf("error opening dir %v backend %v DB: %w", dataDir, dbBackend, err)
	}
	defer db.Close()

	for _, key := range storeKeys {
		prefix := []byte(fmt.Sprintf("s/k:%s/", key))
		prefixDB := dbm.NewPrefixDB(db, prefix)
		log.Printf("Checking.... %v\n", key)

		mutableTree, err := iavl.NewMutableTree(prefixDB, 0)
		if err != nil {
			return err
		}
		if _, err := mutableTree.LoadVersion(0); err != nil {
			return err
		}

		if err := checkIndex(mutableTree); err != nil {
			return fmt.Errorf("%v iavl fast index not match %v", key, err.Error())
		}
	}
	log.Println("Success")

	return nil
}

func checkIndex(mutableTree *iavl.MutableTree) error {
	fastIterator := mutableTree.Iterator(nil, nil, true)
	defer fastIterator.Close()
	iterator := iavl.NewIterator(nil, nil, true, mutableTree.ImmutableTree)
	defer iterator.Close()

	const verboseGap = 50000
	var counter int
	for fastIterator.Valid() && iterator.Valid() {
		if bytes.Compare(fastIterator.Key(), iterator.Key()) != 0 ||
			bytes.Compare(fastIterator.Value(), iterator.Value()) != 0 {
			return fmt.Errorf("fast index key:%x value:%x, iavl node key:%x iavl node value:%x",
				fastIterator.Key(), fastIterator.Value(), iterator.Key(), iterator.Value())
		}
		if counter%verboseGap == 0 {
			log.Printf("Checked count: %v\n", counter)
		}
		counter++
		fastIterator.Next()
		iterator.Next()
	}
	log.Printf("Checked count done: %v\n", counter)

	if fastIterator.Valid() {
		return fmt.Errorf("fast index key:%v value:%v", fastIterator.Key(), fastIterator.Value())
	}

	if iterator.Valid() {
		return fmt.Errorf("iavl node key:%v value:%v", iterator.Key(), iterator.Value())
	}

	return nil
}
