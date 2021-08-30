package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/okex/exchain/x/staking/types"
)

// GetQueryCmd returns the cli query commands for staking module
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	stakingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the staking module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	stakingQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryDelegator(queryRoute, cdc),
		GetCmdQueryValidatorShares(queryRoute, cdc),
		GetCmdQueryValidator(queryRoute, cdc),
		GetCmdQueryValidators(queryRoute, cdc),
		GetCmdQueryProxy(queryRoute, cdc),
		GetCmdQueryParams(queryRoute, cdc),
		GetCmdQueryPool(queryRoute, cdc))...)

	return stakingQueryCmd

}

// GetCmdQueryValidator gets the validator query command.
func GetCmdQueryValidator(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "validator [validator-addr]",
		Short: "query a validator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about an individual validator.

Example:
$ %s query staking validator exvaloper1alq9na49n9yycysh889rl90g9nhe58lcqkfpfg
`,
				version.ClientName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.ValAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryStore(types.GetValidatorKey(addr), storeName)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("no validator found with address %s", addr)
			}

			//return cliCtx.PrintOutput(types.MustUnmarshalValidator(cdc, res))
			return cliCtx.PrintOutput(types.MustUnmarshalValidator(cdc, res).Standardize())
		},
	}
}

// GetCmdQueryValidators gets the query all validators command.
func GetCmdQueryValidators(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "validators",
		Short: "query for all validators",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details about all validators on a network.

Example:
$ %s query staking validators
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			resKVs, _, err := cliCtx.QuerySubspace(types.ValidatorsKey, storeName)
			if err != nil {
				return err
			}

			var validators types.Validators
			for _, kv := range resKVs {
				validators = append(validators, types.MustUnmarshalValidator(cdc, kv.Value))
			}

			//return cliCtx.PrintOutput(validators)
			return cliCtx.PrintOutput(validators.Standardize())
		},
	}
}

// GetCmdQueryPool gets the pool query command.
func GetCmdQueryPool(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "pool",
		Args:  cobra.NoArgs,
		Short: "query the current staking pool values",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values for amounts stored in the staking pool.

Example:
$ %s query staking pool
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			bz, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/pool", storeName), nil)
			if err != nil {
				return err
			}

			var pool types.Pool
			if err := cdc.UnmarshalJSON(bz, &pool); err != nil {
				return err
			}

			return cliCtx.PrintOutput(pool)
		},
	}
}

// GetCmdQueryParams gets the params query command.
func GetCmdQueryParams(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "query the current staking parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as staking parameters.

Example:
$ %s query staking params
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryParameters)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(bz, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}

// GetCmdQueryProxy gets command for querying the delegators by a specific proxy
func GetCmdQueryProxy(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proxy [address]",
		Short: "query the delegator addresses by a proxy",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the addresses of delegators by a specific proxy

Example:
$ %s query staking proxy ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			proxyAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid address：%s", args[0])
			}

			bytes, err := cdc.MarshalJSON(types.NewQueryDelegatorParams(proxyAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryProxy)
			resp, _, err := cliCtx.QueryWithData(route, bytes)
			if err != nil {
				return err
			}

			var delAddrs Delegators
			if err := cdc.UnmarshalJSON(resp, &delAddrs); err != nil {
				return err
			}

			return cliCtx.PrintOutput(delAddrs)
		},
	}
}

// Delegators is a type alias of sdk.AccAddress slice
type Delegators []sdk.AccAddress

// String returns a human readable string representation of Delegators
func (as Delegators) String() (strFormat string) {
	for _, a := range as {
		strFormat = fmt.Sprintf("%s%s\n", strFormat, a.String())
	}

	return strings.TrimSpace(strFormat)
}

// GetCmdQueryDelegator gets command for querying the info of delegator about delegation and shares
func GetCmdQueryDelegator(storeName string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "delegator [address]",
		Short: "query the information about a delegator",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the information of delegations and all shares recently added by a delegator

Example:
$ %s query staking delegator ex1cftp8q8g4aa65nw9s5trwexe77d9t6cr8ndu02
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			delAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("invalid address：%s", args[0])
			}

			delegator, undelegation := types.NewDelegator(delAddr), types.DefaultUndelegation()
			resp, _, err := cliCtx.QueryStore(types.GetDelegatorKey(delAddr), storeName)
			if err != nil {
				return err
			}
			if len(resp) != 0 {
				cdc.MustUnmarshalBinaryLengthPrefixed(resp, &delegator)
			}

			// query for the undelegation info
			bytes, err := cdc.MarshalJSON(types.NewQueryDelegatorParams(delAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", storeName, types.QueryUnbondingDelegation)
			res, _, err := cliCtx.QueryWithData(route, bytes)
			// if err!= nil , we treat it as there's no undelegation of the delegator
			if err == nil {
				if err := cdc.UnmarshalJSON(res, &undelegation); err != nil {
					return err
				}
			}

			return cliCtx.PrintOutput(convertToDelegatorResp(delegator, undelegation))
		},
	}
}

// DelegatorResponse is designed for delegator info query
type DelegatorResponse struct {
	DelegatorAddress     sdk.AccAddress   `json:"delegator_address" yaml:"delegator_address"`
	ValidatorAddresses   []sdk.ValAddress `json:"validator_address" yaml:"validator_address"`
	Shares               sdk.Dec          `json:"shares" yaml:"shares"`
	Tokens               sdk.Dec          `json:"tokens" yaml:"tokens"`
	UnbondedTokens       sdk.Dec          `json:"unbonded_tokens" yaml:"unbonded_tokens"`
	CompletionTime       time.Time        `json:"completion_time" yaml:"completion_time"`
	IsProxy              bool             `json:"is_proxy" yaml:"is_proxy"`
	TotalDelegatedTokens sdk.Dec          `json:"total_delegated_tokens" yaml:"total_delegated_tokens"`
	ProxyAddress         sdk.AccAddress   `json:"proxy_address" yaml:"proxy_address"`
}

// String returns a human readable string representation of DelegatorResponse
func (dr DelegatorResponse) String() (output string) {
	n := len(dr.ValidatorAddresses)
	if n > 0 {
		output = fmt.Sprintf("%s\n", dr.ValidatorAddresses[0].String())
		for i := 1; i < n; i++ {
			output = fmt.Sprintf("%s						%s\n", output, dr.ValidatorAddresses[i].String())
		}
	}

	proxy := "No"
	if dr.IsProxy {
		proxy = "Yes"
	}

	proxied := "No"
	if dr.ProxyAddress != nil {
		proxied = "Yes\n	Proxied by " + dr.ProxyAddress.String() + "\n"
	}

	output = fmt.Sprintf(`Delegator:
	DelegatorAddress: 		%s
	ValidatorAddresses:		%s	
	Shares:					%s
	Tokens:					%s
	UnbondedTokens: 		%s
	CompletionTime:			%s
	IsProxied:				%s
	IsProxy:				%s`,
		dr.DelegatorAddress, output, dr.Shares, dr.Tokens, dr.UnbondedTokens, dr.CompletionTime, proxied, proxy)

	return
}

func convertToDelegatorResp(delegator types.Delegator, undelegation types.UndelegationInfo,
) DelegatorResponse {
	return DelegatorResponse{
		delegator.DelegatorAddress,
		delegator.ValidatorAddresses,
		delegator.Shares,
		delegator.Tokens,
		undelegation.Quantity,
		undelegation.CompletionTime,
		delegator.IsProxy,
		delegator.TotalDelegatedTokens,
		delegator.ProxyAddress,
	}
}

// GetCmdQueryValidatorShares gets command for querying all shares added to a validator
func GetCmdQueryValidatorShares(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "shares-added-to [validator-addr]",
		Short: "query all shares added to a validator",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query all shares added to a validator.

Example:
$ %s query staking shares-added-to exvaloper1alq9na49n9yycysh889rl90g9nhe58lcqkfpfg
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

			bytes, err := cdc.MarshalJSON(types.NewQueryValidatorParams(valAddr))
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryValidatorAllShares)
			resp, _, err := cliCtx.QueryWithData(route, bytes)
			if err != nil {
				return err
			}

			var sharesResponses types.SharesResponses
			if err := cdc.UnmarshalJSON(resp, &sharesResponses); err != nil {
				return err
			}

			return cliCtx.PrintOutput(sharesResponses)
		},
	}
}
