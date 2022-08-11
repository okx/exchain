package global

type ModuleParamsManager interface {
	SetSendEnabled(enable bool)
	GetSendEnabled() bool
	SetSupply(coins interface{}) //coins is sdk.Coins
	GetSupply() interface{}      //sdk.Coins
}

type EmptyManager struct{}

func (e EmptyManager) SetSendEnabled(enable bool)  {}
func (e EmptyManager) GetSendEnabled() bool        { return false }
func (e EmptyManager) SetSupply(coins interface{}) {}
func (e EmptyManager) GetSupply() interface{}      { return nil }

// Manager sets module params to watchDB and avoids golang import cycle
var Manager ModuleParamsManager = EmptyManager{}
