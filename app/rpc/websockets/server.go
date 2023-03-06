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

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/server"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
	"github.com/okx/okbchain/x/common/monitor"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
)

const FlagSubscribeLimit = "ws.max-subscriptions"

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
	maxSubLimit    int
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
		maxSubLimit: viper.GetInt(FlagSubscribeLimit),
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

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.logger.Error("websocket upgrade failed", " error", err)
		return
	}

	s.connPool <- struct{}{}
	s.currentConnNum.Set(float64(len(s.connPool)))
	go s.readLoop(&wsConn{
		mux:  new(sync.Mutex),
		conn: conn,
	})
}

func (s *Server) sendErrResponse(conn *wsConn, msg string) {
	res := makeErrResponse(msg)
	err := conn.WriteJSON(res)
	if err != nil {
		s.logger.Error("websocket failed write message", "error", err)
	}
}

func makeErrResponse(errMsg string) *ErrorResponseJSON {
	return &ErrorResponseJSON{
		Jsonrpc: "2.0",
		Error: &ErrorMessageJSON{
			Code:    big.NewInt(-32600),
			Message: errMsg,
		},
		ID: big.NewInt(1),
	}
}

type wsConn struct {
	conn     *websocket.Conn
	mux      *sync.Mutex
	subCount int
}

func (w *wsConn) GetSubCount() int {
	return w.subCount
}

func (w *wsConn) AddSubCount(delta int) {
	w.subCount += delta
}

func (w *wsConn) WriteJSON(v interface{}) error {
	w.mux.Lock()
	defer w.mux.Unlock()

	return w.conn.WriteJSON(v)
}

func (w *wsConn) Close() error {
	w.mux.Lock()
	defer w.mux.Unlock()

	return w.conn.Close()
}

func (w *wsConn) ReadMessage() (messageType int, p []byte, err error) {
	// not protected by write mutex

	return w.conn.ReadMessage()
}

func (s *Server) readLoop(wsConn *wsConn) {
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
		if err = json.Unmarshal(mb, &msg); err != nil {
			if err = s.batchCall(mb, wsConn); err != nil {
				s.sendErrResponse(wsConn, "invalid request")
			}
			continue
		}

		// check if method == eth_subscribe or eth_unsubscribe
		method := msg["method"]
		methodStr, ok := method.(string)
		if !ok {
			s.sendErrResponse(wsConn, "invalid request")
		}
		if methodStr == "eth_subscribe" {
			if wsConn.GetSubCount() >= s.maxSubLimit {
				s.sendErrResponse(wsConn,
					fmt.Sprintf("subscription has reached the upper limit(%d)", s.maxSubLimit))
				continue
			}
			params, ok := msg["params"].([]interface{})
			if !ok || len(params) == 0 {
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
			wsConn.AddSubCount(1)
			continue
		} else if methodStr == "eth_unsubscribe" {
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
			wsConn.AddSubCount(-1)
			continue
		}

		// otherwise, call the usual rpc server to respond
		data, err := s.getRpcResponse(mb)
		if err != nil {
			s.sendErrResponse(wsConn, err.Error())
		} else {
			wsConn.WriteJSON(data)
		}
	}
}

// getRpcResponse connects to the rest-server over tcp, posts a JSON-RPC request, and return response
func (s *Server) getRpcResponse(mb []byte) (interface{}, error) {
	req, err := http.NewRequest(http.MethodPost, s.rpcAddr, bytes.NewReader(mb))
	if err != nil {
		return nil, fmt.Errorf("failed to request; %s", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to write to rest-server; %s", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body from response; %s", err)
	}

	var wsSend interface{}
	err = json.Unmarshal(body, &wsSend)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal rest-server response; %s", err)
	}
	return wsSend, nil
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

func (s *Server) batchCall(mb []byte, wsConn *wsConn) error {
	var msgs []interface{}
	if err := json.Unmarshal(mb, &msgs); err != nil {
		return err
	}

	for i := 0; i < len(msgs); i++ {
		b, err := json.Marshal(msgs[i])
		if err != nil {
			s.sendErrResponse(wsConn, "invalid request")
			s.logger.Error("web socket batchCall  failed", "error", err)
			break
		}

		data, err := s.getRpcResponse(b)
		if err != nil {
			data = makeErrResponse(err.Error())
		}
		if err := wsConn.WriteJSON(data); err != nil {
			break // connection broken
		}
	}
	return nil
}
