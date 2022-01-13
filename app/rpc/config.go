package rpc

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/okex/exchain/app/crypto/ethsecp256k1"
	"github.com/okex/exchain/app/crypto/hd"
	"github.com/okex/exchain/app/rpc/nacos"
	"github.com/okex/exchain/app/rpc/pendingtx"
	"github.com/okex/exchain/app/rpc/websockets"
	"github.com/okex/exchain/libs/cosmos-sdk/client/flags"
	"github.com/okex/exchain/libs/cosmos-sdk/client/input"
	"github.com/okex/exchain/libs/cosmos-sdk/client/lcd"
	"github.com/okex/exchain/libs/cosmos-sdk/crypto/keys"
	cmserver "github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/spf13/viper"
)

const (
	flagUnlockKey = "unlock-key"
	flagWebsocket = "wsport"

	FlagPersonalAPI           = "personal-api"
	FlagRateLimitAPI          = "rpc.rate-limit-api"
	FlagRateLimitCount        = "rpc.rate-limit-count"
	FlagRateLimitBurst        = "rpc.rate-limit-burst"
	FlagEnableMonitor         = "rpc.enable-monitor"
	FlagDisableAPI            = "rpc.disable-api"
	FlagKafkaAddr             = "pendingtx.kafka-addr"
	FlagKafkaTopic            = "pendingtx.kafka-topic"
	FlagNacosTmrpcUrls        = "rpc.tmrpc_nacos_urls"
	FlagNacosTmrpcNamespaceID = "rpc.tmrpc_nacos_namespace_id"
	FlagNacosTmrpcAppName     = "rpc.tmrpc_application_name"
	FlagRpcExternalAddr       = "rpc.external_laddr"
	FlagRestApplicationName   = "rest.application_name"
	FlagRestNacosUrls         = "rest.nacos_urls"
	FlagRestNacosNamespaceId  = "rest.nacos_namespace_id"
	FlagExternalListenAddr    = "rest.external_laddr"

	MetricsNamespace = "x"
	// MetricsSubsystem is a subsystem shared by all metrics exposed by this package.
	MetricsSubsystem = "rpc"

	MetricsFieldName = "Metrics"
)

// RegisterRoutes creates a new server and registers the `/rpc` endpoint.
// Rpc calls are enabled based on their associated module (eg. "eth").
func RegisterRoutes(rs *lcd.RestServer) {
	// register nacos first
	registerNacos(rs.Logger())

	server := rpc.NewServer()
	accountName := viper.GetString(cmserver.FlagUlockKey)
	accountNames := strings.Split(accountName, ",")

	var privkeys []ethsecp256k1.PrivKey
	if len(accountName) > 0 {
		var err error
		inBuf := bufio.NewReader(os.Stdin)

		keyringBackend := viper.GetString(flags.FlagKeyringBackend)
		passphrase := ""
		switch keyringBackend {
		case keys.BackendOS:
			break
		case keys.BackendFile:
			passphrase, err = input.GetPassword(
				"Enter password to unlock key for RPC API: ",
				inBuf)
			if err != nil {
				panic(err)
			}
		}

		privkeys, err = unlockKeyFromNameAndPassphrase(accountNames, passphrase)
		if err != nil {
			panic(err)
		}
	}

	apis := GetAPIs(rs.CliCtx, rs.Logger(), privkeys...)

	// Register all the APIs exposed by the namespace services
	// TODO: handle allowlist and private APIs
	for _, api := range apis {
		if err := server.RegisterName(api.Namespace, api.Service); err != nil {
			panic(err)
		}
	}

	// Web3 RPC API route
	rs.Mux.HandleFunc("/", server.ServeHTTP).Methods("POST", "OPTIONS")

	// start websockets server
	websocketAddr := viper.GetString(flagWebsocket)
	ws := websockets.NewServer(rs.CliCtx, rs.Logger(), websocketAddr)
	ws.Start()

	// pending tx watcher
	kafkaAddrs := viper.GetString(FlagKafkaAddr)
	kafkaTopic := viper.GetString(FlagKafkaTopic)
	if kafkaAddrs != "" && kafkaTopic != "" {
		kafkaClient := pendingtx.NewKafkaClient(strings.Split(kafkaAddrs, ","), kafkaTopic)
		ptw := pendingtx.NewWatcher(rs.CliCtx, rs.Logger(), kafkaClient)
		ptw.Start()
	}
}

func unlockKeyFromNameAndPassphrase(accountNames []string, passphrase string) ([]ethsecp256k1.PrivKey, error) {
	keybase, err := keys.NewKeyring(
		sdk.KeyringServiceName(),
		viper.GetString(flags.FlagKeyringBackend),
		viper.GetString(cmserver.FlagUlockKeyHome),
		os.Stdin,
		hd.EthSecp256k1Options()...,
	)
	if err != nil {
		return []ethsecp256k1.PrivKey{}, err
	}

	// try the for loop with array []string accountNames
	// run through the bottom code inside the for loop

	keys := make([]ethsecp256k1.PrivKey, len(accountNames))
	for i, acc := range accountNames {
		// With keyring keybase, password is not required as it is pulled from the OS prompt
		privKey, err := keybase.ExportPrivateKeyObject(acc, passphrase)
		if err != nil {
			return []ethsecp256k1.PrivKey{}, err
		}

		var ok bool
		keys[i], ok = privKey.(ethsecp256k1.PrivKey)
		if !ok {
			panic(fmt.Sprintf("invalid private key type %T at index %d", privKey, i))
		}
	}

	return keys, nil
}

func registerNacos(logger log.Logger) {
	nacosUrls := viper.GetString(FlagRestNacosUrls)
	nacosNamespaceId := viper.GetString(FlagRestNacosNamespaceId)
	applicationName := viper.GetString(FlagRestApplicationName)
	externalAddr := viper.GetString(FlagExternalListenAddr)

	// start nacos client for registering restful service
	if nacosUrls != "" {
		nacos.StartNacosClient(logger, nacosUrls, nacosNamespaceId, applicationName, externalAddr)
	}

	nacosTmRpcUrls := viper.GetString(FlagNacosTmrpcUrls)
	nacosTmRpcNamespaceID := viper.GetString(FlagNacosTmrpcNamespaceID)
	nacosTmRpcAppName := viper.GetString(FlagNacosTmrpcAppName)
	rpcExternalAddr := viper.GetString(FlagRpcExternalAddr)
	// start nacos client for tmrpc service
	if nacosTmRpcUrls != "" {
		nacos.StartNacosClient(logger, nacosTmRpcUrls, nacosTmRpcNamespaceID, nacosTmRpcAppName, rpcExternalAddr)
	}
}
