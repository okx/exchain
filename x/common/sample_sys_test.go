package common

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestEnableSampleSystestAll(t *testing.T) {

	e := os.Setenv("SYS_TEST_ALL", "1")
	require.Nil(t, e)

	SkipSysTestChecker(t)
	require.True(t, true)
}

func TestEnableSampleSystestSingle(t *testing.T) {

	e := os.Setenv("SYS_TEST_ALL", "0")
	e = os.Setenv("SAMPLE_SYS_TEST", "1")
	require.Nil(t, e)

	SkipSysTestChecker(t)
	require.True(t, true)
}

func TestDisableSampleSystest(t *testing.T) {
	SkipSysTestChecker(t)
	require.True(t, false)
}
