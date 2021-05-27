package rpc

import (
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/spf13/cobra"
)

// ServeCmd creates a CLI command to start Cosmos REST server with web3 RPC API and
// Cosmos rest-server endpoints
func ServeCmd(cdc *codec.Codec) *cobra.Command {
	cmd := lcd.ServeCommand(cdc, RegisterRoutes)
	cmd.Flags().String(flagUnlockKey, "", "Select a key to unlock on the RPC server")
	cmd.Flags().String(flagWebsocket, "8546", "websocket port to listen to")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block)")

	cmd.Flags().Bool(watcher.FlagFastQuery, false, "Enable the fast query mode for rpc queries")
	cmd.Flags().String(watcher.FlagWatcherDBType, watcher.DBTypeLevel, "config watcher db")
	cmd.Flags().String(watcher.FlagWatcherDisLockUrl, "redis://127.0.0.1:6379", "config watcher dis lock url")
	cmd.Flags().String(watcher.FlagWatcherDisLockUrlPassword, "", "config watcher dis lock password")
	cmd.Flags().String(watcher.FlagHbaseDBUrl, "", "config watcher hbase db url")

	return cmd
}
