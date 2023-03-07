package cli

import (
	"fmt"
	"strings"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	"github.com/okx/okbchain/x/params/types"
	"github.com/spf13/cobra"
)

func GetCmdQueryUpgrade(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "upgrade [name]",
		Args:  cobra.MinimumNArgs(0),
		Short: "Query info of upgrade",
		Long: strings.TrimSpace(`Query info of a upgrade, query all upgrade if 'name' is omitted:

$ okbchaincli query params upgrade <name>
`),
		RunE: func(_ *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			upgradeName := ""
			if len(args) > 0 {
				upgradeName = args[0]
			}

			route := fmt.Sprintf("custom/%s/%s/%s", queryRoute, types.QueryUpgrade, upgradeName)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var infos []types.UpgradeInfo
			cdc.MustUnmarshalJSON(bz, &infos)

			if len(upgradeName) != 0 {
				return cliCtx.PrintOutput(infos[0])
			}

			if len(infos) == 0 {
				return cliCtx.PrintOutput("there's no upgrade")
			}
			return cliCtx.PrintOutput(infos)
		},
	}
}
