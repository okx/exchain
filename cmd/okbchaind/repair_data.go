package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/okx/okbchain/app"
	"github.com/okx/okbchain/app/utils/appstatus"
	"github.com/okx/okbchain/libs/cosmos-sdk/server"
	"github.com/okx/okbchain/libs/cosmos-sdk/store/flatkv"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	tmiavl "github.com/okx/okbchain/libs/iavl"
	"github.com/okx/okbchain/libs/system/trace"
	sm "github.com/okx/okbchain/libs/tendermint/state"
	tmtypes "github.com/okx/okbchain/libs/tendermint/types"
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

			go func() {
				pprofAddress := viper.GetString(pprofAddrFlag)
				err := http.ListenAndServe(pprofAddress, nil)
				if err != nil {
					fmt.Println(err)
				}
			}()
			app.RepairState(ctx, false)
			log.Println("--------- repair data success ---------")
		},
	}
	cmd.Flags().Int64(app.FlagStartHeight, 0, "Set the start block height for repair")
	cmd.Flags().Bool(flatkv.FlagEnable, false, "Enable flat kv storage for read performance")
	cmd.Flags().String(app.Elapsed, app.DefaultElapsedSchemas, "schemaName=1|0,,,")
	cmd.Flags().Bool(trace.FlagEnableAnalyzer, false, "Enable auto open log analyzer")
	cmd.Flags().Int(sm.FlagDeliverTxsExecMode, 0, "execution mode for deliver txs, (0:serial[default], 1:deprecated, 2:parallel)")
	cmd.Flags().String(sdk.FlagDBBackend, tmtypes.DBBackend, "Database backend: goleveldb | rocksdb")
	cmd.Flags().StringP(pprofAddrFlag, "p", "0.0.0.0:6060", "Address and port of pprof HTTP server listening")
	cmd.Flags().Bool(tmiavl.FlagIavlDiscardFastStorage, false, "Discard fast storage")
	cmd.Flags().MarkHidden(tmiavl.FlagIavlDiscardFastStorage)

	return cmd
}

func setExternalPackageValue() {
	tmiavl.SetForceReadIavl(true)
	isFastStorage := appstatus.IsFastStorageStrategy()
	tmiavl.SetEnableFastStorage(isFastStorage)
	if !isFastStorage &&
		!viper.GetBool(tmiavl.FlagIavlDiscardFastStorage) &&
		appstatus.GetFastStorageVersion() >= viper.GetInt64(app.FlagStartHeight) {
		tmiavl.SetEnableFastStorage(true)
		tmiavl.SetIgnoreAutoUpgrade(true)
	}
}
