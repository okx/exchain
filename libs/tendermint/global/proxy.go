package global

import sdk "github.com/okex/exchain/libs/cosmos-sdk/types"

var bankSendEnabled bool

func SetSendEnabled(enable bool) {
	bankSendEnabled = enable
}

func GetSendEnabled() bool {
	return bankSendEnabled
}

var supply sdk.Coins

func SetSupply(coins sdk.Coins) {
	supply = coins
}

func GetSupply() sdk.Coins {
	return supply
}
