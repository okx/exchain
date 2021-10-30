package config

import (
	"path"

	"github.com/okex/exchain/dependence/cosmos-sdk/server"

	"github.com/okex/exchain/x/analyzer"

	"github.com/mosn/holmes"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
)

type PporfConfig struct {
	autoDump              bool
	collectInterval       string
	coolDown              string
	dumpPath              string
	cpuTriggerPercentMin  int
	cpuTriggerPercentDiff int
	cpuTriggerPercentAbs  int
	memTriggerPercentMin  int
	memTriggerPercentDiff int
	memTriggerPercentAbs  int
	triggerAbciElapsed    int64
	useCGroup             bool
}

const (
	FlagPprofAutoDump              = "pprof-auto-dump"
	FlagPprofCollectInterval       = "pprof-collect-interval"
	FlagPprofCpuTriggerPercentMin  = "pprof-cpu-trigger-percent-min"
	FlagPprofCpuTriggerPercentDiff = "pprof-cpu-trigger-percent-diff"
	FlagPprofCpuTriggerPercentAbs  = "pprof-cpu-trigger-percent-abs"
	FlagPprofMemTriggerPercentMin  = "pprof-mem-trigger-percent-min"
	FlagPprofMemTriggerPercentDiff = "pprof-mem-trigger-percent-diff"
	FlagPprofMemTriggerPercentAbs  = "pprof-mem-trigger-percent-abs"
	FlagPprofCoolDown              = "pprof-cool-down"
	FlagPprofAbciElapsed           = "pprof-trigger-abci-elapsed"
	FlagPprofUseCGroup             = "pprof-use-cgroup"
)

// PprofDownload auto dump pprof
func PprofDownload(context *server.Context) {
	c := LoadPprofFromConfig()
	if !c.autoDump {
		return
	}

	// auto download pprof by analyzer
	analyzer.InitializePprofDumper(context.Logger, c.dumpPath, c.coolDown, c.triggerAbciElapsed)

	// auto download pprof by holmes
	h, err := holmes.New(
		holmes.WithCollectInterval(c.collectInterval),
		holmes.WithCoolDown(c.coolDown),
		holmes.WithDumpPath(c.dumpPath),
		holmes.WithCPUDump(c.cpuTriggerPercentMin, c.cpuTriggerPercentDiff, c.cpuTriggerPercentAbs),
		holmes.WithMemDump(c.memTriggerPercentMin, c.memTriggerPercentDiff, c.memTriggerPercentAbs),
		holmes.WithGoroutineDump(2000, 50, 5000),
		holmes.WithBinaryDump(),
		holmes.WithCGroup(c.useCGroup),
	)
	if err != nil {
		tmos.Exit(err.Error())
	}
	h.EnableCPUDump().EnableMemDump().EnableGoroutineDump()
	// start the metrics collect and dump loop
	h.Start()
}

func LoadPprofFromConfig() *PporfConfig {
	autoDump := viper.GetBool(FlagPprofAutoDump)
	collectInterval := viper.GetString(FlagPprofCollectInterval)
	dumpPath := path.Join(viper.GetString(cli.HomeFlag), "pprof")
	cpuTriggerPercentMin := viper.GetInt(FlagPprofCpuTriggerPercentMin)
	cpuTriggerPercentDiff := viper.GetInt(FlagPprofCpuTriggerPercentDiff)
	cpuTriggerPercentAbs := viper.GetInt(FlagPprofCpuTriggerPercentAbs)
	memTriggerPercentMin := viper.GetInt(FlagPprofMemTriggerPercentMin)
	memTriggerPercentDiff := viper.GetInt(FlagPprofMemTriggerPercentDiff)
	memTriggerPercentAbs := viper.GetInt(FlagPprofMemTriggerPercentAbs)
	coolDown := viper.GetString(FlagPprofCoolDown)
	triggerAbciElapsed := viper.GetInt64(FlagPprofAbciElapsed)
	useCGroup := viper.GetBool(FlagPprofUseCGroup)
	c := &PporfConfig{
		autoDump:              autoDump,
		collectInterval:       collectInterval,
		coolDown:              coolDown,
		dumpPath:              dumpPath,
		cpuTriggerPercentMin:  cpuTriggerPercentMin,
		cpuTriggerPercentDiff: cpuTriggerPercentDiff,
		cpuTriggerPercentAbs:  cpuTriggerPercentAbs,
		memTriggerPercentMin:  memTriggerPercentMin,
		memTriggerPercentDiff: memTriggerPercentDiff,
		memTriggerPercentAbs:  memTriggerPercentAbs,
		triggerAbciElapsed:    triggerAbciElapsed,
		useCGroup:             useCGroup,
	}
	return c
}
