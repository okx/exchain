package main

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bankrest "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	supplyrest "github.com/cosmos/cosmos-sdk/x/supply/client/rest"
	"github.com/okex/okexchain/app/rpc"
	"github.com/okex/okexchain/app/types"
	ammswaprest "github.com/okex/okexchain/x/ammswap/client/rest"
	backendrest "github.com/okex/okexchain/x/backend/client/rest"
	dexrest "github.com/okex/okexchain/x/dex/client/rest"
	dist "github.com/okex/okexchain/x/distribution"
	distrest "github.com/okex/okexchain/x/distribution/client/rest"
	farmrest "github.com/okex/okexchain/x/farm/client/rest"
	orderrest "github.com/okex/okexchain/x/order/client/rest"
	stakingrest "github.com/okex/okexchain/x/staking/client/rest"
	"github.com/okex/okexchain/x/token"
	tokensrest "github.com/okex/okexchain/x/token/client/rest"
	"github.com/spf13/viper"
)

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {
	rpc.RegisterRoutes(rs)
	pathPrefix := viper.GetString(server.FlagRestPathPrefix)
	if pathPrefix == "" {
		pathPrefix = types.EthBech32Prefix
	}
	registerRoutesV1(rs, pathPrefix)
	registerRoutesV2(rs, pathPrefix)
}

func registerRoutesV1(rs *lcd.RestServer, pathPrefix string) {
	v1Router := rs.Mux.PathPrefix(fmt.Sprintf("/%s/v1", pathPrefix)).Name("v1").Subrouter()
	client.RegisterRoutes(rs.CliCtx, v1Router)
	authrest.RegisterRoutes(rs.CliCtx, v1Router, auth.StoreKey)
	authrest.RegisterTxRoutes(rs.CliCtx, v1Router)
	bankrest.RegisterRoutes(rs.CliCtx, v1Router)
	stakingrest.RegisterRoutes(rs.CliCtx, v1Router)
	distrest.RegisterRoutes(rs.CliCtx, v1Router, dist.StoreKey)

	orderrest.RegisterRoutes(rs.CliCtx, v1Router)
	tokensrest.RegisterRoutes(rs.CliCtx, v1Router, token.StoreKey)
	backendrest.RegisterRoutes(rs.CliCtx, v1Router)
	dexrest.RegisterRoutes(rs.CliCtx, v1Router)
	ammswaprest.RegisterRoutes(rs.CliCtx, v1Router)
	supplyrest.RegisterRoutes(rs.CliCtx, v1Router)
	farmrest.RegisterRoutes(rs.CliCtx, v1Router)
}

func registerRoutesV2(rs *lcd.RestServer, pathPrefix string) {
	v2Router := rs.Mux.PathPrefix(fmt.Sprintf("/%s/v2", pathPrefix)).Name("v1").Subrouter()
	client.RegisterRoutes(rs.CliCtx, v2Router)
	authrest.RegisterRoutes(rs.CliCtx, v2Router, auth.StoreKey)
	bankrest.RegisterRoutes(rs.CliCtx, v2Router)
	stakingrest.RegisterRoutes(rs.CliCtx, v2Router)
	distrest.RegisterRoutes(rs.CliCtx, v2Router, dist.StoreKey)

	orderrest.RegisterRoutesV2(rs.CliCtx, v2Router)
	tokensrest.RegisterRoutesV2(rs.CliCtx, v2Router, token.StoreKey)
	backendrest.RegisterRoutesV2(rs.CliCtx, v2Router)
}
