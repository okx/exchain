package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnableSampleSystestAll(t *testing.T) {

	e := os.Setenv("SYS_TEST_ALL", "1")
	require.Nil(t, e)

	SkipSysTestChecker(t)
	require.True(t, true)
}

func TestEnableSampleSystestSingle(t *testing.T) {

	e := os.Setenv("SYS_TEST_ALL", "0")
	require.Nil(t, e)
	e = os.Setenv("SAMPLE_SYS_TEST", "1")
	require.Nil(t, e)

	SkipSysTestChecker(t)
	require.True(t, true)
}

func TestDisableSampleSystest(t *testing.T) {
	SkipSysTestChecker(t)
	require.True(t, false)
}
