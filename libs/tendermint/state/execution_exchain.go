package state

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
)

const (
	FlagApplyBlockPprof = "applyblock-pprof"
	FlagApplyBlockPprofTime = "applyblock-pprof-time"
)

var (
	IgnoreSmbCheck bool = false

	ApplyBlockPprof bool = false
	ApplyBlockPprofTime  = 8000
	tmpPprofDir = "tmp.cpu.pprof"
)

func SetIgnoreSmbCheck(check bool) {
	IgnoreSmbCheck = check
}

func PprofStart() (*os.File, time.Time) {
	startTime := time.Now()
	f, err := os.OpenFile(tmpPprofDir, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	runtime.SetCPUProfileRate(300)
	pprof.StartCPUProfile(f)
	return f, startTime
}

func PprofEnd(height int, f *os.File, startTime time.Time) {
	pprof.StopCPUProfile()
	f.Close()
	sec := time.Since(startTime).Milliseconds()
	if int(sec) >= ApplyBlockPprofTime {
		pprofName := "pprof." + strconv.Itoa(height) + "-" + strconv.Itoa(int(sec)) + "ms.bin"
		os.Rename(tmpPprofDir, pprofName)
	} else {
		os.Remove(tmpPprofDir)
	}

}