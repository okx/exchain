package monitor

import (
	"bytes"
	"fmt"
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
		// TODO: add config and cmd flag
		// p2p:26656, rpc:26657, rest:26659
		portMonitor = NewPortMonitor([]string{"26656", "26657", "26659"})
	})

	return portMonitor
}

// PortMonitor - structure of monitor for ports
type PortMonitor struct {
	ports                   []uint64
	maxConnectingNumber     int
	currentConnectingNumber int
	connectingMap           map[uint64]int
}

// NewPortMonitor creates a new instance of PortMonitor
func NewPortMonitor(ports []string) *PortMonitor {
	// check port format
	var portsInt []uint64
	for _, portStr := range ports {
		n, err := strconv.ParseUint(strings.TrimSpace(portStr), 10, 64)
		if err != nil {
			panic(fmt.Sprintf("fail to convert port string %s to integer: %s", portStr, err.Error()))
		}

		if n > 65535 {
			panic(fmt.Sprintf("invalid port %d. It should be between 0 and 65535", n))
		}

		portsInt = append(portsInt, n)
	}

	return &PortMonitor{
		ports:                   portsInt,
		connectingMap:           make(map[uint64]int),
		currentConnectingNumber: -1,
		maxConnectingNumber:     -1,
	}
}

// reset resets the status of PortMonitor
func (pm *PortMonitor) reset() {
	for _, port := range pm.ports {
		pm.connectingMap[port] = -1
	}

	pm.currentConnectingNumber = -1
}

// getConnectingNumbers gets the connecting numbers from ports
func (pm *PortMonitor) getConnectingNumbers() error {
	var connectingNumTotal int
	for _, port := range pm.ports {
		connectingNumber, err := getConnectingNumbersFromPort(port)
		if err != nil {
			return fmt.Errorf("failed to get connecting numbers of port %d: %s", port, err.Error())
		}
		pm.connectingMap[port] = connectingNumber
		connectingNumTotal += connectingNumber
	}

	pm.currentConnectingNumber = connectingNumTotal

	// max check
	if connectingNumTotal > pm.maxConnectingNumber {
		pm.maxConnectingNumber = connectingNumTotal
	}
	return nil
}

func (pm *PortMonitor) Run() error {
	// PortMonitor disabled
	if len(pm.ports) == 0 {
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
	var buffer bytes.Buffer
	for _, port := range pm.ports {
		buffer.WriteString(fmt.Sprintf("%d<%d>, ", port, pm.connectingMap[port]))
	}

	// statistics
	buffer.WriteString(fmt.Sprintf("CurConNum<%d>, MaxConNum<%d>", pm.currentConnectingNumber, pm.maxConnectingNumber))
	return buffer.String()
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
