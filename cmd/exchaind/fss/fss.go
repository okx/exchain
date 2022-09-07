package fss

import (
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/iavl"
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
}

func init() {
	fssCmd.PersistentFlags().StringP(flagDataDir, "d", "./", "The chain data file location")
	fssCmd.PersistentFlags().String(flagDBBackend, "goleveldb", "Database backend: goleveldb | rocksdb")
}
