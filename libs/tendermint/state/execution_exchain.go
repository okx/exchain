package state

import (
	"log"
	"os"
	"path"
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
	HomeDir = ""
	tmpPprofName = "tmp.cpu.pprof"
)

func SetIgnoreSmbCheck(check bool) {
	IgnoreSmbCheck = check
}

func PprofStart() (*os.File, time.Time) {
	startTime := time.Now()
	p := getFilePath("")
	if _, err := os.Stat(p); os.IsNotExist(err) {
		err := os.MkdirAll(p, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	f, err := os.OpenFile(path.Join(p, tmpPprofName), os.O_RDWR|os.O_CREATE, 0644)
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
		newDir := getFilePath(pprofName)
		os.Rename(getFilePath(tmpPprofName), newDir)
	} else {
		os.Remove(getFilePath(tmpPprofName))
	}

}

func getFilePath(fileName string) string {
	return path.Join(HomeDir, "pprof", fileName)
}