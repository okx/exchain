package mix

import (
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	//"github.com/okex/exchain/app"
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
)

func SSS(ctx clientCtx.CLIContext, rtr *runtime.ServeMux) {
	//app.ModuleBasics.RegisterGRPCGatewayRoutes(ctx, rtr)
}
