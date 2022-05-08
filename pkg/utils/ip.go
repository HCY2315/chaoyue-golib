package utils

import (
	"net"
)

// LocalIPs return all non-loopback IPv4 addresses
func LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil && !ipnet.IP.IsLoopback() && !ipnet.IP.IsUnspecified() {
			ips = append(ips, ipnet.IP.String())
		}
	}
	return ips, nil
}

func IsLocalIPv4(ip string) (bool, error) {
	localIps, err := LocalIPv4s()
	if err != nil {
		return false, err
	}
	ipInconsistent := true
	for _, localIp := range localIps {
		if ip == localIp {
			ipInconsistent = false
			break
		}
	}
	if ipInconsistent {
		return false, nil
	}
	return true, nil
}
