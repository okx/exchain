package nacos

import (
	"fmt"
	"strconv"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/okex/exchain/libs/tendermint/libs/log"
	"github.com/okex/exchain/x/stream/common/utils"
)

// StartNacosClient start nacos client and register rest service in nacos
func StartNacosClient(logger log.Logger, urls string, namespace string, name string, externalAddr string) {
	ip, port, err := utils.ResolveIPAndPort(externalAddr)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to resolve %s error: %s", externalAddr, err.Error()))
		return
	}

	serverConfigs, err := getServerConfigs(urls)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to resolve nacos server url %s: %s", urls, err.Error()))
		return
	}
	client, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig": constant.ClientConfig{
			TimeoutMs:           5000,
			ListenInterval:      10000,
			NotLoadCacheAtStart: true,
			NamespaceId:         namespace,
			LogDir:              "/dev/null",
			LogLevel:            "error",
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("failed to create nacos client. error: %s", err.Error()))
		return
	}

	_, err = client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(port),
		ServiceName: name,
		Weight:      10,
		ClusterName: "DEFAULT",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"preserved.register.source": "GO",
			"app_registry_tag":          strconv.FormatInt(time.Now().Unix(), 10),
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("failed to register instance in nacos server. error: %s", err.Error()))
		return
	}
	logger.Info("register application instance in nacos successfully")
}
