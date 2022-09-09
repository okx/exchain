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

func init() {
	fssCmd.AddCommand(createCmd)
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create fast index for IAVL",
	Long: `Create fast index for IAVL:
This command is a tool to generate the IAVL fast index.
It will take long based on the original database size.
When the create lunched, it will show Upgrade to Fast IAVL...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		storeKeys := getStoreKeys()
		return createIndex(storeKeys)
	},
}

func createIndex(storeKeys []string) error {
	dataDir := viper.GetString(flagDataDir)
	dbBackend := viper.GetString(flagDBBackend)
	db, err := base.OpenDB(dataDir+base.AppDBName, dbm.BackendType(dbBackend))
	if err != nil {
		return fmt.Errorf("error opening dir %v backend %v DB: %w", dataDir, dbBackend, err)
	}
	defer db.Close()

	for _, key := range storeKeys {
		log.Printf("Upgrading.... %v\n", key)
		prefix := []byte(fmt.Sprintf("s/k:%s/", key))

		prefixDB := dbm.NewPrefixDB(db, prefix)

		tree, err := iavl.NewMutableTree(prefixDB, 0)
		if err != nil {
			return err
		}
		_, err = tree.LoadVersion(0)
		if err != nil {
			return err
		}
	}

	return nil
}
