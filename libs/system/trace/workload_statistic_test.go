package trace

import (
	"testing"
	"time"
)

func TestWorkload(t *testing.T) {
	abciWorkload := time.Second
	lastRunWorkload := 2 * time.Minute
	persistWorkload := time.Second
	expectWorkload := int64((lastRunWorkload + persistWorkload).Seconds())

	trc := NewTracer(ApplyBlock)
	trc.EnableSummary()
	trc.SetWorkloadStatistic(GetApplyBlockWorkloadSttistic())

	defer func() {
		GetElapsedInfo().AddInfo(RunTx, trc.Format())

		summary := GetApplyBlockWorkloadSttistic().summary()
		for _, sum := range summary {
			workload := int64(sum.workload.Seconds())
			if workload != expectWorkload {
				t.Errorf("period %s: expect workload %v but got %v\n", sum.period.peirod, expectWorkload, workload)
			}
		}
	}()

	trc.Pin(Abci)
	time.Sleep(abciWorkload)
	GetApplyBlockWorkloadSttistic().Add(LastRun, 2*time.Minute)

	trc.Pin(Persist)
	time.Sleep(time.Second)

}
