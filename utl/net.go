package utl

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func GetIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, a := range addrs {
		addr, ok := a.(*net.IPNet)
		if !ok {
			continue
		}
		if addr.IP.IsLoopback() || addr.IP.To4() == nil {
			continue
		}
		return addr.IP.String(), nil
	}
	return "", nil
}

func GetMac() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range interfaces {
		if HasPrefix(i.Name, "en") {
			return i.HardwareAddr.String(), nil
		}
	}
	return "", nil
}

func GetPort(min, max int) int {
	for i := 0; i < 5; i++ {
		port := min + rand.Intn(max-min)
		_, err := net.DialTimeout("udp", fmt.Sprintf(":%d", port), 5*time.Second)
		if err != nil {
			return port
		}
		continue
	}
	return 0
}
