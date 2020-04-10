package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConf(t *testing.T) {

	configDir := "/tmp/not_exists"
	configFile := "tmp.json"

	// 1. Dump & Get Non exists config file
	m := DefaultConfig()
	err := DumpMaintainConf(m, configDir, configFile)
	assert.True(t, err == nil)

	maintainConf, err := LoadMaintainConf(configDir, configFile)
	assert.True(t, maintainConf != nil && err == nil)

	fmt.Printf("%+v \n", maintainConf)

	// 2. Dump & Get already exists config file
	err = DumpMaintainConf(m, configDir, configFile)
	assert.True(t, err == nil)

	maintainConf, err = LoadMaintainConf(configDir, configFile)
	assert.True(t, maintainConf != nil && err == nil)

	fmt.Printf("%+v \n", maintainConf)

	// 3. SafeLoadMaintainConfig
	os.RemoveAll(DefaultTestConfig)
	config := SafeLoadMaintainConfig(DefaultTestConfig)
	assert.True(t, config != nil)
	SafeLoadMaintainConfig(DefaultTestConfig)

}
