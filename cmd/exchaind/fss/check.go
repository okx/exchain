package fss

import (
	"fmt"
	"log"

	"github.com/okex/exchain/cmd/exchaind/base"
	"github.com/okex/exchain/libs/iavl"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// checkCmd represents the create command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check with the fast index with IAVL original nodes",
	Long: `Check fast index with IAVL original nodes:
This command is a tool to check the IAVL fast index.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		iavl.SetEnableFastStorage(true)
		storeKeys := getStoreKeys()
		outputModules(storeKeys)

		return checkIndex(storeKeys)
	},
}

func init() {
	fssCmd.AddCommand(checkCmd)
}

func checkIndex(storeKeys []string) error {
	dataDir := viper.GetString(flagDataDir)
	dbBackend := viper.GetString(flagDBBackend)
	db, err := base.OpenDB(dataDir+base.AppDBName, dbm.BackendType(dbBackend))
	if err != nil {
		return fmt.Errorf("error opening dir %v backend %v DB: %w", dataDir, dbBackend, err)
	}
	defer db.Close()

	for _, key := range storeKeys {
		log.Printf("Checking.... %v\n", key)
		prefix := []byte(fmt.Sprintf("s/k:%s/", key))

		prefixDB := dbm.NewPrefixDB(db, prefix)

		_, err := iavl.NewMutableTree(prefixDB, 0)
		if err != nil {
			return err
		}
	}

	return nil
}
