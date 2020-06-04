package cli

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	flag "github.com/spf13/pflag"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okex/okchain/x/wasm/internal/keeper"
	"github.com/okex/okchain/x/wasm/internal/types"
)

func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the wasm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	queryCmd.AddCommand(flags.GetCommands(
		GetCmdListCode(cdc),
		GetCmdListContractByCode(cdc),
		GetCmdQueryCode(cdc),
		GetCmdGetContractInfo(cdc),
		GetCmdGetContractState(cdc),
	)...)
	return queryCmd
}

// GetCmdListCode lists all wasm code uploaded
func GetCmdListCode(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list-code",
		Short: "List all wasm bytecode on the chain",
		Long:  "List all wasm bytecode on the chain",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, keeper.QueryListCode)
			res, _, err := cliCtx.Query(route)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
}

// GetCmdListContractByCode lists all wasm code uploaded for given code id
func GetCmdListContractByCode(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "list-contract-by-code [code_id]",
		Short: "List wasm all bytecode on the chain for given code id",
		Long:  "List wasm all bytecode on the chain for given code id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s/%d", types.QuerierRoute, keeper.QueryListContractByCode, codeID)
			res, _, err := cliCtx.Query(route)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
}

// GetCmdQueryCode returns the bytecode for a given contract
func GetCmdQueryCode(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "code [code_id] [output filename]",
		Short: "Downloads wasm bytecode for given code id",
		Long:  "Downloads wasm bytecode for given code id",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			codeID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s/%d", types.QuerierRoute, keeper.QueryGetCode, codeID)
			res, _, err := cliCtx.Query(route)
			if err != nil {
				return err
			}
			if len(res) == 0 {
				return fmt.Errorf("contract not found")
			}
			var code keeper.GetCodeResponse
			err = json.Unmarshal(res, &code)
			if err != nil {
				return err
			}

			if len(code.Data) == 0 {
				return fmt.Errorf("contract not found")
			}

			fmt.Printf("Downloading wasm code to %s\n", args[1])
			return ioutil.WriteFile(args[1], code.Data, 0644)
		},
	}
}

// GetCmdGetContractInfo gets details about a given contract
func GetCmdGetContractInfo(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "contract [bech32_address]",
		Short: "Prints out metadata of a contract given its address",
		Long:  "Prints out metadata of a contract given its address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, keeper.QueryGetContract, addr.String())
			res, _, err := cliCtx.Query(route)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
}

// GetCmdGetContractState dumps full internal state of a given contract
func GetCmdGetContractState(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "contract-state",
		Short:                      "Querying commands for the wasm module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	cmd.AddCommand(flags.GetCommands(
		GetCmdGetContractStateAll(cdc),
		GetCmdGetContractStateRaw(cdc),
		GetCmdGetContractStateSmart(cdc),
	)...)
	return cmd

}

func GetCmdGetContractStateAll(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "all [bech32_address]",
		Short: "Prints out all internal state of a contract given its address",
		Long:  "Prints out all internal state of a contract given its address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, keeper.QueryGetContractState, addr.String(), keeper.QueryMethodContractStateAll)
			res, _, err := cliCtx.Query(route)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
}

func GetCmdGetContractStateRaw(cdc *codec.Codec) *cobra.Command {
	decoder := newArgDecoder(hex.DecodeString)
	cmd := &cobra.Command{
		Use:   "raw [bech32_address] [key]",
		Short: "Prints out internal state for key of a contract given its address",
		Long:  "Prints out internal state for of a contract given its address",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			queryData, err := decoder.DecodeString(args[1])
			if err != nil {
				return err
			}
			route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, keeper.QueryGetContractState, addr.String(), keeper.QueryMethodContractStateRaw)
			res, _, err := cliCtx.QueryWithData(route, queryData)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	decoder.RegisterFlags(cmd.PersistentFlags(), "key argument")
	return cmd
}

func GetCmdGetContractStateSmart(cdc *codec.Codec) *cobra.Command {
	decoder := newArgDecoder(asciiDecodeString)

	cmd := &cobra.Command{
		Use:   "smart [bech32_address] [query]",
		Short: "Calls contract with given address  with query data and prints the returned result",
		Long:  "Calls contract with given address  with query data and prints the returned result",
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			addr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			key := args[1]
			if key == "" {
				return errors.New("key must not be empty")
			}
			route := fmt.Sprintf("custom/%s/%s/%s/%s", types.QuerierRoute, keeper.QueryGetContractState, addr.String(), keeper.QueryMethodContractStateSmart)

			queryData, err := decoder.DecodeString(args[1])
			if err != nil {
				return fmt.Errorf("decode query: %s", err)
			}
			res, _, err := cliCtx.QueryWithData(route, queryData)
			if err != nil {
				return err
			}
			fmt.Println(string(res))
			return nil
		},
	}
	decoder.RegisterFlags(cmd.PersistentFlags(), "query argument")
	return cmd
}

type argumentDecoder struct {
	// dec is the default decoder
	dec                func(string) ([]byte, error)
	asciiF, hexF, b64F bool
}

func newArgDecoder(def func(string) ([]byte, error)) *argumentDecoder {
	return &argumentDecoder{dec: def}
}

func (a *argumentDecoder) RegisterFlags(f *flag.FlagSet, argName string) {
	f.BoolVar(&a.asciiF, "ascii", false, "ascii encoded "+argName)
	f.BoolVar(&a.hexF, "hex", false, "hex encoded  "+argName)
	f.BoolVar(&a.b64F, "b64", false, "base64 encoded "+argName)
}

func (a *argumentDecoder) DecodeString(s string) ([]byte, error) {
	found := -1
	for i, v := range []*bool{&a.asciiF, &a.hexF, &a.b64F} {
		if !*v {
			continue
		}
		if found != -1 {
			return nil, errors.New("multiple decoding flags used")
		}
		found = i
	}
	switch found {
	case 0:
		return asciiDecodeString(s)
	case 1:
		return hex.DecodeString(s)
	case 2:
		return base64.StdEncoding.DecodeString(s)
	default:
		return a.dec(s)
	}
}

func asciiDecodeString(s string) ([]byte, error) {
	return []byte(s), nil
}
