package utils

import (
	"errors"
	"net"
	"os"
	"strings"
)

var LocalIP = "unknown"

func init() {
	if ip, err := localIP(); err == nil {
		LocalIP = ip
	}
}

func localIP() (string, error) {
	// 优先获取运维配置机器的环境变量
	hostIP := os.Getenv("HOST_IP")
	if len(hostIP) > 0 {
		return hostIP, nil
	}

	// 优先获取外网ip
	if conn, err := net.Dial("udp", "8.8.8.8:53"); err == nil {
		laddr := conn.LocalAddr().(*net.UDPAddr)
		ip := strings.Split(laddr.String(), ":")[0]
		return ip, nil
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}

	return "", errors.New("can not find the client ip address")
}
