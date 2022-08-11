package global

type ModuleParamsManager interface {
	SetSendEnabled(enable bool)
	GetSendEnabled() bool
	SetSupply(coins interface{}) //coins is sdk.Coins
	GetSupply() interface{}      //sdk.Coins
}

// Manager sets module params to watchDB and avoids golang import cycle
var Manager ModuleParamsManager
