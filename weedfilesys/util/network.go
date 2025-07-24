package util

import (
	"net"
	"strconv"
	"strings"
)

func JoinHostPort(host string, port int) string {
	portStr := strconv.Itoa(port)
	if strings.HasPrefix(host, "[") && strings.HasSuffix(host, "]") {
		return host + ":" + portStr
	}
	return net.JoinHostPort(host, portStr)
}
