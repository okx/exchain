package proxy

import (
	"fmt"
	"log"

	apptypes "github.com/okx/okbchain/app/types"
	types2 "github.com/okx/okbchain/libs/cosmos-sdk/store/types"
	sdk "github.com/okx/okbchain/libs/cosmos-sdk/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/auth"
	authexported "github.com/okx/okbchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/bank"
	capabilitytypes "github.com/okx/okbchain/libs/cosmos-sdk/x/capability/types"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/mint"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/params"
	"github.com/okx/okbchain/libs/cosmos-sdk/x/supply"
	supplyexported "github.com/okx/okbchain/libs/cosmos-sdk/x/supply/exported"
	"github.com/okx/okbchain/libs/tendermint/global"
	distr "github.com/okx/okbchain/x/distribution"
	"github.com/okx/okbchain/x/gov"
	ptypes "github.com/okx/okbchain/x/params/types"
	"github.com/okx/okbchain/x/staking"
	token "github.com/okx/okbchain/x/token/types"
	"github.com/okx/okbchain/x/wasm/types"
	"github.com/okx/okbchain/x/wasm/watcher"
)

const (
	accountBytesLen = 80
)

var gasConfig = types2.KVGasConfig()

// AccountKeeperProxy defines the expected account keeper interface
type AccountKeeperProxy struct {
	cachedAcc map[string]*apptypes.EthAccount
}

func NewAccountKeeperProxy() AccountKeeperProxy {
	return AccountKeeperProxy{}
}

func (a AccountKeeperProxy) SetObserverKeeper(observer auth.ObserverI) {}

func (a AccountKeeperProxy) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account {
	ctx.GasMeter().ConsumeGas(3066, "AccountKeeperProxy NewAccountWithAddress")
	acc := apptypes.EthAccount{
		BaseAccount: &auth.BaseAccount{
			Address: addr,
		},
	}
	return &acc
}

func (a AccountKeeperProxy) GetAllAccounts(ctx sdk.Context) (accounts []authexported.Account) {
	return nil
}

func (a AccountKeeperProxy) IterateAccounts(ctx sdk.Context, cb func(account authexported.Account) bool) {
}

func (a AccountKeeperProxy) GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account {
	ctx.GasMeter().ConsumeGas(gasConfig.ReadCostFlat, types2.GasReadCostFlatDesc)
	ctx.GasMeter().ConsumeGas(gasConfig.ReadCostPerByte*accountBytesLen, types2.GasReadPerByteDesc)
	acc, ok := a.cachedAcc[addr.String()]
	if ok {
		return acc
	}
	return nil
}

func (a AccountKeeperProxy) SetAccount(ctx sdk.Context, account authexported.Account) {
	acc, ok := account.(*apptypes.EthAccount)
	if !ok {
		return
	}
	// delay to make
	if a.cachedAcc == nil {
		a.cachedAcc = make(map[string]*apptypes.EthAccount)
	}
	a.cachedAcc[account.GetAddress().String()] = acc
	ctx.GasMeter().ConsumeGas(gasConfig.WriteCostFlat, types2.GasWriteCostFlatDesc)
	ctx.GasMeter().ConsumeGas(gasConfig.WriteCostPerByte*accountBytesLen, types2.GasWritePerByteDesc)
	return
}

func (a AccountKeeperProxy) RemoveAccount(ctx sdk.Context, account authexported.Account) {
	delete(a.cachedAcc, account.GetAddress().String())
	ctx.GasMeter().ConsumeGas(gasConfig.DeleteCost, types2.GasDeleteDesc)
}

type SubspaceProxy struct{}

func (s SubspaceProxy) GetParamSet(ctx sdk.Context, ps params.ParamSet) {
	ctx.GasMeter().ConsumeGas(2111, "SubspaceProxy GetParamSet")
	if wasmParams, ok := ps.(*types.Params); ok {
		wps := watcher.GetParams()
		wasmParams.CodeUploadAccess = wps.CodeUploadAccess
		wasmParams.InstantiateDefaultPermission = wps.InstantiateDefaultPermission
	}
}
func (s SubspaceProxy) SetParamSet(ctx sdk.Context, ps params.ParamSet) {}

type BankKeeperProxy struct {
	blacklistedAddrs map[string]bool
	akp              AccountKeeperProxy
}

func NewBankKeeperProxy(akp AccountKeeperProxy) BankKeeperProxy {
	modAccAddrs := make(map[string]bool)
	maccPerms := map[string][]string{
		auth.FeeCollectorName:     nil,
		distr.ModuleName:          nil,
		mint.ModuleName:           {supply.Minter},
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
		gov.ModuleName:            nil,
		token.ModuleName:          {supply.Minter, supply.Burner},
	}

	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	return BankKeeperProxy{
		blacklistedAddrs: modAccAddrs,
		akp:              akp,
	}
}

func (b BankKeeperProxy) GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins {
	acc, err := watcher.GetAccount(addr)
	if err == nil {
		return acc.GetCoins()
	}

	bs, err := clientCtx.Codec.MarshalJSON(auth.NewQueryAccountParams(addr.Bytes()))
	if err != nil {
		log.Println("GetAllBalances marshal json error", err)
		return sdk.NewCoins()
	}
	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", auth.QuerierRoute, auth.QueryAccount), bs)
	if err != nil {
		log.Println("GetAllBalances query with data error", err)
		return sdk.NewCoins()
	}
	var account apptypes.EthAccount
	err = clientCtx.Codec.UnmarshalJSON(res, &account)
	if err != nil {
		log.Println("GetAllBalances unmarshal json error", err)
		return sdk.NewCoins()
	}

	if err = watcher.SetAccount(&account); err != nil {
		log.Println("GetAllBalances save account error", err)
	}

	return account.GetCoins()
}

func (b BankKeeperProxy) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	coins := b.GetAllBalances(ctx, addr)
	return sdk.Coin{
		Amount: coins.AmountOf(denom),
		Denom:  denom,
	}
}

func (b BankKeeperProxy) IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error {
	if b.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled
	}
	return nil
}

func (b BankKeeperProxy) GetSendEnabled(ctx sdk.Context) bool {
	ctx.GasMeter().ConsumeGas(1012, "BankKeeperProxy GetSendEnabled")
	return global.Manager.GetSendEnabled()
}

func (b BankKeeperProxy) BlockedAddr(addr sdk.AccAddress) bool {
	return b.BlacklistedAddr(addr)
}

func (b BankKeeperProxy) BlacklistedAddr(addr sdk.AccAddress) bool {
	return b.blacklistedAddrs[addr.String()]
}

func (b BankKeeperProxy) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	ctx.GasMeter().ConsumeGas(16748, "BankKeeperProxy SendCoins")
	return nil
}

type SupplyKeeperProxy struct{}

func (s SupplyKeeperProxy) GetSupply(ctx sdk.Context) supplyexported.SupplyI {
	//TODO: cache total supply in watchDB
	//rarely used, so just query from chain db
	tsParams := supply.NewQueryTotalSupplyParams(1, 0) // no pagination
	bz, err := clientCtx.Codec.MarshalJSON(tsParams)
	if err != nil {
		return supply.DefaultSupply()
	}

	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", supply.QuerierRoute, supply.QueryTotalSupply), bz)
	if err != nil {
		return supply.DefaultSupply()
	}

	var totalSupply sdk.Coins
	err = clientCtx.Codec.UnmarshalJSON(res, &totalSupply)
	if err != nil {
		return supply.DefaultSupply()
	}
	return supply.Supply{
		Total: totalSupply,
	}
}

type CapabilityKeeperProxy struct{}

func (c CapabilityKeeperProxy) GetCapability(ctx sdk.Context, name string) (*capabilitytypes.Capability, bool) {
	return nil, false
}

func (c CapabilityKeeperProxy) ClaimCapability(ctx sdk.Context, cap *capabilitytypes.Capability, name string) error {
	return nil
}

func (c CapabilityKeeperProxy) AuthenticateCapability(ctx sdk.Context, capability *capabilitytypes.Capability, name string) bool {
	return false
}

type PortKeeperProxy struct{}

func (p PortKeeperProxy) BindPort(ctx sdk.Context, portID string) *capabilitytypes.Capability {
	return nil
}

type ParamsKeeperProxy struct{}

func (p ParamsKeeperProxy) ClaimReadyForUpgrade(name string, cb func(ptypes.UpgradeInfo)) {}
