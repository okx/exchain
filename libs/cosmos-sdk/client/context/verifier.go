package context

import (
	"path/filepath"

	"github.com/okex/exchain/libs/tendermint/libs/log"
	tmlite "github.com/okex/exchain/libs/tendermint/lite"
	tmliteproxy "github.com/okex/exchain/libs/tendermint/lite/proxy"
	rpchttp "github.com/okex/exchain/libs/tendermint/rpc/client/http"
	"github.com/pkg/errors"
)

const (
	verifierDir = ".lite_verifier"

	// DefaultVerifierCacheSize defines the default Tendermint cache size.
	DefaultVerifierCacheSize = 10
)

// CreateVerifier returns a Tendermint verifier from a CLIContext object and
// cache size. An error is returned if the CLIContext is missing required values
// or if the verifier could not be created. A CLIContext must at the very least
// have the chain ID and home directory set. If the CLIContext has TrustNode
// enabled, no verifier will be created.
func CreateVerifier(ctx CLIContext, cacheSize int) (tmlite.Verifier, error) {
	if ctx.TrustNode {
		return nil, nil
	}

	switch {
	case ctx.ChainID == "":
		return nil, errors.New("must provide a valid chain ID to create verifier")

	case ctx.HomeDir == "":
		return nil, errors.New("must provide a valid home directory to create verifier")

	case ctx.Client == nil && ctx.NodeURI == "":
		return nil, errors.New("must provide a valid RPC client or RPC URI to create verifier")
	}

	var err error

	// create an RPC client based off of the RPC URI if no RPC client exists
	client := ctx.Client
	if client == nil {
		client, err = rpchttp.New(ctx.NodeURI, "/websocket")
		if err != nil {
			return nil, err
		}
	}

	return tmliteproxy.NewVerifier(
		ctx.ChainID, filepath.Join(ctx.HomeDir, ctx.ChainID, verifierDir),
		client, log.NewNopLogger(), cacheSize,
	)
}
