package rpc

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/okex/exchain/app/rpc/namespaces/eth/txpool"
	evmtypes "github.com/okex/exchain/x/evm/types"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	"golang.org/x/time/rate"

	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/rpc/backend"
	"github.com/okex/exchain/app/rpc/monitor"
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
	TxpoolNamespace   = "txpool"

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

	if viper.GetBool(FlagEnableMonitor) {
		for _, api := range apis {
			makeMonitorMetrics(api.Namespace, api.Service)
		}
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

func makeMonitorMetrics(namespace string, service interface{}) {
	receiver := reflect.ValueOf(service)
	if !hasMetricsField(receiver.Elem()) {
		return
	}
	metricsVal := receiver.Elem().FieldByName(MetricsFieldName)

	monitorMetrics := make(map[string]*monitor.RpcMetrics)
	typ := receiver.Type()
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		if method.PkgPath != "" {
			continue // method not exported
		}
		methodName := formatMethodName(method.Name)
		name := fmt.Sprintf("%s_%s", namespace, methodName)
		monitorMetrics[name] = &monitor.RpcMetrics{
			Counter: prometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: MetricsNamespace,
				Subsystem: MetricsSubsystem,
				Name:      fmt.Sprintf("%s_count", name),
				Help:      fmt.Sprintf("Total request number of %s method.", name),
			}, nil),
			Histogram: prometheus.NewHistogramFrom(stdprometheus.HistogramOpts{
				Namespace: MetricsNamespace,
				Subsystem: MetricsSubsystem,
				Name:      fmt.Sprintf("%s_duration", name),
				Help:      fmt.Sprintf("Request duration of %s method.", name),
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .3, .5, 1, 3, 5, 10},
			}, nil),
		}

	}

	if metricsVal.CanSet() && metricsVal.Type() == reflect.ValueOf(monitorMetrics).Type() {
		metricsVal.Set(reflect.ValueOf(monitorMetrics))
	}
}

// formatMethodName converts to first character of name to lowercase.
func formatMethodName(name string) string {
	ret := []rune(name)
	if len(ret) > 0 {
		ret[0] = unicode.ToLower(ret[0])
	}
	return string(ret)
}

func hasMetricsField(receiver reflect.Value) bool {
	if receiver.Kind() != reflect.Struct {
		return false
	}
	for i := 0; i < receiver.NumField(); i++ {
		if receiver.Type().Field(i).Name == MetricsFieldName {
			return true
		}
	}
	return false
}
