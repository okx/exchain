package rpc

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"

	"github.com/okex/okexchain/app/crypto/ethsecp256k1"
	"github.com/okex/okexchain/app/rpc/backend"
	"github.com/okex/okexchain/app/rpc/namespaces/eth"
	"github.com/okex/okexchain/app/rpc/namespaces/eth/filters"
	"github.com/okex/okexchain/app/rpc/namespaces/net"
	"github.com/okex/okexchain/app/rpc/namespaces/personal"
	"github.com/okex/okexchain/app/rpc/namespaces/web3"
	rpctypes "github.com/okex/okexchain/app/rpc/types"
	"github.com/okex/okexchain/cmd/client"
)

// RPC namespaces and API version
const (
	Web3Namespace     = "web3"
	EthNamespace      = "eth"
	PersonalNamespace = "personal"
	NetNamespace      = "net"

	apiVersion = "1.0"
)

// GetAPIs returns the list of all APIs from the Ethereum namespaces
func GetAPIs(clientCtx context.CLIContext, keys ...ethsecp256k1.PrivKey) []rpc.API {
	nonceLock := new(rpctypes.AddrLocker)
	backend := backend.New(clientCtx)
	ethAPI := eth.NewAPI(clientCtx, backend, nonceLock, keys...)

	apis := []rpc.API{
		{
			Namespace: Web3Namespace,
			Version:   apiVersion,
			Service:   web3.NewAPI(),
			Public:    true,
		},
		{
			Namespace: EthNamespace,
			Version:   apiVersion,
			Service:   ethAPI,
			Public:    true,
		},
		{
			Namespace: EthNamespace,
			Version:   apiVersion,
			Service:   filters.NewAPI(clientCtx, backend),
			Public:    true,
		},
		{
			Namespace: NetNamespace,
			Version:   apiVersion,
			Service:   net.NewAPI(clientCtx),
			Public:    true,
		},
	}

	if viper.GetBool(client.FlagPersonalAPI) {
		apis = append(apis, rpc.API{
			Namespace: PersonalNamespace,
			Version:   apiVersion,
			Service:   personal.NewAPI(ethAPI),
			Public:    false,
		})
	}
	return apis
}
