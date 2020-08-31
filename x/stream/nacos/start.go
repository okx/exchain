package nacos

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/okex/okchain/x/stream/common"
	"github.com/tendermint/tendermint/libs/log"
)

// StartNacosClient start eureka client and register rest service in eureka
func StartNacosClient(logger log.Logger, urls string, namespace string, name string) {
	ip, port, err := common.ResolveRestIpAndPort()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to resolve rest.external_laddr: %s", err.Error()))
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
			LogDir:              "/dev/null",
			NamespaceId:         namespace,
			//Username:			 "nacos",
			//Password:			 "nacos",
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
		ClusterName: "a",
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"preserved.register.source": "GO",
		},
	})
	if err != nil {
		logger.Error(fmt.Sprintf("failed to register instance in nacos server. error: %s", err.Error()))
		return
	}
	logger.Info("register application instance in nacos successfully")
}

func getServerConfigs(urls string) ([]constant.ServerConfig, error) {
	var configs []constant.ServerConfig
	for _, url := range strings.Split(urls, ",") {
		laddr := strings.Split(url, ":")
		serverPort, err := strconv.Atoi(laddr[1])
		if err != nil {
			return nil, err
		}
		configs = append(configs, constant.ServerConfig{
			IpAddr: laddr[0],
			Port:   uint64(serverPort),
		})
	}
	return configs, nil
}
