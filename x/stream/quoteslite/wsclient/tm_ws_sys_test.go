package wsclient

import (
	"context"
	"fmt"
	"github.com/okex/okchain/x/common"
	types2 "github.com/tendermint/tendermint/types"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/rpc/client"
	"testing"
	"time"
)

func TestWebsocketClient(t *testing.T) {
	common.SkipSysTestChecker(t)

	httpCli := client.NewHTTP("tcp://localhost:8888", "/websocket")
	httpCli.Start()

	tmEvtData, err := client.WaitForOneEvent(httpCli.WSEvents, types2.EventNewBlock, time.Minute*10)

	require.True(t, err == nil, err)
	require.True(t, tmEvtData != nil)
	fmt.Printf("%+v", tmEvtData)
}

func TestWebsocketClient2(t *testing.T) {

	common.SkipSysTestChecker(t)

	c := client.NewHTTP("tcp://localhost:26657", "/websocket")
	c.Start()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	subscriber := "TestWebsocketClient2"
	//eventCh, err := c.Subscribe(ctx, subscriber, "backend.channel='ticker' AND backend.product='ETH-USDT'")
	eventCh, err := c.Subscribe(ctx, subscriber, "backend.channel='dex_spot/ticker:tbtc_tusdk'")
	require.Nil(t, err, err)

	// make sure to unregister after the test is over
	defer c.UnsubscribeAll(ctx, subscriber)

	for {
		select {
		case event := <-eventCh:
			fmt.Printf("%+v\n", event.Events["backend.data"])
		case <-ctx.Done():
			fmt.Println("timed out waiting for event")
			return
		}
	}
}
