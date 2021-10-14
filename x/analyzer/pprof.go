package analyzer

import (
	"fmt"
	"os"
	"path"
	"runtime/pprof"
	"time"

	"github.com/tendermint/tendermint/libs/log"
)

type configureType int

const (
	mem configureType = iota
	cpu
	thread
	goroutine
)

const (
	defaultLoggerFlags     = os.O_RDWR | os.O_CREATE | os.O_APPEND
	defaultLoggerPerm      = 0644
	defaultCPUSamplingTime = 30 * time.Second // collect 5s cpu profile
	defaultCoolDown        = 3 * time.Minute
)

var singlePprofDumper *pprofDumper

func InitializePprofDumper(logger log.Logger, dumpPath string, coolDownStr string, abciElapsed int64) {
	if singlePprofDumper != nil {
		return
	}
	coolDown, err := time.ParseDuration(coolDownStr)
	if err != nil {
		coolDown = defaultCoolDown
	}
	singlePprofDumper = &pprofDumper{
		logger:             logger.With("module", "main"),
		dumpPath:           dumpPath,
		coolDown:           coolDown,
		cpuCoolDownTime:    time.Now(),
		triggerAbciElapsed: abciElapsed,
	}
}

type pprofDumper struct {
	logger   log.Logger
	dumpPath string
	// the cool down time after every type of dump
	coolDown           time.Duration
	cpuCoolDownTime    time.Time
	triggerAbciElapsed int64
}

func (dumper *pprofDumper) cpuProfile(height int64) {
	if dumper.cpuCoolDownTime.After(time.Now()) {
		dumper.logger.Info(fmt.Sprintf("height(%d) cpu dump is in coolDown", height))
		return
	}
	dumper.cpuCoolDownTime = time.Now().Add(dumper.coolDown)
	go dumper.dumpCpuPprof(height)
}

func (dumper *pprofDumper) dumpCpuPprof(height int64) {
	fileName := dumper.getBinaryFileName(height, cpu)
	bf, err := os.OpenFile(fileName, defaultLoggerFlags, defaultLoggerPerm)
	if err != nil {
		dumper.logger.Error("height(%d) dump cpu pprof, open file(%s) error:%s", height, fileName, err.Error())
		return
	}
	defer bf.Close()

	err = pprof.StartCPUProfile(bf)
	if err != nil {
		dumper.logger.Error("height(%d) dump cpu pprof, StartCPUProfile error:%s", height, err.Error())
		return
	}

	time.Sleep(defaultCPUSamplingTime)
	pprof.StopCPUProfile()
	dumper.logger.Info("height(%d) dump cpu pprof file(%s)", fileName)
}

func (dumper *pprofDumper) getBinaryFileName(height int64, dumpType configureType) string {
	var (
		binarySuffix = time.Now().Format("20060102150405") + ".bin"
	)
	fileName := fmt.Sprintf("%s_%d_%s", type2name[dumpType], height, binarySuffix)
	return path.Join(dumper.dumpPath, fileName)
}

var type2name = map[configureType]string{
	mem:       "mem",
	cpu:       "cpu",
	thread:    "thread",
	goroutine: "goroutine",
}
