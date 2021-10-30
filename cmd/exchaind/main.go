package main

import (
	"encoding/json"
	"fmt"
	"github.com/okex/exchain/app/rpc"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
	"github.com/okex/exchain/libs/tendermint/crypto/multisig"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	clientkeys "github.com/okex/exchain/libs/cosmos-sdk/client/keys"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/codec"
	appconfig "github.com/okex/exchain/app/config"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	okexchain "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/cmd/client"
	"github.com/okex/exchain/x/genutil"
	genutilcli "github.com/okex/exchain/x/genutil/client/cli"
	genutiltypes "github.com/okex/exchain/x/genutil/types"
	"github.com/okex/exchain/x/staking"
)

const flagInvCheckPeriod = "inv-check-period"

var invCheckPeriod uint

func main() {
	cobra.EnableCommandSorting = false

	cdc := codec.MakeCodec(app.ModuleBasics)

	tmamino.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)
	tmamino.RegisterKeyType(ethsecp256k1.PrivKey{}, ethsecp256k1.PrivKeyName)
	multisig.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)

	keys.CryptoCdc = cdc
	genutil.ModuleCdc = cdc
	genutiltypes.ModuleCdc = cdc
	clientkeys.KeysCdc = cdc

	config := sdk.GetConfig()
	okexchain.SetBech32Prefixes(config)
	okexchain.SetBip44CoinType(config)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "exchaind",
		Short:             "ExChain App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx, appconfig.RegisterDynamicConfig),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		client.ValidateChainID(
			genutilcli.InitCmd(ctx, cdc, app.ModuleBasics, app.DefaultNodeHome),
		),
		genutilcli.CollectGenTxsCmd(ctx, cdc, auth.GenesisAccountIterator{}, app.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(ctx, cdc),
		genutilcli.GenTxCmd(
			ctx, cdc, app.ModuleBasics, staking.AppModuleBasic{}, auth.GenesisAccountIterator{},
			app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, cdc, app.ModuleBasics),
		client.TestnetCmd(ctx, cdc, app.ModuleBasics, auth.GenesisAccountIterator{}),
		replayCmd(ctx),
		repairStateCmd(ctx),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		AddGenesisAccountCmd(ctx, cdc, app.DefaultNodeHome, app.DefaultCLIHome),
		flags.NewCompletionCmd(rootCmd, true),
		dataCmd(ctx),
		exportAppCmd(ctx),
		iaviewerCmd(cdc),
	)

	// Tendermint node base commands
	server.AddCommands(ctx, cdc, rootCmd, newApp, closeApp, exportAppStateAndTMValidators, registerRoutes, client.RegisterAppFlag)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "OKEXCHAIN", app.DefaultNodeHome)
	rootCmd.PersistentFlags().UintVar(&invCheckPeriod, flagInvCheckPeriod,
		0, "Assert registered invariants every N blocks")
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func closeApp(iApp abci.Application) {
	fmt.Println("Close App")
	app := iApp.(*app.OKExChainApp)
	app.StopStore()
	evmtypes.CloseIndexer()
	rpc.CloseEthBackend()
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) (abci.Application) {
	pruningOpts, err := server.GetPruningOptionsFromFlags()
	if err != nil {
		panic(err)
	}

	return app.NewOKExChainApp(
		logger,
		db,
		traceStore,
		true,
		map[int64]bool{},
		0,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
		baseapp.SetHaltHeight(uint64(viper.GetInt(server.FlagHaltHeight))),
	)
}

func exportAppStateAndTMValidators(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	var ethermintApp *app.OKExChainApp
	if height != -1 {
		ethermintApp = app.NewOKExChainApp(logger, db, traceStore, false, map[int64]bool{}, 0)

		if err := ethermintApp.LoadHeight(height); err != nil {
			return nil, nil, err
		} else {
			ethermintApp = app.NewOKExChainApp(logger, db, traceStore, true, map[int64]bool{}, 0)
		}
	}

	return ethermintApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
