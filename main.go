package main

import (
	"fmt"
	"github.com/bytedance/gopkg/util/gopool"
	"goroutine-example/util"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

type HostPort struct {
	host []uint8
	port int
}

type HostPortList []HostPort

func (s HostPortList) Len() int {
	return len(s)
}
func (s HostPortList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s HostPortList) Less(i, j int) bool {
	if s[i].host[0] < s[j].host[0] {
		return true
	}

	if s[i].host[1] < s[j].host[1] {
		return true
	}

	if s[i].host[2] < s[j].host[2] {
		return true
	}

	if s[i].host[3] < s[j].host[3] {
		return true
	}

	return s[i].port < s[j].port
}

// Check 通过ping判断是否存活
func Check(ipRange []string) []string {
	var res []string
	if len(ipRange) == 0 {
		return res
	}

	// 检查是否存活
	var channel = make(chan string)
	for _, ip := range ipRange {
		go func() {
			if util.Ping(ip) {
				channel <- ip
			} else {
				channel <- ""
			}
		}()
	}

	for i := 0; i < len(ipRange); i++ {
		ip := <-channel
		if ip != "" {
			res = append(res, ip)
		}
	}

	return res
}

func main() {
	start := time.Now()
	hostPorts := make(chan string)
	results := make(chan string)
	// 关闭 chan
	defer close(hostPorts)
	defer close(results)

	var openHostPorts HostPortList
	var ipRange []string

	parse := InputParse()

	// 解析
	if util.ValidIpv4(parse.host) {
		ipRange = []string{parse.host}
	} else {
		ipRange = util.GetCidrIpRange(parse.host)
	}

	// 检查是否存活
	ipRange = Check(ipRange)
	totalHost := len(ipRange)
	if totalHost == 0 {
		log.Printf("get empty ip range: %v", ipRange)
		return
	}

	log.Printf("get ip range: %v", ipRange)
	totalPort := (parse.endPort - parse.startPort + 1) * totalHost

	// 多协程扫描
	scannerPool := gopool.NewPool("scannerPool", int32(parse.poolCap), gopool.NewConfig())

	// 写入待扫描的主机加端口号
	scannerPool.Go(func() {
		for p := parse.startPort; p <= parse.endPort; p++ {
			for _, ip := range ipRange {
				address := fmt.Sprintf("%s:%d", ip, p)
				hostPorts <- address
			}
		}
	})

	//// 启动worker
	for i := 0; i < parse.poolCap; i++ {
		scannerPool.Go(func() {
			for address := range hostPorts {
				conn, err := net.DialTimeout(parse.network, address, time.Second*2)
				if err != nil {
					fmt.Printf("get err: %s  port: %s \n", err.Error(), address)
					results <- ""
					continue
				}
				fmt.Printf("get open port: %s\n", address)
				conn.Close()
				results <- address
			}
		})
	}

	// 接收返回的数据
	for i := 0; i < totalPort; i++ {
		hostPort := <-results
		if hostPort != "" {
			split := strings.Split(hostPort, ":")
			hosts := strings.Split(split[0], ".")
			host0, _ := strconv.Atoi(hosts[0])
			host1, _ := strconv.Atoi(hosts[1])
			host2, _ := strconv.Atoi(hosts[2])
			host3, _ := strconv.Atoi(hosts[3])
			port, _ := strconv.Atoi(split[1])
			openHostPorts = append(openHostPorts, HostPort{
				host: []uint8{uint8(host0), uint8(host1), uint8(host2), uint8(host3)},
				port: port,
			})
		}
	}

	// 排序结果
	sort.Sort(openHostPorts)
	for _, hostPort := range openHostPorts {
		fmt.Printf("%d.%d.%d.%d:%d open!\n", hostPort.host[0], hostPort.host[1], hostPort.host[2], hostPort.host[3], hostPort.port)
	}

	// 计算时间
	tc := time.Since(start)
	fmt.Printf("finished total cost: %v \n", tc)
}
