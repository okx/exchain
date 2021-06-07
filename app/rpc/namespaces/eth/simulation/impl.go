package simulation

import (
	"encoding/binary"
	"github.com/okex/exchain/x/evm"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	store "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/okex/exchain/app/types"
	"github.com/okex/exchain/x/ammswap"
	"github.com/okex/exchain/x/backend"
	"github.com/okex/exchain/x/dex"
	distr "github.com/okex/exchain/x/distribution"
	evmtypes "github.com/okex/exchain/x/evm/types"
	"github.com/okex/exchain/x/evm/watcher"
	"github.com/okex/exchain/x/farm"
	"github.com/okex/exchain/x/gov"
	"github.com/okex/exchain/x/order"
	"github.com/okex/exchain/x/staking"
	"github.com/okex/exchain/x/token"
)

type QueryOnChainProxy interface {
	GetAccount(address common.Address) (*types.EthAccount, error)
	GetStorageAtInternal(address common.Address, key []byte) (hexutil.Bytes, error)
	GetCodeByHash(hash common.Hash) (hexutil.Bytes, error)
}

// AccountKeeper defines the expected account keeper interface
type AccountKeeperProxy struct {
	cachedAcc         map[string]*types.EthAccount
	queryOnChainProxy QueryOnChainProxy
	q                 *watcher.Querier
}

func NewAccountKeeperProxy(qoc QueryOnChainProxy) AccountKeeperProxy {
	return AccountKeeperProxy{
		cachedAcc:         make(map[string]*types.EthAccount, 0),
		queryOnChainProxy: qoc,
		q:                 watcher.NewQuerier(),
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
	account, e := a.queryOnChainProxy.GetAccount(common.BytesToAddress(addr.Bytes()))
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
	pr, err := p.q.GetParams()
	if err == nil {
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
	blacklistedAddrs map[string]bool
}

func NewBankKeeperProxy() BankKeeperProxy {
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
		backend.ModuleName:        nil,
		ammswap.ModuleName:        {supply.Minter, supply.Burner},
		farm.ModuleName:           nil,
		farm.YieldFarmingAccount:  nil,
		farm.MintFarmingAccount:   {supply.Burner},
	}

	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}
	return BankKeeperProxy{blacklistedAddrs: modAccAddrs}
}

func (b BankKeeperProxy) BlacklistedAddr(addr sdk.AccAddress) bool {
	return b.blacklistedAddrs[addr.String()]
}

type InternalDba struct {
	dbPrefix []byte
	ocProxy  QueryOnChainProxy
}

var (
	gSimulateCdc         *codec.Codec
	cdcOnce              sync.Once
	gSimulateChainConfig []byte
	configOnce           sync.Once
)

func instanceOfCdc() *codec.Codec {
	cdcOnce.Do(func() {
		module := evm.AppModuleBasic{}
		cdc := codec.New()
		module.RegisterCodec(cdc)
		gSimulateCdc = cdc
	})
	return gSimulateCdc
}

func instanceOfChainConfig() []byte {
	configOnce.Do(func() {
		cdc := instanceOfCdc()
		gSimulateChainConfig = cdc.MustMarshalBinaryBare(evmtypes.DefaultChainConfig())
	})
	return gSimulateChainConfig
}

func NewInternalDba(qoc QueryOnChainProxy) InternalDba {
	return InternalDba{ocProxy: qoc}
}

func (i InternalDba) NewStore(parent store.KVStore, Prefix []byte) evmtypes.StoreProxy {
	i.dbPrefix = Prefix
	if Prefix == nil {
		return nil
	}

	switch Prefix[0] {
	case evmtypes.KeyPrefixChainConfig[0]:
		return ConfigStore{defaultConfig: instanceOfChainConfig()}
	case evmtypes.KeyPrefixBloom[0]:
		return BloomStore{}
	case evmtypes.KeyPrefixStorage[0]:
		if len(Prefix) < 21 {
			return nil
		}
		return StateStore{addr: common.BytesToAddress(Prefix[1:21]), ocProxy: i.ocProxy}
	case evmtypes.KeyPrefixContractBlockedList[0]:
		return ContractBlockedListStore{watcher.NewQuerier()}
	case evmtypes.KeyPrefixContractDeploymentWhitelist[0]:
		return ContractDeploymentWhitelist{watcher.NewQuerier()}
	case evmtypes.KeyPrefixCode[0]:
		return CodeStore{q: watcher.NewQuerier(), ocProxy: i.ocProxy}
	case evmtypes.KeyPrefixHeightHash[0]:
		return HeightHashStore{watcher.NewQuerier()}
	case evmtypes.KeyPrefixBlockHash[0]:
		return BlockHashStore{}
	}
	return nil
}

type HeightHashStore struct {
	q *watcher.Querier
}

func (s HeightHashStore) Set(key, value []byte) {
	//just ignore all set opt
}

func (s HeightHashStore) Get(key []byte) []byte {
	h, _ := s.q.GetBlockHashByNumber(binary.BigEndian.Uint64(key))
	return h.Bytes()
}

func (s HeightHashStore) Has(key []byte) bool {
	return false
}

func (s HeightHashStore) Delete(key []byte) {
	return
}

type BlockHashStore struct {
}

func (s BlockHashStore) Set(key, value []byte) {
	//just ignore all set opt
}

func (s BlockHashStore) Get(key []byte) []byte {

	return nil
}

func (s BlockHashStore) Has(key []byte) bool {
	return false
}

func (s BlockHashStore) Delete(key []byte) {
	return
}

type StateStore struct {
	addr    common.Address
	ocProxy QueryOnChainProxy
}

func (s StateStore) Set(key, value []byte) {
	//just ignore all set opt
}

func (s StateStore) Get(key []byte) []byte {
	//include code and state
	b, e := s.ocProxy.GetStorageAtInternal(s.addr, key)
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
	q       *watcher.Querier
	ocProxy QueryOnChainProxy
}

func (s CodeStore) Set(key, value []byte) {
	//just ignore all set opt
}

func (s CodeStore) Get(key []byte) []byte {
	//include code and state
	b, e := s.ocProxy.GetCodeByHash(common.BytesToHash(key))
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

type ContractBlockedListStore struct {
	q *watcher.Querier
}

func (s ContractBlockedListStore) Set(key, value []byte) {
	//just ignore all set opt
}

func (s ContractBlockedListStore) Get(key []byte) []byte {
	//include code and state
	return nil
}

func (s ContractBlockedListStore) Delete(key []byte) {
	return
}

func (s ContractBlockedListStore) Has(key []byte) bool {
	return s.q.HasContractBlockedList(key)
}

type ContractDeploymentWhitelist struct {
	q *watcher.Querier
}

func (s ContractDeploymentWhitelist) Set(key, value []byte) {
	//just ignore all set opt
}

func (s ContractDeploymentWhitelist) Get(key []byte) []byte {
	//include code and state
	return nil
}

func (s ContractDeploymentWhitelist) Delete(key []byte) {
	return
}

func (s ContractDeploymentWhitelist) Has(key []byte) bool {
	return s.q.HasContractDeploymentWhitelist(key)
}
