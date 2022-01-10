package watcher

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type mockWatcher struct {
	f func(name string, exceptedF func())
	IWatcher
}

func newMockWatcher(watcher IWatcher, f func(name string, exceptedF func())) *mockWatcher {
	return &mockWatcher{f, watcher}
}

type baseWatcherCase struct {
	name string
	f    func()
}

var baseWatchers = []baseWatcherCase{
	{
		f: func() {
			watcherEnable = false
		},
		name: "disableWatcher",
	},
	{
		f: func() {
			watcherEnable = true
			AsyncTxEnable = true
		},
		name: "concurrentWatcher",
	},
	{
		f: func() {
			watcherEnable = true
			AsyncTxEnable = false
		},
		name: "normalWatcher",
	},
}

func TestWatcherType(t *testing.T) {
	IsWatcherEnabled()
	type watcherCase struct {
		base       baseWatcherCase
		exceptType interface{}
	}
	cases := []watcherCase{
		{
			base:       baseWatchers[0],
			exceptType: &disableWatcher{},
		},
		{
			base:       baseWatchers[1],
			exceptType: &concurrentWatcher{},
		},
		{
			base:       baseWatchers[2],
			exceptType: &Watcher{},
		},
	}
	for _, c := range cases {
		t.Run(c.base.name, func(t *testing.T) {
			c.base.f()
			w := NewWatcher()
			require.IsType(t, c.exceptType, w)
		})
	}
}