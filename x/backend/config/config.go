package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	okchaincfg "github.com/cosmos/cosmos-sdk/server/config"

	"github.com/tendermint/tendermint/libs/common"
)

// nolint
var (
	DefaultMaintainConfile = "maintain.conf"
	DefaultNodeHome        = okchaincfg.GetNodeHome()
	DefaultNodeCofig       = filepath.Join(DefaultNodeHome,"config")
	DefaultTestConfig      = filepath.Join(DefaultNodeHome, "test_config")
	DefaultTestDataHome    = filepath.Join( DefaultNodeHome, "test_data")

	DefaultConfig = okchaincfg.DefaultBackendConfig
)

// nolint
type Config = okchaincfg.BackendConfig

func loadMaintainConf(confDir string, fileName string) (*Config, error) {
	fPath := confDir + string(os.PathSeparator) + fileName
	if _, err := os.Stat(fPath); err != nil {
		return nil, err
	}

	bytes := common.MustReadFile(fPath)

	m := Config{}
	err := json.Unmarshal(bytes, &m)
	return &m, err
}

func dumpMaintainConf(maintainConf *Config, confDir string, fileName string) (err error) {
	fPath := confDir + string(os.PathSeparator) + fileName

	if _, err := os.Stat(confDir); err != nil {
		if err = os.MkdirAll(confDir, os.ModePerm); err != nil {
			return err
		}
	}

	bs, err := json.MarshalIndent(maintainConf, "", "  ")
	if err != nil {
		return err
	}
	common.MustWriteFile(fPath, bs, os.ModePerm)

	return nil
}

// nolint
func SafeLoadMaintainConfig(configDir string) (conf *Config, err error) {
	maintainConf, err := loadMaintainConf(configDir, DefaultMaintainConfile)
	if maintainConf == nil || err != nil {
		maintainConf = DefaultConfig()
		if err = dumpMaintainConf(maintainConf, configDir, DefaultMaintainConfile); err != nil {
			return nil, err
		}
	}
	return maintainConf, nil
}
