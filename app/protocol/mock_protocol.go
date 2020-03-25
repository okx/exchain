package protocol

import (
	"encoding/json"
	"github.com/okex/okchain/x/token"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/slashing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/okex/okchain/x/backend"
	distr "github.com/okex/okchain/x/distribution"
	"github.com/okex/okchain/x/staking"
	"github.com/okex/okchain/x/stream"
)

// MockProtocol is designed 4 engine test

var _ Protocol = (*MockProtocol)(nil)

type MockProtocol struct {
	parent  Parent
	version uint64
	logger  log.Logger
}

func (mp *MockProtocol) GetParent() Parent {
	if mp.parent == nil {
		panic("parent is nil in the protocol")
	}
	return mp.parent
}

func (mp *MockProtocol) SetLogger(log log.Logger) Protocol {
	mp.logger = log
	return mp
}

func (mp *MockProtocol) SetParent(parent Parent) Protocol {
	mp.parent = parent
	return mp
}

func (mp *MockProtocol) GetCrisisKeeper() crisis.Keeper {
	return crisis.Keeper{}
}

func (mp *MockProtocol) GetStakingKeeper() staking.Keeper {
	return staking.Keeper{}
}

func (mp *MockProtocol) GetDistrKeeper() distr.Keeper {
	return distr.Keeper{}
}

func (mp *MockProtocol) GetSlashingKeeper() slashing.Keeper {
	return slashing.Keeper{}
}

func (mp *MockProtocol) GetTokenKeeper() token.Keeper {
	return token.Keeper{}
}

func (mp *MockProtocol) ExportGenesis(ctx sdk.Context) map[string]json.RawMessage {
	return nil
}

func NewMockProtocol(version uint64) Protocol {
	return &MockProtocol{
		nil,
		version,
		nil,
	}
}

func (mp *MockProtocol) GetVersion() uint64 {
	return mp.version
}

func (*MockProtocol) ExportAppStateAndValidators(ctx sdk.Context) (appState json.RawMessage, validators []types.GenesisValidator, err error) {
	return nil, nil, nil
}

func (*MockProtocol) LoadContext() {}

func (*MockProtocol) Init() {}

func (*MockProtocol) GetCodec() *codec.Codec {
	return nil
}

func (*MockProtocol) CheckStopped() {}

func (*MockProtocol) GetBackendKeeper() backend.Keeper {
	return backend.Keeper{}
}

func (*MockProtocol) GetStreamKeeper() stream.Keeper {
	return stream.Keeper{}
}

func (*MockProtocol) GetKVStoreKeysMap() map[string]*sdk.KVStoreKey {
	return nil
}

func (*MockProtocol) GetTransientStoreKeysMap() map[string]*sdk.TransientStoreKey {
	return nil
}
