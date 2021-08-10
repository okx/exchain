package config

import (
	"fmt"

	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
)

func NewApollo(log log.Logger) {
	c := &config.AppConfig{
		AppID:          "okexchain",
		Cluster:        "dev",
		IP:             "http://service-apollo-config-server-dev.apollo-dev.svc.base.local:8080",
		NamespaceName:  "rpc-node",
		IsBackupConfig: true,
		//Secret:         "",
	}

	//agollo.SetLogger(&Logger{log})

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("初始化Apollo配置成功")

	//Use your apollo key to test
	cache := client.GetConfigCache(c.NamespaceName)
	fmt.Println(cache.EntryCount())
	cache.Range(func(key, value interface{}) bool {
		fmt.Println("key : ", key, ", value :", value)
		return true
	})
	//ss:=client.GetConfig(c.NamespaceName)

	c2 := &CustomChangeListener{}
	client.AddChangeListener(c2)
}

type CustomChangeListener struct {
}

func (c *CustomChangeListener) OnChange(changeEvent *storage.ChangeEvent) {
	fmt.Println("OnChange", changeEvent.Changes)
	for key, value := range changeEvent.Changes {
		fmt.Println("change key : ", key, ", value :", value.NewValue)
	}
	fmt.Println("OnChange", changeEvent.Namespace)
}

func (c *CustomChangeListener) OnNewestChange(event *storage.FullChangeEvent) {
	fmt.Println("OnNewestChange", event.Changes)
	//for key, value := range changeEvent.Changes {
	//	fmt.Println("change key : ", key, ", value :", value)
	//}
	fmt.Println("OnNewestChange", event.Namespace)

}

type Logger struct {
	logger log.Logger
}

func (l *Logger) Debugf(format string, params ...interface{}) {
	l.logger.Debug(format, params)
}

func (l *Logger) Infof(format string, params ...interface{}) {
	l.logger.Info(format, params)
}

func (l *Logger) Warnf(format string, params ...interface{}) {
	return
}

func (l *Logger) Errorf(format string, params ...interface{}) {
	l.logger.Error(format, params)
}

func (l *Logger) Debug(v ...interface{}) {
	l.logger.Debug("", v)
}

func (l *Logger) Info(v ...interface{}) {
	l.logger.Info("", v)
}

func (l *Logger) Warn(v ...interface{}) {
	return
}

func (l *Logger) Error(v ...interface{}) {
	l.logger.Error("", v)
}
