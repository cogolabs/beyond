package main

import (
	"log"
	"net"
	"time"
)

func init() {
	// prepend file:lineno
	log.SetFlags(log.Flags() | log.Lshortfile)

	// wait for networking
	for i := 0; i < 5; i++ {
		addrs, _ := net.InterfaceAddrs()
		for _, addr := range addrs {
			if a, ok := addr.(*net.IPNet); ok && a.IP.IsGlobalUnicast() {
				return
			}
		}
		time.Sleep(time.Second)
	}
}
