package utils

import (
	"net"
	"strings"
)

// ConvertIPv6ToIPv4 takes an IPv6 address string and converts it to an IPv4 address string.
// If the input is already an IPv4 address, it returns the original input.
func ConvertIPv6ToIPv4(ipAddr string) string {
	// Split the address into host and port
	parts := strings.Split(ipAddr, ":")
	host := parts[0]
	port := ""
	if len(parts) > 1 {
		port = ":" + parts[1]
	}

	// Convert the host part to IPv4
	ip := net.ParseIP(host)
	if ip.To4() != nil {
		// The input is already an IPv4 address
		return ipAddr
	}

	// Convert the IPv6 address to IPv4
	ipv4 := net.IPv4(127, 0, 0, 1)
	return ipv4.String() + port
}
