package protocol

import (
	"errors"
	"fmt"
	"sync"

	"github.com/okex/okchain/x/common/monitor"

	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okex/okchain/x/common/proto"
)

var (
	once            sync.Once
	protocolsEngine *AppProtocolEngine

	// init monitor prometheus metrics
	orderMetrics  = monitor.DefaultOrderMetrics(monitor.DefaultPrometheusConfig())
	streamMetrics = monitor.DefaultStreamMetrics(monitor.DefaultPrometheusConfig())
)

// get the Singleton application protocol engine
func GetEngine() *AppProtocolEngine {
	once.Do(func() {
		protocolKeeper := proto.NewProtocolKeeper(GetMainStoreKey())
		protocolsEngine = NewAppProtocolEngine(protocolKeeper)
		// add protocolV0
		protocolsEngine.Add(NewProtocolV0(nil, 0, nil, 0, protocolKeeper))
	})
	return protocolsEngine
}

type AppProtocolEngine struct {
	protocols map[uint64]Protocol
	current   uint64
	next      uint64
	keeper    proto.VersionKeeper
}

func NewAppProtocolEngine(protocolKeeper proto.ProtocolKeeper) *AppProtocolEngine {
	return &AppProtocolEngine{
		make(map[uint64]Protocol),
		0,
		0,
		protocolKeeper,
	}

}

// add new protocol into engine
func (ape *AppProtocolEngine) Add(p Protocol) {
	if p.GetVersion() != ape.next {
		panic(fmt.Errorf("wrong version being added to the protocol engine: %d; expecting %d", p.GetVersion(), ape.next))
	}
	ape.protocols[ape.next] = p
	ape.next++
}

// get protocol keeper from engine
func (ape *AppProtocolEngine) GetProtocolKeeper() proto.ProtocolKeeper {
	return ape.keeper.(proto.ProtocolKeeper)
}

// get current protolcol from engine
func (ape *AppProtocolEngine) GetCurrentProtocol() Protocol {
	p, flag := ape.protocols[ape.current]
	if !flag {
		panic("Invalid Protocol")
	}
	return p
}

// load the status of current protocol from store
func (ape *AppProtocolEngine) LoadCurrentProtocol(kvStore sdk.KVStore) (bool, uint64) {
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

// get current version from Store
func (ape *AppProtocolEngine) GetCurrentVersionByStore(store sdk.KVStore) uint64 {
	return ape.keeper.GetCurrentVersionByStore(store)
}

// get current version from engine
func (ape *AppProtocolEngine) GetCurrentVersion() uint64 {
	return ape.current
}

// get upgrade config from store
func (ape *AppProtocolEngine) GetUpgradeConfigByStore(store sdk.KVStore) (upgradeConfig proto.AppUpgradeConfig, found bool) {
	return ape.keeper.GetUpgradeConfigByStore(store)
}

// activate new protocol at specific block height
func (ape *AppProtocolEngine) Activate(version uint64) bool {
	protocol, flag := ape.protocols[version]
	if flag {
		//protocol.Init()
		protocol.LoadContext()
		ape.current = version
	}
	return flag
}

// for unittest

// deprecated
func (ape *AppProtocolEngine) clear() {
	ape.protocols = make(map[uint64]Protocol)
	ape.current = 0
	ape.next = 0
}

// deprecated
func (ape *AppProtocolEngine) Clear() {
	ape.clear()
}

// set log and app
func (ape *AppProtocolEngine) FillProtocol(parent Parent, log log.Logger, index uint64) error {
	protocol, ok := ape.protocols[index]
	if !ok {
		return errors.New(fmt.Sprintf("no protocol with version %d in the engine", index))
	}
	protocol.SetParent(parent).SetLogger(log)
	return nil
}
