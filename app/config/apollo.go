package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"

	"github.com/apolloconfig/agollo/v4"
	"github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
)

const FlagApollo = "config.apollo"

type ApolloClient struct {
	Namespace string
	*agollo.Client
	oecConf *OecConfig
}

func NewApolloClient(oecConf *OecConfig) *ApolloClient {
	// IP|Cluster|AppID|NamespaceName
	params := strings.Split(viper.GetString(FlagApollo), "|")
	if len(params) != 4 {
		panic("failed init apollo: invalid connection config")
	}

	c := &config.AppConfig{
		IP:             params[0],
		Cluster:        params[1],
		AppID:          params[2],
		NamespaceName:  params[3],
		IsBackupConfig: false,
		//Secret:         "",
	}

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		panic(fmt.Errorf("failed init apollo: %v", err))
	}

	apc := &ApolloClient{
		params[3],
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
	cache.Clear()
}

func (c *OecConfig) OnChange(changeEvent *storage.ChangeEvent) {
	for key, value := range changeEvent.Changes {
		if value.ChangeType != storage.DELETED {
			c.update(key, value.NewValue)
		}
	}
}

func (c *OecConfig) OnNewestChange(event *storage.FullChangeEvent) {
	return
}
