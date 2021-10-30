package eureka

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/okex/exchain/x/stream/common/utils"
	"github.com/okex/exchain/dependence/tendermint/libs/log"
)

// eurekaClient client for eureka
type eurekaClient struct {
	// for monitor system signal
	signalChan   chan os.Signal
	mutex        sync.RWMutex
	running      bool
	config       *eurekaConfig
	instance     *Instance
	applications *Applications // nolint
}

// eurekaConfig config for eureka
type eurekaConfig struct {
	serverURL             string // server url
	renewalIntervalInSecs int    // the heart-beat interval
	// nolint
	registryFetchIntervalSeconds int    // the fetching interval
	durationInSecs               int    // the expired time
	appName                      string // application name
	appIP                        string // application ip
	port                         int    // server port
	metadata                     map[string]interface{}
}

// sendHeartbeat
func (c *eurekaClient) sendHeartbeat(logger log.Logger) {
	for {
		if c.running {
			if err := heartbeat(c.config.serverURL, c.config.appName, c.instance.InstanceID); err != nil {
				logger.Error(fmt.Sprintf("failed to send heart-beat: %s", err.Error()))
			} else {
				logger.Debug("send heart-beat with application instance successfully")
			}
		} else {
			break
		}
		time.Sleep(time.Duration(c.config.renewalIntervalInSecs) * time.Second)
	}
}

// auto delete register when receive exit signal
func (c *eurekaClient) handleSignal(logger log.Logger) {
	if c.signalChan == nil {
		c.signalChan = make(chan os.Signal)
	}
	signal.Notify(c.signalChan, syscall.SIGTERM, syscall.SIGINT)
	for {
		switch <-c.signalChan {
		case syscall.SIGINT:
			fallthrough
		case syscall.SIGTERM:
			logger.Info("receive exit signal, client instance going to de-egister")
			err := unRegister(c.config.serverURL, c.config.appName, c.instance.InstanceID)
			if err != nil {
				logger.Error(fmt.Sprintf("failed to unregister the instance. error: %s", err.Error()))
			} else {
				logger.Info("unregister application instance successfully")
			}
			os.Exit(0)
		}
	}
}

// nolint
func (c *eurekaClient) refresh(logger log.Logger) {
	for {
		if c.running {
			if applications, err := getAllInstance(c.config.serverURL); err != nil {
				logger.Error(fmt.Sprintf("failed to refresh the instances from server. error: %s", err.Error()))
			} else {
				c.mutex.Lock()
				c.applications = applications
				c.mutex.Unlock()
				logger.Debug("refresh application instance successfully")
			}
		} else {
			break
		}
		time.Sleep(time.Duration(c.config.registryFetchIntervalSeconds) * time.Second)
	}
}

// newClient create a eureka-client
func newClient(config *eurekaConfig) *eurekaClient {
	initConfig(config)
	return &eurekaClient{config: config, instance: newInstance(config)}
}

func initConfig(config *eurekaConfig) {
	if config.serverURL == "" {
		config.serverURL = "http://localhost:8761/eureka"
	}
	if config.renewalIntervalInSecs == 0 {
		config.renewalIntervalInSecs = 30
	}
	if config.durationInSecs == 0 {
		config.durationInSecs = 90
	}
	if config.appName == "" {
		config.appName = "OKDEX-REST-SERVICE"
	} else {
		config.appName = strings.ToLower(config.appName)
	}
	if config.appIP == "" {
		config.appIP = utils.GetLocalIP()
	}
	if config.port == 0 {
		config.port = 80
	}
}
