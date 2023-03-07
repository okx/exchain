package main

import (
	"fmt"
	"github.com/okx/okbchain/libs/system"

	interfacetypes "github.com/okx/okbchain/libs/cosmos-sdk/codec/types"

	authtypes "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/types"

	"github.com/okx/okbchain/app"
	"github.com/okx/okbchain/app/codec"
	"github.com/okx/okbchain/app/crypto/ethsecp256k1"
	chain "github.com/okx/okbchain/app/types"
	"github.com/okx/okbchain/cmd/client"
	sdkclient "github.com/okx/okbchain/libs/cosmos-sdk/client"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/flags"
	clientkeys "github.com/okx/okbchain/libs/cosmos-sdk/client/keys"
	clientrpc "github.com/okx/okbchain/libs/cosmos-sdk/client/rpc"
	sdkcodec "github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/libs/cosmos-sdk/crypto/keys"
	"github.com/okx/okbchain/libs/cosmos-sdk/server"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/version"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	authcmd "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/cli"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/bank"
	tmamino "github.com/okx/okbchain/libs/tendermint/crypto/encoding/amino"
	"github.com/okx/okbchain/libs/tendermint/crypto/multisig"
	"github.com/okx/okbchain/libs/tendermint/libs/cli"
	evmtypes "github.com/okx/okbchain/x/evm/types"
	tokencmd "github.com/okx/okbchain/x/token/client/cli"
	"github.com/spf13/cobra"
)

var (
	cdc          = codec.MakeCodec(app.ModuleBasics)
	interfaceReg = codec.MakeIBC(app.ModuleBasics)
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
	chain.SetBech32Prefixes(config)
	chain.SetBip44CoinType(config)
	config.Seal()

	rootCmd := &cobra.Command{
		Use:   "okbchaincli",
		Short: "Command line interface for interacting with okbchaind",
	}

	// Add --chain-id to persistent flags and mark it required
	rootCmd.PersistentFlags().String(flags.FlagChainID, server.ChainID, "Chain ID of tendermint node")
	rootCmd.PersistentPreRunE = func(_ *cobra.Command, _ []string) error {
		utils.SetParseAppTx(wrapDecoder(parseMsgEthereumTx, parseProtobufTx))
		return client.InitConfig(rootCmd)
	}
	protoCdc := sdkcodec.NewProtoCodec(interfaceReg)
	proxy := sdkcodec.NewCodecProxy(protoCdc, cdc)
	// Construct Root Command
	rootCmd.AddCommand(
		clientrpc.StatusCommand(),
		sdkclient.ConfigCmd(app.DefaultCLIHome),
		queryCmd(proxy, interfaceReg),
		txCmd(proxy, interfaceReg),
		flags.LineBreak,
		client.KeyCommands(),
		client.AddrCommands(),
		flags.LineBreak,
		version.Cmd,
		flags.NewCompletionCmd(rootCmd, true),
	)

	// Add flags and prefix all env exposed with OKBCHAIN
	executor := cli.PrepareMainCmd(rootCmd, system.EnvPrefix, app.DefaultCLIHome)

	err := executor.Execute()
	if err != nil {
		panic(fmt.Errorf("failed executing CLI command: %w", err))
	}
}

func queryCmd(proxy *sdkcodec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}
	cdc := proxy.GetCdc()
	queryCmd.AddCommand(
		authcmd.GetAccountCmd(cdc),
		flags.LineBreak,
		authcmd.QueryTxsByEventsCmd(cdc),
		authcmd.QueryTxCmd(proxy),
		flags.LineBreak,
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)
	app.ModuleBasics.AddQueryCommandsV2(queryCmd, proxy, reg)

	return queryCmd
}

func txCmd(proxy *sdkcodec.CodecProxy, reg interfacetypes.InterfaceRegistry) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}
	cdc := proxy.GetCdc()
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
	app.ModuleBasics.AddTxCommandsV2(txCmd, proxy, reg)

	// remove auth and bank commands as they're mounted under the root tx command
	var cmdsToRemove []*cobra.Command

	for _, cmd := range txCmd.Commands() {
		if cmd.Use == auth.ModuleName ||
			cmd.Use == bank.ModuleName {
			cmdsToRemove = append(cmdsToRemove, cmd)
		}
	}

	txCmd.RemoveCommand(cmdsToRemove...)

	return txCmd
}

func wrapDecoder(handlers ...utils.ParseAppTxHandler) utils.ParseAppTxHandler {
	return func(cdc *sdkcodec.CodecProxy, txBytes []byte) (sdk.Tx, error) {
		var (
			tx  sdk.Tx
			err error
		)
		for _, handler := range handlers {
			tx, err = handler(cdc, txBytes)
			if nil == err && tx != nil {
				return tx, err
			}
		}
		return tx, err
	}
}
func parseMsgEthereumTx(cdc *sdkcodec.CodecProxy, txBytes []byte) (sdk.Tx, error) {
	var tx evmtypes.MsgEthereumTx
	// try to decode through RLP first
	if err := authtypes.EthereumTxDecode(txBytes, &tx); err == nil {
		return &tx, nil
	}
	//try to decode through animo if it is not RLP-encoded
	if err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx); err != nil {
		return nil, err
	}
	return &tx, nil
}

func parseProtobufTx(cdc *sdkcodec.CodecProxy, txBytes []byte) (sdk.Tx, error) {
	tx, err := evmtypes.TxDecoder(cdc)(txBytes, evmtypes.IGNORE_HEIGHT_CHECKING)
	if nil != err {
		return nil, err
	}
	switch realTx := tx.(type) {
	case *authtypes.IbcTx:
		return authtypes.FromProtobufTx(cdc, realTx)
	}
	return tx, err
}
