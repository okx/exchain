package stream

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func doASimpleIO(i int, wg *sync.WaitGroup) {
	fmt.Printf("%d", i)
	wg.Done()
}

func BenchmarkCreateGoRoutine(b *testing.B) {

	wg := &sync.WaitGroup{}

	wg.Add(b.N)

	for i := 0; i < b.N; i++ {
		go doASimpleIO(i, wg)
	}

	wg.Wait()
}

func TestWaitGroup(t *testing.T) {

	af := func(ch chan struct{}) {
		defer func() {
			ch <- struct{}{}
		}()

		time.Sleep(time.Second)
		fmt.Printf("atom task done\n")
	}

	afch := make(chan struct{}, 2)
	go af(afch)
	go af(afch)

	timer := time.NewTimer(2 * time.Second)
	doneCnt := 0

	for {
		select {
		case <-timer.C:
			fmt.Printf("all atom task force stop becoz of timeout\n")
		case <-afch:
			doneCnt++
		}

		if doneCnt == 2 {
			fmt.Printf("all atom task done\n")
			break
		}
	}
	close(afch)

	time.Sleep(4 * time.Second)

}

func TestUtils(t *testing.T) {
	actions := []TaskConst{

		TaskPhase1NextActionRestart,
		TaskPhase1NextActionJumpNextBlock,
		TaskPhase1NextActionNewTask,
		TaskPhase1NextActionReturnTask,

		TaskPhase2NextActionRestart,
		TaskPhase2NextActionJumpNextBlock,
	}

	for _, action := range actions {
		fmt.Printf("%d:%s\n", action, TaskConstDesc[action])
		require.True(t, action != 0, action)
	}

	notExist := StreamKind2EngineKindMap[100]
	fmt.Println(notExist)
}

func TestStreamTask(t *testing.T) {
	st := Task{}
	st.Height = 100
	st.DoneMap = map[Kind]bool{
		StreamRedisKind:  false,
		StreamPulsarKind: false,
	}
	st.UpdatedAt = time.Now().Unix()

	r := st.toJSON()
	fmt.Println(r)

	st2, err := parseTaskFromJSON(r)
	require.True(t, err == nil, err)
	require.True(t, st.Height == st2.Height && st.UpdatedAt == st2.UpdatedAt, st, st2)
	require.EqualValues(t, st, *st2)

	status := st.GetStatus()
	require.True(t, status == TaskStatusStatusFail, status)

	st.DoneMap[StreamRedisKind] = true
	status = st.GetStatus()
	require.True(t, status == TaskStatusPartialSuccess, status)

	st.DoneMap[StreamPulsarKind] = true
	status = st.GetStatus()
	require.True(t, status == TaskStatusSuccess, status)
}
