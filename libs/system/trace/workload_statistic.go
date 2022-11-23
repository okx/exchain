package trace

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/okex/exchain/libs/tendermint/libs/log"
)

var (
	applyBlockWorkloadStatistic = newWorkloadStatistic(
		[]time.Duration{time.Hour, 2 * time.Hour, 4 * time.Hour, 8 * time.Hour}, []string{LastRun, Persist})
)

// TODO: think about a very long work which longer than a statistic period.

// NOTE: CAN NOT be used concurrently for those reasons:
// read/write almost all fields
// workload summary may be wrong if a work is still running(latestBegin.IsZero isn't true)
//
// calling sequence:
// 1. newWorkloadStatistic
// 2. calling begin/end before and after doing some work
// 3. calling summary to get a summary statistic
type WorkloadStatistic struct {
	tags      map[string]struct{}
	workloads []workloadSummary

	latestTag   string
	latestBegin time.Time
	logger      log.Logger

	workCh chan workInfo
}

type workloadSummary struct {
	period   time.Duration
	workload int64
}

type workInfo struct {
	end      time.Time
	workload int64
}

func GetApplyBlockWorkloadSttistic() *WorkloadStatistic {
	return applyBlockWorkloadStatistic
}

func newWorkloadStatistic(periods []time.Duration, tags_ []string) *WorkloadStatistic {
	tags := toTagsMap(tags_)

	workloads := make([]workloadSummary, 0, len(periods))
	for _, period := range periods {
		workloads = append(workloads, workloadSummary{period, 0})
	}

	wls := &WorkloadStatistic{tags: tags, workloads: workloads, workCh: make(chan workInfo, 1000)}
	go wls.shrink_loop()

	return wls
}

func (ws *WorkloadStatistic) SetLogger(logger log.Logger) {
	ws.logger = logger
}

func (ws *WorkloadStatistic) Add(tag string, wl time.Duration) {
	if _, ok := ws.tags[tag]; !ok {
		return
	}

	now := time.Now()
	for i := range ws.workloads {
		atomic.AddInt64(&ws.workloads[i].workload, int64(wl))
	}

	ws.workCh <- workInfo{now, int64(wl)}
}

func (ws *WorkloadStatistic) Format() string {
	var sumItem []string
	for _, summary := range ws.summary() {
		sumItem = append(sumItem, fmt.Sprintf("%.2f", float64(summary.workload)/float64(summary.period)))
	}

	return strings.Join(sumItem, "|")
}

func (ws *WorkloadStatistic) begin(tag string, t time.Time) {
	if _, ok := ws.tags[tag]; !ok {
		return
	}

	ws.latestTag = tag
	ws.latestBegin = t
}

func (ws *WorkloadStatistic) end(tag string, t time.Time) {
	if _, ok := ws.tags[tag]; !ok {
		return
	}

	if ws.latestTag != tag {
		ws.logger.Error("WorkloadStatistic", ": begin tag", ws.latestTag, "end tag", tag)
		return
	}
	if ws.latestBegin.IsZero() {
		ws.logger.Error("WorkloadStatistic", "begin is not called before end")
		return
	}

	dur := t.Sub(ws.latestBegin)
	for i := range ws.workloads {
		atomic.AddInt64(&ws.workloads[i].workload, int64(dur))
	}

	ws.workCh <- workInfo{t, int64(dur)}
	ws.latestBegin = time.Time{}
}

type summaryInfo struct {
	period   time.Duration
	workload time.Duration
}

func (ws *WorkloadStatistic) summary() []summaryInfo {
	if !ws.latestBegin.IsZero() {
		ws.logger.Error("WorkloadStatistic", ": some work is still running when calling summary")
		return nil
	}

	summary := make([]summaryInfo, 0, len(ws.workloads))
	for _, wload := range ws.workloads {
		summary = append(summary, summaryInfo{wload.period, time.Duration(atomic.LoadInt64(&wload.workload))})
	}
	return summary
}

func (ws *WorkloadStatistic) shrink_loop() {
	shrinkInfos := make([]map[int64]int64, 0, len(ws.workloads))
	for i := 0; i < len(ws.workloads); i++ {
		shrinkInfos = append(shrinkInfos, make(map[int64]int64))
	}

	var latest int64
	ticker := time.NewTicker(time.Second)

	for {
		select {
		case work := <-ws.workCh:
			var earilest int64 = int64(^uint64(0) >> 1)

			for i, wload := range ws.workloads {
				period := wload.period
				expiredSec := work.end.Add(period).Unix()
				if expiredSec < earilest {
					earilest = expiredSec
				}

				info := shrinkInfos[i]
				// TODO: it makes recoding workload larger than actual value
				// if a work begin before this period and end during this period
				if _, ok := info[expiredSec]; !ok {
					info[expiredSec] = work.workload
				} else {
					info[expiredSec] += work.workload
				}
			}

			if latest == 0 {
				latest = earilest
			}
		case t := <-ticker.C:
			current := t.Unix()
			if latest == 0 {
				latest = current
			}

			for index, info := range shrinkInfos {
				for i := latest; i < current+1; i++ {
					w, ok := info[i]
					if ok {
						atomic.AddInt64(&ws.workloads[index].workload, -w)
						delete(info, i)
					}
				}
			}

			latest = current
		}
	}

}

func toTagsMap(keys []string) map[string]struct{} {
	tags := make(map[string]struct{})
	for _, tag := range keys {
		tags[tag] = struct{}{}
	}
	return tags
}
