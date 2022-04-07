package proxy

import (
	abci "github.com/okex/exchain/libs/tendermint/abci/types"
	"github.com/okex/exchain/libs/tendermint/types"
	"github.com/okex/exchain/libs/tendermint/version"
)

// RequestInfo contains all the information for sending
// the abci.RequestInfo message during handshake with the app.
// It contains only compile-time version information.
var RequestInfo = abci.RequestInfo{
	Version:      version.Version,
	BlockVersion: version.BlockProtocol.Uint64(),
	P2PVersion:   version.P2PProtocol.Uint64(),
}

var IBCRequestInfo = abci.RequestInfo{
	Version:      version.Version,
	BlockVersion: version.IBCBlockProtocol.Uint64(),
	P2PVersion:   version.IBCP2PProtocol.Uint64(),
}

func GetRequestInfo(h int64) abci.RequestInfo {
	if types.HigherThanVenus1(h) {
		return IBCRequestInfo
	}
	return RequestInfo
}
