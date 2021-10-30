package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	okexchaincfg "github.com/okex/exchain/dependence/cosmos-sdk/server/config"
)

// nolint
var (
	DefaultMaintainConfile = "maintain.conf"
	DefaultNodeHome        = okexchaincfg.GetNodeHome()
	DefaultNodeCofig       = filepath.Join(DefaultNodeHome, "config")
	DefaultTestConfig      = filepath.Join(DefaultNodeHome, "test_config")
	DefaultTestDataHome    = filepath.Join(DefaultNodeHome, "test_data")
	DefaultConfig          = okexchaincfg.DefaultBackendConfig
)

// nolint
type Config = okexchaincfg.BackendConfig

func loadMaintainConf(confDir string, fileName string) (*Config, error) {
	fPath := confDir + string(os.PathSeparator) + fileName
	if _, err := os.Stat(fPath); err != nil {
		return nil, err
	}

	bytes := mustReadFile(fPath)

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
	mustWriteFile(fPath, bs, os.ModePerm)

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

func mustReadFile(filePath string) []byte {
	fileBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf(fmt.Sprintf("mustReadFile failed: %v\n", err))
		os.Exit(1)
		return nil
	}
	return fileBytes
}

func mustWriteFile(filePath string, contents []byte, mode os.FileMode) {
	err := ioutil.WriteFile(filePath, contents, mode)
	if err != nil {
		fmt.Printf(fmt.Sprintf("mustWriteFile failed: %v\n", err))
		os.Exit(1)
	}
}
