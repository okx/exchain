package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/okx/okbchain/libs/tendermint/libs/log"
)

const (
	LocalDynamicConfigPath = "config.dynamic.json"
)

type LocalClient struct {
	path     string
	dir      string
	okbcConf *OkbcConfig
	logger   log.Logger
	watcher  *fsnotify.Watcher
	close    chan struct{}
}

func NewLocalClient(path string, okbcConf *OkbcConfig, logger log.Logger) (*LocalClient, error) {
	if logger == nil {
		logger = log.NewNopLogger()
	}
	dir := filepath.Dir(path)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	client := &LocalClient{
		path:     path,
		dir:      dir,
		okbcConf: okbcConf,
		logger:   logger,
		watcher:  watcher,
		close:    make(chan struct{}),
	}
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// logger.Debug("local config event", "event", event)
				if event.Name == client.path && (event.Has(fsnotify.Write) || event.Has(fsnotify.Create)) {
					logger.Debug("local config changed", "path", path)
					ok = client.LoadConfig()
					if !ok {
						logger.Debug("local config changed but failed to load")
					} else {
						logger.Debug("local config changed and loaded")
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				logger.Error("local config watcher error", "err", err)
			case <-client.close:
				logger.Debug("local client closed")
				return
			}
		}
	}()

	return client, nil
}

func (a *LocalClient) Close() error {
	close(a.close)
	return a.watcher.Close()
}

func (a *LocalClient) Enable() (err error) {
	return a.watcher.Add(a.dir)
}

func (a *LocalClient) configExists() bool {
	_, err := os.Stat(a.path)
	return !os.IsNotExist(err)
}

func (a *LocalClient) LoadConfig() (loaded bool) {
	var conf map[string]string
	bz, err := os.ReadFile(a.path)
	if err != nil {
		a.logger.Error("failed to read local config", "path", a.path, "err", err)
		return false
	}
	err = json.Unmarshal(bz, &conf)
	if err != nil {
		a.logger.Error("failed to unmarshal local config", "path", a.path, "err", err)
		return false
	}
	loaded = true
	for k, v := range conf {
		a.okbcConf.updateFromKVStr(k, v)
	}
	a.logger.Info(a.okbcConf.format())
	return
}
