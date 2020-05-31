package wait

import (
	"net"
	"time"
)

func ForNetwork(tries int, interval time.Duration) {
	// wait for networking
	for i := 0; i < tries; i++ {
		addrs, _ := net.InterfaceAddrs()
		for _, addr := range addrs {
			if a, ok := addr.(*net.IPNet); ok && a.IP.IsGlobalUnicast() {
				return
			}
		}
		time.Sleep(interval)
	}
}
