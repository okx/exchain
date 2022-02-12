package client

import (
	"bufio"
	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/spf13/cobra"
)

func GetContext(cmd *cobra.Command,cdc *codec.Codec)context.CLIContext{
	inBuf := bufio.NewReader(cmd.InOrStdin())
	clientCtx:=context.NewCLIContextWithInput(inBuf).WithCodec(cdc)
	return clientCtx
}
