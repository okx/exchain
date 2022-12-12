package fss

import (
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

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy <src> <dst>",
	Short: "Copy fast index for IAVL",
	Long: `Copy fast index for IAVL:
This command is a tool to Copy the IAVL fast index.
When the copy lunched, it will show Copy Fast IAVL...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("must specify src and dst")
		}

		return copyIndex(args[0], args[1])
	},
}

func init() {
	fssCmd.AddCommand(copyCmd)
}

func copyIndex(src, dst string) error {
	srcDB, err := openDB(src)
	if err != nil {
		return fmt.Errorf("open db %v error %v\n", src, err)
	}
	defer srcDB.Close()
	dstDB, err := openDB(dst)
	if err != nil {
		return fmt.Errorf("open db %v error %v\n", dst, err)
	}
	defer dstDB.Close()

	storeKeys := appstatus.GetAllStoreKeys()
	for _, key := range storeKeys {
		log.Printf("Copying .... %v\n", key)
		err := copyTree(srcDB, dstDB, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func openDB(dataDir string) (dbm.DB, error) {
	dbBackend := viper.GetString(sdk.FlagDBBackend)
	db, err := base.OpenDB(filepath.Join(dataDir, base.AppDBName), dbm.BackendType(dbBackend))
	if err != nil {
		return nil, fmt.Errorf("error opening dir %v backend %v DB: %w", dataDir, dbBackend, err)
	}

	return db, nil
}

func getLatestTree(db dbm.DB, storeKey string) (*iavl.MutableTree, error) {
	tree, err := getTree(db, storeKey)
	if err != nil {
		return nil, err
	}

	_, err = tree.LoadVersion(0)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func getTree(db dbm.DB, storeKey string) (*iavl.MutableTree, error) {
	prefix := []byte(fmt.Sprintf("s/k:%s/", storeKey))
	prefixDB := dbm.NewPrefixDB(db, prefix)

	tree, err := iavl.NewMutableTree(prefixDB, 0)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

func copyTree(srcDB, dstDB dbm.DB, storeKey string) error {
	srcTree, err := getLatestTree(srcDB, storeKey)
	if err != nil {
		return err
	}
	if !srcTree.IsFastCacheEnabled() {
		return fmt.Errorf("src db is not fast storage")
	}

	dstTree, err := getTree(dstDB, storeKey)
	if err != nil {
		return err
	}

	batch := dstTree.NewBatch()

	iter := srcTree.Iterator(nil, nil, true)
	defer iter.Close()

	counter := 0
	const commitGap = 50000
	for ; iter.Valid(); iter.Next() {
		err = dstTree.SaveFastNodeNoCache(iter.Key(), iter.Value(), srcTree.Version(), batch)
		if err != nil {
			return err
		}
		if counter%commitGap == 0 {
			if err = batch.Write(); err != nil {
				return err
			}
			batch.Close()
			batch = dstTree.NewBatch()
			log.Printf("%v copying %v fast nodes.\n", counter, storeKey)
		}
		counter++
	}
	log.Printf("%v copied %v fast nodes done.\n", counter, storeKey)
	fastStorageVersion, err := srcTree.ImmutableTree.DebugFssVersion()
	if err != nil {
		batch.Close()
		return err
	}
	batch.Set(srcTree.GetFssVersionKey(), fastStorageVersion)
	if err = batch.Write(); err != nil {
		return err
	}
	batch.Close()

	return nil
}
