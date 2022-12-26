package fss

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/okex/exchain/cmd/exchaind/base"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/iavl"
	dbm "github.com/okex/exchain/libs/tm-db"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fssCmd.AddCommand(versionCmd)
}

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version  <module>",
	Short: "Read the module's fast storage version",
	Long: `Read the module's fast storage version:
This command is a tool to read the fast storage version.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("must specify module")
		}
		return version(args[0])
	},
}

func version(module string) error {
	dataDir := viper.GetString(flagDataDir)
	dbBackend := viper.GetString(sdk.FlagDBBackend)
	db, err := base.OpenDB(filepath.Join(dataDir, base.AppDBName), dbm.BackendType(dbBackend))
	if err != nil {
		return fmt.Errorf("error opening dir %v backend %v DB: %w", dataDir, dbBackend, err)
	}
	defer db.Close()

	prefix := []byte(fmt.Sprintf("s/k:%s/", module))
	prefixDB := dbm.NewPrefixDB(db, prefix)

	version, err := iavl.GetFastStorageVersion(prefixDB)
	if err != nil {
		return err
	}
	log.Printf("module: %v version: %v\n", module, version)

	return nil
}
