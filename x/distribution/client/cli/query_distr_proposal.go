package cli

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/okx/okbchain/libs/cosmos-sdk/client/context"
	"github.com/okx/okbchain/libs/cosmos-sdk/codec"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/version"

	"github.com/okx/okbchain/x/distribution/client/common"
	"github.com/okx/okbchain/x/distribution/types"
)

// GetCmdQueryDelegatorRewards implements the query delegator rewards command.
func GetCmdQueryDelegatorRewards(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "rewards [delegator-addr] [<validator-addr>]",
		Args:  cobra.RangeArgs(1, 2),
		Short: "Query all distribution delegator rewards or rewards from a particular validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all rewards earned by a delegator, optionally restrict to rewards from a single validator.

Example:
$ %s query distr rewards ex1j5mr2jhr9pf20e7yhln5zkcsgqtdt7cydr8x3y
$ %s query distr rewards ex1j5mr2jhr9pf20e7yhln5zkcsgqtdt7cydr8x3y exvaloper1pt7xrmxul7sx54ml44lvv403r06clrdkehd8z7
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// query for rewards from a particular delegation
			if len(args) == 2 {
				resp, _, err := common.QueryDelegationRewards(cliCtx, queryRoute, args[0], args[1])
				if err != nil {
					return err
				}

				var result sdk.DecCoins
				if err = cdc.UnmarshalJSON(resp, &result); err != nil {
					return fmt.Errorf("failed to unmarshal response: %w", err)
				}

				return cliCtx.PrintOutput(result)
			}

			delegatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := types.NewQueryDelegatorParams(delegatorAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return fmt.Errorf("failed to marshal params: %w", err)
			}

			// query for delegator total rewards
			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryDelegatorTotalRewards)
			res, _, err := cliCtx.QueryWithData(route, bz)
			if err != nil {
				return err
			}

			var result types.QueryDelegatorTotalRewardsResponse
			if err = cdc.UnmarshalJSON(res, &result); err != nil {
				return fmt.Errorf("failed to unmarshal response: %w", err)
			}

			return cliCtx.PrintOutput(result)
		},
	}
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryValidatorOutstandingRewards(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "outstanding-rewards [validator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query distribution outstanding (un-withdrawn) rewards
for a validator and all their delegations.

Example:
$ %s query distr outstanding-rewards exvaloper1pt7xrmxul7sx54ml44lvv403r06clrdkehd8z7
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			valAddr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := types.NewQueryValidatorOutstandingRewardsParams(valAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			resp, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryValidatorOutstandingRewards),
				bz,
			)
			if err != nil {
				return err
			}

			var outstandingRewards types.ValidatorOutstandingRewards
			if err := cdc.UnmarshalJSON(resp, &outstandingRewards); err != nil {
				return err
			}

			return cliCtx.PrintOutput(outstandingRewards)
		},
	}
}

// GetCmdQueryWithdrawAddr implements the query the delegator's withdraw address for commission and reward
func GetCmdQueryWithdrawAddr(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "withdraw-addr [delegator]",
		Args:  cobra.ExactArgs(1),
		Short: "Query delegator's withdraw address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query delegator's withdraw address.

Example:
$ %s query distr withdraw-addr ex17kn7d20d85yymu45h79dqs5pxq9m3nyx2mdmcs
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			delegatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			params := types.NewQueryDelegatorWithdrawAddrParams(delegatorAddr)

			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			resp, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryWithdrawAddr),
				bz,
			)
			if err != nil {
				return err
			}

			var accAddress sdk.AccAddress
			if err := cdc.UnmarshalJSON(resp, &accAddress); err != nil {
				return err
			}

			return cliCtx.PrintOutput(accAddress)
		},
	}
}
