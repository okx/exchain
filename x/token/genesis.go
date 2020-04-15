package token

import (
	"encoding/json"
	"errors"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/okex/okchain/x/common"
	"github.com/okex/okchain/x/token/types"
)

// default owner of okt
const DefaultTokenOwner = "okchain10q0rk5qnyag7wfvvt7rtphlw589m7frsmyq4ya"

// all state that must be provided in genesis file
type GenesisState struct {
	Params    types.Params     `json:"params"`
	Tokens    []types.Token    `json:"tokens"`
	LockCoins []types.AccCoins `json:"locked_asset"`
}

// default GenesisState used by Cosmos Hub
func defaultGenesisState() GenesisState {
	return GenesisState{
		Params:    types.DefaultParams(),
		Tokens:    []types.Token{defaultGenesisStateOKT()},
		LockCoins: nil,
	}
}

// default okt information
func defaultGenesisStateOKT() types.Token {
	addr, err := sdk.AccAddressFromBech32(DefaultTokenOwner)
	if err != nil {
		panic(err)
	}

	totalSupply := sdk.NewDec(1000000000)
	return types.Token{
		Description:         "OK Group Global Utility Token",
		Symbol:              common.NativeToken,
		OriginalSymbol:      "OKT",
		WholeName:           "OKT",
		OriginalTotalSupply: totalSupply,
		TotalSupply:         totalSupply,
		Owner:               addr,
		Mintable:            true,
	}
}

func validateGenesis(data GenesisState) error {
	for _, token := range data.Tokens {
		msg := types.NewMsgTokenIssue(token.Description,
			token.Symbol,
			token.OriginalSymbol,
			token.WholeName,
			token.OriginalTotalSupply.String(),
			token.Owner,
			token.Mintable)

		err := msg.ValidateBasic()
		if err != nil {
			return errors.New(err.Error())
		}
	}
	return nil
}

// initGenesis initialize default parameters
// and the keeper's address to pubkey map
func initGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	keeper.SetParams(ctx, data.Params)

	for _, token := range data.Tokens {
		keeper.NewToken(ctx, token)
	}

	for _, lock := range data.LockCoins {
		if err := keeper.updateLockCoins(ctx, lock.Acc, lock.Coins, true); err != nil {
			panic(err)
		}
	}
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with initGenesis
func ExportGenesis(ctx sdk.Context, keeper Keeper) (data GenesisState) {
	params := keeper.GetParams(ctx)
	tokens := keeper.GetTokensInfo(ctx)
	locks := keeper.GetAllLockCoins(ctx)

	return GenesisState{
		Params:    params,
		Tokens:    tokens,
		LockCoins: locks,
	}
}

// IssueOKT issues okt in initchain
func IssueOKT(ctx sdk.Context, k Keeper, genesisState json.RawMessage, acc auth.Account) error {
	var data GenesisState
	types.ModuleCdc.MustUnmarshalJSON(genesisState, &data)
	for _, t := range data.Tokens {
		if t.Owner.Empty() && acc != nil {
			t.Owner = acc.GetAddress()
		}
		coins := k.GetCoins(ctx, t.Owner)
		if !strings.Contains(coins.String(), t.Symbol) {
			coins = append(coins, sdk.NewDecCoinFromDec(t.Symbol, t.TotalSupply))
			sort.Sort(coins)

			err := k.bankKeeper.SetCoins(ctx, t.Owner, coins)
			if err != nil {
				return err
			}
		}

		k.NewToken(ctx, t)
	}
	return nil
}
