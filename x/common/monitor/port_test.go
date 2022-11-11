package monitor

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPortMonitor(t *testing.T) {
	pm0 := NewPortMonitor([]string{})
	require.NotNil(t, pm0)

	pm := NewPortMonitor([]string{"25559"})
	require.NotNil(t, pm)

	// error check
	require.Panics(t, func() {
		_ = NewPortMonitor([]string{"-1"})
	})

	require.Panics(t, func() {
		_ = NewPortMonitor([]string{"65536"})
	})

	require.Panics(t, func() {
		_ = NewPortMonitor([]string{"abc"})
	})
}
