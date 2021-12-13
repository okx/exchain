package main

import (
	"log"

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
	return cmd
}
