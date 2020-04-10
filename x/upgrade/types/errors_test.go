package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewError(t *testing.T) {
	require.Error(t, NewError(DefaultCodespace, CodeInvalidMsgType, ""))
	require.Error(t, NewError(DefaultCodespace, CodeInvalidMsgType, "test"))
	require.Error(t, NewError(DefaultCodespace, CodeUnSupportedMsgType, "test"))
	require.Error(t, NewError(DefaultCodespace, CodeUnSupportedMsgType, ""))
	require.Error(t, NewError(DefaultCodespace, CodeDoubleSwitch, ""))
}
