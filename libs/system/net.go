package system

import (
	"fmt"
	"net"
)

func GetIpAddr(appendPid bool) (res string, err error) {

	var addrs []net.Addr
	addrs, err = net.InterfaceAddrs()
	if err != nil{
		return
	}
	var comma string
	for _, value := range addrs{
		if ipnet, ok := value.(*net.IPNet); ok && !ipnet.IP.IsLoopback(){
			if ipnet.IP.To4() != nil{
				res += fmt.Sprintf("%s%s", comma, ipnet.IP.String())
				comma = "_"
			}
		}
	}

	if appendPid {
		res = fmt.Sprintf("%s_%d", res, Getpid())
	}
	return
}
