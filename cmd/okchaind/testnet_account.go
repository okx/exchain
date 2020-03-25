package main

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	"github.com/okex/okchain/x/common"
)

func addAccount(address string, amount int64, accs []genaccounts.GenesisAccount) []genaccounts.GenesisAccount {
	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return accs
	}
	decCoins := sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.NewDec(amount))}
	accs = append(accs, genaccounts.GenesisAccount{
		Address: addr,
		Coins:   common.ConvertDecCoinsToCoins(decCoins),
	})
	return accs
}

var (
	testnetAccountList = []string{
		"breeze real effort sail deputy spray offer real injury universe praise common",
		"build injury pool thank property awful seven farm theory crew cruel volcano",
		"replace special wing fade begin one tell tissue decrease wedding wonder crazy",
		"cube sell artist wagon husband carpet volcano salt pupil stove regular shiver",
	}
)

func getTestnetMnemonic(index int) string {
	if len(testnetAccountList)-1 < index {
		return ""
	}

	return testnetAccountList[index]
}
