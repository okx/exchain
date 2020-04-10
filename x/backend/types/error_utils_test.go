package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewErrorsMerged(t *testing.T) {
	e1 := errors.New("e1")
	e2 := errors.New("e2")
	var e3 error

	m1 := NewErrorsMerged(e3)
	require.Nil(t, m1)

	m2 := NewErrorsMerged(e1)
	require.NotNil(t, m2)
	require.Contains(t, m2.Error(), "e1")
	require.NotContains(t, m2.Error(), "e2")
	println(m2.Error())

	m3 := NewErrorsMerged(e1, e2, e3)
	println(m3.Error())
	require.NotNil(t, m3)

	require.Contains(t, m3.Error(), "e1")
	require.Contains(t, m3.Error(), "e2")
	require.NotContains(t, m3.Error(), "e3")

}
