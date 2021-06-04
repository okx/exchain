package main

import (
	"fmt"
	sdkclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clientkeys "github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	clientrpc "github.com/cosmos/cosmos-sdk/client/rpc"
	sdkcdc "github.com/cosmos/cosmos-sdk/codec"
	sdkcodec "github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/okex/exchain/app/rpc"
	"github.com/spf13/cobra"
	tmamino "github.com/tendermint/tendermint/crypto/encoding/amino"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/codec"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	okexchain "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/cmd/client"
	"github.com/okex/exchain/x/dex"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/okex/exchain/x/order"
	tokencmd "github.com/okex/exchain/x/token/client/cli"
)

var (
	cdc = codec.MakeCodec(app.ModuleBasics)
)

func main() {
	// Configure cobra to sort commands
	cobra.EnableCommandSorting = false

	tmamino.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)
	tmamino.RegisterKeyType(ethsecp256k1.PrivKey{}, ethsecp256k1.PrivKeyName)
	multisig.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)

	keys.CryptoCdc = cdc
	clientkeys.KeysCdc = cdc

	// Read in the configuration file for the sdk
	config := sdk.GetConfig()
	okexchain.SetBech32Prefixes(config)
	okexchain.SetBip44CoinType(config)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "exchaincli",
		Short: "Command line interface for interacting with exchaind",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, "", "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		return client.InitConfig(rootCmd)
	}

	// Construct Root Command
	rootCmd.AddCommand(
		clientrpc.StatusCommand(),
		sdkclient.ConfigCmd(app.DefaultCLIHome),
		queryCmd(cdc),
		txCmd(cdc),
		client.ValidateChainID(
			ServeCmd(cdc),
		),
		flags.LineBreak,
		client.KeyCommands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCmd, true),
	)

	// Add flags and prefix all env exposed with OKEXCHAIN
	executor := cli.PrepareMainCmd(rootCmd, "OKEXCHAIN", app.DefaultCLIHome)

	err := executor.Execute()
	if err != nil {
		panic(fmt.Errorf("failed executing CLI command: %w", err))
	}
}

func queryCmd(cdc *sdkcodec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		flags.LineBreak,
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(cdc),
		flags.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *sdkcodec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		tokencmd.SendTxCmd(cdc),
		flags.LineBreak,
		authcmd.GetSignCommand(cdc),
		authcmd.GetMultiSignCommand(cdc),
		flags.LineBreak,
		authcmd.GetBroadcastCommand(cdc),
		authcmd.GetEncodeCommand(cdc),
		authcmd.GetDecodeCommand(cdc),
		flags.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName ||
			cmd.Use == order.ModuleName ||
			cmd.Use == dex.ModuleName ||
			cmd.Use == bank.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

// ServeCmd creates a CLI command to start Cosmos REST server with web3 RPC API and
// Cosmos rest-server endpoints
func ServeCmd(cdc *sdkcdc.Codec) *cobra.Command {
	cmd := lcd.ServeCommand(cdc, client.RegisterRoutes)
	cmd.Flags().Bool(watcher.FlagFastQuery, false, "Enable the fast query mode for rpc queries")
	cmd.Flags().String(watcher.FlagWatcherDBType, watcher.DBTypeLevel, "config watcher db")
	cmd.Flags().String(watcher.FlagWatcherDisLockUrl, "redis://127.0.0.1:6379", "config watcher dis lock url")
	cmd.Flags().String(watcher.FlagWatcherDisLockUrlPassword, "", "config watcher dis lock password")
	cmd.Flags().Int(watcher.FlagFastQueryLru, 1000, "Set the size of LRU cache under fast-query mode")
	cmd.Flags().String(watcher.FlagWatcherDBUrl, "", "config watcher db url")
	cmd.Flags().String(watcher.FlagWatcherDBPassword, "", "config watcher db password")
	cmd.Flags().Int(watcher.FlagFastQueryLru, 1000, "Set the size of LRU cache under fast-query mode")

	// rest-server flags for nacos config
	cmd.Flags().Bool(rpc.FlagRestServerServiceEnable, false, "Rpc`s rest-server service enable config")
	cmd.Flags().String(rpc.FlagRestServerNacosUrls, "", "Rpc`s nacos server urls for getting rest-server service info")
	cmd.Flags().String(rpc.FlagRestServerNacosNamespaceId, "", "Rpc`s nacos name space id for getting rest-server service info")
	cmd.Flags().StringArray(rpc.FlagRestServerNacosClusters, []string{}, "Rpc`s nacos clusters array list for getting rest-server service info")
	cmd.Flags().String(rpc.FlagRestServerNacosServiceName, "", "Rpc`s nacos service name for getting rest-server service info")
	cmd.Flags().String(rpc.FlagRestServerNacosGroupName, "", "Rpc`s nacos group name for getting rest-server service info")

	cmd.Flags().String(server.FlagListenAddr, "tcp://0.0.0.0:26659", "The address for the rest-server to listen on. (0.0.0.0:0 means any interface, any port)")
	cmd.Flags().String(server.FlagUlockKey, "", "Select the keys to unlock on the RPC server")
	cmd.Flags().String(server.FlagUlockKeyHome, "", "The keybase home path")
	cmd.Flags().String(server.FlagRestPathPrefix, "exchain", "Path prefix for registering rest api route.")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(server.FlagCORS, "", "Set the rest-server domains that can make CORS requests (* for all)")
	cmd.Flags().String(server.FlagExternalListenAddr, "127.0.0.1:26659", "Set the rest-server external ip and port, when it is launched by Docker")
	cmd.Flags().String(server.FlagWebsocket, "8546", "websocket port to listen to")
	cmd.Flags().Int(server.FlagWsMaxConnections, 20000, "the max capacity number of websocket client connections")
	cmd.Flags().Int(server.FlagWsSubChannelLength, 100, "the length of subscription channel")
	cmd.Flags().StringP(flags.FlagBroadcastMode, "b", flags.BroadcastSync, "Transaction broadcasting mode (sync|async|block) for web3")

	return cmd
}
