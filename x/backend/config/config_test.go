package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestConf(t *testing.T) {

	configDir := "/tmp/not_exists"
	configFile := "tmp.json"

	// 1. Dump & Get Non exists config file
	m := DefaultConfig()
	err := dumpMaintainConf(m, configDir, configFile)
	assert.True(t, err == nil)

	maintainConf, err := loadMaintainConf(configDir, configFile)
	assert.True(t, maintainConf != nil && err == nil)

	fmt.Printf("%+v \n", maintainConf)

	// 2. Dump & Get already exists config file
	err = dumpMaintainConf(m, configDir, configFile)
	assert.True(t, err == nil)

	maintainConf, err = loadMaintainConf(configDir, configFile)
	assert.True(t, maintainConf != nil && err == nil)

	fmt.Printf("%+v \n", maintainConf)

	// 3. SafeLoadMaintainConfig
	err = os.RemoveAll(DefaultTestConfig)
	require.Nil(t, err)
	config, err := SafeLoadMaintainConfig(DefaultTestConfig)
	assert.True(t, config != nil && err == nil)
}
