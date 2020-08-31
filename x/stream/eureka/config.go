package eureka

import (
	"encoding/xml"
	"fmt"
)

type ApplicationsInfo struct {
	Applications *Applications `json:"applications"`
}

type Applications struct {
	XMLName         xml.Name      `xml:"applications"`
	VersionsDelta   string        `xml:"versions__delta,omitempty"`
	AppsHashcode    string        `xml:"apps__hashcode,omitempty"`
	ApplicationList []Application `xml:"application"`
}

type Application struct {
	XMLName   xml.Name   `xml:"application"`
	Name      string     `xml:"name"`
	Instances []Instance `xml:"instance"`
}

type InstanceInfo struct {
	Instance *Instance `json:"instance"`
}

type Instance struct {
	XMLName                       xml.Name               `xml:"instance"`
	HostName                      string                 `xml:"hostName" json:"hostName"`
	HomePageURL                   string                 `xml:"homePageUrl,omitempty" json:"homePageUrl,omitempty"`
	StatusPageURL                 string                 `xml:"statusPageUrl" json:"statusPageUrl"`
	HealthCheckURL                string                 `xml:"healthCheckUrl,omitempty" json:"healthCheckUrl,omitempty"`
	App                           string                 `xml:"app" json:"app"`
	IPAddr                        string                 `xml:"ipAddr" json:"ipAddr"`
	VipAddress                    string                 `xml:"vipAddress" json:"vipAddress"`
	SecureVipAddress              string                 `xml:"secureVipAddress,omitempty" json:"secureVipAddress,omitempty"`
	Status                        string                 `xml:"status" json:"status"`
	Port                          *Port                  `xml:"port,omitempty" json:"port,omitempty"`
	SecurePort                    *Port                  `xml:"securePort,omitempty" json:"securePort,omitempty"`
	DataCenterInfo                *DataCenterInfo        `xml:"dataCenterInfo" json:"dataCenterInfo"`
	LeaseInfo                     *LeaseInfo             `xml:"leaseInfo,omitempty" json:"leaseInfo,omitempty"`
	Metadata                      map[string]interface{} `xml:"-" json:"metadata,omitempty"`
	IsCoordinatingDiscoveryServer string                 `xml:"isCoordinatingDiscoveryServer,omitempty" json:"isCoordinatingDiscoveryServer,omitempty"`
	LastUpdatedTimestamp          string                 `xml:"lastUpdatedTimestamp,omitempty" json:"lastUpdatedTimestamp,omitempty"`
	LastDirtyTimestamp            string                 `xml:"lastDirtyTimestamp,omitempty" json:"lastDirtyTimestamp,omitempty"`
	ActionType                    string                 `xml:"actionType,omitempty" json:"actionType,omitempty"`
	OverriddenStatus              string                 `xml:"overriddenstatus,omitempty" json:"overriddenstatus,omitempty"`
	CountryID                     int                    `xml:"countryId,omitempty" json:"countryId,omitempty"`
	InstanceID                    string                 `xml:"instanceId,omitempty" json:"instanceId,omitempty"`
}

// Port 端口
type Port struct {
	Port    int    `xml:",chardata" json:"$"`
	Enabled string `xml:"enabled,attr" json:"@enabled"`
}

// DataCenterInfo 数据中心信息
type DataCenterInfo struct {
	Name     string              `xml:"name" json:"name"`
	Class    string              `xml:"class,attr" json:"@class"`
	Metadata *DataCenterMetadata `xml:"metadata,omitempty" json:"metadata,omitempty"`
}

// DataCenterMetadata 数据中心信息元数据
type DataCenterMetadata struct {
	AmiLaunchIndex   string `xml:"ami-launch-index,omitempty" json:"ami-launch-index,omitempty"`
	LocalHostname    string `xml:"local-hostname,omitempty" json:"local-hostname,omitempty"`
	AvailabilityZone string `xml:"availability-zone,omitempty" json:"availability-zone,omitempty"`
	InstanceID       string `xml:"instance-id,omitempty" json:"instance-id,omitempty"`
	PublicIpv4       string `xml:"public-ipv4,omitempty" json:"public-ipv4,omitempty"`
	PublicHostname   string `xml:"public-hostname,omitempty" json:"public-hostname,omitempty"`
	AmiManifestPath  string `xml:"ami-manifest-path,omitempty" json:"ami-manifest-path,omitempty"`
	LocalIpv4        string `xml:"local-ipv4,omitempty" json:"local-ipv4,omitempty"`
	Hostname         string `xml:"hostname,omitempty" json:"hostname,omitempty"`
	AmiID            string `xml:"ami-id,omitempty" json:"ami-id,omitempty"`
	InstanceType     string `xml:"instance-type,omitempty" json:"instance-type,omitempty"`
}

// LeaseInfo 续约信息
type LeaseInfo struct {
	RenewalIntervalInSecs int `xml:"renewalIntervalInSecs,omitempty" json:"renewalIntervalInSecs,omitempty"`
	DurationInSecs        int `xml:"durationInSecs,omitempty" json:"durationInSecs,omitempty"`
}

// newInstance 创建服务实例
func newInstance(config *eurekaConfig) *Instance {
	instance := &Instance{
		InstanceID: fmt.Sprintf("%s:%s:%d", config.appIp, config.appName, config.port),
		HostName:   config.appIp,
		App:        config.appName,
		IPAddr:     config.appIp,
		Port: &Port{
			Port:    config.port,
			Enabled: "true",
		},
		VipAddress:       config.appName,
		SecureVipAddress: config.appName,
		// 续约信息
		LeaseInfo: &LeaseInfo{
			RenewalIntervalInSecs: config.renewalIntervalInSecs,
			DurationInSecs:        config.durationInSecs,
		},
		Status:           "UP",
		OverriddenStatus: "UNKNOWN",
		// 数据中心
		DataCenterInfo: &DataCenterInfo{
			Name:  "MyOwn",
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
		},
		// 元数据
		Metadata: config.metadata,
	}
	instance.HomePageURL = fmt.Sprintf("http://%s:%d", config.appIp, config.port)
	instance.StatusPageURL = fmt.Sprintf("http://%s:%d/info", config.appIp, config.port)
	return instance
}
