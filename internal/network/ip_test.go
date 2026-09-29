package network

import (
	"net"
	"testing"
)

func TestIsPrivateIP(t *testing.T) {
	private := []string{"10.0.0.5", "10.255.255.254", "172.16.0.1", "172.31.255.1", "192.168.1.42"}
	public := []string{"8.8.8.8", "172.15.0.1", "172.32.0.1", "1.1.1.1", "203.0.113.5"}

	for _, s := range private {
		if !isPrivateIP(net.ParseIP(s)) {
			t.Errorf("%s should be detected as private", s)
		}
	}
	for _, s := range public {
		if isPrivateIP(net.ParseIP(s)) {
			t.Errorf("%s should be detected as public", s)
		}
	}
}
