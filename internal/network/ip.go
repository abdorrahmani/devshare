package network

import (
	"fmt"
	"net"
	"strings"
)

// GetLocalIP retrieves the local IP address of the machine.
func GetLocalIP() string {
	// Get all network interfaces
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("❌ Error getting network interfaces:", err)
		return ""
	}

	var candidates []string

	for _, iface := range interfaces {
		// Skip loopback and down interfaces
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				ip := ipnet.IP
				if ip.To4() != nil && !ip.IsLoopback() {
					ipStr := ip.String()

					// Skip router IPs
					if strings.HasSuffix(ipStr, ".1") {
						continue
					}

					// Skip link-local addresses
					if strings.HasPrefix(ipStr, "169.254.") {
						continue
					}

					// Only consider private network ranges
					if isPrivateIP(ip) {
						candidates = append(candidates, ipStr)
						fmt.Printf("🔍 Found candidate IP: %s (interface: %s)\n", ipStr, iface.Name)
					}
				}
			}
		}
	}

	if len(candidates) > 0 {
		selectedIP := candidates[0]
		fmt.Println("🌐 Your LAN IP is:", selectedIP)
		return selectedIP
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("❌ Error getting interface addresses:", err)
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip := ipnet.IP.String()
				if !strings.HasSuffix(ip, ".1") {
					fmt.Println("🌐 Your LAN IP is:", ip)
					return ip
				}
			}
		}
	}

	fmt.Println("❌ No suitable LAN IP found. If you are using a VPN, please disable it to share your dev environment over LAN.")
	return ""
}

// isPrivateIP checks if an IP is in a private network range
func isPrivateIP(ip net.IP) bool {
	// Private network ranges
	privateRanges := []struct {
		start, end net.IP
	}{
		{net.ParseIP("10.0.0.0"), net.ParseIP("10.255.255.255")},
		{net.ParseIP("172.16.0.0"), net.ParseIP("172.31.255.255")},
		{net.ParseIP("192.168.0.0"), net.ParseIP("192.168.255.255")},
	}

	for _, r := range privateRanges {
		if inRange(ip, r.start, r.end) {
			return true
		}
	}
	return false
}

// inRange checks if an IP is within a given range
func inRange(ip, start, end net.IP) bool {
	return bytesToInt(ip) >= bytesToInt(start) && bytesToInt(ip) <= bytesToInt(end)
}

// bytesToInt converts IP bytes to integer for comparison
func bytesToInt(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 + uint32(ip[1])<<16 + uint32(ip[2])<<8 + uint32(ip[3])
}
