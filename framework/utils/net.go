package utils

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

//获取ip
func ExternalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			if strings.Contains(ip.String(), "192.168.1.") ||
				strings.Contains(ip.String(), "192.168.2.") {
				fmt.Println("IPSelect on:", ip.String())
				return ip.String(), nil
			}
			fmt.Println("IPSelect:", ip.String())
		}
	}
	return "", errors.New("connected to the network?")
}

//获取ip
func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
