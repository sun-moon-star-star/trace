package uuid

import (
	"fmt"
	"net"
)

func getIPs(len int) (ips []string) {
	if len == 0 {
		return nil
	}

	interfaceAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("fail to get net interface addrs: %v", err)
		return ips
	}

	cnt := 0
	for _, address := range interfaceAddr {
		ipNet, isValidIpNet := address.(*net.IPNet)
		if isValidIpNet && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			ips = append(ips, ipNet.IP.String())
			cnt++
			if cnt == len {
				return
			}
		}
	}
	return ips
}
