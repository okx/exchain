package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	sdkGovCli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	gcutils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	"github.com/spf13/cobra"

	"github.com/okex/okexchain/x/gov/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	// Group gov queries under a subcommand
	govQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the governance module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	govQueryCmd.AddCommand(client.GetCommands(
		sdkGovCli.GetCmdQueryProposal(queryRoute, cdc),
		sdkGovCli.GetCmdQueryProposals(queryRoute, cdc),
		getCmdQueryVote(queryRoute, cdc),
		getCmdQueryVotes(queryRoute, cdc),
		sdkGovCli.GetCmdQueryParam(queryRoute, cdc),
		sdkGovCli.GetCmdQueryParams(queryRoute, cdc),
		sdkGovCli.GetCmdQueryProposer(queryRoute, cdc),
		getCmdQueryDeposit(queryRoute, cdc),
		getCmdQueryDeposits(queryRoute, cdc),
		sdkGovCli.GetCmdQueryTally(queryRoute, cdc))...)

	return govQueryCmd
}

// Command to Get a Proposal Information
// getCmdQueryVote implements the query proposal vote command.
func getCmdQueryVote(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "vote [proposal-id] [voter-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single vote",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a single vote on a proposal given its identifier.

Example:
$ %s query gov vote 1 okexchain1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			voterAddr, proposalID, _, err := parse(cliCtx, queryRoute, args)
			if err != nil {
				return err
			}

			params := types.NewQueryVoteParams(proposalID, voterAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/vote", queryRoute), bz)
			if err != nil {
				return err
			}

			var vote types.Vote
			if err := cdc.UnmarshalJSON(res, &vote); err != nil {
				return err
			}

			if vote.Empty() {
				res, err = gcutils.QueryVoteByTxQuery(cliCtx, params)
				if err != nil {
					return err
				}
				if err := cdc.UnmarshalJSON(res, &vote); err != nil {
					return err
				}
			}
			return cliCtx.PrintOutput(vote) //nolint:errcheck
		},
	}
}

func getDepositsOrVotes(cdc *codec.Codec, queryRoute string, args []string, isDeposits bool) error {
	cliCtx := context.NewCLIContext().WithCodec(cdc)

	_, proposalID, res, err := parse(cliCtx, queryRoute, args)
	if err != nil {
		return err
	}
	params := types.NewQueryProposalParams(proposalID)
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return err
	}

	var proposal types.Proposal
	cdc.MustUnmarshalJSON(res, &proposal)

	propStatus := proposal.Status
	if isDeposits {
		if !(propStatus == types.StatusVotingPeriod || propStatus == types.StatusDepositPeriod) {
			res, err = gcutils.QueryDepositsByTxQuery(cliCtx, params)
		} else {
			res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/deposits", queryRoute), bz)
		}
	} else {
		if !(propStatus == types.StatusVotingPeriod || propStatus == types.StatusDepositPeriod) {
			res, err = gcutils.QueryVotesByTxQuery(cliCtx, params)
		} else {
			res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/votes", queryRoute), bz)
		}
	}

	if err != nil {
		return err
	}

	type DepositVotes = types.Votes
	if isDeposits {
		type DepositVotes = types.Deposits
	}
	var dep DepositVotes
	cdc.MustUnmarshalJSON(res, &dep)
	return cliCtx.PrintOutput(dep)
}

// getCmdQueryVotes implements the command to query for proposal votes.
func getCmdQueryVotes(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "votes [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query votes on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query vote details for a single proposal by its identifier.

Example:
$ %s query gov votes 1
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getDepositsOrVotes(cdc, queryRoute, args, false)
		},
	}
}

// Command to Get a specific deposit Information
// getCmdQueryDeposit implements the query proposal deposit command.
func getCmdQueryDeposit(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposit [proposal-id] [depositer-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a deposit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a single proposal deposit on a proposal by its identifier.

Example:
$ %s query gov deposit 1 okexchain1skjwj5whet0lpe65qaq4rpq03hjxlwd9nf39lk
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			depositorAddr, proposalID, _, err := parse(cliCtx, queryRoute, args)
			if err != nil {
				return err
			}

			params := types.NewQueryDepositParams(proposalID, depositorAddr)
			bz, err := cdc.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/deposit", queryRoute), bz)
			if err != nil {
				return err
			}

			var deposit types.Deposit
			cdc.MustUnmarshalJSON(res, &deposit)

			if deposit.Empty() {
				res, err = gcutils.QueryDepositByTxQuery(cliCtx, params)
				if err != nil {
					return err
				}
				cdc.MustUnmarshalJSON(res, &deposit)
			}

			return cliCtx.PrintOutput(deposit)
		},
	}
}

// getCmdQueryDeposits implements the command to query for proposal deposits.
func getCmdQueryDeposits(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "deposits [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query deposits on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for all deposits on a proposal.
You can find the proposal-id by running "%s query gov proposals".

Example:
$ %s query gov deposits 1
`,
				version.ClientName, version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return getDepositsOrVotes(cdc, queryRoute, args, true)
		},
	}
}

func parse(cliCtx client.CLIContext, queryRoute string, args []string) (sdk.AccAddress, uint64, []byte, error) {
	// validate that the proposal id is a uint
	proposalID, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, 0, []byte{},
			fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
	}

	// check to see if the proposal is in the store
	res, err := gcutils.QueryProposalByID(proposalID, cliCtx, queryRoute)
	if err != nil {
		return nil, proposalID, []byte{}, fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
	}

	var addr sdk.AccAddress
	if len(args) > 1 {
		addr, err = sdk.AccAddressFromBech32(args[1])
		if err != nil {
			return addr, proposalID, res, fmt.Errorf("invalid address：%s", args[1])
		}
	}
	return addr, proposalID, res, nil
}

// DONTCOVER
