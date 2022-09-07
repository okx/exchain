package fss

import (
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	"github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/x/evm"
	"github.com/spf13/cobra"
)

const (
	flagDataDir   = "data_dir"
	flagDBBackend = "db_backend"
)

func Command(ctx *server.Context) *cobra.Command {
	iavl.SetLogger(ctx.Logger.With("module", "iavl"))
	return fssCmd
}

var fssCmd = &cobra.Command{
	Use:   "fss",
	Short: "FSS is an auxiliary fast storage system to IAVL",
	Long: `IAVL fast storage related commands:
This command include a set of command of the IAVL fast storage.
include create sub command`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		iavl.SetEnableFastStorage(true)
		storeKeys := getStoreKeys()
		outputModules(storeKeys)
	},
}

func init() {
	fssCmd.PersistentFlags().StringP(flagDataDir, "d", "./", "The chain data file location")
	fssCmd.PersistentFlags().String(flagDBBackend, "goleveldb", "Database backend: goleveldb | rocksdb")
}

func getStoreKeys() []string {
	return []string{
		auth.StoreKey,
		evm.StoreKey,
	}
}

func outputModules(storeKeys []string) {
	if iavl.OutputModules == nil {
		iavl.OutputModules = make(map[string]int, len(storeKeys))
	}
	for _, key := range storeKeys {
		iavl.OutputModules[key] = 1
	}
}
