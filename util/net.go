package util

import (
	"fmt"
	"net"
)

// PortIsUsed check port
func PortIsUsed(port uint) bool {
	tcpAddr, _ := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%d", port))
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)
	if err != nil {
		return false
	}
	defer func() { _ = conn.Close() }()
	return true
}

func BuildURL(https bool, host string, port uint) string {
	siteURL := "http"
	if https {
		siteURL += "s"
	}
	siteURL += "://"
	if host == "0.0.0.0" {
		siteURL += "127.0.0.1"
	} else {
		siteURL += host
	}
	if port != 80 && port != 443 {
		siteURL += fmt.Sprintf(":%d", port)
	}
	return siteURL
}
