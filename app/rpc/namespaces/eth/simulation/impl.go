package simulation

import (
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/okex/exchain/app/types"
	evmtypes "github.com/okex/exchain/x/evm/types"
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeperProxy struct {
}

func (a AccountKeeperProxy) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account {
	return types.EthAccount{}
}

func (a AccountKeeperProxy) GetAllAccounts(ctx sdk.Context) (accounts []authexported.Account) {
	return nil
}

func (a AccountKeeperProxy) IterateAccounts(ctx sdk.Context, cb func(account authexported.Account) bool) {
	return
}

func (a AccountKeeperProxy) GetAccount(ctx sdk.Context, addr sdk.AccAddress) authexported.Account {
	return types.EthAccount{}
}

func (a AccountKeeperProxy) SetAccount(ctx sdk.Context, account authexported.Account) {
	return
}

func (a AccountKeeperProxy) RemoveAccount(ctx sdk.Context, account authexported.Account) {
	return
}

type SupplyKeeperProxy struct {
}

func (s SupplyKeeperProxy) SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error {
	return nil
}

type SubspaceProxy struct {
}

func (p SubspaceProxy) GetParamSet(ctx sdk.Context, ps params.ParamSet) {

}
func (p SubspaceProxy) SetParamSet(ctx sdk.Context, ps params.ParamSet) {

}

type BankKeeperProxy struct {
}

func (b BankKeeperProxy) BlacklistedAddr(addr sdk.AccAddress) bool {
	return false
}

type InternalDba struct {
}

func (i InternalDba) NewStore(parent store.KVStore, Prefix []byte) evmtypes.StoreProxy {
	return Store{}
}

type Store struct {
}

func (s Store) Set(key, value []byte) {
	return
}

func (s Store) Get(key []byte) []byte {
	return nil
}

func (s Store) Delete(key []byte) {
	return
}
