package net

import (
	"fmt"

	"github.com/okex/exchain/app/rpc/monitor"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/tendermint/libs/log"
)

// PublicNetAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicNetAPI struct {
	networkVersion uint64
	logger         log.Logger
	Metrics        map[string]*monitor.RpcMetrics
}

// NewAPI creates an instance of the public Net Web3 API.
func NewAPI(clientCtx context.CLIContext, log log.Logger) *PublicNetAPI {
	// parse the chainID from a integer string
	chainIDEpoch, err := ethermint.ParseChainID(clientCtx.ChainID)
	if err != nil {
		panic(err)
	}

	return &PublicNetAPI{
		networkVersion: chainIDEpoch.Uint64(),
		logger:         log.With("module", "json-rpc", "namespace", "net"),
	}
}

// Version returns the current ethereum protocol version.
func (api *PublicNetAPI) Version() string {
	monitor := monitor.GetMonitor("net_version", api.logger, api.Metrics).OnBegin()
	defer monitor.OnEnd()
	return fmt.Sprintf("%d", api.networkVersion)
}
