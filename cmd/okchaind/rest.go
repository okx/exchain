package main

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/lcd"
	authrest "github.com/cosmos/cosmos-sdk/x/auth/client/rest"
	bankrest "github.com/cosmos/cosmos-sdk/x/bank/client/rest"
	supplyrest "github.com/cosmos/cosmos-sdk/x/supply/client/rest"
	backendrest "github.com/okex/okchain/x/backend/client/rest"
	dexrest "github.com/okex/okchain/x/dex/client/rest"
	dist "github.com/okex/okchain/x/distribution"
	distrest "github.com/okex/okchain/x/distribution/client/rest"
	orderrest "github.com/okex/okchain/x/order/client/rest"
	stakingrest "github.com/okex/okchain/x/staking/client/rest"
	"github.com/okex/okchain/x/token"
	tokensrest "github.com/okex/okchain/x/token/client/rest"
	wasmrest "github.com/okex/okchain/x/wasm/client/rest"
)

// registerRoutes registers the routes from the different modules for the LCD.
// NOTE: details on the routes added for each module are in the module documentation
// NOTE: If making updates here you also need to update the test helper in client/lcd/test_helper.go
func registerRoutes(rs *lcd.RestServer) {
	registerRoutesV1(rs)
	registerRoutesV2(rs)
}

func registerRoutesV1(rs *lcd.RestServer) {
	v1Router := rs.Mux.PathPrefix("/okchain/v1").Name("v1").Subrouter()
	client.RegisterRoutes(rs.CliCtx, v1Router)
	authrest.RegisterRoutes(rs.CliCtx, v1Router)
	bankrest.RegisterRoutes(rs.CliCtx, v1Router)
	stakingrest.RegisterRoutes(rs.CliCtx, v1Router)
	distrest.RegisterRoutes(rs.CliCtx, v1Router, dist.StoreKey)

	orderrest.RegisterRoutes(rs.CliCtx, v1Router)
	tokensrest.RegisterRoutes(rs.CliCtx, v1Router, token.ModuleName)
	backendrest.RegisterRoutes(rs.CliCtx, v1Router)
	dexrest.RegisterRoutes(rs.CliCtx, v1Router)
	supplyrest.RegisterRoutes(rs.CliCtx, v1Router)
	wasmrest.RegisterRoutes(rs.CliCtx, v1Router)
}

func registerRoutesV2(rs *lcd.RestServer) {
	v2Router := rs.Mux.PathPrefix("/okchain/v2").Name("v1").Subrouter()
	client.RegisterRoutes(rs.CliCtx, v2Router)
	authrest.RegisterRoutes(rs.CliCtx, v2Router)
	bankrest.RegisterRoutes(rs.CliCtx, v2Router)
	stakingrest.RegisterRoutes(rs.CliCtx, v2Router)
	distrest.RegisterRoutes(rs.CliCtx, v2Router, dist.StoreKey)

	orderrest.RegisterRoutesV2(rs.CliCtx, v2Router)
	tokensrest.RegisterRoutesV2(rs.CliCtx, v2Router, token.ModuleName)
	backendrest.RegisterRoutesV2(rs.CliCtx, v2Router)
}
