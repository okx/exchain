package stream

import (
	"fmt"
	"sync"
	"testing"
	"time"

	appCfg "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/okex/okchain/x/stream/types"

	"github.com/okex/okchain/x/stream/common"

	"github.com/okex/okchain/x/order"
	"github.com/okex/okchain/x/token"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestNewStream(t *testing.T) {
	streamConfig := &appCfg.StreamConfig{
		Engine:         "analysis|mysql|" + MYSQLURL + ",notify|redis|" + REDISURL + ",kline|pulsar|" + PULSARURL,
		RedisScheduler: REDISURL,
		RedisLock:      REDISURL,
		WorkerId:       "worker0",
	}

	mockApp, addrKeysSlice := GetMockApp(t, 2, streamConfig)
	pool, err := common.NewPool(REDISURL, "", mockApp.Logger())
	require.Nil(t, err)
	pool.Get().Do("FLUSHALL")

	require.Equal(t, 3, len(mockApp.streamKeeper.stream.engines))
	mockApp.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 2}})
	ctx := mockApp.NewContext(false, abci.Header{})

	tokenPair := token.GetBuiltInTokenPair()
	mockApp.TokenKeeper.SaveTokenPair(ctx, tokenPair)

	quantity := "1"
	price := "0.1"

	orderMsg0 := order.NewMsgNewOrder(nil, types.TestTokenPair, types.BuyOrder, price, quantity)
	ctx = mockApp.NewContext(true, abci.Header{Height: 2})
	MockApplyBlock(mockApp, int64(2), ProduceOrderTxs(mockApp, ctx, 10, addrKeysSlice[0], &orderMsg0))

	//ctx = mockApp.NewContext(true, abci.Header{Height: 4})
	//orderMsg1 := order.NewMsgNewOrder(nil, types.TestTokenPair, types.SellOrder, price, quantity)
	//MockApplyBlock(mockApp, int64(4), ProduceOrderTxs(mockApp, ctx, 10, addrKeysSlice[1], &orderMsg1))
	//
	//ctx = mockApp.NewContext(true, abci.Header{Height: 5})
}

func doASimpleIO(i int, wg *sync.WaitGroup) {
	fmt.Printf("%d", i)
	wg.Done()
}

func BenchmarkCreateGoRoutine(b *testing.B) {

	wg := &sync.WaitGroup{}

	b.N = 1000000
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

		STREAM_TASK_PHRASE1_NEXT_ACTION_RESTART,
		STREAM_TASK_PHRASE1_NEXT_ACTION_JUMP_NEXT_BLK,
		STREAM_TASK_PHRASE1_NEXT_ACTION_NEW_TASK,
		STREAM_TASK_PHRASE1_NEXT_ACTION_RERUN_TASK,

		STREAM_TASK_PHRASE2_NEXT_ACTION_RESTART,
		STREAM_TASK_PHRASE2_NEXT_ACTION_JUMP_NEXT_BLK,
	}

	for _, action := range actions {
		fmt.Printf("%d:%s\n", action, StreamConstDesc[action])
		require.True(t, action != 0, action)
	}

	notExist := StreamKind2EngineKindMap[100]
	fmt.Println(notExist)
}

func TestStreamTask(t *testing.T) {
	st := Task{}
	st.Height = 100
	st.DoneMap = map[StreamKind]bool{
		StreamRedisKind:  false,
		StreamPulsarKind: false,
	}
	st.UpdatedAt = time.Now().Unix()

	r := st.toJsonStr()
	fmt.Println(r)

	st2, err := parseTaskFromJsonStr(r)
	require.True(t, err == nil, err)
	require.True(t, st.Height == st2.Height && st.UpdatedAt == st2.UpdatedAt, st, st2)
	require.EqualValues(t, st, *st2)

	status := st.GetStatus()
	require.True(t, status == STREAM_TASK_STATUS_FAIL, status)

	st.DoneMap[StreamRedisKind] = true
	status = st.GetStatus()
	require.True(t, status == STREAM_TASK_STATUS_PARTITIAL_SUCCESS, status)

	st.DoneMap[StreamPulsarKind] = true
	status = st.GetStatus()
	require.True(t, status == STREAM_TASK_STATUS_SUCCESS, status)
}
