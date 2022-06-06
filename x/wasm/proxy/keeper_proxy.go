package proxy

import (
	"fmt"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	apptypes "github.com/okex/exchain/app/types"
	ethermint "github.com/okex/exchain/app/types"
	"github.com/okex/exchain/libs/cosmos-sdk/codec"
	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth"
	authexported "github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
	"github.com/okex/exchain/libs/cosmos-sdk/x/bank"
	capabilitytypes "github.com/okex/exchain/libs/cosmos-sdk/x/capability/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/mint"
	"github.com/okex/exchain/libs/cosmos-sdk/x/params"
	"github.com/okex/exchain/libs/cosmos-sdk/x/supply"
	supplyexported "github.com/okex/exchain/libs/cosmos-sdk/x/supply/exported"
	"github.com/okex/exchain/libs/tendermint/global"
	"github.com/okex/exchain/x/ammswap"
	dex "github.com/okex/exchain/x/dex/types"
	distr "github.com/okex/exchain/x/distribution"
	"github.com/okex/exchain/x/farm"
	"github.com/okex/exchain/x/gov"
	"github.com/okex/exchain/x/order"
	"github.com/okex/exchain/x/staking"
	token "github.com/okex/exchain/x/token/types"
	"github.com/okex/exchain/x/wasm/types"
	"github.com/okex/exchain/x/wasm/watcher"
)

type ChainQuerier interface {
	QueryWithData(path string, data []byte) ([]byte, int64, error)
}

// AccountKeeper defines the expected account keeper interface
type AccountKeeperProxy struct {
	cdc       *codec.Codec
	cq        ChainQuerier
	cachedAcc map[string]*apptypes.EthAccount
}

func NewAccountKeeperProxy(cdc *codec.Codec, cq ChainQuerier) AccountKeeperProxy {
	return AccountKeeperProxy{
		cdc: cdc,
		cq:  cq,
	}
}

func (a AccountKeeperProxy) SetObserverKeeper(observer auth.ObserverI) {}

func (a AccountKeeperProxy) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account {
	acc := apptypes.EthAccount{
		BaseAccount: &auth.BaseAccount{
			Address: addr,
		},
		CodeHash: ethcrypto.Keccak256(nil),
	}

	a.SetAccount(ctx, acc)
	return &acc
}

func (a AccountKeeperProxy) GetAllAccounts(ctx sdk.Context) (accounts []authexported.Account) {
	return nil
}

func (a AccountKeeperProxy) IterateAccounts(ctx sdk.Context, cb func(account authexported.Account) bool) {
	return
}

func (a AccountKeeperProxy) GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account {
	acc, ok := a.cachedAcc[addr.String()]
	if ok {
		return acc
	}
	bs, err := a.cdc.MarshalJSON(auth.NewQueryAccountParams(addr.Bytes()))
	if err != nil {
		return nil
	}

	res, _, err := a.cq.QueryWithData(fmt.Sprintf("custom/%s/%s", auth.QuerierRoute, auth.QueryAccount), bs)
	if err != nil {
		return nil
	}

	var account ethermint.EthAccount
	if err := a.cdc.UnmarshalJSON(res, &account); err != nil {
		return nil
	}
	a.SetAccount(ctx, account)
	// TODO:
	// set account to watch db
	return account
}

func (a AccountKeeperProxy) SetAccount(ctx sdk.Context, account authexported.Account, updateState ...bool) {
	acc, ok := account.(apptypes.EthAccount)
	if !ok {
		return
	}
	// delay make
	if a.cachedAcc == nil {
		a.cachedAcc = make(map[string]*apptypes.EthAccount)
	}
	a.cachedAcc[account.GetAddress().String()] = &acc
	return
}

func (a AccountKeeperProxy) RemoveAccount(ctx sdk.Context, account authexported.Account) {
	delete(a.cachedAcc, account.GetAddress().String())
}

type SubspaceProxy struct{}

func (s SubspaceProxy) GetParamSet(ctx sdk.Context, ps params.ParamSet) {
	if wasmParams, ok := ps.(*types.Params); ok {
		wasmParams.CodeUploadAccess = watcher.Params.CodeUploadAccess
		wasmParams.InstantiateDefaultPermission = watcher.Params.InstantiateDefaultPermission
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
		dex.ModuleName:            nil,
		order.ModuleName:          nil,
		ammswap.ModuleName:        {supply.Minter, supply.Burner},
		farm.ModuleName:           nil,
		farm.YieldFarmingAccount:  nil,
		farm.MintFarmingAccount:   {supply.Burner},
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
	acc := b.akp.GetAccount(ctx, addr)
	return acc.GetCoins()
}

func (b BankKeeperProxy) GetBalance(ctx sdk.Context, addr sdk.AccAddress, denom string) sdk.Coin {
	acc := b.akp.GetAccount(ctx, addr)
	return sdk.Coin{
		Denom:  denom,
		Amount: acc.GetCoins().AmountOf(denom),
	}
}

func (b BankKeeperProxy) IsSendEnabledCoins(ctx sdk.Context, coins ...sdk.Coin) error {
	if b.GetSendEnabled(ctx) {
		return bank.ErrSendDisabled
	}
	return nil
}

func (b BankKeeperProxy) GetSendEnabled(ctx sdk.Context) bool {
	return global.GetSendEnabled()
}

func (b BankKeeperProxy) BlockedAddr(addr sdk.AccAddress) bool {
	return b.BlacklistedAddr(addr)
}

func (b BankKeeperProxy) BlacklistedAddr(addr sdk.AccAddress) bool {
	return b.blacklistedAddrs[addr.String()]
}

func (b BankKeeperProxy) SendCoins(ctx sdk.Context, fromAddr sdk.AccAddress, toAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

type SupplyKeeperProxy struct{}

func (s SupplyKeeperProxy) GetSupply(ctx sdk.Context) supplyexported.SupplyI {
	return supply.Supply{
		Total: global.GetSupply(),
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
