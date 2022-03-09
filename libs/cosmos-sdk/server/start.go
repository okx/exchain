package server

// DONTCOVER

import (
	"os"
	"runtime/pprof"

	"github.com/okex/exchain/libs/cosmos-sdk/baseapp"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/client/lcd"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/store/iavl"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/types"
	"github.com/okex/exchain/libs/system"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/mempool"
	"github.com/okex/exchain/libs/tendermint/node"
	"github.com/okex/exchain/libs/tendermint/p2p"
	"github.com/okex/exchain/libs/tendermint/proxy"
	"github.com/okex/exchain/libs/tendermint/rpc/client/local"
	"github.com/okex/exchain/libs/tendermint/state"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	tmiavl "github.com/okex/exchain/libs/iavl"
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	tcmd "github.com/okex/exchain/libs/tendermint/cmd/tendermint/commands"
	tmos "github.com/okex/exchain/libs/tendermint/libs/os"
	pvm "github.com/okex/exchain/libs/tendermint/privval"
	tmtypes "github.com/okex/exchain/libs/tendermint/types"
)

// Tendermint full-node start flags
const (
	flagAddress            = "address"
	flagTraceStore         = "trace-store"
	flagCPUProfile         = "cpu-profile"
	FlagMinGasPrices       = "minimum-gas-prices"
	FlagHaltHeight         = "halt-height"
	FlagHaltTime           = "halt-time"
	FlagInterBlockCache    = "inter-block-cache"
	FlagUnsafeSkipUpgrades = "unsafe-skip-upgrades"
	FlagTrace              = "trace"

	FlagPruning           = "pruning"
	FlagPruningKeepRecent = "pruning-keep-recent"
	FlagPruningKeepEvery  = "pruning-keep-every"
	FlagPruningInterval   = "pruning-interval"
	FlagLocalRpcPort      = "local-rpc-port"
	FlagPortMonitor       = "netstat"
	FlagEvmImportPath     = "evm-import-path"
	FlagEvmImportMode     = "evm-import-mode"
	FlagGoroutineNum      = "goroutine-num"

	FlagPruningMaxWsNum = "pruning-max-worldstate-num"
	FlagExportKeystore  = "export-keystore"
	FlagLogServerUrl    = "log-server"
)

// StartCmd runs the service passed in, either stand-alone or in-process with
// Tendermint.
func StartCmd(ctx *Context,
	cdc *codec.Codec,
	appCreator AppCreator,
	appStop AppStop,
	registerRoutesFn func(restServer *lcd.RestServer),
	registerAppFlagFn func(cmd *cobra.Command),
	appPreRun func(ctx *Context) error,
	subFunc func(logger log.Logger) log.Subscriber,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run the full node",
		Long: `Run the full node application with Tendermint in or out of process. By
default, the application will run with Tendermint in process.

Pruning options can be provided via the '--pruning' flag or alternatively with '--pruning-keep-recent',
'pruning-keep-every', and 'pruning-interval' together.

For '--pruning' the options are as follows:

default: the last 100 states are kept in addition to every 500th state; pruning at 10 block intervals
nothing: all historic states will be saved, nothing will be deleted (i.e. archiving node)
everything: all saved states will be deleted, storing only the current state; pruning at 10 block intervals
custom: allow pruning options to be manually specified through 'pruning-keep-recent', 'pruning-keep-every', and 'pruning-interval'

Node halting configurations exist in the form of two flags: '--halt-height' and '--halt-time'. During
the ABCI Commit phase, the node will check if the current block height is greater than or equal to
the halt-height or if the current block time is greater than or equal to the halt-time. If so, the
node will attempt to gracefully shutdown and the block will not be committed. In addition, the node
will not be able to commit subsequent blocks.

For profiling and benchmarking purposes, CPU profiling can be enabled via the '--cpu-profile' flag
which accepts a path for the resulting pprof file.
`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// app pre run
			if err := appPreRun(ctx); err != nil {
				return err
			}
			// set external package flags
			SetExternalPackageValue(cmd)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx.Logger.Info("Starting ABCI with Tendermint")

			sub := subFunc(ctx.Logger)
			log.SetSubscriber(sub)

			setPID(ctx)
			_, err := startInProcess(ctx, cdc, appCreator, appStop, registerRoutesFn)
			if err != nil {
				tmos.Exit(err.Error())
			}
			return nil
		},
	}

	RegisterServerFlags(cmd)
	registerAppFlagFn(cmd)
	// add support for all Tendermint-specific command line options
	tcmd.AddNodeFlags(cmd)
	cmd.AddCommand(nodeModeCmd(ctx))
	return cmd
}

func startInProcess(ctx *Context, cdc *codec.Codec, appCreator AppCreator, appStop AppStop,
	registerRoutesFn func(restServer *lcd.RestServer)) (*node.Node, error) {

	cfg := ctx.Config
	home := cfg.RootDir
	//startInProcess hooker
	callHooker(FlagHookstartInProcess, ctx)

	traceWriterFile := viper.GetString(flagTraceStore)
	db, err := openDB(home)
	if err != nil {
		return nil, err
	}

	traceWriter, err := openTraceWriter(traceWriterFile)
	if err != nil {
		return nil, err
	}

	app := appCreator(ctx.Logger, db, traceWriter)

	nodeKey, err := p2p.LoadOrGenNodeKey(cfg.NodeKeyFile())
	if err != nil {
		return nil, err
	}

	// create & start tendermint node
	tmNode, err := node.NewNode(
		cfg,
		pvm.LoadOrGenFilePV(cfg.PrivValidatorKeyFile(), cfg.PrivValidatorStateFile()),
		nodeKey,
		proxy.NewLocalClientCreator(app),
		node.DefaultGenesisDocProviderFunc(cfg),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(cfg.Instrumentation),
		ctx.Logger.With("module", "node"),
	)
	if err != nil {
		return nil, err
	}

	app.SetOption(abci.RequestSetOption{
		Key:   "CheckChainID",
		Value: tmNode.ConsensusState().GetState().ChainID,
	})

	ctx.Logger.Info("startInProcess",
		"ConsensusStateChainID", tmNode.ConsensusState().GetState().ChainID,
		"GenesisDocChainID", tmNode.GenesisDoc().ChainID,
	)
	if err := tmNode.Start(); err != nil {
		return nil, err
	}

	var cpuProfileCleanup func()

	if cpuProfile := viper.GetString(flagCPUProfile); cpuProfile != "" {
		f, err := os.Create(cpuProfile)
		if err != nil {
			return nil, err
		}

		ctx.Logger.Info("starting CPU profiler", "profile", cpuProfile)
		if err := pprof.StartCPUProfile(f); err != nil {
			return nil, err
		}

		cpuProfileCleanup = func() {
			ctx.Logger.Info("stopping CPU profiler", "profile", cpuProfile)
			pprof.StopCPUProfile()
			f.Close()
		}
	}

	TrapSignal(func() {
		if tmNode.IsRunning() {
			_ = tmNode.Stop()
		}
		appStop(app)

		if cpuProfileCleanup != nil {
			cpuProfileCleanup()
		}

		ctx.Logger.Info("exiting...")
	})

	if registerRoutesFn != nil {
		go lcd.StartRestServer(cdc, registerRoutesFn, tmNode, viper.GetString(FlagListenAddr))
	}

	baseapp.SetGlobalMempool(tmNode.Mempool(), cfg.Mempool.SortTxByGp, cfg.Mempool.EnablePendingPool)

	if cfg.Mempool.EnablePendingPool {
		cliCtx := context.NewCLIContext().WithCodec(cdc)
		cliCtx.Client = local.New(tmNode)
		cliCtx.TrustNode = true
		accRetriever := types.NewAccountRetriever(cliCtx)
		tmNode.Mempool().SetAccountRetriever(accRetriever)
	}

	if parser, ok := app.(mempool.TxInfoParser); ok {
		tmNode.Mempool().SetTxInfoParser(parser)
	}

	// run forever (the node will not be returned)
	select {}
}

// Use SetExternalPackageValue to set external package config value.
func SetExternalPackageValue(cmd *cobra.Command) {
	iavl.IavlCacheSize = viper.GetInt(iavl.FlagIavlCacheSize)
	tmiavl.IavlCacheInitRatio = viper.GetFloat64(tmiavl.FlagIavlCacheInitRatio)
	tmiavl.OutputModules, _ = cmd.Flags().GetStringToInt(tmiavl.FlagOutputModules)
	tmiavl.CommitIntervalHeight = viper.GetInt64(tmiavl.FlagIavlCommitIntervalHeight)
	tmiavl.MinCommitItemCount = viper.GetInt64(tmiavl.FlagIavlMinCommitItemCount)
	tmiavl.HeightOrphansCacheSize = viper.GetInt(tmiavl.FlagIavlHeightOrphansCacheSize)
	tmiavl.MaxCommittedHeightNum = viper.GetInt(tmiavl.FlagIavlMaxCommittedHeightNum)
	tmiavl.EnableAsyncCommit = viper.GetBool(tmiavl.FlagIavlEnableAsyncCommit)
	system.EnableGid = viper.GetBool(system.FlagEnableGid)

	state.ApplyBlockPprofTime = viper.GetInt(state.FlagApplyBlockPprofTime)
	state.HomeDir = viper.GetString(cli.HomeFlag)

	abci.SetDisableABCIQueryMutex(viper.GetBool(abci.FlagDisableABCIQueryMutex))
	abci.SetDisableCheckTx(viper.GetBool(abci.FlagDisableCheckTx))

	tmtypes.DownloadDelta = viper.GetBool(tmtypes.FlagDownloadDDS)
	tmtypes.UploadDelta = viper.GetBool(tmtypes.FlagUploadDDS)
	tmtypes.FastQuery = viper.GetBool(tmtypes.FlagFastQuery)
	tmtypes.DeltaVersion = viper.GetInt(tmtypes.FlagDeltaVersion)
}
