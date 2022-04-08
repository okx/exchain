package main

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/okex/exchain/app/logevents"
	"github.com/okex/exchain/cmd/exchaind/mpt"

	"github.com/okex/exchain/app/rpc"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tmamino "github.com/okex/exchain/libs/tendermint/crypto/encoding/amino"
	"github.com/okex/exchain/libs/tendermint/crypto/multisig"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
	dbm "github.com/okex/exchain/libs/tm-db"

	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	clientkeys "github.com/okex/exchain/libs/cosmos-sdk/client/keys"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"

	"github.com/okex/exchain/app"
	"github.com/okex/exchain/app/codec"
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

	codecProxy, _ := codec.MakeCodecSuit(app.ModuleBasics)

	tmamino.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)
	tmamino.RegisterKeyType(ethsecp256k1.PrivKey{}, ethsecp256k1.PrivKeyName)
	multisig.RegisterKeyType(ethsecp256k1.PubKey{}, ethsecp256k1.PubKeyName)

	keys.CryptoCdc = codecProxy.GetCdc()
	genutil.ModuleCdc = codecProxy.GetCdc()
	genutiltypes.ModuleCdc = codecProxy.GetCdc()
	clientkeys.KeysCdc = codecProxy.GetCdc()

	config := sdk.GetConfig()
	okexchain.SetBech32Prefixes(config)
	okexchain.SetBip44CoinType(config)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "exchaind",
		Short:             "ExChain App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		client.ValidateChainID(
			genutilcli.InitCmd(ctx, codecProxy.GetCdc(), app.ModuleBasics, app.DefaultNodeHome),
		),
		genutilcli.CollectGenTxsCmd(ctx, codecProxy.GetCdc(), auth.GenesisAccountIterator{}, app.DefaultNodeHome),
		genutilcli.MigrateGenesisCmd(ctx, codecProxy.GetCdc()),
		genutilcli.GenTxCmd(
			ctx, codecProxy.GetCdc(), app.ModuleBasics, staking.AppModuleBasic{}, auth.GenesisAccountIterator{},
			app.DefaultNodeHome, app.DefaultCLIHome,
		),
		genutilcli.ValidateGenesisCmd(ctx, codecProxy.GetCdc(), app.ModuleBasics),
		client.TestnetCmd(ctx, codecProxy.GetCdc(), app.ModuleBasics, auth.GenesisAccountIterator{}),
		replayCmd(ctx, client.RegisterAppFlag),
		repairStateCmd(ctx),
		displayStateCmd(ctx),
		mpt.MptCmd(ctx),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		AddGenesisAccountCmd(ctx, codecProxy.GetCdc(), app.DefaultNodeHome, app.DefaultCLIHome),
		flags.NewCompletionCmd(rootCmd, true),
		dataCmd(ctx),
		exportAppCmd(ctx),
		iaviewerCmd(ctx, codecProxy.GetCdc()),
		subscribeCmd(codecProxy.GetCdc()),
	)

	subFunc := func(logger log.Logger) log.Subscriber {
		return logevents.NewProvider(logger)
	}
	// Tendermint node base commands
	server.AddCommands(ctx, codecProxy, rootCmd, newApp, closeApp, exportAppStateAndTMValidators,
		registerRoutes, client.RegisterAppFlag, app.PreRun, subFunc)

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
	app.StopBaseApp()
	evmtypes.CloseIndexer()
	rpc.CloseEthBackend()
}

func newApp(logger log.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
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
		}
	} else {
		ethermintApp = app.NewOKExChainApp(logger, db, traceStore, true, map[int64]bool{}, 0)
	}

	return ethermintApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
