package utils

import (
	"fmt"
	"net"
)

// PeerServiceDiscovery looks up dns records for this host name and returns ip concatenated with the port
func PeerServiceDiscovery(host string, port string) []string {
	addrs, _ := net.LookupHost(host)

	for i, addr := range addrs {
		addrs[i] = fmt.Sprintf("%s:%s", addr, port)
	}

	return addrs
}
