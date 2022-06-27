package utils

import (
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	clientCtx "github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/server"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/types/rest"
	authUtils "github.com/okex/exchain/libs/cosmos-sdk/x/auth/client/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"net/http"
)

var (
	ClientContextKey = sdk.ContextKey("client.context")
	ServerContextKey = sdk.ContextKey("server.context")
)

func GetClientContextFromCmd(cmd *cobra.Command) clientCtx.CLIContext {
	// TODO need cdc like context.NewCLIContext().WithCodec(cdc)
	if v := cmd.Context().Value(ClientContextKey); v != nil {
		clientCtxPtr := v.(*context.CLIContext)
		return *clientCtxPtr
	}
	return context.NewCLIContext()
}

func GetServerContextFromCmd(cmd *cobra.Command) server.Context {
	// TODO need real server.context
	if v := cmd.Context().Value(ServerContextKey); v != nil {
		clientCtxPtr := v.(*server.Context)
		return *clientCtxPtr
	}
	ctx := server.NewDefaultContext()
	return *ctx
}

func GetClientQueryContext(cmd *cobra.Command) (clientCtx.CLIContext, error) {
	// TODO need cdc like context.NewCLIContext().WithCodec(cdc)
	if v := cmd.Context().Value(ClientContextKey); v != nil {
		clientCtxPtr := v.(*context.CLIContext)
		return *clientCtxPtr, nil
	}
	return context.NewCLIContext(), nil
}

func GetClientTxContext(cmd *cobra.Command) (clientCtx.CLIContext, error) {
	// TODO need cdc like context.NewCLIContext().WithCodec(cdc)
	if v := cmd.Context().Value(ClientContextKey); v != nil {
		clientCtxPtr := v.(*context.CLIContext)
		return *clientCtxPtr, nil
	}
	return context.NewCLIContext(), nil
}

func GenerateOrBroadcastTxCLI(cliCtx clientCtx.CLIContext, flagSet *pflag.FlagSet, msgs ...sdk.Msg) error {
	//TODO need cmd and cdc
	//inBuf := bufio.NewReader(cmd.InOrStdin())
	//txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(authUtils.GetTxEncoder(cdc))
	//if cliCtx.GenerateOnly {
	//	return authUtils.PrintUnsignedStdTx(txBldr, cliCtx, msgs)
	//}
	//
	//return authUtils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
	return nil
}

func WriteGeneratedTxResponse(
	clientCtx clientCtx.CLIContext, w http.ResponseWriter, br rest.BaseReq, msgs ...sdk.Msg,
) {
	authUtils.WriteGenerateStdTxResponse(w, clientCtx, br, msgs)
}
