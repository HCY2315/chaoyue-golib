package tcp

import (
	"net"
	"time"
)

func IsTCPAddrOpen(host, port string) bool {
	to := time.Second * 2
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), to)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
