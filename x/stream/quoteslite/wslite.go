package quoteslite

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/lite/proxy"
	"github.com/tendermint/tendermint/rpc/core"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	rpcserver "github.com/tendermint/tendermint/rpc/lib/server"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	cmn "github.com/tendermint/tendermint/libs/common"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

// LiteCmd represents the base command when called without any subcommands
var LiteCmd = &cobra.Command{
	Use:   "lite",
	Short: "Run lite-client proxy server, verifying tendermint rpc",
	Long: `This node will run a secure proxy to a tendermint rpc server.

All calls that can be tracked back to a block header by a proof
will be verified before passing them back to the caller. Other that
that it will present the same interface as a full tendermint node,
just with added trust and running locally.`,
	RunE:         runProxy,
	SilenceUsage: true,
}

var (
	listenAddr         string
	nodeAddr           string
	chainID            string
	home               string
	maxOpenConnections int
	cacheSize          int

	opendexWSAddr string // websocket for opendex desktop which
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

func init() {
	LiteCmd.Flags().StringVar(&listenAddr, "laddr", "tcp://localhost:8888", "Serve the proxy on the given address")
	LiteCmd.Flags().StringVar(&nodeAddr, "node", "tcp://localhost:26657", "Connect to a Tendermint node at this address")
	LiteCmd.Flags().StringVar(&chainID, "chain-id", "tendermint", "Specify the Tendermint chain ID")
	LiteCmd.Flags().StringVar(&home, "home-dir", ".tendermint-lite", "Specify the home directory")
	LiteCmd.Flags().IntVar(
		&maxOpenConnections,
		"max-open-connections",
		900,
		"Maximum number of simultaneous connections (including WebSocket).")
	LiteCmd.Flags().IntVar(&cacheSize, "cache-size", 10, "Specify the memory trust store cache size")

	LiteCmd.Flags().StringVar(&opendexWSAddr, "open-wss-addr", "tcp://localhost:6666", "Specify the address of opendex websocket")
}

func EnsureAddrHasSchemeOrDefaultToTCP(addr string) (string, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "tcp", "unix":
	case "":
		u.Scheme = "tcp"
	default:
		return "", fmt.Errorf("unknown scheme %q, use either tcp or unix", u.Scheme)
	}
	return u.String(), nil
}

func runProxy(cmd *cobra.Command, args []string) error {
	// Stop upon receiving SIGTERM or CTRL-C.
	cmn.TrapSignal(logger, func() {
	})

	go StartWSServer(logger, nodeAddr)

	nodeAddr, err := EnsureAddrHasSchemeOrDefaultToTCP(nodeAddr)
	if err != nil {
		return err
	}
	listenAddr, err := EnsureAddrHasSchemeOrDefaultToTCP(listenAddr)
	if err != nil {
		return err
	}

	// First, connect a client
	logger.Info("Connecting to source HTTP client...")
	node := rpcclient.NewHTTP(nodeAddr, "/websocket")

	logger.Info("Constructing Verifier...")
	cert, err := proxy.NewVerifier(chainID, home, node, logger, cacheSize)
	if err != nil {
		return errors.Wrap(err, "constructing Verifier")
	}
	cert.SetLogger(logger)
	sc := proxy.SecureClient(node, cert)

	logger.Info("Starting proxy...")
	err = StartProxy(sc, listenAddr, logger, maxOpenConnections)
	if err != nil {
		return errors.Wrap(err, "starting proxy")
	}

	// Run forever
	select {}
}

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
