package simulation

import (
	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/app/types"
	"github.com/okex/exchain/x/evm"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
)

// AccountKeeper defines the expected account keeper interface
type AccountKeeperProxy struct {
	cachedAcc map[string]*types.EthAccount
	q         *watcher.Querier
}

func NewAccountKeeperProxy() AccountKeeperProxy {
	return AccountKeeperProxy{
		cachedAcc: make(map[string]*types.EthAccount, 0),
		q:         watcher.NewQuerier(),
	}
}

func (a AccountKeeperProxy) SetObserverKeeper(observer auth.ObserverI) {
}

func (a AccountKeeperProxy) NewAccountWithAddress(ctx sdk.Context, addr sdk.AccAddress) authexported.Account {
	acc := types.EthAccount{
		BaseAccount: &auth.BaseAccount{},
		CodeHash:    ethcrypto.Keccak256(nil),
	}
	acc.SetAddress(addr)
	a.cachedAcc[addr.String()] = &acc
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
	account, e := a.q.GetAccount(addr)
	if e != nil {
		//query account from chain
		return nil
	}
	return account
}

func (a AccountKeeperProxy) SetAccount(ctx sdk.Context, account authexported.Account) {
	acc, ok := account.(types.EthAccount)
	if !ok {
		return
	}
	a.cachedAcc[account.GetAddress().String()] = &acc
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
	q *watcher.Querier
}

func NewSubspaceProxy() SubspaceProxy {
	return SubspaceProxy{
		q: watcher.NewQuerier(),
	}
}

func (p SubspaceProxy) GetParamSet(ctx sdk.Context, ps params.ParamSet) {
	pr, e := p.q.GetParams()
	if e == nil {
		evmParam := ps.(*evmtypes.Params)
		evmParam.MaxGasLimitPerTx = pr.MaxGasLimitPerTx
		evmParam.EnableCall = pr.EnableCall
		evmParam.EnableContractBlockedList = pr.EnableContractBlockedList
		evmParam.EnableCreate = pr.EnableCreate
		evmParam.ExtraEIPs = pr.ExtraEIPs
		evmParam.EnableContractDeploymentWhitelist = pr.EnableContractDeploymentWhitelist
	}

}
func (p SubspaceProxy) SetParamSet(ctx sdk.Context, ps params.ParamSet) {

}

type BankKeeperProxy struct {
}

func (b BankKeeperProxy) BlacklistedAddr(addr sdk.AccAddress) bool {
	return false
}

type InternalDba struct {
	dbPrefix []byte
}

func newCdc() *codec.Codec {
	module := evm.AppModuleBasic{}
	cdc := codec.New()
	module.RegisterCodec(cdc)
	return cdc
}

func (i InternalDba) NewStore(parent store.KVStore, Prefix []byte) evmtypes.StoreProxy {
	i.dbPrefix = Prefix
	if Prefix == nil {
		return nil
	}
	if len(Prefix) >= 21 {
		return StateStore{addr: common.BytesToAddress(Prefix[1:21]), q: watcher.NewQuerier()}
	}
	cdc := newCdc()
	switch Prefix[0] {
	case evmtypes.KeyPrefixChainConfig[0]:
		return ConfigStore{defaultConfig: cdc.MustMarshalBinaryBare(evmtypes.DefaultChainConfig())}
	case evmtypes.KeyPrefixBloom[0]:
		return BloomStore{}
	}
	return CodeStore{q: watcher.NewQuerier()}
}

type StateStore struct {
	addr common.Address
	q    *watcher.Querier
}

func (s StateStore) Set(key, value []byte) {
	//just ignore all set opt
	return
}

func (s StateStore) Get(key []byte) []byte {
	//include code and state
	b, e := s.q.GetState(s.addr, key)
	if e != nil {
		return nil
	}
	return b
}

func (s StateStore) Has(key []byte) bool {
	return false
}

func (s StateStore) Delete(key []byte) {
	return
}

type ConfigStore struct {
	defaultConfig []byte
}

func (s ConfigStore) Set(key, value []byte) {
	//just ignore all set opt
	return
}

func (s ConfigStore) Get(key []byte) []byte {
	return s.defaultConfig
}

func (s ConfigStore) Delete(key []byte) {
	return
}

func (s ConfigStore) Has(key []byte) bool {
	return false
}

type BloomStore struct {
}

func (s BloomStore) Set(key, value []byte) {
	//just ignore all set opt
	return
}

func (s BloomStore) Get(key []byte) []byte {
	return nil
}

func (s BloomStore) Delete(key []byte) {
	return
}

func (s BloomStore) Has(key []byte) bool {
	return false
}

type CodeStore struct {
	q *watcher.Querier
}

func (s CodeStore) Set(key, value []byte) {
	//just ignore all set opt
	return
}

func (s CodeStore) Get(key []byte) []byte {
	//include code and state
	b, e := s.q.GetCodeByHash(key)
	if e != nil {
		return nil
	}
	return b
}

func (s CodeStore) Delete(key []byte) {
	return
}

func (s CodeStore) Has(key []byte) bool {
	return false
}
