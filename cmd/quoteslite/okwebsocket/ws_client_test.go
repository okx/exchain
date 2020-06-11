package okex

import (
	"bytes"
	"compress/flate"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"testing"
	"time"
)

func TestOKexWebsocketSmoke(t *testing.T) {
	agent := OKWSAgent{}
	config := GetDefaultConfig()
	config.WSEndpoint = "wss://dexcomreal.bafang.com:8443/"
	config.IsPrint = true
	agent.config = config

	// Step1: Start agent.
	agent.Start(config)

	// Step2: Subscribe Channel
	// Step2.0: Subscribe public Channel swap/ticker successfully.
	agent.Subscribe("dex_spot/ticker", "tbtc_tusdk", defaultPrintData)

	// Step3 stop
	time.Sleep(time.Second * 60)
	agent.Stop()
}


func TestLocalNodeWebsocketSmoke(t *testing.T) {
	agent := OKWSAgent{}
	config := GetDefaultConfig()
	config.WSEndpoint = "ws://localhost:6666/"
	config.IsPrint = false
	agent.config = config

	// Step1: Start agent.
	agent.Start(config)

	// Step2: Subscribe Channel
	// Step2.0: Subscribe public Channel swap/ticker successfully.
	agent.Subscribe("dex_spot/ticker", "tbtc_tusdk", defaultPrintData)

	// Step3 stop
	time.Sleep(time.Second * 30)
	agent.Stop()
}

func GzipDecode(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func TestOKExWSSubscribe(t *testing.T) {
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
