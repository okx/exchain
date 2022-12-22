package fss

import (
	"bytes"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/okex/exchain/x/evm"
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
		return check([]string{evm.StoreKey})
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

		if err := checkIndex(key, mutableTree); err != nil {
			return fmt.Errorf("%v iavl fast index not match %v", key, err.Error())
		}
	}
	log.Println("Success")

	return nil
}

func checkIndex(key string, mutableTree *iavl.MutableTree) error {
	fastIterator := mutableTree.Iterator(nil, nil, true)
	defer fastIterator.Close()

	const verboseGap = 50000
	var total int
	var xen int
	for fastIterator.Valid() {
		if total%verboseGap == 0 {
			log.Printf("Checked count: %v\n", total)
			log.Printf("Checked xen count: %v\n", xen)
		}
		k := fastIterator.Key()
		if len(k) == 53 && bytes.Equal(k[1:21], common.HexToAddress("1cc4d981e897a3d2e7785093a648c0a75fad0453").Bytes()) {
			xen++
		}
		if len(k) > 4 && (k[1] == 0x1c || k[2] == 0xc4 || k[3] == 0xd9) {
			log.Printf("%x\n", k)
		}

		total++
		fastIterator.Next()
	}
	log.Printf("Checked %v count done: %v\n", key, total)
	log.Printf("Checked xen %v count done: %v\n", key, xen)

	if fastIterator.Valid() {
		return fmt.Errorf("fast index key:%v value:%v", fastIterator.Key(), fastIterator.Value())
	}

	return nil
}
