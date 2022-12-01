package fss

import (
	"encoding/hex"
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
	fssCmd.AddCommand(readCmd)
}

// checkCmd represents the create command
var readCmd = &cobra.Command{
	Use:   "read <module> <key>",
	Short: "Read the fast node",
	Long: `Read fast node value:
This command is a tool to read the IAVL fast node.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("must specify module and key")
		}
		return read(args[0], args[1])
	},
}

func read(module, key string) error {
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

	prefix := []byte(fmt.Sprintf("s/k:%s/", module))
	prefixDB := dbm.NewPrefixDB(db, prefix)

	mutableTree, err := iavl.NewMutableTree(prefixDB, 0)
	if err != nil {
		return err
	}
	if _, err := mutableTree.LoadVersion(0); err != nil {
		return err
	}

	keyByte, err := hex.DecodeString(key)
	if err != nil {
		return fmt.Errorf("error decoding key: %v %w", key, err)
	}
	v := mutableTree.Get(keyByte)
	log.Printf("key :%v\n", key)
	log.Printf("value :%x\n", v)

	return nil
}
