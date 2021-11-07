package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"github.com/spf13/viper"

	"github.com/gorilla/websocket"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	rpccli "github.com/okex/exchain/libs/tendermint/rpc/client/http"
	ctypes "github.com/okex/exchain/libs/tendermint/rpc/core/types"
)

type Context struct {
	interruptedCh chan interface{}
	signalCh      chan os.Signal
}

func (ctx *Context) closeAll() {
	close(ctx.signalCh)
	close(ctx.interruptedCh)
}

func newContext() *Context {
	ctx := Context{
		interruptedCh: make(chan interface{}, 4),
		signalCh:      make(chan os.Signal),
	}

	return &ctx
}

type Conn struct {
	cliConn      *websocket.Conn
	rpcConn      *rpccli.HTTP
	ctx          *Context
	logger       log.Logger
	loginAddress string

	cliInChan    chan []byte
	cliOutChan   chan interface{}
	rpcEventChan chan ctypes.ResultEvent
	rpcStopChan  chan interface{}
}

func newOKWSConn(ctx *Context, cliConn *websocket.Conn, logger log.Logger) *Conn {
	if ctx == nil || cliConn == nil {
		return nil
	}

	conn := Conn{
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

func (conn *Conn) start() {
	conn.logger.Debug("starting bi-direction communication")

	go conn.handleFinalise()
	go conn.handleCliRead()
	go conn.handleCliWrite()
	go conn.handleRPCEventReceived()
	go conn.handleConvert()
}

func (conn *Conn) stopAll() {
	conn.logger.Debug("okWSConn.stopAll start")

	// 1. close all the connection
	conn.logger.Debug("okWSConn.stopAll close bi-connection")
	if conn.cliConn != nil {
		conn.cliConn.Close()
		conn.logger.Debug("okWSConn'connection to client is closed.")
	}

	if conn.rpcConn != nil {
		err := conn.rpcConn.Stop()
		if err != nil {
			conn.logger.Error("rpcConn stop error", "msg", err.Error())
		}
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

func (conn *Conn) handleFinalise() {
	conn.logger.Debug("handleFinalise start")

	msg := <-conn.ctx.interruptedCh
	conn.logger.Debug("handleFinalise get interrupted signal", "msg", msg)
	conn.stopAll()

	conn.logger.Debug("handleFinalise finished")
}

func (conn *Conn) handleCliRead() {
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
				txtMsg, err = gzipDecode(msg)
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

func (conn *Conn) handleCliWrite() {
	defer func() {
		if err := recover(); err != nil {
			conn.logger.Error(fmt.Sprintf("handleCliWrite recover panic:%v", err))
		}
	}()

	conn.logger.Debug("handleCliWrite start")
	for outMsg := range conn.cliOutChan {
		var err error
		if msg, ok := outMsg.(string); ok {
			err = conn.cliConn.WriteMessage(websocket.TextMessage, []byte(msg))
		} else {
			err = conn.cliConn.WriteJSON(outMsg)
		}

		conn.logger.Debug("handleCliWrite write", "OutMsg", outMsg)

		if err != nil {
			conn.ctx.interruptedCh <- err
			break
		}
	}
}

func (conn *Conn) convert2WSTableResponseFromMap(resultEvt ctypes.ResultEvent, topic *SubscriptionTopic) (r interface{}, e error) {
	innerChannel, e := topic.ToString()
	if e != nil {
		return nil, e
	}

	resp := TableResponse{
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

func (conn *Conn) convertWSTableResponseFromList(resultEvt ctypes.ResultEvent, topic *SubscriptionTopic) (r interface{}, e error) {
	innerChannel, e := topic.ToString()
	if e != nil {
		return nil, e
	}

	resp := TableResponse{
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

func (conn *Conn) handleRPCEventReceived() {
	defer func() {
		if err := recover(); err != nil {
			conn.logger.Error(fmt.Sprintf("handleRPCEventReceived recover panic:%v", err))
		}
	}()
	conn.logger.Debug("handleRPCEventReceived start")

	convertors := map[string]func(event ctypes.ResultEvent, topic *SubscriptionTopic) (interface{}, error){
		DexSpotAccount:     conn.convert2WSTableResponseFromMap,
		DexSpotTicker:      conn.convert2WSTableResponseFromMap,
		DexSpotOrder:       conn.convertWSTableResponseFromList,
		DexSpotAllTicker3s: conn.convertWSTableResponseFromList,
	}

	for evt := range conn.rpcEventChan {
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

func (conn *Conn) cliPing() (err error) {
	msg := "pong"
	conn.cliOutChan <- msg
	return err
}

func (conn *Conn) cliSubscribe(op *BaseOp) (err error) {

	// 1. get all of the subscription info
	var topics []*SubscriptionTopic
	if op != nil && op.Op == eventSubscribe && op.Args != nil && len(op.Args) > 0 {
		subStrs := op.Args
		for _, subStr := range subStrs {
			topic := FormSubscriptionTopic(subStr)
			if topic == nil {
				continue
			}
			// private channel
			if topic.NeedLogin() {
				if conn.loginAddress == "" {
					errResp := ErrorResponse{
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
		// nolint
		err = fmt.Errorf("BaseOp {%+v} is not a valid one, expected type: %s", op, eventSubscribe)
	}

	// 2. if rpc client does not exist, create one
	if err == nil && conn.rpcConn == nil {
		rpcAddr := viper.GetString("rpc.laddr")
		conn.logger.Debug("cliSubscribe", "rpc.laddr", rpcAddr)
		// HTTP client can be replaced with LocalClient
		c, err := rpccli.New(rpcAddr, "/websocket")
		if err != nil {
			conn.logger.Error("cliSubscribe", "error", err.Error())
			return err
		}
		conn.rpcConn = c
		err = c.Start()
	}

	// 3. do rpc subscription
	if err == nil && conn.rpcConn != nil {
		subscriber := conn.getSubsciber()
		for _, topic := range topics {
			ctx, cancel := context.WithTimeout(context.Background(), maxRPCContextTimeout)

			channel, query := subscriptionTopic2Query(topic)
			eventCh, rpcErr := conn.rpcConn.Subscribe(ctx, subscriber, query)

			if rpcErr == nil {
				conn.logger.Debug(fmt.Sprintf("%s subscribe to %s", subscriber, query))
				eventResp := EventResponse{
					Event:   op.Op,
					Channel: channel,
				}
				conn.cliOutChan <- eventResp

				// a goroutine receive result event from rpc websocket client
				go conn.receiveRPCResultEvents(eventCh, subscriber, channel)

			} else {
				errResp := ErrorResponse{
					Event:     "error",
					Message:   fmt.Sprintf("fail to subscribe %s, error: %s", channel, rpcErr.Error()),
					ErrorCode: 30043,
				}
				conn.cliOutChan <- errResp
			}
			cancel()
		}
	}

	// 4. push initial data
	initialDataMap := map[string]func(topic *SubscriptionTopic){
		DexSpotDepthBook: conn.initialDepthBook,
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

func (conn *Conn) initialDepthBook(topic *SubscriptionTopic) {
	depthBookRes, ok := GetDepthBookFromCache(topic.Filter)
	conn.logger.Debug("initialDepthBook", "depthBookRes", depthBookRes, "ok", ok)
	if !ok {
		return
	}
	resp := TableResponse{
		Table:  topic.Channel,
		Action: "partial",
		Data:   []interface{}{depthBookRes},
	}
	conn.cliOutChan <- resp
}

func (conn *Conn) receiveRPCResultEvents(eventCh <-chan ctypes.ResultEvent, subscriber, channel string) {
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

		case <-conn.rpcStopChan:
			break
		}
	}
}

func (conn *Conn) getSubsciber() string {
	return conn.cliConn.RemoteAddr().String()
}

func (conn *Conn) cliUnSubscribe(op *BaseOp) (err error) {
	// 1. check op is a valid unsubscribe op
	var topics []*SubscriptionTopic
	if op != nil && op.Op == eventUnsubscribe && op.Args != nil && len(op.Args) > 0 {
		subStrs := op.Args
		for _, subStr := range subStrs {
			topic := FormSubscriptionTopic(subStr)
			if topic == nil {
				continue
			}
			// private channel
			if topic.NeedLogin() {
				if conn.loginAddress == "" {
					errResp := ErrorResponse{
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
		// nolint
		err = fmt.Errorf("BaseOp {%+v} is not a valid one, expected type: %s", op, eventUnsubscribe)
	}

	if conn.rpcConn == nil {
		// 2. if rpcConn is not initialized, raise error
		err = fmt.Errorf("RPC WS Client hasn't been initialized properly")
	} else {
		// 3. do unsubscribe work
		subscriber := conn.getSubsciber()
		for _, topic := range topics {
			ctx, cancel := context.WithTimeout(context.Background(), maxRPCContextTimeout)
			channel, query := subscriptionTopic2Query(topic)
			rpcErr := conn.rpcConn.Unsubscribe(ctx, subscriber, query)
			if rpcErr == nil {
				conn.logger.Debug(fmt.Sprintf("%s unsubscribe to %s", subscriber, query))
				eventResp := EventResponse{
					Event:   op.Op,
					Channel: channel,
				}
				conn.cliOutChan <- eventResp

			} else {
				errResp := ErrorResponse{
					Event:     "error",
					Message:   fmt.Sprintf("fail to unsubscribe %s, error: %s", channel, rpcErr.Error()),
					ErrorCode: 30043,
				}
				conn.cliOutChan <- errResp
			}
			cancel()
		}
	}

	return err
}

func (conn *Conn) cliLogin(op *BaseOp) error {
	if op == nil || op.Op != eventLogin || len(op.Args) != 1 {
		err := fmt.Errorf("invalid request, when doing: %s", eventLogin)
		errResp := ErrorResponse{
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

func (conn *Conn) handleConvert() {
	defer func() {
		if err := recover(); err != nil {
			conn.logger.Error(fmt.Sprintf("handleConvert recover panic:%v", err))
		}
		conn.logger.Debug("handleConvert finished")
	}()
	conn.logger.Debug("handleConvert start")

	cliEventMap := map[string]func(op *BaseOp) error{
		eventSubscribe:   conn.cliSubscribe,
		eventUnsubscribe: conn.cliUnSubscribe,
		eventLogin:       conn.cliLogin,
	}

	for cliInMsg := range conn.cliInChan {
		var err error
		op := BaseOp{}
		if jsonErr := json.Unmarshal(cliInMsg, &op); jsonErr == nil {
			conn.logger.Debug(fmt.Sprintf("handleConvert BaseOp: %+v", op))
			f := cliEventMap[op.Op]
			err = f(&op)
		} else if string(cliInMsg) == "ping" {
			err = conn.cliPing()
		}

		if err != nil {
			conn.ctx.interruptedCh <- err
			break
		}
	}
}
