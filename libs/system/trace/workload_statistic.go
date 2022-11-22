package trace

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	applyBlockWorkloadStatistic = newWorkloadStatistic([]time.Duration{time.Hour, 2 * time.Hour, 4 * time.Hour, 8 * time.Hour}, []string{LastRun, Persist})
)

type workloadSummary struct {
	period   time.Duration
	workload atomic.Int64 // nano seconds
	//oldest   *list.Element // element in `works` list
}

// NOTE: CAN NOT be used concurrently for those reasons:
// read/write almost all fields
// workload summary may be wrong if a work is still running(latestBegin.IsZero isn't true)
//
// calling sequence:
// 1. newWorkloadStatistic
// 2. calling begin/end before and after doing some work
// 3. calling summary to get a summary statistic
type WorkloadStatistic struct {
	tags          map[string]struct{}
	maximumPeriod time.Duration
	workloads     []workloadSummary

	latestTag   string
	latestBegin time.Time

	workCh chan workInfo
}

type workInfo struct {
	t        time.Time // TODO: begin time or end time?
	workload int64
}

func GetApplyBlockWorkloadSttistic() *WorkloadStatistic {
	return applyBlockWorkloadStatistic
}

func newWorkloadStatistic(periods []time.Duration, tags_ []string) *WorkloadStatistic {
	tags := toTagsMap(tags_)

	var maximumPeriod time.Duration
	workloads := make([]workloadSummary, len(periods))
	for _, period := range periods {
		workloads = append(workloads, workloadSummary{period, atomic.Int64{}})
		if period > maximumPeriod {
			maximumPeriod = period
		}
	}

	return &WorkloadStatistic{tags: tags, maximumPeriod: maximumPeriod, workloads: workloads, workCh: make(chan workInfo)}
}

func (ws *WorkloadStatistic) begin(tag string, t time.Time) {
	if _, ok := ws.tags[tag]; !ok {
		return
	}

	ws.latestTag = tag
	ws.latestBegin = t
}

func (ws *WorkloadStatistic) end(tag string, t time.Time) {
	assert(ws.latestTag == tag, "WorkloadStatistic: begin tag is %s but end tag is %s", ws.latestTag, tag)
	assert(!ws.latestBegin.IsZero(), "WorkloadStatistic: begin is not called before end")

	dur := t.Sub(ws.latestBegin)
	for _, wload := range ws.workloads {
		wload.workload.Add(int64(dur))
	}

	ws.latestBegin = time.Time{}
}

type summaryInfo struct {
	period   time.Duration
	workload time.Duration
}

func (ws *WorkloadStatistic) summary() []summaryInfo {
	assert(ws.latestBegin.IsZero(), "WorkloadStatistic: some work is still running when calling summary")

	summary := make([]summaryInfo, len(ws.workloads))
	for _, wload := range ws.workloads {
		summary = append(summary, summaryInfo{wload.period, time.Duration(wload.workload.Load())})
	}
	return summary
}

func (ws *WorkloadStatistic) shrink_loop() {
	type mywork struct {
		t    time.Time
		load int64
	}

	shrinkInfos := make([]map[int64]int64, len(ws.workloads))
	for i := 0; i < len(ws.workloads); i++ {
		shrinkInfos = append(shrinkInfos, make(map[int64]int64))
	}
	var latest int64

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case work := <-ws.workCh:
			if latest == 0 {
				latest = work.t.Unix()
			}

			for i, wload := range ws.workloads {
				period := wload.period
				expired := work.t.Add(period)
				expiredSec := expired.Unix()

				info := shrinkInfos[i]
				if _, ok := info[expiredSec]; !ok {
					info[expiredSec] = work.workload
				} else {
					info[expiredSec] += work.workload
				}
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
						ws.workloads[index].workload.Add(-w)
					}
				}
			}
		}
	}

}

func toTagsMap(keys []string) map[string]struct{} {
	tags := make(map[string]struct{}, len(keys))
	for _, tag := range keys {
		tags[tag] = struct{}{}
	}
	return tags
}

func assert(b bool, msg string, args ...interface{}) {
	if !b {
		panic(fmt.Sprintf(msg, args...))
	}
}
