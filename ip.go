package main

import (
	"fmt"
	"math"
	"net"
)

type IP uint32

// 实现 String 方法
func (ip IP) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// []byte --> IP
func ParseIP(b []byte) IP {
	return IP(IP(b[0])<<24 + IP(b[1])<<16 + IP(b[2])<<8 + IP(b[3]))
}

// /////////////////////////////////////////////////////////////////////////////////
// 根据IP和mask换算内网IP范围
func Table(ipNet *net.IPNet) []IP {
	ip := ipNet.IP.To4()
	log.Info("本机IP:", ip)
	var min, max IP
	var data []IP

	for i := 0; i < 4; i++ {
		b := IP(ip[i] & ipNet.Mask[i])
		//每次移动8bit
		// index顺序 0 1 2 3
		min += b << ((3 - uint(i)) * 8)
	}

	one, _ := ipNet.Mask.Size() //In 32 format
	max = min | IP((1<<(32-one))-1)
	max2 := min | IP(math.Pow(2, float64(32-one))-1)

	if max != max2 {
		log.Fatal("Assertion Failed")
	}

	log.Infof("内网IP范围:%s --- %s", min, max)

	// max 是广播地址，忽略
	// i & 0x000000ff  == 0 是尾段为0的IP，根据RFC的规定，忽略

	for i := min; i < max; i++ {
		if i&0x000000ff == 0 {
			continue
		}
		data = append(data, i)
	}
	return data
}
