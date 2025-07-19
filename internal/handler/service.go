package handler

import "net"

func isInSubNet(address string, cidr string) bool {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	ip := net.ParseIP(address)
	return ipNet.Contains(ip)
}
