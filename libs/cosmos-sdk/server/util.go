package server

import (
	"encoding/json"
	"errors"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/google/gops/agent"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/client/lcd"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/libs/cosmos-sdk/server/config"
	"github.com/okex/exchain/libs/cosmos-sdk/version"
	tcmd "github.com/okex/exchain/libs/tendermint/cmd/tendermint/commands"
	cfg "github.com/okex/exchain/libs/tendermint/config"
	"github.com/okex/exchain/libs/tendermint/libs/cli"
	tmflags "github.com/okex/exchain/libs/tendermint/libs/cli/flags"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/libs/tendermint/state"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const FlagGops = "gops"

// server context
type Context struct {
	Config *cfg.Config
	Logger log.Logger
}

func NewDefaultContext() *Context {
	return NewContext(
		cfg.DefaultConfig(),
		log.NewTMLogger(log.NewSyncWriter(os.Stdout)),
	)
}

func NewContext(config *cfg.Config, logger log.Logger) *Context {
	return &Context{config, logger}
}

//___________________________________________________________________________________

// PersistentPreRunEFn returns a PersistentPreRunE function for cobra
// that initailizes the passed in context with a properly configured
// logger and config object.
func PersistentPreRunEFn(context *Context) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if cmd.Name() == version.Cmd.Name() {
			return nil
		}
		config, err := interceptLoadConfig()
		if err != nil {
			return err
		}
		if !viper.IsSet(state.FlagDeliverTxsExecMode) {
			if viper.GetBool(state.FlagEnableConcurrency) {
				viper.Set(state.FlagDeliverTxsExecMode, state.DeliverTxsExecModeParallel)
			}
		}
		// okchain
		output := os.Stdout
		if !config.LogStdout {
			output, err = os.OpenFile(config.LogFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
			if err != nil {
				return err
			}
		}

		logger := log.NewTMLogger(log.NewSyncWriter(output))
		logger, err = tmflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel())
		if err != nil {
			return err
		}
		if viper.GetBool(cli.TraceFlag) {
			logger = log.NewTracingLogger(logger)
		}
		logger = logger.With("module", "main")
		context.Config = config
		context.Logger = logger

		if viper.GetBool(FlagGops) {
			err = agent.Listen(agent.Options{ShutdownCleanup: true})
			if err != nil {
				logger.Error("gops agent error", "err", err)
			}
		}

		return nil
	}
}

// If a new config is created, change some of the default tendermint settings
func interceptLoadConfig() (conf *cfg.Config, err error) {
	tmpConf := cfg.DefaultConfig()
	err = viper.Unmarshal(tmpConf)
	if err != nil {
		// TODO: Handle with #870
		panic(err)
	}
	rootDir := tmpConf.RootDir
	configFilePath := filepath.Join(rootDir, "config/config.toml")
	// Intercept only if the file doesn't already exist

	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// the following parse config is needed to create directories
		conf, _ = tcmd.ParseConfig() // NOTE: ParseConfig() creates dir/files as necessary.
		conf.ProfListenAddress = "localhost:6060"
		conf.P2P.RecvRate = 5120000
		conf.P2P.SendRate = 5120000
		conf.TxIndex.IndexAllKeys = true
		conf.Consensus.TimeoutCommit = 3 * time.Second
		conf.Consensus.TimeoutConsensus = 1 * time.Second
		cfg.WriteConfigFile(configFilePath, conf)
		// Fall through, just so that its parsed into memory.
	}

	if conf == nil {
		conf, err = tcmd.ParseConfig() // NOTE: ParseConfig() creates dir/files as necessary.
		if err != nil {
			panic(err)
		}
	}

	appConfigFilePath := filepath.Join(rootDir, "config/exchaind.toml")
	if _, err := os.Stat(appConfigFilePath); os.IsNotExist(err) {
		appConf, _ := config.ParseConfig()
		config.WriteConfigFile(appConfigFilePath, appConf)
	}

	viper.SetConfigName("exchaind")
	err = viper.MergeInConfig()

	return conf, err
}

// add server commands
func AddCommands(
	ctx *Context, cdc *codec.CodecProxy,
	registry jsonpb.AnyResolver,
	rootCmd *cobra.Command,
	appCreator AppCreator, appStop AppStop, appExport AppExporter,
	registerRouters func(rs *lcd.RestServer),
	registerAppFlagFn func(cmd *cobra.Command),
	appPreRun func(ctx *Context, cmd *cobra.Command) error,
	subFunc func(logger log.Logger) log.Subscriber) {

	rootCmd.PersistentFlags().String("log_level", ctx.Config.LogLevel, "Log level")
	rootCmd.PersistentFlags().String("log_file", ctx.Config.LogFile, "Log file")
	rootCmd.PersistentFlags().Bool("log_stdout", ctx.Config.LogStdout, "Print log to stdout, rather than a file")

	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint subcommands",
	}

	tendermintCmd.AddCommand(
		ShowNodeIDCmd(ctx),
		ShowValidatorCmd(ctx),
		ShowAddressCmd(ctx),
		VersionCmd(ctx),
	)

	rootCmd.AddCommand(
		StartCmd(ctx, cdc, registry, appCreator, appStop, registerRouters, registerAppFlagFn, appPreRun, subFunc),
		StopCmd(ctx),
		UnsafeResetAllCmd(ctx),
		flags.LineBreak,
		tendermintCmd,
		ExportCmd(ctx, cdc.GetCdc(), appExport),
		flags.LineBreak,
		version.Cmd,
	)
}

//___________________________________________________________________________________

// InsertKeyJSON inserts a new JSON field/key with a given value to an existing
// JSON message. An error is returned if any serialization operation fails.
//
// NOTE: The ordering of the keys returned as the resulting JSON message is
// non-deterministic, so the client should not rely on key ordering.
func InsertKeyJSON(cdc *codec.Codec, baseJSON []byte, key string, value json.RawMessage) ([]byte, error) {
	var jsonMap map[string]json.RawMessage

	if err := cdc.UnmarshalJSON(baseJSON, &jsonMap); err != nil {
		return nil, err
	}

	jsonMap[key] = value
	bz, err := codec.MarshalJSONIndent(cdc, jsonMap)

	return json.RawMessage(bz), err
}

// https://stackoverflow.com/questions/23558425/how-do-i-get-the-local-ip-address-in-go
// TODO there must be a better way to get external IP
func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if skipInterface(iface) {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip := addrToIP(addr)
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}

// TrapSignal traps SIGINT and SIGTERM and terminates the server correctly.
func TrapSignal(cleanupFunc func(int)) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs

		exitCode := 128
		switch sig {
		case syscall.SIGINT:
			exitCode += int(syscall.SIGINT)
		case syscall.SIGTERM:
			exitCode += int(syscall.SIGTERM)
		}

		if cleanupFunc != nil {
			cleanupFunc(exitCode)
		}

		os.Exit(exitCode)
	}()
}

func skipInterface(iface net.Interface) bool {
	if iface.Flags&net.FlagUp == 0 {
		return true // interface down
	}
	if iface.Flags&net.FlagLoopback != 0 {
		return true // loopback interface
	}
	return false
}

func addrToIP(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	return ip
}

// DONTCOVER
