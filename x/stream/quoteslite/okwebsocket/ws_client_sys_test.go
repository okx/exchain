package okex

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/okex/okchain/x/common"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
)

func TestOKexWebsocketSmoke(t *testing.T) {
	//common.SkipSysTestChecker(t)

	agent := OKWSAgent{}
	config := GetDefaultConfig()
	config.WSEndpoint = "wss://dexcomreal.bafang.com:8443/"
	config.IsPrint = true
	agent.config = config

	// Step1: Start agent.
	agent.Start(config)

	// Step2: Subscribe Channel
	// Step2.0: Subscribe public Channel swap/ticker successfully.
	agent.Subscribe("dex_spot/ticker", "tbtc_tusdk", printReceivedData)

	// Step3 stop
	time.Sleep(time.Minute * 10)
	agent.Stop()
}

func TestLocalNodeWebsocketSmoke(t *testing.T) {

	//common.SkipSysTestChecker(t)

	agent := OKWSAgent{}
	config := GetDefaultConfig()
	config.WSEndpoint = "ws://0.0.0.0:6666/"
	config.IsPrint = true
	agent.config = config

	// Step1: Start agent.
	agent.Start(config)

	// Step2: Subscribe Channel
	// Step2.0: Subscribe public Channel swap/ticker successfully.
	//agent.Subscribe("dex_spot/ticker", "tbtc_tusdk", defaultPrintData)
	//agent.Subscribe("dex_spot/account",
	//	"okt:okchain10q0rk5qnyag7wfvvt7rtphlw589m7frsmyq4ya", nil)
	//agent.Subscribe("dex_spot/candle60s", "eos-37c_okt", nil)
	//agent.Subscribe("dex_spot/candle60s", "eos-d78_okt", nil)
	//agent.Subscribe("dex_spot/candle180s", "eos-37c_okt", nil)
	//agent.Subscribe("dex_spot/candle300s", "eos-37c_okt", nil)
	//agent.Subscribe("dex_spot/candle900s", "eos-37c_okt", nil)
	//agent.Subscribe("dex_spot/candle1800s", "eos-37c_okt", nil)
	//agent.Subscribe("dex_spot/candle3600s", "eos-37c_okt", nil)
	//agent.Subscribe("dex_spot/candle7200s", "eos-37c_okt", nil)
	//agent.Subscribe("dex_spot/ticker", "eos-37c_okt", nil)
	agent.Subscribe("dex_spot/all_ticker_3s", "", nil)

	// Step3 stop
	time.Sleep(time.Second * 601)
	agent.Stop()
}

func TestOKExWSSubscribe(t *testing.T) {

	common.SkipSysTestChecker(t)

	// https://www.okex.me/dex-test/spot/trade?debug_push=true
	c, _, _ := websocket.DefaultDialer.Dial("wss://dexcomreal.bafang.com:8443/ws/v3", nil)
	defer c.Close()

	m := map[string]string{
		"op":   "subscribe",
		"args": "dex_spot/ticker:tbtc_tusdk",
	}

	c.EnableWriteCompression(true)

	sMsg, err := json.Marshal(m)
	require.Nil(t, err, err)
	err = c.WriteMessage(websocket.TextMessage, sMsg)
	require.Nil(t, err, err)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			r, _ := GzipDecode(message)

			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %+v", string(r))
			//log.Printf("recv: %s", string(message))
		}
	}()

	ticker := time.NewTicker(time.Second * 3)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.PingMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}

}
