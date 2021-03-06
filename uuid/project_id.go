package uuid

import (
	"fmt"
	"math/rand"
	"net"
	"trace"
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

func getIdentifier() (identifier string) {
	identifier = trace.Config.Server.ModuleName
	// ip
	if trace.Config.Server.TraceId.ProjectIdStrategy == 1 {
		ips := getIPs(1)
		if len(ips) > 0 {
			identifier += "_" + ips[0]
		}
	}
	return identifier
}

func checkProjectId(projectId uint16) bool {
	return true
}

func applyProjectId() uint16 {
	return uint16(rand.Intn(MaxProjectId) + 1)
}

func GetProjectId() uint16 {
	check := trace.Config.Server.TraceId.ProjectIdCheck
	projectId := trace.Config.Server.TraceId.ProjectId

	if projectId > 0 {
		if check == 0 || checkProjectId(projectId) {
			return projectId
		}
	}

	return applyProjectId()
}
