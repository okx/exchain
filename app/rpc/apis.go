package rpc

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/ethereum/go-ethereum/rpc"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"golang.org/x/time/rate"
	"strings"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/rpc/backend"
	"github.com/okex/exchain/app/rpc/namespaces/eth"
	"github.com/okex/exchain/app/rpc/namespaces/eth/filters"
	"github.com/okex/exchain/app/rpc/namespaces/net"
	"github.com/okex/exchain/app/rpc/namespaces/personal"
	"github.com/okex/exchain/app/rpc/namespaces/web3"
	rpctypes "github.com/okex/exchain/app/rpc/types"
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
func GetAPIs(clientCtx context.CLIContext, log log.Logger, keys ...ethsecp256k1.PrivKey) []rpc.API {
	nonceLock := new(rpctypes.AddrLocker)
	rateLimiters := getRateLimiter()
	ethBackend := backend.New(clientCtx, log, rateLimiters)
	ethAPI := eth.NewAPI(clientCtx, log, ethBackend, nonceLock, keys...)
	if evmtypes.GetEnableBloomFilter() {
		server.TrapSignal(func() {
			if ethBackend != nil {
				ethBackend.Close()
			}
		})
		ethBackend.StartBloomHandlers(evmtypes.BloomBitsBlocks, evmtypes.GetIndexer().GetDB())
	}

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
			Service:   filters.NewAPI(clientCtx, ethBackend),
			Public:    true,
		},
		{
			Namespace: NetNamespace,
			Version:   apiVersion,
			Service:   net.NewAPI(clientCtx),
			Public:    true,
		},
	}

	if viper.GetBool(FlagPersonalAPI) {
		apis = append(apis, rpc.API{
			Namespace: PersonalNamespace,
			Version:   apiVersion,
			Service:   personal.NewAPI(ethAPI, log),
			Public:    false,
		})
	}
	return apis
}

func getRateLimiter() map[string]*rate.Limiter {
	rateLimitApi := viper.GetString(FlagRateLimitApi)
	rateLimitCount := viper.GetInt(FlagRateLimitCount)
	rateLimitBurst := viper.GetInt(FlagRateLimitBurst)
	if rateLimitApi == "" || rateLimitCount == 0 {
		return nil
	}
	rateLimiters := make(map[string]*rate.Limiter)
	apis := strings.Split(rateLimitApi, ",")
	for _, api := range apis {
		rateLimiters[api] = rate.NewLimiter(rate.Limit(rateLimitCount), rateLimitBurst)
	}
	return rateLimiters
}
