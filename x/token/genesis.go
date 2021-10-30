package token

import (
	"errors"
	"fmt"

	sdk "github.com/okex/exchain/dependence/cosmos-sdk/types"
	"github.com/okex/exchain/x/common"
	"github.com/okex/exchain/x/token/types"
)

// default owner of okt
const DefaultTokenOwner = "ex10q0rk5qnyag7wfvvt7rtphlw589m7frs3hvqmf"

// all state that must be provided in genesis file
type GenesisState struct {
	Params       types.Params     `json:"params"`
	Tokens       []types.Token    `json:"tokens"`
	LockedAssets []types.AccCoins `json:"locked_assets"`
	LockedFees   []types.AccCoins `json:"locked_fees"`
}

// default GenesisState used by Cosmos Hub
func defaultGenesisState() GenesisState {
	return GenesisState{
		Params:       types.DefaultParams(),
		Tokens:       []types.Token{defaultGenesisStateOKT()},
		LockedAssets: nil,
		LockedFees:   nil,
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
		OriginalSymbol:      common.NativeToken,
		WholeName:           "OKT",
		OriginalTotalSupply: totalSupply,
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
	// if module account dosen't exist, it will create automatically
	moduleAcc := keeper.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
	if moduleAcc == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set params
	keeper.SetParams(ctx, data.Params)

	for _, token := range data.Tokens {
		keeper.NewToken(ctx, token)
	}

	for _, lock := range data.LockedAssets {
		if err := keeper.updateLockedCoins(ctx, lock.Acc, lock.Coins, true, types.LockCoinsTypeQuantity); err != nil {
			panic(err)
		}
	}
	for _, lock := range data.LockedFees {
		if err := keeper.updateLockedCoins(ctx, lock.Acc, lock.Coins, true, types.LockCoinsTypeFee); err != nil {
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
	lockedAsset := keeper.GetAllLockedCoins(ctx)

	var lockedFees []types.AccCoins
	keeper.IterateLockedFees(ctx, func(acc sdk.AccAddress, coins sdk.SysCoins) bool {
		lockedFees = append(lockedFees,
			types.AccCoins{
				Acc:   acc,
				Coins: coins,
			})
		return false
	})

	return GenesisState{
		Params:       params,
		Tokens:       tokens,
		LockedAssets: lockedAsset,
		LockedFees:   lockedFees,
	}
}
