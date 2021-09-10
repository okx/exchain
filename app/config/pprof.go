package config

import (
	"github.com/mosn/holmes"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	tmos "github.com/tendermint/tendermint/libs/os"
	"path"
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
}

const (
	FlagPprofAutoDump              = "pprof-auto-dump"
	FlagPprofCpuTriggerPercentMin  = "pprof-cpu-trigger-percent-min"
	FlagPprofCpuTriggerPercentDiff = "pprof-cpu-trigger-percent-diff"
	FlagPprofCpuTriggerPercentAbs  = "pprof-cpu-trigger-percent-abs"
	FlagPprofMemTriggerPercentMin  = "pprof-mem-trigger-percent-min"
	FlagPprofMemTriggerPercentDiff = "pprof-mem-trigger-percent-diff"
	FlagPprofMemTriggerPercentAbs  = "pprof-mem-trigger-percent-abs"
)

// PprofDown auto dump pprof
func PprofDown() {
	c := LoadPprofFromConfig()
	if c.autoDump {
		h, err := holmes.New(
			holmes.WithCollectInterval(c.collectInterval),
			holmes.WithCoolDown(c.coolDown),
			holmes.WithDumpPath(c.dumpPath),
			holmes.WithCPUDump(c.cpuTriggerPercentMin, c.cpuTriggerPercentDiff, c.cpuTriggerPercentAbs),
			holmes.WithMemDump(c.memTriggerPercentMin, c.memTriggerPercentDiff, c.memTriggerPercentAbs),
			holmes.WithGoroutineDump(2000, 50, 5000),
			holmes.WithBinaryDump(),
		)
		if err != nil {
			tmos.Exit(err.Error())
		}
		h.EnableCPUDump().EnableMemDump().EnableGoroutineDump()
		// start the metrics collect and dump loop
		h.Start()
	}
}

func LoadPprofFromConfig() *PporfConfig {
	autoDump := viper.GetBool(FlagPprofAutoDump)
	dumpPath := path.Join(viper.GetString(cli.HomeFlag), "pprof")
	cpuTriggerPercentMin := viper.GetInt(FlagPprofCpuTriggerPercentMin)
	cpuTriggerPercentDiff := viper.GetInt(FlagPprofCpuTriggerPercentDiff)
	cpuTriggerPercentAbs := viper.GetInt(FlagPprofCpuTriggerPercentAbs)
	memTriggerPercentMin := viper.GetInt(FlagPprofMemTriggerPercentMin)
	memTriggerPercentDiff := viper.GetInt(FlagPprofMemTriggerPercentDiff)
	memTriggerPercentAbs := viper.GetInt(FlagPprofMemTriggerPercentAbs)

	c := &PporfConfig{
		autoDump:              autoDump,
		collectInterval:       "5s",
		coolDown:              "3m",
		dumpPath:              dumpPath,
		cpuTriggerPercentMin:  cpuTriggerPercentMin,
		cpuTriggerPercentDiff: cpuTriggerPercentDiff,
		cpuTriggerPercentAbs:  cpuTriggerPercentAbs,
		memTriggerPercentMin:  memTriggerPercentMin,
		memTriggerPercentDiff: memTriggerPercentDiff,
		memTriggerPercentAbs:  memTriggerPercentAbs,
	}
	return c
}
