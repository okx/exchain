package mpt

import (
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/mpt/types"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	"github.com/spf13/cobra"
)

func MptCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mpt",
		Short: "migrate iavl state to mpt state (if use migrate mpt data, then you should set `--use-composite-key true` when you decide to use mpt to store the coming data)",
	}

	cmd.AddCommand(
		iavl2mptCmd(ctx),
		cleanIavlStoreCmd(ctx),
		mptViewerCmd(ctx),
		AccountGetCmd(ctx),
	)
	cmd.PersistentFlags().UintVar(&types.TrieRocksdbBatchSize, types.FlagTrieRocksdbBatchSize, 100, "Concurrent rocksdb batch size for mpt")
	cmd.PersistentFlags().String(sdk.FlagDBBackend, tmtypes.DBBackend, "Database backend: goleveldb | rocksdb")

	return cmd
}
