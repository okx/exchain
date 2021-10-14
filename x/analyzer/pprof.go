package analyzer

import (
	"errors"
	"fmt"
	"os"
	"path"
	"runtime/pprof"
	"time"
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

func InitializePprofDumper(dumpPath string, coolDownStr string, abciElapsed int64) {
	if singlePprofDumper != nil {
		return
	}
	coolDown, err := time.ParseDuration(coolDownStr)
	if err != nil {
		coolDown = defaultCoolDown
	}
	singlePprofDumper = &pprofDumper{
		dumpPath:           dumpPath,
		coolDown:           coolDown,
		cpuCoolDownTime:    time.Now(),
		triggerAbciElapsed: abciElapsed,
	}
}

type pprofDumper struct {
	dumpPath string
	// the cool down time after every type of dump
	coolDown           time.Duration
	cpuCoolDownTime    time.Time
	triggerAbciElapsed int64
}

func (dumper *pprofDumper) cpuProfile(height int64) error {
	if dumper.cpuCoolDownTime.After(time.Now()) {
		return errors.New("cpu dump is in coolDown")
	}
	fileName := dumper.getBinaryFileName(height, cpu)
	bf, err := os.OpenFile(fileName, defaultLoggerFlags, defaultLoggerPerm)
	if err != nil {
		return err
	}
	defer bf.Close()

	err = pprof.StartCPUProfile(bf)
	if err != nil {
		return err
	}

	time.Sleep(defaultCPUSamplingTime)
	pprof.StopCPUProfile()
	dumper.cpuCoolDownTime = time.Now().Add(dumper.coolDown)
	return nil
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
