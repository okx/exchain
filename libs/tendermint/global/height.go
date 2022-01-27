package global

import "sync/atomic"

var lastCommittedHeight int64

// SetGlobalHeight sets lastCommittedHeight safely.
func SetGlobalHeight(height int64) {
	atomic.StoreInt64(&lastCommittedHeight, height)
}

// GetGlobalHeight gets lastCommittedHeight safely.
func GetGlobalHeight() int64 {
	return atomic.LoadInt64(&lastCommittedHeight)
}

var maxPeerHeight int64

// SetGlobalMaxPeerHeight sets maxPeerHeight safely.
func SetGlobalMaxPeerHeight(height int64) {
	atomic.StoreInt64(&maxPeerHeight, height)
}

func GetGlobalMaxPeerHeight() int64 {
	return atomic.LoadInt64(&maxPeerHeight)
}
