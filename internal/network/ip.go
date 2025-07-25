package network

import (
	"fmt"
	"net"
)

// GetLocalIP retrieves the local IP address of the machine.
func GetLocalIP() string {
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("🌐 Your LAN IP is:", ipnet.IP.String())
				return ipnet.IP.String()
			}
		}
	}

	fmt.Println("❌ If you are using a VPN, please disable it to share your dev environment over LAN.")
	return ""
}
