package quoteslite

import (
	"context"
	"encoding/json"
	"fmt"
	okex "github.com/okex/okchain/cmd/quoteslite/okwebsocket"
	"net/http"
	"time"

	//"log"

	amino "github.com/tendermint/go-amino"

	"github.com/gorilla/websocket"
	"github.com/tendermint/tendermint/libs/log"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
	"github.com/tendermint/tendermint/rpc/core"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"
)

const (
	wsEndpoint = "/websocket"
)

var (
	upgrader = websocket.Upgrader{}
)

// StartProxy will start the websocket manager on the client,
// set up the rpc routes to proxy via the given client,
// and start up an http/rpc server on the location given by bind (eg. :1234)
// NOTE: This function blocks - you may want to call it in a go-routine.
func StartProxy(c rpcclient.Client, listenAddr string, logger log.Logger, maxOpenConnections int) error {
	err := c.Start()
	if err != nil {
		return err
	}

	cdc := amino.NewCodec()
	ctypes.RegisterAmino(cdc)
	r := RPCRoutes(c)

	// build the handler...
	mux := http.NewServeMux()
	rpcserver.RegisterRPCFuncs(mux, r, cdc, logger)

	unsubscribeFromAllEvents := func(remoteAddr string) {
		if err := c.UnsubscribeAll(context.Background(), remoteAddr); err != nil {
			logger.Error("Failed to unsubscribe from events", "err", err)
		}
	}
	wm := rpcserver.NewWebsocketManager(r, cdc, rpcserver.OnDisconnect(unsubscribeFromAllEvents))
	wm.SetLogger(logger)
	core.SetLogger(logger)
	mux.HandleFunc(wsEndpoint, wm.WebsocketHandler)

	config := rpcserver.DefaultConfig()
	config.MaxOpenConnections = maxOpenConnections
	l, err := rpcserver.Listen(listenAddr, config)
	if err != nil {
		return err
	}
	return rpcserver.StartHTTPServer(l, mux, logger, config)
}

// RPCRoutes just routes everything to the given client, as if it were
// a tendermint fullnode.
//
// if we want security, the client must implement it as a secure client
func RPCRoutes(c rpcclient.Client) map[string]*rpcserver.RPCFunc {
	return map[string]*rpcserver.RPCFunc{
		// Subscribe/unsubscribe are reserved for websocket events.
		//"subscribe":       rpcserver.NewWSRPCFunc(c.(Wrapper).SubscribeWS, "query"),
		//"unsubscribe":     rpcserver.NewWSRPCFunc(c.(Wrapper).UnsubscribeWS, "query"),
		//"unsubscribe_all": rpcserver.NewWSRPCFunc(c.(Wrapper).UnsubscribeAllWS, ""),
	}
}


func msgHandler(w http.ResponseWriter, r *http.Request)  {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("upgrade:%+v\n", err)
		return
	}
	defer c.Close()

	//1. receive a subscribe message from
	mt, message, err := c.ReadMessage()
	if err != nil && mt != websocket.TextMessage {
		fmt.Printf("read:+v\n", err)
	}

	op := okex.BaseOp{}
	json.Unmarshal(message, &op)


	// 2. connect to okchaind websocket rpc so as to get the notification of specific channel
	rpcCli := rpcclient.NewHTTP("tcp://localhost:26657", "/websocket")
	err = rpcCli.Start()
	defer rpcCli.Stop()


	// 3. send subscribe response
	if err == nil {
		rep := okex.WSEventResponse{}
		rep.Event = op.Op
		rep.Channel = op.Args[0]
		writeBytes, _ := json.Marshal(rep)
		c.WriteMessage(websocket.TextMessage, writeBytes)
	}


	// 4. receive filter event from abci.event websocket rpc
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	subTopic := okex.FormSubscriptionTopic(op.Args[0])
	subscriber := c.RemoteAddr().String()
	query := fmt.Sprintf("backend.channel='%s' AND backend.filter='%s'", subTopic.Channel, subTopic.Filter)
	eventCh, err := rpcCli.Subscribe(ctx, subscriber, query)


	if err == nil {
		// 4.1 receive events from rpc websocket and write the message back to opendex websocket client.
		//go func() {
			for {
				select {
				case event := <-eventCh:
					data := event.Events["backend.data"]
					tr := okex.WSTableResponse{Table: subTopic.Channel, Data: []interface{}{data}}
					c.WriteJSON(tr)

				case <-ctx.Done():
					fmt.Printf("subscriber: %s timed out waiting for event\n", subscriber)
					return
				}
			}
		//}()
	}
}

func StartWSServer(logger log.Logger) {
	http.HandleFunc("/ws/v3", msgHandler)
	logger.Info("Starting WebSocket server on 127.0.0.1:6666")
	http.ListenAndServe("localhost:6666", nil)
}