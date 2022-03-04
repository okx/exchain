package main

import (
	"github.com/okex/exchain/x/common/analyzer"
	types2 "github.com/okex/exchain/x/evm/types"
	"log"

	"github.com/okex/exchain/libs/cosmos-sdk/store/flatkv"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sm "github.com/okex/exchain/libs/tendermint/state"
	"github.com/spf13/cobra"
)

func repairStateCmd(ctx *server.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repair-state",
		Short: "Repair the SMB(state machine broken) data of node",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("--------- repair data start ---------")

			app.RepairState(ctx, false)
			log.Println("--------- repair data success ---------")
		},
	}
	cmd.Flags().Bool(sm.FlagParalleledTx, false, "parallel execution for evm txs")
	cmd.Flags().Int64(app.FlagStartHeight, 0, "Set the start block height for repair")
	cmd.Flags().Bool(flatkv.FlagEnable, false, "Enable flat kv storage for read performance")
	cmd.Flags().String(app.Elapsed, app.DefaultElapsedSchemas, "schemaName=1|0,,,")
	cmd.Flags().Bool(analyzer.FlagEnableAnalyzer, true, "Enable auto open log analyzer")
	cmd.Flags().BoolVar(&types2.UseCompositeKey, types2.FlagUseCompositeKey,false, "Use composite key to store contract state")

	return cmd
}
