package cli

import (
	"fmt"
	"strings"

	"github.com/okex/exchain/libs/cosmos-sdk/client/context"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	"github.com/okex/exchain/x/params/types"
	"github.com/spf13/cobra"
)

func GetCmdQueryUpgrade(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade [name]",
		Args:  cobra.ExactArgs(1),
		Short: "Query parameters of a upgrade",
		Long: strings.TrimSpace(`Query parameters of upgrade:

$ exchaincli query params upgrade <upgrade-name>
`),
		RunE: func(_ *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryUpgrade, args[0])
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var infos types.UpgradeInfo
			cdc.MustUnmarshalJSON(bz, &infos)
			return cliCtx.PrintOutput(infos)
		},
	}
}
