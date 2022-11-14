package fss

import (
	"github.com/okex/exchain/app/utils/appstatus"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/iavl"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagDataDir = "data_dir"
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
		storeKeys := appstatus.GetAllStoreKeys()
		outputModules(storeKeys)
	},
}

func init() {
	fssCmd.PersistentFlags().StringP(flagDataDir, "d", "./", "The chain data file location")
	fssCmd.PersistentFlags().String(sdk.FlagDBBackend, tmtypes.DBBackend, "Database backend: goleveldb | rocksdb")
	viper.BindPFlag(flagDataDir, fssCmd.PersistentFlags().Lookup(flagDataDir))
	viper.BindPFlag(sdk.FlagDBBackend, fssCmd.PersistentFlags().Lookup(sdk.FlagDBBackend))
}

func outputModules(storeKeys []string) {
	if iavl.OutputModules == nil {
		iavl.OutputModules = make(map[string]int, len(storeKeys))
	}
	for _, key := range storeKeys {
		iavl.OutputModules[key] = 1
	}
}
