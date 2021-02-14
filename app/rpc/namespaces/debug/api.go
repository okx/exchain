package debug

import (
	"github.com/okex/okexchain/app/rpc/namespaces/eth"
)

// PrivateDebugAPI is the collection of Ethereum full node APIs exposed over
// the private debugging endpoint.
type PrivateDebugAPI struct {
	eth   *eth.PublicEthereumAPI
}

// NewPrivateDebugAPI creates a new API definition for the full node-related
// private debug methods of the Ethereum service.
func NewPrivateDebugAPI(eth *eth.PublicEthereumAPI) *PrivateDebugAPI {
	return &PrivateDebugAPI{eth: eth}
}
