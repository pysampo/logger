package util

import (
	"sync"

	tnet "github.com/toolkits/net"
)

var (
	once     sync.Once
	clientIP = "127.0.0.1"
)

func GetLocalIP() string {
	once.Do(func() {
		ips, _ := tnet.IntranetIP()
		if len(ips) > 0 {
			clientIP = ips[0]
		} else {
			clientIP = "127.0.0.1"
		}
	})
	return clientIP
}
