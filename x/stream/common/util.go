package common

import (
	"net"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/viper"
)

func ResolveRestIPAndPort() (string, int, error) {
	laddr := strings.Split(viper.GetString(server.FlagExternalListenAddr), ":")
	ip := laddr[0]
	if ip == "127.0.0.1" {
		return GetLocalIP(), 26659, nil
	}
	port, err := strconv.Atoi(laddr[1])
	if err != nil {
		return "", 0, err
	}
	return ip, port, nil
}

// GetLocalIP get local ip
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
