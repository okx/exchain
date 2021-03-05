package monitor

import (
	"bytes"
	"fmt"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/viper"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

var (
	portMonitor     *PortMonitor
	initPortMonitor sync.Once
)

// GetPortMonitor gets the global instance of PortMonitor
func GetPortMonitor() *PortMonitor {
	initPortMonitor.Do(func() {
		portMonitor = NewPortMonitor(parsePorts(viper.GetString(server.FlagPortMonitor)))
	})

	return portMonitor
}

// PortMonitor - structure of monitor for ports
type PortMonitor struct {
	enable bool
	ports  []uint64
	// max total connecting numbers in one round
	maxConnectingNumberTotal int
	// connecting number of each port in one round
	connectingMap map[uint64]int
	// max connecting number record of each port
	connectingMaxMap map[uint64]int
}

// NewPortMonitor creates a new instance of PortMonitor
func NewPortMonitor(ports []string) *PortMonitor {
	if len(ports) == 0 {
		// disable the port monitor
		return &PortMonitor{
			enable: false,
		}
	}
	// check port format
	var portsUint64 []uint64
	connectingMaxMap := make(map[uint64]int)
	for _, portStr := range ports {
		port := ParsePort(portStr)
		portsUint64 = append(portsUint64, port)
		// init connectingMaxMap with -1
		connectingMaxMap[port] = -1
	}

	return &PortMonitor{
		enable:                   true,
		ports:                    portsUint64,
		connectingMap:            make(map[uint64]int),
		connectingMaxMap:         connectingMaxMap,
		maxConnectingNumberTotal: -1,
	}
}

// reset resets the status of PortMonitor
func (pm *PortMonitor) reset() {
	for _, port := range pm.ports {
		pm.connectingMap[port] = -1
	}
}

// getConnectingNumbers gets the connecting numbers from ports
func (pm *PortMonitor) getConnectingNumbers() error {
	var connectingNumTotal int
	for _, port := range pm.ports {
		connectingNumber, err := getConnectingNumbersFromPort(port)
		if err != nil {
			return fmt.Errorf("failed to get connecting numbers of port %d: %s", port, err.Error())
		}

		// update max connecting map
		if connectingNumber > pm.connectingMaxMap[port] {
			pm.connectingMaxMap[port] = connectingNumber
		}

		// update connecting map for this round
		pm.connectingMap[port] = connectingNumber
		connectingNumTotal += connectingNumber
	}

	// max total check
	if connectingNumTotal > pm.maxConnectingNumberTotal {
		pm.maxConnectingNumberTotal = connectingNumTotal
	}
	return nil
}

// Run starts monitoring
func (pm *PortMonitor) Run() error {
	// PortMonitor disabled
	if !pm.enable {
		return nil
	}

	pm.reset()
	err := pm.getConnectingNumbers()
	if err != nil {
		return err
	}

	return nil
}

// GetResultString gets the format string result
func (pm *PortMonitor) GetResultString() string {
	// PortMonitor disabled
	if !pm.enable {
		return ""
	}

	var buffer bytes.Buffer

	// connecting number of each port in this round
	for _, port := range pm.ports {
		buffer.WriteString(fmt.Sprintf("%d<%d>, ", port, pm.connectingMap[port]))
	}

	// max connecting number of each port
	for _, port := range pm.ports {
		buffer.WriteString(fmt.Sprintf("%d-Max<%d>, ", port, pm.connectingMaxMap[port]))
	}

	// statistics
	buffer.WriteString(fmt.Sprintf("MaxConNum<%d>", pm.maxConnectingNumberTotal))
	return buffer.String()
}

//GetConnectingMap gets connectingMap
func (pm *PortMonitor) GetConnectingMap() map[uint64]int {
	return pm.connectingMap
}

// tools function
func getConnectingNumbersFromPort(port uint64) (int, error) {
	// get connecting number from a shell command running
	shellCmd := fmt.Sprintf("netstat -nat | grep -i %d | wc -l", port)
	resBytes, err := exec.Command("/bin/sh", "-c", shellCmd).Output()
	if err != nil {
		return -1, err
	}

	// data washing
	return strconv.Atoi(string(bytes.TrimSpace(resBytes)))
}

func parsePorts(inputStr string) []string {
	inputStr = strings.TrimSpace(inputStr)
	if len(inputStr) == 0 {
		// nothing input
		return nil
	}

	return strings.Split(inputStr, ",")
}

// ParsePort parses port into uint from a string
func ParsePort(inputStr string) uint64 {
	inputStr = strings.TrimSpace(inputStr)
	port, err := strconv.ParseUint(inputStr, 10, 64)
	if err != nil {
		panic(fmt.Sprintf("fail to convert port string %s to integer: %s", inputStr, err.Error()))
	}

	if port > 65535 {
		panic(fmt.Sprintf("invalid port %d. It should be between 0 and 65535", port))
	}

	return port
}
