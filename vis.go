package main

import (
	"fmt"
	"net"
	"sort"
	"sync"
)

type DisplayInfo struct {
	// 地址
	Mac net.HardwareAddr
	// 主机名
	Hostname string
	// 厂商信息
	Manuf string
}

// 存放最终的数据，key[IP] 存放的是IP地址
var data map[IP]DisplayInfo
var mu sync.RWMutex

func pushData(ip IP, mac net.HardwareAddr, hostname, manuf string) {
	mu.RLock()

	if _, ok := data[ip]; !ok {
		data[ip] = DisplayInfo{Mac: mac, Hostname: hostname, Manuf: manuf}
	}

	info := data[ip]
	if len(hostname) > 0 && len(info.Hostname) == 0 {
		info.Hostname = hostname
	}

	if len(manuf) > 0 && len(info.Manuf) == 0 {
		info.Manuf = manuf
	}
	if mac != nil {
		info.Mac = mac
	}
	data[ip] = info

	mu.RUnlock()
}

// 格式化输出结果
// xxx.xxx.xxx.xxx  xx:xx:xx:xx:xx:xx  hostname  manuf
// xxx.xxx.xxx.xxx  xx:xx:xx:xx:xx:xx  hostname  manuf
func PrintData() {
	var keys IPSlice
	for k := range data {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	for _, k := range keys {
		d := data[k]
		mac := ""
		if d.Mac != nil {
			mac = d.Mac.String()
		}
		fmt.Printf("%-15s %-17s %-30s %-10s\n", k.String(), mac, d.Hostname, d.Manuf)
	}
}
