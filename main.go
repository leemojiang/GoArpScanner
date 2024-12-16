package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// ipNet 存放 IP地址和子网掩码
var ipNet *net.IPNet

// 本机的mac地址，发以太网包需要用到
var localHaddr net.HardwareAddr

// Net interface Name
var ifaceName string
var pcapName string
var ifaceIdx int

func init() {
	log.SetLevel(logrus.InfoLevel)
	formatter := &logrus.TextFormatter{
		ForceColors: true,
	}
	log.SetFormatter(formatter)

	flag.IntVar(&ifaceIdx, "I", 16, "Net Interface index")

	flag.Parse()
	log.Info("Program Start")

	// allow non root user to execute by compare with euid
	// if os.Geteuid() != 0 {
	// 	log.Fatal("goscan must run as root.")
	// }
}

func main() {
	setupNetInfo(ifaceIdx)
	setupPCAPInfo()

	// Send ARPs on Ip addr list
	for _, ipAddr := range Table(ipNet) {
		fmt.Println(ipAddr)
		sendArpPackage(ipAddr)
		break
	}

}

func setupNetInfo(idx int) {
	iface, err := net.InterfaceByIndex(idx)
	if err != nil {
		log.Fatal("无法获取此本地网络:", err)
	}
	addrs, _ := iface.Addrs()

	for _, a := range addrs {
		if ip, ok := a.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				ipNet = ip
				localHaddr = iface.HardwareAddr
				ifaceName = iface.Name

				log.Info("网络IPv4地址:", ip)
				log.Info("MAC地址:", localHaddr)
				log.Info("网络名称:", ifaceName)

			}
		}

	}

	// END:
	if ipNet == nil || len(localHaddr) == 0 {
		log.Fatal("无法获取本地网络信息")
	}
}

func showNetInfo() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Error(err)
	}

	for _, a := range addrs {
		if ip, ok := a.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				log.Info("IP:", ip.IP)
				log.Info("子网掩码:", ip.Mask)
				// it, _ := net.InterfaceByIndex(i)
				// log.Info("MAC:", it.HardwareAddr.String())
			}
		}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Error(err)
	}

	for _, iface := range interfaces {
		log.Info("Interface name:", iface.Name)
		log.Info("Interface MAC:", iface.HardwareAddr.String())
		log.Info("Interface Index:", iface.Index)
	}
}

func showPCAPInfo() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	for _, device := range devices {
		fmt.Printf("Name: %s, Description: %s ,Address: %s\n", device.Name, device.Description, device.Addresses)
	}
}

func setupPCAPInfo() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	for _, device := range devices {
		if device.Description == "WAN Miniport (IP)" {
			pcapName = device.Name
			break
		}
	}

	if pcapName == "" {
		log.Fatal("未找到合适的设备")
	}

	// 使用设备名称打开网卡
	handle, err := pcap.OpenLive(pcapName, 2048, true, pcap.BlockForever)
	if err != nil {
		log.Fatal("pcap打开失败:", err)
	}
	defer handle.Close()
}
