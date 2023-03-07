package bank

import (
	"github.com/okx/okbchain/libs/cosmos-sdk/types/module"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/bank/internal/keeperadapter"
)

var (
	_ module.AppModuleAdapter = AppModule{}
)

func (am AppModule) RegisterServices(cfg module.Configurator) {
	RegisterBankMsgServer(cfg.MsgServer(), keeperadapter.NewMsgServerImpl(am.keeper))
	RegisterQueryServer(cfg.QueryServer(), keeperadapter.NewBankQueryServer(*am.adapterKeeper, am.supplyKeeper))
}
