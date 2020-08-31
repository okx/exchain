package okex

/*
 OKEX websocket API agent
 @author Lingting Fu
 @date 2018-12-27
 @version 1.0.0
*/

import (
	"bytes"
	"compress/flate"
	"io/ioutil"

	"github.com/gorilla/websocket"

	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

type OKWSAgent struct {
	baseUrl string
	config  *Config
	conn    *websocket.Conn

	wsEvtCh  chan interface{}
	wsErrCh  chan interface{}
	wsTbCh   chan interface{}
	stopCh   chan interface{}
	errCh    chan error
	signalCh chan os.Signal

	subMap         map[string][]ReceivedDataCallback
	activeChannels map[string]bool
	hotDepthsMap   map[string]*WSHotDepths

	processMut sync.Mutex
}

func (a *OKWSAgent) SetConfig(config *Config) {
	a.config = config
}

func (a *OKWSAgent) Start(config *Config) error {
	a.config = config
	a.baseUrl = config.WSEndpoint + "ws/v3?compress=true"
	log.Printf("Connecting to %s", a.baseUrl)
	c, _, err := websocket.DefaultDialer.Dial(a.baseUrl, nil)

	if err != nil {
		log.Fatalf("dial:%+v", err)
		return err
	} else {
		if a.config.IsPrint {
			log.Printf("Connected to %s", a.baseUrl)
		}
		a.conn = c
		a.config = config

		a.wsEvtCh = make(chan interface{})
		a.wsErrCh = make(chan interface{})
		a.wsTbCh = make(chan interface{})
		a.errCh = make(chan error)
		a.stopCh = make(chan interface{}, 16)
		a.signalCh = make(chan os.Signal)
		a.activeChannels = make(map[string]bool)
		a.subMap = make(map[string][]ReceivedDataCallback)
		a.hotDepthsMap = make(map[string]*WSHotDepths)

		signal.Notify(a.signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		go a.work()
		go a.receive()
		go a.finalize()
	}

	return nil
}

func (a *OKWSAgent) Subscribe(channel, filter string, cb ReceivedDataCallback) error {
	a.processMut.Lock()
	defer a.processMut.Unlock()

	st := SubscriptionTopic{channel, filter}
	bo, err := subscribeOp([]*SubscriptionTopic{&st})
	if err != nil {
		return err
	}

	msg, err := Struct2JsonString(bo)
	if a.config.IsPrint {
		log.Printf("Send Msg: %s", msg)
	}
	if err := a.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return err
	}

	cbs := a.subMap[st.Channel]
	if cbs == nil {
		cbs = []ReceivedDataCallback{}
		a.activeChannels[st.Channel] = false
	}

	if cb != nil {
		cbs = append(cbs, cb)
		fullTopic, _ := st.ToString()
		a.subMap[st.Channel] = cbs
		a.subMap[fullTopic] = cbs
	}

	return nil
}

func (a *OKWSAgent) UnSubscribe(channel, filter string) error {
	a.processMut.Lock()
	defer a.processMut.Unlock()

	st := SubscriptionTopic{channel, filter}
	bo, err := unsubscribeOp([]*SubscriptionTopic{&st})
	if err != nil {
		return err
	}

	msg, err := Struct2JsonString(bo)
	if a.config.IsPrint {
		log.Printf("Send Msg: %s", msg)
	}
	if err := a.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return err
	}

	a.subMap[channel] = nil
	a.activeChannels[channel] = false

	return nil
}

func (a *OKWSAgent) Login(apiKey, passphrase string) error {

	timestamp := EpochTime()

	preHash := PreHashString(timestamp, GET, "/users/self/verify", "")
	if sign, err := HmacSha256Base64Signer(preHash, a.config.SecretKey); err != nil {
		return err
	} else {
		op, err := loginOp(apiKey, passphrase, timestamp, sign)
		data, err := Struct2JsonString(op)
		log.Printf("Send Msg: %s", data)
		err = a.conn.WriteMessage(websocket.TextMessage, []byte(data))
		if err != nil {
			return err
		}
		time.Sleep(time.Millisecond * 100)
	}
	return nil
}

func (a *OKWSAgent) DexLogin(address string) error {
	op := BaseOp{
		Op:   "dex_jwt",
		Args: []string{address},
	}
	data, err := Struct2JsonString(op)
	log.Printf("Send Msg: %s", data)
	err = a.conn.WriteMessage(websocket.TextMessage, []byte(data))
	if err != nil {
		return err
	}
	time.Sleep(time.Millisecond * 100)
	return nil
}

func (a *OKWSAgent) keepalive() {
	a.ping()
}

func (a *OKWSAgent) Stop() error {
	defer func() {
		a := recover()
		log.Printf("Stop End. Recover msg: %+v", a)
	}()

	a.stopCh <- nil
	return nil
}

func (a *OKWSAgent) finalize() error {
	defer func() {
		log.Printf("Finalize End. Connection to WebSocket is closed.")
	}()

	select {
	case <-a.stopCh:
		if a.conn != nil {
			//close(a.errCh)
			//close(a.wsTbCh)
			//close(a.wsEvtCh)
			//close(a.wsErrCh)
			return a.conn.Close()
		}
	}

	return nil
}

func (a *OKWSAgent) ping() {
	msg := "ping"
	//log.Printf("Send Msg: %s", msg)
	a.conn.WriteMessage(websocket.TextMessage, []byte(msg))
}

func GzipDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return ioutil.ReadAll(reader)
}

func (a *OKWSAgent) handleErrResponse(r interface{}) error {
	log.Printf("handleErrResponse %+v \n", r)
	return nil
}

func (a *OKWSAgent) handleEventResponse(r interface{}) error {
	er := r.(*WSEventResponse)
	a.activeChannels[er.Channel] = (er.Event == CHNL_EVENT_SUBSCRIBE)
	return nil
}

func (a *OKWSAgent) handleTableResponse(r interface{}) error {
	tb := ""
	switch r.(type) {
	case *WSTableResponse:
		tb = r.(*WSTableResponse).Table
	case *WSDepthTableResponse:
		tb = r.(*WSDepthTableResponse).Table
	}

	cbs := a.subMap[tb]
	if cbs != nil {
		for i := 0; i < len(cbs); i++ {
			cb := cbs[i]
			if err := cb(r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *OKWSAgent) work() {
	defer func() {
		e := recover()
		log.Printf("Work End. Recover msg: %+v", e)
		if e != nil {
			debug.PrintStack()
		}
	}()

	ticker := time.NewTicker(29 * time.Second)
	//defer ticker.Stop()

	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				break
			}
			a.keepalive()
		case errR := <-a.wsErrCh:
			a.handleErrResponse(errR)
		case evtR := <-a.wsEvtCh:
			a.handleEventResponse(evtR)
		case tb := <-a.wsTbCh:
			a.handleTableResponse(tb)
		case <-a.signalCh:
			break
		case err := <-a.errCh:
			DefaultDataCallBack(err)
			break
		case <-a.stopCh:
			return

		}
	}
}

func (a *OKWSAgent) receive() {
	defer func() {
		a := recover()
		if a != nil {
			log.Printf("Receive End. Recover msg: %+v", a)
			//debug.PrintStack()
		}
	}()

	for {

		messageType, message, err := a.conn.ReadMessage()
		if err != nil {
			a.errCh <- err
			log.Printf("ReadMessage error:%s", err.Error())
			break
		}

		txtMsg := message
		switch messageType {
		case websocket.TextMessage:
		case websocket.BinaryMessage:
			txtMsg, err = GzipDecode(message)
		}

		rsp, err := loadResponse(txtMsg)
		if a.config.IsPrint {
			if rsp != nil {
				log.Printf("LoadedRep: %+v, err: %+v", rsp, err)
			} else {
				log.Printf("TextMsg: %s", txtMsg)
			}
		}

		if err != nil {
			log.Printf("loadResponse error: %+v, quit receive goroutine, textMsg: %s",
				err, string(txtMsg))
			break
		}

		switch rsp.(type) {
		case *WSErrorResponse:
			a.wsErrCh <- rsp
		case *WSEventResponse:
			er := rsp.(*WSEventResponse)
			a.wsEvtCh <- er
		case *WSDepthTableResponse:
			var err error
			dtr := rsp.(*WSDepthTableResponse)
			hotDepths := a.hotDepthsMap[dtr.Table]
			if hotDepths == nil {
				hotDepths = NewWSHotDepths(dtr.Table)
				err = hotDepths.loadWSDepthTableResponse(dtr)
				if err == nil {
					a.hotDepthsMap[dtr.Table] = hotDepths
				}
			} else {
				err = hotDepths.loadWSDepthTableResponse(dtr)
			}

			if err == nil {
				a.wsTbCh <- dtr
			} else {
				log.Printf("Failed to loadWSDepthTableResponse, dtr: %+v, err: %+v", dtr, err)
			}

		case *WSTableResponse:
			tb := rsp.(*WSTableResponse)
			a.wsTbCh <- tb
		default:
			//log.Println(rsp)
		}
	}

	log.Printf("OKWSAgent.receive goroutine quit")
}
