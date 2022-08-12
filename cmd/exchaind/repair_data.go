package main

import (
	"log"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/cosmos-sdk/store/flatkv"
	tmiavl "github.com/okex/exchain/libs/iavl"
	"github.com/okex/exchain/libs/system/trace"
	sm "github.com/okex/exchain/libs/tendermint/state"
	types2 "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func repairStateCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair-state",
		Short: "Repair the SMB(state machine broken) data of node",
		PreRun: func(_ *cobra.Command, _ []string) {
			setExternalPackageValue()
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- repair data start ---------")

			app.RepairState(ctx, false)
			log.Println("--------- repair data success ---------")
		},
	}
	cmd.Flags().Int64(app.FlagStartHeight, 0, "Set the start block height for repair")
	cmd.Flags().Bool(flatkv.FlagEnable, false, "Enable flat kv storage for read performance")
	cmd.Flags().String(app.Elapsed, app.DefaultElapsedSchemas, "schemaName=1|0,,,")
	cmd.Flags().Bool(trace.FlagEnableAnalyzer, false, "Enable auto open log analyzer")
	cmd.Flags().BoolVar(&types2.TrieUseCompositeKey, types2.FlagTrieUseCompositeKey, true, "Use composite key to store contract state")
	cmd.Flags().Int(sm.FlagDeliverTxsExecMode, 0, "execution mode for deliver txs, (0:serial[default], 1:deprecated, 2:parallel)")
	cmd.Flags().Bool(tmiavl.FlagIavlEnableFastStorage, false, "Enable fast storage")
	cmd.Flags().Int(tmiavl.FlagIavlFastStorageCacheSize, 100000, "Max size of iavl fast storage cache")

	return cmd
}

func setExternalPackageValue() {
	tmiavl.SetEnableFastStorage(viper.GetBool(tmiavl.FlagIavlEnableFastStorage))
	tmiavl.SetFastNodeCacheSize(viper.GetInt(tmiavl.FlagIavlFastStorageCacheSize))
}
