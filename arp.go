package main

import (
	"context"
	"net"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	manuf "github.com/timest/gomanuf"
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
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
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
		DstHwAddress:      net.HardwareAddr{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
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

	log.Info("ARP Send to: ", dstIP)
	t = time.NewTicker(3 * time.Second)
}

func listenARP(ctx context.Context) {
	handle, err := pcap.OpenLive(pcapName, 1024, false, 20*time.Second)
	if err != nil {
		log.Fatal("ARP接受 Pcap打开失败: ", err)
	}
	defer handle.Close()

	// 设置BPF过滤器，并检查是否有错误
	// if err := handle.SetBPFFilter("arp"); err != nil {
	// 	log.Fatal("设置BPF过滤器失败: ", err)
	// }
	// handle.SetBPFFilter("arp or rarp")

	ps := gopacket.NewPacketSource(handle, handle.LinkType())

	log.Info("启动ARP监听")
	for {
		select {
		//外部调用 终止监听
		case <-ctx.Done():
			log.Info("ARP listen 结束")
			return
		case p := <-ps.Packets():
			if arpLayer := p.Layer(layers.LayerTypeARP); arpLayer != nil {
				arp, _ := arpLayer.(*layers.ARP)
				if arp.Operation == 2 {
					mac := net.HardwareAddr(arp.SourceHwAddress)
					ip := ParseIP(arp.SourceProtAddress)
					m := manuf.Search(mac.String())
					// log.Info("IP: MAC:", ip, mac, m)
					pushData(ip, mac, "", m)
				}
			}
			// arp := p.Layer(layers.LayerTypeARP).(*layers.ARP)
		}
	}
}
