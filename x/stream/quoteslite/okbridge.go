package quoteslite

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/spf13/viper"

	"github.com/gorilla/websocket"
	okex "github.com/okex/okchain/x/stream/quoteslite/okwebsocket"
	"github.com/tendermint/tendermint/libs/log"
	rpccli "github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

type okWSContext struct {
	interruptedCh chan interface{}
	signalCh      chan os.Signal
}

func (ctx *okWSContext) closeAll() {
	close(ctx.signalCh)
	close(ctx.interruptedCh)
}

func newOKWSContext() *okWSContext {
	ctx := okWSContext{
		interruptedCh: make(chan interface{}, 4),
		signalCh:      make(chan os.Signal),
	}

	return &ctx
}

type okWSConn struct {
	cliConn      *websocket.Conn
	rpcConn      *rpccli.HTTP
	ctx          *okWSContext
	logger       log.Logger
	loginAddress string

	cliInChan    chan []byte
	cliOutChan   chan interface{}
	rpcEventChan chan ctypes.ResultEvent
	rpcStopChan  chan interface{}
}

func newOKWSConn(ctx *okWSContext, cliConn *websocket.Conn, logger log.Logger) *okWSConn {

	if ctx == nil || cliConn == nil {
		return nil
	}

	conn := okWSConn{
		cliConn:      cliConn,
		rpcConn:      nil,
		ctx:          ctx,
		logger:       logger,
		cliInChan:    make(chan []byte),
		cliOutChan:   make(chan interface{}),
		rpcEventChan: make(chan ctypes.ResultEvent, 64),
		rpcStopChan:  make(chan interface{}),
	}

	conn.start()

	return &conn
}

func (conn *okWSConn) start() {
	conn.logger.Debug("starting bi-direction communication")

	go conn.handleFinalise()
	go conn.handleCliRead()
	go conn.handleCliWrite()
	go conn.handleRPCEventReceived()
	go conn.handleConvert()
}

func (conn *okWSConn) stopAll() {
	conn.logger.Debug("okWSConn.stopAll start")

	// 1. close all the connection
	conn.logger.Debug("okWSConn.stopAll close bi-connection")
	if conn.cliConn != nil {
		conn.cliConn.Close()
		conn.logger.Debug("okWSConn'connection to client is closed.")
	}

	if conn.rpcConn != nil {
		conn.rpcConn.Stop()
		conn.logger.Debug("okWSConn'connection to rpc websocket is closed.")
	}

	// 2. close all the channel
	conn.logger.Debug("okWSConn.stopAll close connection channel")
	close(conn.rpcStopChan)
	close(conn.cliInChan)
	close(conn.cliOutChan)
	close(conn.rpcEventChan)

	conn.logger.Debug("okWSConn.stopAll close connection context channel")
	conn.ctx.closeAll()

	conn.logger.Debug("okWSConn.stopAll finished")
}

func (conn *okWSConn) handleFinalise() {
	conn.logger.Debug("handleFinalise start")

	select {
	case msg := <-conn.ctx.interruptedCh:
		conn.logger.Debug("handleFinalise get interrupted signal", "msg", msg)
		conn.stopAll()
	}

	conn.logger.Debug("handleFinalise finished")
}

func (conn *okWSConn) handleCliRead() {
	defer func() {
		if err := recover(); err != nil {
			conn.logger.Error(fmt.Sprintf("handleCliRead recover panic:%v", err))
			debug.PrintStack()
		}
	}()

	conn.logger.Debug("handleCliRead start")

	for {
		inMsgType, msg, err := conn.cliConn.ReadMessage()
		conn.logger.Debug("handleCliRead ReadMessage", "msg", string(msg), "msgType", inMsgType, "error", err)
		if err == nil {
			txtMsg := msg
			switch inMsgType {
			case websocket.TextMessage:
			case websocket.BinaryMessage:
				txtMsg, err = okex.GzipDecode(msg)
			default:
				continue
			}

			if err == nil {
				conn.cliInChan <- txtMsg
			}
		}

		if err != nil {
			conn.ctx.interruptedCh <- err
			break
		}
	}

	conn.logger.Debug("handleCliRead finished")
}

func (conn *okWSConn) handleCliWrite() {
	defer func() {
		if err := recover(); err != nil {
			conn.logger.Error(fmt.Sprintf("handleCliWrite recover panic:%v", err))
		}
	}()

	conn.logger.Debug("handleCliWrite start")
	for {
		select {
		case outMsg, ok := <-conn.cliOutChan:
			if !ok {
				break
			}

			var err error
			switch outMsg.(type) {
			case string:
				err = conn.cliConn.WriteMessage(websocket.TextMessage, []byte(outMsg.(string)))
			default:
				err = conn.cliConn.WriteJSON(outMsg)
			}

			conn.logger.Debug("handleCliWrite write", "OutMsg", outMsg)

			if err != nil {
				conn.ctx.interruptedCh <- err
				break
			}
		}
	}
	conn.logger.Debug("handleCliWrite finished")
}

func (conn *okWSConn) convert2WSTableResponseFromMap(resultEvt ctypes.ResultEvent, topic *okex.SubscriptionTopic) (r interface{}, e error) {

	//for k, v := range resultEvt.Events {
	//	conn.logger.Debug("verbose event items", "k", k, "v", v)
	//}
	innerChannel, e := topic.ToString()
	if e != nil {
		return nil, e
	}

	resp := okex.WSTableResponse{
		Table:  topic.Channel,
		Action: "update",
		Data:   nil,
	}

	matchFilterIdx := 0
	for idx, channel := range resultEvt.Events[rpcChannelKey] {
		if channel == innerChannel {
			matchFilterIdx = idx
			break
		}
	}

	eventItemStr := resultEvt.Events[rpcChannelDataKey][matchFilterIdx]
	var obj map[string]interface{}
	jerr := json.Unmarshal([]byte(eventItemStr), &obj)
	if jerr == nil {
		resp.Data = []interface{}{obj}
	} else {
		e = fmt.Errorf("error info: %s, eventItemStr: %s", jerr.Error(), eventItemStr)
	}

	return resp, e
}

func (conn *okWSConn) convertWSTableResponseFromList(resultEvt ctypes.ResultEvent, topic *okex.SubscriptionTopic) (r interface{}, e error) {
	innerChannel, e := topic.ToString()
	if e != nil {
		return nil, e
	}

	resp := okex.WSTableResponse{
		Table:  topic.Channel,
		Action: "update",
		Data:   nil,
	}

	matchFilterIdx := 0
	for idx, channel := range resultEvt.Events[rpcChannelKey] {
		if channel == innerChannel {
			matchFilterIdx = idx
			break
		}
	}

	eventItemStr := resultEvt.Events[rpcChannelDataKey][matchFilterIdx]
	var obj []interface{}
	e = json.Unmarshal([]byte(eventItemStr), &obj)
	if e == nil {
		resp.Data = obj
	}
	return resp, e
}

func (conn *okWSConn) handleRPCEventReceived() {
	defer func() {
		if err := recover(); err != nil {
			conn.logger.Error(fmt.Sprintf("handleRPCEventReceived recover panic:%v", err))
		}
	}()
	conn.logger.Debug("handleRPCEventReceived start")

	convertors := map[string]func(event ctypes.ResultEvent, topic *okex.SubscriptionTopic) (interface{}, error){
		okex.DexSpotAccount:     conn.convert2WSTableResponseFromMap,
		okex.DexSpotTicker:      conn.convert2WSTableResponseFromMap,
		okex.DexSpotOrder:       conn.convertWSTableResponseFromList,
		okex.DexSpotAllTicker3s: conn.convertWSTableResponseFromList,
	}

	for {
		select {
		case evt, ok := <-conn.rpcEventChan:
			if !ok {
				break
			}

			topic := query2SubscriptionTopic(evt.Query)
			if topic != nil {
				convertFunc := convertors[topic.Channel]
				if convertFunc == nil {
					convertFunc = conn.convert2WSTableResponseFromMap
				}

				if convertFunc != nil {
					r, e := convertFunc(evt, topic)
					if e == nil {
						conn.cliOutChan <- r
					} else {
						conn.ctx.interruptedCh <- e
						break
					}
				}
			} else {
				conn.logger.Debug("handleRPCEventReceived get event", "event", evt.Events)
				conn.ctx.interruptedCh <- evt
			}
		}
	}

	conn.logger.Debug("handleRPCEventReceived finished")
}

func (conn *okWSConn) cliPing() (err error) {
	msg := "pong"
	conn.cliOutChan <- msg
	return err
}

func (conn *okWSConn) cliSubscribe(op *okex.BaseOp) (err error) {

	// 1. get all of the subscription info
	topics := []*okex.SubscriptionTopic{}
	if op != nil && op.Op == okex.CHNL_EVENT_SUBSCRIBE && op.Args != nil && len(op.Args) > 0 {
		subStrs := op.Args
		for _, subStr := range subStrs {
			topic := okex.FormSubscriptionTopic(subStr)
			if topic == nil {
				continue
			}
			// private channel
			if topic.NeedLogin() {
				if conn.loginAddress == "" {
					errResp := okex.WSErrorResponse{
						Event:     "error",
						Message:   fmt.Sprintf("User not logged in / User must be logined in, before subscribe:%s", topic.Channel),
						ErrorCode: 30041,
					}
					conn.cliOutChan <- errResp
					continue
				}
				topic.Filter = fmt.Sprintf("%s:%s", topic.Filter, conn.loginAddress)
			}
			topics = append(topics, topic)

		}
	} else {
		err = fmt.Errorf("BaseOp {%+v} is not a valid one, expected type: %s", op, okex.CHNL_EVENT_SUBSCRIBE)
	}

	// 2. if rpc client does not exist, create one
	if err == nil && conn.rpcConn == nil {
		rpcAddr := viper.GetString("rpc.laddr")
		conn.logger.Debug("cliSubscribe", "rpc.laddr", rpcAddr)
		// HTTP client can be replaced with LocalClient
		c := rpccli.NewHTTP(rpcAddr, "/websocket")
		conn.rpcConn = c
		err = c.Start()
	}

	// 3. do rpc subscription
	if err == nil && conn.rpcConn != nil {
		subscriber := conn.getSubsciber()
		for _, topic := range topics {

			ctx, _ := context.WithTimeout(context.Background(), maxRpcContextTimeout)
			channel, query := subscriptionTopic2Query(topic)
			eventCh, rpcErr := conn.rpcConn.Subscribe(ctx, subscriber, query)

			if rpcErr == nil {
				conn.logger.Debug(fmt.Sprintf("%s subscribe to %s", subscriber, query))
				eventResp := okex.WSEventResponse{
					Event:   op.Op,
					Channel: channel,
				}
				conn.cliOutChan <- eventResp

				// a goroutine receive result event from rpc websocket client
				go conn.receiveRPCResultEvents(eventCh, subscriber, channel)

			} else {
				errResp := okex.WSErrorResponse{
					Event:     "error",
					Message:   fmt.Sprintf("fail to subscribe %s, error: %s", channel, rpcErr.Error()),
					ErrorCode: 30043,
				}
				conn.cliOutChan <- errResp
			}
		}
	}

	// 4. push initial data
	initialDataMap := map[string]func(topic *okex.SubscriptionTopic){
		okex.DexSpotDepthBook: conn.initialDepthBook,
	}
	for _, topic := range topics {
		initialDataFunc, ok := initialDataMap[topic.Channel]
		if !ok {
			continue
		}
		initialDataFunc(topic)
	}

	return err
}

func (conn *okWSConn) initialDepthBook(topic *okex.SubscriptionTopic) {
	depthBookRes, ok := GetDepthBookFromCache(topic.Filter)
	conn.logger.Debug("initialDepthBook", "depthBookRes", depthBookRes, "ok", ok)
	if !ok {
		return
	}
	resp := okex.WSTableResponse{
		Table:  topic.Channel,
		Action: "partial",
		Data:   []interface{}{depthBookRes},
	}
	conn.cliOutChan <- resp
}

func (conn *okWSConn) receiveRPCResultEvents(eventCh <-chan ctypes.ResultEvent, subscriber, channel string) {
	conn.logger.Debug("receiveRPCResultEvents start", subscriber, channel)

	time.Sleep(time.Millisecond)

	for {
		select {
		case event, ok := <-eventCh:
			if !ok {
				conn.logger.Debug("receiveRPCResultEvents's eventCh is closed or something else", conn.getSubsciber(), event.Query)
				break
			}
			conn.rpcEventChan <- event
			//conn.logger.Debug("receiveRPCResultEvents get event", conn.getSubsciber(), event.Query)

		case <-conn.rpcStopChan:
			break
		}
	}

	conn.logger.Debug("receiveRPCResultEvents finished", subscriber, channel)
}

func (conn *okWSConn) getSubsciber() string {
	return conn.cliConn.RemoteAddr().String()
}

func (conn *okWSConn) cliUnSubscribe(op *okex.BaseOp) (err error) {
	// 1. check op is a valid unsubscribe op
	topics := []*okex.SubscriptionTopic{}
	if op != nil && op.Op == okex.CHNL_EVENT_UNSUBSCRIBE && op.Args != nil && len(op.Args) > 0 {
		subStrs := op.Args
		for _, subStr := range subStrs {
			topic := okex.FormSubscriptionTopic(subStr)
			if topic == nil {
				continue
			}
			// private channel
			if topic.NeedLogin() {
				if conn.loginAddress == "" {
					errResp := okex.WSErrorResponse{
						Event:     "error",
						Message:   fmt.Sprintf("User not logged in / User must be logined in, before subscribe:%s", topic.Channel),
						ErrorCode: 30041,
					}
					conn.cliOutChan <- errResp
					continue
				}
				topic.Filter = fmt.Sprintf("%s:%s", topic.Filter, conn.loginAddress)
			}
			topics = append(topics, topic)
		}
	} else {
		err = fmt.Errorf("BaseOp {%+v} is not a valid one, expected type: %s", op, okex.CHNL_EVENT_UNSUBSCRIBE)
	}

	if conn.rpcConn == nil {
		// 2. if rpcConn is not initialized, raise error
		err = fmt.Errorf("RPC WS Client hasn't been initialized properly")
	} else {
		// 3. do unsubscibe work
		subscriber := conn.getSubsciber()
		for _, topic := range topics {
			ctx, _ := context.WithTimeout(context.Background(), maxRpcContextTimeout)
			channel, query := subscriptionTopic2Query(topic)
			rpcErr := conn.rpcConn.Unsubscribe(ctx, subscriber, query)
			if rpcErr == nil {
				conn.logger.Debug(fmt.Sprintf("%s unsubscribe to %s", subscriber, query))
				eventResp := okex.WSEventResponse{
					Event:   op.Op,
					Channel: channel,
				}
				conn.cliOutChan <- eventResp

			} else {
				errResp := okex.WSErrorResponse{
					Event:     "error",
					Message:   fmt.Sprintf("fail to unsubscribe %s, error: %s", channel, rpcErr.Error()),
					ErrorCode: 30043,
				}
				conn.cliOutChan <- errResp
			}
		}
	}

	return err
}

func (conn *okWSConn) cliLogin(op *okex.BaseOp) error {
	if op == nil || op.Op != okex.CHNL_EVENT_LOGIN || len(op.Args) != 1 {
		err := fmt.Errorf("invalid request, when doing: %s", okex.CHNL_EVENT_LOGIN)
		errResp := okex.WSErrorResponse{
			Event:     "error",
			Message:   err.Error(),
			ErrorCode: 30043,
		}
		conn.cliOutChan <- errResp

		conn.logger.Error(err.Error())
		return err
	}
	conn.loginAddress = op.Args[0]
	return nil
}

func (conn *okWSConn) handleConvert() {
	defer func() {
		if err := recover(); err != nil {
			conn.logger.Error(fmt.Sprintf("handleConvert recover panic:%v", err))
		}
	}()
	conn.logger.Debug("handleConvert start")

	//conn.wg.Add(1)
	//defer conn.wg.Done()

	cliEventMap := map[string]func(op *okex.BaseOp) error{
		okex.CHNL_EVENT_SUBSCRIBE:   conn.cliSubscribe,
		okex.CHNL_EVENT_UNSUBSCRIBE: conn.cliUnSubscribe,
		okex.CHNL_EVENT_LOGIN:       conn.cliLogin,
	}

	for {
		var err error
		select {
		case cliInMsg, ok := <-conn.cliInChan:
			if !ok {
				break
			}

			op := okex.BaseOp{}
			if jsonErr := json.Unmarshal(cliInMsg, &op); jsonErr == nil {
				conn.logger.Debug(fmt.Sprintf("handleConvert BaseOp: %+v", op))
				f := cliEventMap[op.Op]
				err = f(&op)
			} else {
				if string(cliInMsg) == "ping" {
					err = conn.cliPing()
				}
			}
		}

		if err != nil {
			conn.ctx.interruptedCh <- err
			break
		}
	}

	conn.logger.Debug("handleConvert finished")
}
