package quoteslite

import (
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	wsEndpoint = "/websocket"
)

var (
	upgrader = websocket.Upgrader{
		EnableCompression: true,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func bridgeMsgHandler(w http.ResponseWriter, r *http.Request, logger log.Logger) {
	logger.Debug(fmt.Sprintf("bridgeMsgHandler remoteAddr: %s", r.RemoteAddr))
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Debug(fmt.Sprintf("bridgeMsgHandler error: %s", err.Error()))
		return
	}

	c.SetPingHandler(func(appData string) error {
		return c.WriteControl(websocket.PongMessage, []byte(string("pong")), time.Now().Add(writeWait))
	})

	connCtx := newOKWSContext()
	signal.Notify(connCtx.signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	newOKWSConn(connCtx, c, logger)
}

func bridgeMsgHandlerWithLogger(logger log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bridgeMsgHandler(w, r, logger)
	}
}

func StartWSServer(logger log.Logger, endpoint string) {
	http.HandleFunc("/ws/v3", bridgeMsgHandlerWithLogger(logger))
	logger.Info("Starting WebSocket server on ", endpoint)
	http.ListenAndServe(endpoint, nil)
}
