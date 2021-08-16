package config

import (
	"fmt"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
)

type ApolloClient struct {
	Namespace string
	*agollo.Client
	oecConf *OecConfig
}

func NewApolloClient(oecConf *OecConfig) *ApolloClient {
	c := &config.AppConfig{
		AppID:          "okexchain",
		Cluster:        "dev",
		IP:             "http://service-apollo-config-server-dev.apollo-dev.svc.base.local:8080",
		NamespaceName:  "rpc-node",
		IsBackupConfig: true,
		//Secret:         "",
	}

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		panic(fmt.Errorf("failed init apollo: %v", err))
	}

	apc := &ApolloClient{
		"rpc-node",
		client,
		oecConf,
	}
	client.AddChangeListener(oecConf)

	return apc
}

func (a *ApolloClient) LoadConfig() {
	cache := a.GetConfigCache(a.Namespace)
	cache.Range(func(key, value interface{}) bool {
		a.oecConf.update(key, value)
		return true
	})
}

func (c *OecConfig) OnChange(changeEvent *storage.ChangeEvent) {
	for key, value := range changeEvent.Changes {
		c.update(key, value.NewValue)
	}
}

func (c *OecConfig) OnNewestChange(event *storage.FullChangeEvent) {
	return
}
