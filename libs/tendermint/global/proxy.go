package global

var bankSendEnabled bool

func SetSendEnabled(enable bool) {
	bankSendEnabled = enable
}

func GetSendEnabled() bool {
	return bankSendEnabled
}

var supply interface{} //sdk.Coins

func SetSupply(coins interface{}) {
	supply = coins
}

func GetSupply() interface{} {
	return supply
}
