package fss

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/okx/okbchain/app/utils/appstatus"
	"github.com/okx/okbchain/cmd/okbchaind/base"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/iavl"
	dbm "github.com/okx/okbchain/libs/tm-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create fast index for IAVL",
	Long: `Create fast index for IAVL:
This command is a tool to generate the IAVL fast index.
It will take long based on the original database size.
When the create lunched, it will show Upgrade to Fast IAVL...`,
	RunE: func(cmd *cobra.Command, args []string) error {
		iavl.SetEnableFastStorage(true)
		storeKeys := appstatus.GetAllStoreKeys()
		outputModules(storeKeys)

		return createIndex(storeKeys)
	},
}

func init() {
	fssCmd.AddCommand(createCmd)
}

func outputModules(storeKeys []string) {
	if iavl.OutputModules == nil {
		iavl.OutputModules = make(map[string]int, len(storeKeys))
	}
	for _, key := range storeKeys {
		iavl.OutputModules[key] = 1
	}
}

func createIndex(storeKeys []string) error {
	dataDir := viper.GetString(flagDataDir)
	dbBackend := viper.GetString(sdk.FlagDBBackend)
	db, err := base.OpenDB(filepath.Join(dataDir, base.AppDBName), dbm.BackendType(dbBackend))
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
