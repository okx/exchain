package cache

import (
	"sync"

	sdk "github.com/okex/exchain/libs/cosmos-sdk/types"
	"github.com/okex/exchain/libs/cosmos-sdk/x/auth/exported"
)

type cValue struct {
	acc     exported.Account
	deleted bool
	dirty   bool
}

type AccCache interface {
	GetAcc(addr sdk.AccAddress) (acc exported.Account)
	SetAcc(addr sdk.AccAddress, acc exported.Account)
	DeleteAcc(addr sdk.AccAddress)
	Write()
}

type Manager struct {
	mtx      sync.Mutex
	accounts map[string]*cValue
	parent   AccCache
}

func (m *Manager) GetAcc(addr sdk.AccAddress) (acc exported.Account) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	cv, ok := m.accounts[addr.String()]
	if !ok {
		if m.parent != nil {
			acc = m.parent.GetAcc(addr)
			m.setCacheAcc(addr, acc, false, false)
		}
	} else {
		acc = cv.acc
	}

	return acc
}

func (m *Manager) SetAcc(addr sdk.AccAddress, acc exported.Account) {
	m.mtx.Lock()
	defer m.mtx.Unlock()
	m.setCacheAcc(addr, acc, false, true)
}

func (m *Manager) DeleteAcc(addr sdk.AccAddress) {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	m.setCacheAcc(addr, nil, true, true)
}

func (m *Manager) Write() {
	m.mtx.Lock()
	defer m.mtx.Unlock()

	// We need a copy of all of the keys.
	// Not the best, but probably not a bottleneck depending.
	keys := make([]string, 0, len(m.accounts))
	for key, v := range m.accounts {
		if v.dirty {
			keys = append(keys, key)
		}
	}

	// TODO: Consider allowing usage of Batch, which would allow the write to
	// at least happen atomically.
	for _, key := range keys {
		cacheValue := m.accounts[key]
		switch {
		case cacheValue.deleted:
			if addr, err := sdk.AccAddressFromBech32(key); err == nil {
				m.parent.DeleteAcc(addr)
			}
		case cacheValue.acc == nil:
			// Skip, it already doesn't exist in parent.
		default:
			if addr, err := sdk.AccAddressFromBech32(key); err == nil {
				m.parent.SetAcc(addr, cacheValue.acc)
			}
		}
	}

	// Clear the cache
	m.accounts = make(map[string]*cValue)
}

func (m *Manager) setCacheAcc(addr sdk.AccAddress, acc exported.Account, deleted bool, dirty bool) {
	m.accounts[addr.String()] = &cValue{
		acc:     acc,
		deleted: deleted,
		dirty:   dirty,
	}
}
