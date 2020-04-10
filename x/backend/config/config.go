package config

import (
	"bytes"
	"encoding/json"
	okchaincfg "github.com/cosmos/cosmos-sdk/server/config"
	"os"

	"github.com/tendermint/tendermint/libs/common"
)

var (
	DefaultMaintainConfile = "maintain.conf"
	DefaultNodeHome        = okchaincfg.DefaultBackendNodeHome
	DefaultNodeCofig       = DefaultNodeHome + "/config"
	DefaultTestConfig      = DefaultNodeHome + "/test_config"
	DefaultTestDataHome    = DefaultNodeHome + "/test_data"

	DefaultConfig = okchaincfg.DefaultBackendConfig
)

type Config = okchaincfg.BackendConfig

func LoadMaintainConf(confDir string, fileName string) (*Config, error) {
	fPath := confDir + string(os.PathSeparator) + fileName
	if _, err := os.Stat(fPath); err != nil {
		return nil, err
	}

	bytes := common.MustReadFile(fPath)

	m := Config{}
	err := json.Unmarshal(bytes, &m)
	return &m, err
}

func DumpMaintainConf(maintainConf *Config, confDir string, fileName string) error {
	fPath := confDir + string(os.PathSeparator) + fileName

	if _, err := os.Stat(confDir); err != nil {
		os.MkdirAll(confDir, os.ModePerm)
	}

	if bs, err := json.Marshal(maintainConf); err != nil {
		return err
	} else {
		var out bytes.Buffer
		json.Indent(&out, bs, "", "  ")
		common.MustWriteFile(fPath, out.Bytes(), os.ModePerm)
	}

	return nil
}

func SafeLoadMaintainConfig(configDir string) *Config {
	maintainConf, err := LoadMaintainConf(configDir, DefaultMaintainConfile)
	if maintainConf == nil || err != nil {
		maintainConf = DefaultConfig()
		DumpMaintainConf(maintainConf, configDir, DefaultMaintainConfile)
	}
	return maintainConf
}
