package p2p

import (
	tmp2p "github.com/okex/exchain/libs/tendermint/proto/p2p"
)

func (info DefaultNodeInfo) ToProto() *tmp2p.DefaultNodeInfo {
	dni := new(tmp2p.DefaultNodeInfo)
	dni.ProtocolVersion = tmp2p.ProtocolVersion{
		P2P:   uint64(info.ProtocolVersion.P2P),
		Block: uint64(info.ProtocolVersion.Block),
		App:   uint64(info.ProtocolVersion.App),
	}

	dni.DefaultNodeID = string(info.DefaultNodeID)
	dni.ListenAddr = info.ListenAddr
	dni.Network = info.Network
	dni.Version = info.Version
	dni.Channels = info.Channels
	dni.Moniker = info.Moniker
	dni.Other = tmp2p.DefaultNodeInfoOther{
		TxIndex:     info.Other.TxIndex,
		RPCAdddress: info.Other.RPCAddress,
	}

	return dni
}
