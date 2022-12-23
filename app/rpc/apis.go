package rpc

import (
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/okex/exchain/app/rpc/namespaces/eth/txpool"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/rpc/backend"
	"github.com/okex/exchain/app/rpc/namespaces/debug"
	"github.com/okex/exchain/app/rpc/namespaces/eth"
	"github.com/okex/exchain/app/rpc/namespaces/eth/filters"
	"github.com/okex/exchain/app/rpc/namespaces/net"
	"github.com/okex/exchain/app/rpc/namespaces/personal"
	"github.com/okex/exchain/app/rpc/namespaces/web3"
	rpctypes "github.com/okex/exchain/app/rpc/types"
	cosmost "github.com/okex/exchain/libs/cosmos-sdk/store/types"
)

// RPC namespaces and API version
const (
	Web3Namespace     = "web3"
	EthNamespace      = "eth"
	PersonalNamespace = "personal"
	NetNamespace      = "net"
	TxpoolNamespace   = "txpool"
	DebugNamespace    = "debug"

	apiVersion = "1.0"
)

var ethBackend *backend.EthermintBackend

func CloseEthBackend() {
	if ethBackend != nil {
		ethBackend.Close()
	}
}

// GetAPIs returns the list of all APIs from the Ethereum namespaces
func GetAPIs(clientCtx context.CLIContext, log log.Logger, keys ...ethsecp256k1.PrivKey) []rpc.API {
	nonceLock := new(rpctypes.AddrLocker)
	rateLimiters := getRateLimiter()
	disableAPI := getDisableAPI()
	ethBackend = backend.New(clientCtx, log, rateLimiters, disableAPI)
	ethAPI := eth.NewAPI(clientCtx, log, ethBackend, nonceLock, keys...)
	if evmtypes.GetEnableBloomFilter() {
		ethBackend.StartBloomHandlers(evmtypes.BloomBitsBlocks, evmtypes.GetIndexer().GetDB())
	}

	apis := []rpc.API{
		{
			Namespace: Web3Namespace,
			Version:   apiVersion,
			Service:   web3.NewAPI(log),
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
			Service:   filters.NewAPI(clientCtx, log, ethBackend),
			Public:    true,
		},
		{
			Namespace: NetNamespace,
			Version:   apiVersion,
			Service:   net.NewAPI(clientCtx, log),
			Public:    true,
		},
		{
			Namespace: TxpoolNamespace,
			Version:   apiVersion,
			Service:   txpool.NewAPI(clientCtx, log, ethBackend),
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

	if viper.GetBool(FlagDebugAPI) && viper.GetString(server.FlagPruning) == cosmost.PruningOptionNothing {
		apis = append(apis, rpc.API{
			Namespace: DebugNamespace,
			Version:   apiVersion,
			Service:   debug.NewAPI(clientCtx, log, ethBackend),
			Public:    true,
		})
	}

	return apis
}

func getRateLimiter() map[string]*rate.Limiter {
	rateLimitApi := viper.GetString(FlagRateLimitAPI)
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

func getDisableAPI() map[string]bool {
	disableAPI := viper.GetString(FlagDisableAPI)
	apiMap := make(map[string]bool)
	apis := strings.Split(disableAPI, ",")
	for _, api := range apis {
		apiMap[api] = true
	}
	return apiMap
}
