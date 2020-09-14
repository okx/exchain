package eureka

import (
	"fmt"

	"github.com/okex/okexchain/x/stream/common"
	"github.com/tendermint/tendermint/libs/log"
)

// StartEurekaClient start eureka client and register rest service in eureka
func StartEurekaClient(logger log.Logger, url string, name string) {
	ip, port, err := common.ResolveRestIPAndPort()
	if err != nil {
		logger.Error(fmt.Sprintf("failed to resolve rest.external_laddr: %s", err.Error()))
		return
	}

	c := newClient(&eurekaConfig{
		serverURL:             url,
		appName:               name,
		appIP:                 ip,
		port:                  port,
		renewalIntervalInSecs: 30,
		durationInSecs:        90,
	})
	err = registerEurekaInstance(c, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("failed to register application instance in eureka server: %s", err.Error()))
		return
	}

	// sendHeartbeat
	go c.sendHeartbeat(logger)
	// handle signal, auto delete register when receive exit signal
	go c.handleSignal(logger)
}

func registerEurekaInstance(c *eurekaClient, logger log.Logger) error {
	c.mutex.Lock()
	c.running = true
	c.mutex.Unlock()
	// register
	err := register(c.instance, c.config.serverURL, c.config.appName)
	if err != nil {
		return err
	}
	logger.Info("register application instance in eureka successfully")
	return nil
}
