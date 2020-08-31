package wsclient

import (
	"fmt"
	okex "github.com/okex/okchain/x/stream/quoteslite/okwebsocket"
	"reflect"
	"time"
)

type IWSService interface {
	GetEndPoint() string
	ReceivedDataCallback(interface{}) error
}

type BaseWS struct {
	endPoint       string
	receivedEvents []interface{}
}

func (ws *BaseWS) GetEndPoint() string {
	return ws.endPoint
}

func (ws *BaseWS) ReceivedDataCallback(r interface{}) (e error) {
	ws.receivedEvents = append(ws.receivedEvents, r)
	return
}

func CapturePublicChannelNotices(ws IWSService, channel, filter string) {
	// 0. start WebSocketAgent agent

	agent := okex.OKWSAgent{}
	config := okex.GetDefaultConfig()
	config.WSEndpoint = ws.GetEndPoint()
	config.IsPrint = true

	// 1. subscibe channel described in {ws}
	agent.Start(config)
	agent.Subscribe(channel, filter, ws.ReceivedDataCallback)

	// 2. wait {captureSeconds} second, capture websocket responses at the same time
	timer := time.NewTimer(70 * time.Second)

	select {
	case <-timer.C:
		fmt.Println("70s time's up.")
	}

	// 3. unsubscribe
	agent.UnSubscribe(channel, filter)
	time.Sleep(3 * time.Second)

	// 4. stop agent
	agent.Stop()

	return
}

func mergeEventTypes(events []interface{}) map[reflect.Type][]interface{} {
	r := map[reflect.Type][]interface{}{}

	for _, e := range events {
		t := reflect.TypeOf(e)
		items := r[t]
		if items == nil {
			r[t] = []interface{}{e}
		} else {
			r[t] = append(r[t], e)
		}
	}
	return r
}

func mergeTableTypes(events []interface{}) map[string][]interface{} {
	r := map[string][]interface{}{}
	for _, e := range events {
		if reflect.TypeOf(e) == reflect.TypeOf(&okex.WSTableResponse{}) {
			tr := e.(*okex.WSTableResponse)
			tbEvts := r[tr.Table]
			if tbEvts == nil {
				r[tr.Table] = []interface{}{e}
			} else {
				r[tr.Table] = append(tbEvts, e)
			}
		}
	}
	return r
}

func CompareEventsDiff(expectedEvents, capturedEvents []interface{}, verbose bool) bool {

	eMEvents := mergeEventTypes(expectedEvents)
	cMEvents := mergeEventTypes(capturedEvents)

	noDiff := true
	if len(eMEvents) != len(cMEvents) {
		fmt.Printf("[E] Event type count mismatched. expected Cnt: %d, captured Cnt: %d\n",
			len(eMEvents), len(cMEvents))
		noDiff = false
	}

	for t, expected := range eMEvents {
		captured := cMEvents[t]
		fmt.Printf("[I] EventType: %s, expected Cnt: %d, captured Cnt: %d\n", t, len(expected), len(captured))

		// 1. compare event types count
		if captured == nil || len(captured) == 0 {
			fmt.Printf("[E] Event type %t not found in the captured events.\n", t)
			noDiff = false
			continue
		}

		eTableEvents := mergeEventTypes(expected)
		cTableEvents := mergeEventTypes(captured)

		for expTbName, expTbEvents := range eTableEvents {

			// 2. check if every single event in expectedEvents can be found in capturedEvents
			capTbEvts := cTableEvents[expTbName]
			fmt.Printf("[I] TableName: %s, expected Cnt: %d, captured Cnt: %d\n",
				expTbName, len(expTbEvents), len(capTbEvts))

			if capTbEvts == nil || len(capTbEvts) == 0 {
				fmt.Printf("[E] No %s found in Captured Table Response\n", expTbName)
				noDiff = false
				continue
			}

			// 3. check every event type if the response data type is the same
			switch (expectedEvents[0]).(type) {
			case *okex.WSTableResponse:
				expEvt := (expectedEvents[0]).(*okex.WSTableResponse)
				capEvt := (capTbEvts[0]).(*okex.WSTableResponse)

				eEvtData := expEvt.Data[0].(map[string]interface{})
				cEvtData := capEvt.Data[0].(map[string]interface{})

				for ek, _ := range eEvtData {
					_, ok := cEvtData[ek]
					if !ok {
						fmt.Printf("[E] %s not found in the captured event data list, ExpData: %+v, CapData: %+v\n",
							ek, eEvtData, cEvtData)
						noDiff = false
					}
				}

			case *okex.WSDepthTableResponse:
				expEvt := (expectedEvents[0]).(*okex.WSDepthTableResponse)
				capEvt := (capTbEvts[0]).(*okex.WSDepthTableResponse)

				if len(expEvt.Data) > 0 && len(capEvt.Data) > 0 {
					continue
				} else {
					fmt.Printf("[E] Both WSDepthTableResponse 's Data is empty. Expected: %+v, Captured: %+v",
						expEvt, capEvt)
				}
			}
		}

	}

	return noDiff
}
