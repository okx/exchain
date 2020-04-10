package protocol

import (
	"fmt"
	"sync"

	"github.com/okex/okchain/x/common/monitor"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common/proto"
)

var (
	// Singleton Pattern: protocolEngine
	once            sync.Once
	protocolsEngine *appProtocolEngine

	// init monitor prometheus metrics
	orderMetrics  = monitor.DefaultOrderMetrics(monitor.DefaultPrometheusConfig())
	streamMetrics = monitor.DefaultStreamMetrics(monitor.DefaultPrometheusConfig())
)

// GetEngine gets the Singleton application protocol engine
func GetEngine() *appProtocolEngine {
	once.Do(func() {
		protocolKeeper := proto.NewProtocolKeeper(GetMainStoreKey())
		protocolsEngine = NewAppProtocolEngine(protocolKeeper)
		// add protocolV0
		protocolsEngine.Add(NewProtocolV0(nil, 0, nil, 0, protocolKeeper))
	})
	return protocolsEngine
}

type appProtocolEngine struct {
	protocols map[uint64]Protocol
	current   uint64
	next      uint64
	keeper    proto.VersionKeeper
}

// NewAppProtocolEngine returns a pointer of a new appProtocolEngine object
func NewAppProtocolEngine(protocolKeeper proto.ProtocolKeeper) *appProtocolEngine {
	return &appProtocolEngine{
		make(map[uint64]Protocol),
		0,
		0,
		protocolKeeper,
	}
}

// Add adds new protocol into engine
func (ape *appProtocolEngine) Add(p Protocol) {
	if p.GetVersion() != ape.next {
		panic(fmt.Errorf("wrong version being added to the protocol engine: %d; expecting %d", p.GetVersion(), ape.next))
	}
	ape.protocols[ape.next] = p
	ape.next++
}

// GetProtocolKeeper gets protocol keeper from engine
func (ape *appProtocolEngine) GetProtocolKeeper() proto.ProtocolKeeper {
	return ape.keeper.(proto.ProtocolKeeper)
}

// GetCurrentProtocol gets current protolcol from engine
func (ape *appProtocolEngine) GetCurrentProtocol() Protocol {
	p, flag := ape.protocols[ape.current]
	if !flag {
		panic("Invalid Protocol")
	}
	return p
}

// LoadCurrentProtocol loads the status of current protocol from store
func (ape *appProtocolEngine) LoadCurrentProtocol(kvStore sdk.KVStore) (bool, uint64) {
	// find the current version from store
	current := ape.GetCurrentVersionByStore(kvStore)
	p, flag := ape.protocols[current]
	if flag {
		//p.Init()
		p.LoadContext()
		ape.current = current
	}
	return flag, current
}

// GetCurrentVersionByStore gets current version from Store
func (ape *appProtocolEngine) GetCurrentVersionByStore(store sdk.KVStore) uint64 {
	return ape.keeper.GetCurrentVersionByStore(store)
}

// GetCurrentVersion gets current version from engine
func (ape *appProtocolEngine) GetCurrentVersion() uint64 {
	return ape.current
}

// GetUpgradeConfigByStore gets upgrade config from store
func (ape *appProtocolEngine) GetUpgradeConfigByStore(store sdk.KVStore) (upgradeConfig proto.AppUpgradeConfig,
	found bool) {
	return ape.keeper.GetUpgradeConfigByStore(store)
}

// Activate activates new protocol at specific block height
func (ape *appProtocolEngine) Activate(version uint64) bool {
	protocol, flag := ape.protocols[version]
	if flag {
		//protocol.Init()
		protocol.LoadContext()
		ape.current = version
	}
	return flag
}

// deprecated
func (ape *appProtocolEngine) Clear() {
	ape.protocols = make(map[uint64]Protocol)
	ape.current = 0
	ape.next = 0
}

// FillProtocol sets logger and app
func (ape *appProtocolEngine) FillProtocol(parent Parent, log log.Logger, index uint64) error {
	protocol, ok := ape.protocols[index]
	if !ok {
		return fmt.Errorf("no protocol with version %d in the engine", index)
	}
	protocol.SetParent(parent).SetLogger(log)
	return nil
}
