package main

import "net"

func sendArpPackage(ip IP) {
	srcIP := net.ParseIP(ipNet.IP.String()).To4()
	dstIP := net.ParseIP(ip.String()).To4()

	if srcIP == nil || dstIP == nil {
		log.Fatal("IP 解析出现问题")
	}

}
