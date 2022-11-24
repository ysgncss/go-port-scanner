package util

import (
	"errors"
	"github.com/go-ping/ping"
	"goroutine-example/constant"
	"log"
	"net"
	"strings"
	"time"
)

func Ping(target string) bool {
	pingServer, err := ping.NewPinger(target)
	if err != nil {
		log.Printf(err.Error())
		return false
	}

	pingServer.Count = constant.IcmpCount
	pingServer.Timeout = time.Duration(constant.PingTime * time.Millisecond)
	pingServer.SetPrivileged(true)
	pingServer.Run() // blocks until finished
	stats := pingServer.Statistics()

	// log.Print(stats)
	// 有回包，就是说明IP是可用的
	if stats.PacketsRecv >= 1 {
		return true
	}
	return false
}

// ValidIpv4 判断ip是否是真实IPV4
func ValidIpv4(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, p := range parts {
		if len(p) == 0 {
			return false
		}
		if p[0] == '0' && len(p) > 1 {
			return false
		}
		num := 0
		for _, c := range p {
			if !('0' <= c && c <= '9') {
				return false
			}
			num = num*10 + int(c-'0')
		}
		if !(0 <= num && num <= 255) {
			return false
		}
	}
	return true
}

// ValidIpv6 判断ip是否是真实IPV6
func ValidIpv6(ip string) bool {
	parts := strings.Split(ip, ":")
	if len(parts) != 8 {
		return false
	}
	for _, p := range parts {
		if !(1 <= len(p) && len(p) <= 4) {
			return false
		}
		for _, c := range p {
			if !(('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')) {
				return false
			}
		}
	}
	return true
}

// GetHostIP 获取本机IP地址
func GetHostIP() (net.IP, error) {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, a := range address {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP, nil
			}
		}
	}
	return nil, errors.New("no valid ipv4 address founded")
}

// GetCidrIpRange 根据网段识别出所有IP地址，"cidr" 代表网段
func GetCidrIpRange(cidr string) []string {
	var ips []string

	ipAddr, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Print(err)
		return ips
	}

	if ipAddr == nil {
		log.Printf("empty parsed ip address with input: %s", cidr)
		return ips
	}

	for ip := ipAddr.Mask(ipNet.Mask); ipNet.Contains(ip); increment(ip) {
		ips = append(ips, ip.String())
	}

	return ips
}

func increment(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] != 0 {
			break
		}
	}
}
