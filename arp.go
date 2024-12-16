package main

import (
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func sendArpPackage(ip IP) {
	srcIP := net.ParseIP(ipNet.IP.String()).To4()
	dstIP := net.ParseIP(ip.String()).To4()

	if srcIP == nil || dstIP == nil {
		log.Fatal("IP 解析出现问题")
	}

	// 以太网首部
	// EthernetType 0x0806  ARP
	ether := &layers.Ethernet{
		SrcMAC:       localHaddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	a := &layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     uint8(6),
		ProtAddressSize:   uint8(4),
		Operation:         uint16(1), // 0x0001 arp request 0x0002 arp response
		SourceHwAddress:   localHaddr,
		SourceProtAddress: srcIP,
		DstHwAddress:      net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		DstProtAddress:    dstIP,
	}

	buffer := gopacket.NewSerializeBuffer()
	var opt gopacket.SerializeOptions
	gopacket.SerializeLayers(buffer, opt, ether, a)
	outgoingPacket := buffer.Bytes()

	handle, err := pcap.OpenLive(pcapName, 2048, false, 30*time.Second)
	if err != nil {
		log.Fatal("pcap打开失败:", err)
	}
	defer handle.Close()

	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		log.Fatal("发送arp数据包失败..")
	}

	log.Info("Send to: ", dstIP)
}
