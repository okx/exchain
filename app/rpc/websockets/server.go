package websockets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/okex/exchain/x/common/monitor"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
)

// Server defines a server that handles Ethereum websockets.
type Server struct {
	rpcAddr string // listen address of rest-server
	wsAddr  string // listen address of ws server
	api     *PubSubAPI
	logger  log.Logger

	connPool       chan struct{}
	connPoolLock   *sync.Mutex
	currentConnNum metrics.Gauge
	maxConnNum     metrics.Gauge
}

// NewServer creates a new websocket server instance.
func NewServer(clientCtx context.CLIContext, log log.Logger, wsAddr string) *Server {
	restServerAddr := viper.GetString(server.FlagListenAddr)
	parts := strings.SplitN(restServerAddr, "://", 2)
	if len(parts) != 2 {
		panic(fmt.Errorf("invalid listening address %s (use fully formed addresses, including the tcp:// or unix:// prefix)", restServerAddr))
	}
	url := parts[1]
	urlParts := strings.SplitN(url, ":", 2)
	if len(urlParts) != 2 {
		panic(fmt.Errorf("invalid listening address %s (use ip:port as an url)", url))
	}
	port := urlParts[1]

	return &Server{
		rpcAddr:      "http://localhost:" + port,
		wsAddr:       wsAddr,
		api:          NewAPI(clientCtx, log),
		logger:       log.With("module", "websocket-server"),
		connPool:     make(chan struct{}, viper.GetInt(server.FlagWsMaxConnections)),
		connPoolLock: new(sync.Mutex),
		currentConnNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: monitor.XNameSpace,
			Subsystem: "websocket",
			Name:      "connection_number",
			Help:      "the number of current websocket client connections",
		}, nil),
		maxConnNum: prometheus.NewGaugeFrom(stdprometheus.GaugeOpts{
			Namespace: monitor.XNameSpace,
			Subsystem: "websocket",
			Name:      "connection_capacity",
			Help:      "the capacity number of websocket client connections",
		}, nil),
	}
}

// Start runs the websocket server
func (s *Server) Start() {
	ws := mux.NewRouter()
	ws.Handle("/", s)
	s.maxConnNum.Set(float64(viper.GetInt(server.FlagWsMaxConnections)))
	s.currentConnNum.Set(0)

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", s.wsAddr), ws)
		if err != nil {
			s.logger.Error("http error:", err)
		}
	}()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.connPoolLock.Lock()
	defer s.connPoolLock.Unlock()
	if len(s.connPool) >= cap(s.connPool) {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("websocket upgrade failed", " error", err)
		return
	}

	s.connPool <- struct{}{}
	s.currentConnNum.Set(float64(len(s.connPool)))
	go s.readLoop(wsConn)
}

func (s *Server) sendErrResponse(conn *websocket.Conn, msg string) {
	res := &ErrorResponseJSON{
		Jsonrpc: "2.0",
		Error: &ErrorMessageJSON{
			Code:    big.NewInt(-32600),
			Message: msg,
		},
		ID: big.NewInt(1),
	}
	err := conn.WriteJSON(res)
	if err != nil {
		s.logger.Error("websocket failed write message", "error", err)
	}
}

func (s *Server) readLoop(wsConn *websocket.Conn) {
	subIds := make(map[rpc.ID]struct{})
	for {
		_, mb, err := wsConn.ReadMessage()
		if err != nil {
			_ = wsConn.Close()
			s.logger.Error("failed to read message, close the websocket connection.", "error", err)
			s.closeWsConnection(subIds)
			return
		}

		var msg map[string]interface{}
		err = json.Unmarshal(mb, &msg)
		if err != nil {
			s.sendErrResponse(wsConn, "invalid request")
			continue
		}

		// check if method == eth_subscribe or eth_unsubscribe
		method := msg["method"]
		if method.(string) == "eth_subscribe" {
			params := msg["params"].([]interface{})
			if len(params) == 0 {
				s.sendErrResponse(wsConn, "invalid parameters")
				continue
			}

			reqId, ok := msg["id"].(float64)
			if !ok {
				s.sendErrResponse(wsConn, "invaild id in request message")
				continue
			}

			id, err := s.api.subscribe(wsConn, params)
			if err != nil {
				s.sendErrResponse(wsConn, err.Error())
				continue
			}

			res := &SubscriptionResponseJSON{
				Jsonrpc: "2.0",
				ID:      reqId,
				Result:  id,
			}

			err = wsConn.WriteJSON(res)
			if err != nil {
				s.logger.Error("failed to write json response", "ID", id, "error", err)
				continue
			}
			s.logger.Debug("successfully subscribe", "ID", id)
			subIds[id] = struct{}{}
			continue
		} else if method.(string) == "eth_unsubscribe" {
			ids, ok := msg["params"].([]interface{})
			if len(ids) == 0 {
				s.sendErrResponse(wsConn, "invalid parameters")
				continue
			}
			id, idok := ids[0].(string)
			if !ok || !idok {
				s.sendErrResponse(wsConn, "invalid parameters")
				continue
			}

			reqId, ok := msg["id"].(float64)
			if !ok {
				s.sendErrResponse(wsConn, "invaild id in request message")
				continue
			}

			ok = s.api.unsubscribe(rpc.ID(id))
			res := &SubscriptionResponseJSON{
				Jsonrpc: "2.0",
				ID:      reqId,
				Result:  ok,
			}

			err = wsConn.WriteJSON(res)
			if err != nil {
				s.logger.Error("failed to write json response", "ID", id, "error", err)
				continue
			}
			s.logger.Debug("successfully unsubscribe", "ID", id)
			delete(subIds, rpc.ID(id))
			continue
		}

		// otherwise, call the usual rpc server to respond
		err = s.tcpGetAndSendResponse(wsConn, mb)
		if err != nil {
			s.sendErrResponse(wsConn, err.Error())
		}
	}
}

// tcpGetAndSendResponse connects to the rest-server over tcp, posts a JSON-RPC request, and sends the response
// to the client over websockets
func (s *Server) tcpGetAndSendResponse(conn *websocket.Conn, mb []byte) error {
	req, err := http.NewRequest(http.MethodPost, s.rpcAddr, bytes.NewReader(mb))
	if err != nil {
		return fmt.Errorf("failed to request; %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to write to rest-server; %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read body from response; %s", err)
	}

	var wsSend interface{}
	err = json.Unmarshal(body, &wsSend)
	if err != nil {
		return fmt.Errorf("failed to unmarshal rest-server response; %s", err)
	}

	return conn.WriteJSON(wsSend)
}

func (s *Server) closeWsConnection(subIds map[rpc.ID]struct{}) {
	for id := range subIds {
		s.api.unsubscribe(id)
		delete(subIds, id)
	}
	s.connPoolLock.Lock()
	defer s.connPoolLock.Unlock()
	<-s.connPool
	s.currentConnNum.Set(float64(len(s.connPool)))
}
