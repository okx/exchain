package protocol

import (
	"encoding/json"

	"github.com/okex/okexchain/x/token"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/slashing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okexchain/x/backend"
	distr "github.com/okex/okexchain/x/distribution"
	"github.com/okex/okexchain/x/staking"
	"github.com/okex/okexchain/x/stream"
)

var _ Protocol = (*MockProtocol)(nil)

// MockProtocol is designed for engine test
type MockProtocol struct {
	parent  Parent
	version uint64
	logger  log.Logger
}

// NewMockProtocol creates a new instance of MockProtocol
func NewMockProtocol(version uint64) Protocol {
	return &MockProtocol{
		nil,
		version,
		nil,
	}
}

// GetParent returns the Parent implements
func (mp *MockProtocol) GetParent() Parent {
	if mp.parent == nil {
		panic("parent is nil in the protocol")
	}
	return mp.parent
}

// SetLogger sets logger
func (mp *MockProtocol) SetLogger(log log.Logger) Protocol {
	mp.logger = log
	return mp
}

// SetParent sets Parent and return the Protocol implements
func (mp *MockProtocol) SetParent(parent Parent) Protocol {
	mp.parent = parent
	return mp
}

// GetVersion gets version
func (mp *MockProtocol) GetVersion() uint64 {
	return mp.version
}

// nolint
func (*MockProtocol) LoadContext()                                                {}
func (*MockProtocol) Init()                                                       {}
func (*MockProtocol) CheckStopped()                                               {}
func (*MockProtocol) GetCodec() *codec.Codec                                      { return nil }
func (*MockProtocol) ExportGenesis(ctx sdk.Context) map[string]json.RawMessage    { return nil }
func (*MockProtocol) GetKVStoreKeysMap() map[string]*sdk.KVStoreKey               { return nil }
func (*MockProtocol) GetTransientStoreKeysMap() map[string]*sdk.TransientStoreKey { return nil }
func (*MockProtocol) GetBackendKeeper() backend.Keeper                            { return backend.Keeper{} }
func (*MockProtocol) GetStreamKeeper() stream.Keeper                              { return stream.Keeper{} }
func (*MockProtocol) GetCrisisKeeper() crisis.Keeper                              { return crisis.Keeper{} }
func (*MockProtocol) GetStakingKeeper() staking.Keeper                            { return staking.Keeper{} }
func (*MockProtocol) GetDistrKeeper() distr.Keeper                                { return distr.Keeper{} }
func (*MockProtocol) GetSlashingKeeper() slashing.Keeper                          { return slashing.Keeper{} }
func (*MockProtocol) GetTokenKeeper() token.Keeper                                { return token.Keeper{} }
